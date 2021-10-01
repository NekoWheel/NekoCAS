package context

import (
	"net/http"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoCAS/internal/conf"
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/form"
	"github.com/NekoWheel/NekoCAS/internal/helper"
	"github.com/NekoWheel/NekoCAS/internal/web/template"
)

// Context 为请求上下文。
type Context struct {
	*macaron.Context
	Cache   cache.Cache
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store
	Setting *setting

	User       *db.User
	IsLogged   bool
	Service    *db.Service
	ServiceURL string
}

type setting struct {
	OpenRegister  string
	SiteLogo      string
	MailWhitelist string
	Privacy       string
}

// Success 返回模板，状态码 200。
func (c *Context) Success(name string) {
	c.HTML(http.StatusOK, name)
}

// Error 返回模板错误页。
func (c *Context) Error(err error) {
	c.Data["ErrorMsg"] = err
	c.HTML(http.StatusOK, "error")
}

// RenderWithErr 返回表单报错。
func (c *Context) RenderWithErr(msg, tpl string, f interface{}) {
	if f != nil {
		form.Assign(f, c.Data)
	}
	c.Flash.ErrorMsg = msg
	c.Data["Flash"] = c.Flash
	c.HTML(http.StatusOK, tpl)
}

// HasError 返回表单验证是否有错误。
func (c *Context) HasError() bool {
	hasErr, ok := c.Data["HasError"]
	if !ok {
		return false
	}
	c.Flash.ErrorMsg = c.Data["ErrorMsg"].(string)
	c.Data["Flash"] = c.Flash
	return hasErr.(bool)
}

// Contexter 初始化一个请求上下文实例。
func Contexter() macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store, f *session.Flash, x csrf.CSRF, cache cache.Cache) {
		c := &Context{
			Context: ctx,
			Cache:   cache,
			csrf:    x,
			Flash:   f,
			Session: sess,
			Setting: &setting{},
		}

		// 获取登录用户信息
		c.User = authenticatedUser(c.Session)

		if c.User != nil {
			c.IsLogged = true
			c.Data["LoggedUser"] = c.User
			c.Data["IsAdmin"] = c.User.IsAdmin

			// 检查用户账号是否已激活
			if !c.User.IsActive &&
				ctx.Req.URL.Path != "/activate" &&
				ctx.Req.URL.Path != "/activate_code" &&
				ctx.Req.URL.Path != "/logout" { // 允许未激活用户登出
				c.Redirect("/activate")
				return
			}
			// 账号已激活
			if c.User.IsActive {
				if ctx.Req.URL.Path == "/activate" || ctx.Req.URL.Path == "/activate_code" {
					c.Redirect("/")
				}
			}
		}

		// 站点设置
		c.Setting = &setting{
			OpenRegister:  db.MustGetSetting("open_register", "on"),
			SiteLogo:      db.MustGetSetting("site_logo", "https://cas.n3ko.co/static/NekoWheel.png"),
			MailWhitelist: db.MustGetSetting("mail_whitelist"),
			Privacy:       db.MustGetSetting("privacy"),
		}
		c.Data["Setting"] = c.Setting

		// 后台菜单
		c.Data["Tab"] = c.Flash.Get("Tab")

		c.Data["SiteName"] = conf.Site.Name
		c.Data["CommitSha"] = helper.Substr(conf.CommitSHA, 0, 8)
		c.Data["CommitLink"] = "https://github.com/NekoWheel/NekoCAS/commit/" + conf.CommitSHA
		c.Data["CSRFToken"] = x.GetToken()
		c.Data["CSRFTokenHTML"] = template.Safe(`<input type="hidden" name="_csrf" value="` + x.GetToken() + `">`)
		log.Trace("Session ID: %s", sess.ID())
		log.Trace("CSRF Token: %v", c.Data["CSRFToken"])

		ctx.Map(c)
	}
}
