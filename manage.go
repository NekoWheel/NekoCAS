package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"github.com/wuhan005/govalid"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	var services []service
	cas.DB.Model(&service{}).Find(&services)
	var users []user
	cas.DB.Model(&user{}).Find(&users)
	c.HTML(http.StatusOK, "manageIndex.tmpl", gin.H{
		"_csrf":    c.GetString("_csrf"),
		"services": services,
		"users":    users,
		"register": cas.Conf.Register,
	})
}

func (cas *cas) switchBanServiceHandler(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(302, "/manage")
		return
	}
	var s service
	cas.DB.Model(&service{}).Where("id = ?", id).Find(&s)
	tx := cas.DB.Begin()
	if tx.Model(&service{}).Where("id = ?", id).Update(map[string]interface{}{"ban": !s.Ban}).RowsAffected != 1 {
		tx.Rollback()
		c.Redirect(302, "/manage")
		return
	}
	tx.Commit()
	c.Redirect(302, "/manage")
}

func (cas *cas) deleteServiceHandler(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(302, "/manage")
		return
	}
	tx := cas.DB.Begin()
	if tx.Where("id = ?", id).Delete(&service{}).RowsAffected != 1 {
		tx.Rollback()
		c.Redirect(302, "/manage")
		return
	}
	tx.Commit()
	c.Redirect(302, "/manage")
}

func (cas *cas) newServiceViewHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "manageServiceNew.tmpl", gin.H{
		"_csrf": c.GetString("_csrf"),
		"error": "",
	})
}

func (cas *cas) newServiceHandler(c *gin.Context) {
	inputForm := struct {
		Name   string `form:"name" valid:"required;" label:"服务名"`
		Avatar string `form:"avatar" label:"服务 Logo"`
		Domain string `form:"domain" valid:"required;" label:"白名单域名"`
	}{}

	errs := c.ShouldBind(&inputForm)
	if errs != nil {
		c.HTML(http.StatusOK, "manageServiceNew.tmpl", gin.H{
			"error":  "数据格式不正确",
			"_csrf":  c.GetString("_csrf"),
			"name":   inputForm.Name,
			"avatar": inputForm.Avatar,
			"domain": inputForm.Domain,
		})
		return
	}
	// check form
	v := govalid.New(inputForm)
	if !v.Check() {
		c.HTML(http.StatusOK, "manageServiceNew.tmpl", gin.H{
			"error":  v.Errors[0].Message,
			"_csrf":  c.GetString("_csrf"),
			"name":   inputForm.Name,
			"avatar": inputForm.Avatar,
			"domain": inputForm.Domain,
		})
		return
	}

	newService := &service{
		Name:   inputForm.Name,
		Secret: randstr.Hex(16),
		Avatar: inputForm.Avatar,
		Ban:    false,
	}
	tx := cas.DB.Begin()
	if tx.Create(&newService).RowsAffected != 1 {
		tx.Rollback()
		c.HTML(http.StatusOK, "manageServiceNew.tmpl", gin.H{
			"error":  "服务器错误，添加失败！",
			"_csrf":  c.GetString("_csrf"),
			"name":   inputForm.Name,
			"avatar": inputForm.Avatar,
			"domain": inputForm.Domain,
		})
		return
	}
	tx.Commit()
	// add domain
	domains := strings.Split(inputForm.Domain, "\r\n")
	for _, v := range domains {
		d := v
		u, err := url.Parse(d)
		if err == nil {
			cas.DB.Create(&domain{
				Domain:    u.Host,
				ServiceID: newService.ID,
			})
		}
	}

	c.Redirect(302, "/manage/")
}
