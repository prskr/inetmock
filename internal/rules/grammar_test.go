package rules

import (
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

//nolint:funlen
func TestParse(t *testing.T) {
	t.Parallel()
	type args struct {
		rule string
	}
	tests := []struct {
		name    string
		args    args
		want    *Routing
		wantErr bool
	}{
		{
			name: "Terminator only - string argument",
			args: args{
				rule: `=> File("default.html")`,
			},
			want: &Routing{
				Terminator: &Method{
					Name: "File",
					Params: []Param{
						{
							String: stringRef("default.html"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Terminator only - no argument",
			args: args{
				rule: `=> NoContent()`,
			},
			want: &Routing{
				Terminator: &Method{
					Name: "NoContent",
				},
			},
			wantErr: false,
		},
		{
			name: "Terminator with module - no argument",
			args: args{
				rule: `=> http.NoContent()`,
			},
			want: &Routing{
				Terminator: &Method{
					Module: "http",
					Name:   "NoContent",
				},
			},
			wantErr: false,
		},
		{
			name: "Terminator only do not support method name not starting with capital letter",
			args: args{
				rule: `=> file("default.html")`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Terminator with module - string argument",
			args: args{
				rule: `=> http.File("default.html")`,
			},
			want: &Routing{
				Terminator: &Method{
					Module: "http",
					Name:   "File",
					Params: []Param{
						{
							String: stringRef("default.html"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Terminator only - int argument",
			args: args{
				rule: `=> ReturnInt(1)`,
			},
			want: &Routing{
				Terminator: &Method{
					Name: "ReturnInt",
					Params: []Param{
						{
							Int: intRef(1),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Terminator with module - int argument",
			args: args{
				rule: `=> http.ReturnInt(1)`,
			},
			want: &Routing{
				Terminator: &Method{
					Module: "http",
					Name:   "ReturnInt",
					Params: []Param{
						{
							Int: intRef(1),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Terminator only - int argument, multiple digits",
			args: args{
				rule: `=> ReturnInt(1337)`,
			},
			want: &Routing{
				Terminator: &Method{
					Name: "ReturnInt",
					Params: []Param{
						{
							Int: intRef(1337),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Terminator with Module - int argument, multiple digits",
			args: args{
				rule: `=> http.ReturnInt(1337)`,
			},
			want: &Routing{
				Terminator: &Method{
					Module: "http",
					Name:   "ReturnInt",
					Params: []Param{
						{
							Int: intRef(1337),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Terminator only - float argument",
			args: args{
				rule: `=> ReturnFloat(13.37)`,
			},
			want: &Routing{
				Terminator: &Method{
					Name: "ReturnFloat",
					Params: []Param{
						{
							Float: floatRef(13.37),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "path pattern and terminator",
			args: args{
				rule: `PathPattern(".*\\.(?i)png") => ReturnFile("default.html")`,
			},
			want: &Routing{
				Terminator: &Method{
					Name: "ReturnFile",
					Params: []Param{
						{
							String: stringRef("default.html"),
						},
					},
				},
				Filters: &Filters{
					Chain: []Method{
						{
							Name: "PathPattern",
							Params: []Param{
								{
									String: stringRef(`.*\.(?i)png`),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "path pattern and terminator with Modules",
			args: args{
				rule: `http.PathPattern(".*\\.(?i)png") => http.ReturnFile("default.html")`,
			},
			want: &Routing{
				Terminator: &Method{
					Module: "http",
					Name:   "ReturnFile",
					Params: []Param{
						{
							String: stringRef("default.html"),
						},
					},
				},
				Filters: &Filters{
					Chain: []Method{
						{
							Module: "http",
							Name:   "PathPattern",
							Params: []Param{
								{
									String: stringRef(`.*\.(?i)png`),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "HTTP method, path pattern and terminator",
			args: args{
				rule: `Method("GET") -> PathPattern("/index.html") => ReturnFile("default.html")`,
			},
			want: &Routing{
				Terminator: &Method{
					Name: "ReturnFile",
					Params: []Param{
						{
							String: stringRef("default.html"),
						},
					},
				},
				Filters: &Filters{
					Chain: []Method{
						{
							Name: "Method",
							Params: []Param{
								{
									String: stringRef(http.MethodGet),
								},
							},
						},
						{
							Name: "PathPattern",
							Params: []Param{
								{
									String: stringRef("/index.html"),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "HTTP method, path pattern and terminator with modules",
			args: args{
				rule: `http.Method("GET") -> http.PathPattern("/index.html") => http.ReturnFile("default.html")`,
			},
			want: &Routing{
				Terminator: &Method{
					Module: "http",
					Name:   "ReturnFile",
					Params: []Param{
						{
							String: stringRef("default.html"),
						},
					},
				},
				Filters: &Filters{
					Chain: []Method{
						{
							Module: "http",
							Name:   "Method",
							Params: []Param{
								{
									String: stringRef(http.MethodGet),
								},
							},
						},
						{
							Module: "http",
							Name:   "PathPattern",
							Params: []Param{
								{
									String: stringRef("/index.html"),
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Parse(tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}

func stringRef(s string) *string {
	return &s
}

func intRef(i int) *int {
	return &i
}

func floatRef(f float64) *float64 {
	return &f
}
