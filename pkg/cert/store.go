package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"net"

	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

const (
	ipv4Loopback = "127.0.0.1"
)

var (
	defaultKeyProvider = func(options CertOptions) func() (key interface{}, err error) {
		return func() (key interface{}, err error) {
			return privateKeyForCurve(options)
		}
	}
)

type KeyProvider func() (key interface{}, err error)

type Store interface {
	CACert() *tls.Certificate
	GetCertificate(serverName string, ip string) (*tls.Certificate, error)
	TLSConfig() *tls.Config
}

func NewDefaultStore(
	options CertOptions,
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
	options CertOptions,
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
	options    CertOptions
	caCert     *tls.Certificate
	cache      Cache
	timeSource TimeSource
	generator  Generator
	logger     logging.Logger
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

	return &tls.Config{
		CipherSuites:             suites,
		MinVersion:               s.options.MinTLSVersion.TLSVersion(),
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		GetCertificate: func(info *tls.ClientHelloInfo) (cert *tls.Certificate, err error) {
			var localIp string
			if localIp, err = extractIPFromAddress(info.Conn.LocalAddr().String()); err != nil {
				localIp = ipv4Loopback
			}

			if cert, err = s.GetCertificate(info.ServerName, localIp); err != nil {
				s.logger.Error(
					"error while resolving certificate",
					zap.String("serverName", info.ServerName),
					zap.String("localAddr", localIp),
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

func (s *store) GetCertificate(serverName string, ip string) (cert *tls.Certificate, err error) {
	if crt, ok := s.cache.Get(serverName); ok {
		return crt, nil
	}

	if cert, err = s.generator.ServerCert(GenerationOptions{
		CommonName:  serverName,
		DNSNames:    []string{serverName},
		IPAddresses: []net.IP{net.ParseIP(ip)},
	}, s.caCert); err == nil {
		s.cache.Put(cert)
	}

	return
}

func privateKeyForCurve(options CertOptions) (privateKey interface{}, err error) {
	switch options.Curve {
	case CurveTypeP224:
		privateKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case CurveTypeP256:
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case CurveTypeP384:
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case CurveTypeP521:
		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}

	return
}
