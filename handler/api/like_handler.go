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
	TargetType int8  `json:"target_type" binding:"required,oneof=1 2" example:"1"`
	TargetID   int64 `json:"target_id" binding:"required" example:"123456"`
}

// Toggle 点赞/取消点赞
// @Summary 点赞/取消点赞
// @Description 对帖子或评论进行点赞或取消点赞操作
// @Tags 点赞
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body toggleLikeReq true "操作请求"
// @Success 200 {object} response.Response "操作成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/likes/toggle [post]
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

// Status 获取点赞状态
// @Summary 获取点赞状态
// @Description 查询当前用户对帖子或评论的点赞状态
// @Tags 点赞
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param target_type query int8 true "目标类型(1-帖子,2-评论)" example:"1"
// @Param target_id query int64 true "目标ID" example:"123456"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/likes/status [get]
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
