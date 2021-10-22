package audit_mock

import (
	"testing"
	"time"

	"gitlab.com/inetmock/inetmock/pkg/audit"
)

type EmitterMockEmitCallParams struct {
	Ev audit.Event
}

type EmitterMockEmitCall struct {
	Timestamp time.Time
	Params    EmitterMockEmitCallParams
}

type EmitterMockCalls struct {
	Emit []EmitterMockEmitCall
}

type EmitterMockCallsContext struct {
	EmitterMockCalls
	TB testing.TB
}

type EmitterMock struct {
	TB     testing.TB
	Calls  EmitterMockCalls
	OnEmit func(state EmitterMockCallsContext, ev audit.Event)
}

func (em *EmitterMock) Emit(ev audit.Event) {
	em.Calls.Emit = append(em.Calls.Emit, EmitterMockEmitCall{
		Timestamp: time.Now(),
		Params: EmitterMockEmitCallParams{
			Ev: ev,
		},
	})
	if em.OnEmit != nil {
		ctx := EmitterMockCallsContext{
			EmitterMockCalls: em.Calls,
			TB:               em.TB,
		}
		em.OnEmit(ctx, ev)
	}
}
