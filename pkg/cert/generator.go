package cert

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
)

type GenerationOptions struct {
	CommonName         string
	Organization       []string
	OrganizationalUnit []string
	IPAddresses        []net.IP
	DNSNames           []string
	Country            []string
	Province           []string
	Locality           []string
	StreetAddress      []string
	PostalCode         []string
}

type Generator interface {
	CACert(options GenerationOptions) (*tls.Certificate, error)
	ServerCert(options GenerationOptions, ca *tls.Certificate) (*tls.Certificate, error)
}

func NewDefaultGenerator(options Options) Generator {
	return NewGenerator(options, NewTimeSource(), defaultKeyProvider(options))
}

func NewGenerator(options Options, source TimeSource, provider KeyProvider) Generator {
	return &generator{
		options:    options,
		provider:   provider,
		timeSource: source,
	}
}

type generator struct {
	options    Options
	provider   KeyProvider
	timeSource TimeSource
}

func (g *generator) privateKey() (key interface{}, err error) {
	if g.provider != nil {
		return g.provider()
	} else {
		return defaultKeyProvider(g.options)()
	}
}

//nolint:gocritic
func (g *generator) ServerCert(options GenerationOptions, ca *tls.Certificate) (*tls.Certificate, error) {
	var err = applyDefaultGenerationOptions(&options)
	if err != nil {
		return nil, err
	}
	var serialNumber *big.Int
	if serialNumber, err = generateSerialNumber(); err != nil {
		return nil, err
	}

	var privateKey interface{}
	if privateKey, err = g.privateKey(); err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         options.CommonName,
			Organization:       options.Organization,
			OrganizationalUnit: options.OrganizationalUnit,
			Country:            options.Country,
			Province:           options.Province,
			Locality:           options.Locality,
			StreetAddress:      options.StreetAddress,
			PostalCode:         options.PostalCode,
		},
		IPAddresses: options.IPAddresses,
		DNSNames:    options.DNSNames,
		NotBefore:   g.timeSource.UTCNow().Add(-g.options.Validity.Server.NotBeforeRelative),
		NotAfter:    g.timeSource.UTCNow().Add(g.options.Validity.Server.NotAfterRelative),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	var caCrt *x509.Certificate
	if caCrt, err = x509.ParseCertificate(ca.Certificate[0]); err != nil {
		return nil, err
	}

	var derBytes []byte
	if derBytes, err = x509.CreateCertificate(rand.Reader, &template, caCrt, publicKey(privateKey), ca.PrivateKey); err != nil {
		return nil, err
	}

	var privateKeyBytes []byte
	if privateKeyBytes, err = x509.MarshalPKCS8PrivateKey(privateKey); err != nil {
		return nil, err
	}

	var cert *tls.Certificate
	if cert, err = parseCert(derBytes, privateKeyBytes); err != nil {
		return nil, err
	}

	return cert, nil
}

//nolint:gocritic
func (g *generator) CACert(options GenerationOptions) (*tls.Certificate, error) {
	var err = applyDefaultGenerationOptions(&options)
	if err != nil {
		return nil, err
	}

	var privateKey interface{}
	var serialNumber *big.Int
	if serialNumber, err = generateSerialNumber(); err != nil {
		return nil, err
	}

	if privateKey, err = g.privateKey(); err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:    options.CommonName,
			Organization:  options.Organization,
			Country:       options.Country,
			Province:      options.Province,
			Locality:      options.Locality,
			StreetAddress: options.StreetAddress,
			PostalCode:    options.PostalCode,
		},
		IsCA:                  true,
		NotBefore:             g.timeSource.UTCNow().Add(-g.options.Validity.CA.NotBeforeRelative),
		NotAfter:              g.timeSource.UTCNow().Add(g.options.Validity.CA.NotAfterRelative),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	var derBytes []byte
	if derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privateKey), privateKey); err != nil {
		return nil, err
	}

	var privateKeyBytes []byte
	if privateKeyBytes, err = x509.MarshalPKCS8PrivateKey(privateKey); err != nil {
		return nil, err
	}

	var cert *tls.Certificate
	if cert, err = parseCert(derBytes, privateKeyBytes); err != nil {
		return nil, err
	}

	return cert, nil
}

func generateSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, serialNumberLimit)
}

func publicKey(privateKey interface{}) interface{} {
	switch k := privateKey.(type) {
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func parseCert(derBytes, privateKeyBytes []byte) (*tls.Certificate, error) {
	pemEncodedPublicKey := pem.EncodeToMemory(&pem.Block{Type: certificateBlockType, Bytes: derBytes})
	pemEncodedPrivateKey := pem.EncodeToMemory(&pem.Block{Type: privateKeyBlockType, Bytes: privateKeyBytes})
	cert, err := tls.X509KeyPair(pemEncodedPublicKey, pemEncodedPrivateKey)
	return &cert, err
}
