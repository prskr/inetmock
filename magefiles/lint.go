package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
)

func Lint(ctx context.Context) {
	mg.CtxDeps(ctx, LintGo, LintProtobuf)
}

func LintProtobuf(ctx context.Context) error {
	bufCmd := exec.CommandContext(ctx, "buf", "lint")

	bufCmd.Stdout = os.Stdout
	bufCmd.Dir = filepath.Join(workingDir, "api", "proto")

	return bufCmd.Run()
}

func LintGo(context.Context) error {
	return GoLangCiLint(
		"run",
		"-v",
		"--issues-exit-code=1",
	)
}
