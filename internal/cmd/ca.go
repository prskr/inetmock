package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/baez90/inetmock/pkg/cert"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

const (
	generateCACommonName             = "cn"
	generateCaOrganizationName       = "o"
	generateCaOrganizationalUnitName = "ou"
	generateCaCountryName            = "c"
	generateCaLocalityName           = "l"
	generateCaStateName              = "st"
	generateCaStreetAddressName      = "street-address"
	generateCaPostalCodeName         = "postal-code"
	generateCACertOutPath            = "out-dir"
	generateCACurveName              = "curve"
	generateCANotBeforeRelative      = "not-before"
	generateCANotAfterRelative       = "not-after"
)

var (
	generateCaCmd *cobra.Command
	caCertOptions cert.GenerationOptions
)

func init() {
	generateCaCmd = &cobra.Command{
		Use:   "generate-ca",
		Short: "Generate a new CA certificate and corresponding key",
		Long:  ``,
		Run:   runGenerateCA,
	}

	generateCaCmd.Flags().StringVar(&caCertOptions.CommonName, generateCACommonName, "INetMock", "Certificate Common Name that will also be used as file name during generation.")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.Organization, generateCaOrganizationName, nil, "Organization information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.OrganizationalUnit, generateCaOrganizationalUnitName, nil, "Organizational unit information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.Country, generateCaCountryName, nil, "Country information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.Province, generateCaStateName, nil, "State information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.Locality, generateCaLocalityName, nil, "Locality information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.StreetAddress, generateCaStreetAddressName, nil, "Street address information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.PostalCode, generateCaPostalCodeName, nil, "Postal code information to append to certificate")
	generateCaCmd.Flags().String(generateCACertOutPath, "", "Path where CA files should be stored")
	generateCaCmd.Flags().String(generateCACurveName, "", "Name of the curve to use, if empty ED25519 is used, other valid values are [P224, P256,P384,P521]")
	generateCaCmd.Flags().Duration(generateCANotBeforeRelative, 17520*time.Hour, "Relative time value since when in the past the CA certificate should be valid. The value has a time unit, the greatest time unit is h for hour.")
	generateCaCmd.Flags().Duration(generateCANotAfterRelative, 17520*time.Hour, "Relative time value until when in the future the CA certificate should be valid. The value has a time unit, the greatest time unit is h for hour.")
}

func runGenerateCA(_ *cobra.Command, _ []string) {
	var certOutPath, curveName string
	var notBefore, notAfter time.Duration
	var err error

	if certOutPath, err = getStringFlag(generateCaCmd, generateCACertOutPath, logger); err != nil {
		return
	}
	if curveName, err = getStringFlag(generateCaCmd, generateCACurveName, logger); err != nil {
		return
	}
	if notBefore, err = getDurationFlag(generateCaCmd, generateCANotBeforeRelative, logger); err != nil {
		return
	}
	if notAfter, err = getDurationFlag(generateCaCmd, generateCANotAfterRelative, logger); err != nil {
		return
	}

	logger, _ := logging.CreateLogger()

	logger = logger.With(
		zap.String(generateCACurveName, curveName),
		zap.String(generateCACertOutPath, certOutPath),
	)

	generator := cert.NewDefaultGenerator(config.CertOptions{
		CertCachePath: certOutPath,
		Curve:         config.CurveType(curveName),
		Validity: config.ValidityByPurpose{
			CA: config.ValidityDuration{
				NotAfterRelative:  notAfter,
				NotBeforeRelative: notBefore,
			},
		},
	})

	var caCrt *tls.Certificate
	if caCrt, err = generator.CACert(caCertOptions); err != nil {
		logger.Error(
			"failed to generate CA certificate",
			zap.Error(err),
		)
		return
	}

	if len(caCrt.Certificate) < 1 {
		logger.Error("no public key given for generated CA certificate")
		return
	}

	var pubKey *x509.Certificate
	if pubKey, err = x509.ParseCertificate(caCrt.Certificate[0]); err != nil {
		logger.Error(
			"failed to parse public key from generated CA",
			zap.Error(err),
		)
		return
	}

	pemCrt := cert.NewPEM(caCrt)
	if err = pemCrt.Write(pubKey.Subject.CommonName, certOutPath); err != nil {
		logger.Error(
			"failed to write Ca files",
			zap.Error(err),
		)
	}
	logger.Info("completed certificate generation")
}

func getDurationFlag(cmd *cobra.Command, flagName string, logger logging.Logger) (val time.Duration, err error) {
	if val, err = cmd.Flags().GetDuration(flagName); err != nil {
		logger.Error(
			"failed to parse parse flag",
			zap.String("flag", flagName),
			zap.Error(err),
		)
	}
	return
}

func getStringFlag(cmd *cobra.Command, flagName string, logger logging.Logger) (val string, err error) {
	if val, err = cmd.Flags().GetString(flagName); err != nil {
		logger.Error(
			"failed to parse parse flag",
			zap.String("flag", flagName),
			zap.Error(err),
		)
	}
	return
}
