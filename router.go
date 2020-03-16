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

	r.GET("/", cas.indexViewHandler)

	login := r.Group("/login")
	{
		login.Use(cas.loginPreCheck)
		login.GET("/", cas.loginViewHandler)
		login.POST("/", cas.loginActionHandler)
	}

	register := r.Group("/register")
	{
		register.Use(cas.registerPreCheck)
		register.GET("/", cas.registerViewHandler)
		register.POST("/", cas.registerActionHandler)
	}

	r.GET("/validate", cas.validateHandler)

	panic(r.Run(":" + strconv.Itoa(cas.Conf.Port)))
}
