package main

import (
	"fmt"
	"path/filepath"

	_ "gotest.tools/gotestsum/cmd"
)

func TestShort() error {
	return GoTestSum(
		"--packages=./...",
		"--rerun-fails=5",
		"--",
		"-timeout=5m",
		"-short",
		"-race",
		"-shuffle=on",
		fmt.Sprintf("-coverprofile=%s", filepath.Join(outDir, "cov-raw.out")),
		"-covermode=atomic",
	)
}

func TestAll() error {
	return GoTestSum(
		"--packages=./...",
		"--rerun-fails=5",
		"--",
		"-timeout=10m",
		"--tags=sudo",
		"-race",
		"-shuffle=on",
		fmt.Sprintf("-coverprofile=%s", filepath.Join(outDir, "cov-raw.out")),
		"-covermode=atomic",
	)
}
