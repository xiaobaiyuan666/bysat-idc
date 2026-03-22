package handler

import (
	"net/http"
	"strconv"

	automationDomain "idc-finance/internal/modules/automation/domain"
	automationService "idc-finance/internal/modules/automation/service"
	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *automationService.Service
	retryer automationService.RetryExecutor
}

func New(service *automationService.Service, retryer automationService.RetryExecutor) *Handler {
	return &Handler{service: service, retryer: retryer}
}

func (handler *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/automation/tasks", handler.listTasks)
	router.GET("/automation/tasks/:id", handler.getTask)
	router.POST("/automation/tasks/:id/retry", handler.retryTask)
	router.GET("/automation/settings", handler.getSettings)
	router.PATCH("/automation/settings", handler.updateSettings)
}

func (handler *Handler) listTasks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	sourceID, _ := parseOptionalInt64(c.Query("sourceId"))
	orderID, _ := parseOptionalInt64(c.Query("orderId"))
	invoiceID, _ := parseOptionalInt64(c.Query("invoiceId"))
	serviceID, _ := parseOptionalInt64(c.Query("serviceId"))

	items, summary := handler.service.List(automationDomain.TaskFilter{
		Status:     c.Query("status"),
		TaskType:   c.Query("taskType"),
		Channel:    c.Query("channel"),
		Stage:      c.Query("stage"),
		SourceType: c.Query("sourceType"),
		SourceID:   sourceID,
		OrderID:    orderID,
		InvoiceID:  invoiceID,
		ServiceID:  serviceID,
		Keyword:    c.Query("keyword"),
		Limit:      limit,
	})

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items":   items,
		"summary": summary,
		"total":   len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) getTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "任务编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	task, ok := handler.service.GetByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "TASK_NOT_FOUND",
			"message":   "自动化任务不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(task, middleware.GetRequestID(c)))
}

func (handler *Handler) retryTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "任务编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, err := handler.service.Retry(id, handler.retryer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "AUTOMATION_RETRY_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) getSettings(c *gin.Context) {
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.GetSettings(), middleware.GetRequestID(c)))
}

func (handler *Handler) updateSettings(c *gin.Context) {
	var request automationDomain.AutomationSettings
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "自动化策略参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	settings, err := handler.service.SaveSettings(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      "AUTOMATION_SETTINGS_SAVE_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(settings, middleware.GetRequestID(c)))
}

func parseOptionalInt64(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}
