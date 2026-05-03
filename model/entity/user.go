package entity

import (
	"time"
)

type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UnionID   *string   `gorm:"size:64;uniqueIndex" json:"union_id,omitempty"`
	OpenID    *string   `gorm:"size:64;index" json:"open_id,omitempty"`
	Phone     *string   `gorm:"size:20;uniqueIndex" json:"phone,omitempty"`
	Nickname  string    `gorm:"size:100" json:"nickname"`
	Avatar    string    `gorm:"size:500" json:"avatar"`
	Role      int8      `gorm:"default:1" json:"role"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type UserRole int8

const (
	UserRoleResident UserRole = 1
	UserRoleManager  UserRole = 2
	UserRoleAdmin    UserRole = 3
)

type UserStatus int8

const (
	UserStatusBanned  UserStatus = 0
	UserStatusActive  UserStatus = 1
)

type UserProfile struct {
	UserID    int64     `gorm:"primaryKey" json:"user_id"`
	Building  string    `gorm:"size:100" json:"building"`
	Unit      string    `gorm:"size:50" json:"unit"`
	Room      string    `gorm:"size:50" json:"room"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
