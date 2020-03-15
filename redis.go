package main

import "github.com/go-redis/redis/v7"

func (cas *cas) initRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     cas.Conf.RedisAddr,
		Password: cas.Conf.RedisPassword,
		DB:       0,
	})
	cas.Redis = client
}
