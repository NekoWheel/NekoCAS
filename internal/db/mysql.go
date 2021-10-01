package db

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/NekoWheel/NekoCAS/internal/conf"
)

var db *gorm.DB

func ConnDB() error {
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.MySQL.User,
		conf.MySQL.Password,
		conf.MySQL.Addr,
		conf.MySQL.Name,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.Wrap(err, "connect database")
	}

	err = db.AutoMigrate(&User{}, &Service{}, &Setting{})
	if err != nil {
		return errors.Wrap(err, "auto migrate")
	}

	return nil
}
