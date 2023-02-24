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

	containerRequest := tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image: imageReference.Name(),
			Cmd: []string{
				"serve",
				"--config",
				"/etc/inetmock/config.yaml",
			},
			User: "root",
			Mounts: tc.Mounts(
				tc.BindMount(
					filepath.Join(WorkingDir, "config-container.yaml"),
					"/etc/inetmock/config.yaml",
				),
				tc.BindMount(
					filepath.Join(WorkingDir, "assets", "fakeFiles"),
					"/var/lib/inetmock/fakeFiles",
				),
				tc.BindMount(
					filepath.Join(WorkingDir, "assets", "demoCA"),
					"/var/lib/inetmock/ca",
				),
				tc.BindMount("/sys/kernel/debug", "/sys/kernel/debug"),
				tc.ContainerMount{
					Target: "/var/lib/inetmock/data",
					Source: tc.DockerTmpfsMountSource{},
				},
			),
			ExposedPorts: []string{
				"53/udp",
				"80/tcp",
				"443/tcp",
			},
			Privileged: true,
			WaitingFor: wait.ForLog("App startup completed"),
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

	var (
		dns   nat.Port = "53/udp"
		http  nat.Port = "80/tcp"
		https nat.Port = "443/tcp"
	)

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

func mappedPorts(ctx context.Context, container tc.Container, ports ...nat.Port) (mapped map[nat.Port]nat.Port, err error) {
	mapped = make(map[nat.Port]nat.Port, len(ports))
	for i := range ports {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		p := ports[i]
		if mapped[p], err = container.MappedPort(ctx, p); err != nil {
			return nil, err
		}
	}

	return mapped, nil
}
