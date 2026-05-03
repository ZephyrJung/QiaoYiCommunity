package entity

import "time"

type Feedback struct {
	ID           int64            `gorm:"primaryKey" json:"id"`
	UserID       int64            `gorm:"not null;index" json:"user_id"`
	Type         int8             `gorm:"not null" json:"type"`
	Title        string           `gorm:"size:200;not null" json:"title"`
	Description  string           `gorm:"type:text;not null" json:"description"`
	Location     string           `gorm:"size:200" json:"location"`
	Status       int8             `gorm:"default:1" json:"status"`
	ContactPhone string           `gorm:"size:20" json:"contact_phone"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Images       []FeedbackImage  `gorm:"foreignKey:FeedbackID" json:"images,omitempty"`
	Replies      []FeedbackReply  `gorm:"foreignKey:FeedbackID" json:"replies,omitempty"`
	User         *User            `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}

type FeedbackType int8

const (
	FeedbackTypeFacility  FeedbackType = 1
	FeedbackTypeSecurity  FeedbackType = 2
	FeedbackTypeHygiene   FeedbackType = 3
	FeedbackTypeNoise     FeedbackType = 4
	FeedbackTypeOther     FeedbackType = 5
)

type FeedbackStatus int8

const (
	FeedbackStatusDeleted    FeedbackStatus = 0
	FeedbackStatusPending    FeedbackStatus = 1
	FeedbackStatusProcessing FeedbackStatus = 2
	FeedbackStatusResolved   FeedbackStatus = 3
)

type FeedbackImage struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	FeedbackID int64     `gorm:"not null;index" json:"feedback_id"`
	URL        string    `gorm:"size:500;not null" json:"url"`
	CreatedAt  time.Time `json:"created_at"`
}

func (FeedbackImage) TableName() string {
	return "feedback_images"
}

type FeedbackReply struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	FeedbackID int64     `gorm:"not null;index" json:"feedback_id"`
	ReplierID  int64     `gorm:"not null;index" json:"replier_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	Replier    *User     `gorm:"foreignKey:ReplierID;references:ID" json:"replier,omitempty"`
}

func (FeedbackReply) TableName() string {
	return "feedback_replies"
}
