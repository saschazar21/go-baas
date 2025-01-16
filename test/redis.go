package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	IMAGE_NAME = "redis/redis-stack:latest"
	HOSTNAME   = "redis-testcontainer"
	REDIS_PORT = 6379

	REDIS_PORT_ENV = "REDIS_PORT"
	REDIS_URL_ENV  = "REDIS_URL"

	CWD = "../"
)

func CreateContainer(ctx context.Context, t *testing.T) (container *redis.RedisContainer, err error) {
	containerReq := testcontainers.ContainerRequest{
		Image:        IMAGE_NAME,
		ExposedPorts: []string{fmt.Sprintf("%d", REDIS_PORT)},
		Hostname:     HOSTNAME,
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	}

	c, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		t.Error(err)
	}

	container = &redis.RedisContainer{
		Container: c,
	}

	var uri string
	if uri, err = container.ConnectionString(ctx); err != nil {
		t.Error(err)
	}

	t.Setenv(REDIS_URL_ENV, uri)

	var p nat.Port
	if p, err = container.MappedPort(ctx, nat.Port(fmt.Sprintf("%d", REDIS_PORT))); err != nil {
		t.Error(err)
	}

	t.Setenv(REDIS_PORT_ENV, p.Port())

	return
}

func TerminateContainer(container *redis.RedisContainer, t *testing.T) {
	if err := container.Terminate(context.Background()); err != nil {
		t.Error(err)
	}
}
