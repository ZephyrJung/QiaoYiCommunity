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

// Create 创建反馈
// @Summary 创建反馈
// @Description 创建用户反馈
// @Tags 反馈
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreateFeedbackReq true "创建请求"
// @Success 200 {object} response.Response{data=entity.Feedback} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/feedbacks [post]
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

// GetByID 获取反馈详情
// @Summary 获取反馈详情
// @Description 根据ID获取反馈详情
// @Tags 反馈
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "反馈ID"
// @Success 200 {object} response.Response{data=entity.Feedback} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "反馈不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/feedbacks/{id} [get]
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

// Update 更新反馈
// @Summary 更新反馈
// @Description 更新已提交的反馈
// @Tags 反馈
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "反馈ID"
// @Param body body service.UpdateFeedbackReq true "更新请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/feedbacks/{id} [put]
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

// Delete 删除反馈
// @Summary 删除反馈
// @Description 删除指定反馈
// @Tags 反馈
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "反馈ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/feedbacks/{id} [delete]
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

// List 获取反馈列表
// @Summary 获取反馈列表
// @Description 获取反馈列表，支持分页和查询自己的反馈
// @Tags 反馈
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mine query string false "是否只查询自己的(1-是)"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=utils.PageResult{list=[]entity.Feedback}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/feedbacks [get]
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

// CreateReply 回复反馈
// @Summary 回复反馈
// @Description 对反馈进行回复
// @Tags 反馈
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "反馈ID"
// @Param body body service.CreateReplyReq true "回复内容"
// @Success 200 {object} response.Response "回复成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/feedbacks/{id}/replies [post]
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
