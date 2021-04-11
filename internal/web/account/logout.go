package account

import (
	"github.com/NekoWheel/NekoCAS/internal/web/context"
)

func LogoutViewHandler(c *context.Context) {
	c.Data["Service"] = c.Service
	c.Success("logout")
}

func LogoutActionHandler(c *context.Context) {
	c.Data["Service"] = c.Service
	_ = c.Session.Destory(c.Context)
	c.Success("logout")
}
