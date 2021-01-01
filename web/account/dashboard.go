package account

import "github.com/NekoWheel/NekoCAS/web/context"

func DashboardViewHandler(c *context.Context) {
	c.Success("dashboard")
}
