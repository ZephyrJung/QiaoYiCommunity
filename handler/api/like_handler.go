package api

import (
	"strconv"

	"github.com/ZephyrJung/QiaoYiCommunity/middleware"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/ZephyrJung/QiaoYiCommunity/service"
	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

type toggleLikeReq struct {
	TargetType int8  `json:"target_type" binding:"required,oneof=1 2"`
	TargetID   int64 `json:"target_id" binding:"required"`
}

func (h *LikeHandler) Toggle(c *gin.Context) {
	var req toggleLikeReq
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

	liked, err := h.likeService.Toggle(c.Request.Context(), uid, req.TargetType, req.TargetID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, gin.H{"liked": liked})
}

func (h *LikeHandler) Status(c *gin.Context) {
	targetType, err := strconv.ParseInt(c.Query("target_type"), 10, 8)
	if err != nil {
		response.BadRequest(c, "invalid target_type")
		return
	}
	targetID, err := strconv.ParseInt(c.Query("target_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid target_id")
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid user")
		return
	}

	liked, err := h.likeService.Status(c.Request.Context(), uid, int8(targetType), targetID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, gin.H{"liked": liked})
}
