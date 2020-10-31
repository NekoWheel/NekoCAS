package db

import "gorm.io/gorm"



type Service struct {
	gorm.Model

	Name   string
	Secret string
	Avatar string
	Ban    bool
}

type Domain struct {
	gorm.Model

	Domain    string
	ServiceID uint
}

type ServiceAuth struct {
	gorm.Model

	ServiceID uint
	UserID    uint
	Token     string
}
