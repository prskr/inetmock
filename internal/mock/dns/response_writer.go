package dnsmock

import (
	"net"

	mdns "github.com/miekg/dns"

	"inetmock.icb4dc0.de/inetmock/internal/mock"
)

type ResponseWriterMock struct {
	*mock.CallCounter
	Local        net.Addr
	Remote       net.Addr
	OnWriteMsg   func(msg *mdns.Msg) error
	OnWrite      func(data []byte) (int, error)
	OnClose      func() error
	OnTsigStatus func() error
	OnHijack     func()
}

func (rw ResponseWriterMock) LocalAddr() net.Addr {
	if rw.CallCounter != nil {
		rw.Inc(rw.LocalAddr)
	}
	return rw.Local
}

func (rw ResponseWriterMock) RemoteAddr() net.Addr {
	if rw.CallCounter != nil {
		rw.Inc(rw.RemoteAddr)
	}
	return rw.Remote
}

func (rw ResponseWriterMock) WriteMsg(msg *mdns.Msg) error {
	if rw.CallCounter != nil {
		rw.Inc(rw.WriteMsg)
	}
	if rw.OnWriteMsg != nil {
		return rw.OnWriteMsg(msg)
	}
	return nil
}

func (rw ResponseWriterMock) Write(data []byte) (int, error) {
	if rw.CallCounter != nil {
		rw.Inc(rw.Write)
	}
	if rw.OnWrite != nil {
		return rw.OnWrite(data)
	}
	return len(data), nil
}

func (rw ResponseWriterMock) Close() error {
	if rw.CallCounter != nil {
		rw.Inc(rw.Close)
	}
	if rw.OnClose != nil {
		return rw.OnClose()
	}
	return nil
}

func (rw ResponseWriterMock) TsigStatus() error {
	if rw.CallCounter != nil {
		rw.Inc(rw.TsigStatus)
	}
	if rw.OnTsigStatus != nil {
		return rw.OnTsigStatus()
	}
	return nil
}

func (rw ResponseWriterMock) TsigTimersOnly(bool) {
	if rw.CallCounter != nil {
		rw.Inc(rw.TsigTimersOnly)
	}
}

func (rw ResponseWriterMock) Hijack() {
	if rw.CallCounter != nil {
		rw.Inc(rw.Hijack)
	}
	if rw.OnHijack != nil {
		rw.OnHijack()
	}
}
