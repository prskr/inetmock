package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/health"
)

var (
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Run various commands for checks e.g. run checks, parse checks for validation,...",
	}

	runCheckCmd = &cobra.Command{
		Use:     "run",
		Short:   "Run a check script",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"exec"},
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument, got %d", len(args))
			}
			return runCheck(args[0])
		},
		SilenceUsage: true,
	}

	runCheckArgs = &struct {
		TargetIP      net.IP
		HTTPPort      uint16
		HTTPSPort     uint16
		Timeout       time.Duration
		TLSSkipVerify bool
		CACertPath    string
	}{}
)

//nolint:gomnd // 127.0.0.1 is well known
func init() {
	runCheckCmd.Flags().IPVar(&runCheckArgs.TargetIP, "target-ip", net.IPv4(127, 0, 0, 1), "target IP used to connect for the check execution")
	runCheckCmd.Flags().Uint16Var(&runCheckArgs.HTTPPort, "http-port", 80, "Port to connect to for 'http://' requests")
	runCheckCmd.Flags().Uint16Var(&runCheckArgs.HTTPSPort, "https-port", 443, "Port to connect to for 'https://' requests")
	runCheckCmd.Flags().DurationVar(&runCheckArgs.Timeout, "check-timeout", 1*time.Second, "timeout to execute the check")
	runCheckCmd.Flags().StringVar(&runCheckArgs.CACertPath, "ca-cert", "", "Path to CA cert file to trust additionally to system cert pool")
	runCheckCmd.Flags().BoolVarP(&runCheckArgs.TLSSkipVerify, "insecure", "i", false, "Skip TLS server certificate verification")
	checkCmd.AddCommand(runCheckCmd)
}

func runCheck(script string) error {
	var healthCfg = health.Config{
		Client: health.HTTPClientConfig{
			HTTP: health.Server{
				IP:   runCheckArgs.TargetIP.String(),
				Port: runCheckArgs.HTTPPort,
			},
			HTTPS: health.Server{
				IP:   runCheckArgs.TargetIP.String(),
				Port: runCheckArgs.HTTPSPort,
			},
		},
	}

	var certPool *x509.CertPool

	switch strings.ToLower(runtime.GOOS) {
	case "linux", "darwin", "freebsd", "netbsd", "openbsd", "solaris":
		var err error
		if certPool, err = x509.SystemCertPool(); err != nil {
			return err
		}
	default:
		certPool = x509.NewCertPool()
	}

	if runCheckArgs.CACertPath != "" {
		if err := addCACertToPool(certPool); err != nil {
			cliApp.Logger().Warn("failed to load CA cert", zap.Error(err))
		}
	}

	var tlsConfig = &tls.Config{
		RootCAs: certPool,
		//nolint:gosec
		InsecureSkipVerify: runCheckArgs.TLSSkipVerify,
	}

	var client = health.HTTPClient(healthCfg, tlsConfig)

	check := new(rules.Check)
	if err := rules.Parse(script, check); err != nil {
		return err
	}

	if compiledCheck, err := health.NewHTTPRuleCheck("CLI", client, cliApp.Logger().Named("check"), check); err != nil {
		return err
	} else {
		ctx, cancel := context.WithTimeout(cliApp.Context(), runCheckArgs.Timeout)
		defer cancel()
		if err := compiledCheck.Status(ctx); err != nil {
			return err
		}
	}
	cliApp.Logger().Info("Successfully executed")
	return nil
}

func addCACertToPool(pool *x509.CertPool) (err error) {
	var buffer = bytes.NewBuffer(nil)
	var reader io.ReadCloser
	if reader, err = os.Open(runCheckArgs.CACertPath); err != nil {
		return err
	}

	defer func() {
		err = multierr.Append(err, reader.Close())
	}()

	if _, err = io.Copy(buffer, reader); err != nil {
		return err
	}

	pool.AppendCertsFromPEM(buffer.Bytes())
	return nil
}
