package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	generateCACertOutPath = "cert-out"
	generateCAKeyOutPath  = "key-out"
	generateCACurveName   = "curve"
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

	return cmd
}

func runGenerateCA(logger *zap.Logger) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var certOutPath, keyOutPath, curveName string
		var err error
		if certOutPath, err = cmd.Flags().GetString(generateCACertOutPath); err != nil {
			logger.Error(
				"failed to parse parse flag",
				zap.String("flag", generateCACertOutPath),
				zap.Error(err),
			)
			return
		}
		if keyOutPath, err = cmd.Flags().GetString(generateCAKeyOutPath); err != nil {
			logger.Error(
				"failed to parse parse flag",
				zap.String("flag", generateCAKeyOutPath),
				zap.Error(err),
			)
			return
		}
		if curveName, err = cmd.Flags().GetString(generateCACurveName); err != nil {
			logger.Error(
				"failed to parse parse flag",
				zap.String("flag", generateCACurveName),
				zap.Error(err),
			)
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
