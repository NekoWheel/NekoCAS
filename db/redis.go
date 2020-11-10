package db

import (
	"github.com/NekoWheel/NekoCAS/conf"
	"github.com/go-redis/redis/v7"
	log "unknwon.dev/clog/v2"
)

var red *redis.Client

func ConnRedis() {
	red = redis.NewClient(&redis.Options{
		Addr:     conf.Get().Redis.Addr,
		Password: conf.Get().Redis.Password,
		DB:       0,
	})

	cb := red.Ping()
	if cb.Err() != nil {
		log.Fatal("Failed to ping Redis server: %v", cb.Err())
	}
}
