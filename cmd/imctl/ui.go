package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/CAFxX/httpcompression"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"gitlab.com/inetmock/inetmock/internal/ui/api"
	"gitlab.com/inetmock/inetmock/internal/ui/cert"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/web"
)

const (
	defaultListeningPort uint16 = 8732
)

var (
	host  string
	port  uint16
	uiCmd = &cobra.Command{
		Use:   "ui",
		Short: "Client UI to monitor inetmock activity",
	}
	uiServeCmd = &cobra.Command{
		Use:          "serve",
		Short:        "Start a local web server to serve the client UI",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI(web.WebFS)
		},
	}
	uiServeDevCmd = &cobra.Command{
		Use:          "serve-dev",
		Short:        "Start a local web server to serve the client UI",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wd, err := os.Getwd(); err != nil {
				return err
			} else {
				return runAPI(os.DirFS(filepath.Join(wd, "web")))
			}
		},
	}
)

func runAPI(webFS fs.FS) error {
	var socketURL *url.URL
	if u, err := url.Parse(cfg.SocketPath); err != nil {
		return err
	} else {
		socketURL = u
	}

	tlsCfg := new(tls.Config)
	certRequest := certRequestForHost(host)
	if srvCert, err := cert.GenerateServerCert(certRequest); err != nil {
		return err
	} else {
		tlsCfg.Certificates = append(tlsCfg.Certificates, *srvCert)
	}

	mux := http.NewServeMux()
	api.RegisterGRPCWebSocketProxy(cliApp.Context(), mux, socketURL, tlsCfg, cliApp.Logger())
	if err := api.RegisterViews(mux, webFS, cliApp.Logger().Named("views")); err != nil {
		return err
	}
	handler := api.RegisterStaticFileHandlingMiddleware(mux, webFS, cliApp.Logger().Named("static-files-middleware"))
	hostPort := fmt.Sprintf("%s:%d", host, port)

	return cmuxListen(cliApp.Context(), hostPort, tlsCfg, handler, cliApp.Logger())
}

func init() {
	uiServeCmd.Flags().StringVar(&host, "host", "127.0.0.1", "Host to bind web server on")
	uiServeCmd.Flags().Uint16Var(&port, "port", defaultListeningPort, "Port to bind the web server to")
	uiCmd.AddCommand(uiServeCmd, uiServeDevCmd)
}

func cmuxListen(ctx context.Context, hostPort string, tlsCfg *tls.Config, handler http.Handler, logger logging.Logger) error {
	var (
		errGrp, _                  = errgroup.WithContext(ctx)
		plainListener, tlsListener net.Listener
		compressorMiddleware       func(http.Handler) http.Handler
		err                        error
	)
	if plainListener, err = net.Listen("tcp", hostPort); err != nil {
		return err
	}

	if compressorMiddleware, err = httpcompression.DefaultAdapter(); err != nil {
		return err
	}

	muxer := cmux.New(plainListener)
	tlsListener = tls.NewListener(muxer.Match(cmux.TLS()), tlsCfg)
	plainListener = muxer.Match(cmux.Any())

	muxer.HandleError(func(err error) bool {
		logger.Error("error while serving cmux", zap.Error(err))
		return false
	})

	defer plainListener.Close()
	defer tlsListener.Close()

	defer muxer.Close()

	errGrp.Go(func() error {
		if err := muxer.Serve(); err != nil && !errors.Is(err, cmux.ErrServerClosed) {
			return err
		}
		return nil
	})

	errGrp.Go(func() error {
		if err := http.Serve(plainListener, handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	errGrp.Go(func() error {
		if err := http.Serve(tlsListener, compressorMiddleware(handler)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	return errGrp.Wait()
}

func certRequestForHost(host string) cert.Request {
	dnsNames := []string{"localhost"}
	if hostname, err := os.Hostname(); err == nil {
		dnsNames = append(dnsNames, hostname)
	}

	var allAddrs []net.IP
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipAddr, ok := addr.(*net.IPNet); ok {
				allAddrs = append(allAddrs, ipAddr.IP)
			}
		}
	}

	if parsedIP := net.ParseIP(host); parsedIP != nil {
		ipsToRequest := []net.IP{parsedIP}
		if parsedIP.IsUnspecified() {
			ipsToRequest = allAddrs
		}
		return cert.Request{
			DNSNames:    dnsNames,
			IPAddresses: ipsToRequest,
		}
	} else {
		dnsNames = append(dnsNames, host)
		return cert.Request{
			DNSNames:    dnsNames,
			IPAddresses: allAddrs,
		}
	}
}
