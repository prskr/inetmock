package dns

import (
	"time"
)

const (
	defaultTimerDuration  = 500 * time.Millisecond
	timerDurationRounding = 50 * time.Millisecond
)

type autoEvictingQueue struct {
	TTLQueue
	timer *time.Timer
}

func WrapToAutoEvict(existing TTLQueue) TTLQueue {
	queue := &autoEvictingQueue{
		timer:    time.NewTimer(defaultTimerDuration),
		TTLQueue: existing,
	}

	queue.startEvictionTimer()

	return queue
}

func (a *autoEvictingQueue) startEvictionTimer() {
	go func() {
		for {
			<-a.timer.C
			a.TTLQueue.Evict()
			if front := a.TTLQueue.PeekFront(); front == nil {
				a.timer.Reset(defaultTimerDuration)
			} else if front.timeout.After(time.Now().UTC()) {
				a.timer.Reset(front.timeout.Sub(time.Now().UTC()).Round(timerDurationRounding))
			} else {
				a.timer.Reset(defaultTimerDuration)
			}
		}
	}()
}
