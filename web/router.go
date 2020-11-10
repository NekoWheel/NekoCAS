package web

import (
	"github.com/NekoWheel/NekoCAS/conf"
	v1 "github.com/NekoWheel/NekoCAS/spec/v1"
	v2 "github.com/NekoWheel/NekoCAS/spec/v2"
	"github.com/NekoWheel/NekoCAS/web/account"
	"github.com/NekoWheel/NekoCAS/web/context"
	"github.com/NekoWheel/NekoCAS/web/form"
	"github.com/NekoWheel/NekoCAS/web/middleware"
	"github.com/go-macaron/binding"
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

	bindIgnErr := binding.BindIgnErr

	r.Group("", func() {
		// 登录前访问
		r.Group("", func() {
			r.Combo("/register").
				Get(account.RegisterViewHandler).
				Post(bindIgnErr(form.Register{}), account.RegisterActionHandler)
		}, reqSignOut)

		// 无论是否已经登录都可以访问
		r.Combo("/login", middleware.ServicePreCheck).
			Get(account.LoginViewHandler).
			Post(bindIgnErr(form.Login{}), account.LoginActionHandler)

		// 登录后访问
		r.Group("", func() {
			r.Get("/")
			r.Get("/logout", account.LogoutHandler)
		}, reqSignIn)

		// CAS 协议实现
		r.Get("/validate", middleware.ServicePreCheck, v1.ValidateHandler)        // v1
		r.Get("/serviceValidate", middleware.ServicePreCheck, v2.ValidateHandler) // v2
		//r.Get("/proxy", )
		//r.Get("/proxyValidate", )
	},
		session.Sessioner(session.Options{
			CookieName: "nekocas",
		}),

		csrf.Csrfer(csrf.Options{
			Secret: conf.Get().CSRFKey,
			Header: "X-CSRF-Token",
		}),

		context.Contexter(),
	)

	//r.HTMLRender = middleware.LoadTemplates("./templates")
	//store := sessions.NewCookieStore([]byte(config.Get().SessionKey))
	//r.Use(sessions.Middleware("session", store))
	//// CSRF Token
	//r.Use(middleware.CSRF())
	//
	//// Account
	//r.GET("/register", middleware.NotLoggedIn, middleware.OpenForRegister, account.RegisterViewHandler)
	//r.POST("/register", middleware.NotLoggedIn, middleware.OpenForRegister, account.RegisterActionHandler)
	//r.POST("/login", middleware.NotLoggedIn, middleware.ServicePreCheck, account.LoginActionHandler)
	//r.GET("/", middleware.LoggedIn)
	////r.GET("/profile", middleware.LoggedIn, cas.profileViewHandler)
	////r.POST("/profile", middleware.LoggedIn, cas.profileActionHandler)
	//
	//// Service
	//
	//// CAS Protocol
	//r.GET("/login", middleware.ServicePreCheck, cas.LoginViewHandler) // Login view
	//r.POST("/logout", middleware.LoggedIn, cas.LogoutHandler)         // Logout action
	//r.POST("/validate", middleware.LoggedIn, cas.ValidateHandler)     // Server-side validate action
	//
	//// service first time login, ask user for permission.
	////r.POST("/authorize", cas.authRequired, cas.authorizeHandler)
	////r.POST("/revoke", cas.authRequired, cas.revoke)
	////
	////r.GET("/validate", cas.validateHandler)
	//
	//// manager
	////manage := r.Group("/manage", cas.authRequired, cas.adminRequired)
	////manage.GET("/", cas.managerViewHandler)
	////manage.POST("/service/ban", cas.switchBanServiceHandler)
	////manage.POST("/service/delete", cas.deleteServiceHandler)
	////manage.GET("/service/new", cas.newServiceViewHandler)
	////manage.POST("/service/new", cas.newServiceHandler)
	////
	////manage.POST("/user/admin", cas.setAdminHandler)
	////manage.POST("/user/delete", cas.deleteAdminHandler)
	//
	//err := r.Run(":" + strconv.Itoa(config.Get().Port))
	//log.Fatal("Failed to start web server: %v", err)

	r.Run()
}
