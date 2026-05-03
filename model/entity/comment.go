package entity

import "time"

type Comment struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	PostID     int64     `gorm:"not null;index" json:"post_id"`
	UserID     int64     `gorm:"not null;index" json:"user_id"`
	ParentID   int64     `gorm:"default:0;index" json:"parent_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	LikeCount  uint32    `gorm:"default:0" json:"like_count"`
	Status     int8      `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       *User     `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentStatus int8

const (
	CommentStatusDeleted CommentStatus = 0
	CommentStatusNormal  CommentStatus = 1
)
