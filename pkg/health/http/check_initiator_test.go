package http_test

import (
	"context"
	"errors"
	"net"
	gohttp "net/http"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/health/http"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/http/mock"
)

func TestInitiatorForRule(t *testing.T) {
	t.Parallel()
	type args struct {
		rule *rules.Check
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Initiator for GET request",
			args: args{
				rule: &rules.Check{
					Initiator: &rules.Call{
						Module: "http",
						Name:   "Get",
						Params: []rules.Param{
							{
								String: rules.StringP("https://www.google.com"),
							},
						},
					},
				},
			},
			want:    td.NotNil(),
			wantErr: false,
		},
		{
			name: "Initiator for POST request without parameter",
			args: args{
				rule: &rules.Check{
					Initiator: &rules.Call{
						Module: "http",
						Name:   "Post",
						Params: []rules.Param{
							{
								String: rules.StringP("https://www.google.com"),
							},
						},
					},
				},
			},
			want:    td.NotNil(),
			wantErr: false,
		},
		{
			name: "Initiator for POST request with JSON parameter",
			args: args{
				rule: &rules.Check{
					Initiator: &rules.Call{
						Module: "http",
						Name:   "Post",
						Params: []rules.Param{
							{
								String: rules.StringP("https://www.google.com"),
							},
							{
								String: rules.StringP(`{"firstName":"Ted"}`),
							},
						},
					},
				},
			},
			want:    td.NotNil(),
			wantErr: false,
		},
		{
			name: "Error wrong module",
			args: args{
				rule: &rules.Check{
					Initiator: &rules.Call{
						Module: "dns",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error unknown initiator",
			args: args{
				rule: &rules.Check{
					Initiator: &rules.Call{
						Module: "http",
						Name:   "proxy",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error missing argument",
			args: args{
				rule: &rules.Check{
					Initiator: &rules.Call{
						Module: "http",
						Name:   "get",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error missing initiator",
			args: args{
				rule: new(rules.Check),
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := logging.CreateTestLogger(t)
			got, err := http.InitiatorForRule(tt.args.rule, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitiatorForRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}

func TestRequestInitiatorForMethod(t *testing.T) {
	t.Parallel()
	type args struct {
		method       string
		ruleParams   []rules.Param
		routingRules []string
	}
	tests := []struct {
		name               string
		args               args
		wantInitiatorError bool
		wantRequestError   bool
		wantResponse       interface{}
	}{
		{
			name: "Run a GET request",
			args: args{
				method: gohttp.MethodGet,
				ruleParams: []rules.Param{
					{
						String: rules.StringP("https://www.google.de/search"),
					},
				},
				routingRules: []string{
					`Method("GET") -> PathPattern("^\\/search$") => Status(204)`,
				},
			},
			wantInitiatorError: false,
			wantRequestError:   false,
			wantResponse: td.Struct(new(gohttp.Response), td.StructFields{
				"StatusCode": 204,
			}),
		},
		{
			name: "Run a GET request without matching rule",
			args: args{
				method: gohttp.MethodGet,
				ruleParams: []rules.Param{
					{
						String: rules.StringP("https://www.google.de/s3arch"),
					},
				},
				routingRules: []string{
					`Method("GET") -> PathPattern("^\\/search$") => Status(204)`,
				},
			},
			wantInitiatorError: false,
			wantRequestError:   false,
			wantResponse: td.Struct(new(gohttp.Response), td.StructFields{
				"StatusCode": 404,
			}),
		},
		{
			name: "Run a POST request",
			args: args{
				method: gohttp.MethodPost,
				ruleParams: []rules.Param{
					{
						String: rules.StringP("https://stackoverflow.com/api/v1/question"),
					},
				},
				routingRules: []string{
					`Method("POST") => Status(204)`,
				},
			},
			wantInitiatorError: false,
			wantRequestError:   false,
			wantResponse: td.Struct(new(gohttp.Response), td.StructFields{
				"StatusCode": 204,
			}),
		},
		{
			name: "Run a POST request with a JSON body",
			args: args{
				method: gohttp.MethodPost,
				ruleParams: []rules.Param{
					{
						String: rules.StringP("https://stackoverflow.com/api/v1/question"),
					},
					{
						String: rules.StringP(`{"question": "How to do 'Hello, world!' in Go?"}`),
					},
				},
				routingRules: []string{
					`Method("POST") => Status(204)`,
				},
			},
			wantInitiatorError: false,
			wantRequestError:   false,
			wantResponse: td.Struct(new(gohttp.Response), td.StructFields{
				"StatusCode": 204,
			}),
		},
		{
			name: "Run a POST request without matching rule",
			args: args{
				method: gohttp.MethodPost,
				ruleParams: []rules.Param{
					{
						String: rules.StringP("https://stackoverflow.com/api/v1/question"),
					},
				},
				routingRules: []string{
					`Method("GET") -> PathPattern("^\\/.*") => Status(204)`,
				},
			},
			wantInitiatorError: false,
			wantRequestError:   false,
			wantResponse: td.Struct(new(gohttp.Response), td.StructFields{
				"StatusCode": 404,
			}),
		},
		{
			name: "Run a PUT request",
			args: args{
				method: gohttp.MethodPut,
				ruleParams: []rules.Param{
					{
						String: rules.StringP("https://stackoverflow.com/api/v1/question"),
					},
					{
						String: rules.StringP(`{"question": "How to do 'Hello, world!' in Go?"}`),
					},
				},
				routingRules: []string{
					`Method("PUT") -> PathPattern("^/api/v1/question$") => Status(204)`,
				},
			},
			wantInitiatorError: false,
			wantRequestError:   false,
			wantResponse: td.Struct(new(gohttp.Response), td.StructFields{
				"StatusCode": 204,
			}),
		},
		{
			name: "Run a PUT request without matching rule",
			args: args{
				method: gohttp.MethodPut,
				ruleParams: []rules.Param{
					{
						String: rules.StringP("https://stackoverflow.com/api/v1/question"),
					},
					{
						String: rules.StringP(`{"question": "How to do 'Hello, world!' in Go?"}`),
					},
				},
				routingRules: []string{
					`Method("PUT") -> PathPattern("^\\/search$") => Status(204)`,
				},
			},
			wantInitiatorError: false,
			wantRequestError:   false,
			wantResponse: td.Struct(new(gohttp.Response), td.StructFields{
				"StatusCode": 404,
			}),
		},
		{
			name: "Run a DELETE request",
			args: args{
				method: gohttp.MethodDelete,
				ruleParams: []rules.Param{
					{
						String: rules.StringP("https://stackoverflow.com/api/v1/question/ccd4f873-a244-4ef2-b53d-f08efb35db2d"),
					},
				},
				routingRules: []string{
					`Method("DELETE") -> PathPattern("^/api/v1/question") => Status(204)`,
				},
			},
			wantInitiatorError: false,
			wantRequestError:   false,
			wantResponse: td.Struct(new(gohttp.Response), td.StructFields{
				"StatusCode": 204,
			}),
		},
		{
			name: "Run a DELETE request without matching rule",
			args: args{
				method: gohttp.MethodDelete,
				ruleParams: []rules.Param{
					{
						String: rules.StringP("https://stackoverflow.com/api/v1/question/ccd4f873-a244-4ef2-b53d-f08efb35db2d"),
					},
				},
				routingRules: []string{
					`Method("DELETE") -> PathPattern("^/search$") => Status(204)`,
				},
			},
			wantInitiatorError: false,
			wantRequestError:   false,
			wantResponse: td.Struct(new(gohttp.Response), td.StructFields{
				"StatusCode": 404,
			}),
		},
		{
			name: "Error argument type mismatch",
			args: args{
				method: gohttp.MethodDelete,
				ruleParams: []rules.Param{
					{
						Int: rules.IntP(404),
					},
				},
			},
			wantInitiatorError: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := logging.CreateTestLogger(t)
			got := http.RequestInitiatorForMethod(tt.args.method)

			var err error
			var initiator http.Initiator
			if initiator, err = got(logger, tt.args.ruleParams...); err != nil {
				if !tt.wantInitiatorError {
					t.Errorf("Creating initiator got unexpected error = %v", err)
				}
				return
			}

			client := setupServer(t, tt.args.routingRules)

			ctx, cancel := context.WithTimeout(test.Context(t), 50*time.Millisecond)
			t.Cleanup(cancel)

			var resp *gohttp.Response
			if resp, err = initiator.Do(ctx, client); err != nil {
				if !tt.wantRequestError {
					t.Errorf("Executing initiator returned unexpected error - error = %v", err)
				}
				return
			}

			td.Cmp(t, resp, tt.wantResponse)
		})
	}
}

func setupServer(tb testing.TB, routingRules []string) *gohttp.Client {
	tb.Helper()
	logger := logging.CreateTestLogger(tb)

	var err error
	var stream audit.EventStream
	if stream, err = audit.NewEventStream(logger); err != nil {
		tb.Fatalf("Failed to init EventStream = error %v", err)
	}

	tb.Cleanup(func() {
		_ = stream.Close()
	})

	router := mock.Router{
		HandlerName: "HealthMockHandler",
		Logger:      logger,
	}

	for _, rule := range routingRules {
		if err := router.RegisterRule(rule); err != nil {
			tb.Fatalf("failed to setup rule %s - error = %v", rule, err)
		}
	}

	listener := test.NewInMemoryListener(tb)

	go func(tb testing.TB, listener net.Listener, handler gohttp.Handler) {
		tb.Helper()
		switch err := gohttp.Serve(listener, handler); {
		case errors.Is(err, nil), errors.Is(err, gohttp.ErrServerClosed):
		default:
			tb.Errorf("srv.Serve() error = %v", err)
		}
	}(tb, listener, &router)

	return test.HTTPClientForInMemListener(listener)
}
