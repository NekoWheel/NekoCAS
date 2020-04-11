package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (cas *cas) initRouter() {
	r := gin.Default()
	r.HTMLRender = cas.loadTemplates("./templates")
	r.Use(sessions.Sessions("session", memstore.NewStore([]byte(cas.Conf.Key))))
	r.Use(cas.csrfMiddleware)

	// register, middleware prevent register after login.
	r.GET("/register", cas.registerPreCheck, cas.registerViewHandler)
	r.POST("/register", cas.registerPreCheck, cas.registerActionHandler)

	// login, middleware check the service data if exists.
	r.GET("/login", cas.loginPreCheck, cas.loginViewHandler)
	r.POST("/login", cas.loginPreCheck, cas.loginActionHandler)
	r.POST("/logout", cas.logoutHandler)

	r.GET("/", cas.authRequired, cas.indexViewHandler)
	r.GET("/profile", cas.authRequired, cas.profileViewHandler)
	r.POST("/profile", cas.authRequired, cas.profileActionHandler)

	// service first time login, ask user for permission.
	r.POST("/authorize", cas.authRequired, cas.authorizeHandler)
	r.POST("/revoke", cas.authRequired, cas.revoke)

	r.GET("/validate", cas.validateHandler)

	// manager
	manage := r.Group("/manage", cas.authRequired, cas.adminRequired)
	manage.GET("/", cas.managerViewHandler)
	manage.POST("/service/ban", cas.switchBanServiceHandler)
	manage.POST("/service/delete", cas.deleteServiceHandler)
	manage.GET("/service/new", cas.newServiceViewHandler)
	manage.POST("/service/new", cas.newServiceHandler)

	panic(r.Run(":" + strconv.Itoa(cas.Conf.Port)))
}
