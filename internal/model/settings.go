package model

import "gorm.io/gorm"

type Settings struct {
	gorm.Model
	UserID    uint   `gorm:"uniqueIndex"`
	IsPrivate bool   `gorm:"default:false"`
	DarkMode  bool   `gorm:"default:false"`
	Language  string `gorm:"default:'en'"`
}
