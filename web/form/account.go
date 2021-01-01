package form

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type Register struct {
	Mail     string `binding:"Required;Email;MaxSize(254)" locale:"电子邮箱"`
	Name     string `binding:"Required;AlphaDashDot;MaxSize(35)" locale:"用户名"`
	NickName string `binding:"Required;MaxSize(20)" locale:"昵称"`
	Password string `binding:"Required;MaxSize(255)" locale:"密码"`
	Retype   string
}

func (f *Register) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}

type Login struct {
	Mail     string `binding:"Required;Email;MaxSize(254)" locale:"电子邮箱"`
	Password string `binding:"Required;MaxSize(255)" locale:"密码"`
}

func (f *Login) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}

type UpdateProfile struct {
	NickName string `binding:"Required;MaxSize(20)" locale:"昵称"`
	Password string `binding:"MaxSize(255)" locale:"密码"`
	Retype   string
}

func (f *UpdateProfile) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}
