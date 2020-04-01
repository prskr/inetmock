package path

import (
	"os"
	"path/filepath"
)

func WorkingDirectory() (cwd string) {
	cwd, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	return
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
