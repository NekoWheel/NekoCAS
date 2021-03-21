package mail

import (
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/NekoWheel/NekoCAS/internal/conf"
	"github.com/NekoWheel/NekoCAS/internal/filesystem"
	"github.com/NekoWheel/NekoCAS/templates"
	"gopkg.in/gomail.v2"
	"gopkg.in/macaron.v1"
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
			Directory:  "templates",
			IndentJSON: macaron.Env != macaron.PROD,
			Funcs: []template.FuncMap{map[string]interface{}{
				"Year": func() int {
					return time.Now().Year()
				},
			}},
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
		"SiteName": conf.Get().Site.Name,
		"Email":    to,
		"Link":     conf.Get().Site.BaseURL + "/activate_code?code=" + code,
	}
	body, err := render("activate", data)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("激活您的 %s 账号", conf.Get().Site.Name)
	return send(to, title, body)
}

func SendLostPasswordMail(to, code string) error {
	data := map[string]interface{}{
		"SiteName": conf.Get().Site.Name,
		"Email":    to,
		"Link":     conf.Get().Site.BaseURL + "/reset_password?code=" + code,
	}
	body, err := render("reset_password", data)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("您正在找回您的 %s 账号密码", conf.Get().Site.Name)
	return send(to, title, body)
}

func send(to, title, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", conf.Get().Mail.Account)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(
		conf.Get().Mail.SMTP,
		conf.Get().Mail.Port,
		conf.Get().Mail.Account,
		conf.Get().Mail.Password,
	)
	return d.DialAndSend(m)
}
