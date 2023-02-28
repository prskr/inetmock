package test

import (
	"errors"
	"os"
	"path/filepath"
)

func DiscoverRepoRoot() (root string, err error) {
	root, err = os.Getwd()
	if err != nil {
		return
	}

	for root != "" {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			return root, nil
		}

		root = filepath.Dir(root)
	}

	return "", errors.New("failed to discover repo root")
}
