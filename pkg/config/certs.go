package config

import (
	"crypto/tls"
	"strings"
	"time"
)

type CurveType string

type TLSVersion string

func (x TLSVersion) Value() string {
	return strings.ToUpper(string(x))
}

func (x TLSVersion) TLSVersion() uint16 {
	switch TLSVersion(x.Value()) {
	case TLSVersionSSL3:
		return tls.VersionSSL30
	case TLSVersionTLS10:
		return tls.VersionTLS10
	case TLSVersionTLS11:
		return tls.VersionTLS11
	case TLSVersionTLS12:
		return tls.VersionTLS12
	default:
		return tls.VersionTLS13
	}
}

type File struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

type ValidityDuration struct {
	NotBeforeRelative time.Duration
	NotAfterRelative  time.Duration
}

type ValidityByPurpose struct {
	CA     ValidityDuration
	Server ValidityDuration
}

type CertOptions struct {
	RootCACert                  File
	CertCachePath               string
	Curve                       CurveType
	Validity                    ValidityByPurpose
	IncludeInsecureCipherSuites bool
	MinTLSVersion               TLSVersion
}
