package manager

import (
	"strings"

	"github.com/NekoWheel/NekoCAS/db"
	"github.com/NekoWheel/NekoCAS/web/context"
	"github.com/NekoWheel/NekoCAS/web/form"
	"gorm.io/gorm"
)

func ServicesViewHandler(c *context.Context) {
	total := db.CountServices()
	limit := 10
	page := c.QueryInt("p")
	if page <= 0 {
		page = 1
	}

	totalPage := int(total/int64(limit)) + 1
	if page > totalPage {
		page = totalPage
	}

	c.Data["From"] = (page-1)*limit + 1
	c.Data["To"] = page * limit

	c.Data["NextPage"] = page + 1
	c.Data["PreviousPage"] = page - 1
	c.Data["LastPage"] = total/int64(limit) + 1

	c.Data["Total"] = total
	c.Data["Services"] = db.GetServices((page-1)*limit, limit)
	c.Success("manage/services")
}

func NewServiceViewHandler(c *context.Context) {
	c.Success("manage/new_service")
}

func NewServiceActionHandler(c *context.Context, f form.NewService) {
	// 表单报错
	if c.HasError() {
		c.Success("manage/new_service")
		return
	}

	f.Domain = strings.TrimPrefix(f.Domain, "http://")
	f.Domain = strings.TrimPrefix(f.Domain, "https://")
	f.Domain = strings.TrimRight(f.Domain, "/")

	s := &db.Service{
		Name:   f.Name,
		Avatar: f.Avatar,
		Domain: f.Domain,
		Ban:    false,
	}

	if err := db.CreateService(s); err != nil {
		c.RenderWithErr(err.Error(), "manage/new_service", &f)
		return
	}
	c.Flash.Success("添加服务成功")
	c.Redirect("/manage/services")
}

func EditServiceViewHandler(c *context.Context) {
	id := c.QueryInt("id")
	service := db.GetServiceByID(uint(id))
	if service == nil {
		c.Flash.Error("服务不存在")
		c.Redirect("/manage/services")
		return
	}

	c.Data["Service"] = service
	c.Success("manage/edit_service")
}

func EditServiceActionHandler(c *context.Context, f form.EditService) {
	id := c.QueryInt("id")
	service := db.GetServiceByID(uint(id))
	if service == nil {
		c.Flash.Error("服务不存在")
		c.Redirect("/manage/services")
		return
	}

	// 表单报错
	if c.HasError() {
		c.Success("manage/edit_service")
		return
	}

	f.Domain = strings.TrimPrefix(f.Domain, "http://")
	f.Domain = strings.TrimPrefix(f.Domain, "https://")
	f.Domain = strings.TrimRight(f.Domain, "/")

	s := &db.Service{
		Model:  gorm.Model{ID: uint(id)},
		Name:   f.Name,
		Avatar: f.Avatar,
		Domain: f.Domain,
		Ban:    false,
	}

	if err := db.UpdateService(s); err != nil {
		c.RenderWithErr(err.Error(), "manage/new_service", &f)
		return
	}
	c.Flash.Success("修改服务成功")
	c.Redirect("/manage/services")
}

func DeleteServiceViewHandler(c *context.Context) {
	id := c.QueryInt("id")
	service := db.GetServiceByID(uint(id))
	if service == nil {
		c.Redirect("/manage/services")
		return
	}

	c.Data["Service"] = service
	c.Success("manage/delete_service")
}

func DeleteServiceActionHandler(c *context.Context) {
	id := c.QueryInt("id")
	service := db.GetServiceByID(uint(id))
	if service == nil {
		c.Redirect("/manage/services")
		return
	}

	if err := db.DeleteService(uint(id)); err != nil {
		c.Flash.Error("删除失败")
	} else {
		c.Flash.Success("删除服务成功")
	}
	c.Redirect("/manage/services")
}
