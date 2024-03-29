package mail

import (
	"crypto/tls"
	"fmt"
	"sync"

	"gopkg.in/gomail.v2"
	"gopkg.in/macaron.v1"

	"github.com/NekoWheel/NekoCAS/internal/conf"
	"github.com/NekoWheel/NekoCAS/internal/filesystem"
	templ "github.com/NekoWheel/NekoCAS/internal/web/template"
	"github.com/NekoWheel/NekoCAS/templates"
)

var (
	tplRender     *macaron.TplRender
	tplRenderOnce sync.Once
)

// render 根据给定的信息渲染邮件模板
func render(tpl string, data map[string]interface{}) (string, error) {
	tplRenderOnce.Do(func() {
		var templateFS macaron.TemplateFileSystem
		if macaron.Env == macaron.PROD {
			templateFS = filesystem.NewFS(templates.FS)
		}

		opt := &macaron.RenderOptions{
			IndentJSON:         macaron.Env != macaron.PROD,
			Funcs:              templ.FuncMap(),
			TemplateFileSystem: templateFS,
		}

		ts := macaron.NewTemplateSet()
		ts.Set(macaron.DEFAULT_TPL_SET_NAME, opt)
		tplRender = &macaron.TplRender{
			TemplateSet: ts,
			Opt:         opt,
		}
	})

	return tplRender.HTMLString(tpl, data)
}

func SendActivationMail(to, code string) error {
	data := map[string]interface{}{
		"SiteName": conf.Site.Name,
		"Email":    to,
		"Link":     conf.Site.BaseURL + "/activate_code?code=" + code,
	}
	body, err := render("mail/activate", data)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("激活您的 %s 账号", conf.Site.Name)
	return send(to, title, body)
}

func SendLostPasswordMail(to, code string) error {
	data := map[string]interface{}{
		"SiteName": conf.Site.Name,
		"Email":    to,
		"Link":     conf.Site.BaseURL + "/reset_password?code=" + code,
	}
	body, err := render("mail/reset_password", data)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("您正在找回您的 %s 账号密码", conf.Site.Name)
	return send(to, title, body)
}

func send(to, title, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", conf.Mail.Account)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(
		conf.Mail.SMTP,
		conf.Mail.Port,
		conf.Mail.Account,
		conf.Mail.Password,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(m)
}
