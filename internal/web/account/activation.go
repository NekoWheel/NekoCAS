package account

import (
	"github.com/go-macaron/cache"
	"github.com/unknwon/com"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoCAS/internal/context"
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/mail"
)

func ActivationViewHandler(c *context.Context, cache cache.Cache) {
	key := "Activate_Mail_" + com.ToStr(c.User.ID)
	if !cache.IsExist(key) {
		code := c.User.GetActivationCode()
		err := mail.SendActivationMail(c.User.Email, code)
		if err != nil {
			log.Error("Failed to send activation email to %q with error: %v", c.User.Email, err)
		}

		_ = cache.Put(key, true, 120)
	}

	c.Success("activate")
}

func ActivationActionHandler(c *context.Context, cache cache.Cache) {
	key := "Activate_Mail_" + com.ToStr(c.User.ID)
	if !cache.IsExist(key) {
		code := c.User.GetActivationCode()
		err := mail.SendActivationMail(c.User.Email, code)
		if err != nil {
			log.Error("Failed to send activation email to %q with error: %v", c.User.Email, err)
			c.RenderWithErr("服务内部错误，发送邮件失败！", "activate", nil)
			return
		}
		
		_ = cache.Put(key, 1, 120)
	} else {
		c.Flash.Error("邮件发送过于频繁，请等待 2 分钟后再尝试。")
	}
	c.Redirect("/activate")
}

func VerifyUserActiveCodeHandler(c *context.Context) {
	code := c.QueryTrim("code")
	if code == "" {
		c.Redirect("/")
		return
	}

	defer c.Redirect("/login")

	user := db.VerifyUserActiveCode(code)
	if user == nil {
		c.Flash.Error("账号激活码无效。")
		return
	}

	err := db.UpdateUserProfile(&db.User{
		Model: gorm.Model{
			ID: user.ID,
		},
		IsActive: true,
	})
	if err != nil {
		c.Flash.Error("账号激活失败。")
		return
	}

	c.Flash.Success("账号激活成功，请登录。")
}
