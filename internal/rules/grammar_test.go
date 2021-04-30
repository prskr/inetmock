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
			name: "HTTP method, path pattern and terminator",
			args: args{
				rule: `HTTPMethod("GET") -> PathPattern("/index.html") => ReturnFile("default.html")`,
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
							Name: "HTTPMethod",
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
