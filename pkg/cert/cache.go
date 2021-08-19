//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/cert/cert_cache.mock.go -package=cert_mock

package cert

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"sync"
)

type Cache interface {
	Put(cert *tls.Certificate) error
	Get(cn string) (*tls.Certificate, bool)
}

func NewFileSystemCache(certCachePath string, source TimeSource) Cache {
	m := new(sync.RWMutex)
	return &fileSystemCache{
		certCachePath: certCachePath,
		inMemCache:    make(map[string]*tls.Certificate),
		timeSource:    source,
		readLock:      m.RLocker(),
		writeLock:     m,
	}
}

type fileSystemCache struct {
	certCachePath string
	readLock      sync.Locker
	writeLock     sync.Locker
	inMemCache    map[string]*tls.Certificate
	timeSource    TimeSource
}

func (f *fileSystemCache) Put(cert *tls.Certificate) (err error) {
	f.writeLock.Lock()
	defer f.writeLock.Unlock()
	if cert == nil {
		err = errors.New("cert may not be nil")
		return
	}
	var cn string
	if len(cert.Certificate) > 0 {
		var pubKey *x509.Certificate
		if pubKey, err = x509.ParseCertificate(cert.Certificate[0]); err != nil {
			return err
		} else {
			cn = pubKey.Subject.CommonName
		}

		f.inMemCache[cn] = cert
		pemCrt := NewPEM(cert)
		err = pemCrt.Write(cn, f.certCachePath)
	} else {
		err = errors.New("no public key present for certificate")
	}
	return
}

func (f *fileSystemCache) Get(cn string) (*tls.Certificate, bool) {
	f.readLock.Lock()
	defer f.readLock.Unlock()

	if crt, ok := f.inMemCache[cn]; ok {
		return crt, true
	}

	pemCrt := NewPEM(nil)
	if err := pemCrt.Read(cn, f.certCachePath); err != nil || pemCrt.Cert() == nil {
		return nil, false
	}

	x509Cert, err := x509.ParseCertificate(pemCrt.Cert().Certificate[0])
	if err == nil && !certShouldBeRenewed(f.timeSource, x509Cert) {
		return pemCrt.Cert(), true
	}

	return nil, false
}

func certShouldBeRenewed(timeSource TimeSource, cert *x509.Certificate) bool {
	lifetime := cert.NotAfter.Sub(cert.NotBefore)
	// if the cert is closer to the end of the lifetime than lifetime/2 it should be renewed
	return cert.NotAfter.Sub(timeSource.UTCNow()) < lifetime/4
}
