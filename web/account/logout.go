package account

import (
	"github.com/NekoWheel/NekoCAS/web/context"
)

func LogoutHandler(c *context.Context) {
	// TODO: url
	_ = c.Session.Destory(c.Context)
	c.Success("logout")
}
