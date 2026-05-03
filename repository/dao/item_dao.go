package dao

import (
	"context"
	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"gorm.io/gorm"
)

type ItemDAO interface {
	Create(ctx context.Context, item *entity.Item) error
	GetByID(ctx context.Context, id int64) (*entity.Item, error)
	Update(ctx context.Context, item *entity.Item) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter ItemFilter, page utils.Pagination) ([]entity.Item, int64, error)
}

type ItemFilter struct {
	CategoryID *uint32
	Keyword    string
	MinPrice   *float64
	MaxPrice   *float64
	Status     *int8
	UserID     *int64
}

type itemDAO struct {
	db *gorm.DB
}

func NewItemDAO(db *gorm.DB) ItemDAO {
	return &itemDAO{db: db}
}

func (d *itemDAO) Create(ctx context.Context, item *entity.Item) error {
	return d.db.WithContext(ctx).Create(item).Error
}

func (d *itemDAO) GetByID(ctx context.Context, id int64) (*entity.Item, error) {
	var item entity.Item
	if err := d.db.WithContext(ctx).Preload("Images").Preload("User").Preload("Category").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *itemDAO) Update(ctx context.Context, item *entity.Item) error {
	return d.db.WithContext(ctx).Save(item).Error
}

func (d *itemDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Model(&entity.Item{}).Where("id = ?", id).Update("status", entity.ItemStatusDeleted).Error
}

func (d *itemDAO) List(ctx context.Context, filter ItemFilter, page utils.Pagination) ([]entity.Item, int64, error) {
	var items []entity.Item
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.Item{}).Where("status != ?", entity.ItemStatusDeleted)
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+filter.Keyword+"%")
	}
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Images").Preload("User").Order("created_at DESC").Offset(page.Offset()).Limit(page.PageSize).Find(&items).Error
	return items, total, err
}
