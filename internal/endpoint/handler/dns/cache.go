package dns

import (
	"net"
	"sync"
	"time"

	"gitlab.com/inetmock/inetmock/internal/queue"
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

type Record struct {
	Name    string
	Address net.IP
}

type IPResolver interface {
	Lookup(host string) net.IP
}

type IPResolverFunc func(host string) net.IP

func (f IPResolverFunc) Lookup(host string) net.IP {
	return f(host)
}

type CacheOption func(cfg *cacheConfig)

func NewCache(opts ...CacheOption) *Cache {
	var cfg = cacheConfig{
		ttl:         defaultTTL,
		initialSize: defaultInitialSize,
	}
	for idx := range opts {
		opts[idx](&cfg)
	}

	var rwMutex = new(sync.RWMutex)

	var cache = &Cache{
		cfg:          cfg,
		readLock:     rwMutex.RLocker(),
		writeLock:    rwMutex,
		forwardIndex: make(map[string]*queue.Entry),
		reverseIndex: make(map[uint32]*queue.Entry),
		queue:        queue.WrapToAutoEvict(queue.NewTTL(cfg.initialSize)),
	}

	cache.queue.OnEvicted(queue.EvictionCallbackFunc(cache.onCacheEvicted))

	return cache
}

type cacheConfig struct {
	ttl         time.Duration
	initialSize int
}

type cacheQueue interface {
	Push(name string, value interface{}, ttl time.Duration) *queue.Entry
	UpdateTTL(e *queue.Entry, newTTL time.Duration)
	OnEvicted(callback queue.EvictionCallback)
}

type Cache struct {
	cfg          cacheConfig
	readLock     sync.Locker
	writeLock    sync.Locker
	forwardIndex map[string]*queue.Entry
	reverseIndex map[uint32]*queue.Entry
	queue        cacheQueue
}

func (c *Cache) PutRecord(host string, address net.IP) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	rec := &Record{
		Name:    host,
		Address: address,
	}
	i := IPToInt32(address)
	e := c.queue.Push(host, rec, c.cfg.ttl)
	c.forwardIndex[host] = e
	c.reverseIndex[i] = e
}

func (c *Cache) ForwardLookup(host string) net.IP {
	c.readLock.Lock()
	if e, cached := c.forwardIndex[host]; cached {
		c.queue.UpdateTTL(e, c.cfg.ttl)
		c.readLock.Unlock()
		return e.Value.(*Record).Address
	} else {
		c.readLock.Unlock()
		return nil
	}
}

func (c *Cache) ReverseLookup(address net.IP) (host string, miss bool) {
	c.readLock.Lock()
	defer c.readLock.Unlock()
	if e, cached := c.reverseIndex[IPToInt32(address)]; cached {
		c.queue.UpdateTTL(e, c.cfg.ttl)
		return e.Value.(*Record).Name, false
	} else {
		return "", true
	}
}

func (c *Cache) onCacheEvicted(evictedItems []*queue.Entry) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	for idx := range evictedItems {
		var record = evictedItems[idx].Value.(*Record)
		delete(c.forwardIndex, record.Name)
		delete(c.reverseIndex, IPToInt32(record.Address))
	}
}
