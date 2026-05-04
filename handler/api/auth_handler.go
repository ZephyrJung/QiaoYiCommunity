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
	JsCode string `json:"js_code" binding:"required" example:"0134567890abcdef"`
}

type phoneSendCodeReq struct {
	Phone string `json:"phone" binding:"required" example:"13800138000"`
}

type phoneLoginReq struct {
	Phone string `json:"phone" binding:"required" example:"13800138000"`
	Code  string `json:"code" binding:"required" example:"123456"`
}

type refreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// WeChatLogin 微信小程序登录
// @Summary 微信小程序登录
// @Description 使用微信小程序code进行登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body wechatLoginReq true "登录请求"
// @Success 200 {object} response.Response{data=service.LoginResult} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/auth/wechat/login [post]
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

// PhoneSendCode 发送短信验证码
// @Summary 发送短信验证码
// @Description 向指定手机号发送验证码
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body phoneSendCodeReq true "请求参数"
// @Success 200 {object} response.Response "发送成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/auth/phone/send-code [post]
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

// PhoneLogin 手机验证码登录
// @Summary 手机验证码登录
// @Description 使用手机号和验证码进行登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body phoneLoginReq true "登录请求"
// @Success 200 {object} response.Response{data=service.LoginResult} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/auth/phone/login [post]
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

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 使用refresh token获取新的access token
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body refreshTokenReq true "刷新请求"
// @Success 200 {object} response.Response{data=service.TokenPair} "刷新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/auth/refresh [post]
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

// Logout 退出登录
// @Summary 退出登录
// @Description 用户退出登录，使refresh token失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "退出成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/auth/logout [post]
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
