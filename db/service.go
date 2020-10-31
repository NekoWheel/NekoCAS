package db

import (
	"net/url"

	"github.com/pkg/errors"
)

// ServiceByURL 通过 ServiceURL 查找对应的服务
func ServiceByURL(u string) (*Service, error) {
	serviceURL, err := url.ParseRequestURI(u)
	if err != nil || serviceURL.Hostname() == "" {
		return nil, errors.New("参数无效")
	}

	// HTTPS 检测
	//if serviceURL.Scheme != "https" {
	//	return nil, errors.New("非 HTTPS 协议")
	//}

	// Check service whitelist
	trustDomain := new(Domain)
	db.Model(&Domain{}).Where("domain = ?", serviceURL.Hostname()).Find(&trustDomain)
	if trustDomain.ID == 0 {
		return nil, errors.New("域名不在白名单内")
	}

	serviceData := new(Service)
	db.Model(&Service{}).Where("id = ? and ban = ?", trustDomain.ServiceID, false).Find(&serviceData)
	if serviceData.ID == 0 {
		return nil, errors.New("域名不在白名单内")
	}
	return serviceData, nil
}
