package account

import (
	"github.com/NekoWheel/NekoCAS/internal/web/context"
)

func DashboardViewHandler(c *context.Context) {
	c.Success("dashboard")
}
