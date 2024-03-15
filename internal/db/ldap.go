package db

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NekoWheel/NekoCAS/internal/conf"
	"github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
)

func ldapDialAndBind() (*ldap.Conn, error) {
	c, err := ldap.DialURL(conf.Ldap.URL)
	if err != nil {
		return nil, errors.Wrap(err, "connect ldap")
	}

	err = c.Bind(conf.Ldap.BindDN, conf.Ldap.BindPassword)
	if err != nil {
		return nil, errors.Wrap(err, "auth ldap")
	}

	return c, nil
}

func AutoSyncLdap() {
	dur, err := time.ParseDuration(conf.Ldap.SyncInterval)
	if err != nil {
		log.Println("parse time duration error, ldap sync is not working")
	}

	if err := syncLdap(); err != nil {
		log.Println("sync ldap error:", err)
	}
	t := time.NewTicker(dur)
	for range t.C {
		if err := syncLdap(); err != nil {
			log.Println("sync ldap error:", err)
		}
	}
}

func syncLdap() error {
	c, err := ldapDialAndBind()
	if err != nil {
		return errors.Wrap(err, "sync ldap")
	}

	userFilter := fmt.Sprintf(conf.Ldap.UserFilter, "*")
	req := ldap.NewSearchRequest(conf.Ldap.BaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, userFilter,
		[]string{conf.Ldap.Mapping.Nickname, conf.Ldap.Mapping.Email, conf.Ldap.Mapping.Avatar}, nil)
	res, err := c.Search(req)
	if err != nil {
		return errors.Wrap(err, "search ldap")
	}

	for _, e := range res.Entries {
		email := strings.ToLower(e.GetAttributeValue(conf.Ldap.Mapping.Email))
		nickname := e.GetAttributeValue(conf.Ldap.Mapping.Nickname)
		if !IsEmailUsed(email) && !IsNickNameUsed(nickname) {
			user := &User{
				NickName: nickname,
				Email:    email,
				IsLdap:   true,
				IsActive: true,
			}
			if err := CreateUser(user); err != nil {
				log.Println("create user error:", err)
			}
		}
	}
	return nil
}

func ldapAuthenticate(email string, password string) (bool, error) {
	c, err := ldapDialAndBind()
	if err != nil {
		return false, nil
	}

	filter, err := sanitizedUserQuery(email)
	if err != nil {
		return false, errors.Wrap(err, "sanitize input email")
	}

	req := ldap.NewSearchRequest(conf.Ldap.BaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, nil, nil)
	res, err := c.Search(req)
	if err != nil {
		return false, errors.Wrap(err, "serach ldap")
	}
	if len(res.Entries) != 1 {
		return false, errors.Wrap(err, "entry not unique or not existed")
	}

	userDN := res.Entries[0].DN
	if err := c.Bind(userDN, password); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultAuthorizationDenied) {
			return false, nil
		}
		return false, errors.Wrap(err, "bind DN")
	}

	return true, nil
}

func sanitizedUserQuery(email string) (string, error) {
	// See http://tools.ietf.org/search/rfc4515
	badCharacters := "\x00()*\\"
	if strings.ContainsAny(email, badCharacters) {
		return "", fmt.Errorf("'%s' contains invalid query characters. Aborting.", email)
	}

	return fmt.Sprintf(conf.Ldap.UserFilter, email), nil
}
