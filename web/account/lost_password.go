package account

import (
	"fmt"

	"github.com/NekoWheel/NekoCAS/db"
	"github.com/NekoWheel/NekoCAS/mail"
	"github.com/NekoWheel/NekoCAS/web/context"
	"github.com/NekoWheel/NekoCAS/web/form"
	"github.com/go-macaron/cache"
	"github.com/unknwon/com"
)

func LostPasswordHandler(c *context.Context) {
	c.Success("lost_password")
}

func LostPasswordActionHandler(c *context.Context, f form.LostPassword, cache cache.Cache) {
	user := db.GetUserByEmail(f.Email)
	if user == nil {
		c.RenderWithErr("该邮箱不存在", "lost_password", nil)
		return
	}

	c.Flash.Success(fmt.Sprintf("找回密码邮件发送成功，请检查邮箱 %s。", user.Email))

	key := "Lost_Password_" + com.ToStr(user.ID)
	if !cache.IsExist(key) {
		code := user.GetActivationCode()
		err := mail.SendLostPasswordMail(user.Email, code)
		if err != nil {
			c.RenderWithErr("服务内部错误，发送邮件失败！", "lost_password", nil)
			return
		}
		_ = cache.Put(key, true, 120)
	} else {
		c.Flash.Error("邮件发送过于频繁，请等待 2 分钟后再尝试。")
	}

	c.Redirect("/lost_password")
}

func ResetPasswordHandler(c *context.Context) {
	code := c.QueryTrim("code")
	if code == "" {
		c.Redirect("/")
		return
	}

	user := db.VerifyUserActiveCode(code)
	if user == nil {
		c.Flash.Error("重置密码链接无效。")
		c.Redirect("/login")
		return
	}

	c.Data["Email"] = user.Email

	c.Success("reset_password")
}

func ResetPasswordActionHandler(c *context.Context, f form.ResetPassword) {
	code := c.QueryTrim("code")
	if code == "" {
		c.Redirect("/")
		return
	}

	user := db.VerifyUserActiveCode(code)
	if user == nil {
		c.Flash.Error("重置密码链接无效。")
		c.Redirect("/login")
		return
	}

	c.Data["Email"] = user.Email

	// 表单报错
	if c.HasError() {
		c.Success("reset_password")
		return
	}

	if f.Password != f.Retype {
		c.Flash.Error("两次输入的密码不相同。")
		c.Redirect(c.Req.URL.String())
		return
	}
	
	user.Password = f.Password
	err := db.UpdateUserProfile(user)
	if err != nil {
		c.Flash.Error("服务内部错误，密码重置失败。")
		c.Redirect(c.Req.URL.String())
		return
	}

	c.Flash.Success("密码重置成功，请登录。")
	c.Redirect("/login")
}
