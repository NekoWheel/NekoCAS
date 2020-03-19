package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
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

func (cas *cas) md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
