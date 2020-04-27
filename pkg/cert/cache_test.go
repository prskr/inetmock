package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	certmock "github.com/baez90/inetmock/internal/mock/cert"
	config2 "github.com/baez90/inetmock/pkg/config"
	"github.com/golang/mock/gomock"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		timeSource TimeSource
		cert       *x509.Certificate
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Expect cert should not be renewed right after creation",
			args: args{
				timeSource: func() TimeSource {
					tsMock := certmock.NewMockTimeSource(ctrl)
					tsMock.
						EXPECT().
						UTCNow().
						Return(time.Now().UTC()).
						Times(1)
					return tsMock
				}(),
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
				timeSource: func() TimeSource {
					tsMock := certmock.NewMockTimeSource(ctrl)
					tsMock.
						EXPECT().
						UTCNow().
						Return(time.Now().UTC().Add(serverRelativeValidity/2 + 1*time.Hour)).
						Times(1)
					return tsMock
				}(),
				cert: &x509.Certificate{
					NotAfter:  time.Now().UTC().Add(serverRelativeValidity),
					NotBefore: time.Now().UTC().Add(-serverRelativeValidity),
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := certShouldBeRenewed(tt.args.timeSource, tt.args.cert); got != tt.want {
				t.Errorf("certShouldBeRenewed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileSystemCache_Get(t *testing.T) {
	type fields struct {
		certCachePath string
		inMemCache    map[string]*tls.Certificate
		timeSource    TimeSource
	}
	type args struct {
		cn string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantOk bool
	}{
		{
			name: "Get a miss when no cert is present",
			fields: fields{
				certCachePath: os.TempDir(),
				inMemCache:    make(map[string]*tls.Certificate),
				timeSource:    NewTimeSource(),
			},
			args: args{
				cnLocalhost,
			},
			wantOk: false,
		},
		{
			name: "Get a prepared certificate from the memory cache",
			fields: func() fields {
				certGen := setupCertGen()

				caCrt, _ := certGen.CACert(GenerationOptions{
					CommonName: caCN,
				})

				srvCrt, _ := certGen.ServerCert(GenerationOptions{
					CommonName: cnLocalhost,
				}, caCrt)

				return fields{
					certCachePath: os.TempDir(),
					inMemCache: map[string]*tls.Certificate{
						cnLocalhost: srvCrt,
					},
					timeSource: NewTimeSource(),
				}
			}(),
			args: args{
				cn: cnLocalhost,
			},
			wantOk: true,
		},
		{
			name: "Get a prepared certificate from the file system",
			fields: func() fields {
				certGen := setupCertGen()

				caCrt, _ := certGen.CACert(GenerationOptions{
					CommonName: "INetMock",
				})

				srvCrt, _ := certGen.ServerCert(GenerationOptions{
					CommonName: serverCN,
				}, caCrt)

				pem := NewPEM(srvCrt)
				if err := pem.Write(serverCN, os.TempDir()); err != nil {
					panic(err)
				}

				t.Cleanup(func() {
					for _, f := range []string{
						path.Join(os.TempDir(), fmt.Sprintf("%s.pem", serverCN)),
						path.Join(os.TempDir(), fmt.Sprintf("%s.key", serverCN)),
					} {
						_ = os.Remove(f)
					}
				})

				return fields{
					certCachePath: os.TempDir(),
					inMemCache:    make(map[string]*tls.Certificate),
					timeSource:    NewTimeSource(),
				}
			}(),
			args: args{
				cn: serverCN,
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileSystemCache{
				certCachePath: tt.fields.certCachePath,
				inMemCache:    tt.fields.inMemCache,
				timeSource:    tt.fields.timeSource,
			}
			gotCrt, gotOk := f.Get(tt.args.cn)

			if gotOk && (gotCrt == nil || !reflect.DeepEqual(reflect.TypeOf(new(tls.Certificate)), reflect.TypeOf(gotCrt))) {
				t.Errorf("Wanted propert certificate but got %v", gotCrt)
			}

			if gotOk != tt.wantOk {
				t.Errorf("Get() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_fileSystemCache_Put(t *testing.T) {
	type fields struct {
		certCachePath string
		inMemCache    map[string]*tls.Certificate
		timeSource    TimeSource
	}
	type args struct {
		cert *tls.Certificate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	return NewDefaultGenerator(config2.CertOptions{
		Validity: config2.ValidityByPurpose{
			Server: config2.ValidityDuration{
				NotBeforeRelative: serverRelativeValidity,
				NotAfterRelative:  serverRelativeValidity,
			},
			CA: config2.ValidityDuration{
				NotBeforeRelative: caRelativeValidity,
				NotAfterRelative:  caRelativeValidity,
			},
		},
	})
}
