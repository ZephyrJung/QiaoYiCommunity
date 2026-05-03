package dao

import (
	"context"
	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"gorm.io/gorm"
)

type LikeDAO interface {
	Create(ctx context.Context, like *entity.Like) error
	Delete(ctx context.Context, userID int64, targetType int8, targetID int64) error
	Exists(ctx context.Context, userID int64, targetType int8, targetID int64) (bool, error)
}

type likeDAO struct {
	db *gorm.DB
}

func NewLikeDAO(db *gorm.DB) LikeDAO {
	return &likeDAO{db: db}
}

func (d *likeDAO) Create(ctx context.Context, like *entity.Like) error {
	return d.db.WithContext(ctx).Create(like).Error
}

func (d *likeDAO) Delete(ctx context.Context, userID int64, targetType int8, targetID int64) error {
	return d.db.WithContext(ctx).Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Delete(&entity.Like{}).Error
}

func (d *likeDAO) Exists(ctx context.Context, userID int64, targetType int8, targetID int64) (bool, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&entity.Like{}).Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Count(&count).Error
	return count > 0, err
}
