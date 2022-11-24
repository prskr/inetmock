package audit_mock

import (
	"sync"
	"testing"
	"time"

	"inetmock.icb4dc0.de/inetmock/pkg/audit"
)

var _ audit.Emitter = (*EmitterMock)(nil)

type EmitterMockEmitCallParams struct {
	Ev *audit.Event
}

type EmitterMockEmitCall struct {
	Timestamp time.Time
	Params    EmitterMockEmitCallParams
}

type EmitterMockCalls struct {
	lock sync.RWMutex
	emit []EmitterMockEmitCall
}

func (c *EmitterMockCalls) AddEmit(emitCall EmitterMockEmitCall) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.emit = append(c.emit, emitCall)
}

func (c *EmitterMockCalls) Emit() []EmitterMockEmitCall {
	c.lock.RLock()
	defer c.lock.RUnlock()

	result := make([]EmitterMockEmitCall, len(c.emit))
	copy(result, c.emit)
	return result
}

type EmitterMockCallsContext struct {
	*EmitterMockCalls
	TB testing.TB
}

type EmitterMock struct {
	lock   sync.Mutex
	calls  EmitterMockCalls
	TB     testing.TB
	OnEmit func(state *EmitterMockCallsContext, ev *audit.Event)
}

func (em *EmitterMock) WithCalls(f func(calls *EmitterMockCalls)) {
	em.lock.Lock()
	defer em.lock.Unlock()

	f(&em.calls)
}

func (em *EmitterMock) Emit(ev *audit.Event) {
	em.lock.Lock()
	defer em.lock.Unlock()

	em.calls.AddEmit(EmitterMockEmitCall{
		Timestamp: time.Now(),
		Params: EmitterMockEmitCallParams{
			Ev: ev,
		},
	})

	if em.OnEmit != nil {
		ctx := &EmitterMockCallsContext{
			EmitterMockCalls: &em.calls,
			TB:               em.TB,
		}
		em.OnEmit(ctx, ev)
	}
}

func (em *EmitterMock) Builder() audit.EventBuilder {
	return audit.BuilderForEmitter(em)
}
