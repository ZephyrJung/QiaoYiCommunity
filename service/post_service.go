package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
)

type PostService struct {
	postDAO     dao.PostDAO
	categoryDAO dao.CategoryDAO
}

func NewPostService(postDAO dao.PostDAO, categoryDAO dao.CategoryDAO) *PostService {
	return &PostService{
		postDAO:     postDAO,
		categoryDAO: categoryDAO,
	}
}

type CreatePostReq struct {
	CategoryID uint32   `json:"category_id" binding:"required"`
	Title      string   `json:"title" binding:"required,max=200"`
	Content    string   `json:"content" binding:"required"`
	ImageURLs  []string `json:"image_urls"`
}

type UpdatePostReq struct {
	CategoryID uint32   `json:"category_id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	ImageURLs  []string `json:"image_urls"`
}

func (s *PostService) Create(ctx context.Context, userID int64, req *CreatePostReq) (*entity.Post, error) {
	if _, err := s.categoryDAO.GetByID(ctx, req.CategoryID); err != nil {
		return nil, errors.New("category not found")
	}

	post := &entity.Post{
		ID:         utils.GenerateID(),
		UserID:     userID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
		Status:     int8(entity.PostStatusNormal),
	}

	if len(req.ImageURLs) > 0 {
		images := make([]entity.PostImage, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			images[i] = entity.PostImage{ID: utils.GenerateID(), PostID: post.ID, URL: url, SortOrder: i}
		}
		post.Images = images
	}

	if err := s.postDAO.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("create post failed: %w", err)
	}
	return post, nil
}

func (s *PostService) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	post, err := s.postDAO.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.postDAO.IncrementViewCount(ctx, id)
	return post, nil
}

func (s *PostService) Update(ctx context.Context, userID, id int64, req *UpdatePostReq) error {
	post, err := s.postDAO.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return errors.New("not owner")
	}

	if req.CategoryID > 0 {
		if _, err := s.categoryDAO.GetByID(ctx, req.CategoryID); err != nil {
			return errors.New("category not found")
		}
		post.CategoryID = req.CategoryID
	}
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if req.ImageURLs != nil {
		images := make([]entity.PostImage, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			images[i] = entity.PostImage{ID: utils.GenerateID(), PostID: post.ID, URL: url, SortOrder: i}
		}
		post.Images = images
	}

	return s.postDAO.Update(ctx, post)
}

func (s *PostService) Delete(ctx context.Context, userID, id int64) error {
	post, err := s.postDAO.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return errors.New("not owner")
	}
	return s.postDAO.Delete(ctx, id)
}

func (s *PostService) List(ctx context.Context, categoryID *uint32, page utils.Pagination) (*utils.PageResult, error) {
	page.Normalize()
	posts, total, err := s.postDAO.List(ctx, categoryID, page)
	if err != nil {
		return nil, err
	}
	return utils.NewPageResult(posts, total, page.Page, page.PageSize), nil
}
