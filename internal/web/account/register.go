package account

import (
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/mail"
	"github.com/NekoWheel/NekoCAS/internal/web/context"
	"github.com/NekoWheel/NekoCAS/internal/web/form"
	"github.com/go-macaron/cache"
	log "unknwon.dev/clog/v2"
)

func RegisterViewHandler(c *context.Context) {
	c.Success("register")
}

func RegisterActionHandler(c *context.Context, f form.Register, cache cache.Cache) {
	if c.Setting.OpenRegister != "on" {
		c.RenderWithErr("当前未开放注册", "register", &f)
		return
	}
	
	if c.HasError() {
		c.Success("register")
		return
	}

	if f.Password != f.Retype {
		c.RenderWithErr("两次输入的密码不匹配", "register", &f)
		return
	}

	u := &db.User{
		NickName: f.NickName,
		Email:    f.Mail,
		Password: f.Password,
	}

	users := db.GetUsers(0, 1)
	if len(users) == 0 {
		// 第一个注册的用户设置成管理员。
		u.IsAdmin = true
		log.Info("Set %q as admin", f.Mail)
	}

	if err := db.CreateUser(u); err != nil {
		switch {
		case db.IsErrUserAlreadyExist(err), db.IsErrEmailAlreadyUsed(err), db.IsErrNameNotAllowed(err):
			c.RenderWithErr(err.Error(), "register", &f)
		default:
			c.RenderWithErr(err.Error(), "register", &f)
		}
		return
	}
	log.Trace("Account created: %s", u.Email)

	// 发送账号激活邮件
	code := u.GetActivationCode()
	go func() {
		err := mail.SendActivationMail(u.Email, code)
		if err != nil {
			log.Error("Failed to send activation email to %q with error %v", u.Email, err)
		}
	}()

	c.Flash.Success("注册成功！")
	c.Redirect("/login")
}
