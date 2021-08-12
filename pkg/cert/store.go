package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"net"

	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/netutils"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var (
	defaultKeyProvider = func(options Options) func() (key interface{}, err error) {
		return func() (key interface{}, err error) {
			return privateKeyForCurve(options)
		}
	}

	//nolint:gomnd // IPv4 loopback address is well known
	ipv4LoopbackIP = net.IPv4(127, 0, 0, 1)
)

type KeyProvider func() (key interface{}, err error)

type Store interface {
	CACert() *tls.Certificate
	GetCertificate(serverName string, ip net.IP) (*tls.Certificate, error)
	TLSConfig() *tls.Config
}

func MustDefaultStore(
	options Options,
	logger logging.Logger,
) Store {
	if store, err := NewDefaultStore(options, logger); err != nil {
		panic(err)
	} else {
		return store
	}
}

func NewDefaultStore(
	options Options,
	logger logging.Logger,
) (Store, error) {
	timeSource := NewTimeSource()
	return NewStore(
		options,
		NewFileSystemCache(options.CertCachePath, timeSource),
		NewDefaultGenerator(options),
		logger,
	)
}

func NewStore(
	options Options,
	cache Cache,
	generator Generator,
	logger logging.Logger,
) (Store, error) {
	store := &store{
		options:   options,
		cache:     cache,
		generator: generator,
		logger:    logger,
	}

	if err := store.loadCACert(); err != nil {
		return nil, err
	}

	return store, nil
}

type store struct {
	options   Options
	caCert    *tls.Certificate
	cache     Cache
	generator Generator
	logger    logging.Logger
}

func (s *store) TLSConfig() *tls.Config {
	rootCaPool := x509.NewCertPool()
	rootCaPubKey, _ := x509.ParseCertificate(s.caCert.Certificate[0])
	rootCaPool.AddCert(rootCaPubKey)

	suites := make([]uint16, 0)

	for _, suite := range tls.CipherSuites() {
		suites = append(suites, suite.ID)
	}

	if s.options.IncludeInsecureCipherSuites {
		for _, suite := range tls.InsecureCipherSuites() {
			suites = append(suites, suite.ID)
		}
	}

	//nolint:gosec
	return &tls.Config{
		CipherSuites:             suites,
		MinVersion:               s.options.MinTLSVersion.TLSVersion(),
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		GetCertificate: func(info *tls.ClientHelloInfo) (cert *tls.Certificate, err error) {
			var localIP *netutils.IPPort
			if localIP, err = netutils.IPPortFromAddress(info.Conn.LocalAddr()); err != nil {
				localIP = &netutils.IPPort{
					IP: ipv4LoopbackIP,
				}
			}

			if cert, err = s.GetCertificate(info.ServerName, localIP.IP); err != nil {
				s.logger.Error(
					"error while resolving certificate",
					zap.String("serverName", info.ServerName),
					zap.String("localAddr", localIP.String()),
					zap.Error(err),
				)
			}

			return
		},
		RootCAs: rootCaPool,
	}
}

func (s *store) loadCACert() (err error) {
	pemCrt := NewPEM(nil)
	if err = pemCrt.ReadFrom(s.options.RootCACert.PublicKeyPath, s.options.RootCACert.PrivateKeyPath); err != nil {
		return
	}
	s.caCert = pemCrt.Cert()
	return
}

func (s *store) CACert() *tls.Certificate {
	return s.caCert
}

func (s *store) GetCertificate(serverName string, ip net.IP) (cert *tls.Certificate, err error) {
	if crt, ok := s.cache.Get(serverName); ok {
		return crt, nil
	}

	if cert, err = s.generator.ServerCert(GenerationOptions{
		CommonName:  serverName,
		DNSNames:    []string{serverName},
		IPAddresses: []net.IP{ip},
	}, s.caCert); err == nil {
		_ = s.cache.Put(cert)
	}

	return
}

func privateKeyForCurve(options Options) (privateKey interface{}, err error) {
	switch options.Curve {
	case CurveTypeP224:
		privateKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case CurveTypeP256:
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case CurveTypeP384:
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case CurveTypeP521:
		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	case CurveTypeED25519:
		fallthrough
	default:
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}

	return
}
