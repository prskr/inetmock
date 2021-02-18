package cert

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func readViper(cfg string) *viper.Viper {
	vpr := viper.New()
	vpr.SetConfigType("yaml")
	if err := vpr.ReadConfig(strings.NewReader(cfg)); err != nil {
		panic(err)
	}
	return vpr
}

//nolint:funlen
func Test_loadFromConfig(t *testing.T) {
	type args struct {
		config *viper.Viper
	}
	tests := []struct {
		name    string
		args    args
		want    Options
		wantErr bool
	}{
		{
			name:    "Parse valid TLS configuration",
			wantErr: false,
			args: args{
				config: readViper(`
tls:
  ecdsaCurve: P256
  validity:
    ca:
      notBeforeRelative: 17520h
      notAfterRelative: 17520h
    server:
      NotBeforeRelative: 168h
      NotAfterRelative: 168h
  rootCaCert:
    publicKey: ./ca.pem
    privateKey: ./ca.key
  certCachePath: /tmp/inetmock/
`),
			},
			want: Options{
				RootCACert: File{
					PublicKeyPath:  "./ca.pem",
					PrivateKeyPath: "./ca.key",
				},
				CertCachePath: "/tmp/inetmock/",
				Curve:         CurveTypeP256,
				Validity: ValidityByPurpose{
					CA: ValidityDuration{
						NotBeforeRelative: 17520 * time.Hour,
						NotAfterRelative:  17520 * time.Hour,
					},
					Server: ValidityDuration{
						NotBeforeRelative: 168 * time.Hour,
						NotAfterRelative:  168 * time.Hour,
					},
				},
			},
		},
		{
			name: "Get an error if CA public key path is missing",
			args: args{
				readViper(`
tls:
  rootCaCert:
    privateKey: ./ca.key
`),
			},
			want:    Options{},
			wantErr: true,
		},
		{
			name: "Get an error if CA private key path is missing",
			args: args{
				readViper(`
tls:
  rootCaCert:
    publicKey: ./ca.pem
`),
			},
			want:    Options{},
			wantErr: true,
		},
		{
			name: "Get default options if all required fields are set",
			args: args{
				readViper(`
tls:
  rootCaCert:
    publicKey: ./ca.pem
    privateKey: ./ca.key
`),
			},
			want: Options{
				RootCACert: File{
					PublicKeyPath:  "./ca.pem",
					PrivateKeyPath: "./ca.key",
				},
				CertCachePath: os.TempDir(),
				Curve:         CurveTypeED25519,
				Validity: ValidityByPurpose{
					CA: ValidityDuration{
						NotBeforeRelative: 17520 * time.Hour,
						NotAfterRelative:  17520 * time.Hour,
					},
					Server: ValidityDuration{
						NotBeforeRelative: 168 * time.Hour,
						NotAfterRelative:  168 * time.Hour,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
