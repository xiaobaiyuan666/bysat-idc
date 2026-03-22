package handler

import (
	"errors"
	"net/http"
	"strconv"

	"idc-finance/internal/modules/customer/dto"
	"idc-finance/internal/modules/customer/repository"
	"idc-finance/internal/modules/customer/service"
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

func (handler *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/customers", handler.list)
	router.POST("/customers", handler.create)
	router.PATCH("/customers/:id", handler.update)
	router.GET("/customers/:id", handler.detail)
	router.GET("/customers/:id/contacts", handler.listContacts)
	router.POST("/customers/:id/contacts", handler.addContact)
	router.PUT("/customers/:id/contacts/:contactId", handler.updateContact)
	router.PATCH("/customers/:id/contacts/:contactId", handler.updateContact)
	router.DELETE("/customers/:id/contacts/:contactId", handler.deleteContact)
	router.GET("/customers/:id/identities", handler.listIdentities)
	router.POST("/customers/:id/identities/:identityId/review", handler.reviewIdentity)
	router.POST("/customers/:id/identity/review", handler.reviewIdentity)
	router.GET("/customers/:id/services", handler.listServices)
	router.GET("/customers/:id/invoices", handler.listInvoices)
	router.GET("/customers/:id/tickets", handler.listTickets)
	router.GET("/customers/:id/audit-logs", handler.listAuditLogs)

	router.GET("/customer-groups", handler.listGroups)
	router.POST("/customer-groups", handler.createGroup)
	router.PATCH("/customer-groups/:id", handler.updateGroup)
	router.DELETE("/customer-groups/:id", handler.deleteGroup)

	router.GET("/customer-levels", handler.listLevels)
	router.POST("/customer-levels", handler.createLevel)
	router.PATCH("/customer-levels/:id", handler.updateLevel)
	router.DELETE("/customer-levels/:id", handler.deleteLevel)

	router.GET("/customer-identities", handler.listIdentityOverview)
}

func (handler *Handler) list(c *gin.Context) {
	items := handler.service.List()
	c.JSON(http.StatusOK, appErrors.Ok(dto.CustomerListResponse{
		Items: items,
		Total: len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) create(c *gin.Context) {
	var request dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "客户参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	customer := handler.service.CreateWithAudit(request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	c.JSON(http.StatusCreated, appErrors.Ok(customer, middleware.GetRequestID(c)))
}

func (handler *Handler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": "客户编号不正确", "requestId": middleware.GetRequestID(c)})
		return
	}

	var request dto.UpdateCustomerRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": "客户参数不正确", "requestId": middleware.GetRequestID(c)})
		return
	}

	customer, exists := handler.service.Update(id, request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": "CUSTOMER_NOT_FOUND", "message": "客户不存在", "requestId": middleware.GetRequestID(c)})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(customer, middleware.GetRequestID(c)))
}

func (handler *Handler) detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "客户编号不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	customer, exists := handler.service.GetByID(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "CUSTOMER_NOT_FOUND",
			"message":   "客户不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(customer, middleware.GetRequestID(c)))
}

func (handler *Handler) listGroups(c *gin.Context) {
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.ListGroups(), middleware.GetRequestID(c)))
}

func (handler *Handler) createGroup(c *gin.Context) {
	var request dto.SaveCustomerGroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": "客户分组参数不正确", "requestId": middleware.GetRequestID(c)})
		return
	}

	group, err := handler.service.CreateGroup(request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "GROUP_CREATE_FAILED", "message": "客户分组创建失败", "requestId": middleware.GetRequestID(c)})
		return
	}

	c.JSON(http.StatusCreated, appErrors.Ok(group, middleware.GetRequestID(c)))
}

func (handler *Handler) updateGroup(c *gin.Context) {
	id, ok := parseID(c, "id", "客户分组编号不正确")
	if !ok {
		return
	}

	var request dto.SaveCustomerGroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": "客户分组参数不正确", "requestId": middleware.GetRequestID(c)})
		return
	}

	group, err := handler.service.UpdateGroup(id, request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if err != nil {
		handler.handleGroupError(c, err)
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(group, middleware.GetRequestID(c)))
}

