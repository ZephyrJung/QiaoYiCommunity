package dao

import (
	"context"
	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"gorm.io/gorm"
)

type CommentDAO interface {
	Create(ctx context.Context, comment *entity.Comment) error
	GetByID(ctx context.Context, id int64) (*entity.Comment, error)
	Delete(ctx context.Context, id int64) error
	ListByPostID(ctx context.Context, postID int64, page utils.Pagination) ([]entity.Comment, int64, error)
}

type commentDAO struct {
	db *gorm.DB
}

func NewCommentDAO(db *gorm.DB) CommentDAO {
	return &commentDAO{db: db}
}

func (d *commentDAO) Create(ctx context.Context, comment *entity.Comment) error {
	return d.db.WithContext(ctx).Create(comment).Error
}

func (d *commentDAO) GetByID(ctx context.Context, id int64) (*entity.Comment, error) {
	var c entity.Comment
	if err := d.db.WithContext(ctx).Preload("User").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *commentDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Model(&entity.Comment{}).Where("id = ?", id).Update("status", entity.CommentStatusDeleted).Error
}

func (d *commentDAO) ListByPostID(ctx context.Context, postID int64, page utils.Pagination) ([]entity.Comment, int64, error) {
	var list []entity.Comment
	var total int64
	query := d.db.WithContext(ctx).Model(&entity.Comment{}).Where("post_id = ? AND status = ? AND parent_id = 0", postID, entity.CommentStatusNormal)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Preload("User").Order("created_at DESC").Offset(page.Offset()).Limit(page.PageSize).Find(&list).Error
	return list, total, err
}
