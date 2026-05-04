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
