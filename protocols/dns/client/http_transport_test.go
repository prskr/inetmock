package client_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"

	. "github.com/maxatome/go-testdeep/td"
	mdns "github.com/miekg/dns"

	protocolmock "gitlab.com/inetmock/inetmock/internal/mock/protocol"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/protocols/dns/client"
)

func TestRequestPackerPOST(t *testing.T) {
	t.Parallel()
	type args struct {
		question *mdns.Msg
	}
	tests := []struct {
		name    string
		packer  client.RequestPacker
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "GET packer",
			packer: client.RequestPackerGET,
			args: args{
				question: func() *mdns.Msg {
					q := new(mdns.Msg)
					q.SetQuestion(mdns.Fqdn("gitlab.com"), mdns.TypeA)
					return q
				}(),
			},
			want: Struct(&http.Request{Method: http.MethodGet}, StructFields{
				"Body": Nil(),
				"URL": Code(func(u *url.URL) error {
					if u == nil {
						return errors.New("URL is nil")
					}
					values := u.Query()
					if !values.Has("dns") {
						return errors.New("expected 'dns' key to be present in query values")
					}

					return nil
				}),
			}),
			wantErr: false,
		},
		{
			name:   "POST packer",
			packer: client.RequestPackerPOST,
			args: args{
				question: func() *mdns.Msg {
					q := new(mdns.Msg)
					q.SetQuestion(mdns.Fqdn("gitlab.com"), mdns.TypeA)
					return q
				}(),
			},
			want: Struct(&http.Request{Method: http.MethodPost, URL: test.MustParseURL("https://quad9.com/dns-query")}, StructFields{
				"Body": NotNil(),
			}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			queryURL := test.MustParseURL("https://quad9.com/dns-query")
			if got, err := tt.packer.Pack(queryURL, tt.args.question); (err != nil) != tt.wantErr {
				t.Errorf("WantErr = %t, but got %v", tt.wantErr, err)
				return
			} else {
				Cmp(t, got, tt.want)
			}
		})
	}
}

func TestHTTPTransport_RoundTrip(t *testing.T) {
	t.Parallel()
	type fields struct {
		Client func(tb testing.TB, question *mdns.Msg) client.HTTPClient
		Scheme string
		Server string
	}
	tests := []struct {
		name     string
		fields   fields
		question *mdns.Msg
		wantResp interface{}
		wantErr  bool
	}{
		{
			name: "Execute an A query",
			fields: fields{
				Client: func(tb testing.TB, question *mdns.Msg) client.HTTPClient {
					tb.Helper()
					return &protocolmock.HTTPClientMock{
						OnDo: func(state protocolmock.HTTPClientMockContext, req *http.Request) (*http.Response, error) {
							resp := new(mdns.Msg)
							resp.SetReply(question)
							q := question.Question[0]
							resp.Answer = append(resp.Answer, &mdns.A{
								A:   net.IPv4(10, 0, 0, 1),
								Hdr: mdns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: mdns.ClassINET, Ttl: 30},
							})

							data, err := resp.Pack()
							CmpNoError(tb, err)

							return &http.Response{
								StatusCode: http.StatusOK,
								Body:       io.NopCloser(bytes.NewReader(data)),
							}, nil
						},
					}
				},
				Scheme: "https",
				Server: "quad9.com",
			},
			question: new(mdns.Msg).SetQuestion(mdns.Fqdn("gitlab.com"), mdns.TypeA),
			wantResp: NotNil(),
			wantErr:  false,
		},
		{
			name: "Return non-success code",
			fields: fields{
				Client: func(tb testing.TB, question *mdns.Msg) client.HTTPClient {
					tb.Helper()
					return &protocolmock.HTTPClientMock{
						OnDo: func(state protocolmock.HTTPClientMockContext, req *http.Request) (*http.Response, error) {
							return &http.Response{
								StatusCode: http.StatusBadRequest,
							}, nil
						},
					}
				},
				Scheme: "https",
				Server: "quad9.com",
			},
			question: new(mdns.Msg).SetQuestion(mdns.Fqdn("gitlab.com"), mdns.TypeA),
			wantErr:  true,
		},
		{
			name: "Return error from client",
			fields: fields{
				Client: func(tb testing.TB, question *mdns.Msg) client.HTTPClient {
					tb.Helper()
					return &protocolmock.HTTPClientMock{
						OnDo: func(state protocolmock.HTTPClientMockContext, req *http.Request) (*http.Response, error) {
							return nil, errors.New("something wrong")
						},
					}
				},
				Scheme: "https",
				Server: "quad9.com",
			},
			question: new(mdns.Msg).SetQuestion(mdns.Fqdn("gitlab.com"), mdns.TypeA),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := client.HTTPTransport{
				Packer: client.RequestPackerPOST,
				Client: tt.fields.Client(t, tt.question),
				Scheme: tt.fields.Scheme,
				Server: tt.fields.Server,
			}
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			gotResp, err := h.RoundTrip(ctx, tt.question)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			Cmp(t, gotResp, tt.wantResp)
		})
	}
}
