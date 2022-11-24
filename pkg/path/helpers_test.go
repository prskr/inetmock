package path_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/pkg/path"
)

func TestFileExists(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name string
		want bool
	}
	tests := []testCase{
		{
			name: "Ensure temp file exists",
			want: true,
		},
		{
			name: "Ensure random file name does not exist",
			want: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "nonexistent")
			if tt.want {
				tmpFile, err := os.CreateTemp(tempDir, "testFileExists.*.tmp")
				if td.CmpNoError(t, err) {
					t.Cleanup(func() {
						_ = tmpFile.Close()
					})
					filePath = filepath.Join(tempDir, filepath.Base(tmpFile.Name()))
				}
			}
			if got := path.FileExists(filePath); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
