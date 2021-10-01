package db

import (
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"

	"github.com/NekoWheel/NekoCAS/internal/conf"
)

var red *redis.Client

func ConnRedis() error {
	red = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       0,
	})

	cb := red.Ping()
	if err := cb.Err(); err != nil {
		return errors.Wrap(err, "ping")
	}
	return nil
}
