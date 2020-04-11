package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"github.com/baez90/inetmock/pkg/path"
	"go.uber.org/zap"
	"math/big"
	"net"
	"path/filepath"
)

type keyProvider func() (key interface{}, err error)

type certStore struct {
	options            *tlsOptions
	keyProvider        keyProvider
	caCert             *x509.Certificate
	caPrivateKey       interface{}
	certCache          map[string]*tls.Certificate
	timeSourceInstance timeSource
	logger             *zap.Logger
}

func (cs *certStore) timeSource() timeSource {
	if cs.timeSourceInstance == nil {
		cs.timeSourceInstance = createTimeSource()
	}
	return cs.timeSourceInstance
}

func (cs *certStore) initCaCert() (err error) {
	if path.FileExists(cs.options.rootCaCert.publicKeyPath) && path.FileExists(cs.options.rootCaCert.privateKeyPath) {
		var tlsCert tls.Certificate
		if tlsCert, err = tls.LoadX509KeyPair(cs.options.rootCaCert.publicKeyPath, cs.options.rootCaCert.privateKeyPath); err != nil {
			return
		} else if len(tlsCert.Certificate) > 0 {
			cs.caPrivateKey = tlsCert.PrivateKey
			cs.caCert, err = x509.ParseCertificate(tlsCert.Certificate[0])
		}
	} else {
		cs.caCert, cs.caPrivateKey, err = cs.generateCaCert()
	}
	return
}

func (cs *certStore) initProvider() {
	if cs.keyProvider == nil {
		cs.keyProvider = func() (key interface{}, err error) {
			return privateKeyForCurve(cs.options)
		}
	}
}

func (cs *certStore) generateCaCert() (pubKey *x509.Certificate, privateKey interface{}, err error) {
	cs.initProvider()
	timeSource := cs.timeSource()
	var serialNumber *big.Int
	if serialNumber, err = generateSerialNumber(); err != nil {
		return
	}

	if privateKey, err = cs.keyProvider(); err != nil {
		return
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"INetMock"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IsCA:                  true,
		NotBefore:             timeSource.UTCNow().Add(-cs.options.validity.ca.notBeforeRelative),
		NotAfter:              timeSource.UTCNow().Add(cs.options.validity.ca.notAfterRelative),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	var derBytes []byte
	if derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privateKey), privateKey); err != nil {
		return
	}

	pubKey, err = x509.ParseCertificate(derBytes)

	if err = writePublicKey(cs.options.rootCaCert.publicKeyPath, derBytes); err != nil {
		return
	}

	var privateKeyBytes []byte
	if privateKeyBytes, err = x509.MarshalPKCS8PrivateKey(privateKey); err != nil {
		return
	}

	if err = writePrivateKey(cs.options.rootCaCert.privateKeyPath, privateKeyBytes); err != nil {
		return
	}

	return
}

func (cs *certStore) getCertificate(serverName string, ip string) (cert *tls.Certificate, err error) {
	if cert, ok := cs.certCache[serverName]; ok {
		return cert, nil
	}

	certPath := filepath.Join(cs.options.certCachePath, fmt.Sprintf("%s.pem", serverName))
	keyPath := filepath.Join(cs.options.certCachePath, fmt.Sprintf("%s.key", serverName))

	if path.FileExists(certPath) && path.FileExists(keyPath) {
		if tlsCert, loadErr := tls.LoadX509KeyPair(certPath, keyPath); loadErr == nil {
			cs.certCache[serverName] = &tlsCert
			x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
			if err == nil && !certShouldBeRenewed(cs.timeSource(), x509Cert) {
				return &tlsCert, nil
			}
		}
	}

	if cert, err = cs.generateDomainCert(serverName, ip); err == nil {
		cs.certCache[serverName] = cert
	}

	return
}

func (cs *certStore) generateDomainCert(
	serverName string,
	localIp string,
) (cert *tls.Certificate, err error) {
	cs.initProvider()

	if cs.caCert == nil {
		if err = cs.initCaCert(); err != nil {
			return
		}
	}

	var serialNumber *big.Int
	if serialNumber, err = generateSerialNumber(); err != nil {
		return
	}

	var privateKey interface{}
	if privateKey, err = cs.keyProvider(); err != nil {
		return
	}

	notBefore := cs.timeSource().UTCNow().Add(-cs.options.validity.domain.notBeforeRelative)
	notAfter := cs.timeSource().UTCNow().Add(cs.options.validity.domain.notAfterRelative)

	cs.logger.Info(
		"generate domain certificate",
		zap.Time("notBefore", notBefore),
		zap.Time("notAfter", notAfter),
	)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"INetMock"},
		},
		IPAddresses: []net.IP{net.ParseIP(localIp)},
		DNSNames:    []string{serverName},
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}

	var derBytes []byte
	if derBytes, err = x509.CreateCertificate(rand.Reader, &template, cs.caCert, publicKey(privateKey), cs.caPrivateKey); err != nil {
		return
	}

	if err = writePublicKey(filepath.Join(cs.options.certCachePath, fmt.Sprintf("%s.pem", serverName)), derBytes); err != nil {
		return
	}

	var privateKeyBytes []byte
	if privateKeyBytes, err = x509.MarshalPKCS8PrivateKey(privateKey); err != nil {
		return
	}

	if cert, err = parseCert(derBytes, privateKeyBytes); err != nil {
		return
	}

	if err = writePrivateKey(filepath.Join(cs.options.certCachePath, fmt.Sprintf("%s.key", serverName)), privateKeyBytes); err != nil {
		return
	}

	return
}
