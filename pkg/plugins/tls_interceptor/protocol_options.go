package main

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const (
	certCachePathConfigKey         = "certCachePath"
	ecdsaCurveConfigKey            = "ecdsaCurve"
	targetIpAddressConfigKey       = "target.ipAddress"
	targetPortConfigKey            = "target.port"
	publicKeyConfigKey             = "rootCaCert.publicKey"
	privateKeyPathConfigKey        = "rootCaCert.privateKey"
	caCertValidityNotBeforeKey     = "validity.ca.notBeforeRelative"
	caCertValidityNotAfterKey      = "validity.ca.notAfterRelative"
	domainCertValidityNotBeforeKey = "validity.domain.notBeforeRelative"
	domainCertValidityNotAfterKey  = "validity.domain.notAfterRelative"
)

type cert struct {
	publicKeyPath  string
	privateKeyPath string
}

type certValidity struct {
	notBeforeRelative time.Duration
	notAfterRelative  time.Duration
}

type validity struct {
	ca     certValidity
	domain certValidity
}

type redirectionTarget struct {
	ipAddress string
	port      uint16
}

func (rt redirectionTarget) address() string {
	return fmt.Sprintf("%s:%d", rt.ipAddress, rt.port)
}

type tlsOptions struct {
	rootCaCert        cert
	certCachePath     string
	redirectionTarget redirectionTarget
	ecdsaCurve        curveType
	validity          validity
}

func loadFromConfig(config *viper.Viper) *tlsOptions {

	return &tlsOptions{
		certCachePath: config.GetString(certCachePathConfigKey),
		ecdsaCurve:    curveType(config.GetString(ecdsaCurveConfigKey)),
		redirectionTarget: redirectionTarget{
			ipAddress: config.GetString(targetIpAddressConfigKey),
			port:      uint16(config.GetInt(targetPortConfigKey)),
		},
		validity: validity{
			ca: certValidity{
				notBeforeRelative: config.GetDuration(caCertValidityNotBeforeKey),
				notAfterRelative:  config.GetDuration(caCertValidityNotAfterKey),
			},
			domain: certValidity{
				notBeforeRelative: config.GetDuration(domainCertValidityNotBeforeKey),
				notAfterRelative:  config.GetDuration(domainCertValidityNotAfterKey),
			},
		},
		rootCaCert: cert{
			publicKeyPath:  config.GetString(publicKeyConfigKey),
			privateKeyPath: config.GetString(privateKeyPathConfigKey),
		},
	}
}
