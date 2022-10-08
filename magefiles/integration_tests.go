package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/magefile/mage/mg"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/multierr"
)

const defaultDirPermissions = 0o755

var (
	workingDir string
	outDir     string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	workingDir = wd
	outDir = filepath.Join(workingDir, "out")

	if err = os.MkdirAll(outDir, defaultDirPermissions); err != nil {
		panic(err)
	}
}

func IntegrationTests(ctx context.Context) error {
	mg.Deps(SnapshotBuild)

	imctl, err := lookupImCtl()
	if err != nil {
		return err
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image: "code.icb4dc0.de/inetmock/inetmock:latest",
			ExposedPorts: []string{
				"80/tcp",
				"443/tcp",
				"53/udp",
			},
			SkipReaper: parseBoolFromEnv("DISABLE_REAPER"),
			Privileged: true,
			WaitingFor: wait.ForExec([]string{
				"/usr/lib/inetmock/bin/imctl",
				"health",
				"container",
			}),
		},
		Started: true,
	})
	if err != nil {
		return err
	}

	defer func() {
		if termErr := container.Terminate(ctx); err != nil {
			err = multierr.Append(err, termErr)
		}
	}()

	containerHost, err := container.Host(ctx)
	if err != nil {
		return err
	}

	dnsPort, err := container.MappedPort(ctx, "53/udp")
	if err != nil {
		return err
	}

	httpPort, err := container.MappedPort(ctx, "80/tcp")
	if err != nil {
		return err
	}

	httpsPort, err := container.MappedPort(ctx, "443/tcp")
	if err != nil {
		return err
	}

	imctlCmd := exec.CommandContext(
		ctx,
		imctl,
		"check",
		"run",
		"--insecure",
		"--dns-proto=udp",
		fmt.Sprintf("--dns-port=%d", dnsPort.Int()),
		fmt.Sprintf("--http-port=%d", httpPort.Int()),
		fmt.Sprintf("--https-port=%d", httpsPort.Int()),
		fmt.Sprintf("--target=%s", containerHost),
		"--log-level=debug",
	)

	scriptBytes, err := os.ReadFile(filepath.Join(workingDir, "testdata", "integration.imcs"))
	if err != nil {
		return err
	}

	imctlCmd.Stdin = bytes.NewReader(scriptBytes)
	imctlCmd.Stdout = os.Stdout
	imctlCmd.Stderr = os.Stderr

	return imctlCmd.Run()
}

func lookupImCtl() (string, error) {
	matches, err := filepath.Glob(filepath.Join(workingDir, "dist", "imctl_linux_amd64*", "imctl"))
	if err != nil {
		return "", err
	}

	if len(matches) < 1 {
		return "", errors.New("imctl not found")
	}

	return matches[0], nil
}

func parseBoolFromEnv(envName string) bool {
	val := os.Getenv(envName)
	if val == "" {
		return false
	}

	if parsed, err := strconv.ParseBool(val); err != nil {
		return false
	} else {
		return parsed
	}
}
