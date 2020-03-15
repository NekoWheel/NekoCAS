package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/url"
	"time"
)

type cas struct {
	Conf   *config
	DB     *gorm.DB
	Router *gin.Engine
	Redis  *redis.Client
}

func (cas *cas) init() {
	cas.initConfig()
	cas.initRedis()
	cas.initDatabase()
	cas.initRouter()
}

func (cas *cas) loginPreCheck(c *gin.Context) {
	serviceQuery, ok := c.GetQuery("service")
	if !ok {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"title":   "非法访问",
			"message": "缺少请求参数 service",
		})
		c.Abort()
		return
	}

	serviceURL, err := url.ParseRequestURI(serviceQuery)
	if err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"title":   "非法访问",
			"message": "参数 service 无效",
		})
		c.Abort()
		return
	}
	if serviceURL.Scheme != "https" {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"title":   "非法访问",
			"message": "service 非 HTTPS 协议",
		})
		c.Abort()
		return
	}

	// check service whitelist
	trustDomain := new(domain)
	cas.DB.Model(&domain{}).Where(&domain{Domain: serviceURL.Hostname()}).Find(&trustDomain)
	if trustDomain.ID == 0 {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"title":   "非法访问",
			"message": "域名不在白名单内",
		})
		c.Abort()
		return
	}

	// get service id
	serviceData := new(service)
	cas.DB.Model(&service{}).Where(&service{Model: gorm.Model{ID: trustDomain.ServiceID}}).Find(&serviceData)
	c.Set("serviceID", serviceData.ID)
	c.Set("serviceURL", serviceURL)
	c.Next()
}

func (cas *cas) loginViewHandler(c *gin.Context) {
	// TODO 500 error handler
	serviceURLInterface, _ := c.Get("serviceURL")
	serviceURL, _ := serviceURLInterface.(*url.URL)
	serviceIDInterface, _ := c.Get("serviceID")
	serviceID, _ := serviceIDInterface.(uint)

	session := sessions.Default(c)
	if session.Get("userID") != nil {
		userID := session.Get("userID").(uint)
		c.Redirect(302, cas.newServiceTicketCallBack(serviceURL, userID, serviceID))
		return
	}

	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"error": "",
	})
}

func (cas *cas) loginActionHandler(c *gin.Context) {
	// TODO 500 error handler
	serviceURLInterface, _ := c.Get("serviceURL")
	serviceURL, _ := serviceURLInterface.(*url.URL)
	serviceIDInterface, _ := c.Get("serviceID")
	serviceID, _ := serviceIDInterface.(uint)

	loginForm := struct {
		Email    string `form:"mail" binding:"required"`
		Password string `form:"password" binding:"required"`
	}{}

	err := c.ShouldBind(&loginForm)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
		c.Abort()
		return
	}

	u := new(user)
	cas.DB.Model(&user{}).Where(&user{Email: loginForm.Email}).Find(&u)
	if u.ID == 0 || u.Password != cas.addSalt(loginForm.Password) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": "登录失败！电子邮箱或密码错误！",
		})
		c.Abort()
		return
	}

	session := sessions.Default(c)
	session.Set("userID", u.ID)
	_ = session.Save()

	c.Redirect(302, cas.newServiceTicketCallBack(serviceURL, u.ID, serviceID))
}

func (cas *cas) newServiceTicketCallBack(serviceURL *url.URL, userID uint, serviceID uint) string {
	// generate service ticket
	st := cas.generateServiceToken()
	// save the service ticket
	cas.Redis.Set(st, fmt.Sprintf("%d|%d", userID, serviceID), 5*time.Minute)

	// add query data into url
	query := serviceURL.Query()
	query.Set("ticket", st)
	serviceURL.RawQuery = query.Encode()
	return serviceURL.String()
}
