package cert

const (
	CurveTypeP224    CurveType = "P224"
	CurveTypeP256    CurveType = "P256"
	CurveTypeP384    CurveType = "P384"
	CurveTypeP521    CurveType = "P521"
	CurveTypeED25519 CurveType = "ED25519"

	defaultServerValidityDuration = "168h"
	defaultCAValidityDuration     = "17520h"

	certCachePathConfigKey         = "tls.certCachePath"
	ecdsaCurveConfigKey            = "tls.ecdsaCurve"
	publicKeyConfigKey             = "tls.rootCaCert.publicKey"
	privateKeyPathConfigKey        = "tls.rootCaCert.privateKey"
	caCertValidityNotBeforeKey     = "tls.validity.ca.notBeforeRelative"
	caCertValidityNotAfterKey      = "tls.validity.ca.notAfterRelative"
	serverCertValidityNotBeforeKey = "tls.validity.server.notBeforeRelative"
	serverCertValidityNotAfterKey  = "tls.validity.server.notAfterRelative"
)
