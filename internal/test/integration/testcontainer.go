package integration

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const containerShutdownTimeout = 5 * time.Second

func SetupINetMockContainer(ctx context.Context, tb testing.TB, exposedPorts ...string) (testcontainers.Container, error) {
	//nolint:dogsled
	_, fileName, _, _ := runtime.Caller(0)

	var err error
	var repoRoot string
	if repoRoot, err = filepath.Abs(filepath.Join(filepath.Dir(fileName), "..", "..", "..")); err != nil {
		return nil, err
	}

	var waitStrategies []wait.Strategy

	var tcpPortPresent = false
	for _, port := range exposedPorts {
		if strings.Contains(port, "tcp") {
			tcpPortPresent = true
			waitStrategies = append(waitStrategies, wait.ForListeningPort(nat.Port(port)))
		}
	}

	if !tcpPortPresent {
		exposedPorts = append(exposedPorts, "80/tcp")
		waitStrategies = append(waitStrategies, wait.ForListeningPort("80/tcp"))
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    repoRoot,
			Dockerfile: filepath.Join(".", "testdata", "integration.dockerfile"),
		},
		ExposedPorts: exposedPorts,
		WaitingFor:   wait.ForAll(waitStrategies...),
	}

	var imContainer testcontainers.Container
	imContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	tb.Cleanup(func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), containerShutdownTimeout)
		defer cancel()
		_ = imContainer.Terminate(shutdownCtx)
	})

	return imContainer, nil
}
