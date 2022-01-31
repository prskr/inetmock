package endpoint

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/soheilhy/cmux"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/netutils"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var (
	ErrNoSuchGroup      = errors.New("no group with given name configured")
	_              Host = (*Server)(nil)
)

type Link struct {
	Listener   net.Listener
	PacketConn net.PacketConn
}

func NewServer(logger logging.Logger, tlsConfig *tls.Config) *Server {
	return &Server{
		groups:    make(map[string]*ListenerGroup),
		TLSConfig: tlsConfig,
		Logger:    logger,
	}
}

type Server struct {
	lock         sync.Mutex
	groups       map[string]*ListenerGroup
	TLSConfig    *tls.Config
	ErrorHandler []ErrorHandler
	Logger       logging.Logger
}

func (s *Server) ConfigureGroup(grp *ListenerGroup) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.groups[grp.Name] = grp
}

func (s *Server) ConfiguredGroups() []GroupInfo {
	s.lock.Lock()
	defer s.lock.Unlock()

	infos := make([]GroupInfo, 0, len(s.groups))

	for name, grp := range s.groups {
		info := GroupInfo{
			Name:      name,
			Endpoints: grp.ConfiguredEndpoints(),
			Serving:   grp.isServing,
		}

		infos = append(infos, info)
	}

	return infos
}

func (s *Server) ServeGroup(ctx context.Context, groupName string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if g, exists := s.groups[groupName]; !exists {
		return fmt.Errorf("%w: %s", ErrNoSuchGroup, groupName)
	} else {
		return s.serveGroup(ctx, g)
	}
}

func (s *Server) ServeGroups(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, grp := range s.groups {
		if err := s.serveGroup(ctx, grp); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) ShutdownGroup(ctx context.Context, groupName string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if grpToShutDown, exists := s.groups[groupName]; !exists {
		return ErrNoSuchGroup
	} else {
		return grpToShutDown.Shutdown(ctx)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, grp := range s.groups {
		if err := grp.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) ShutdownOnCancel(ctx context.Context) {
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			s.lock.Lock()
			defer s.lock.Unlock()
			for idx := range s.ErrorHandler {
				s.ErrorHandler[idx].OnError(err)
			}
		}
	}()
}

func (s *Server) serveGroup(ctx context.Context, grpToStart *ListenerGroup) error {
	grpToStart.errorHandlers = make([]ErrorHandler, len(s.ErrorHandler))
	copy(grpToStart.errorHandlers, s.ErrorHandler)

	s.Logger.Debug("Starting endpoint group", zap.String("group_name", grpToStart.Name), zap.String("addr", grpToStart.Addr.String()))
	if err := s.setupGroup(grpToStart, s.TLSConfig); err != nil {
		return err
	}
	if err := grpToStart.Serve(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Server) setupGroup(grp *ListenerGroup, tlsConfig *tls.Config) error {
	var uplink *Uplink
	if u, err := s.setupUplink(grp); err != nil {
		return err
	} else {
		uplink = u
	}

	setupLogger := s.Logger.With(zap.String("group_name", grp.Name))

	if len(grp.endpoints) <= 1 {
		for name, le := range grp.endpoints {
			setupLogger.Debug("Preparing single handler group",
				zap.String("handler_name", name),
				zap.Bool("tls", le.TLS),
			)
			if le.TLS {
				uplink.Listener = tls.NewListener(uplink.Listener, tlsConfig)
			}
			le.Name = fmt.Sprintf("%s:%s", grp.Name, name)
			le.Uplink = *uplink
			return nil
		}
	}

	if uplink.IsUDP() {
		return ErrUDPMultiplexer
	}

	plainGrp, tlsGrp, err := grp.GroupByTLS()
	if err != nil {
		return err
	}

	setupLogger.Debug("Preparing multiplexing group")

	lis := uplink.Listener

	if !plainGrp.IsEmpty() {
		setupLogger.Debug("Configuring plain text endpoints")
		plainMux := cmux.New(lis)
		grp.SetupMux(plainMux, plainGrp)
		lis = plainMux.Match(cmux.Any())
	}

	if !tlsGrp.IsEmpty() {
		setupLogger.Debug("Configuring TLS endpoints")
		tlsMux := cmux.New(tls.NewListener(lis, tlsConfig))
		grp.SetupMux(tlsMux, tlsGrp)
	}

	return nil
}

func (s *Server) setupUplink(grp *ListenerGroup) (u *Uplink, err error) {
	u = &Uplink{Unmanaged: grp.Unmanaged, Addr: grp.Addr}
	if grp.Unmanaged {
		return
	}

	switch a := grp.Addr.(type) {
	case *net.UDPAddr:
		u.PacketConn, err = net.ListenUDP("udp", a)
	case *net.TCPAddr:
		u.Listener, err = netutils.ListenTCP(a)
	}
	return
}
