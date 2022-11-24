package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"code.gitea.io/sdk/gitea"
	"github.com/magefile/mage/mg"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func IntegrationTests(ctx context.Context) (err error) {
	mg.CtxDeps(ctx, BuildInetmock)

	notification := commitStatusOption("concourse-ci/test/integration", "integration tests")
	if err = setCommitStatus(ctx, notification); err != nil {
		return err
	}

	defer func() {
		if err == nil {
			notification.State = gitea.StatusSuccess
		} else {
			notification.State = gitea.StatusFailure
		}

		err = multierr.Append(err, setCommitStatus(ctx, notification))
	}()

	//nolint:gosec // that's alright
	inetmockCmd := exec.Command(
		filepath.Join(OutDir, "inetmock"),
		"serve",
		"--config=testdata/config-integration.yaml",
	)

	stderrPipe, err := inetmockCmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = inetmockCmd.Start(); err != nil {
		return err
	}

	go func() {
		if err = inetmockCmd.Wait(); err != nil {
			zap.L().Error("Error occurred while waiting for inetmock process", zap.Error(err))
		}
	}()

	stdoutScanner := bufio.NewScanner(stderrPipe)
	stdoutScanner.Split(bufio.ScanLines)

	const sleepDuration = 100 * time.Millisecond

	for {
		if stdoutScanner.Scan() {
			if strings.Contains(stdoutScanner.Text(), "App startup completed") {
				break
			} else {
				fmt.Println(stdoutScanner.Text())
			}
		} else if inetmockCmd.ProcessState != nil && inetmockCmd.ProcessState.Exited() {
			return fmt.Errorf("inetmock process exited with exit code %d", inetmockCmd.ProcessState.ExitCode())
		} else {
			time.Sleep(sleepDuration)
		}
	}

	defer func() {
		err = multierr.Append(err, inetmockCmd.Process.Signal(syscall.SIGTERM))
	}()

	imctlCmd := exec.Command(
		"go",
		"run",
		"inetmock.icb4dc0.de/inetmock/cmd/imctl",
		"check",
		"run",
		"--insecure",
		"--dns-proto=udp",
		"--dns-port=1053",
		"--http-port=80",
		"--https-port=443",
		"--target=127.0.0.1",
		"--log-level=debug",
	)

	scriptBytes, err := os.ReadFile(filepath.Join(WorkingDir, "testdata", "integration.imcs"))
	if err != nil {
		return err
	}

	imctlCmd.Stdin = bytes.NewReader(scriptBytes)
	imctlCmd.Stdout = os.Stdout
	imctlCmd.Stderr = os.Stderr

	return imctlCmd.Run()
}
