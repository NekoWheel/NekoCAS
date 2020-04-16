package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func (cas *cas) initDatabase() {
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local",
			cas.Conf.DBUser,
			cas.Conf.DBPassword,
			cas.Conf.DBAddr,
			cas.Conf.DBName,
		))

	if err != nil {
		log.Fatalln(err)
	}
	cas.DB = db

	cas.DB.AutoMigrate(&user{}, &service{}, &domain{}, &serviceAuth{}, )
}
