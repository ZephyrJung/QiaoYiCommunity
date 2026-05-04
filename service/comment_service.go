package service

import (
	"context"
	"fmt"

	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
)

type CommentService struct {
	commentDAO dao.CommentDAO
}

func NewCommentService(commentDAO dao.CommentDAO) *CommentService {
	return &CommentService{commentDAO: commentDAO}
}

type CreateCommentReq struct {
	PostID   int64  `json:"post_id" binding:"required"`
	ParentID int64  `json:"parent_id"`
	Content  string `json:"content" binding:"required"`
}

func (s *CommentService) Create(ctx context.Context, userID int64, req *CreateCommentReq) (*entity.Comment, error) {
	comment := &entity.Comment{
		ID:       utils.GenerateID(),
		PostID:   req.PostID,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  req.Content,
		Status:   int8(entity.CommentStatusNormal),
	}
	if comment.ParentID == 0 {
		comment.ParentID = 0
	}
	if err := s.commentDAO.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("create comment failed: %w", err)
	}
	return comment, nil
}

func (s *CommentService) ListByPostID(ctx context.Context, postID int64, page utils.Pagination) (*utils.PageResult, error) {
	page.Normalize()
	list, total, err := s.commentDAO.ListByPostID(ctx, postID, page)
	if err != nil {
		return nil, err
	}
	return utils.NewPageResult(list, total, page.Page, page.PageSize), nil
}

func (s *CommentService) Delete(ctx context.Context, userID, id int64) error {
	comment, err := s.commentDAO.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return fmt.Errorf("not owner")
	}
	return s.commentDAO.Delete(ctx, id)
}
