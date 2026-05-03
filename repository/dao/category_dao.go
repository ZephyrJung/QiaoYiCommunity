package dao

import (
	"context"
	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"gorm.io/gorm"
)

type CategoryDAO interface {
	ListByType(ctx context.Context, typ int8) ([]entity.Category, error)
	GetByID(ctx context.Context, id uint32) (*entity.Category, error)
}

type categoryDAO struct {
	db *gorm.DB
}

func NewCategoryDAO(db *gorm.DB) CategoryDAO {
	return &categoryDAO{db: db}
}

func (d *categoryDAO) ListByType(ctx context.Context, typ int8) ([]entity.Category, error) {
	var list []entity.Category
	err := d.db.WithContext(ctx).Where("type = ? AND status = ?", typ, 1).Order("sort_order").Find(&list).Error
	return list, err
}

func (d *categoryDAO) GetByID(ctx context.Context, id uint32) (*entity.Category, error) {
	var c entity.Category
	if err := d.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
