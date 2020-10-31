package context

import (
	"net/http"

	"github.com/NekoWheel/NekoCAS/db"
	"github.com/NekoWheel/NekoCAS/web/form"
	"github.com/NekoWheel/NekoCAS/web/template"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"
)

// Context 请求上下文
type Context struct {
	*macaron.Context
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store

	User       *db.User
	IsLogged   bool
	Service    *db.Service
	ServiceURL string
}

// Success 返回模板，状态码 200
func (c *Context) Success(name string) {
	c.HTML(http.StatusOK, name)
}

// Error 返回模板错误页
func (c *Context) Error(err error) {
	c.Data["ErrorMsg"] = err
	c.HTML(http.StatusOK, "error")
}

// RenderWithErr 返回表单报错
func (c *Context) RenderWithErr(msg, tpl string, f interface{}) {
	if f != nil {
		form.Assign(f, c.Data)
	}
	c.Flash.ErrorMsg = msg
	c.Data["Flash"] = c.Flash
	c.HTML(http.StatusOK, tpl)
}

// HasError 返回表单验证是否有错误
func (c *Context) HasError() bool {
	hasErr, ok := c.Data["HasError"]
	if !ok {
		return false
	}
	c.Flash.ErrorMsg = c.Data["ErrorMsg"].(string)
	c.Data["Flash"] = c.Flash
	return hasErr.(bool)
}

// Contexter initializes a classic context for a request.
func Contexter() macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store, f *session.Flash, x csrf.CSRF) {
		c := &Context{
			Context: ctx,
			csrf:    x,
			Flash:   f,
			Session: sess,
		}

		// Get user from session
		c.User = authenticatedUser(c.Session)

		if c.User != nil {
			c.IsLogged = true
			c.Data["LoggedUser"] = c.User
		}

		c.Data["CSRFToken"] = x.GetToken()
		c.Data["CSRFTokenHTML"] = template.Safe(`<input type="hidden" name="_csrf" value="` + x.GetToken() + `">`)
		log.Trace("Session ID: %s", sess.ID())
		log.Trace("CSRF Token: %v", c.Data["CSRFToken"])

		ctx.Map(c)
	}
}
