package account

import (
	"github.com/NekoWheel/NekoCAS/internal/context"
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/form"
)

func ProfileViewHandler(c *context.Context) {
	c.Success("profile")
}

func ProfileEditViewHandler(c *context.Context) {
	c.Success("profile_edit")
}

func ProfileEditActionHandler(c *context.Context, f form.UpdateProfile) {
	// 表单报错
	if c.HasError() {
		c.Success("profile_edit")
		return
	}

	// 检查昵称是否使用
	u, err := db.GetUserByNickName(f.NickName)
	if err == nil && u.ID != c.User.ID {
		c.Flash.Error("用户昵称已被使用，换一个吧~")
		c.Redirect("/profile/edit", 302)
		return
	}

	if f.Password != f.Retype {
		c.Flash.Error("两次输入的密码不匹配")
		c.Redirect("/profile/edit", 302)
		return
	}

	c.User.NickName = f.NickName
	c.User.Password = f.Password

	if err := db.UpdateUserProfile(c.User); err != nil {
		c.Flash.Error("修改个人信息失败")
		c.Redirect("/profile/edit", 302)
		return
	}

	c.Flash.Success("修改个人信息成功")
	c.Redirect("/profile/edit", 302)
}
