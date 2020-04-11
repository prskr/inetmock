package path

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
)

func TestFileExists(t *testing.T) {

	tmpFile, err := ioutil.TempFile("", "inetmock")

	if err != nil {
		t.Errorf("failed to create temp file: %v", err)
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Ensure temp file exists",
			want: true,
			args: args{
				filename: tmpFile.Name(),
			},
		},
		{
			name: "Ensure random file name does not exist",
			want: false,
			args: args{
				filename: path.Join(os.TempDir(), fmt.Sprintf("asdf-%d", rand.Uint32())),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.filename); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
