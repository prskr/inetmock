package cert

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

var (
	requiredConfigKeys = []string{
		publicKeyConfigKey,
		privateKeyPathConfigKey,
	}
)

type CurveType string

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

type Options struct {
	RootCACert    File
	CertCachePath string
	Curve         CurveType
	Validity      ValidityByPurpose
}

func loadFromConfig(config *viper.Viper) (Options, error) {
	missingKeys := make([]string, 0)
	for _, requiredKey := range requiredConfigKeys {
		if !config.IsSet(requiredKey) {
			missingKeys = append(missingKeys, requiredKey)
		}
	}

	if len(missingKeys) > 0 {
		return Options{}, fmt.Errorf("config keys are missing: %s", strings.Join(missingKeys, ", "))
	}

	config.SetDefault(certCachePathConfigKey, os.TempDir())
	config.SetDefault(ecdsaCurveConfigKey, string(CurveTypeED25519))
	config.SetDefault(caCertValidityNotBeforeKey, defaultCAValidityDuration)
	config.SetDefault(caCertValidityNotAfterKey, defaultCAValidityDuration)
	config.SetDefault(serverCertValidityNotBeforeKey, defaultServerValidityDuration)
	config.SetDefault(serverCertValidityNotAfterKey, defaultServerValidityDuration)

	return Options{
		CertCachePath: config.GetString(certCachePathConfigKey),
		Curve:         CurveType(config.GetString(ecdsaCurveConfigKey)),
		Validity: ValidityByPurpose{
			CA: ValidityDuration{
				NotBeforeRelative: config.GetDuration(caCertValidityNotBeforeKey),
				NotAfterRelative:  config.GetDuration(caCertValidityNotAfterKey),
			},
			Server: ValidityDuration{
				NotBeforeRelative: config.GetDuration(serverCertValidityNotBeforeKey),
				NotAfterRelative:  config.GetDuration(serverCertValidityNotAfterKey),
			},
		},
		RootCACert: File{
			PublicKeyPath:  config.GetString(publicKeyConfigKey),
			PrivateKeyPath: config.GetString(privateKeyPathConfigKey),
		},
	}, nil
}
