package netflow

import (
	"errors"
	"sync"
	"time"
)

var ErrAlreadySyncing = errors.New("already syncing epoch")

const DefaultEpochSyncWindow = 1 * time.Second

type (
	epochOptions struct {
		Start        time.Time
		ErrorHandler ErrorSink
	}

	EpochOption interface {
		Apply(opt *epochOptions)
	}

	EpochOptionFunc func(opt *epochOptions)
)

func (f EpochOptionFunc) Apply(opt *epochOptions) {
	f(opt)
}

func WithStart(startTime time.Time) EpochOption {
	return EpochOptionFunc(func(opt *epochOptions) {
		opt.Start = startTime
	})
}

func NewEpoch(m *Map[NATConfigKey, uint32], opts ...EpochOption) *Epoch {
	o := epochOptions{
		Start:        time.Now(),
		ErrorHandler: noOpErrorSink,
	}

	for i := range opts {
		opts[i].Apply(&o)
	}

	return &Epoch{
		ConfigMap:    m,
		Start:        o.Start,
		ErrorHandler: o.ErrorHandler,
	}
}

type Epoch struct {
	done         chan struct{}
	lock         sync.Mutex
	ticker       *time.Ticker
	Start        time.Time
	ConfigMap    *Map[NATConfigKey, uint32]
	ErrorHandler ErrorSink
}

func (e *Epoch) StartSync(windowSize time.Duration) error {
	e.lock.Lock()
	if e.ticker != nil {
		return ErrAlreadySyncing
	}
	e.done = make(chan struct{})
	e.ticker = time.NewTicker(windowSize)
	e.lock.Unlock()

	go e.doSync()

	return nil
}

func (e *Epoch) Sync() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.isDone() {
		return nil
	}

	currentEpochValue := uint32(time.Since(e.Start).Seconds())
	return e.ConfigMap.Put(natConfigKeyCurrentEpoch, currentEpochValue)
}

func (e *Epoch) Stop() {
	e.lock.Lock()
	defer e.lock.Unlock()

	close(e.done)
	e.ticker.Stop()
}

func (e *Epoch) doSync() {
	for range e.ticker.C {
		if err := e.Sync(); err != nil {
			if errHandler := e.ErrorHandler; errHandler != nil {
				errHandler.OnError(err)
			}
		}
	}
}

func (e *Epoch) isDone() bool {
	select {
	case _, more := <-e.done:
		return !more
	default:
		return false
	}
}
