package integration

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupINetMockContainer(ctx context.Context, exposedPorts ...string) (testcontainers.Container, error) {
	//nolint:dogsled
	_, fileName, _, _ := runtime.Caller(0)

	var err error
	var repoRoot string
	if repoRoot, err = filepath.Abs(filepath.Join(filepath.Dir(fileName), "..", "..", "..")); err != nil {
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

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       repoRoot,
			Dockerfile:    filepath.Join(".", "testdata", "integration.dockerfile"),
			PrintBuildLog: true,
		},
		ExposedPorts: exposedPorts,
		SkipReaper:   true,
		WaitingFor:   wait.ForLog("Startup of all endpoints completed"),
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
