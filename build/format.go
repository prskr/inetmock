//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Format() {
	mg.Deps(GoImports)
	mg.Deps(GoFumpt)
}

func GoImports() error {
	if err := ensureGoTool("goimports", "golang.org/x/tools/cmd/goimports", "latest"); err != nil {
		return err
	}

	return sh.RunV(
		"goimports",
		"-local=inetmock.icb4dc0.de/inetmock",
		"-w",
		WorkingDir,
	)
}

func GoFumpt() error {
	if err := ensureGoTool("gofumpt", "mvdan.cc/gofumpt", "latest"); err != nil {
		return err
	}
	return sh.RunV("gofumpt", "-l", "-w", WorkingDir)
}
