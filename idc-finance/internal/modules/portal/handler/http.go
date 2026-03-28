package handler

import (
	"net/http"
	"strconv"
	"strings"

	customerDTO "idc-finance/internal/modules/customer/dto"
	portalService "idc-finance/internal/modules/portal/service"
	"idc-finance/internal/platform/auth"
	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *portalService.Service
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

func New(service *portalService.Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard", handler.dashboard)
	router.GET("/account", handler.account)
	router.GET("/account/profile", handler.account)
	router.PUT("/account/profile", handler.updateProfile)
	router.GET("/account/contacts", handler.contacts)
	router.POST("/account/contacts", handler.createContact)
	router.PATCH("/account/contacts/:id", handler.updateContact)
	router.DELETE("/account/contacts/:id", handler.deleteContact)
	router.GET("/account/identity", handler.identity)
	router.POST("/account/identity", handler.submitIdentity)
	router.POST("/account/security/password", handler.changePassword)
	router.GET("/wallet", handler.wallet)
	router.GET("/wallet/transactions", handler.walletTransactions)
	router.GET("/wallet/recharges", handler.walletRecharges)
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

func (handler *Handler) updateProfile(c *gin.Context) {
	var request customerDTO.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "资料参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok := handler.service.UpdateProfile(getPortalCustomerID(c), request, middleware.GetRequestID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "客户不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	_ = auth.UpdatePortalIdentity(result.Customer.ID, result.Customer.Name, result.Customer.Name)
	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) contacts(c *gin.Context) {
	result, ok := handler.service.GetAccount(getPortalCustomerID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "客户不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": result.Customer.Contacts,
		"total": len(result.Customer.Contacts),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) createContact(c *gin.Context) {
	var request customerDTO.CreateContactRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "联系人参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok := handler.service.CreateContact(getPortalCustomerID(c), request, middleware.GetRequestID(c))
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

func (handler *Handler) updateContact(c *gin.Context) {
	contactID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "联系人编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request customerDTO.UpdateContactRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "联系人参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok := handler.service.UpdateContact(getPortalCustomerID(c), contactID, request, middleware.GetRequestID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CONTACT_NOT_FOUND",
			"message":   "联系人不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) deleteContact(c *gin.Context) {
	contactID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "联系人编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	if !handler.service.DeleteContact(getPortalCustomerID(c), contactID, middleware.GetRequestID(c)) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CONTACT_NOT_FOUND",
			"message":   "联系人不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{"deleted": true}, middleware.GetRequestID(c)))
}

func (handler *Handler) identity(c *gin.Context) {
	result, ok := handler.service.GetIdentity(getPortalCustomerID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "客户不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"identity": result,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) submitIdentity(c *gin.Context) {
	var request customerDTO.SubmitIdentityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "实名认证参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok := handler.service.SubmitIdentity(getPortalCustomerID(c), request, middleware.GetRequestID(c))
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

func (handler *Handler) changePassword(c *gin.Context) {
	var request changePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "密码参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	if err := auth.ChangePortalPassword(strings.TrimSpace(request.CurrentPassword), strings.TrimSpace(request.NewPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "PASSWORD_CHANGE_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"changed": true,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) wallet(c *gin.Context) {
	result, ok := handler.service.GetWallet(getPortalCustomerID(c))
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

func (handler *Handler) walletTransactions(c *gin.Context) {
	items, ok := handler.service.ListWalletTransactions(
		getPortalCustomerID(c),
		c.Query("transactionType"),
		c.Query("keyword"),
		c.Query("channel"),
		parsePositiveInt(c.DefaultQuery("limit", "200"), 200),
	)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "客户不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) walletRecharges(c *gin.Context) {
	items, ok := handler.service.ListWalletTransactions(
		getPortalCustomerID(c),
		"RECHARGE",
		c.Query("keyword"),
		c.Query("channel"),
		parsePositiveInt(c.DefaultQuery("limit", "200"), 200),
	)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "客户不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func getPortalCustomerID(c *gin.Context) int64 {
	value, exists := c.Get("customerID")
	if !exists {
		return 1
	}
	customerID, _ := value.(int64)
	return customerID
}

func parsePositiveInt(raw string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
