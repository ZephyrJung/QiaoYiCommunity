package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// Response 通用API响应结构体
// @Description 通用API响应包装
// @Description code为0表示成功，非0表示错误
// @Description data字段在成功时包含业务数据，错误时可能为空
// swagger:model
//
//nolint:unused
var _swagger_Response = Response{}

type Response struct {
	Code    int         `json:"code" example:"0"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess       = 0
	CodeError         = 1
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeValidationErr = 422
	CodeTooMany       = 429
	CodeInternal      = 500
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

func ErrorWithStatus(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusBadRequest, CodeValidationErr, message)
}

func Unauthorized(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusForbidden, CodeForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusNotFound, CodeNotFound, message)
}

func InternalError(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusInternalServerError, CodeInternal, message)
}
