package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"

	certmock "gitlab.com/inetmock/inetmock/internal/mock/cert"
)

const (
	cnLocalhost            = "localhost"
	caCN                   = "UnitTests"
	serverRelativeValidity = 24 * time.Hour
	caRelativeValidity     = 168 * time.Hour
)

var (
	serverCN = fmt.Sprintf("%s-%d", cnLocalhost, time.Now().Unix())
)

func Test_certShouldBeRenewed(t *testing.T) {
	t.Parallel()
	type args struct {
		timeSourceSetup func(ctrl *gomock.Controller) TimeSource
		cert            *x509.Certificate
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{
		{
			name: "Expect cert should not be renewed right after creation",
			args: args{
				timeSourceSetup: func(ctrl *gomock.Controller) TimeSource {
					tsMock := certmock.NewMockTimeSource(ctrl)
					tsMock.
						EXPECT().
						UTCNow().
						Return(time.Now().UTC()).
						Times(1)
					return tsMock
				},
				cert: &x509.Certificate{
					NotAfter:  time.Now().UTC().Add(serverRelativeValidity),
					NotBefore: time.Now().UTC().Add(-serverRelativeValidity),
				},
			},
			want: false,
		},
		{
			name: "Expect cert should be renewed if the remaining lifetime is less than a quarter of the total lifetime",
			args: args{
				timeSourceSetup: func(ctrl *gomock.Controller) TimeSource {
					tsMock := certmock.NewMockTimeSource(ctrl)
					tsMock.
						EXPECT().
						UTCNow().
						Return(time.Now().UTC().Add(serverRelativeValidity/2 + 1*time.Hour)).
						Times(1)
					return tsMock
				},
				cert: &x509.Certificate{
					NotAfter:  time.Now().UTC().Add(serverRelativeValidity),
					NotBefore: time.Now().UTC().Add(-serverRelativeValidity),
				},
			},
			want: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			if got := certShouldBeRenewed(tt.args.timeSourceSetup(ctrl), tt.args.cert); got != tt.want {
				t.Errorf("certShouldBeRenewed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileSystemCache_Get(t *testing.T) {
	t.Parallel()

	certGen := setupCertGen()

	caCrt, _ := certGen.CACert(GenerationOptions{
		CommonName: caCN,
	})

	srvCrt, _ := certGen.ServerCert(GenerationOptions{
		CommonName: cnLocalhost,
	}, caCrt)

	type fields struct {
		inMemCache map[string]*tls.Certificate
		timeSource TimeSource
	}
	type args struct {
		cn string
	}
	type testCase struct {
		name             string
		fields           fields
		polluteCertCache bool
		args             args
		wantOk           bool
	}
	tests := []testCase{
		{
			name: "Get a miss when no cert is present",
			fields: fields{
				inMemCache: make(map[string]*tls.Certificate),
				timeSource: NewTimeSource(),
			},
			args: args{
				cnLocalhost,
			},
			wantOk: false,
		},
		{
			name: "Get a prepared certificate from the memory cache",
			fields: fields{
				inMemCache: map[string]*tls.Certificate{
					cnLocalhost: srvCrt,
				},
				timeSource: NewTimeSource(),
			},
			args: args{
				cn: cnLocalhost,
			},
			wantOk: true,
		},
		{
			name: "Get a prepared certificate from the file system",
			fields: fields{
				inMemCache: make(map[string]*tls.Certificate),
				timeSource: NewTimeSource(),
			},
			args: args{
				cn: serverCN,
			},
			polluteCertCache: true,
			wantOk:           true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			f := &fileSystemCache{
				certCachePath: dir,
				inMemCache:    tt.fields.inMemCache,
				timeSource:    tt.fields.timeSource,
			}

			if tt.polluteCertCache {
				pem := NewPEM(srvCrt)
				if err := pem.Write(serverCN, dir); err != nil {
					t.Fatalf("polluteCertCache error = %v", err)
				}
			}

			gotCrt, gotOk := f.Get(tt.args.cn)

			if gotOk && (gotCrt == nil || !td.CmpIsa(t, gotCrt, new(tls.Certificate))) {
				t.Errorf("Wanted propert certificate but got %v", gotCrt)
			}

			if gotOk != tt.wantOk {
				t.Errorf("Get() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_fileSystemCache_Put(t *testing.T) {
	t.Parallel()
	type fields struct {
		certCachePath string
		inMemCache    map[string]*tls.Certificate
		timeSource    TimeSource
	}
	type args struct {
		cert *tls.Certificate
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	tests := []testCase{
		{
			name: "Want error if nil cert is passed to put",
			fields: fields{
				certCachePath: os.TempDir(),
				inMemCache:    make(map[string]*tls.Certificate),
				timeSource:    NewTimeSource(),
			},
			args: args{
				cert: nil,
			},
			wantErr: true,
		},
		{
			name: "Want error if empty cert is passed to put",
			fields: fields{
				certCachePath: os.TempDir(),
				inMemCache:    make(map[string]*tls.Certificate),
				timeSource:    NewTimeSource(),
			},
			args: args{
				cert: &tls.Certificate{},
			},
			wantErr: true,
		},
		{
			name: "No error if valid cert is passed",
			fields: fields{
				certCachePath: os.TempDir(),
				inMemCache:    make(map[string]*tls.Certificate),
				timeSource:    NewTimeSource(),
			},
			args: args{
				cert: func() *tls.Certificate {
					gen := setupCertGen()
					ca, _ := gen.CACert(GenerationOptions{
						CommonName: caCN,
					})

					srvCN := fmt.Sprintf("%s-%d", cnLocalhost, time.Now().Unix())

					t.Cleanup(func() {
						for _, f := range []string{
							path.Join(os.TempDir(), fmt.Sprintf("%s.pem", srvCN)),
							path.Join(os.TempDir(), fmt.Sprintf("%s.key", srvCN)),
						} {
							_ = os.Remove(f)
						}
					})

					srvCrt, _ := gen.ServerCert(GenerationOptions{
						CommonName: srvCN,
					}, ca)
					return srvCrt
				}(),
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &fileSystemCache{
				certCachePath: tt.fields.certCachePath,
				inMemCache:    tt.fields.inMemCache,
				timeSource:    tt.fields.timeSource,
			}
			if err := f.Put(tt.args.cert); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func setupCertGen() Generator {
	return NewDefaultGenerator(Options{
		Validity: ValidityByPurpose{
			Server: ValidityDuration{
				NotBeforeRelative: serverRelativeValidity,
				NotAfterRelative:  serverRelativeValidity,
			},
			CA: ValidityDuration{
				NotBeforeRelative: caRelativeValidity,
				NotAfterRelative:  caRelativeValidity,
			},
		},
	})
}
