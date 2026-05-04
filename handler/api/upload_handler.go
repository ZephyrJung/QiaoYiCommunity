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
