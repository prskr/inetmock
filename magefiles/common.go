package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/magefile/mage/sh"
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
	GitCommit              string
	GiteaClient            *gitea.Client
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

	if currentCommit, err := sh.Output("git", "rev-parse", "HEAD"); err != nil {
		panic(err)
	} else {
		GitCommit = currentCommit
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

	if giteaToken := os.Getenv("GITEA_TOKEN"); giteaToken != "" {
		if client, err := gitea.NewClient("https://code.icb4dc0.de", gitea.SetToken(giteaToken)); err == nil {
			GiteaClient = client
		}
	}

	zap.L().Info("Completed initialization", zap.String("commit", GitCommit))
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

func commitStatusOption(context, description string) gitea.CreateStatusOption {
	return gitea.CreateStatusOption{
		Context:     context,
		Description: description,
		State:       gitea.StatusPending,
		TargetURL:   "https://concourse.icb4dc0.de/teams/inetmock/pipelines",
	}
}
