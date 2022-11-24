package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/sh"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/net/context/ctxhttp"
)

var (
	GoReleaser = sh.RunCmd("goreleaser")
	GoInstall  = sh.RunCmd("go", "install")
)

func ensureURLTool(ctx context.Context, toolName, downloadURL string) error {
	return checkForTool(toolName, func() error {
		resp, err := ctxhttp.Get(ctx, http.DefaultClient, downloadURL)
		if err != nil {
			return err
		}

		defer multierr.AppendInvoke(&err, multierr.Close(resp.Body))

		outFile, err := os.Create(filepath.Join("usr", "local", "bin", toolName))
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, resp.Body)

		return err
	})
}

//nolint:unparam // subject to be changed in the future
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
