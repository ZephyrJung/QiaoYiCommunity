package api

import (
	"strconv"

	"github.com/ZephyrJung/QiaoYiCommunity/middleware"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
	"github.com/ZephyrJung/QiaoYiCommunity/service"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	itemService *service.ItemService
}

func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{itemService: itemService}
}

// Create 创建物品
// @Summary 创建物品
// @Description 创建二手物品发布
// @Tags 物品
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreateItemReq true "创建请求"
// @Success 200 {object} response.Response{data=entity.Item} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/items [post]
func (h *ItemHandler) Create(c *gin.Context) {
	var req service.CreateItemReq
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

	item, err := h.itemService.Create(c.Request.Context(), uid, &req)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, item)
}

// GetByID 获取物品详情
// @Summary 获取物品详情
// @Description 根据ID获取物品详情
// @Tags 物品
// @Accept json
// @Produce json
// @Param id path int64 true "物品ID"
// @Success 200 {object} response.Response{data=entity.Item} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "物品不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/items/{id} [get]
func (h *ItemHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	item, err := h.itemService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "item not found")
		return
	}

	response.Success(c, item)
}

// Update 更新物品
// @Summary 更新物品
// @Description 更新已发布的物品信息
// @Tags 物品
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "物品ID"
// @Param body body service.UpdateItemReq true "更新请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/items/{id} [put]
func (h *ItemHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.UpdateItemReq
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

	if err := h.itemService.Update(c.Request.Context(), uid, id, &req); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除物品
// @Summary 删除物品
// @Description 删除指定物品
// @Tags 物品
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "物品ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/items/{id} [delete]
func (h *ItemHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid user")
		return
	}

	if err := h.itemService.Delete(c.Request.Context(), uid, id); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// MarkSold 标记物品已售出
// @Summary 标记物品已售出
// @Description 将物品状态标记为已售出
// @Tags 物品
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int64 true "物品ID"
// @Success 200 {object} response.Response "标记成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/items/{id}/sold [patch]
func (h *ItemHandler) MarkSold(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	userID, _ := c.Get(middleware.ContextUserIDKey)
	uid, ok := userID.(int64)
	if !ok {
		response.Unauthorized(c, "invalid user")
		return
	}

	if err := h.itemService.MarkSold(c.Request.Context(), uid, id); err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

// List 获取物品列表
// @Summary 获取物品列表
// @Description 获取二手物品列表，支持分页、分类和价格筛选
// @Tags 物品
// @Accept json
// @Produce json
// @Param category_id query uint32 false "分类ID"
// @Param keyword query string false "关键词"
// @Param min_price query float64 false "最低价格"
// @Param max_price query float64 false "最高价格"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=utils.PageResult{list=[]entity.Item}} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/items [get]
func (h *ItemHandler) List(c *gin.Context) {
	var page utils.Pagination
	if err := c.ShouldBindQuery(&page); err != nil {
		page.Page = utils.DefaultPage
		page.PageSize = utils.DefaultPageSize
	}
	page.Normalize()

	filter := dao.ItemFilter{}
	if catID := c.Query("category_id"); catID != "" {
		if id, err := strconv.ParseUint(catID, 10, 32); err == nil {
			v := uint32(id)
			filter.CategoryID = &v
		}
	}
	filter.Keyword = c.Query("keyword")
	if minPrice := c.Query("min_price"); minPrice != "" {
		if v, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = &v
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if v, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = &v
		}
	}

	result, err := h.itemService.List(c.Request.Context(), filter, page)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, result)
}
