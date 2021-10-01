package account

import (
	"github.com/NekoWheel/NekoCAS/internal/context"
)

func DashboardViewHandler(c *context.Context) {
	c.Success("dashboard")
}
