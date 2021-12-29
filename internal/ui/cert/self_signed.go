package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"
)

type Request struct {
	DNSNames    []string
	IPAddresses []net.IP
}

func GenerateServerCert(req Request) (cert *tls.Certificate, err error) {
	const defaultCertValidFor = 6 * time.Hour
	var priv *ecdsa.PrivateKey
	if priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil {
		return nil, err
	}

	keyUsage := x509.KeyUsageDigitalSignature
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	var serialNumber *big.Int
	if serialNumber, err = rand.Int(rand.Reader, serialNumberLimit); err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(defaultCertValidFor),

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              req.DNSNames,
		IPAddresses:           req.IPAddresses,
	}

	var derBytes []byte
	if derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv); err != nil {
		return nil, err
	}

	cert = new(tls.Certificate)
	cert.Certificate = append(cert.Certificate, derBytes)
	cert.PrivateKey = priv

	return cert, nil
}
