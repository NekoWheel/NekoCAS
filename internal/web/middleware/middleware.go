package middleware

import (
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/web/context"
)

// ServicePreCheck 获取 Service 信息中间件
func ServicePreCheck(c *context.Context) {
	serviceURL := c.Query("service")
	if serviceURL == "" {
		c.Service = &db.Service{}
		return
	}

	service, err := db.ServiceByURL(serviceURL)
	if err != nil {
		c.Error(err)
		return
	}
	c.ServiceURL = serviceURL
	c.Service = service
	c.Next()
}
