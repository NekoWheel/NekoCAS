package manager

import (
	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoCAS/internal/context"
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/form"
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
		err := db.SetSetting("open_register", "on")
		if err != nil {
			log.Error("Failed to set %q to %q", "open_register", "on")
		}
	} else {
		err := db.SetSetting("open_register", "off")
		if err != nil {
			log.Error("Failed to set %q to %q", "open_register", "off")
		}
	}

	_ = db.SetSetting("site_logo", f.SiteLogo)
	_ = db.SetSetting("mail_whitelist", f.MailWhitelist)
	_ = db.SetSetting("privacy", f.Privacy)

	c.Redirect("/manage/site")
}
