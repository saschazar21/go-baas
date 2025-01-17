package db

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

const REDIS_URL_ENV string = "REDIS_URL"

func NewRedis() (rdb *redis.Client, err error) {
	if os.Getenv(REDIS_URL_ENV) == "" {
		log.Fatal("No REDIS_URL env provided.")

		return
	}

	var options *redis.Options
	options, err = redis.ParseURL(os.Getenv(REDIS_URL_ENV))

	if err != nil {
		log.Println(err)

		return
	}

	rdb = redis.NewClient(options)

	return
}
