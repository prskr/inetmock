package cert

import (
	"github.com/baez90/inetmock/pkg/config"
	"os"
)

func init() {
	config.AddDefaultValue(certCachePathConfigKey, os.TempDir())
	config.AddDefaultValue(ecdsaCurveConfigKey, string(config.CurveTypeED25519))
	config.AddDefaultValue(caCertValidityNotBeforeKey, defaultCAValidityDuration)
	config.AddDefaultValue(caCertValidityNotAfterKey, defaultCAValidityDuration)
	config.AddDefaultValue(serverCertValidityNotBeforeKey, defaultServerValidityDuration)
	config.AddDefaultValue(serverCertValidityNotAfterKey, defaultServerValidityDuration)
}
