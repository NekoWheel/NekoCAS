package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func (cas *cas) indexViewHandler(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("userID") == nil {
		c.Redirect(302, "/login")
		return
	}

	u := new(user)
	cas.DB.Model(&user{}).Where(&user{Model: gorm.Model{ID: session.Get("userID").(uint)}}).Find(&u)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"name": u.Name,
		"email": u.Email,
		"avatar": "https://cdn.v2ex.com/gravatar/" + cas.md5(u.Email),
	})
}
