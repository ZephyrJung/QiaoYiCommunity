package api

import (
	"strconv"

	"github.com/ZephyrJung/QiaoYiCommunity/middleware"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/ZephyrJung/QiaoYiCommunity/service"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// Create 创建评论
// @Summary 创建评论
// @Description 为帖子创建评论
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreateCommentReq true "创建请求"
// @Success 200 {object} response.Response{data=entity.Comment} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/comments [post]
func (h *CommentHandler) Create(c *gin.Context) {
	var req service.CreateCommentReq
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

	comment, err := h.commentService.Create(c.Request.Context(), uid, &req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, comment)
}

// ListByPostID 获取帖子评论列表
// @Summary 获取帖子评论列表
// @Description 根据帖子ID获取评论列表
// @Tags 评论
// @Accept json
// @Produce json
// @Param post_id path int64 true "帖子ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=utils.PageResult{list=[]entity.Comment}} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/comments/post/{post_id} [get]
func (h *CommentHandler) ListByPostID(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("post_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post_id")
		return
	}

	var page utils.Pagination
	if err := c.ShouldBindQuery(&page); err != nil {
		page.Page = utils.DefaultPage
		page.PageSize = utils.DefaultPageSize
	}
	page.Normalize()

	result, err := h.commentService.ListByPostID(c.Request.Context(), postID, page)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, result)
}

// Delete 删除评论
// @Summary 删除评论
// @Description 删除指定评论
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "评论ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/comments/{id} [delete]
func (h *CommentHandler) Delete(c *gin.Context) {
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

	if err := h.commentService.Delete(c.Request.Context(), uid, id); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}
