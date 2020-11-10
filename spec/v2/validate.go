package v2

import (
	"github.com/NekoWheel/NekoCAS/db"
	"github.com/NekoWheel/NekoCAS/web/context"
)

func ValidateHandler(c *context.Context) {
	c.Header().Set("Content-Type", "text/xml")

	ticket := c.Query("ticket")
	service := c.Service
	if service == nil || ticket == "" {
		c.PlainText(200, NewCASFailureResponse("INVALID_REQUEST", "Both ticket and service parameters must be given"))
		return
	}

	ticketUser, ticketService, ok := db.ValidateTicket(ticket)
	if !ok {
		c.PlainText(200, NewCASFailureResponse("INVALID_TICKET", "Ticket not recognized"))
		return
	}
	if ticketService.ID != service.ID {
		c.PlainText(200, NewCASFailureResponse("INVALID_SERVICE", "Ticket was used for another service than it was generated for"))
		return
	}

	c.PlainText(200, NewCASSuccessResponse(ticketUser.Name, ""))
}
