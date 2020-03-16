package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"github.com/thanhpk/randstr"
	"io"
	"strings"
)

func (cas *cas) addSalt(raw string) string {
	return cas.hmacSha1Encode(raw, cas.Conf.Salt)
}

func (cas *cas) hmacSha1Encode(input string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	_, _ = io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (cas *cas) generateServiceToken() string {
	return fmt.Sprintf("ST-%s", randstr.String(32))
}

func (cas *cas) getErrorMessage(field string, value ...interface{}) string {
	// value - tag, param...
	fieldName := map[string]string{
		"Mail":     "电子邮箱",
		"Name":     "昵称",
		"Password": "密码",
	}
	tagMessage := map[string]string{
		"required": "%s不能为空",
		"email":    "%s格式不合法",
		"min":      "%s最小长度为%v",
		"max":      "%s最大长度为%v",
	}

	tag := tagMessage[value[0].(string)]

	if len(value) == 0 {
		return tag
	} else if len(value) == 1 {
		// tag, field
		return fmt.Sprintf(tag, fieldName[field])
	}
	value[0] = fieldName[field]
	// TODO here is a not good method, change it later.
	return strings.Split(fmt.Sprintf(tag, value...), "%")[0]
}
