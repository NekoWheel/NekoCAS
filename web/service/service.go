package service

import "github.com/NekoWheel/NekoCAS/db"

// Local 来源从 CAS 登录，跳转到 CAS 个人信息
func Local() *db.Service {
	return &db.Service{}
}
