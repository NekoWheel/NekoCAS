package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

func (cas *cas) registerPreCheck(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("userID") != nil {
		c.Redirect(302, "/")
	}
}

func (cas *cas) registerViewHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", gin.H{
		"error": "",
	})
}

func (cas *cas) registerActionHandler(c *gin.Context) {
	registerForm := struct {
		Mail     string `form:"mail" binding:"required,email,max=30"`
		Name     string `form:"name" binding:"required,min=5,max=20"`
		Password string `form:"password" binding:"required,min=8,max=30"`
	}{}

	// check form
	errs := c.ShouldBind(&registerForm)
	if errs != nil {
		err := errs.(validator.ValidationErrors)[0]
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"error": cas.getErrorMessage(err.Field(), err.Tag(), err.Value()),
			"name":  registerForm.Name,
			"mail":  registerForm.Mail,
		})
		return
	}

	u := new(user)
	cas.DB.Model(&user{}).Where(&user{Email: registerForm.Mail}).Find(&u)
	if u.ID != 0 {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"error": "该电子邮箱已经注册过了！",
			"email": registerForm.Mail,
			"name":  registerForm.Name,
		})
		c.Abort()
		return
	}

	cas.DB.Create(&user{
		Name:     registerForm.Name,
		Email:    registerForm.Mail,
		Password: cas.addSalt(registerForm.Password),
		Token:    randstr.String(32),
	})

	c.Redirect(302, "/login")
}
