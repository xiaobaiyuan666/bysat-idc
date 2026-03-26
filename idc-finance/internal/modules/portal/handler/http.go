package handler

import (
	"net/http"

	portalService "idc-finance/internal/modules/portal/service"
	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *portalService.Service
}

func New(service *portalService.Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard", handler.dashboard)
	router.GET("/account", handler.account)
	router.GET("/wallet", handler.wallet)
}

func (handler *Handler) dashboard(c *gin.Context) {
	result, ok := handler.service.GetDashboard(getPortalCustomerID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "客户不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) tickets(c *gin.Context) {
	result, ok := handler.service.ListTickets(getPortalCustomerID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "客户不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": result,
		"total": len(result),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) account(c *gin.Context) {
	result, ok := handler.service.GetAccount(getPortalCustomerID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "客户不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) wallet(c *gin.Context) {
	result, ok := handler.service.GetWallet(getPortalCustomerID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "瀹㈡埛涓嶅瓨鍦?",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func getPortalCustomerID(c *gin.Context) int64 {
	value, exists := c.Get("customerID")
	if !exists {
		return 1
	}
	customerID, _ := value.(int64)
	return customerID
}
