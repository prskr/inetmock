package queue

import (
	"math"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	stretchFactor      float64 = 0.2
	minimumCapIncrease float64 = 5.0
	minimumCapReserve  float64 = 5.0
	minimumCapacity    int     = 10
	evictionCacheSized int     = 10
)

var (
	cacheEvictionCounter prometheus.Counter
)

func init() {
	cacheEvictionCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "dns",
		Subsystem: "cache",
		Name:      "evictions_total",
		Help:      "Number of times how often the cache got evicted",
	})
	prometheus.MustRegister(cacheEvictionCounter)
}

type EvictionCallback interface {
	OnEvicted(evictedEntries []*Entry)
}

type EvictionCallbackFunc func(evictedEntries []*Entry)

func (f EvictionCallbackFunc) OnEvicted(evictedEntries []*Entry) {
	f(evictedEntries)
}

type TTL struct {
	modLock       sync.Locker
	readLock      sync.Locker
	offset        Offset
	backing       []*Entry
	virtual       []*Entry
	evictionCache chan []*Entry
}

func NewTTLFromSeed(seed []*Entry) *TTL {
	var mutex = new(sync.RWMutex)

	for idx := range seed {
		seed[idx].index = idx
	}

	return &TTL{
		modLock:  mutex,
		readLock: mutex.RLocker(),
		backing:  seed,
		virtual:  seed,
	}
}

func NewTTL(capacity int) *TTL {
	if capacity < minimumCapacity {
		capacity = minimumCapacity
	}
	return NewTTLFromSeed(make([]*Entry, 0, capacity))
}

func (t *TTL) Len() int {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	return len(t.virtual)
}

func (t *TTL) Cap() int {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	return cap(t.virtual)
}

func (t *TTL) Get(idx int) *Entry {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	if idx >= len(t.virtual) {
		return nil
	}
	return t.virtual[idx]
}

func (t *TTL) IndexOf(e *Entry) int {
	return e.index - t.offset.CurrentOffset
}

func (t *TTL) UpdateTTL(e *Entry, newTTL time.Duration) {
	t.modLock.Lock()
	defer t.modLock.Unlock()

	var (
		length    = len(t.virtual)
		insertIdx int
	)

	e.timeout = time.Now().UTC().Add(newTTL)
	insertIdx = sort.Search(length, func(i int) bool {
		return t.virtual[i].timeout.After(e.timeout)
	})

	if insertIdx == length {
		insertIdx -= 1
	}

	if insertIdx == e.index-t.offset.CurrentOffset {
		return
	}

	t.move(e.index-t.offset.CurrentOffset, insertIdx, 1)

	t.virtual[insertIdx] = e
	e.index = insertIdx + t.offset.CurrentOffset
}

func (t *TTL) Evict() {
	t.modLock.Lock()
	defer t.modLock.Unlock()

	t.doEvict()
}

func (t *TTL) Push(name string, value interface{}, ttl time.Duration) *Entry {
	t.modLock.Lock()
	defer t.modLock.Unlock()

	var (
		length   = len(t.virtual)
		capacity = cap(t.virtual)
		entry    = &Entry{
			Key:     name,
			Value:   value,
			timeout: time.Now().UTC().Add(ttl),
			index:   length + t.offset.CurrentOffset,
		}
	)

	if length == capacity {
		t.doEvict()
		length = len(t.virtual)
		entry.index = length + t.offset.CurrentOffset
	}

	/*
	 * Shortcut if the TTL is already the latest one
	 */
	if length == 0 || entry.timeout.After(t.virtual[length-1].timeout) {
		t.virtual = append(t.virtual, entry)
		return entry
	}

	var insertIdx = sort.Search(length, func(i int) bool {
		return t.virtual[i].timeout.After(entry.timeout)
	})

	if insertIdx >= length {
		t.virtual = append(t.virtual, entry)
	} else {
		t.virtual = append(t.virtual, nil)
		t.move(length, insertIdx, -1)
		t.virtual[insertIdx] = entry
		entry.index = insertIdx + t.offset.CurrentOffset
	}

	return entry
}

func (t *TTL) PeekFront() *Entry {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	if len(t.virtual) > 0 {
		return t.virtual[0]
	}
	return nil
}

func (t *TTL) OnEvicted(callback EvictionCallback) {
	t.modLock.Lock()
	defer t.modLock.Unlock()
	if t.evictionCache != nil {
		close(t.evictionCache)
	}

	t.evictionCache = make(chan []*Entry, evictionCacheSized)
	go func(in <-chan []*Entry, callback EvictionCallback) {
		for e := range in {
			callback.OnEvicted(e)
		}
	}(t.evictionCache, callback)
}

func (t *TTL) doEvict() {
	var (
		virtualLength   = len(t.virtual)
		now             = time.Now().UTC()
		firstToNotEvict int
	)
	if virtualLength < 1 {
		return
	}

	for firstToNotEvict = 0; firstToNotEvict < virtualLength; firstToNotEvict++ {
		if t.virtual[firstToNotEvict].timeout.After(now) {
			break
		}
	}

	if firstToNotEvict == 0 {
		return
	}

	if t.evictionCache != nil {
		go func(evictedItems []*Entry) {
			t.evictionCache <- evictedItems
		}(t.virtual[:firstToNotEvict])
	}

	t.virtual = t.virtual[firstToNotEvict:]

	var newOffset, overflow = t.offset.Inc(firstToNotEvict)

	var (
		minimumReserve = int(math.Max(math.Ceil((1.0-stretchFactor)*float64(len(t.virtual))), minimumCapReserve))
		virtualReserve = cap(t.virtual) - len(t.virtual)
		backingReserve = cap(t.backing) - cap(t.virtual)
	)

	switch {
	// no need to copy elements
	case virtualReserve > minimumReserve:
		return

	/*
	 * re-acquire reserve from backing field by copying from virtual back to
	 * the beginning of the backing field, resetting virtual to original size and
	 * resetting the backing field
	 */
	case virtualReserve+backingReserve > minimumReserve:
		cacheEvictionCounter.Inc()
		t.backing = append(t.backing, t.virtual...)
		t.virtual = t.backing
		t.backing = t.backing[0:0]
	/*
	 * if neither the actual/virtual reserve nor the backing reserve are enough to
	 * get a buffer of 20% back increase the underlying structures and then do the same as above
	 */
	default:
		cacheEvictionCounter.Inc()
		t.backing = make([]*Entry, 0, cap(t.backing)+int(math.Max(math.Ceil(float64(cap(t.backing))*stretchFactor), minimumCapIncrease)))
		t.backing = append(t.backing, t.virtual...)
		t.virtual = t.backing
		t.backing = t.backing[0:0]
	}
	/*
	 * if the offset variable overflowed it is necessary to reset all cached indices to ensure the 'touch' functionality still works
	 */
	if overflow {
		virtualLength = len(t.virtual)
		for idx := 0; idx < virtualLength; idx++ {
			t.virtual[idx].index = idx + newOffset
		}
	}
}

func (t *TTL) move(startIdx, endIdx, offset int) {
	length := len(t.virtual)
	for idx := startIdx; idx != endIdx && idx < length; idx += offset {
		t.virtual[idx] = t.virtual[idx+offset]
		if t.virtual[idx] != nil {
			t.virtual[idx].index = idx + t.offset.CurrentOffset
		}
	}
}
