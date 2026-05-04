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

// List 获取分类列表
// @Summary 获取分类列表
// @Description 获取帖子或物品的分类列表
// @Tags 分类
// @Accept json
// @Produce json
// @Param type query int8 false "分类类型(1-帖子分类,2-物品分类)"
// @Success 200 {object} response.Response{data=[]entity.Category} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	typ, _ := strconv.ParseInt(c.Query("type"), 10, 8)
	list, err := h.categoryDAO.ListByType(c.Request.Context(), int8(typ))
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}
	response.Success(c, list)
}
