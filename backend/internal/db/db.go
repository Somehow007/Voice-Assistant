package db

import (
	"Voice-Assistant/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// 自动迁移表结构
	db.AutoMigrate(&model.Assistant{}, &model.History{})
	return db, err
}
