//nolint:dupl
package rules_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/rules"
)

func TestParse(t *testing.T) {
	t.Parallel()
	type args struct {
		rule string
	}
	tests := []struct {
		name    string
		args    args
		target  interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "Routing - Terminator only - string argument",
			args: args{
				rule: `=> File("default.html")`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Name: "File",
					Params: []rules.Param{
						{
							String: rules.StringP("default.html"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator only - no argument",
			args: args{
				rule: `=> NoContent()`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Name: "NoContent",
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator with module - no argument",
			args: args{
				rule: `=> http.NoContent()`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Module: "http",
					Name:   "NoContent",
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator only do not support method name not starting with capital letter",
			args: args{
				rule: `=> file("default.html")`,
			},
			target:  new(rules.Routing),
			wantErr: true,
		},
		{
			name: "Routing - Terminator with module - string argument",
			args: args{
				rule: `=> http.File("default.html")`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Module: "http",
					Name:   "File",
					Params: []rules.Param{
						{
							String: rules.StringP("default.html"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator only - int argument",
			args: args{
				rule: `=> ReturnInt(1)`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Name: "ReturnInt",
					Params: []rules.Param{
						{
							Int: rules.IntP(1),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator with module - int argument",
			args: args{
				rule: `=> http.ReturnInt(1)`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Module: "http",
					Name:   "ReturnInt",
					Params: []rules.Param{
						{
							Int: rules.IntP(1),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator only - int argument, multiple digits",
			args: args{
				rule: `=> ReturnInt(1337)`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Name: "ReturnInt",
					Params: []rules.Param{
						{
							Int: rules.IntP(1337),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator with Module - int argument, multiple digits",
			args: args{
				rule: `=> http.ReturnInt(1337)`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Module: "http",
					Name:   "ReturnInt",
					Params: []rules.Param{
						{
							Int: rules.IntP(1337),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator only - float argument",
			args: args{
				rule: `=> ReturnFloat(13.37)`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Name: "ReturnFloat",
					Params: []rules.Param{
						{
							Float: rules.FloatP(13.37),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator only - IP argument",
			args: args{
				rule: `=> IP(8.8.8.8)`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Name: "IP",
					Params: []rules.Param{
						{
							IP: net.ParseIP("8.8.8.8"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - Terminator only - CIDR argument",
			args: args{
				rule: `=> IP(8.8.8.8/32)`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Name: "IP",
					Params: []rules.Param{
						{
							CIDR: rules.MustParseCIDR("8.8.8.8/32"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - path pattern and terminator",
			args: args{
				rule: `PathPattern(".*\\.(?i)png") => ReturnFile("default.html")`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Name: "ReturnFile",
					Params: []rules.Param{
						{
							String: rules.StringP("default.html"),
						},
					},
				},
				Filters: &rules.Filters{
					Chain: []rules.Method{
						{
							Name: "PathPattern",
							Params: []rules.Param{
								{
									String: rules.StringP(`.*\.(?i)png`),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - path pattern and terminator with Modules",
			args: args{
				rule: `http.PathPattern(".*\\.(?i)png") => http.ReturnFile("default.html")`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Module: "http",
					Name:   "ReturnFile",
					Params: []rules.Param{
						{
							String: rules.StringP("default.html"),
						},
					},
				},
				Filters: &rules.Filters{
					Chain: []rules.Method{
						{
							Module: "http",
							Name:   "PathPattern",
							Params: []rules.Param{
								{
									String: rules.StringP(`.*\.(?i)png`),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - HTTP method, path pattern and terminator",
			args: args{
				rule: `Method("GET") -> PathPattern("/index.html") => ReturnFile("default.html")`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Name: "ReturnFile",
					Params: []rules.Param{
						{
							String: rules.StringP("default.html"),
						},
					},
				},
				Filters: &rules.Filters{
					Chain: []rules.Method{
						{
							Name: "Method",
							Params: []rules.Param{
								{
									String: rules.StringP(http.MethodGet),
								},
							},
						},
						{
							Name: "PathPattern",
							Params: []rules.Param{
								{
									String: rules.StringP("/index.html"),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Routing - HTTP method, path pattern and terminator with modules",
			args: args{
				rule: `http.Method("GET") -> http.PathPattern("/index.html") => http.ReturnFile("default.html")`,
			},
			target: new(rules.Routing),
			want: &rules.Routing{
				Terminator: &rules.Method{
					Module: "http",
					Name:   "ReturnFile",
					Params: []rules.Param{
						{
							String: rules.StringP("default.html"),
						},
					},
				},
				Filters: &rules.Filters{
					Chain: []rules.Method{
						{
							Module: "http",
							Name:   "Method",
							Params: []rules.Param{
								{
									String: rules.StringP(http.MethodGet),
								},
							},
						},
						{
							Module: "http",
							Name:   "PathPattern",
							Params: []rules.Param{
								{
									String: rules.StringP("/index.html"),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Check - Initiator only - string argument",
			args: args{
				rule: `http.Get("https://www.microsoft.com/")`,
			},
			target: new(rules.Check),
			want: &rules.Check{
				Initiator: &rules.Method{
					Module: "http",
					Name:   "Get",
					Params: []rules.Param{
						{
							String: rules.StringP("https://www.microsoft.com/"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Check - Initiator only - raw string argument",
			args: args{
				rule: "http.Post(\"https://www.microsoft.com/\", `{\"Name\":\"Ted.Tester\"}`)",
			},
			target: new(rules.Check),
			want: &rules.Check{
				Initiator: &rules.Method{
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
		{
			name: "Check - Initiator and single filter",
			args: args{
				rule: `http.Get("https://www.microsoft.com/") => Status(200)`,
			},
			target: new(rules.Check),
			want: &rules.Check{
				Initiator: &rules.Method{
					Module: "http",
					Name:   "Get",
					Params: []rules.Param{
						{
							String: rules.StringP("https://www.microsoft.com/"),
						},
					},
				},
				Validators: &rules.Filters{
					Chain: []rules.Method{
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
		{
			name: "Check - Initiator and multiple filters",
			args: args{
				rule: `http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text/html")`,
			},
			target: new(rules.Check),
			want: &rules.Check{
				Initiator: &rules.Method{
					Module: "http",
					Name:   "Get",
					Params: []rules.Param{
						{
							String: rules.StringP("https://www.microsoft.com/"),
						},
					},
				},
				Validators: &rules.Filters{
					Chain: []rules.Method{
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
		{
			name: "CheckScript without comments",
			args: args{
				rule: `
http.Get("https://www.gogol.com/") => Status(404)
http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text.html")
`,
			},
			target: new(rules.CheckScript),
			want: td.Struct(new(rules.CheckScript), td.StructFields{
				"Checks": td.Len(2),
			}),
		},
		{
			name: "CheckScript with comments",
			args: args{
				rule: `
# GET https://www.gogol.com/ expect a not found response
http.Get("https://www.gogol.com/") => Status(404)

// GET https://www.microsoft.com/ - expect status OK and HTML content
http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text/html")
`,
			},
			target: new(rules.CheckScript),
			want: td.Struct(new(rules.CheckScript), td.StructFields{
				"Checks": td.Len(2),
			}),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := rules.Parse(tt.args.rule, tt.target)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			td.Cmp(t, tt.target, tt.want)
		})
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
