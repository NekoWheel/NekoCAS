package web

import (
	"net/http"

	"github.com/NekoWheel/NekoCAS/internal/conf"
	"github.com/NekoWheel/NekoCAS/internal/filesystem"
	"github.com/NekoWheel/NekoCAS/internal/spec/v1"
	"github.com/NekoWheel/NekoCAS/internal/spec/v2"
	"github.com/NekoWheel/NekoCAS/internal/web/account"
	"github.com/NekoWheel/NekoCAS/internal/web/context"
	"github.com/NekoWheel/NekoCAS/internal/web/form"
	"github.com/NekoWheel/NekoCAS/internal/web/manager"
	"github.com/NekoWheel/NekoCAS/internal/web/middleware"
	"github.com/NekoWheel/NekoCAS/internal/web/template"
	"github.com/NekoWheel/NekoCAS/public"
	"github.com/NekoWheel/NekoCAS/templates"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

// newMacaron 初始化一个新的 Macaron 实例。
func newMacaron() *macaron.Macaron {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Statics(macaron.StaticOptions{
		FileSystem: http.FS(public.FS),
	}, "."))

	return m
}

func Run() {
	r := newMacaron()

	var templateFS macaron.TemplateFileSystem
	if macaron.Env == macaron.PROD {
		templateFS = filesystem.NewFS(templates.FS)
	}

	// 登录登出状态
	reqSignIn := context.Toggle(&context.ToggleOptions{SignInRequired: true})
	reqSignOut := context.Toggle(&context.ToggleOptions{SignOutRequired: true})
	reqManager := context.Toggle(&context.ToggleOptions{SignInRequired: true, AdminRequired: true})

	renderOpt := macaron.RenderOptions{
		Directory:          "templates",
		IndentJSON:         macaron.Env != macaron.PROD,
		Funcs:              template.FuncMap(),
		TemplateFileSystem: templateFS,
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
			r.Combo("/lost_password").Get(account.LostPasswordHandler).Post(bindIgnErr(form.LostPassword{}), account.LostPasswordActionHandler)
			r.Combo("/reset_password").Get(account.ResetPasswordHandler).Post(bindIgnErr(form.ResetPassword{}), account.ResetPasswordActionHandler)
		}, reqSignOut)

		// 无论是否已经登录都可以访问
		r.Combo("/login", middleware.ServicePreCheck).
			Get(account.LoginViewHandler).
			Post(bindIgnErr(form.Login{}), account.LoginActionHandler)
		r.Any("/activate_code", account.VerifyUserActiveCodeHandler)
		r.Get("/privacy", func(c *context.Context) {
			c.Success("privacy")
		})

		// 登录后访问
		r.Group("", func() {
			r.Get("/", account.DashboardViewHandler)
			r.Combo("/profile").Get(account.ProfileViewHandler)
			r.Combo("/profile/edit").Get(account.ProfileEditViewHandler).Post(bindIgnErr(form.UpdateProfile{}), account.ProfileEditActionHandler)
			r.Combo("/logout", middleware.ServicePreCheck).Get(account.LogoutViewHandler).Post(account.LogoutActionHandler)
			r.Combo("/activate").Get(account.ActivationViewHandler).Post(account.ActivationActionHandler)
		}, reqSignIn)

		// 管理页面
		r.Group("/manage", func() {
			// 用户
			r.Get("/users", manager.UsersViewHandler)

			// 服务
			r.Get("/services", manager.ServicesViewHandler)
			r.Combo("/services/new").Get(manager.NewServiceViewHandler).Post(bindIgnErr(form.NewService{}), manager.NewServiceActionHandler)
			r.Combo("/services/edit").Get(manager.EditServiceViewHandler).Post(bindIgnErr(form.EditService{}), manager.EditServiceActionHandler)
			r.Combo("/services/delete").Get(manager.DeleteServiceViewHandler).Post(manager.DeleteServiceActionHandler)

			// 站点设置
			r.Combo("/site").Get(manager.SiteViewHandler).Post(bindIgnErr(form.Site{}), manager.SiteActionHandler)
		}, reqManager)

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
