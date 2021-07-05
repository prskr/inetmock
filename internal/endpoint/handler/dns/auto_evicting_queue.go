package dns

import (
	"net"
	"time"
)

const (
	defaultTimerDuration = 500 * time.Millisecond
)

type autoEvictingQueue struct {
	backing TTLQueue
	timer   *time.Timer
}

func WrapToAutoEvict(existing TTLQueue) TTLQueue {
	queue := &autoEvictingQueue{
		timer:   time.NewTimer(defaultTimerDuration),
		backing: existing,
	}

	queue.startEvictionTimer()

	return queue
}

func (a *autoEvictingQueue) Push(name string, address net.IP, ttl time.Duration) *Entry {
	return a.backing.Push(name, address, ttl)
}

func (a *autoEvictingQueue) UpdateTTL(entry *Entry, newTTL time.Duration) {
	a.backing.UpdateTTL(entry, newTTL)
}

func (a *autoEvictingQueue) Evict() {
	a.backing.Evict()
}

func (a *autoEvictingQueue) PeekFront() *Entry {
	return a.backing.PeekFront()
}

func (a autoEvictingQueue) OnEvicted(callback EvictionCallback) {
	a.backing.OnEvicted(callback)
}

func (a autoEvictingQueue) Cap() int {
	return a.backing.Cap()
}

func (a autoEvictingQueue) Len() int {
	return a.backing.Len()
}

func (a *autoEvictingQueue) startEvictionTimer() {
	go func() {
		<-time.After(50 * time.Millisecond)
		for {
			<-a.timer.C
			a.backing.Evict()
			if front := a.backing.PeekFront(); front == nil {
				a.timer.Reset(defaultTimerDuration)
			} else if front.timeout.After(time.Now().UTC()) {
				a.timer.Reset(front.timeout.Sub(time.Now().UTC()).Round(50 * time.Millisecond))
			} else {
				a.timer.Reset(defaultTimerDuration)
			}
		}
	}()
}
