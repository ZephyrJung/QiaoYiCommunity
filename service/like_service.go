package service

import (
	"context"
	"fmt"

	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
)

type LikeService struct {
	likeDAO dao.LikeDAO
}

func NewLikeService(likeDAO dao.LikeDAO) *LikeService {
	return &LikeService{likeDAO: likeDAO}
}

func (s *LikeService) Toggle(ctx context.Context, userID int64, targetType int8, targetID int64) (liked bool, err error) {
	exists, err := s.likeDAO.Exists(ctx, userID, targetType, targetID)
	if err != nil {
		return false, fmt.Errorf("check like failed: %w", err)
	}

	if exists {
		if err := s.likeDAO.Delete(ctx, userID, targetType, targetID); err != nil {
			return false, fmt.Errorf("unlike failed: %w", err)
		}
		return false, nil
	}

	like := &entity.Like{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
	}
	if err := s.likeDAO.Create(ctx, like); err != nil {
		return false, fmt.Errorf("like failed: %w", err)
	}
	return true, nil
}

func (s *LikeService) Status(ctx context.Context, userID int64, targetType int8, targetID int64) (bool, error) {
	return s.likeDAO.Exists(ctx, userID, targetType, targetID)
}
