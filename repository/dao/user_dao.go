package dao

import (
	"context"
	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"gorm.io/gorm"
)

type UserDAO interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByUnionID(ctx context.Context, unionID string) (*entity.User, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

type userDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAO{db: db}
}

func (d *userDAO) Create(ctx context.Context, user *entity.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *userDAO) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	if err := d.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAO) GetByUnionID(ctx context.Context, unionID string) (*entity.User, error) {
	var user entity.User
	if err := d.db.WithContext(ctx).Where("union_id = ?", unionID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAO) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	if err := d.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAO) Update(ctx context.Context, user *entity.User) error {
	return d.db.WithContext(ctx).Save(user).Error
}
