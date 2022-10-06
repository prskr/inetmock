package main

import "github.com/magefile/mage/sh"

var (
	GoReleaser = sh.RunCmd("goreleaser")
	GoTestSum  = sh.RunCmd("go", "run", "gotest.tools/gotestsum")
)
