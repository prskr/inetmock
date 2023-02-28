package integration

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"inetmock.icb4dc0.de/inetmock/internal/test"
)

func SetupINetMockContainer(ctx context.Context, exposedPorts ...string) (testcontainers.Container, error) {
	var (
		repoRoot string
		err      error
	)

	if repoRoot, err = test.DiscoverRepoRoot(); err != nil {
		return nil, err
	}

	tcpPortPresent := false
	for _, port := range exposedPorts {
		if strings.Contains(port, "tcp") {
			tcpPortPresent = true
		}
	}

	if !tcpPortPresent {
		exposedPorts = append(exposedPorts, "80/tcp")
	}

	printBuildLog, _ := strconv.ParseBool(os.Getenv("INETMOCK_PRINT_CONTAINER_BUILD_LOG"))

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       repoRoot,
			Dockerfile:    filepath.Join(".", "testdata", "integration.dockerfile"),
			PrintBuildLog: printBuildLog,
		},
		Privileged:   true,
		ExposedPorts: exposedPorts,
		Mounts: testcontainers.Mounts(testcontainers.ContainerMount{
			Source:   testcontainers.DockerBindMountSource{HostPath: "/sys"},
			Target:   "/sys",
			ReadOnly: true,
		}),
		WaitingFor: wait.ForLog("App startup completed"),
	}

	var imContainer testcontainers.Container
	imContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	return imContainer, nil
}
