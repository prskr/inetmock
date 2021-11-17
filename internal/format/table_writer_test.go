package format_test

import (
	"strings"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/format"
)

func Test_tblWriter_Write(t *testing.T) {
	t.Parallel()
	type s1 struct {
		Name string
		Age  int
	}

	type s2 struct {
		Name string `table:"Full name"`
		Age  int    `table:"Age in years"`
	}

	type s3 struct {
		Name        string
		Age         int
		privateCity string
	}

	type args struct {
		in interface{}
	}
	type testCase struct {
		name       string
		args       args
		wantErr    bool
		wantResult string
	}
	tests := []testCase{
		{
			name: "Test write table without errors",
			args: args{
				in: s1{
					Name: "Ted Tester",
					Age:  28,
				},
			},
			wantErr: false,
			wantResult: `
|    NAME    | AGE |
|------------|-----|
| Ted Tester |  28 |
`,
		},
		{
			name: "Test write table without errors with pointer value",
			args: args{
				in: &s1{
					Name: "Ted Tester",
					Age:  28,
				},
			},
			wantErr: false,
			wantResult: `
|    NAME    | AGE |
|------------|-----|
| Ted Tester |  28 |
`,
		},
		{
			name: "Test write table without errors with multiple rows",
			args: args{
				in: []s1{
					{
						Name: "Ted Tester",
						Age:  28,
					},
					{
						Name: "Heinz",
						Age:  33,
					},
				},
			},
			wantErr: false,
			wantResult: `
|    NAME    | AGE |
|------------|-----|
| Ted Tester |  28 |
| Heinz      |  33 |
`,
		},
		{
			name: "Test write table without errors with multiple pointer rows",
			args: args{
				in: []*s1{
					{
						Name: "Ted Tester",
						Age:  28,
					},
					{
						Name: "Heinz",
						Age:  33,
					},
				},
			},
			wantErr: false,
			wantResult: `
|    NAME    | AGE |
|------------|-----|
| Ted Tester |  28 |
| Heinz      |  33 |
`,
		},
		{
			name: "Test write table without errors and with custom headers",
			args: args{
				in: s2{
					Name: "Ted Tester",
					Age:  28,
				},
			},
			wantErr: false,
			wantResult: `
| FULL NAME  | AGE IN YEARS |
|------------|--------------|
| Ted Tester |           28 |
`,
		},
		{
			name: "Test write table without errors and with private field",
			args: args{
				in: s3{
					Name:        "Ted Tester",
					Age:         28,
					privateCity: "Munich",
				},
			},
			wantErr: false,
			wantResult: `
|    NAME    | AGE |
|------------|-----|
| Ted Tester |  28 |
`,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bldr := new(strings.Builder)

			// hack to be able to format expected strings pretty
			bldr.WriteByte('\n')
			tw := format.Writer("table", bldr)
			if err := tw.Write(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bldr.String() != tt.wantResult {
				t.Errorf("Write() got = %s, want %s", bldr.String(), tt.wantResult)
			}
		})
	}
}
