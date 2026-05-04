package api

import (
	"strconv"

	"github.com/ZephyrJung/QiaoYiCommunity/middleware"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/ZephyrJung/QiaoYiCommunity/service"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"github.com/gin-gonic/gin"
)

type FeedbackHandler struct {
	feedbackService *service.FeedbackService
}

func NewFeedbackHandler(feedbackService *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{feedbackService: feedbackService}
}

func (h *FeedbackHandler) Create(c *gin.Context) {
	var req service.CreateFeedbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid user")
		return
	}

	fb, err := h.feedbackService.Create(c.Request.Context(), uid, &req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, fb)
}

func (h *FeedbackHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	fb, err := h.feedbackService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "feedback not found")
		return
	}

	response.Success(c, fb)
}

func (h *FeedbackHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.UpdateFeedbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid user")
		return
	}

	if err := h.feedbackService.Update(c.Request.Context(), uid, id, &req); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *FeedbackHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid user")
		return
	}

	if err := h.feedbackService.Delete(c.Request.Context(), uid, id); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *FeedbackHandler) List(c *gin.Context) {
	var page utils.Pagination
	if err := c.ShouldBindQuery(&page); err != nil {
		page.Page = utils.DefaultPage
		page.PageSize = utils.DefaultPageSize
	}
	page.Normalize()

	var userID *int64
	if c.Query("mine") == "1" {
		uid, ok := c.Get(middleware.ContextUserIDKey)
		if id, ok2 := uid.(int64); ok && ok2 {
			userID = &id
		}
	}

	result, err := h.feedbackService.List(c.Request.Context(), userID, page)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *FeedbackHandler) CreateReply(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.CreateReplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid user")
		return
	}

	if err := h.feedbackService.CreateReply(c.Request.Context(), id, uid, &req); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}
