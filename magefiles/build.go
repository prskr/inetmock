package main

import (
	"context"
	"path/filepath"

	"github.com/magefile/mage/mg"
)

func BuildInetmock(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate)

	return GoBuild("-o", filepath.Join(OutDir, "inetmock"), "-trimpath", "inetmock.icb4dc0.de/inetmock/cmd/inetmock")
}

func SnapshotBuild(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate)

	if err := ensureGoTool("goreleaser", "github.com/goreleaser/goreleaser", "latest"); err != nil {
		return err
	}

	return GoReleaser("release", "--snapshot", "--skip-publish", "--rm-dist")
}
