package cert

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
)

const (
	certificateBlockType = "CERTIFICATE"
	privateKeyBlockType  = "PRIVATE KEY"
)

type PEMCert struct {
	*tls.Certificate
}

func (p PEMCert) Write(cn, outDir string) (err error) {
	var (
		certOut      io.WriteCloser
		keyOut       io.WriteCloser
		privKeyBytes []byte
	)
	if certOut, err = os.Create(filepath.Join(outDir, fmt.Sprintf("%s.pem", cn))); err != nil {
		return
	}
	defer multierr.AppendInvoke(&err, multierr.Close(certOut))
	if err = pem.Encode(certOut, &pem.Block{Type: certificateBlockType, Bytes: p.Certificate.Certificate[0]}); err != nil {
		return
	}

	if keyOut, err = os.Create(filepath.Join(outDir, fmt.Sprintf("%s.key", cn))); err != nil {
		return
	}
	privKeyBytes, err = x509.MarshalPKCS8PrivateKey(p.Certificate.PrivateKey)
	err = pem.Encode(keyOut, &pem.Block{Type: privateKeyBlockType, Bytes: privKeyBytes})
	return
}

func Read(cn, inDir string) (*PEMCert, error) {
	certPath := filepath.Join(inDir, fmt.Sprintf("%s.pem", cn))
	keyPath := filepath.Join(inDir, fmt.Sprintf("%s.key", cn))

	return ReadFrom(certPath, keyPath)
}

func ReadFrom(pubKeyPath, privateKeyPath string) (*PEMCert, error) {
	if tlsCrt, err := tls.LoadX509KeyPair(pubKeyPath, privateKeyPath); err != nil {
		return nil, err
	} else {
		return &PEMCert{Certificate: &tlsCrt}, nil
	}
}
