package entity

import "time"

type Post struct {
	ID           int64         `gorm:"primaryKey" json:"id"`
	UserID       int64         `gorm:"not null;index" json:"user_id"`
	CategoryID   uint32        `gorm:"not null;index" json:"category_id"`
	Title        string        `gorm:"size:200;not null" json:"title"`
	Content      string        `gorm:"type:text;not null" json:"content"`
	ViewCount    uint32        `gorm:"default:0" json:"view_count"`
	LikeCount    uint32        `gorm:"default:0" json:"like_count"`
	CommentCount uint32        `gorm:"default:0" json:"comment_count"`
	IsTop        int8          `gorm:"default:0" json:"is_top"`
	Status       int8          `gorm:"default:1" json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Images       []PostImage   `gorm:"foreignKey:PostID" json:"images,omitempty"`
	User         *User         `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Category     *Category     `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

func (Post) TableName() string {
	return "posts"
}

type PostStatus int8

const (
	PostStatusDeleted PostStatus = 0
	PostStatusNormal  PostStatus = 1
)

type PostImage struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	PostID    int64     `gorm:"not null;index" json:"post_id"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

func (PostImage) TableName() string {
	return "post_images"
}
