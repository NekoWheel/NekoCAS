package manager

import (
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/web/context"
)

func UsersViewHandler(c *context.Context) {
	total := db.CountUsers()
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
	c.Data["Users"] = db.GetUsers((page-1)*limit, limit)
	c.Success("manage/users")
}
