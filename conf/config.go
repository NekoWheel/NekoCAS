package conf

import (
	log "unknwon.dev/clog/v2"

	"github.com/BurntSushi/toml"
)

var COMMIT_SHA = "debug"

type config struct {
	Site struct {
		Name        string `toml:"name"`
		BaseURL     string `toml:"base_url"`
		Port        int    `toml:"port"`
		ICP         string `toml:"icp"`
		SecurityKey string `toml:"security_key"`
		CSRFKey     string `toml:"csrf_key"`
	} `toml:"site"`

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

	Mail struct {
		Account  string `toml:"account"`
		Password string `toml:"password"`
		SMTP     string `toml:"smtp"`
		Port     int    `toml:"port"`
	} `toml:"mail"`
}

var conf *config

func init() {
	conf = new(config)
	_, err := toml.DecodeFile("./config/nekocas.toml", &conf)
	if err != nil {
		log.Fatal("Failed to decode config file: %v", err)
	}
}

// Get returns the config struct.
func Get() *config {
	return conf
}
