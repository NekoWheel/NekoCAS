package context

import (
	"net/url"

	"github.com/NekoWheel/NekoCAS/db"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(c *Context) {
		// 已登录用户尝试访问未登录页面
		if options.SignOutRequired && c.IsLogged && c.Req.RequestURI != "/" {
			c.Redirect("/")
			return
		}

		// 未授权访问
		if options.SignInRequired {
			if !c.IsLogged {
				c.SetCookie("redirect_to", url.QueryEscape(c.Req.RequestURI))
				c.Redirect("/login")
				return
			}
		}

		// 管理员访问
		if options.AdminRequired {
			if !c.User.IsAdmin {
				c.Redirect("/")
				return
			}
		}
	}
}

func authenticatedUser(sess session.Store) *db.User {
	id := sess.Get("uid")
	if id == nil {
		return nil
	}
	uid := id.(uint)

	if uid == 0 {
		return nil
	}
	return db.GetUserByID(uid)
}
