package api

import (
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/ZephyrJung/QiaoYiCommunity/service"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

// UploadImages 上传图片
// @Summary 上传图片
// @Description 上传多张图片
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param files formData file true "图片文件"
// @Success 200 {object} response.Response "上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/upload/images [post]
func (h *UploadHandler) UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.BadRequest(c, "files field is required")
		return
	}

	urls, err := h.uploadService.UploadImages(c.Request.Context(), files)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, gin.H{"urls": urls})
}
