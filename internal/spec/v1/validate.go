package v1

import (
	"fmt"
	"net/http"

	"github.com/NekoWheel/NekoCAS/internal/context"
	"github.com/NekoWheel/NekoCAS/internal/db"
)

func ValidateHandler(c *context.Context) {
	ticket := c.Query("ticket")
	user, ok := db.ValidateServiceTicket(c.Service, ticket)
	if ok {
		c.PlainText(http.StatusOK, []byte(fmt.Sprintf("yes\n%s\n", user.NickName)))
	}
	c.PlainText(http.StatusOK, []byte("no\n"))
}
