package dns

import (
	"math"
	"net"
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

type TTLQueue interface {
	Push(name string, address net.IP, ttl time.Duration) *Entry
	UpdateTTL(e *Entry, newTTL time.Duration)
	Evict()
	PeekFront() *Entry
	OnEvicted(callback EvictionCallback)
	Cap() int
	Len() int
}

type ttlQueue struct {
	modLock       sync.Locker
	readLock      sync.Locker
	backing       []*Entry
	virtual       []*Entry
	evictionCache chan []*Entry
}

func NewFromSeed(seed []*Entry) TTLQueue {
	var mutex = new(sync.RWMutex)
	return &ttlQueue{
		modLock:  mutex,
		readLock: mutex.RLocker(),
		backing:  seed,
		virtual:  seed,
	}
}

func NewQueue(capacity int) TTLQueue {
	if capacity < minimumCapacity {
		capacity = minimumCapacity
	}
	return NewFromSeed(make([]*Entry, 0, capacity))
}

func (t *ttlQueue) Len() int {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	return len(t.virtual)
}

func (t *ttlQueue) Cap() int {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	return cap(t.virtual)
}

func (t *ttlQueue) Get(idx int) *Entry {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	if idx >= len(t.virtual) {
		return nil
	}
	return t.virtual[idx]
}

func (t *ttlQueue) Sort() {
	t.modLock.Lock()
	defer t.modLock.Unlock()
	sort.Slice(t.virtual, func(i, j int) bool {
		return t.virtual[i].timeout.Before(t.virtual[j].timeout)
	})
}

func (t *ttlQueue) UpdateTTL(e *Entry, newTTL time.Duration) {
	e.timeout = time.Now().UTC().Add(newTTL)
	t.Sort()
}

func (t *ttlQueue) Evict() {
	t.modLock.Lock()
	defer t.modLock.Unlock()

	t.doEvict()
}

func (t *ttlQueue) Push(name string, address net.IP, ttl time.Duration) *Entry {
	t.modLock.Lock()
	defer t.Sort()
	defer t.modLock.Unlock()

	entry := &Entry{
		Name:    name,
		Address: address,
		timeout: time.Now().UTC().Add(ttl),
	}
	if len(t.virtual) == cap(t.virtual) {
		t.doEvict()
	}
	t.virtual = append(t.virtual, entry)

	return entry
}

func (t *ttlQueue) PeekFront() *Entry {
	t.readLock.Lock()
	defer t.readLock.Unlock()
	if len(t.virtual) > 0 {
		return t.virtual[0]
	}
	return nil
}

func (t *ttlQueue) OnEvicted(callback EvictionCallback) {
	t.modLock.Lock()
	defer t.modLock.Unlock()
	if t.evictionCache != nil {
		close(t.evictionCache)
	}

	t.evictionCache = make(chan []*Entry, 10)
	go func(in <-chan []*Entry, callback EvictionCallback) {
		for e := range in {
			callback.OnEvicted(e)
		}
	}(t.evictionCache, callback)
}

func (t *ttlQueue) doEvict() {
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
		t.evictionCache <- t.virtual[:firstToNotEvict]
	}
	t.virtual = t.virtual[firstToNotEvict:]

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
		return
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
		return
	}
}
