package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func (cas *cas) setAdminHandler(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(302, "/manage")
		return
	}

	var u user
	cas.DB.Model(&user{}).Where("id = ?", id).Find(&u)
	tx := cas.DB.Begin()
	if u.Permission == 0 {
		u.Permission = 1
	} else {
		u.Permission = 0
	}
	if tx.Model(&user{}).Where("id = ?", id).Update(map[string]interface{}{"permission": u.Permission}).RowsAffected != 1 {
		tx.Rollback()
		c.Redirect(302, "/manage")
		return
	}
	tx.Commit()
	c.Redirect(302, "/manage")
}

func (cas *cas) deleteAdminHandler(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(302, "/manage")
		return
	}

	tx := cas.DB.Begin()
	if tx.Where("id = ?", id).Delete(&user{}).RowsAffected != 1 {
		tx.Rollback()
		c.Redirect(302, "/manage")
		return
	}
	tx.Commit()
	c.Redirect(302, "/manage")
}
