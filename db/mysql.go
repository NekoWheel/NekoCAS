package db

import (
	"fmt"

	"github.com/NekoWheel/NekoCAS/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"
)

var db *gorm.DB

func ConnDB() {
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Get().MySQL.User,
		conf.Get().MySQL.Password,
		conf.Get().MySQL.Addr,
		conf.Get().MySQL.Name,
	)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL database: %v", err)
	}

	err = db.AutoMigrate(&User{}, &Service{})
	if err != nil {
		log.Fatal("Failed to auto migrate database tables: %v", err)
	}
}
