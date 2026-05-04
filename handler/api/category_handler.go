package api

import (
	"strconv"

	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryDAO dao.CategoryDAO
}

func NewCategoryHandler(categoryDAO dao.CategoryDAO) *CategoryHandler {
	return &CategoryHandler{categoryDAO: categoryDAO}
}

func (h *CategoryHandler) List(c *gin.Context) {
	typ, _ := strconv.ParseInt(c.Query("type"), 10, 8)
	list, err := h.categoryDAO.ListByType(c.Request.Context(), int8(typ))
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}
	response.Success(c, list)
}
