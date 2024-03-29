package rules_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
)

type testCase interface {
	Run(t *testing.T)
	Name() string
}

type parseTest[T any] struct {
	name    string
	rule    string
	parser  func(rule string) (*T, error)
	want    any
	wantErr bool
}

func (pt parseTest[T]) Name() string {
	return pt.name
}

func (pt parseTest[T]) Run(t *testing.T) {
	t.Helper()
	t.Parallel()
	got, err := pt.parser(pt.rule)

	if (err != nil) != pt.wantErr {
		t.Errorf("pt.wantErr = %v but got error %v", pt.wantErr, err)
	}

	if pt.wantErr {
		return
	}

	td.Cmp(t, got, pt.want)
}

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []testCase{
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response only - string argument",
			rule:   `=> File("default.html")`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Name:   "File",
					Params: params(rules.Param{String: rules.StringP("default.html")}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response only - no argument",
			rule:   `=> NoContent()`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{Name: "NoContent"},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response with module - no argument",
			rule:   `=> http.NoContent()`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Module: "http",
					Name:   "NoContent",
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:    "SingleResponsePipeline - Response only do not support method name not starting with capital letter",
			rule:    `=> file("default.html")`,
			parser:  rules.Parse[rules.SingleResponsePipeline],
			wantErr: true,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response with module - string argument",
			rule:   `=> http.File("default.html")`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Module: "http",
					Name:   "File",
					Params: params(rules.Param{String: rules.StringP("default.html")}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response only - int argument",
			rule:   `=> ReturnInt(1)`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Name:   "ReturnInt",
					Params: params(rules.Param{Int: rules.IntP(1)}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response with module - int argument",
			rule:   `=> http.ReturnInt(1)`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Module: "http",
					Name:   "ReturnInt",
					Params: params(rules.Param{Int: rules.IntP(1)}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response only - int argument, multiple digits",
			rule:   `=> ReturnInt(1337)`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Name:   "ReturnInt",
					Params: params(rules.Param{Int: rules.IntP(1337)}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response with Module - int argument, multiple digits",
			rule:   `=> http.ReturnInt(1337)`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Module: "http",
					Name:   "ReturnInt",
					Params: params(rules.Param{Int: rules.IntP(1337)}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response only - float argument",
			rule:   `=> ReturnFloat(13.37)`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Name:   "ReturnFloat",
					Params: params(rules.Param{Float: rules.FloatP(13.37)}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response only - IP argument",
			rule:   `=> IP(8.8.8.8)`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Name:   "IP",
					Params: params(rules.Param{IP: net.ParseIP("8.8.8.8")}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - Response only - CIDR argument",
			rule:   `=> IP(8.8.8.8/32)`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Name:   "IP",
					Params: params(rules.Param{CIDR: rules.MustParseCIDR("8.8.8.8/32")}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - path pattern and terminator",
			rule:   `PathPattern(".*\\.(?i)png") => ReturnFile("default.html")`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Name:   "ReturnFile",
					Params: params(rules.Param{String: rules.StringP("default.html")}),
				},
				FilterChain: &rules.Filters{
					Chain: []rules.Call{
						{
							Name:   "PathPattern",
							Params: params(rules.Param{String: rules.StringP(`.*\.(?i)png`)}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - path pattern and terminator with Modules",
			rule:   `http.PathPattern(".*\\.(?i)png") => http.ReturnFile("default.html")`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Module: "http",
					Name:   "ReturnFile",
					Params: params(rules.Param{String: rules.StringP("default.html")}),
				},
				FilterChain: &rules.Filters{
					Chain: []rules.Call{
						{
							Module: "http",
							Name:   "PathPattern",
							Params: params(rules.Param{String: rules.StringP(`.*\.(?i)png`)}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - HTTP method, path pattern and terminator",
			rule:   `Method("GET") -> PathPattern("/index.html") => ReturnFile("default.html")`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Name:   "ReturnFile",
					Params: params(rules.Param{String: rules.StringP("default.html")}),
				},
				FilterChain: &rules.Filters{
					Chain: []rules.Call{
						{
							Name:   "Method",
							Params: params(rules.Param{String: rules.StringP(http.MethodGet)}),
						},
						{
							Name:   "PathPattern",
							Params: params(rules.Param{String: rules.StringP("/index.html")}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[rules.SingleResponsePipeline]{
			name:   "SingleResponsePipeline - HTTP method, path pattern and terminator with modules",
			rule:   `http.Method("GET") -> http.PathPattern("/index.html") => http.ReturnFile("default.html")`,
			parser: rules.Parse[rules.SingleResponsePipeline],
			want: &rules.SingleResponsePipeline{
				Response: &rules.Call{
					Module: "http",
					Name:   "ReturnFile",
					Params: params(rules.Param{String: rules.StringP("default.html")}),
				},
				FilterChain: &rules.Filters{
					Chain: []rules.Call{
						{
							Module: "http",
							Name:   "Method",
							Params: params(rules.Param{String: rules.StringP(http.MethodGet)}),
						},
						{
							Module: "http",
							Name:   "PathPattern",
							Params: params(rules.Param{String: rules.StringP("/index.html")}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[rules.ChainedResponsePipeline]{
			name:   "ChainedResponsePipeline - single response only - IP argument",
			rule:   `=> IP(1.2.3.4)`,
			parser: rules.Parse[rules.ChainedResponsePipeline],
			want: &rules.ChainedResponsePipeline{
				Response: []rules.Call{{Name: "IP", Params: params(rules.Param{IP: net.IPv4(1, 2, 3, 4)})}},
			},
			wantErr: false,
		},
		parseTest[rules.ChainedResponsePipeline]{
			name:   "ChainedResponsePipeline - multi-response only - IP arguments",
			rule:   `=> IP(1.2.3.4) => Router(1.2.3.1)`,
			parser: rules.Parse[rules.ChainedResponsePipeline],
			want: &rules.ChainedResponsePipeline{
				Response: calls(
					rules.Call{Name: "IP", Params: params(rules.Param{IP: net.IPv4(1, 2, 3, 4)})},
					rules.Call{Name: "Router", Params: params(rules.Param{IP: net.IPv4(1, 2, 3, 1)})},
				),
			},
			wantErr: false,
		},
		parseTest[rules.ChainedResponsePipeline]{
			name:   "ChainedResponsePipeline - single response - single filter - IP argument",
			rule:   `MatchMAC("00:06:7C:.*") => IP(1.2.3.4)`,
			parser: rules.Parse[rules.ChainedResponsePipeline],
			want: &rules.ChainedResponsePipeline{
				FilterChain: &rules.Filters{
					Chain: []rules.Call{
						{
							Name:   "MatchMAC",
							Params: params(rules.Param{String: rules.StringP(`00:06:7C:.*`)}),
						},
					},
				},
				Response: []rules.Call{{Name: "IP", Params: params(rules.Param{IP: net.IPv4(1, 2, 3, 4)})}},
			},
			wantErr: false,
		},
		parseTest[rules.ChainedResponsePipeline]{
			name:   "ChainedResponsePipeline - single response - single filter - IP argument",
			rule:   `MatchMAC("00:06:7C:.*") => IP(1.2.3.4)`,
			parser: rules.Parse[rules.ChainedResponsePipeline],
			want: &rules.ChainedResponsePipeline{
				FilterChain: &rules.Filters{
					Chain: []rules.Call{
						{
							Name:   "MatchMAC",
							Params: params(rules.Param{String: rules.StringP(`00:06:7C:.*`)}),
						},
					},
				},
				Response: []rules.Call{{Name: "IP", Params: params(rules.Param{IP: net.IPv4(1, 2, 3, 4)})}},
			},
			wantErr: false,
		},
		parseTest[rules.ChainedResponsePipeline]{
			name:   "ChainedResponsePipeline - single response - multiple filters - IP argument",
			rule:   `RequestType('Discover') -> MatchMAC("00:06:7C:.*") => IP(1.2.3.4)`,
			parser: rules.Parse[rules.ChainedResponsePipeline],
			want: &rules.ChainedResponsePipeline{
				FilterChain: &rules.Filters{
					Chain: []rules.Call{
						{
							Name:   "RequestType",
							Params: params(rules.Param{String: rules.StringP("Discover")}),
						},
						{
							Name:   "MatchMAC",
							Params: params(rules.Param{String: rules.StringP(`00:06:7C:.*`)}),
						},
					},
				},
				Response: []rules.Call{{Name: "IP", Params: params(rules.Param{IP: net.IPv4(1, 2, 3, 4)})}},
			},
			wantErr: false,
		},
		parseTest[rules.ChainedResponsePipeline]{
			name:   "ChainedResponsePipeline - multiple responses - multiple filters - IP argument",
			rule:   `RequestType('Discover') -> MatchMAC("00:06:7C:.*") => IP(1.2.3.4) => Router(1.2.3.1)`,
			parser: rules.Parse[rules.ChainedResponsePipeline],
			want: &rules.ChainedResponsePipeline{
				FilterChain: &rules.Filters{
					Chain: []rules.Call{
						{
							Name:   "RequestType",
							Params: params(rules.Param{String: rules.StringP("Discover")}),
						},
						{
							Name:   "MatchMAC",
							Params: params(rules.Param{String: rules.StringP(`00:06:7C:.*`)}),
						},
					},
				},
				Response: []rules.Call{
					{Name: "IP", Params: params(rules.Param{IP: net.IPv4(1, 2, 3, 4)})},
					{Name: "Router", Params: params(rules.Param{IP: net.IPv4(1, 2, 3, 1)})},
				},
			},
			wantErr: false,
		},
		parseTest[rules.Check]{
			name:   "Check - Initiator only - string argument",
			rule:   `http.Get("https://www.microsoft.com/")`,
			parser: rules.Parse[rules.Check],
			want: &rules.Check{
				Initiator: &rules.Call{
					Module: "http",
					Name:   "Get",
					Params: params(rules.Param{String: rules.StringP("https://www.microsoft.com/")}),
				},
			},
			wantErr: false,
		},
		parseTest[rules.Check]{
			name:   "Check - Initiator only - raw string argument",
			rule:   "http.Post(\"https://www.microsoft.com/\", `{\"Name\":\"Ted.Tester\"}`)",
			parser: rules.Parse[rules.Check],
			want: &rules.Check{
				Initiator: &rules.Call{
					Module: "http",
					Name:   "Post",
					Params: []rules.Param{
						{
							String: rules.StringP("https://www.microsoft.com/"),
						},
						{
							String: rules.StringP(`{"Name":"Ted.Tester"}`),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[rules.Check]{
			name:   "Check - Initiator and single filter",
			rule:   `http.Get("https://www.microsoft.com/") => Status(200)`,
			parser: rules.Parse[rules.Check],
			want: &rules.Check{
				Initiator: &rules.Call{
					Module: "http",
					Name:   "Get",
					Params: params(rules.Param{String: rules.StringP("https://www.microsoft.com/")}),
				},
				Validators: &rules.Filters{
					Chain: []rules.Call{
						{
							Name: "Status",
							Params: []rules.Param{
								{
									Int: rules.IntP(200),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[rules.Check]{
			name:   "Check - Initiator and multiple filters",
			rule:   `http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text/html")`,
			parser: rules.Parse[rules.Check],
			want: &rules.Check{
				Initiator: &rules.Call{
					Module: "http",
					Name:   "Get",
					Params: params(rules.Param{String: rules.StringP("https://www.microsoft.com/")}),
				},
				Validators: &rules.Filters{
					Chain: []rules.Call{
						{
							Name: "Status",
							Params: []rules.Param{
								{
									Int: rules.IntP(200),
								},
							},
						},
						{
							Name: "Header",
							Params: []rules.Param{
								{
									String: rules.StringP("Content-Type"),
								},
								{
									String: rules.StringP("text/html"),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[rules.CheckScript]{
			name: "CheckScript without comments",
			rule: `
http.Get("https://www.gogol.com/") => Status(404)
http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text.html")
`,
			parser: rules.Parse[rules.CheckScript],
			want: td.Struct(new(rules.CheckScript), td.StructFields{
				"Checks": td.Len(2),
			}),
		},
		parseTest[rules.CheckScript]{
			name: "CheckScript with comments",
			rule: `
# GET https://www.gogol.com/ expect a not found response
http.Get("https://www.gogol.com/") => Status(404)

// GET https://www.microsoft.com/ - expect status OK and HTML content
http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text/html")
`,
			parser: rules.Parse[rules.CheckScript],
			want: td.Struct(new(rules.CheckScript), td.StructFields{
				"Checks": td.Len(2),
			}),
		},
	}
	//nolint:paralleltest // is actually called in Run function
	for _, tc := range tests {
		tt := tc
		t.Run(tt.Name(), tt.Run)
	}
}

func TestParam_AsString(t *testing.T) {
	t.Parallel()
	type fields struct {
		String *string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Empty string",
			fields: fields{
				String: rules.StringP(""),
			},
			want: "",
		},
		{
			name: "Any string",
			fields: fields{
				String: rules.StringP("Hello, world!"),
			},
			want: "Hello, world!",
		},
		{
			name:    "nil value",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := rules.Param{
				String: tt.fields.String,
			}
			got, err := p.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AsString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParam_AsInt(t *testing.T) {
	t.Parallel()
	type fields struct {
		Int *int
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "zero value",
			fields: fields{
				Int: rules.IntP(0),
			},
			want: 0,
		},
		{
			name: "Any int",
			fields: fields{
				Int: rules.IntP(42),
			},
			want: 42,
		},
		{
			name:    "nil value",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := rules.Param{
				Int: tt.fields.Int,
			}
			got, err := p.AsInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AsInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParam_AsFloat(t *testing.T) {
	t.Parallel()
	type fields struct {
		Float *float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{
			name: "Zero value",
			fields: fields{
				Float: rules.FloatP(0),
			},
			want: 0,
		},
		{
			name: "Any value",
			fields: fields{
				Float: rules.FloatP(13.37),
			},
			want: 13.37,
		},
		{
			name:    "nil value",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := rules.Param{
				Float: tt.fields.Float,
			}
			got, err := p.AsFloat()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AsFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func params(p ...rules.Param) []rules.Param {
	return p
}

func calls(c ...rules.Call) []rules.Call {
	return c
}
