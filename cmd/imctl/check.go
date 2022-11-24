package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/valyala/bytebufferpool"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
	"inetmock.icb4dc0.de/inetmock/pkg/health"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

type runCheckArgs struct {
	Target        string
	HTTPPort      uint16
	HTTPSPort     uint16
	DNSPort       uint16
	DNSProto      string
	DoTPort       uint16
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

			clientSetupCtx, cancel := context.WithTimeout(cliApp.Context(), args.Timeout)
			httpClient, dnsResolver, err := setupClients(clientSetupCtx, cliApp.Logger(), args)
			cancel()
			if err != nil {
				return err
			}

			switch len(params) {
			case noArgs:
				stdinReader := bufio.NewReader(os.Stdin)
				scriptBuffer := strings.Builder{}
				for {
					if line, err := stdinReader.ReadString('\n'); err != nil {
						if errors.Is(err, io.EOF) {
							return runCheck(cliApp.Context(), cliApp.Logger(), scriptBuffer.String(), httpClient, dnsResolver)
						}
						return err
					} else {
						scriptBuffer.WriteString(line)
					}
				}
			case singleArg:
				return runCheck(cliApp.Context(), cliApp.Logger(), params[0], httpClient, dnsResolver)
			default:
				return errors.New("missing script")
			}
		},
		SilenceUsage: true,
	}

	args = new(runCheckArgs)
)

//nolint:lll // still better readable than breaking these lines
func init() {
	const (
		defaultHTTPPort  int = 80
		defaultHTTPSPort     = 443
		defaultDNSPort       = 53
		defaultDoTPort       = 853
		defaultTimeout       = 5 * time.Second
	)

	runCheckCmd.Flags().StringVar(&args.Target, "target", "localhost", "target IP used to connect for the check execution")
	runCheckCmd.Flags().Uint16Var(&args.HTTPPort, "http-port", uint16(defaultHTTPPort), "Port to connect to for 'http://' requests")
	runCheckCmd.Flags().Uint16Var(&args.HTTPSPort, "https-port", defaultHTTPSPort, "Port to connect to for 'https://' requests")
	runCheckCmd.Flags().Uint16Var(&args.DNSPort, "dns-port", defaultDNSPort, "Port to connect to for DNS requests")
	runCheckCmd.Flags().StringVar(&args.DNSProto, "dns-proto", "tcp", "Protocol to use for DNS requests one of [tcp, tcp4, tcp6, udp, udp4, udp6]")
	runCheckCmd.Flags().Uint16Var(&args.DoTPort, "dot-port", defaultDoTPort, "Port to use for DoT requests")
	runCheckCmd.Flags().DurationVar(&args.Timeout, "check-timeout", defaultTimeout, "timeout to execute the check")
	runCheckCmd.Flags().StringVar(&args.CACertPath, "ca-cert", "", "Path to CA cert file to trust additionally to system cert pool")
	runCheckCmd.Flags().BoolVarP(&args.TLSSkipVerify, "insecure", "i", false, "Skip TLS server certificate verification")
	checkCmd.AddCommand(runCheckCmd)
}

func runCheck(
	ctx context.Context,
	logger logging.Logger,
	script string,
	httpClients health.HTTPClientForModule,
	dnsResolvers health.ResolverForModule,
) error {
	checkLogger := logger.Named("check")
	checkScript, err := rules.Parse[rules.CheckScript](script)
	if err != nil {
		return err
	}

	compiledChecks := make([]health.Check, 0, len(checkScript.Checks))

	for idx := range checkScript.Checks {
		check := checkScript.Checks[idx]
		switch module := strings.ToLower(check.Initiator.Module); module {
		case "http", "http2":
			if compiledCheck, err := health.NewHTTPRuleCheck("CLI", httpClients, checkLogger, &check); err != nil {
				return err
			} else {
				compiledChecks = append(compiledChecks, compiledCheck)
			}
		case "dns", "dot", "doh", "doh2":
			if compiledCheck, err := health.NewDNSRuleCheck("CLI", dnsResolvers, checkLogger, &check); err != nil {
				return err
			} else {
				compiledChecks = append(compiledChecks, compiledCheck)
			}
		default:
			return fmt.Errorf("unmatched check module: %s", module)
		}
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, args.Timeout)
	defer cancel()
	grp, grpcCtx := errgroup.WithContext(timeoutCtx)

	for idx := range compiledChecks {
		idx := idx
		grp.Go(func() error {
			if err := compiledChecks[idx].Status(grpcCtx); err != nil {
				checkLogger.Error("Failed to execute check", zap.Error(err))
				return err
			}
			return nil
		})
	}

	if err := grp.Wait(); err != nil {
		checkLogger.Error("Failed to execute check", zap.Error(err))
		return err
	}

	checkLogger.Info("Successfully executed")
	return nil
}

func setupClients(
	ctx context.Context,
	logger logging.Logger,
	args *runCheckArgs,
) (health.HTTPClientForModule, health.ResolverForModule, error) {
	target, err := resolveTargetIP(ctx, logger, args.Target)
	if err != nil {
		return nil, nil, err
	}
	healthCfg := health.Config{
		Client: health.ClientsConfig{
			HTTP: health.Server{
				IP:   target,
				Port: args.HTTPPort,
			},
			HTTPS: health.Server{
				IP:   target,
				Port: args.HTTPSPort,
			},
			DNS: health.Server{
				IP:   target,
				Port: args.DNSPort,
			},
			DoT: health.Server{
				IP:   target,
				Port: args.DoTPort,
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
			return nil, nil, err
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

	return health.HTTPClients(healthCfg, tlsConfig), health.Resolvers(healthCfg, tlsConfig), nil
}

func addCACertToPool(pool *x509.CertPool) (err error) {
	buffer := bytebufferpool.Get()
	defer bytebufferpool.Put(buffer)

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

func resolveTargetIP(ctx context.Context, logger logging.Logger, target string) (string, error) {
	logger = logger.With(zap.String("target", target))
	if targetIP := net.ParseIP(target); len(targetIP) >= net.IPv4len {
		logger.Debug("target is an IP address - will be set for clients")
		return target, nil
	}

	logger.Debug("target is apparently not an IP - resolving IP address")
	if addrs, err := net.DefaultResolver.LookupHost(ctx, args.Target); err != nil {
		return "", err
	} else {
		logger.Debug("Resolved target addresses", zap.Strings("resolvedAddresses", addrs))
		//nolint:gosec // no need for cryptographic security when picking a random IP address to contact
		pickedAddr := addrs[rand.Intn(len(addrs))]
		logger.Debug("Picked random address", zap.String("newTargetAddress", pickedAddr))
		if parsed := net.ParseIP(pickedAddr); parsed == nil {
			logger.Error("Could not parse resolved IP", zap.String("addr", pickedAddr))
			return "", errors.New("could not parse resolved IP")
		} else if parsed.To4() == nil {
			pickedAddr = fmt.Sprintf("[%s]", pickedAddr)
		}
		args.Target = pickedAddr
	}
	return target, nil
}
