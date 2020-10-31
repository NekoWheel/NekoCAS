package account

import (
	"errors"
	"net/url"

	"github.com/NekoWheel/NekoCAS/db"
	"github.com/NekoWheel/NekoCAS/web/context"
	"github.com/NekoWheel/NekoCAS/web/form"
)

func LoginViewHandler(c *context.Context) {
	c.Success("login")
}

func LoginActionHandler(c *context.Context, f form.Login) {
	if c.HasError() {
		c.Success("login")
		return
	}

	u, err := db.UserAuthenticate(f.Mail, f.Password)
	if err != nil {
		c.RenderWithErr(err.Error(), "login", f)
		return
	}

	c.User = u
	_ = c.Session.Set("uid", u.ID)
	_ = c.Session.Set("uname", u.Name)

	// 携带 Ticket 跳转到对应服务
	if c.Service.ID != 0 {
		ticket, err := db.NewServiceTicket(c.Service, c.User)
		if err != nil {
			c.Error(err)
			return
		}
		// 解析跳转 URL
		redirectURL, err := url.Parse(c.ServiceURL)
		if err != nil {
			c.Error(errors.New("解析 URL 失败"))
			return
		}
		query := redirectURL.Query()
		query.Set("ticket", ticket)
		redirectURL.RawQuery = query.Encode()
		c.Redirect(redirectURL.String())
		return
	}

	c.Redirect("/")
}
