package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"io"
	"net/url"
	"time"
)

func (cas *cas) addSalt(raw string) string {
	return cas.hmacSha1Encode(raw, cas.Conf.Salt)
}

func (cas *cas) hmacSha1Encode(input string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	_, _ = io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (cas *cas) generateServiceToken() string {
	return fmt.Sprintf("ST-%s", randstr.String(32))
}

// return the service's auth data, contains the service token.
func (cas *cas) getServiceAuth(serviceID uint, userID uint) *serviceAuth {
	auth := new(serviceAuth)
	cas.DB.Model(&serviceAuth{}).Where(&serviceAuth{
		ServiceID: serviceID,
		UserID:    userID,
	}).Find(&auth)
	return auth
}

// return the service redirect url with the Server Ticket.
func (cas *cas) newServiceTicketCallBack(serviceURL string, userID uint, serviceID int) string {
	// generate service ticket
	st := cas.generateServiceToken()
	// save the service ticket
	cas.Redis.Set(st, fmt.Sprintf("%d|%d", userID, serviceID), 5*time.Minute)

	// add query data into url
	requestURL, err := url.ParseRequestURI(serviceURL)
	if err != nil {
		return ""
	}
	query := requestURL.Query()
	query.Set("ticket", st)
	requestURL.RawQuery = query.Encode()
	return requestURL.String()
}

func (cas *cas) md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func (cas *cas) getServiceByURL(urlStr string) (*service, error) {
	serviceURL, err := url.ParseRequestURI(urlStr)
	if err != nil || serviceURL.Hostname() == "" {
		return nil, errors.New("参数无效")
	}
	if serviceURL.Scheme != "https" {
		return nil, errors.New("非 HTTPS 协议")
	}

	// check service whitelist
	trustDomain := new(domain)
	cas.DB.Model(&domain{}).Where("domain = ?", serviceURL.Hostname()).Find(&trustDomain)
	if trustDomain.ID == 0 {
		return nil, errors.New("域名不在白名单内")
	}

	// get service id
	serviceData := new(service)
	cas.DB.Model(&service{}).Where("id = ? and ban = ?", trustDomain.ServiceID, false).Find(&serviceData)
	if serviceData.ID == 0 {
		return nil, errors.New("域名不在白名单内")
	}
	return serviceData, nil
}

func (cas *cas) csrfMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	// set the csrf token
	if session.Get("_csrf") == nil {
		if c.Request.Method == "GET" {
			csrfToken := randstr.String(32)
			session.Set("_csrf", csrfToken)
			session.Save()
			c.Set("_csrf", csrfToken)
			c.Next()
			return
		}
		c.Redirect(302, "/")
		c.Abort()
		return
	}

	c.Set("_csrf", session.Get("_csrf"))

	// check the csrf token
	if c.Request.Method == "POST" {
		csrfForm, exist := c.GetPostForm("_csrf")
		csrfToken := session.Get("_csrf").(string)
		if !exist || csrfForm != csrfToken {
			c.Redirect(302, "/")
			c.Abort()
			return
		}
	}
	c.Next()
}
