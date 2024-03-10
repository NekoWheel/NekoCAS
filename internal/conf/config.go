package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// CommitSHA 在编译时注入，为当前 Git Commit 哈希值。
var CommitSHA = "debug"

var Site SiteSegment

type SiteSegment struct {
	Name        string `toml:"name"`
	BaseURL     string `toml:"base_url"`
	Port        int    `toml:"port"`
	ICP         string `toml:"icp"`
	SecurityKey string `toml:"security_key"`
	CSRFKey     string `toml:"csrf_key"`
}

var MySQL MySQLSegment

type MySQLSegment struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Addr     string `toml:"addr"`
	Name     string `toml:"name"`
}

var Redis RedisSegment

type RedisSegment struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

var Mail MailSegment

type MailSegment struct {
	Account  string `toml:"account"`
	Password string `toml:"password"`
	SMTP     string `toml:"smtp"`
	Port     int    `toml:"port"`
}

var Ldap LdapSegment

type LdapSegment struct {
	Enabled      bool   `toml:"enable"`
	URL          string `toml:"url"`
	BindDN       string `toml:"bind_dn"`
	BindPassword string `toml:"bind_password"`
	UserFilter   string `toml:"user_filter"`
	BaseDN       string `toml:"base_dn"`
	SyncInterval string `toml:"sync_interval"`
	Mapping      struct {
		Nickname string `toml:"nickname"`
		Email    string `toml:"email"`
		Avatar   string `toml:"avatar"`
	} `toml:"mapping"`
}

// Load 从配置文件中加载配置。
func Load() error {
	var config struct {
		Site  SiteSegment  `toml:"site"`
		MySQL MySQLSegment `toml:"mysql"`
		Redis RedisSegment `toml:"redis"`
		Mail  MailSegment  `toml:"mail"`
		Ldap  LdapSegment  `toml:"ldap"`
	}

	_, err := toml.DecodeFile("./config/nekocas.toml", &config)
	if err != nil {
		return errors.Wrap(err, "decode config file")
	}

	Site = config.Site
	MySQL = config.MySQL
	Redis = config.Redis
	Mail = config.Mail
	Ldap = config.Ldap

	return nil
}
