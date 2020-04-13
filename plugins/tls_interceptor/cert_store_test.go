package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func Test_generateCaCert(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "*-inetmock")
	if err != nil {
		t.Errorf("failed to create temp dir %v", err)
		return
	}

	options := &tlsOptions{
		ecdsaCurve: "P256",
		rootCaCert: cert{
			publicKeyPath:  filepath.Join(tmpDir, "localhost.pem"),
			privateKeyPath: filepath.Join(tmpDir, "localhost.key"),
		},
		validity: validity{
			ca: certValidity{
				notBeforeRelative: time.Hour * 24 * 30,
				notAfterRelative:  time.Hour * 24 * 30,
			},
		},
	}

	certStore := certStore{
		options: options,
	}

	defer func() {
		_ = os.Remove(tmpDir)
	}()

	_, _, err = certStore.generateCaCert()

	if err != nil {
		t.Errorf("failed to generate CA cert %v", err)
	}

	if _, err = os.Stat(options.rootCaCert.publicKeyPath); err != nil {
		t.Errorf("cert file was not created")
	}

	if _, err = os.Stat(options.rootCaCert.privateKeyPath); err != nil {
		t.Errorf("cert file was not created")
	}
}

func Test_generateDomainCert(t *testing.T) {

	tmpDir, err := ioutil.TempDir(os.TempDir(), "*-inetmock")
	if err != nil {
		t.Errorf("failed to create temp dir %v", err)
		return
	}
	defer func() {
		_ = os.Remove(tmpDir)
	}()

	caTlsCert, _ := loadPEMCert(testCaCrt, testCaKey)
	caCert, _ := x509.ParseCertificate(caTlsCert.Certificate[0])

	options := &tlsOptions{
		ecdsaCurve:    "P256",
		certCachePath: tmpDir,
		validity: validity{
			domain: certValidity{
				notAfterRelative:  time.Hour * 24 * 30,
				notBeforeRelative: time.Hour * 24 * 30,
			},
			ca: certValidity{
				notAfterRelative:  time.Hour * 24 * 30,
				notBeforeRelative: time.Hour * 24 * 30,
			},
		},
	}

	zapLogger, _ := zap.NewDevelopment()
	logger := logging.NewLogger(zapLogger)

	certStore := certStore{
		options:      options,
		caCert:       caCert,
		logger:       logger,
		caPrivateKey: caTlsCert.PrivateKey,
	}

	type args struct {
		domain string
		ip     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test create google.com cert",
			args: args{
				domain: "google.com",
				ip:     "127.0.0.1",
			},
			wantErr: false,
		},
		{
			name: "Test create golem.de cert",
			args: args{
				domain: "golem.de",
				ip:     "127.0.0.1",
			},
			wantErr: false,
		},
		{
			name: "Test create golem.de cert with any IP address",
			args: args{
				domain: "golem.de",
				ip:     "10.10.0.10",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if domainCert, err := certStore.generateDomainCert(
				tt.args.domain,
				tt.args.ip,
			); (err != nil) != tt.wantErr || reflect.DeepEqual(domainCert, tls.Certificate{}) {
				t.Errorf("generateDomainCert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_certStore_initCaCert(t *testing.T) {

	tmpDir, err := ioutil.TempDir(os.TempDir(), "*-inetmock")
	if err != nil {
		t.Errorf("failed to create temp dir %v", err)
		return
	}
	defer func() {
		_ = os.Remove(tmpDir)
	}()

	caCertPath := filepath.Join(tmpDir, "cacert.pem")
	caKeyPath := filepath.Join(tmpDir, "cacert.key")

	if err := ioutil.WriteFile(caCertPath, testCaCrt, 0600); err != nil {
		t.Errorf("failed to write cacert.pem %v", err)
		return
	}
	if err := ioutil.WriteFile(caKeyPath, testCaKey, 0600); err != nil {
		t.Errorf("failed to write cacert.key %v", err)
		return
	}

	type fields struct {
		options *tlsOptions
		caCert  *x509.Certificate
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Init CA cert from file",
			wantErr: false,
			fields: fields{
				options: &tlsOptions{
					rootCaCert: cert{
						publicKeyPath:  caCertPath,
						privateKeyPath: caKeyPath,
					},
				},
			},
		},
		{
			name:    "Init CA with new cert",
			wantErr: false,
			fields: fields{
				options: &tlsOptions{
					rootCaCert: cert{
						publicKeyPath:  filepath.Join(tmpDir, "nonexistent.pem"),
						privateKeyPath: filepath.Join(tmpDir, "nonexistent.key"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &certStore{
				options: tt.fields.options,
				caCert:  tt.fields.caCert,
			}
			if err := cs.initCaCert(); (err != nil) != tt.wantErr || cs.caCert == nil {
				t.Errorf("initCaCert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
