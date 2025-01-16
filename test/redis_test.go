package test

import (
	"context"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisContainer(t *testing.T) {
	var err error
	ctx := context.Background()

	// Create a new Redis container
	container, err := CreateContainer(ctx, t)
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		TerminateContainer(container, t)
	})

	var opts *redis.Options
	if opts, err = redis.ParseURL(os.Getenv(REDIS_URL_ENV)); err != nil {
		t.Fatalf("%v", err)
	}

	rdb := redis.NewClient(opts)

	var pong string
	if pong, err = rdb.Ping(ctx).Result(); err != nil {
		t.Fatalf("%v", err)
	}

	assert.Equal(t, "PONG", pong)
}
