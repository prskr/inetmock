package cert

const (
	defaultServerValidityDuration = "168h"
	defaultCAValidityDuration     = "17520h"

	certCachePathConfigKey         = "tls.certCachePath"
	ecdsaCurveConfigKey            = "tls.ecdsaCurve"
	caCertValidityNotBeforeKey     = "tls.validity.ca.notBeforeRelative"
	caCertValidityNotAfterKey      = "tls.validity.ca.notAfterRelative"
	serverCertValidityNotBeforeKey = "tls.validity.server.notBeforeRelative"
	serverCertValidityNotAfterKey  = "tls.validity.server.notAfterRelative"
)
