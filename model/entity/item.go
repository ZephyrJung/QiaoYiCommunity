package entity

import (
	"time"
)

type Item struct {
	ID            int64          `gorm:"primaryKey" json:"id"`
	UserID        int64          `gorm:"not null;index" json:"user_id"`
	CategoryID    uint32         `gorm:"not null;index" json:"category_id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Price         float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice *float64       `gorm:"type:decimal(10,2)" json:"original_price,omitempty"`
	Status        int8           `gorm:"default:1" json:"status"`
	ContactPhone  string         `gorm:"size:20" json:"contact_phone"`
	ContactWechat string         `gorm:"size:100" json:"contact_wechat"`
	ViewCount     uint32         `gorm:"default:0" json:"view_count"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Images        []ItemImage    `gorm:"foreignKey:ItemID" json:"images,omitempty"`
	User          *User          `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Category      *Category      `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

func (Item) TableName() string {
	return "items"
}

type ItemStatus int8

const (
	ItemStatusDeleted   ItemStatus = 0
	ItemStatusAvailable ItemStatus = 1
	ItemStatusSold      ItemStatus = 2
)

type ItemImage struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ItemID    int64     `gorm:"not null;index" json:"item_id"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

func (ItemImage) TableName() string {
	return "item_images"
}
