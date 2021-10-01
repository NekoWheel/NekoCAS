package db

import (
	"net/url"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Service 为接入的服务。
type Service struct {
	gorm.Model

	Name   string
	Avatar string // 服务 Logo
	Domain string // 白名单域名
	Ban    bool   // 是否封禁
}

var ErrServiceExists = errors.New("服务已存在")

// CreateService 新建一个新的服务。
func CreateService(s *Service) error {
	isExist := IsServiceExist(s.Name)
	if isExist {
		return ErrServiceExists
	}

	if err := db.Create(s).Error; err != nil {
		return errors.Wrap(err, "添加新服务")
	}
	return nil
}

var ErrorServiceNotFound = errors.New("服务不存在")

// UpdateService 更新一个服务。
func UpdateService(s *Service) error {
	_, err := GetServiceByID(s.ID)
	if err != nil {
		return err
	}

	if err := db.Model(&Service{}).Where("id = ?", s.ID).Updates(map[string]interface{}{
		"Name":   s.Name,
		"Avatar": s.Avatar,
		"Domain": s.Domain,
		"Ban":    s.Ban,
	}).Error; err != nil {
		return errors.Wrap(err, "更新服务")
	}
	return nil
}

// ServiceByURL 通过 ServiceURL 查找对应的服务。
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

// GetServiceByID 根据对应的 ServiceID 查找对应的服务。
func GetServiceByID(id uint) (*Service, error) {
	var s Service
	if err := db.Model(&Service{}).Where(&User{
		Model: gorm.Model{
			ID: id,
		},
	}).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrorServiceNotFound
		}
		return nil, err
	}
	return &s, nil
}

// DeleteService 根据指定 Service ID 删除对应的服务。
func DeleteService(id uint) error {
	return db.Model(&Service{}).Where("id = ?", id).Delete(&Service{}).Error
}

// GetServices 批量获取服务。
// options[0] offset
// options[1] limit
func GetServices(options ...int) ([]*Service, error) {
	var services []*Service

	if len(options) == 0 {
		if err := db.Model(&Service{}).Find(&services).Error; err != nil {
			return nil, err
		}
	} else {
		offset := 0
		if len(options) > 1 && options[0] > 0 {
			offset = options[0]
		}

		limit := 0
		if len(options) == 2 && options[1] > 0 {
			limit = options[1]
		}

		if err := db.Model(&Service{}).Offset(offset).Limit(limit).Find(&services).Error; err != nil {
			return nil, err
		}
	}

	return services, nil
}

// CountServices 返回服务的总数。
func CountServices() int64 {
	var count int64
	db.Model(&Service{}).Count(&count)
	return count
}

// IsServiceExist 检查服务名是否重复。
func IsServiceExist(name string) bool {
	if name == "" {
		return false
	}

	var s Service
	err := db.Model(&Service{}).Where(&Service{Name: name}).First(&s).Error
	return err == nil
}
