package main

import (
	"context"

	"github.com/magefile/mage/mg"
)

func SnapshotBuild(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate)

	return GoReleaser("release", "--snapshot", "--skip-publish", "--rm-dist")
}
