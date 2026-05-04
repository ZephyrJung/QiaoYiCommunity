package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ZephyrJung/QiaoYiCommunity/model/entity"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
)

type ItemService struct {
	itemDAO     dao.ItemDAO
	categoryDAO dao.CategoryDAO
}

func NewItemService(itemDAO dao.ItemDAO, categoryDAO dao.CategoryDAO) *ItemService {
	return &ItemService{
		itemDAO:     itemDAO,
		categoryDAO: categoryDAO,
	}
}

type CreateItemReq struct {
	CategoryID    uint32   `json:"category_id" binding:"required"`
	Title         string   `json:"title" binding:"required,max=200"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required,min=0"`
	OriginalPrice *float64 `json:"original_price,omitempty"`
	ContactPhone  string   `json:"contact_phone"`
	ContactWechat string   `json:"contact_wechat"`
	ImageURLs     []string `json:"image_urls"`
}

type UpdateItemReq struct {
	CategoryID    uint32   `json:"category_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	OriginalPrice *float64 `json:"original_price,omitempty"`
	ContactPhone  string   `json:"contact_phone"`
	ContactWechat string   `json:"contact_wechat"`
	ImageURLs     []string `json:"image_urls"`
}

func (s *ItemService) Create(ctx context.Context, userID int64, req *CreateItemReq) (*entity.Item, error) {
	if _, err := s.categoryDAO.GetByID(ctx, req.CategoryID); err != nil {
		return nil, errors.New("category not found")
	}

	item := &entity.Item{
		ID:            utils.GenerateID(),
		UserID:        userID,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		ContactPhone:  req.ContactPhone,
		ContactWechat: req.ContactWechat,
		Status:        int8(entity.ItemStatusAvailable),
	}

	if len(req.ImageURLs) > 0 {
		images := make([]entity.ItemImage, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			images[i] = entity.ItemImage{ID: utils.GenerateID(), ItemID: item.ID, URL: url, SortOrder: i}
		}
		item.Images = images
	}

	if err := s.itemDAO.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("create item failed: %w", err)
	}
	return item, nil
}

func (s *ItemService) GetByID(ctx context.Context, id int64) (*entity.Item, error) {
	return s.itemDAO.GetByID(ctx, id)
}

func (s *ItemService) Update(ctx context.Context, userID, id int64, req *UpdateItemReq) error {
	item, err := s.itemDAO.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return errors.New("not owner")
	}

	if req.CategoryID > 0 {
		if _, err := s.categoryDAO.GetByID(ctx, req.CategoryID); err != nil {
			return errors.New("category not found")
		}
		item.CategoryID = req.CategoryID
	}
	if req.Title != "" {
		item.Title = req.Title
	}
	item.Description = req.Description
	if req.Price > 0 {
		item.Price = req.Price
	}
	item.OriginalPrice = req.OriginalPrice
	item.ContactPhone = req.ContactPhone
	item.ContactWechat = req.ContactWechat

	if req.ImageURLs != nil {
		images := make([]entity.ItemImage, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			images[i] = entity.ItemImage{ID: utils.GenerateID(), ItemID: item.ID, URL: url, SortOrder: i}
		}
		item.Images = images
	}

	return s.itemDAO.Update(ctx, item)
}

func (s *ItemService) Delete(ctx context.Context, userID, id int64) error {
	item, err := s.itemDAO.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return errors.New("not owner")
	}
	return s.itemDAO.Delete(ctx, id)
}

func (s *ItemService) MarkSold(ctx context.Context, userID, id int64) error {
	item, err := s.itemDAO.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return errors.New("not owner")
	}
	item.Status = int8(entity.ItemStatusSold)
	return s.itemDAO.Update(ctx, item)
}

func (s *ItemService) List(ctx context.Context, filter dao.ItemFilter, page utils.Pagination) (*utils.PageResult, error) {
	page.Normalize()
	items, total, err := s.itemDAO.List(ctx, filter, page)
	if err != nil {
		return nil, err
	}
	return utils.NewPageResult(items, total, page.Page, page.PageSize), nil
}
