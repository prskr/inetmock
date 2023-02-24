//go:build mage

package main

import (
	"fmt"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	_ "gotest.tools/gotestsum/cmd"
)

func TestShort() error {
	mg.Deps(Generate)

	return sh.RunWithV(
		map[string]string{
			"CGO_ENABLED": "1",
		},
		"go", "run",
		"gotest.tools/gotestsum",
		"--packages=./...",
		"--rerun-fails=5",
		"--",
		"-timeout=5m",
		"-short",
		"-race",
		"-shuffle=on",
		fmt.Sprintf("-coverprofile=%s", filepath.Join(OutDir, "cov-raw.out")),
		"-covermode=atomic",
	)
}

func TestAll() error {
	mg.Deps(Generate)

	return sh.RunWithV(
		map[string]string{
			"CGO_ENABLED": "1",
		},
		"go", "run",
		"gotest.tools/gotestsum",
		"--packages=./...",
		"--rerun-fails=5",
		"--",
		"-timeout=10m",
		"--tags=sudo",
		"-race",
		"-shuffle=on",
		fmt.Sprintf("-coverprofile=%s", filepath.Join(OutDir, "cov-raw.out")),
		"-covermode=atomic",
	)
}
