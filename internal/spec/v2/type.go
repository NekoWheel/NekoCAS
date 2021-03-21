package v2

import (
	"encoding/xml"

	"github.com/NekoWheel/NekoCAS/internal/db"
)

type CASServiceResponse struct {
	XMLName      xml.Name `xml:"cas:serviceResponse"`
	Xmlns        string   `xml:"xmlns:cas,attr"`
	Success      *CASAuthenticationSuccess
	Failure      *CASAuthenticationFailure
	ProxySuccess *CASProxySuccess
	ProxyFailure *CASProxyFailure
}

type CASAuthenticationSuccess struct {
	XMLName    xml.Name `xml:"cas:authenticationSuccess"`
	User       CASUser
	Attributes CASAttributes
}

type CASAuthenticationFailure struct {
	XMLName xml.Name `xml:"cas:authenticationFailure"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:",chardata"`
}

type CASUser struct {
	XMLName xml.Name `xml:"cas:user"`
	User    string   `xml:",chardata"`
}

type CASPgtIou struct {
	XMLName xml.Name `xml:"cas:proxyGrantingTicket"`
	Ticket  string   `xml:",chardata"`
}

type CASAttributes struct {
	XMLName xml.Name `xml:"cas:attributes"`
	Email   string
}

type CASProxySuccess struct {
	XMLName xml.Name `xml:"cas:proxyTicket"`
	Ticket  string   `xml:",chardata"`
}

type CASProxyFailure struct {
	XMLName xml.Name `xml:"cas:proxyFailure"`
	Code    string   `xml:"string"`
	Message string   `xml:",chardata"`
}

// newCASResponse 创建一个新的 CAS XML 返回
func newCASResponse() CASServiceResponse {
	return CASServiceResponse{
		Xmlns: "http://www.yale.edu/tp/cas",
	}
}

// NewCASSuccessResponse 创建一个 CAS XML 成功返回，包含用户信息
func NewCASSuccessResponse(u *db.User) []byte {
	s := newCASResponse()
	s.Success = &CASAuthenticationSuccess{
		User: CASUser{User: u.NickName},
		Attributes: CASAttributes{
			Email: u.Email,
		},
	}
	x, _ := xml.Marshal(s)
	return x
}

// NewCASFailureResponse 创建一个 CAS XML 失败返回，包含错误码以及错误信息
func NewCASFailureResponse(c string, msg string) []byte {
	f := newCASResponse()
	f.Failure = &CASAuthenticationFailure{
		Code:    c,
		Message: msg,
	}
	x, _ := xml.Marshal(f)
	return x
}

// NewCASProxySuccessResponse 创建一个 CAS Proxy XML 成功返回，包含 Ticket
func NewCASProxySuccessResponse(pt string) []byte {
	s := newCASResponse()
	s.ProxySuccess = &CASProxySuccess{
		Ticket: pt,
	}
	x, _ := xml.Marshal(s)
	return x
}

// NewCASProxyFailureResponse 创建一个 CAS Proxy XML 失败返回，包含错误码以及错误信息
func NewCASProxyFailureResponse(c string, msg string) []byte {
	f := newCASResponse()
	f.ProxyFailure = &CASProxyFailure{
		Code:    c,
		Message: msg,
	}
	x, _ := xml.Marshal(f)
	return x
}
