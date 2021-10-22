package dnsmock

import (
	"net"
	"testing"
	"time"
)

type CacheMockPutRecordCallParams struct {
	Host    string
	Address net.IP
}

type CacheMockForwardLookupCallParams struct {
	Host string
}

type CacheMockForwardLookupCallResults struct {
	Res0 net.IP
}

type CacheMockReverseLookupCallParams struct {
	Address net.IP
}

type CacheMockReverseLookupCallResults struct {
	Host string
	Miss bool
}

type CacheMockPutRecordCall struct {
	Timestamp time.Time
	Params    CacheMockPutRecordCallParams
}

type CacheMockForwardLookupCall struct {
	Timestamp time.Time
	Params    CacheMockForwardLookupCallParams
	Results   CacheMockForwardLookupCallResults
}

type CacheMockReverseLookupCall struct {
	Timestamp time.Time
	Params    CacheMockReverseLookupCallParams
	Results   CacheMockReverseLookupCallResults
}

type CacheMockCalls struct {
	PutRecord     []CacheMockPutRecordCall
	ForwardLookup []CacheMockForwardLookupCall
	ReverseLookup []CacheMockReverseLookupCall
}

type CacheMockCallsContext struct {
	CacheMockCalls
	TB testing.TB
}

type CacheMock struct {
	TB              testing.TB
	Calls           CacheMockCalls
	OnPutRecord     func(state CacheMockCallsContext, host string, address net.IP)
	OnForwardLookup func(state CacheMockCallsContext, host string) net.IP
	OnReverseLookup func(state CacheMockCallsContext, address net.IP) (host string, miss bool)
}

func (m *CacheMock) PutRecord(host string, address net.IP) {
	if m.OnPutRecord != nil {
		ctx := CacheMockCallsContext{
			CacheMockCalls: m.Calls,
			TB:             m.TB,
		}
		m.OnPutRecord(ctx, host, address)
	}

	m.Calls.PutRecord = append(m.Calls.PutRecord, CacheMockPutRecordCall{
		Timestamp: time.Now(),
		Params: CacheMockPutRecordCallParams{
			Host:    host,
			Address: address,
		},
	})

	return
}

func (m *CacheMock) ForwardLookup(host string) (res0 net.IP) {
	if m.OnForwardLookup != nil {
		ctx := CacheMockCallsContext{
			CacheMockCalls: m.Calls,
			TB:             m.TB,
		}
		res0 = m.OnForwardLookup(ctx, host)
	}

	m.Calls.ForwardLookup = append(m.Calls.ForwardLookup, CacheMockForwardLookupCall{
		Timestamp: time.Now(),
		Params: CacheMockForwardLookupCallParams{
			Host: host,
		},
		Results: CacheMockForwardLookupCallResults{
			Res0: res0,
		},
	})

	return
}

func (m *CacheMock) ReverseLookup(address net.IP) (host string, miss bool) {
	if m.OnReverseLookup != nil {
		ctx := CacheMockCallsContext{
			CacheMockCalls: m.Calls,
			TB:             m.TB,
		}
		host, miss = m.OnReverseLookup(ctx, address)
	}

	m.Calls.ReverseLookup = append(m.Calls.ReverseLookup, CacheMockReverseLookupCall{
		Timestamp: time.Now(),
		Params: CacheMockReverseLookupCallParams{
			Address: address,
		}, Results: CacheMockReverseLookupCallResults{
			Host: host,
			Miss: miss,
		},
	})

	return
}
