package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/thanhpk/randstr"
)

// generateServiceToken generates a random Service Ticket.
func generateServiceToken() string {
	return fmt.Sprintf("ST-%s", randstr.String(32))
}

// HashEmail 将邮箱地址 Md5 转换成 Avatar 头像哈希
// https://en.gravatar.com/site/implement/hash/
func HashEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	h := md5.New()
	_, _ = h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}

//func ValidateResponse(c *gin.Context, ok bool, userName string) {
//	const LF = string(rune(10))
//	if ok {
//		c.Data(http.StatusOK, "text/plain", []byte("yes"+LF+" "+userName+LF))
//		return
//	}
//	c.Data(http.StatusOK, "text/plain", []byte("no"+LF))
//}
