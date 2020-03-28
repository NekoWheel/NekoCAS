package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/thanhpk/randstr"
	"net/http"
	"strconv"
	"strings"
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

	// every app get the different token.
	auth := cas.getServiceAuth(serviceData.ID, userData.ID)
	if auth.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"data":    nil,
			"message": "No permission.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"name":  userData.Name,
			"email": userData.Email,
			"token": auth.Token,
		},
		"message": "ok",
	})
}

func (cas *cas) authorizeHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	serviceURL, ok := c.GetPostForm("serviceURL")
	if !ok {
		c.Redirect(302, "/")
		return
	}
	serviceData, err := cas.getServiceByURL(serviceURL)
	if err != nil {
		c.Redirect(302, "/")
		return
	}

	auth := new(serviceAuth)
	cas.DB.Model(&serviceAuth{}).Where(&serviceAuth{
		ServiceID: serviceData.ID,
		UserID:    userID,
	}).Find(&auth)
	if auth.ID == 0 {
		tx := cas.DB.Begin()
		if tx.Create(&serviceAuth{
			ServiceID: serviceData.ID,
			UserID:    userID,
			Token:     randstr.String(32),
		}).RowsAffected != 1 {
			tx.Rollback()
			c.Redirect(302, "/")
			return
		}
		tx.Commit()
	}
	c.Redirect(302, cas.newServiceTicketCallBack(serviceURL, userID, int(serviceData.ID)))
	return
}

func (cas *cas) revoke(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	serviceIDStr := c.PostForm("id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		c.Redirect(302, "/")
		return
	}
	tx := cas.DB.Begin()
	if tx.Where("service_id = ? AND user_id = ?", serviceID, userID).Delete(&serviceAuth{}).RowsAffected != 1 {
		tx.Rollback()
		c.Redirect(302, "/")
		return
	}
	tx.Commit()
	c.Redirect(302, "/")
}

func (cas *cas) logoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(302, "/login")
}
