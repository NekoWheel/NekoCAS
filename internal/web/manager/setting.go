package manager

import (
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/web/context"
	"github.com/NekoWheel/NekoCAS/internal/web/form"
)

func SiteViewHandler(c *context.Context) {
	c.Success("manage/site")
}

func SiteActionHandler(c *context.Context, f form.Site) {
	// 表单报错
	if c.HasError() {
		c.Success("manage/site")
		return
	}

	if f.OpenRegister {
		db.SetSetting("open_setting", "on")
	} else {
		db.SetSetting("open_setting", "off")
	}

	db.SetSetting("site_logo", f.SiteLogo)
	db.SetSetting("mail_whitelist", f.MailWhitelist)
	db.SetSetting("privacy", f.Privacy)

	c.Redirect("/manage/site")
}
