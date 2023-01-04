package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Lint(ctx context.Context) {
	mg.CtxDeps(ctx, Generate)
	mg.CtxDeps(ctx, Format)
	mg.CtxDeps(ctx, LintGo, LintProtobuf)
}

func LintProtobuf(ctx context.Context) (err error) {
	bufCmd := exec.CommandContext(ctx, "buf", "lint")

	bufCmd.Stdout = os.Stdout
	bufCmd.Dir = filepath.Join(WorkingDir, "api", "proto")

	return bufCmd.Run()
}

func LintGo(ctx context.Context) (err error) {
	return sh.RunV(
		"golangci-lint",
		"run",
		"-v",
		"--issues-exit-code=1",
	)
}
