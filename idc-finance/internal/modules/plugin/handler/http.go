package handler

import (
	"net/http"
	"strconv"
	"strings"

	"idc-finance/internal/modules/plugin/dto"
	pluginService "idc-finance/internal/modules/plugin/service"
	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *pluginService.Service
}

func New(service *pluginService.Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterAdminRoutes(router *gin.RouterGroup) {
	router.GET("/plugins/definitions", handler.listDefinitions)
	router.GET("/plugins/configs", handler.listConfigs)
	router.GET("/plugins/:code/config", handler.getConfig)
	router.PUT("/plugins/:code/config", handler.saveConfig)
}

func (handler *Handler) RegisterPortalRoutes(router *gin.RouterGroup) {
	router.POST("/plugins/kyc/:code/verify", handler.verifyKYC)
	router.GET("/plugins/kyc/records", handler.listKYCRecords)
	router.POST("/plugins/sms/:code/send", handler.dispatchSMS)
	router.GET("/plugins/sms/records", handler.listSMSRecords)
	router.POST("/plugins/payment/:code/orders", handler.createPaymentOrder)
	router.GET("/plugins/payment/orders", handler.listPaymentOrders)
}

func (handler *Handler) RegisterPublicRoutes(router *gin.RouterGroup) {
	router.POST("/plugins/payment/:code/notify", handler.notifyPayment)
	router.GET("/plugins/payment/:code/notify", handler.notifyPayment)
	router.POST("/plugins/sms/:code/send-public", handler.dispatchSMSPublic)
}

func (handler *Handler) listDefinitions(c *gin.Context) {
	items := handler.service.ListDefinitions()
	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listConfigs(c *gin.Context) {
	items := handler.service.ListConfigs()
	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) getConfig(c *gin.Context) {
	item, ok := handler.service.GetConfig(c.Param("code"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "PLUGIN_NOT_FOUND",
			"message":   "插件不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(item, middleware.GetRequestID(c)))
}

func (handler *Handler) saveConfig(c *gin.Context) {
	var request dto.SavePluginConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "插件配置参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	item, err := handler.service.SaveConfig(
		c.Param("code"),
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "PLUGIN_CONFIG_SAVE_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(item, middleware.GetRequestID(c)))
}

func (handler *Handler) verifyKYC(c *gin.Context) {
	var request dto.KYCVerifyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "实名认证参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	record, err := handler.service.VerifyKYC(c.Param("code"), getPortalCustomerID(c), request, middleware.GetRequestID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "KYC_VERIFY_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(record, middleware.GetRequestID(c)))
}

func (handler *Handler) listKYCRecords(c *gin.Context) {
	limit := parsePositiveInt(c.DefaultQuery("limit", "20"), 20)
	items := handler.service.ListKYCRecords(getPortalCustomerID(c), limit)
	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) dispatchSMS(c *gin.Context) {
	var request dto.SMSDispatchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "短信发送参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	record, err := handler.service.DispatchSMS(c.Param("code"), getPortalCustomerID(c), request, middleware.GetRequestID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "SMS_SEND_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(record, middleware.GetRequestID(c)))
}

func (handler *Handler) listSMSRecords(c *gin.Context) {
	limit := parsePositiveInt(c.DefaultQuery("limit", "20"), 20)
	items := handler.service.ListSMSRecords(getPortalCustomerID(c), limit)
	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) createPaymentOrder(c *gin.Context) {
	var request dto.CreatePaymentOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "支付下单参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	order, err := handler.service.CreatePaymentOrder(c.Param("code"), getPortalCustomerID(c), request, middleware.GetRequestID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "PAYMENT_ORDER_CREATE_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusCreated, appErrors.Ok(order, middleware.GetRequestID(c)))
}

func (handler *Handler) listPaymentOrders(c *gin.Context) {
	limit := parsePositiveInt(c.DefaultQuery("limit", "20"), 20)
	items := handler.service.ListPaymentOrders(getPortalCustomerID(c), limit)
	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) notifyPayment(c *gin.Context) {
	var request dto.PaymentNotifyRequest

	if outTradeNo := c.Query("out_trade_no"); outTradeNo != "" {
		request.OrderNo = outTradeNo
		request.TradeNo = c.DefaultQuery("trade_no", "")
		request.Status = "PAID"
		request.Sign = c.DefaultQuery("sign", "")
	} else if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "支付回调参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	_, ok, err := handler.service.NotifyPayment(c.Param("code"), request, middleware.GetRequestID(c))
	if err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}
	if !ok {
		c.String(http.StatusNotFound, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

func (handler *Handler) dispatchSMSPublic(c *gin.Context) {
	var request dto.SMSDispatchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "短信发送参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	record, err := handler.service.DispatchSMS(c.Param("code"), 0, request, middleware.GetRequestID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "SMS_SEND_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(record, middleware.GetRequestID(c)))
}

func getPortalCustomerID(c *gin.Context) int64 {
	value, exists := c.Get("customerID")
	if !exists {
		return 1
	}
	customerID, _ := value.(int64)
	return customerID
}

func getAdminID(c *gin.Context) int64 {
	value, exists := c.Get("adminID")
	if !exists {
		return 1
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

func parsePositiveInt(raw string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
