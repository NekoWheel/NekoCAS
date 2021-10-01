package v2

import (
	"net/http"

	"github.com/NekoWheel/NekoCAS/internal/context"
	"github.com/NekoWheel/NekoCAS/internal/db"
)

func ValidateHandler(c *context.Context) {
	c.Header().Set("Content-Type", "text/xml")

	ticket := c.Query("ticket")
	service := c.Service
	if service == nil || ticket == "" {
		c.PlainText(http.StatusOK, NewCASFailureResponse("INVALID_REQUEST", "Both ticket and service parameters must be given"))
		return
	}

	ticketUser, ticketService, ok := db.ValidateTicket(ticket)
	if !ok {
		c.PlainText(http.StatusOK, NewCASFailureResponse("INVALID_TICKET", "Ticket not recognized"))
		return
	}
	if ticketService.ID != service.ID {
		c.PlainText(http.StatusOK, NewCASFailureResponse("INVALID_SERVICE", "Ticket was used for another service than it was generated for"))
		return
	}

	c.PlainText(http.StatusOK, NewCASSuccessResponse(ticketUser))
}
