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

	login := r.Group("/login")
	{
		login.Use(cas.loginPreCheck)
		login.GET("/", cas.loginViewHandler)
		login.POST("/", cas.loginActionHandler)
	}

	panic(r.Run(":" + strconv.Itoa(cas.Conf.Port)))
}
