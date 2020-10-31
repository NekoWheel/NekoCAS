package config

import (
	log "unknwon.dev/clog/v2"

	"github.com/BurntSushi/toml"
)

type config struct {
	Port       int    `toml:"port"`
	Salt       string `toml:"salt"`
	SessionKey string `toml:"session_key"`
	CSRFKey    string `toml:"csrf_key"`

	DBUser     string `toml:"db_user"`
	DBPassword string `toml:"db_password"`
	DBAddr     string `toml:"db_addr"`
	DBName     string `toml:"db_name"`

	RedisAddr     string `toml:"redis_addr"`
	RedisPassword string `toml:"redis_password"`
}

var conf *config

func init() {
	conf = new(config)
	_, err := toml.DecodeFile("./nekocas.toml", &conf)
	if err != nil {
		log.Fatal("Failed to decode config file: %v", err)
	}
}

// Get returns the config struct.
func Get() *config {
	return conf
}
