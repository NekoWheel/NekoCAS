package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/thanhpk/randstr"
)

// NewServiceTicket 生成一个 Ticket。
func NewServiceTicket(service *Service, user *User) (string, error) {
	ticket := "ST-" + randstr.String(32)
	err := red.Set(ticket, fmt.Sprintf("%d|%d", service.ID, user.ID), -1).Err()
	if err != nil {
		return "", err
	}
	return ticket, nil
}

// ValidateServiceTicket 验证 Ticket 是否正确。
func ValidateServiceTicket(service *Service, ticket string) (*User, bool) {
	u, s, ok := ValidateTicket(ticket)
	if !ok {
		return nil, false
	}
	if s.ID != service.ID {
		return nil, false
	}
	return u, true
}

// ValidateTicket 验证 Ticket 是否正确。
func ValidateTicket(ticket string) (*User, *Service, bool) {
	ticketData, err := red.Get(ticket).Result()
	if ticketData == "" || err != nil {
		return nil, nil, false
	}
	ticketPart := strings.Split(ticketData, "|")
	if len(ticketPart) != 2 {
		return nil, nil, false
	}

	serviceID := ticketPart[0]
	userID := ticketPart[1]

	sid, err := strconv.Atoi(serviceID)
	if err != nil {
		return nil, nil, false
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		return nil, nil, false
	}

	user := MustGetUserByID(uint(uid))
	service, err := GetServiceByID(uint(sid))
	return user, service, err == nil
}
