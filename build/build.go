//go:build mage

package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/publish"
	"github.com/magefile/mage/mg"
	"go.uber.org/zap"
)

var (
	imageBuildResult build.Result
	imageReference   name.Reference
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

func BuildImage(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate)

	log.SetOutput(nil)

	ko, err := build.NewGo(
		ctx,
		WorkingDir,
		build.WithPlatforms("linux/amd64"),
		build.WithBaseImages(getBaseImage("gcr.io/distroless/static:nonroot")),
	)
	if err != nil {
		return err
	}

	imageBuildResult, err = ko.Build(ctx, "./cmd/inetmock")
	if err != nil {
		return err
	}

	staticNamer := func(s1, s2 string) string {
		return "code.icb4dc0.de/inetmock/inetmock/base"
	}

	daemonPublisher, err := publish.NewDaemon(staticNamer, []string{"latest"})
	if err != nil {
		return err
	}

	imageReference, err = daemonPublisher.Publish(ctx, imageBuildResult, "")
	if err != nil {
		return err
	}

	zap.L().Info("Built image", zap.String("ref", imageReference.Name()))

	return nil
}

func getBaseImage(imageRef string) build.GetBase {
	var cache sync.Map
	fetch := func(ctx context.Context, ref name.Reference) (build.Result, error) {
		// For ko.local, look in the daemon.
		if ref.Context().RegistryStr() == publish.LocalDomain {
			return daemon.Image(ref)
		}

		ropt := []remote.Option{
			remote.WithUserAgent("golang-ko"),
			remote.WithContext(ctx),
		}

		desc, err := remote.Get(ref, ropt...)
		if err != nil {
			return nil, err
		}
		if desc.MediaType.IsIndex() {
			return desc.ImageIndex()
		}
		return desc.Image()
	}
	return func(ctx context.Context, s string) (name.Reference, build.Result, error) {
		s = strings.TrimPrefix(s, build.StrictScheme)
		ref, err := name.ParseReference(imageRef, name.WithDefaultRegistry("docker.io"))
		if err != nil {
			return nil, nil, fmt.Errorf("parsing base image (%q): %w", imageRef, err)
		}

		if v, ok := cache.Load(ref.String()); ok {
			return ref, v.(build.Result), nil
		}

		result, err := fetch(ctx, ref)
		if err != nil {
			return ref, result, err
		}

		if _, ok := ref.(name.Digest); ok {
			log.Printf("Using base %s for %s", ref, s)
		} else {
			dig, err := result.Digest()
			if err != nil {
				return ref, result, err
			}
			log.Printf("Using base %s@%s for %s", ref, dig, s)
		}

		cache.Store(ref.String(), result)
		return ref, result, nil
	}
}
