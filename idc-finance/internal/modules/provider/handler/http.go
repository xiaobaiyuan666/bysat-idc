package handler

import (
	"net/http"
	"strconv"

	providerService "idc-finance/internal/modules/provider/service"
	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *providerService.Service
}

func New(service *providerService.Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterAdminRoutes(router *gin.RouterGroup) {
	router.GET("/providers/accounts", handler.listAccounts)
	router.GET("/providers/accounts/:id", handler.getAccount)
	router.POST("/providers/accounts", handler.createAccount)
	router.PATCH("/providers/accounts/:id", handler.updateAccount)
	router.DELETE("/providers/accounts/:id", handler.deleteAccount)
	router.GET("/providers/accounts/:id/health", handler.accountHealth)
	router.GET("/providers/mofang/health", handler.health)
	router.GET("/providers/mofang/instances", handler.listInstances)
	router.GET("/providers/mofang/instances/:id", handler.getInstanceDetail)
	router.POST("/providers/mofang/instances/:id/actions/:action", handler.instanceAction)
	router.POST("/providers/mofang/sync", handler.pullSync)
	router.POST("/providers/mofang/services/:id/sync", handler.syncService)
	router.GET("/providers/mofang/services/:id/resources", handler.getServiceResources)
	router.POST("/providers/mofang/services/:id/resource-actions/:action", handler.runServiceResourceAction)
	router.GET("/providers/mofang/sync-logs", handler.listSyncLogs)
	router.GET("/providers/upstream/health", handler.upstreamHealth)
	router.GET("/providers/upstream/products", handler.listUpstreamProducts)
	router.GET("/providers/upstream/products/:id/template", handler.getUpstreamProductTemplate)
}

func (handler *Handler) health(c *gin.Context) {
	accountID, err := parseOptionalIDParam(c.Query("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "accountId 参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.CheckHealth(accountID), middleware.GetRequestID(c)))
}

func (handler *Handler) listInstances(c *gin.Context) {
	accountID, queryErr := parseOptionalIDParam(c.Query("accountId"))
	if queryErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "accountId 参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	items, err := handler.service.ListInstances(accountID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "MOFANG_LIST_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) getInstanceDetail(c *gin.Context) {
	accountID, queryErr := parseOptionalIDParam(c.Query("accountId"))
	if queryErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "accountId 参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	item, err := handler.service.GetInstanceDetail(c.Param("id"), accountID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "MOFANG_DETAIL_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(item, middleware.GetRequestID(c)))
}

func (handler *Handler) instanceAction(c *gin.Context) {
	var request providerService.ActionRequest
	_ = c.ShouldBindJSON(&request)

	accountID, queryErr := parseOptionalIDParam(c.Query("accountId"))
	if queryErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "accountId 参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, err := handler.service.ExecuteAction(c.Param("id"), c.Param("action"), request, accountID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "MOFANG_ACTION_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) pullSync(c *gin.Context) {
	var request providerService.PullSyncOptions
	_ = c.ShouldBindJSON(&request)

	result, err := handler.service.PullSync(request)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "MOFANG_SYNC_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) syncService(c *gin.Context) {
	serviceID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "服务编号不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, err := handler.service.SyncServiceByID(serviceID, true)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "MOFANG_SYNC_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) getServiceResources(c *gin.Context) {
	serviceID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "服务编号不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, err := handler.service.GetServiceResources(serviceID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "MOFANG_RESOURCES_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) runServiceResourceAction(c *gin.Context) {
	serviceID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "服务编号不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request providerService.ResourceActionRequest
	_ = c.ShouldBindJSON(&request)

	result, err := handler.service.ExecuteServiceResourceAction(serviceID, c.Param("action"), request)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "MOFANG_RESOURCE_ACTION_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) listSyncLogs(c *gin.Context) {
	limit, _ := parseOptionalIDParam(c.Query("limit"))
	serviceID, _ := parseOptionalIDParam(c.Query("serviceId"))

	result, err := handler.service.ListSyncLogs(int(limit), serviceID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "MOFANG_SYNC_LOG_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": result,
		"total": len(result),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) upstreamHealth(c *gin.Context) {
	accountID, err := parseOptionalIDParam(c.Query("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "accountId 参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.CheckFinanceHealth(accountID), middleware.GetRequestID(c)))
}

func (handler *Handler) listUpstreamProducts(c *gin.Context) {
	accountID, queryErr := parseOptionalIDParam(c.Query("accountId"))
	if queryErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "accountId 参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	items, err := handler.service.ListUpstreamProducts(accountID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "UPSTREAM_LIST_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	total := 0
	for _, group := range items {
		total += len(group.Products)
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"groups": items,
		"total":  total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) getUpstreamProductTemplate(c *gin.Context) {
	accountID, queryErr := parseOptionalIDParam(c.Query("accountId"))
	if queryErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "accountId 参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	item, err := handler.service.GetUpstreamProductTemplate(c.Param("id"), accountID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":      "UPSTREAM_TEMPLATE_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(item, middleware.GetRequestID(c)))
}

func (handler *Handler) listAccounts(c *gin.Context) {
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.ListAccounts(c.Query("providerType")), middleware.GetRequestID(c)))
}

func (handler *Handler) getAccount(c *gin.Context) {
	accountID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "接口编号不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	account, exists := handler.service.GetAccount(accountID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "NOT_FOUND",
			"message":   "接口不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(account, middleware.GetRequestID(c)))
}

func (handler *Handler) createAccount(c *gin.Context) {
	var request providerService.AccountUpsertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "接口参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	account, err := handler.service.CreateAccount(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "PROVIDER_ACCOUNT_CREATE_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(account, middleware.GetRequestID(c)))
}

func (handler *Handler) updateAccount(c *gin.Context) {
	accountID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "接口编号不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	var request providerService.AccountUpsertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "接口参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	account, updateErr := handler.service.UpdateAccount(accountID, request)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "PROVIDER_ACCOUNT_UPDATE_FAILED",
			"message":   updateErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(account, middleware.GetRequestID(c)))
}

func (handler *Handler) deleteAccount(c *gin.Context) {
	accountID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "接口编号不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if err := handler.service.DeleteAccount(accountID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "PROVIDER_ACCOUNT_DELETE_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(gin.H{"deleted": true}, middleware.GetRequestID(c)))
}

func (handler *Handler) accountHealth(c *gin.Context) {
	accountID, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "接口编号不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	account, exists := handler.service.GetAccount(accountID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "NOT_FOUND",
			"message":   "接口不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if account.ProviderType == "MOFANG_CLOUD" {
		c.JSON(http.StatusOK, appErrors.Ok(handler.service.CheckHealth(accountID), middleware.GetRequestID(c)))
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.CheckFinanceHealth(accountID), middleware.GetRequestID(c)))
}

func parseIDParam(value string) (int64, error) {
	return parseOptionalIDParam(value)
}

func parseOptionalIDParam(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}
