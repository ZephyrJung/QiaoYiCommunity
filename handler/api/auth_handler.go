package api

import (
	"github.com/ZephyrJung/QiaoYiCommunity/middleware"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/ZephyrJung/QiaoYiCommunity/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type wechatLoginReq struct {
	JsCode string `json:"js_code" binding:"required"`
}

type phoneSendCodeReq struct {
	Phone string `json:"phone" binding:"required"`
}

type phoneLoginReq struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type refreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) WeChatLogin(c *gin.Context) {
	var req wechatLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.WeChatLogin(c.Request.Context(), req.JsCode)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *AuthHandler) PhoneSendCode(c *gin.Context) {
	var req phoneSendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.authService.PhoneSendCode(c.Request.Context(), req.Phone); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *AuthHandler) PhoneLogin(c *gin.Context) {
	var req phoneLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.PhoneLogin(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, tp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserIDKey)
	tokenID, _ := c.Get(middleware.ContextTokenIDKey)
	uid, ok1 := userID.(int64)
	tid, ok2 := tokenID.(string)
	if !ok1 || !ok2 {
		response.Unauthorized(c, "invalid token context")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), uid, tid); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}
