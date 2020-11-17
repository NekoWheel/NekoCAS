package helper

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// HashEmail 将邮箱地址 Md5 转换成 Avatar 头像哈希
// https://en.gravatar.com/site/implement/hash/
func HashEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	h := md5.New()
	_, _ = h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}
