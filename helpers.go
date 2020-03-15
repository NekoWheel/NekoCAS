package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"github.com/thanhpk/randstr"
	"io"
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
