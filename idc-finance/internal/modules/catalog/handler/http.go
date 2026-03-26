package handler

import (
	"net/http"
	"strconv"

	"idc-finance/internal/modules/catalog/dto"
	"idc-finance/internal/modules/catalog/service"
	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterAdminRoutes(router *gin.RouterGroup) {
	router.GET("/products", handler.listAdmin)
	router.GET("/products/:id", handler.detail)
	router.POST("/products", handler.create)
	router.POST("/products/import-upstream", handler.importUpstream)
	router.GET("/products/import-upstream/history", handler.listImportUpstreamHistory)
	router.GET("/products/import-upstream/history/:id", handler.detailImportUpstreamHistory)
	router.PATCH("/products/:id", handler.update)
	router.POST("/products/:id/upstream/sync", handler.syncUpstream)
	router.GET("/product-groups", handler.groups)
}

func (handler *Handler) RegisterPortalRoutes(router *gin.RouterGroup) {
	router.GET("/products", handler.listPortal)
}

func (handler *Handler) listAdmin(c *gin.Context) {
	items := handler.service.List()
	c.JSON(http.StatusOK, appErrors.Ok(dto.ProductListResponse{
		Items: items,
		Total: len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listPortal(c *gin.Context) {
	items := handler.service.ListActive()
	c.JSON(http.StatusOK, appErrors.Ok(dto.ProductListResponse{
		Items: items,
		Total: len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "商品编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok := handler.service.GetDetail(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "PRODUCT_NOT_FOUND",
			"message":   "商品不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) create(c *gin.Context) {
	var request dto.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "商品参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	product := handler.service.Create(request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	c.JSON(http.StatusCreated, appErrors.Ok(product, middleware.GetRequestID(c)))
}

func (handler *Handler) importUpstream(c *gin.Context) {
	var request dto.ImportUpstreamProductsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "上游导入参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, err := handler.service.ImportUpstreamProducts(request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "PRODUCT_IMPORT_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) listImportUpstreamHistory(c *gin.Context) {
	var query dto.UpstreamImportHistoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "上游同步记录查询参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(handler.service.ListUpstreamImportHistory(query), middleware.GetRequestID(c)))
}

func (handler *Handler) detailImportUpstreamHistory(c *gin.Context) {
	result, ok := handler.service.GetUpstreamImportHistory(c.Param("id"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "UPSTREAM_IMPORT_HISTORY_NOT_FOUND",
			"message":   "未找到对应的上游同步记录",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "商品编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "商品参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	product, ok := handler.service.Update(id, request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "PRODUCT_NOT_FOUND",
			"message":   "商品不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(product, middleware.GetRequestID(c)))
}

func (handler *Handler) syncUpstream(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "商品编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.SyncProductUpstreamRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "商品同步参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	product, ok, syncErr := handler.service.SyncUpstream(id, request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if syncErr != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "PRODUCT_SYNC_FAILED",
			"message":   syncErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "PRODUCT_NOT_FOUND",
			"message":   "商品不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(product, middleware.GetRequestID(c)))
}

func (handler *Handler) groups(c *gin.Context) {
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.ListGroups(), middleware.GetRequestID(c)))
}

func getAdminID(c *gin.Context) int64 {
	value, exists := c.Get("adminID")
	if !exists {
		return 0
	}
	adminID, _ := value.(int64)
	return adminID
}

func getAdminName(c *gin.Context) string {
	value, exists := c.Get("adminName")
	if !exists {
		return "系统管理员"
	}
	adminName, _ := value.(string)
	return adminName
}
