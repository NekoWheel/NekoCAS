package db

import (
	"gorm.io/gorm"
)

// Setting 为应用设置
type Setting struct {
	gorm.Model

	Key   string
	Value string
}

func GetSetting(key string, defaultValue ...string) (string, error) {
	var setting Setting
	if err := db.Model(&Setting{}).Where("`key` = ?", key).First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if len(defaultValue) > 0 {
				return defaultValue[0], nil
			}
		}
		return "", err
	}

	return setting.Value, nil
}

func MustGetSetting(key string, defaultValue ...string) string {
	value, _ := GetSetting(key, defaultValue...)
	return value
}

func SetSetting(key, value string) error {
	_, err := GetSetting(key)
	if err == gorm.ErrRecordNotFound {
		return db.Model(&Setting{}).Create(&Setting{
			Key:   key,
			Value: value,
		}).Error
	}
	return db.Model(&Setting{}).Where("`key` = ?", key).Update("value", value).Error
}
