package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

const (
	generateCACertOutPath       = "cert-out"
	generateCAKeyOutPath        = "key-out"
	generateCACurveName         = "curve"
	generateCANotBeforeRelative = "not-before"
	generateCANotAfterRelative  = "not-after"
)

func generateCACmd(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-ca",
		Short: "Generate a new CA certificate and corresponding key",
		Long:  ``,
		Run:   runGenerateCA(logger),
	}

	cmd.Flags().String(generateCACertOutPath, "", "Path where CA cert file should be stored")
	cmd.Flags().String(generateCAKeyOutPath, "", "Path where CA key file should be stored")
	cmd.Flags().String(generateCACurveName, "", "Name of the curve to use, if empty ED25519 is used, other valid values are [P224, P256,P384,P521]")
	cmd.Flags().Duration(generateCANotBeforeRelative, 17520*time.Hour, "Relative time value since when in the past the CA certificate should be valid. The value has a time unit, the greatest time unit is h for hour.")
	cmd.Flags().Duration(generateCANotAfterRelative, 17520*time.Hour, "Relative time value until when in the future the CA certificate should be valid. The value has a time unit, the greatest time unit is h for hour.")

	return cmd
}

func getDurationFlag(cmd *cobra.Command, flagName string, logger *zap.Logger) (val time.Duration, err error) {
	if val, err = cmd.Flags().GetDuration(flagName); err != nil {
		logger.Error(
			"failed to parse parse flag",
			zap.String("flag", flagName),
			zap.Error(err),
		)
	}
	return
}

func getStringFlag(cmd *cobra.Command, flagName string, logger *zap.Logger) (val string, err error) {
	if val, err = cmd.Flags().GetString(flagName); err != nil {
		logger.Error(
			"failed to parse parse flag",
			zap.String("flag", flagName),
			zap.Error(err),
		)
	}
	return
}

func runGenerateCA(logger *zap.Logger) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var certOutPath, keyOutPath, curveName string
		var notBefore, notAfter time.Duration
		var err error

		if certOutPath, err = getStringFlag(cmd, generateCACertOutPath, logger); err != nil {
			return
		}

		if keyOutPath, err = getStringFlag(cmd, generateCAKeyOutPath, logger); err != nil {
			return
		}
		if curveName, err = getStringFlag(cmd, generateCACurveName, logger); err != nil {
			return
		}

		if notBefore, err = getDurationFlag(cmd, generateCANotBeforeRelative, logger); err != nil {
			return
		}

		if notAfter, err = getDurationFlag(cmd, generateCANotAfterRelative, logger); err != nil {
			return
		}

		logger := logger.With(
			zap.String(generateCACurveName, curveName),
			zap.String(generateCACertOutPath, certOutPath),
			zap.String(generateCAKeyOutPath, keyOutPath),
		)

		certStore := certStore{
			options: &tlsOptions{
				ecdsaCurve: curveType(curveName),
				validity: validity{
					ca: certValidity{
						notAfterRelative:  notAfter,
						notBeforeRelative: notBefore,
					},
				},
				rootCaCert: cert{
					publicKeyPath:  certOutPath,
					privateKeyPath: keyOutPath,
				},
			},
		}

		if _, _, err := certStore.generateCaCert(); err != nil {
			logger.Error(
				"failed to generate CA cert",
				zap.Error(err),
			)
		} else {
			logger.Info("Successfully generated CA cert")
		}
	}
}
