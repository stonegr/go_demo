package models

import (
	"go_blog/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	db, err := gorm.Open(mysql.Open(utils.GetDSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	// 自动迁移
	return DB.AutoMigrate(&Post{})
}
