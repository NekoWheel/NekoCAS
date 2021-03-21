package form

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type Site struct {
	OpenRegister  bool   `locale:"是否开放注册"`
	SiteLogo      string `binding:"MaxSize(255)" locale:"站点图标"`
	MailWhitelist string `binding:"MaxSize(255)" locale:"注册邮箱白名单"`
	Privacy       string `locale:"隐私政策"`
}

func (f *Site) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}
