package account

import (
	"github.com/NekoWheel/NekoCAS/db"
	"github.com/NekoWheel/NekoCAS/web/context"
	"github.com/NekoWheel/NekoCAS/web/form"
	log "unknwon.dev/clog/v2"
)

func RegisterViewHandler(c *context.Context) {
	c.Success("register")
}

func RegisterActionHandler(c *context.Context, f form.Register) {
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

	c.Flash.Success("注册成功！")
	c.Redirect("/login")
}
