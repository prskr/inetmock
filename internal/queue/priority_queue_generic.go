// Package queue - bla
//nolint // Go 1.18 is not yet supported

//go:build go1.18
// +build go1.18

package queue

import (
	"math"
	"sort"
	"sync"
	"time"
)

type GEvictionCallback[T any] interface {
	OnEvicted(evictedEntries []*GEntry[T])
}

func NewGTTLFromSeed[T any](seed []*GEntry[T]) *GTTL[T] {
	mutex := new(sync.RWMutex)

	for idx := range seed {
		seed[idx].index = idx
	}

	return &GTTL[T]{
		modLock:  mutex,
		readLock: mutex.RLocker(),
		backing:  seed,
		virtual:  seed,
	}
}

func NewGTTL[T any](capacity int) *GTTL[T] {
	if capacity < minimumCapacity {
		capacity = minimumCapacity
	}
	return NewGTTLFromSeed[T](make([]*GEntry[T], 0, capacity))
}

type GEntry[TValue any] struct {
	Key     string
	Value   TValue
	timeout time.Time
	index   int
}

func (e GEntry[TValue]) WithTTL(ttl time.Duration) *GEntry[TValue] {
	e.timeout = time.Now().UTC().Add(ttl)
	return &e
}

func (e GEntry[TValue]) TTL() time.Time {
	return e.timeout
}

type GTTL[T any] struct {
	modLock       sync.Locker
	readLock      sync.Locker
	offset        Offset
	backing       []*GEntry[T]
	virtual       []*GEntry[T]
	evictionCache chan []*GEntry[T]
}

func (t *GTTL[T]) Len() int {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	return len(t.virtual)
}

func (t *GTTL[T]) Cap() int {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	return cap(t.virtual)
}

func (t *GTTL[T]) Get(idx int) *GEntry[T] {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	if idx >= len(t.virtual) {
		return nil
	}
	return t.virtual[idx]
}

func (t *GTTL[T]) IndexOf(e *GEntry[T]) int {
	return e.index - t.offset.CurrentOffset
}

func (t *GTTL[T]) UpdateTTL(e *GEntry[T], newTTL time.Duration) {
	t.modLock.Lock()
	defer t.modLock.Unlock()

	if e.index == evictedEntryIndex {
		newEntry := t.Push(e.Key, e.Value, newTTL)
		*e = *newEntry
		return
	}

	var (
		length    = len(t.virtual)
		insertIdx int
	)

	e.timeout = time.Now().UTC().Add(newTTL)
	insertIdx = sort.Search(length, func(i int) bool {
		return t.virtual[i].timeout.After(e.timeout)
	})

	if insertIdx >= length {
		insertIdx = length - 1
	}

	// if actual index is desired index
	// e.index is relative index to the current offset
	if insertIdx == e.index-t.offset.CurrentOffset {
		return
	}

	t.move(e.index-t.offset.CurrentOffset, insertIdx, 1)

	t.virtual[insertIdx] = e
	e.index = insertIdx + t.offset.CurrentOffset
}

func (t *GTTL[T]) Evict() {
	t.modLock.Lock()
	defer t.modLock.Unlock()

	t.doEvict()
}

func (t *GTTL[T]) Push(name string, value T, ttl time.Duration) *GEntry[T] {
	t.modLock.Lock()
	defer t.modLock.Unlock()

	var (
		length   = len(t.virtual)
		capacity = cap(t.virtual)
		entry    = &GEntry[T]{
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

	insertIdx := sort.Search(length, func(i int) bool {
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

func (t *GTTL[T]) PeekFront() *GEntry[T] {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	if len(t.virtual) > 0 {
		return t.virtual[0]
	}
	return nil
}

func (t *GTTL[T]) OnEvicted(callback GEvictionCallback[T]) {
	t.modLock.Lock()
	defer t.modLock.Unlock()
	if t.evictionCache != nil {
		close(t.evictionCache)
	}

	t.evictionCache = make(chan []*GEntry[T], evictionCacheSized)
	go func(in <-chan []*GEntry[T], callback GEvictionCallback[T]) {
		for e := range in {
			callback.OnEvicted(e)
		}
	}(t.evictionCache, callback)
}

func (t *GTTL[T]) doEvict() {
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

	evictedEntries := t.virtual[:firstToNotEvict]

	for idx := range evictedEntries {
		evictedEntries[idx].index = evictedEntryIndex
	}

	if t.evictionCache != nil {
		go func(evictedItems []*GEntry[T]) {
			t.evictionCache <- evictedItems
		}(evictedEntries)
	}

	t.virtual = t.virtual[firstToNotEvict:]

	newOffset, overflow := t.offset.Inc(firstToNotEvict)

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
		t.backing = make([]*GEntry[T], 0, cap(t.backing)+int(math.Max(math.Ceil(float64(cap(t.backing))*stretchFactor), minimumCapIncrease)))
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

func (t *GTTL[T]) move(startIdx, endIdx, offset int) {
	if endIdx < 0 {
		return
	}

	if startIdx < 0 {
		startIdx = 0
	}

	length := len(t.virtual)
	for idx := startIdx; idx != endIdx && idx < length; idx += offset {
		t.virtual[idx] = t.virtual[idx+offset]
		if t.virtual[idx] != nil {
			t.virtual[idx].index = idx + t.offset.CurrentOffset
		}
	}
}