func (handler *Handler) deleteGroup(c *gin.Context) {
	id, ok := parseID(c, "id", "客户分组编号不正确")
	if !ok {
		return
	}

	if err := handler.service.DeleteGroup(id, getAdminID(c), getAdminName(c), middleware.GetRequestID(c)); err != nil {
		handler.handleGroupError(c, err)
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{"deleted": true}, middleware.GetRequestID(c)))
}

func (handler *Handler) listLevels(c *gin.Context) {
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.ListLevels(), middleware.GetRequestID(c)))
}

func (handler *Handler) createLevel(c *gin.Context) {
	var request dto.SaveCustomerLevelRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": "客户等级参数不正确", "requestId": middleware.GetRequestID(c)})
		return
	}

	level, err := handler.service.CreateLevel(request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "LEVEL_CREATE_FAILED", "message": "客户等级创建失败", "requestId": middleware.GetRequestID(c)})
		return
	}

	c.JSON(http.StatusCreated, appErrors.Ok(level, middleware.GetRequestID(c)))
}

func (handler *Handler) updateLevel(c *gin.Context) {
	id, ok := parseID(c, "id", "客户等级编号不正确")
	if !ok {
		return
	}

	var request dto.SaveCustomerLevelRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": "客户等级参数不正确", "requestId": middleware.GetRequestID(c)})
		return
	}

	level, err := handler.service.UpdateLevel(id, request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if err != nil {
		handler.handleLevelError(c, err)
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(level, middleware.GetRequestID(c)))
}

func (handler *Handler) deleteLevel(c *gin.Context) {
	id, ok := parseID(c, "id", "客户等级编号不正确")
	if !ok {
		return
	}

	if err := handler.service.DeleteLevel(id, getAdminID(c), getAdminName(c), middleware.GetRequestID(c)); err != nil {
		handler.handleLevelError(c, err)
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{"deleted": true}, middleware.GetRequestID(c)))
}

func (handler *Handler) listIdentityOverview(c *gin.Context) {
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.ListIdentityOverview(), middleware.GetRequestID(c)))
}

func (handler *Handler) listContacts(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}
	contacts, exists := handler.service.ListContacts(customerID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": "CUSTOMER_NOT_FOUND", "message": "客户不存在", "requestId": middleware.GetRequestID(c)})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(contacts, middleware.GetRequestID(c)))
}

func (handler *Handler) addContact(c *gin.Context) {
	customerID, ok := parseID(c, "id", "客户编号不正确")
	if !ok {
		return
	}

	var request dto.CreateContactRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": "联系人参数不正确", "requestId": middleware.GetRequestID(c)})
		return
	}

	contact, result := handler.service.AddContact(customerID, request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if !result {
		c.JSON(http.StatusNotFound, gin.H{"code": "CUSTOMER_NOT_FOUND", "message": "客户不存在", "requestId": middleware.GetRequestID(c)})
		return
	}

	c.JSON(http.StatusCreated, appErrors.Ok(contact, middleware.GetRequestID(c)))
}

func (handler *Handler) updateContact(c *gin.Context) {
	customerID, ok := parseID(c, "id", "客户编号不正确")
	if !ok {
		return
	}
	contactID, ok := parseID(c, "contactId", "联系人编号不正确")
	if !ok {
		return
	}

	var request dto.UpdateContactRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": "联系人参数不正确", "requestId": middleware.GetRequestID(c)})
		return
	}

	contact, result := handler.service.UpdateContact(customerID, contactID, request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if !result {
		c.JSON(http.StatusNotFound, gin.H{"code": "CONTACT_NOT_FOUND", "message": "联系人不存在", "requestId": middleware.GetRequestID(c)})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(contact, middleware.GetRequestID(c)))
}

