package helper

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/NekoWheel/NekoCAS/conf"
	"github.com/unknwon/com"
)

// HashEmail 将邮箱地址 Md5 转换成 Avatar 头像哈希
// https://en.gravatar.com/site/implement/hash/
func HashEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	h := md5.New()
	_, _ = h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}

const TIME_LIMIT_CODE_LENGTH = 12 + 6 + 40

// CreateTimeLimitCode 根据给定的数据生成带过期时间的码。
// 格式： 12 日期 + 6 分钟 + 40 sha1
func CreateTimeLimitCode(data string, minutes int, startInf interface{}) string {
	format := "200601021504"

	var start, end time.Time
	var startStr, endStr string

	if startInf == nil {
		start = time.Now()
		startStr = start.Format(format)
	} else {
		startStr = startInf.(string)
		start, _ = time.ParseInLocation(format, startStr, time.Local)
		startStr = start.Format(format)
	}

	end = start.Add(time.Minute * time.Duration(minutes))
	endStr = end.Format(format)

	sh := sha1.New()
	_, _ = sh.Write([]byte(data + conf.Get().Site.SecurityKey + startStr + endStr + com.ToStr(minutes)))
	encoded := hex.EncodeToString(sh.Sum(nil))

	code := fmt.Sprintf("%s%06d%s", startStr, minutes, encoded)
	return code
}

// VerifyTimeLimitCode 检查限时验证码是否正确。
func VerifyTimeLimitCode(data string, minutes int, code string) bool {
	if len(code) <= 18 {
		return false
	}

	start := code[:12]
	lives := code[12:18]
	if d, err := com.StrTo(lives).Int(); err == nil {
		minutes = d
	}

	retCode := CreateTimeLimitCode(data, minutes, start)
	if retCode == code && minutes > 0 {
		// 检查是否超时
		before, _ := time.ParseInLocation("200601021504", start, time.Local)
		now := time.Now()
		if before.Add(time.Minute*time.Duration(minutes)).Unix() > now.Unix() {
			return true
		}
	}

	return false
}
