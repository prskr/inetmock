package dns

import (
	"net"
	"sync"
	"time"
)

const (
	defaultTTL         = 30 * time.Second
	defaultInitialSize = 1000
	minimumInitialSize = 100
)

var (
	WithTTL = func(ttl time.Duration) CacheOption {
		return func(cfg *cacheConfig) {
			if ttl > 0 {
				cfg.ttl = ttl
			}
		}
	}
	WithInitialSize = func(initialSize int) CacheOption {
		return func(cfg *cacheConfig) {
			if initialSize >= minimumInitialSize {
				cfg.initialSize = initialSize
			}
		}
	}
)

type Entry struct {
	Name    string
	Address net.IP
	timeout time.Time
}

func (e Entry) WithTTL(ttl time.Duration) *Entry {
	e.timeout = time.Now().UTC().Add(ttl)
	return &e
}

type Cache interface {
	PutRecord(host string, address net.IP)
	ForwardLookup(host string, resolver IPResolver) net.IP
	ReverseLookup(address net.IP) (host string, miss bool)
}

type IPResolver interface {
	Lookup(host string) net.IP
}

type IPResolverFunc func(host string) net.IP

func (f IPResolverFunc) Lookup(host string) net.IP {
	return f(host)
}

type CacheOption func(cfg *cacheConfig)

func NewCache(opts ...CacheOption) Cache {
	var cfg = cacheConfig{
		ttl:         defaultTTL,
		initialSize: defaultInitialSize,
	}
	for idx := range opts {
		opts[idx](&cfg)
	}

	var rwMutex = new(sync.RWMutex)

	var cache = &cache{
		cfg:          cfg,
		readLock:     rwMutex.RLocker(),
		writeLock:    rwMutex,
		forwardIndex: make(map[string]*Entry),
		reverseIndex: make(map[uint32]*Entry),
		queue:        WrapToAutoEvict(NewQueue(cfg.initialSize)),
	}

	cache.queue.OnEvicted(EvictionCallbackFunc(cache.onCacheEvicted))

	return cache
}

type cacheConfig struct {
	ttl         time.Duration
	initialSize int
}

type cache struct {
	cfg          cacheConfig
	readLock     sync.Locker
	writeLock    sync.Locker
	forwardIndex map[string]*Entry
	reverseIndex map[uint32]*Entry
	queue        TTLQueue
}

func (c *cache) PutRecord(host string, address net.IP) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	e := &Entry{
		Name:    host,
		Address: address,
		timeout: time.Now().UTC().Add(c.cfg.ttl),
	}
	i := IPToInt32(address)
	c.forwardIndex[host] = e
	c.reverseIndex[i] = e
}

func (c *cache) ForwardLookup(host string, resolver IPResolver) net.IP {
	c.readLock.Lock()
	if e, cached := c.forwardIndex[host]; cached {
		c.queue.UpdateTTL(e, c.cfg.ttl)
		c.readLock.Unlock()
		return e.Address
	} else {
		ip := resolver.Lookup(host)
		e = c.queue.Push(host, ip, c.cfg.ttl)
		/* need to update the indexes - acquire write-lock */
		c.readLock.Unlock()
		c.writeLock.Lock()
		defer c.writeLock.Unlock()
		c.forwardIndex[host] = e
		c.reverseIndex[IPToInt32(ip)] = e
		return ip
	}
}

func (c *cache) ReverseLookup(address net.IP) (host string, miss bool) {
	c.readLock.Lock()
	defer c.readLock.Unlock()
	if e, cached := c.reverseIndex[IPToInt32(address)]; cached {
		c.queue.UpdateTTL(e, c.cfg.ttl)
		return e.Name, false
	} else {
		return "", true
	}
}

func (c *cache) onCacheEvicted(evictedItems []*Entry) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	for idx := range evictedItems {
		delete(c.forwardIndex, evictedItems[idx].Name)
		delete(c.reverseIndex, IPToInt32(evictedItems[idx].Address))
	}
}
