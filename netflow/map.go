package netflow

import (
	"encoding"
	"errors"
	"fmt"
	"sync"

	manager "github.com/DataDog/ebpf-manager"
	"github.com/cilium/ebpf"
	"golang.org/x/exp/maps"
)

type BinarySizer interface {
	BinarySize() int
}

type BinaryCollectionUnmarshaler interface {
	BinarySizer
	encoding.BinaryUnmarshaler
}

type MapKey interface {
	comparable
}

func MapOf[K MapKey, V any](m *ebpf.Map) *Map[K, V] {
	return &Map[K, V]{
		underlying: m,
	}
}

func WithBatchSize(size int) BatchOption {
	return BatchOptionFunc(func(opt *batchOptions) {
		opt.BatchSize = size
	})
}

func WithUseFallback(useFallback bool) BatchOption {
	return BatchOptionFunc(func(opt *batchOptions) {
		opt.UseFallback = useFallback
	})
}

type batchOptions struct {
	BatchSize   int
	UseFallback bool
}

func (o batchOptions) Fallback(feat EBPFFeature) bool {
	return o.UseFallback || !CheckForFeature(feat)
}

type (
	BatchOption interface {
		Apply(opt *batchOptions)
	}
	BatchOptionFunc func(opt *batchOptions)
)

func (f BatchOptionFunc) Apply(opt *batchOptions) {
	f(opt)
}

type Map[K MapKey, V any] struct {
	lock       sync.RWMutex
	underlying *ebpf.Map
}

func (m *Map[K, V]) Iterate(accept func(key K, val V)) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.doIterate(accept)
}

func (m *Map[K, V]) Get(key K) (val V, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if err := m.underlying.Lookup(key, &val); err != nil {
		return val, err
	}

	return val, nil
}

func (m *Map[K, V]) Put(key K, val V) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.underlying.Put(key, val)
}

func (m *Map[K, V]) PutAll(toInsert map[K]V, opts ...BatchOption) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.doPutAll(toInsert, opts)
}

func (m *Map[K, V]) GetAll(opts ...BatchOption) (result map[K]V, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.doGetAll(opts)
}

func (m *Map[K, V]) Sync(desired map[K]V, opts ...BatchOption) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	var (
		current  map[K]V
		toInsert = maps.Clone(desired)
		toDelete = make([]K, 0)
	)

	if current, err = m.doGetAll(opts); err != nil {
		return err
	}

	for k := range current {
		if _, stillPresent := desired[k]; stillPresent {
			delete(toInsert, k)
		} else {
			toDelete = append(toDelete, k)
		}
	}

	return errors.Join(
		m.doDeleteAll(toDelete, opts),
		m.doPutAll(toInsert, opts),
	)
}

func (m *Map[K, V]) DeleteAll(keys []K, opts ...BatchOption) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.doDeleteAll(keys, opts)
}

func (m *Map[K, V]) Cap() uint32 {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.underlying.MaxEntries()
}

func (m *Map[K, V]) doGetAll(opts []BatchOption) (result map[K]V, err error) {
	bo := batchOptionsFor(opts)
	if bo.Fallback(FeatureBatchAPI) {
		return m.getAllFallback()
	}

	return m.getAllUsingBatch(bo)
}

func (m *Map[K, V]) doPutAll(toInsert map[K]V, opts []BatchOption) error {
	bo := batchOptionsFor(opts)
	if bo.Fallback(FeatureBatchAPI) {
		return m.putAllFallback(toInsert)
	}

	return m.putAllUsingBatch(toInsert)
}

func (m *Map[K, V]) doDeleteAll(keys []K, opts []BatchOption) error {
	bo := batchOptionsFor(opts)
	if bo.Fallback(FeatureBatchAPI) {
		return m.deleteAllFallback(keys)
	}

	return m.deleteAllUsingBatch(keys)
}

func (m *Map[K, V]) doIterate(accept func(key K, val V)) error {
	iterator := m.underlying.Iterate()
	var (
		key = new(K)
		val = new(V)
	)
	for iterator.Next(key, val) {
		accept(*key, *val)
	}

	if err := iterator.Err(); err != nil {
		return err
	}
	return nil
}

func (m *Map[K, V]) getAllUsingBatch(bo batchOptions) (result map[K]V, err error) {
	var (
		prevKey    *K
		nextKeyOut = new(K)
		keysOut    = initiatedColOf[K](bo.BatchSize)
		valuesOut  = initiatedColOf[V](bo.BatchSize)
		elemsRead  int
	)

	result = make(map[K]V)

	for err == nil {
		if prevKey == nil {
			elemsRead, err = m.underlying.BatchLookup(nil, nextKeyOut, keysOut, valuesOut, nil)
		} else {
			elemsRead, err = m.underlying.BatchLookup(prevKey, nextKeyOut, keysOut, valuesOut, nil)
		}

		for i := 0; i < elemsRead; i++ {
			result[keysOut[i]] = valuesOut[i]
		}

		if err != nil && !errors.Is(err, ebpf.ErrKeyNotExist) {
			return nil, err
		}

		prevKey = nextKeyOut
	}

	return result, nil
}

func (m *Map[K, V]) putAllUsingBatch(toInsert map[K]V) (err error) {
	var (
		keys = make(valueCollection[K], 0, len(toInsert))
		vals = make(valueCollection[V], 0, len(toInsert))
	)

	for k, v := range toInsert {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	_, err = m.underlying.BatchUpdate(keys, vals, nil)
	return err
}

func (m *Map[K, V]) putAllFallback(toInsert map[K]V) error {
	for k, v := range toInsert {
		if err := m.underlying.Put(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (m *Map[K, V]) getAllFallback() (result map[K]V, err error) {
	result = make(map[K]V)

	err = m.doIterate(func(key K, val V) {
		result[key] = val
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Map[K, V]) deleteAllUsingBatch(keys []K) (err error) {
	vc := valueCollection[K](keys)
	_, err = m.underlying.BatchDelete(vc, nil)
	return err
}

func (m *Map[K, V]) deleteAllFallback(keys []K) (err error) {
	for i := range keys {
		if err := m.underlying.Delete(keys[i]); err != nil {
			return err
		}
	}
	return nil
}

func initiatedColOf[T any](length int) valueCollection[T] {
	result := make(valueCollection[T], length)
	for i := 0; i < length; i++ {
		t := new(T)
		result[i] = *t
	}

	return result
}

func batchOptionsFor(opts []BatchOption) batchOptions {
	const (
		defaultBatchSize = 50
	)

	inst := batchOptions{
		BatchSize: defaultBatchSize,
	}

	for i := range opts {
		opts[i].Apply(&inst)
	}

	return inst
}

func getMapFromManager(mgr *manager.Manager, mapName string) (*ebpf.Map, error) {
	if m, present, err := mgr.GetMap(mapName); err != nil {
		return nil, err
	} else if !present {
		return nil, fmt.Errorf("%w: %s", ErrMissingMap, mapName)
	} else {
		return m, nil
	}
}

func mapOfManager[K MapKey, V any](mgr *manager.Manager, name string) (*Map[K, V], error) {
	if m, err := getMapFromManager(mgr, name); err != nil {
		return nil, err
	} else {
		return MapOf[K, V](m), nil
	}
}
