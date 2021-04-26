package rules

import (
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

var (
	defaultHtml = "default.html"
	defaultPath = "/index.html"
	methodGet   = http.MethodGet
)

func TestParse(t *testing.T) {
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
			name: "Terminator only",
			args: args{
				rule: `=> ReturnFile("default.html")`,
			},
			want: &Routing{
				Terminator: &Method{
					Name: "ReturnFile",
					Params: []Param{
						{
							String: &defaultHtml,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "path pattern and terminator",
			args: args{
				rule: `PathPattern("/index.html") => ReturnFile("default.html")`,
			},
			want: &Routing{
				Terminator: &Method{
					Name: "ReturnFile",
					Params: []Param{
						{
							String: &defaultHtml,
						},
					},
				},
				Filters: &Filters{
					Chain: []Method{
						{
							Name: "PathPattern",
							Params: []Param{
								{
									String: &defaultPath,
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
							String: &defaultHtml,
						},
					},
				},
				Filters: &Filters{
					Chain: []Method{
						{
							Name: "HTTPMethod",
							Params: []Param{
								{
									String: &methodGet,
								},
							},
						},
						{
							Name: "PathPattern",
							Params: []Param{
								{
									String: &defaultPath,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}
