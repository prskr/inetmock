package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
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

type runCheckArgs struct {
	TargetIP      net.IP
	HTTPPort      uint16
	HTTPSPort     uint16
	DNSPort       uint16
	DNSProto      string
	Timeout       time.Duration
	TLSSkipVerify bool
	CACertPath    string
}

var (
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Run various commands for checks e.g. run checks, parse checks for validation,...",
	}

	runCheckCmd = &cobra.Command{
		Use:     "run",
		Short:   "Run a check script",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"exec"},
		RunE: func(_ *cobra.Command, params []string) error {
			const (
				noArgs = iota
				singleArg
			)
			if len(params) > 1 {
				return fmt.Errorf("expected 1 argument, got %d", len(params))
			}

			switch len(params) {
			case noArgs:
				stdinReader := bufio.NewReader(os.Stdin)
				script := make([]string, 0)
				for {
					if line, err := stdinReader.ReadString('\n'); err != nil {
						if errors.Is(err, io.EOF) {
							return runCheck(script, args)
						}
						return err
					} else {
						script = append(script, line)
					}
				}
			case singleArg:
				return runCheck(params, args)
			default:
				return errors.New("missing script")
			}
		},
		SilenceUsage: true,
	}

	args = new(runCheckArgs)
)

// nolint:lll // still better readable than breaking these lines
func init() {
	const (
		defaultHTTPPort  int = 80
		defaultHTTPSPort     = 443
		defaultDNSPort       = 53
		defaultTimeout       = 1 * time.Second
	)

	//nolint:gomnd // 127.0.0.1 is well known
	runCheckCmd.Flags().IPVar(&args.TargetIP, "target-ip", net.IPv4(127, 0, 0, 1), "target IP used to connect for the check execution")
	runCheckCmd.Flags().Uint16Var(&args.HTTPPort, "http-port", uint16(defaultHTTPPort), "Port to connect to for 'http://' requests")
	runCheckCmd.Flags().Uint16Var(&args.HTTPSPort, "https-port", defaultHTTPSPort, "Port to connect to for 'https://' requests")
	runCheckCmd.Flags().Uint16Var(&args.DNSPort, "dns-port", defaultDNSPort, "Port to connect to for DNS requests")
	runCheckCmd.Flags().StringVar(&args.DNSProto, "dns-proto", "udp", "Protocol to use for DNS requests one of [tcp, tcp4, tcp6, udp, udp4, udp6]")
	runCheckCmd.Flags().DurationVar(&args.Timeout, "check-timeout", defaultTimeout, "timeout to execute the check")
	runCheckCmd.Flags().StringVar(&args.CACertPath, "ca-cert", "", "Path to CA cert file to trust additionally to system cert pool")
	runCheckCmd.Flags().BoolVarP(&args.TLSSkipVerify, "insecure", "i", false, "Skip TLS server certificate verification")
	checkCmd.AddCommand(runCheckCmd)
}

func runCheck(script []string, args *runCheckArgs) error {
	healthCfg := health.Config{
		Client: health.ClientsConfig{
			HTTP: health.Server{
				IP:   args.TargetIP.String(),
				Port: args.HTTPPort,
			},
			HTTPS: health.Server{
				IP:   args.TargetIP.String(),
				Port: args.HTTPSPort,
			},
			DNS: health.Server{
				IP:   args.TargetIP.String(),
				Port: args.DNSPort,
			},
		},
	}

	switch proto := strings.ToLower(args.DNSProto); proto {
	case "tcp", "tcp4", "tcp6":
		healthCfg.Client.DNS.Proto = proto
	case "udp", "udp4", "udp6":
		fallthrough
	default:
		healthCfg.Client.DNS.Proto = proto
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

	if args.CACertPath != "" {
		if err := addCACertToPool(certPool); err != nil {
			cliApp.Logger().Warn("failed to load CA cert", zap.Error(err))
		}
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
		//nolint:gosec
		InsecureSkipVerify: args.TLSSkipVerify,
	}

	client := health.HTTPClient(healthCfg, tlsConfig)
	resolver := health.DNSResolver(healthCfg)
	check := new(rules.Check)
	checkLogger := cliApp.Logger().Named("check")

	for idx := range script {
		rawRule := script[idx]

		if err := rules.Parse(rawRule, check); err != nil {
			return err
		}

		var (
			compiledCheck health.Check
			err           error
		)

		switch module := strings.ToLower(check.Initiator.Module); module {
		case "http":
			if compiledCheck, err = health.NewHTTPRuleCheck("CLI", client, checkLogger, check); err != nil {
				return err
			}
		case "dns":
			if compiledCheck, err = health.NewDNSRuleCheck("CLI", resolver, checkLogger, check); err != nil {
				return err
			}
		}

		ctx, cancel := context.WithTimeout(cliApp.Context(), args.Timeout)
		err = compiledCheck.Status(ctx)
		cancel()
		if err != nil {
			return err
		}
	}
	cliApp.Logger().Info("Successfully executed")
	return nil
}

func addCACertToPool(pool *x509.CertPool) (err error) {
	buffer := bytes.NewBuffer(nil)
	var reader io.ReadCloser
	if reader, err = os.Open(args.CACertPath); err != nil {
		return err
	}

	defer multierr.AppendInvoke(&err, multierr.Close(reader))

	if _, err = io.Copy(buffer, reader); err != nil {
		return err
	}

	pool.AppendCertsFromPEM(buffer.Bytes())
	return nil
}
