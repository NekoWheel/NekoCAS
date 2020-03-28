package main

import "github.com/jinzhu/gorm"

type user struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Token    string `json:"token"`
}

type serviceTicket struct {
	gorm.Model
	AppID  uint
	Ticket string
	UserID uint
}

type service struct {
	gorm.Model
	Name   string
	Secret string
	Avatar string
}

type domain struct {
	gorm.Model
	Domain    string
	ServiceID uint
}

type serviceAuth struct {
	gorm.Model
	ServiceID uint
	UserID    uint
	Token     string
}
