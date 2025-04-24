package model

import "time"

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusBlocked  = "blocked"
)

type FollowerRelation struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	FollowerID uint
	CreatedAt  time.Time
	Status     string `gorm:"default:'pending'"`
}

func NewFollowerRelation(userID, followerID uint) FollowerRelation {
	return FollowerRelation{
		UserID:     userID,
		FollowerID: followerID,
		Status:     StatusPending,
		CreatedAt:  time.Now(),
	}
}
