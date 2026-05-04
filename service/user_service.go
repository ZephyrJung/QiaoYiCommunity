package service

import (
	"context"
	"errors"

	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
	"gorm.io/gorm"
)

type UserService struct {
	userDAO dao.UserDAO
}

func NewUserService(userDAO dao.UserDAO) *UserService {
	return &UserService{userDAO: userDAO}
}

type UserProfileResult struct {
	User    *entity.User        `json:"user"`
	Profile *entity.UserProfile `json:"profile,omitempty"`
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*UserProfileResult, error) {
	user, err := s.userDAO.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := s.userDAO.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &UserProfileResult{User: user}, nil
		}
		return nil, err
	}

	return &UserProfileResult{User: user, Profile: profile}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int64, building, unit, room string) error {
	profile := &entity.UserProfile{
		UserID:   userID,
		Building: building,
		Unit:     unit,
		Room:     room,
	}
	return s.userDAO.SaveProfile(ctx, profile)
}
