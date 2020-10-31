package profile

//
//func (cas *cas) authRequired(c *gin.Context) {
//	session := sessions.Default(c)
//	if session.Get("userID") == nil {
//		c.Redirect(302, "/login")
//		c.Abort()
//		return
//	}
//	c.Set("userID", session.Get("userID"))
//	c.Next()
//}
//
//func (cas *cas) indexViewHandler(c *gin.Context) {
//	u := new(user)
//	cas.DB.Model(&user{}).Where(&user{Model: gorm.Model{ID: c.MustGet("userID").(uint)}}).Find(&u)
//
//	// get service
//	auths := make([]serviceAuth, 0)
//	cas.DB.Model(&serviceAuth{}).Where(&serviceAuth{UserID: u.ID}).Find(&auths)
//	authIDs := make([]uint, len(auths))
//	for index, auth := range auths {
//		authIDs[index] = auth.ServiceID
//	}
//	services := make([]service, 0)
//	cas.DB.Model(&service{}).Where("id in (?)", authIDs).Find(&services)
//
//	c.HTML(http.StatusOK, "index.tmpl", gin.H{
//		"error":      "",
//		"_csrf":      c.GetString("_csrf"),
//		"name":       u.Name,
//		"email":      u.Email,
//		"permission": u.Permission,
//		"services":   services,
//		"auths":      auths,
//		"avatar":     "https://cdn.v2ex.com/gravatar/" + cas.md5(u.Email),
//	})
//}
//
//func (cas *cas) profileViewHandler(c *gin.Context) {
//	u := new(user)
//	cas.DB.Model(&user{}).Where(&user{Model: gorm.Model{ID: c.MustGet("userID").(uint)}}).Find(&u)
//	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
//		"error":    "",
//		"_csrf":    c.GetString("_csrf"),
//		"name":     u.Name,
//		"email":    u.Email,
//		"nameForm": u.Name,
//		"avatar":   "https://cdn.v2ex.com/gravatar/" + cas.md5(u.Email),
//	})
//}
//
//func (cas *cas) profileActionHandler(c *gin.Context) {
//	u := new(user)
//	cas.DB.Model(&user{}).Where(&user{Model: gorm.Model{ID: c.MustGet("userID").(uint)}}).Find(&u)
//
//	updateForm := struct {
//		Name     string `form:"name" valid:"required;minlen:5;maxlen:20" label:"昵称"`
//		Password string `form:"password" valid:"minlen:8;maxlen:30" label:"密码"`
//	}{}
//
//	errs := c.ShouldBind(&updateForm)
//	if errs != nil {
//		c.HTML(http.StatusOK, "profile.tmpl", gin.H{
//			"error":    "数据格式不正确",
//			"_csrf":    c.GetString("_csrf"),
//			"name":     u.Name,
//			"email":    u.Email,
//			"nameForm": updateForm.Name,
//			"avatar":   "https://cdn.v2ex.com/gravatar/" + cas.md5(u.Email),
//		})
//		return
//	}
//	// check form
//	v := govalid.New(updateForm)
//	if !v.Check() {
//		c.HTML(http.StatusOK, "profile.tmpl", gin.H{
//			"error":    v.Errors[0].Message,
//			"_csrf":    c.GetString("_csrf"),
//			"name":     u.Name,
//			"email":    u.Email,
//			"nameForm": updateForm.Name,
//			"avatar":   "https://cdn.v2ex.com/gravatar/" + cas.md5(u.Email),
//		})
//		return
//	}
//
//	updateUser := user{
//		Name: updateForm.Name,
//	}
//	if updateForm.Password != "" {
//		updateUser.Password = cas.addSalt(updateForm.Password)
//	}
//
//	tx := cas.DB.Begin()
//	if tx.Model(&user{}).Where(&user{Model: gorm.Model{ID: u.ID}}).Update(&updateUser).RowsAffected != 1 {
//		tx.Rollback()
//		c.HTML(http.StatusOK, "profile.tmpl", gin.H{
//			"error":    "服务器错误，修改失败！",
//			"_csrf":    c.GetString("_csrf"),
//			"name":     u.Name,
//			"email":    u.Email,
//			"nameForm": updateForm.Name,
//			"avatar":   "https://cdn.v2ex.com/gravatar/" + cas.md5(u.Email),
//		})
//		return
//	}
//	tx.Commit()
//	c.Redirect(302, "/profile")
//	return
//}
