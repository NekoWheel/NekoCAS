package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
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

func Load() error {
	c := config{}

	_, err := toml.DecodeFile("./config/nekocas.toml", &c)
	if err != nil {
		return errors.Wrap(err, "decode config file")
	}

	conf = &c
	return nil
}

// Get returns the config struct.
func Get() *config {
	return conf
}
