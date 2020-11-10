package conf

import (
	log "unknwon.dev/clog/v2"

	"github.com/BurntSushi/toml"
)

type config struct {
	Port       int    `toml:"port"`
	Salt       string `toml:"salt"`
	SessionKey string `toml:"session_key"`
	CSRFKey    string `toml:"csrf_key"`

	MySQL struct {
		User     string `toml:"user"`
		Password string `toml:"password"`
		Addr     string `toml:"addr"`
		Name     string `toml:"name"`
	} `toml:"mysql"`

	Redis struct {
		Addr     string `toml:"addr"`
		Password string `toml:"password"`
	} `toml:"redis"`
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
