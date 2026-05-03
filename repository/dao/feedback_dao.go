package dao

import (
	"context"
	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"gorm.io/gorm"
)

type FeedbackDAO interface {
	Create(ctx context.Context, fb *entity.Feedback) error
	GetByID(ctx context.Context, id int64) (*entity.Feedback, error)
	Update(ctx context.Context, fb *entity.Feedback) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, userID *int64, page utils.Pagination) ([]entity.Feedback, int64, error)
	CreateReply(ctx context.Context, reply *entity.FeedbackReply) error
}

type feedbackDAO struct {
	db *gorm.DB
}

func NewFeedbackDAO(db *gorm.DB) FeedbackDAO {
	return &feedbackDAO{db: db}
}

func (d *feedbackDAO) Create(ctx context.Context, fb *entity.Feedback) error {
	return d.db.WithContext(ctx).Create(fb).Error
}

func (d *feedbackDAO) GetByID(ctx context.Context, id int64) (*entity.Feedback, error) {
	var fb entity.Feedback
	if err := d.db.WithContext(ctx).Preload("Images").Preload("Replies.Replier").Preload("User").First(&fb, id).Error; err != nil {
		return nil, err
	}
	return &fb, nil
}

func (d *feedbackDAO) Update(ctx context.Context, fb *entity.Feedback) error {
	return d.db.WithContext(ctx).Save(fb).Error
}

func (d *feedbackDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Model(&entity.Feedback{}).Where("id = ?", id).Update("status", entity.FeedbackStatusDeleted).Error
}

func (d *feedbackDAO) List(ctx context.Context, userID *int64, page utils.Pagination) ([]entity.Feedback, int64, error) {
	var list []entity.Feedback
	var total int64
	query := d.db.WithContext(ctx).Model(&entity.Feedback{}).Where("status != ?", entity.FeedbackStatusDeleted)
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Preload("User").Order("created_at DESC").Offset(page.Offset()).Limit(page.PageSize).Find(&list).Error
	return list, total, err
}

func (d *feedbackDAO) CreateReply(ctx context.Context, reply *entity.FeedbackReply) error {
	return d.db.WithContext(ctx).Create(reply).Error
}
