package db

import (
	"fmt"

	"github.com/NekoWheel/NekoCAS/config"
	"github.com/go-redis/redis/v7"
	"github.com/thanhpk/randstr"
	log "unknwon.dev/clog/v2"
)

var red *redis.Client

func ConnRedis() {
	red = redis.NewClient(&redis.Options{
		Addr:     config.Get().RedisAddr,
		Password: config.Get().RedisPassword,
		DB:       0,
	})

	cb := red.Ping()
	if cb.Err() != nil {
		log.Fatal("Failed to ping Redis server: %v", cb.Err())
	}
}

func NewServiceTicket(service *Service, user *User) (string, error) {
	ticket := randstr.String(32)
	err := red.Set(ticket, fmt.Sprintf("%d|%d", service.ID, user.ID), -1).Err()
	if err != nil {
		return "", err
	}
	return ticket, nil
}
