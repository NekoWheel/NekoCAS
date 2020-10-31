package db

import (
	"fmt"

	"github.com/NekoWheel/NekoCAS/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"
)

var db *gorm.DB

func ConnDB() {
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Get().DBUser,
		config.Get().DBPassword,
		config.Get().DBAddr,
		config.Get().DBName,
	)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL database: %v", err)
	}
	
	err = db.AutoMigrate(&User{}, &Service{}, &Domain{}, &ServiceAuth{})
	if err != nil {
		log.Fatal("Failed to auto migrate database tables: %v", err)
	}
}
