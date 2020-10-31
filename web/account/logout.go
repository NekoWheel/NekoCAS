package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tommy351/gin-sessions"
)

func LogoutHandler(c *gin.Context) {
	// TODO: url
	session := sessions.Get(c)
	session.Clear()
	_ = session.Save()
	// TODO: logout.tpl
	c.HTML(http.StatusOK, "logout.tpl", nil)
}
