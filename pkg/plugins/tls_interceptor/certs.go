package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
)

type curveType string

const (
	certificateBlockType = "CERTIFICATE"
	privateKeyBlockType  = "PRIVATE KEY"

	curveTypeP224    curveType = "P224"
	curveTypeP256    curveType = "P256"
	curveTypeP384    curveType = "P384"
	curveTypeP521    curveType = "P521"
	curveTypeED25519 curveType = "ED25519"
)

func loadPEMCert(certPEMBytes []byte, keyPEMBytes []byte) (*tls.Certificate, error) {
	cert, err := tls.X509KeyPair(certPEMBytes, keyPEMBytes)
	return &cert, err
}

func parseCert(derBytes []byte, privateKeyBytes []byte) (*tls.Certificate, error) {
	pemEncodedPublicKey := pem.EncodeToMemory(&pem.Block{Type: certificateBlockType, Bytes: derBytes})
	pemEncodedPrivateKey := pem.EncodeToMemory(&pem.Block{Type: privateKeyBlockType, Bytes: privateKeyBytes})
	cert, err := tls.X509KeyPair(pemEncodedPublicKey, pemEncodedPrivateKey)
	return &cert, err
}

func writePublicKey(crtOutPath string, derBytes []byte) (err error) {
	var certOut *os.File
	if certOut, err = os.Create(crtOutPath); err != nil {
		return
	}
	if err = pem.Encode(certOut, &pem.Block{Type: certificateBlockType, Bytes: derBytes}); err != nil {
		return
	}
	err = certOut.Close()
	return
}

func writePrivateKey(keyOutPath string, privateKeyBytes []byte) (err error) {
	var keyOut *os.File
	if keyOut, err = os.OpenFile(keyOutPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return
	}

	if err = pem.Encode(keyOut, &pem.Block{Type: privateKeyBlockType, Bytes: privateKeyBytes}); err != nil {
		return
	}

	err = keyOut.Close()
	return
}

func certShouldBeRenewed(timeSource timeSource, cert *x509.Certificate) bool {
	lifetime := cert.NotAfter.Sub(cert.NotBefore)
	// if the cert is closer to the end of the lifetime than lifetime/2 it should be renewed
	if cert.NotAfter.Sub(timeSource.UTCNow()) < lifetime/4 {
		return true
	}
	return false
}

func generateSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, serialNumberLimit)
}

func privateKeyForCurve(options *tlsOptions) (privateKey interface{}, err error) {
	switch options.ecdsaCurve {
	case "P224":
		privateKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		_, privateKey, err = ed25519.GenerateKey(rand.Reader)
	}

	return
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
