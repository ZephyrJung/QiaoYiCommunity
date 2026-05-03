package dao

import (
	"context"
	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"gorm.io/gorm"
)

type PostDAO interface {
	Create(ctx context.Context, post *entity.Post) error
	GetByID(ctx context.Context, id int64) (*entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, categoryID *uint32, page utils.Pagination) ([]entity.Post, int64, error)
	IncrementViewCount(ctx context.Context, id int64) error
}

type postDAO struct {
	db *gorm.DB
}

func NewPostDAO(db *gorm.DB) PostDAO {
	return &postDAO{db: db}
}

func (d *postDAO) Create(ctx context.Context, post *entity.Post) error {
	return d.db.WithContext(ctx).Create(post).Error
}

func (d *postDAO) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	var post entity.Post
	if err := d.db.WithContext(ctx).Preload("Images").Preload("User").Preload("Category").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (d *postDAO) Update(ctx context.Context, post *entity.Post) error {
	return d.db.WithContext(ctx).Save(post).Error
}

func (d *postDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Model(&entity.Post{}).Where("id = ?", id).Update("status", entity.PostStatusDeleted).Error
}

func (d *postDAO) List(ctx context.Context, categoryID *uint32, page utils.Pagination) ([]entity.Post, int64, error) {
	var posts []entity.Post
	var total int64
	query := d.db.WithContext(ctx).Model(&entity.Post{}).Where("status = ?", entity.PostStatusNormal)
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Preload("User").Preload("Category").Order("is_top DESC, created_at DESC").Offset(page.Offset()).Limit(page.PageSize).Find(&posts).Error
	return posts, total, err
}

func (d *postDAO) IncrementViewCount(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Model(&entity.Post{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
