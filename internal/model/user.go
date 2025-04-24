package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string
	AvatarURL string
	Bio       string
	Followers []FollowerRelation `gorm:"foreignKey:UserID"`
	Following []FollowerRelation `gorm:"foreignKey:FollowerID"`
	Settings  Settings           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
