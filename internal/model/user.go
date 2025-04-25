package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string
	AvatarURL      string
	Bio            string
	FollowersCount uint               `gorm:"default:0"`
	FollowingCount uint               `gorm:"default:0"`
	Followers      []FollowerRelation `gorm:"foreignKey:UserID" json:"omitempty"`
	Following      []FollowerRelation `gorm:"foreignKey:FollowerID" json:"omitempty"`
	Settings       Settings           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"omitempty"`
}
