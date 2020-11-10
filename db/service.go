package db

import (
	"net/url"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Service 接入的服务
type Service struct {
	gorm.Model

	Name   string
	Avatar string // 服务 Logo
	Domain string // 白名单域名
	Ban    bool   // 是否封禁
}

// ServiceByURL 通过 ServiceURL 查找对应的服务
func ServiceByURL(u string) (*Service, error) {
	serviceURL, err := url.ParseRequestURI(u)
	if err != nil || serviceURL.Hostname() == "" {
		return nil, errors.New("参数无效")
	}

	if serviceURL.Scheme != "https" && serviceURL.Scheme != "http" {
		return nil, errors.New("不支持的协议")
	}

	var service Service
	if err := db.Model(&Service{}).Where("domain = ? and ban = ?", serviceURL.Hostname(), false).First(&service).Error; err != nil {
		return nil, errors.New("域名不在白名单内")
	}
	return &service, nil
}

// GetServiceByID 通过 Service ID 查找对应的服务
func GetServiceByID(id uint) *Service {
	var s Service
	db.Model(&Service{}).Where(&User{
		Model: gorm.Model{
			ID: id,
		},
	}).Find(&s)
	if s.ID == 0 {
		return nil
	}
	return &s
}