func (handler *Handler) deleteContact(c *gin.Context) {
	customerID, ok := parseID(c, "id", "客户编号不正确")
	if !ok {
		return
	}
	contactID, ok := parseID(c, "contactId", "联系人编号不正确")
	if !ok {
		return
	}

	ok = handler.service.DeleteContact(customerID, contactID, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": "CONTACT_NOT_FOUND", "message": "联系人不存在", "requestId": middleware.GetRequestID(c)})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{"deleted": true}, middleware.GetRequestID(c)))
}

func (handler *Handler) reviewIdentity(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}

	var request dto.ReviewIdentityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": "实名认证审核参数不正确", "requestId": middleware.GetRequestID(c)})
		return
	}

	identity, exists := handler.service.ReviewIdentity(customerID, request, getAdminID(c), getAdminName(c), middleware.GetRequestID(c))
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": "IDENTITY_NOT_FOUND", "message": "实名认证信息不存在", "requestId": middleware.GetRequestID(c)})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(identity, middleware.GetRequestID(c)))
}

func (handler *Handler) listIdentities(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}
	identities, exists := handler.service.ListIdentityDetails(customerID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": "CUSTOMER_NOT_FOUND", "message": "客户不存在", "requestId": middleware.GetRequestID(c)})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(identities, middleware.GetRequestID(c)))
}

func (handler *Handler) listServices(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}
	items, exists := handler.service.ListServices(customerID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": "CUSTOMER_NOT_FOUND", "message": "客户不存在", "requestId": middleware.GetRequestID(c)})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(items, middleware.GetRequestID(c)))
}

func (handler *Handler) listInvoices(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}
	items, exists := handler.service.ListInvoices(customerID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": "CUSTOMER_NOT_FOUND", "message": "客户不存在", "requestId": middleware.GetRequestID(c)})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(items, middleware.GetRequestID(c)))
}

func (handler *Handler) listTickets(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}
	items, exists := handler.service.ListTickets(customerID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": "CUSTOMER_NOT_FOUND", "message": "客户不存在", "requestId": middleware.GetRequestID(c)})
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(items, middleware.GetRequestID(c)))
}

func (handler *Handler) listAuditLogs(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.ListCustomerAuditLogs(customerID), middleware.GetRequestID(c)))
}

func (handler *Handler) handleGroupError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrGroupNotFound):
		c.JSON(http.StatusNotFound, gin.H{"code": "GROUP_NOT_FOUND", "message": "客户分组不存在", "requestId": middleware.GetRequestID(c)})
	case errors.Is(err, repository.ErrGroupInUse):
		c.JSON(http.StatusConflict, gin.H{"code": "GROUP_IN_USE", "message": "该客户分组仍被客户使用，不能删除", "requestId": middleware.GetRequestID(c)})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": "GROUP_OPERATION_FAILED", "message": "客户分组操作失败", "requestId": middleware.GetRequestID(c)})
	}
}

func (handler *Handler) handleLevelError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrLevelNotFound):
		c.JSON(http.StatusNotFound, gin.H{"code": "LEVEL_NOT_FOUND", "message": "客户等级不存在", "requestId": middleware.GetRequestID(c)})
	case errors.Is(err, repository.ErrLevelInUse):
		c.JSON(http.StatusConflict, gin.H{"code": "LEVEL_IN_USE", "message": "该客户等级仍被客户使用，不能删除", "requestId": middleware.GetRequestID(c)})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": "LEVEL_OPERATION_FAILED", "message": "客户等级操作失败", "requestId": middleware.GetRequestID(c)})
	}
}

func parseCustomerID(c *gin.Context) (int64, bool) {
	return parseID(c, "id", "客户编号不正确")
}

func parseID(c *gin.Context, param, message string) (int64, bool) {
	value, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ARGUMENT", "message": message, "requestId": middleware.GetRequestID(c)})
		return 0, false
	}
	return value, true
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
