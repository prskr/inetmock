//go:build mage

package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slices"
)

const defaultDirPermissions = 0o755

var (
	GoSourceFiles          []string
	ProtobufSourceFiles    []string
	GeneratedMockFiles     []string
	GeneratedProtobufFiles []string
	WorkingDir             string
	OutDir                 string
	GenerateAlways         bool
	dirsToIgnore           = []string{
		".git",
		"magefiles",
		".gitlab",
		".github",
		".run",
		".task",
	}
)

func init() {
	if b, err := strconv.ParseBool(os.Getenv("GENERATE_ALWAYS")); err == nil {
		GenerateAlways = b
	}
	if wd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		WorkingDir = wd
	}

	OutDir = filepath.Join(WorkingDir, "out")

	if err := os.MkdirAll(OutDir, defaultDirPermissions); err != nil {
		panic(err)
	}

	if err := initLogging(); err != nil {
		panic(err)
	}

	if err := initSourceFiles(); err != nil {
		panic(err)
	}

	zap.L().Info("Completed initialization")
}

func initLogging() error {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	cfg.Encoding = "console"
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	if logger, err := cfg.Build(); err != nil {
		return err
	} else {
		zap.ReplaceGlobals(logger)
	}

	return nil
}

func initSourceFiles() error {
	return filepath.WalkDir(WorkingDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if slices.Contains(dirsToIgnore, filepath.Base(path)) {
				return fs.SkipDir
			}
			return nil
		}

		_, ext, found := strings.Cut(filepath.Base(path), ".")
		if !found {
			return nil
		}

		switch ext {
		case "proto":
			ProtobufSourceFiles = append(ProtobufSourceFiles, path)
		case "pb.go":
			GeneratedProtobufFiles = append(GeneratedProtobufFiles, path)
		case "mock.go":
			GeneratedMockFiles = append(GeneratedMockFiles, path)
		case "go":
			GoSourceFiles = append(GoSourceFiles, path)
		}

		return nil
	})
}
