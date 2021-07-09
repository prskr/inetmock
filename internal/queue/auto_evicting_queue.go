package queue

import (
	"time"
)

const (
	defaultTimerDuration  = 500 * time.Millisecond
	timerDurationRounding = 50 * time.Millisecond
)

type autoEvictingQueue struct {
	TTL
	timer *time.Timer
}

func WrapToAutoEvict(existing TTL) TTL {
	queue := &autoEvictingQueue{
		timer: time.NewTimer(defaultTimerDuration),
		TTL:   existing,
	}

	queue.startEvictionTimer()

	return queue
}

func (a *autoEvictingQueue) startEvictionTimer() {
	go func() {
		for {
			<-a.timer.C
			a.TTL.Evict()
			if front := a.TTL.PeekFront(); front == nil {
				a.timer.Reset(defaultTimerDuration)
			} else if front.timeout.After(time.Now().UTC()) {
				a.timer.Reset(front.timeout.Sub(time.Now().UTC()).Round(timerDurationRounding))
			} else {
				a.timer.Reset(defaultTimerDuration)
			}
		}
	}()
}
