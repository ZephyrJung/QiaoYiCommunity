package api

import (
	"github.com/ZephyrJung/QiaoYiCommunity/middleware"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/ZephyrJung/QiaoYiCommunity/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type updateProfileReq struct {
	Building string `json:"building"`
	Unit     string `json:"unit"`
	Room     string `json:"room"`
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserIDKey)
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid user")
		return
	}

	result, err := h.userService.GetProfile(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserIDKey)
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid user")
		return
	}

	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.UpdateProfile(c.Request.Context(), uid, req.Building, req.Unit, req.Room); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}
