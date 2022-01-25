package endpoint

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/soheilhy/cmux"
	"go.uber.org/multierr"
)

const (
	defaultClosingTimeout = 100 * time.Millisecond
	defaultReadTimeout    = 100 * time.Millisecond
)

type (
	ErrorHandler interface {
		OnError(err error)
	}
	MuxServer interface {
		Serve() error
		Close()
		HandleError(cmux.ErrorHandler)
		SetReadTimeout(time.Duration)
	}
	ErrorHandlerFunc func(err error)
)

func (f ErrorHandlerFunc) OnError(err error) {
	f(err)
}

func NewListenerGroup(spec ListenerSpec) (grp *ListenerGroup, err error) {
	grp = &ListenerGroup{
		Name:      spec.Name,
		endpoints: make(map[string]*ListenerEndpoint),
		Unmanaged: spec.Unmanaged,
	}

	if grp.Addr, err = spec.Addr(); err != nil {
		return nil, err
	}

	if grp.Name == "" {
		switch a := grp.Addr.(type) {
		case *net.TCPAddr:
			grp.Name = fmt.Sprintf("%d/tcp", a.Port)
		case *net.UDPAddr:
			grp.Name = fmt.Sprintf("%d/udp", a.Port)
		}
	}

	return grp, nil
}

func NewListenerEndpoint(spec Spec, handler ProtocolHandler) *ListenerEndpoint {
	return &ListenerEndpoint{
		TLS:     spec.TLS,
		Handler: handler,
		Options: spec.Options,
	}
}

type Group struct {
	Names    []string
	Handlers map[string]MultiplexHandler
}

func (g Group) IsEmpty() bool {
	return len(g.Names) == 0
}

type ListenerGroup struct {
	lock          sync.Mutex
	multiplexers  []MuxServer
	errorHandlers []ErrorHandler
	endpoints     map[string]*ListenerEndpoint
	isServing     bool
	CloseTimeout  time.Duration
	Name          string
	Unmanaged     bool
	Addr          net.Addr
}

func (lg *ListenerGroup) AddErrorHandler(eh ErrorHandler) {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	lg.errorHandlers = append(lg.errorHandlers, eh)
}

func (lg *ListenerGroup) ConfigureEndpoint(name string, le *ListenerEndpoint) {
	if le == nil {
		return
	}

	lg.lock.Lock()
	defer lg.lock.Unlock()

	lg.endpoints[name] = le
}

func (lg *ListenerGroup) ConfiguredEndpoints() (eps []string) {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	eps = make([]string, 0, len(lg.endpoints))
	for _, ep := range lg.endpoints {
		eps = append(eps, ep.Name)
	}

	return eps
}

func (lg *ListenerGroup) Serve(ctx context.Context) error {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	for _, ep := range lg.endpoints {
		if err := ep.Handler.Start(ctx, NewStartupSpec(ep.Name, ep.Uplink, ep.Options)); err != nil {
			return err
		}
	}

	for idx := range lg.multiplexers {
		go func(mux MuxServer) {
			mux.HandleError(func(err error) bool {
				if err = IgnoreShutdownError(err); err != nil {
					lg.notifyErrorHandlers(err)
				}
				return true
			})
			mux.SetReadTimeout(defaultReadTimeout)
			if err := IgnoreShutdownError(mux.Serve()); err != nil {
				lg.notifyErrorHandlers(err)
			}
		}(lg.multiplexers[idx])
	}

	lg.isServing = true

	return nil
}

func (lg *ListenerGroup) SetupMux(mux cmux.CMux, grp *Group) {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	lg.multiplexers = append(lg.multiplexers, mux)
	for idx := range grp.Names {
		name := grp.Names[idx]
		ep := lg.endpoints[name]
		ep.Name = fmt.Sprintf("%s:%s", lg.Name, name)
		ep.Uplink.Addr = lg.Addr
		ep.Uplink.Listener = mux.Match(grp.Handlers[name].Matchers()...)
	}
}

func (lg *ListenerGroup) GroupByTLS() (plainGrp, tlsGrp *Group, err error) {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	if plainGrp, err = groupEndpoints(lg.endpoints, func(s *ListenerEndpoint) bool { return !s.TLS }); err != nil {
		return nil, nil, err
	}

	if tlsGrp, err = groupEndpoints(lg.endpoints, func(s *ListenerEndpoint) bool { return s.TLS }); err != nil {
		return nil, nil, err
	}

	return
}

func (lg *ListenerGroup) Shutdown(ctx context.Context) (err error) {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	closingTimeout := lg.CloseTimeout
	if closingTimeout <= 0 {
		closingTimeout = defaultClosingTimeout
	}

	for _, le := range lg.endpoints {
		closingCtx, cancel := context.WithTimeout(ctx, closingTimeout)
		err = IgnoreShutdownError(le.Close(closingCtx))
		cancel()
		if err != nil {
			break
		}
	}

	for idx := range lg.multiplexers {
		lg.multiplexers[idx].Close()
	}

	lg.multiplexers = nil
	lg.isServing = false

	return err
}

func (lg *ListenerGroup) notifyErrorHandlers(err error) {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	for idx := range lg.errorHandlers {
		lg.errorHandlers[idx].OnError(err)
	}
}

type ListenerEndpoint struct {
	Name    string
	TLS     bool
	Handler ProtocolHandler
	Uplink  Uplink
	Options map[string]interface{}
}

func (le ListenerEndpoint) Close(ctx context.Context) (err error) {
	if stoppable, ok := le.Handler.(StoppableHandler); ok {
		err = stoppable.Stop(ctx)
	}

	return multierr.Append(err, le.Uplink.Close())
}

func groupEndpoints(endpoints map[string]*ListenerEndpoint, predicate func(s *ListenerEndpoint) bool) (*Group, error) {
	grp := &Group{
		Names:    make([]string, 0, len(endpoints)),
		Handlers: make(map[string]MultiplexHandler),
	}

	for name, spec := range endpoints {
		var e MultiplexHandler
		if ep, ok := spec.Handler.(MultiplexHandler); !ok {
			return nil, fmt.Errorf("handler %s %w", spec.Name, ErrMultiplexingNotSupported)
		} else {
			e = ep
		}

		if predicate(spec) {
			grp.Names = append(grp.Names, name)
			grp.Handlers[name] = e
		}
	}
	sort.Strings(grp.Names)
	return grp, nil
}
