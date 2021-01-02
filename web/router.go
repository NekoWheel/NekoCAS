package web

import (
	"net/http"

	"github.com/NekoWheel/NekoCAS/conf"
	v1 "github.com/NekoWheel/NekoCAS/spec/v1"
	v2 "github.com/NekoWheel/NekoCAS/spec/v2"
	"github.com/NekoWheel/NekoCAS/web/account"
	"github.com/NekoWheel/NekoCAS/web/context"
	"github.com/NekoWheel/NekoCAS/web/form"
	"github.com/NekoWheel/NekoCAS/web/middleware"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func Run() {
	r := macaron.Classic()

	// 登录登出状态
	reqSignIn := context.Toggle(&context.ToggleOptions{SignInRequired: true})
	reqSignOut := context.Toggle(&context.ToggleOptions{SignOutRequired: true})

	renderOpt := macaron.RenderOptions{
		Directory:  "templates",
		IndentJSON: macaron.Env != macaron.PROD,
	}
	r.Use(macaron.Renderer(renderOpt))
	r.Use(cache.Cacher())

	bindIgnErr := binding.BindIgnErr

	r.Group("", func() {
		// 登录前访问
		r.Group("", func() {
			r.Combo("/register").
				Get(account.RegisterViewHandler).
				Post(bindIgnErr(form.Register{}), account.RegisterActionHandler)
			r.Any("/activate_code", account.VerifyUserActiveCodeHandler)
		}, reqSignOut)

		// 无论是否已经登录都可以访问
		r.Combo("/login", middleware.ServicePreCheck).
			Get(account.LoginViewHandler).
			Post(bindIgnErr(form.Login{}), account.LoginActionHandler)

		// 登录后访问
		r.Group("", func() {
			r.Get("/", account.DashboardViewHandler)
			r.Combo("/profile").Get(account.ProfileViewHandler)
			r.Combo("/profile/edit").Get(account.ProfileEditViewHandler).Post(bindIgnErr(form.UpdateProfile{}), account.ProfileEditActionHandler)
			r.Post("/logout", account.LogoutHandler)
			r.Combo("/activate").Get(account.ActivationViewHandler).Post(account.ActivationActionHandler)
		}, reqSignIn)

		// CAS 协议实现
		r.Get("/validate", middleware.ServicePreCheck, v1.ValidateHandler)        // v1
		r.Get("/serviceValidate", middleware.ServicePreCheck, v2.ValidateHandler) // v2
	},
		session.Sessioner(session.Options{
			CookieName: "nekocas",
		}),

		csrf.Csrfer(csrf.Options{
			Secret: conf.Get().Site.CSRFKey,
			Header: "X-CSRF-Token",
		}),

		context.Contexter(),
	)

	r.NotFound(func(c *macaron.Context) {
		c.Data["Title"] = "页面不存在"
		c.HTML(http.StatusNotFound, "404")
	})

	r.Run("0.0.0.0", conf.Get().Site.Port)
}
