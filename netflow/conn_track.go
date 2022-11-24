package netflow

import (
	"errors"
	"math"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

const (
	connMetaBinaryLength  = 16
	connIdentBinaryLength = 12
)

var ErrCleanupAlreadyRunning = errors.New("conn_track cleanup already running")

func NewConnTrackCleaner(
	m *Map[ConnIdent, ConnMeta],
	errHandler ErrorSink,
	highWaterMark float64,
	interfaceName string,
) *ConnTrackCleaner {
	if errHandler == nil {
		errHandler = noOpErrorSink
	}
	return &ConnTrackCleaner{
		interfaceName: interfaceName,
		connTrackMap:  m,
		ErrorHandler:  errHandler,
		HighWaterMark: highWaterMark,
	}
}

type ConnTrackCleaner struct {
	lock          sync.Mutex
	connTrackMap  *Map[ConnIdent, ConnMeta]
	cleanupTicker *time.Ticker
	interfaceName string
	HighWaterMark float64
	ErrorHandler  ErrorSink
}

func (c *ConnTrackCleaner) Start(interval time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cleanupTicker != nil {
		return ErrCleanupAlreadyRunning
	}

	c.cleanupTicker = time.NewTicker(interval)
	go c.doCleanup()
	return nil
}

func (c *ConnTrackCleaner) Stop() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cleanupTicker.Stop()
}

func (c *ConnTrackCleaner) doCleanup() {
	for range c.cleanupTicker.C {
		c.Cleanup()
	}
}

func (c *ConnTrackCleaner) Cleanup() {
	if !c.lock.TryLock() {
		return
	}
	defer c.lock.Unlock()

	var (
		entriesByTimestamp = make(map[uint32][]ConnIdent)
		timestamps         []uint32
		current            map[ConnIdent]ConnMeta
	)

	if all, err := c.connTrackMap.GetAll(); c.handleErr(err) {
		return
	} else {
		current = all
	}

	for k, v := range current {
		if _, ok := entriesByTimestamp[v.LastObserved]; !ok {
			timestamps = append(timestamps, v.LastObserved)
		}
		entriesByTimestamp[v.LastObserved] = append(entriesByTimestamp[v.LastObserved], k)
	}

	mapCap := float64(c.connTrackMap.Cap())
	currentFillLevel := float64(len(current)) / mapCap
	if currentFillLevel < c.HighWaterMark {
		connTrackGauge.WithLabelValues(c.interfaceName).Set(float64(len(current)))
		return
	}

	minNumberOfEntriesToDelete := int(math.Ceil((currentFillLevel - c.HighWaterMark) * mapCap))
	slices.Sort(timestamps)

	keysToDelete := make([]ConnIdent, 0)
	for i := 0; minNumberOfEntriesToDelete-len(keysToDelete) > 0 && i < len(timestamps); i++ {
		keysToDelete = append(keysToDelete, entriesByTimestamp[timestamps[i]]...)
	}

	if err := c.connTrackMap.DeleteAll(keysToDelete); !c.handleErr(err) {
		connTrackGauge.WithLabelValues(c.interfaceName).Set(float64(len(current) - len(keysToDelete)))
	}
}

func (c *ConnTrackCleaner) handleErr(err error) (errorOccurred bool) {
	if errorOccurred = err != nil; !errorOccurred {
		return
	}

	if errHandler := c.ErrorHandler; errHandler != nil {
		errHandler.OnError(err)
	}
	return errorOccurred
}
