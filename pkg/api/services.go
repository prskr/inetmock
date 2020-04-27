package api

import (
	"github.com/baez90/inetmock/pkg/cert"
	config2 "github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
)

var (
	svcs Services
)

type Services interface {
	CertStore() cert.Store
}

type services struct {
	certStore cert.Store
}

func InitServices(
	config config2.Config,
	logger logging.Logger,
) error {
	certStore, err := cert.NewDefaultStore(config, logger)
	if err != nil {
		return err
	}
	svcs = &services{
		certStore: certStore,
	}
	return nil
}

func ServicesInstance() Services {
	return svcs
}

func (s *services) CertStore() cert.Store {
	return s.certStore
}
