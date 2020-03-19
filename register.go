package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"github.com/wuhan005/govalid"
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
		Mail     string `form:"mail" valid:"required;email;maxlen=30"`
		Name     string `form:"name" valid:"required;minlen=5;maxlen=20" label:"昵称"`
		Password string `form:"password" valid:"required;minlen=8;maxlen=30" label:"密码"`
	}{}

	// check form
	err := c.ShouldBind(&registerForm)
	if err != nil {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"error": "数据格式不正确",
			"name":  registerForm.Name,
			"mail":  registerForm.Mail,
		})
		return
	}
	// check form
	v := govalid.New(registerForm)
	if !v.Check() {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"error": v.Errors[0].Message,
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
