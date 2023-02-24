//go:build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/sh"
	"go.uber.org/zap"
	"golang.org/x/net/context/ctxhttp"
)

var (
	GoReleaser = sh.RunCmd("goreleaser")
	GoInstall  = sh.RunCmd("go", "install")
	GoBuild    = sh.RunCmd("go", "build")
)

func ensureURLTool(ctx context.Context, toolName, downloadURL string) error {
	return checkForTool(toolName, func() error {
		resp, err := ctxhttp.Get(ctx, http.DefaultClient, downloadURL)
		if err != nil {
			return err
		}

		defer func() {
			err = errors.Join(err, resp.Body.Close())
		}()

		const ownerExecute = 0o755
		outFile, err := os.OpenFile(filepath.Join("/", "usr", "local", "bin", toolName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, ownerExecute)
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, resp.Body)

		return err
	})
}

func ensureGoTool(toolName, importPath, version string) error {
	return checkForTool(toolName, func() error {
		logger := zap.L()
		toolToInstall := fmt.Sprintf("%s@%s", importPath, version)
		logger.Info("Installing Go tool", zap.String("toolToInstall", toolToInstall))
		return GoInstall(toolToInstall)
	})
}

func checkForTool(toolName string, fallbackAction func() error) error {
	logger := zap.L()
	if _, err := exec.LookPath(toolName); err != nil {
		logger.Warn("tool is missing", zap.String("toolName", toolName))
		return fallbackAction()
	}

	return nil
}
