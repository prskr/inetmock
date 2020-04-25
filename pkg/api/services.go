package api

import (
	"github.com/baez90/inetmock/pkg/cert"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/spf13/viper"
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
	config *viper.Viper,
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
