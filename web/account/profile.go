package account

import (
	"fmt"

	"github.com/NekoWheel/NekoCAS/db"
	"github.com/NekoWheel/NekoCAS/web/context"
	"github.com/NekoWheel/NekoCAS/web/form"
)

func ProfileViewHandler(c *context.Context) {
	fmt.Println(c.Flash.Get("Tab"))
	c.Success("dashboard")
}

func ProfileActionHandler(c *context.Context, f form.UpdateProfile) {
	c.Flash.Set("Tab", "update_info")

	// 表单报错
	if c.HasError() {
		c.Success("dashboard")
		return
	}

	if f.Password != f.Retype {
		c.Flash.Error("两次输入的密码不匹配")
		c.Redirect("/", 302)
		return
	}

	c.User.NickName = f.NickName
	c.User.Password = f.Password

	if err := db.UpdateUserProfile(c.User); err != nil {
		c.Flash.Error("修改个人信息失败")
		c.Redirect("/", 302)
		return
	}

	c.Flash.Success("修改个人信息成功")
	c.Redirect("/", 302)
}
