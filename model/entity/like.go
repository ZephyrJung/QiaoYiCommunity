package entity

import "time"

type Like struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     int64     `gorm:"not null" json:"user_id"`
	TargetType int8      `gorm:"not null" json:"target_type"`
	TargetID   int64     `gorm:"not null" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Like) TableName() string {
	return "likes"
}

type LikeTargetType int8

const (
	LikeTargetPost    LikeTargetType = 1
	LikeTargetComment LikeTargetType = 2
)
