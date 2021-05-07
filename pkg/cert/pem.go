package cert

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/multierr"

	"gitlab.com/inetmock/inetmock/pkg/path"
)

const (
	certificateBlockType = "CERTIFICATE"
	privateKeyBlockType  = "PRIVATE KEY"
)

type PEM interface {
	Cert() *tls.Certificate
	Write(cn string, outDir string) error
	Read(cn string, inDir string) error
	ReadFrom(pubKeyPath, privateKeyPath string) error
}

func NewPEM(crt *tls.Certificate) PEM {
	return &pemCrt{
		crt: crt,
	}
}

type pemCrt struct {
	crt *tls.Certificate
}

func (p *pemCrt) Cert() *tls.Certificate {
	return p.crt
}

func (p pemCrt) Write(cn, outDir string) (err error) {
	var certOut *os.File
	if certOut, err = os.Create(filepath.Join(outDir, fmt.Sprintf("%s.pem", cn))); err != nil {
		return
	}
	defer multierr.AppendInvoke(&err, multierr.Close(certOut))
	if err = pem.Encode(certOut, &pem.Block{Type: certificateBlockType, Bytes: p.crt.Certificate[0]}); err != nil {
		return
	}

	var keyOut *os.File
	if keyOut, err = os.OpenFile(filepath.Join(outDir, fmt.Sprintf("%s.key", cn)), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return
	}
	var privKeyBytes []byte
	privKeyBytes, err = x509.MarshalPKCS8PrivateKey(p.crt.PrivateKey)
	err = pem.Encode(keyOut, &pem.Block{Type: privateKeyBlockType, Bytes: privKeyBytes})
	return
}

func (p *pemCrt) Read(cn, inDir string) error {
	certPath := filepath.Join(inDir, fmt.Sprintf("%s.pem", cn))
	keyPath := filepath.Join(inDir, fmt.Sprintf("%s.key", cn))

	return p.ReadFrom(certPath, keyPath)
}

func (p *pemCrt) ReadFrom(pubKeyPath, privateKeyPath string) (err error) {
	var tlsCrt tls.Certificate
	if path.FileExists(pubKeyPath) && path.FileExists(privateKeyPath) {
		if tlsCrt, err = tls.LoadX509KeyPair(pubKeyPath, privateKeyPath); err == nil {
			p.crt = &tlsCrt
		}
	} else {
		err = errors.New("either public or private key file do not exist")
	}

	return
}
