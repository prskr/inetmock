package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"code.gitea.io/sdk/gitea"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"go.uber.org/multierr"
)

func Lint(ctx context.Context) {
	mg.CtxDeps(ctx, Generate)
	mg.CtxDeps(ctx, Format, LintGo, LintProtobuf)
}

func LintProtobuf(ctx context.Context) (err error) {
	status := commitStatusOption("concourse-ci/lint/protobuf", "Lint protobufs")
	if err = setCommitStatus(ctx, status); err != nil {
		return err
	}

	defer func() {
		if err == nil {
			status.State = gitea.StatusSuccess
		} else {
			status.State = gitea.StatusFailure
		}

		err = multierr.Append(err, setCommitStatus(ctx, status))
	}()

	bufCmd := exec.CommandContext(ctx, "buf", "lint")

	bufCmd.Stdout = os.Stdout
	bufCmd.Dir = filepath.Join(WorkingDir, "api", "proto")

	return bufCmd.Run()
}

func LintGo(ctx context.Context) (err error) {
	status := commitStatusOption("concourse-ci/lint/golangci-lint", "Lint Go files")
	if err = setCommitStatus(ctx, status); err != nil {
		return err
	}

	defer func() {
		if err == nil {
			status.State = gitea.StatusSuccess
		} else {
			status.State = gitea.StatusFailure
		}

		err = multierr.Append(err, setCommitStatus(ctx, status))
	}()

	return sh.RunV(
		"golangci-lint",
		"run",
		"-v",
		"--issues-exit-code=1",
	)
}
