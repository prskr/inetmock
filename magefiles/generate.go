package main

import (
	"context"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Generate(ctx context.Context) error {
	errorGroup, groupCtx := errgroup.WithContext(ctx)

	errorGroup.Go(func() error {
		return ensureGoTool("protoc-gen-go", "google.golang.org/protobuf/cmd/protoc-gen-go", "latest")
	})

	errorGroup.Go(func() error {
		return ensureGoTool("protoc-gen-go-grpc", "google.golang.org/grpc/cmd/protoc-gen-go-grpc", "latest")
	})

	errorGroup.Go(func() error {
		return ensureURLTool(groupCtx, "buf", "https://github.com/bufbuild/buf/releases/latest/download/buf-Linux-x86_64")
	})

	if err := errorGroup.Wait(); err != nil {
		return err
	}

	mg.Deps(GenerateProtobuf, GenerateGo)

	return nil
}

func GenerateProtobuf() error {
	logger := zap.L()

	lastProtobufGeneration, err := target.NewestModTime(GeneratedProtobufFiles...)
	if err != nil {
		return err
	}

	lastProtobufModification, err := target.NewestModTime(ProtobufSourceFiles...)
	if lastProtobufGeneration.After(lastProtobufModification) {
		logger.Info("Skipping unnecessary protobuf generation")
		return nil
	}

	return sh.RunV("buf", "generate")
}

func GenerateGo() error {
	logger := zap.L()

	lastMockGeneration, err := target.NewestModTime(GeneratedMockFiles...)
	if err != nil {
		return err
	}

	lastSourceModification, err := target.NewestModTime(GoSourceFiles...)
	if err != nil {
		return err
	}

	if lastMockGeneration.After(lastSourceModification) {
		logger.Info("Skipping unnecessary 'go generate' invocation")
		return nil
	}

	return sh.RunV("go", "generate", "-x", "./...")
}
