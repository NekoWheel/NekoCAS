package form

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type NewService struct {
	Name   string `binding:"Required;MaxSize(254)" locale:"服务名"`
	Avatar string `binding:"Required;MaxSize(255)" locale:"服务 Logo 链接"`
	Domain string `binding:"Required;MaxSize(255)" locale:"白名单域名"`
}

func (f *NewService) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}

type EditService struct {
	Name   string `binding:"Required;MaxSize(254)" locale:"服务名"`
	Avatar string `binding:"Required;MaxSize(255)" locale:"服务 Logo 链接"`
	Domain string `binding:"Required;MaxSize(255)" locale:"白名单域名"`
}

func (f *EditService) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}
