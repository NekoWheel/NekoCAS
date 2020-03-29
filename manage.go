package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (cas *cas) adminRequired(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("userType") == nil {
		c.Redirect(302, "/login")
		c.Abort()
		return
	}
	userType, ok := session.Get("userType").(int)
	if !ok {
		c.Redirect(302, "/login")
		c.Abort()
		return
	}
	if userType == 1 {
		c.Next()
		return
	}
	c.Redirect(302, "/")
	c.Abort()
}

func (cas *cas) managerViewHandler(c *gin.Context) {

}
