package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

type config struct {
	Port int    `toml:"httpport"`
	ICP  string `toml:"icp"`
	Salt string `toml:"salt"`
	Key  string `toml:"key"`

	DBUser     string `toml:"db_user"`
	DBPassword string `toml:"db_password"`
	DBAddr     string `toml:"db_addr"`
	DBName     string `toml:"db_name"`

	RedisAddr     string `toml:"redis_addr"`
	RedisPassword string `toml:"redis_password"`
}

func (cas *cas) initConfig() {
	conf := new(config)
	_, err := toml.DecodeFile("./conf/nekocas.toml", &conf)
	if err != nil {
		log.Fatalln(err)
	}
	cas.Conf = conf
}
