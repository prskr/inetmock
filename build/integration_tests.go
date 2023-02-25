//go:build mage

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/go-connections/nat"
	"github.com/magefile/mage/mg"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func IntegrationTests(ctx context.Context) (err error) {
	mg.CtxDeps(ctx, BuildImage)

	baseImage := imageReference.Name()

	var (
		dns   = "1053/udp"
		http  = "80/tcp"
		https = "443/tcp"
	)

	containerRequest := tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			FromDockerfile: tc.FromDockerfile{
				Context:    WorkingDir,
				Dockerfile: filepath.Join("deploy", "docker", "integration.dockerfile"),
				BuildArgs: map[string]*string{
					"BASE_IMAGE": &baseImage,
				},
			},
			Mounts: tc.Mounts(
				tc.ContainerMount{
					Target: "/var/lib/inetmock/data",
					Source: tc.DockerTmpfsMountSource{},
				},
			),
			ExposedPorts: []string{
				dns,
				http,
				https,
			},
			WaitingFor: wait.ForLog("App startup completed"),
			SkipReaper: DisableReaper,
		},
		Started: true,
	}
	container, err := tc.GenericContainer(ctx, containerRequest)
	if err != nil {
		return err
	}

	defer func() {
		_ = container.Terminate(context.Background())
	}()

	mapped, err := mappedPorts(ctx, container, dns, http, https)
	if err != nil {
		return err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return err
	}

	imctlCmd := exec.CommandContext(
		ctx,
		"go",
		"run",
		"inetmock.icb4dc0.de/inetmock/cmd/imctl",
		"check",
		"run",
		"--insecure",
		"--dns-proto=udp",
		fmt.Sprintf("--dns-port=%d", mapped[dns].Int()),
		fmt.Sprintf("--http-port=%d", mapped[http].Int()),
		fmt.Sprintf("--https-port=%d", mapped[https].Int()),
		fmt.Sprintf("--target=%s", host),
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

func mappedPorts(ctx context.Context, container tc.Container, ports ...string) (mapped map[string]nat.Port, err error) {
	mapped = make(map[string]nat.Port, len(ports))
	for i := range ports {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		p := ports[i]
		if mapped[p], err = container.MappedPort(ctx, nat.Port(p)); err != nil {
			return nil, err
		}
	}

	return mapped, nil
}
