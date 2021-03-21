package db

import (
	"github.com/NekoWheel/NekoCAS/internal/conf"
	"github.com/go-redis/redis/v7"

	"github.com/pkg/errors"
)

var red *redis.Client

func ConnRedis() error {
	red = redis.NewClient(&redis.Options{
		Addr:     conf.Get().Redis.Addr,
		Password: conf.Get().Redis.Password,
		DB:       0,
	})

	cb := red.Ping()
	if err := cb.Err(); err != nil {
		return errors.Wrap(err, "ping")
	}
	return nil
}
