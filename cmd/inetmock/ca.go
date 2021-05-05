package main

import (
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/cert"
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
	defaultValidityOffset            = 17520 * time.Hour
)

var (
	generateCaCmd          *cobra.Command
	caCertOptions          cert.GenerationOptions
	notBefore, notAfter    time.Duration
	certOutPath, curveName string
)

//nolint:lll
func init() {
	generateCaCmd = &cobra.Command{
		Use:          "generate-ca",
		Short:        "Generate a new CA certificate and corresponding key",
		Long:         ``,
		Run:          runGenerateCA,
		SilenceUsage: true,
	}

	generateCaCmd.Flags().StringVar(&caCertOptions.CommonName, generateCACommonName, "INetMock", "Certificate Common Name that will also be used as file name during generation.")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.Organization, generateCaOrganizationName, nil, "Organization information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.OrganizationalUnit, generateCaOrganizationalUnitName, nil, "Organizational unit information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.Country, generateCaCountryName, nil, "Country information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.Province, generateCaStateName, nil, "State information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.Locality, generateCaLocalityName, nil, "Locality information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.StreetAddress, generateCaStreetAddressName, nil, "Street address information to append to certificate")
	generateCaCmd.Flags().StringSliceVar(&caCertOptions.PostalCode, generateCaPostalCodeName, nil, "Postal code information to append to certificate")
	generateCaCmd.Flags().StringVar(&certOutPath, generateCACertOutPath, "", "Path where CA files should be stored")
	generateCaCmd.Flags().StringVar(&curveName, generateCACurveName, "", "Name of the curve to use, if empty ED25519 is used, other valid values are [P224, P256,P384,P521]")
	generateCaCmd.Flags().DurationVar(&notBefore, generateCANotBeforeRelative, defaultValidityOffset, "Relative time value since when in the past the CA certificate should be valid. The value has a time unit, the greatest time unit is h for hour.")
	generateCaCmd.Flags().DurationVar(&notAfter, generateCANotAfterRelative, defaultValidityOffset, "Relative time value until when in the future the CA certificate should be valid. The value has a time unit, the greatest time unit is h for hour.")
}

func runGenerateCA(_ *cobra.Command, _ []string) {
	logger := serverApp.Logger().Named("generate-ca")

	logger = logger.With(
		zap.String(generateCACurveName, curveName),
		zap.String(generateCACertOutPath, certOutPath),
	)

	generator := cert.NewDefaultGenerator(cert.Options{
		CertCachePath: certOutPath,
		Curve:         cert.CurveType(curveName),
		Validity: cert.ValidityByPurpose{
			CA: cert.ValidityDuration{
				NotAfterRelative:  notAfter,
				NotBeforeRelative: notBefore,
			},
		},
	})

	var err error
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
