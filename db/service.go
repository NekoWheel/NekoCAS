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

// CreateService 新建一个新的服务
func CreateService(s *Service) error {
	isExist := IsServiceExist(s.Name)
	if isExist {
		return errors.Errorf("服务名已存在")
	}

	tx := db.Begin()
	if tx.Create(s).RowsAffected != 1 {
		tx.Rollback()
		return errors.Errorf("数据库错误")
	}
	tx.Commit()
	return nil
}

// UpdateService 更新一个服务
func UpdateService(s *Service) error {
	service := GetServiceByID(s.ID)
	if service == nil {
		return errors.Errorf("服务不存在")
	}

	tx := db.Begin()
	if err := tx.Model(&Service{}).Where("id = ?", s.ID).Updates(map[string]interface{}{
		"Name":   s.Name,
		"Avatar": s.Avatar,
		"Domain": s.Domain,
		"Ban":    s.Ban,
	}).Error; err != nil {
		tx.Rollback()
		return errors.Errorf("数据库错误")
	}
	tx.Commit()
	return nil
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

// DeleteService 根据指定 Service ID 删除对应的服务
func DeleteService(id uint) error {
	return db.Model(&Service{}).Where("id = ?", id).Delete(&Service{}).Error
}

// GetServices 批量获取服务
// options[0] offset
// options[1] limit
func GetServices(options ...int) []*Service {
	var services []*Service

	if len(options) == 0 {
		db.Model(&Service{}).Find(&services)
	} else {
		offset := 0
		if len(options) > 1 && options[0] > 0 {
			offset = options[0]
		}

		limit := 0
		if len(options) == 2 && options[1] > 0 {
			limit = options[1]
		}
		db.Model(&Service{}).Offset(offset).Limit(limit).Find(&services)
	}

	return services
}

// CountServices 返回服务的总数
func CountServices() int64 {
	var count int64
	db.Model(&Service{}).Count(&count)
	return count
}

// IsServiceExist 检查服务名是否重复
func IsServiceExist(name string) bool {
	if name == "" {
		return false
	}
	var s Service
	db.Model(&Service{}).Where(&Service{Name: name}).Find(&s)
	return s.ID != 0
}
