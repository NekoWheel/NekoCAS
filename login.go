package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wuhan005/govalid"
	"net/http"
)

func (cas *cas) loginPreCheck(c *gin.Context) {
	serviceQuery, _ := c.GetQuery("service")
	if serviceQuery != "" {
		serviceData, err := cas.getServiceByURL(serviceQuery)
		if err != nil {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"title":   "非法访问",
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("serviceID", int(serviceData.ID))
		c.Set("serviceURL", serviceQuery)
	} else {
		// login to cas
		c.Set("serviceID", 0)
		c.Set("serviceURL", "")
	}
	c.Next()
}

func (cas *cas) loginViewHandler(c *gin.Context) {
	serviceURL := c.GetString("serviceURL")
	serviceID := c.GetInt("serviceID")

	// is login
	session := sessions.Default(c)
	if session.Get("userID") != nil {
		if serviceID == 0 {
			// login cas
			c.Redirect(302, "/")
		} else {
			// login service
			userID := session.Get("userID").(uint)
			auth := cas.getServiceAuth(uint(serviceID), userID)
			if auth.ID == 0 {
				// service first time login, ask user for permission.
				serviceData := new(service)
				cas.DB.Model(&service{}).Where(&service{Model: gorm.Model{ID: uint(serviceID)}}).Find(&serviceData)
				c.HTML(http.StatusOK, "authorize.tmpl", gin.H{
					"_csrf":       c.GetString("_csrf"),
					"serviceName": serviceData.Name,
					"serviceURL":  serviceURL,
				})
				return
			}
			c.Redirect(302, cas.newServiceTicketCallBack(serviceURL, userID, serviceID))
		}
		return
	}

	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"_csrf": c.GetString("_csrf"),
		"error": "",
	})
}

func (cas *cas) loginActionHandler(c *gin.Context) {
	serviceURL := c.GetString("serviceURL")
	serviceID := c.GetInt("serviceID")

	loginForm := struct {
		Email    string `form:"mail" valid:"required"`
		Password string `form:"password" valid:"required"`
	}{}

	err := c.ShouldBind(&loginForm)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
		c.Abort()
		return
	}

	v := govalid.New(loginForm)
	if !v.Check() {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": "登录失败！电子邮箱或密码错误！",
			"_csrf": c.GetString("_csrf"),
		})
		c.Abort()
		return
	}

	u := new(user)
	cas.DB.Model(&user{}).Where(&user{Email: loginForm.Email}).Find(&u)
	if u.ID == 0 || u.Password != cas.addSalt(loginForm.Password) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": "登录失败！电子邮箱或密码错误！",
			"_csrf": c.GetString("_csrf"),
		})
		c.Abort()
		return
	}

	session := sessions.Default(c)
	session.Set("userID", u.ID)
	session.Set("userType", u.Permission)
	_ = session.Save()

	if serviceID == 0 {
		c.Redirect(302, "/")
	} else {
		userID := session.Get("userID").(uint)
		auth := cas.getServiceAuth(uint(serviceID), userID)
		if auth.ID == 0 {
			serviceData := new(service)
			cas.DB.Model(&service{}).Where(&service{Model: gorm.Model{ID: uint(serviceID)}}).Find(&serviceData)
			// service first time login, ask user for permission.
			c.HTML(http.StatusOK, "authorize.tmpl", gin.H{
				"_csrf":       c.GetString("_csrf"),
				"serviceName": serviceData.Name,
				"serviceURL":  serviceURL,
			})
			return
		}
		c.Redirect(302, cas.newServiceTicketCallBack(serviceURL, userID, serviceID))
	}
}
