package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slices"
)

var (
	GoSourceFiles          []string
	ProtobufSourceFiles    []string
	GeneratedMockFiles     []string
	GeneratedProtobufFiles []string
	WorkingDir             string
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
	if wd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		WorkingDir = wd
	}

	if err := initLogging(); err != nil {
		panic(err)
	}

	if err := initSourceFiles(); err != nil {
		panic(err)
	}
}

func initLogging() error {
	cfg := zap.NewDevelopmentConfig()
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
