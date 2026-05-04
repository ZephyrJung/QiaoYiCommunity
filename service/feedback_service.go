package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
)

type FeedbackService struct {
	feedbackDAO dao.FeedbackDAO
}

func NewFeedbackService(feedbackDAO dao.FeedbackDAO) *FeedbackService {
	return &FeedbackService{feedbackDAO: feedbackDAO}
}

type CreateFeedbackReq struct {
	Type         int8     `json:"type" binding:"required"`
	Title        string   `json:"title" binding:"required,max=200"`
	Description  string   `json:"description" binding:"required"`
	Location     string   `json:"location"`
	ContactPhone string   `json:"contact_phone"`
	ImageURLs    []string `json:"image_urls"`
}

type UpdateFeedbackReq struct {
	Type         int8     `json:"type"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Location     string   `json:"location"`
	ContactPhone string   `json:"contact_phone"`
	Status       int8     `json:"status"`
	ImageURLs    []string `json:"image_urls"`
}

type CreateReplyReq struct {
	Content string `json:"content" binding:"required"`
}

func (s *FeedbackService) Create(ctx context.Context, userID int64, req *CreateFeedbackReq) (*entity.Feedback, error) {
	fb := &entity.Feedback{
		ID:           utils.GenerateID(),
		UserID:       userID,
		Type:         req.Type,
		Title:        req.Title,
		Description:  req.Description,
		Location:     req.Location,
		ContactPhone: req.ContactPhone,
		Status:       int8(entity.FeedbackStatusPending),
	}

	if len(req.ImageURLs) > 0 {
		images := make([]entity.FeedbackImage, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			images[i] = entity.FeedbackImage{ID: utils.GenerateID(), FeedbackID: fb.ID, URL: url}
		}
		fb.Images = images
	}

	if err := s.feedbackDAO.Create(ctx, fb); err != nil {
		return nil, fmt.Errorf("create feedback failed: %w", err)
	}
	return fb, nil
}

func (s *FeedbackService) GetByID(ctx context.Context, id int64) (*entity.Feedback, error) {
	return s.feedbackDAO.GetByID(ctx, id)
}

func (s *FeedbackService) Update(ctx context.Context, userID, id int64, req *UpdateFeedbackReq) error {
	fb, err := s.feedbackDAO.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if fb.UserID != userID {
		return errors.New("not owner")
	}

	if req.Type > 0 {
		fb.Type = req.Type
	}
	if req.Title != "" {
		fb.Title = req.Title
	}
	if req.Description != "" {
		fb.Description = req.Description
	}
	fb.Location = req.Location
	fb.ContactPhone = req.ContactPhone
	if req.Status > 0 {
		fb.Status = req.Status
	}

	if req.ImageURLs != nil {
		images := make([]entity.FeedbackImage, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			images[i] = entity.FeedbackImage{ID: utils.GenerateID(), FeedbackID: fb.ID, URL: url}
		}
		fb.Images = images
	}

	return s.feedbackDAO.Update(ctx, fb)
}

func (s *FeedbackService) Delete(ctx context.Context, userID, id int64) error {
	fb, err := s.feedbackDAO.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if fb.UserID != userID {
		return errors.New("not owner")
	}
	return s.feedbackDAO.Delete(ctx, id)
}

func (s *FeedbackService) List(ctx context.Context, userID *int64, page utils.Pagination) (*utils.PageResult, error) {
	page.Normalize()
	list, total, err := s.feedbackDAO.List(ctx, userID, page)
	if err != nil {
		return nil, err
	}
	return utils.NewPageResult(list, total, page.Page, page.PageSize), nil
}

func (s *FeedbackService) CreateReply(ctx context.Context, feedbackID int64, replierID int64, req *CreateReplyReq) error {
	reply := &entity.FeedbackReply{
		ID:         utils.GenerateID(),
		FeedbackID: feedbackID,
		ReplierID:  replierID,
		Content:    req.Content,
	}
	return s.feedbackDAO.CreateReply(ctx, reply)
}
