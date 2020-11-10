package v1

import (
	"fmt"

	"github.com/NekoWheel/NekoCAS/db"
	"github.com/NekoWheel/NekoCAS/web/context"
)

func ValidateHandler(c *context.Context) {
	ticket := c.Query("ticket")
	user, ok := db.ValidateServiceTicket(c.Service, ticket)
	if ok {
		c.PlainText(200, []byte(fmt.Sprintf("yes\n%s\n", user.Name)))
	}
	c.PlainText(200, []byte("no\n"))
}
