package api

import (
	"strconv"

	"github.com/ZephyrJung/QiaoYiCommunity/middleware"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/ZephyrJung/QiaoYiCommunity/service"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// Create 创建帖子
// @Summary 创建帖子
// @Description 创建社区帖子
// @Tags 帖子
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreatePostReq true "创建请求"
// @Success 200 {object} response.Response{data=entity.Post} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	var req service.CreatePostReq
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

	post, err := h.postService.Create(c.Request.Context(), uid, &req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, post)
}

// GetByID 获取帖子详情
// @Summary 获取帖子详情
// @Description 根据ID获取帖子详情
// @Tags 帖子
// @Accept json
// @Produce json
// @Param id path int64 true "帖子ID"
// @Success 200 {object} response.Response{data=entity.Post} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "帖子不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/posts/{id} [get]
func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	post, err := h.postService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "post not found")
		return
	}

	response.Success(c, post)
}

// Update 更新帖子
// @Summary 更新帖子
// @Description 更新已发布的帖子
// @Tags 帖子
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "帖子ID"
// @Param body body service.UpdatePostReq true "更新请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.UpdatePostReq
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

	if err := h.postService.Update(c.Request.Context(), uid, id, &req); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除帖子
// @Summary 删除帖子
// @Description 删除指定帖子
// @Tags 帖子
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "帖子ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
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

	if err := h.postService.Delete(c.Request.Context(), uid, id); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// List 获取帖子列表
// @Summary 获取帖子列表
// @Description 获取帖子列表，支持分页和分类筛选
// @Tags 帖子
// @Accept json
// @Produce json
// @Param category_id query uint32 false "分类ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=utils.PageResult{list=[]entity.Post}} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/posts [get]
func (h *PostHandler) List(c *gin.Context) {
	var page utils.Pagination
	if err := c.ShouldBindQuery(&page); err != nil {
		page.Page = utils.DefaultPage
		page.PageSize = utils.DefaultPageSize
	}
	page.Normalize()

	var categoryID *uint32
	if cid := c.Query("category_id"); cid != "" {
		if id, err := strconv.ParseUint(cid, 10, 32); err == nil {
			v := uint32(id)
			categoryID = &v
		}
	}

	result, err := h.postService.List(c.Request.Context(), categoryID, page)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, result)
}
