package entity

import "time"

type Category struct {
	ID         uint32    `gorm:"primaryKey" json:"id"`
	Type       int8      `gorm:"not null;index" json:"type"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Icon       string    `gorm:"size:500" json:"icon"`
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	Status     int8      `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Category) TableName() string {
	return "categories"
}

type CategoryType int8

const (
	CategoryTypeItem CategoryType = 1
	CategoryTypePost CategoryType = 2
)
