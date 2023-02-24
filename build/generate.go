//go:build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const compilationCommandFormat = `clang \
-Wno-unused-value \
-Wno-pointer-sign \
-Wno-compare-distinct-pointer-types \
-Wunused \
-Wall \
-fno-stack-protector \
-fno-ident \
-g \
-O2 \
-emit-llvm %s -c -o - | llc -march=bpf -mcpu=probe -filetype=obj -o %s`

func Generate(ctx context.Context) error {
	errorGroup, groupCtx := errgroup.WithContext(ctx)

	errorGroup.Go(func() error {
		return ensureGoTool("mockgen", "github.com/golang/mock/mockgen", "v1.6.0")
	})

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

	mg.Deps(GenerateProtobuf, GenerateGo, CompileEBPF)

	return nil
}

func GenerateProtobuf() error {
	logger := zap.L()

	lastProtobufGeneration, err := target.NewestModTime(GeneratedProtobufFiles...)
	if err != nil {
		return err
	}

	lastProtobufModification, err := target.NewestModTime(ProtobufSourceFiles...)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	logger.Debug("Determined last time protobuf files where modified", zap.Time("lastProtobufModification", lastProtobufModification))

	if lastProtobufGeneration.After(lastProtobufModification) && !GenerateAlways {
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

	logger.Debug("Determined last time mocks where generated", zap.Time("lastMockGeneration", lastMockGeneration))

	if lastMockGeneration.After(lastSourceModification) && !GenerateAlways {
		logger.Info("Skipping unnecessary 'go generate' invocation")
		return nil
	}

	return sh.RunV("go", "generate", "-x", "./...")
}

func CompileEBPF() error {
	return errors.Join(
		compileEBPFTarget("nat"),
		compileEBPFTarget("firewall"),
		compileEBPFTarget("tests"),
	)
}

func compileEBPFTarget(targetName string) error {
	var (
		eBPFSourceDirectory = filepath.Join("netflow", "ebpf")
		outFilePath         = filepath.Join(eBPFSourceDirectory, fmt.Sprintf("%s.o", targetName))
		sourceFilePath      = filepath.Join(eBPFSourceDirectory, fmt.Sprintf("%s.c", targetName))
		logger              = zap.L().With(zap.String("target", targetName))
	)

	if compilationRequired, err := target.Glob(outFilePath, sourceFilePath, filepath.Join(eBPFSourceDirectory, "*.h")); err != nil {
		return err
	} else if compilationRequired || GenerateAlways {
		compilationCmd := fmt.Sprintf(compilationCommandFormat, sourceFilePath, outFilePath)
		logger.Debug("Compile eBPF", zap.String("cmd", compilationCmd))
		return sh.RunV("sh", "-c", compilationCmd)
	}

	logger.Info("Skipping eBPF recompilation")

	return nil
}
