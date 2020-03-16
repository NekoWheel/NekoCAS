package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	serviceQuery, _ := c.GetQuery("service")
	if serviceQuery != "" {
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
		c.Set("serviceID", int(serviceData.ID))
		c.Set("serviceURL", serviceURL)
	} else {
		// login to cas
		c.Set("serviceID", 0)
		c.Set("serviceURL", &url.URL{})
	}
	c.Next()
}

func (cas *cas) loginViewHandler(c *gin.Context) {
	// TODO 500 error handler
	serviceURLInterface, _ := c.Get("serviceURL")
	serviceURL, _ := serviceURLInterface.(*url.URL)
	serviceID := c.GetInt("serviceID")

	// is login
	session := sessions.Default(c)
	if session.Get("userID") != nil {
		if serviceID == 0 {
			c.Redirect(302, "/")
		} else {
			userID := session.Get("userID").(uint)
			c.Redirect(302, cas.newServiceTicketCallBack(serviceURL, userID, serviceID))
		}
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
	serviceID := c.GetInt("serviceID")

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

	if serviceID == 0 {
		c.Redirect(302, "/")
	} else {
		userID := session.Get("userID").(uint)
		c.Redirect(302, cas.newServiceTicketCallBack(serviceURL, userID, serviceID))
	}
}

func (cas *cas) newServiceTicketCallBack(serviceURL *url.URL, userID uint, serviceID int) string {
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

func (cas *cas) validateHandler(c *gin.Context) {
	serviceQuery, _ := c.GetQuery("service")
	serviceTicketQuery, ok := c.GetQuery("ticket")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    nil,
			"message": "Missing `ticket`.",
		})
		return
	}

	// get ticket first
	stDataString := cas.Redis.Get(serviceTicketQuery).Val()
	cas.Redis.Del(serviceTicketQuery)
	if stDataString == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"data":    nil,
			"message": "Invalid ticket.",
		})
		return
	}

	stData := strings.Split(stDataString, "|")
	if len(stData) != 2 {
		c.JSON(http.StatusBadGateway, gin.H{
			"success": false,
			"data":    nil,
			"message": "Server error.",
		})
		return
	}

	// get service
	serviceData := new(service)
	cas.DB.Model(&service{}).Where("secret = ?", serviceQuery).Find(&serviceData)

	userIDStr := stData[0]
	serviceIDStr := stData[1]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"success": false,
			"data":    nil,
			"message": "Server error.",
		})
		return
	}
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"success": false,
			"data":    nil,
			"message": "Server error.",
		})
		return
	}

	if serviceData.ID != uint(serviceID) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"data":    nil,
			"message": "Invalid service.",
		})
		return
	}

	userData := new(user)
	cas.DB.Model(&user{}).Where(&user{Model: gorm.Model{ID: uint(userID)}}).Find(&userData)
	if userData.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"data":    nil,
			"message": "User not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"name":  userData.Name,
			"email": userData.Email,
			"token": userData.Token,
		},
		"message": "ok",
	})
}
