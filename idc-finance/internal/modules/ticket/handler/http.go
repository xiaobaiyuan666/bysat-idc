package handler

import (
	"net/http"
	"strconv"
	"strings"

	"idc-finance/internal/modules/ticket/domain"
	"idc-finance/internal/modules/ticket/dto"
	ticketService "idc-finance/internal/modules/ticket/service"
	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *ticketService.Service
}

func New(service *ticketService.Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterAdminRoutes(router *gin.RouterGroup) {
	router.GET("/tickets", handler.listAdminTickets)
	router.POST("/tickets", handler.adminCreateTicket)
	router.GET("/tickets/departments", handler.listTicketDepartments)
	router.PUT("/tickets/departments", handler.saveTicketDepartments)
	router.GET("/tickets/presets", handler.adminTicketPresets)
	router.PUT("/tickets/presets", handler.adminSaveTicketPresets)
	router.GET("/tickets/statistics", handler.adminTicketStatistics)
	router.POST("/tickets/auto-close/run", handler.adminRunAutoCloseSweep)
	router.POST("/tickets/:id/claim", handler.adminClaimTicket)
	router.GET("/tickets/:id", handler.adminTicketDetail)
	router.PATCH("/tickets/:id", handler.adminUpdateTicket)
	router.POST("/tickets/:id/replies", handler.adminReplyTicket)
}

func (handler *Handler) RegisterPortalRoutes(router *gin.RouterGroup) {
	router.GET("/tickets", handler.listPortalTickets)
	router.GET("/tickets/departments", handler.listPortalTicketDepartments)
	router.GET("/tickets/:id", handler.portalTicketDetail)
	router.POST("/tickets", handler.portalCreateTicket)
	router.POST("/tickets/:id/replies", handler.portalReplyTicket)
	router.POST("/tickets/:id/close", handler.portalCloseTicket)
}

func (handler *Handler) listAdminTickets(c *gin.Context) {
	items, total := handler.service.List(parseTicketFilter(c))
	c.JSON(http.StatusOK, appErrors.Ok(dto.ListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) adminCreateTicket(c *gin.Context) {
	var request dto.CreateAdminTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid ticket payload")
		return
	}

	result, err := handler.service.CreateAdminTicket(
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "TICKET_CREATE_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusCreated, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminTicketPresets(c *gin.Context) {
	items := handler.service.ListPresetReplies(strings.TrimSpace(c.Query("departmentName")))
	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listTicketDepartments(c *gin.Context) {
	items := handler.service.ListDepartments()
	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) saveTicketDepartments(c *gin.Context) {
	var request dto.UpdateDepartmentsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid department payload")
		return
	}

	items, err := handler.service.SaveDepartments(request.Items)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TICKET_DEPARTMENTS_SAVE_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) adminSaveTicketPresets(c *gin.Context) {
	var request dto.UpdatePresetRepliesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid preset payload")
		return
	}

	items, err := handler.service.SavePresetReplies(request.Items)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TICKET_PRESETS_SAVE_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) adminRunAutoCloseSweep(c *gin.Context) {
	result, err := handler.service.RunAutoCloseSweep(
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TICKET_AUTO_CLOSE_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminTicketStatistics(c *gin.Context) {
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.GetStatistics(), middleware.GetRequestID(c)))
}

func (handler *Handler) adminClaimTicket(c *gin.Context) {
	ticketID, ok := parseID(c)
	if !ok {
		return
	}

	result, exists, err := handler.service.ClaimTicket(
		ticketID,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "TICKET_CLAIM_FAILED", err.Error())
		return
	}
	if !exists {
		respondError(c, http.StatusNotFound, "TICKET_NOT_FOUND", "Ticket not found")
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminTicketDetail(c *gin.Context) {
	ticketID, ok := parseID(c)
	if !ok {
		return
	}
	result, exists := handler.service.GetDetail(ticketID)
	if !exists {
		respondError(c, http.StatusNotFound, "TICKET_NOT_FOUND", "Ticket not found")
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminUpdateTicket(c *gin.Context) {
	ticketID, ok := parseID(c)
	if !ok {
		return
	}

	var request dto.UpdateTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid ticket update payload")
		return
	}

	result, exists, err := handler.service.UpdateAdminTicket(
		ticketID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "TICKET_UPDATE_FAILED", err.Error())
		return
	}
	if !exists {
		respondError(c, http.StatusNotFound, "TICKET_NOT_FOUND", "Ticket not found")
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminReplyTicket(c *gin.Context) {
	ticketID, ok := parseID(c)
	if !ok {
		return
	}

	var request dto.ReplyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid reply payload")
		return
	}

	result, exists, err := handler.service.ReplyAsAdmin(
		ticketID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "TICKET_REPLY_FAILED", err.Error())
		return
	}
	if !exists {
		respondError(c, http.StatusNotFound, "TICKET_NOT_FOUND", "Ticket not found")
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) listPortalTickets(c *gin.Context) {
	items, total, err := handler.service.ListByCustomer(getPortalCustomerID(c), parseTicketFilter(c))
	if err != nil {
		respondError(c, http.StatusNotFound, "CUSTOMER_NOT_FOUND", err.Error())
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(dto.ListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listPortalTicketDepartments(c *gin.Context) {
	items := make([]domain.Department, 0)
	for _, item := range handler.service.ListDepartments() {
		if !item.Enabled {
			continue
		}
		items = append(items, item)
	}
	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"items": items,
		"total": len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) portalTicketDetail(c *gin.Context) {
	ticketID, ok := parseID(c)
	if !ok {
		return
	}
	result, exists := handler.service.GetDetailByCustomer(getPortalCustomerID(c), ticketID)
	if !exists {
		respondError(c, http.StatusNotFound, "TICKET_NOT_FOUND", "Ticket not found")
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) portalCreateTicket(c *gin.Context) {
	var request dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid ticket payload")
		return
	}

	result, err := handler.service.CreatePortalTicket(
		getPortalCustomerID(c),
		getPortalCustomerName(c),
		request,
		middleware.GetRequestID(c),
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "TICKET_CREATE_FAILED", err.Error())
		return
	}
	c.JSON(http.StatusCreated, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) portalReplyTicket(c *gin.Context) {
	ticketID, ok := parseID(c)
	if !ok {
		return
	}

	var request dto.ReplyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid reply payload")
		return
	}

	result, exists, err := handler.service.ReplyAsCustomer(
		getPortalCustomerID(c),
		getPortalCustomerName(c),
		ticketID,
		request,
		middleware.GetRequestID(c),
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "TICKET_REPLY_FAILED", err.Error())
		return
	}
	if !exists {
		respondError(c, http.StatusNotFound, "TICKET_NOT_FOUND", "Ticket not found")
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) portalCloseTicket(c *gin.Context) {
	ticketID, ok := parseID(c)
	if !ok {
		return
	}

	result, exists, err := handler.service.CloseByCustomer(
		getPortalCustomerID(c),
		getPortalCustomerName(c),
		ticketID,
		middleware.GetRequestID(c),
	)
	if err != nil {
		respondError(c, http.StatusBadRequest, "TICKET_CLOSE_FAILED", err.Error())
		return
	}
	if !exists {
		respondError(c, http.StatusNotFound, "TICKET_NOT_FOUND", "Ticket not found")
		return
	}
	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func parseTicketFilter(c *gin.Context) domain.ListFilter {
	return domain.ListFilter{
		Page:       parsePositiveInt(c.DefaultQuery("page", "1"), 1),
		Limit:      parsePositiveInt(c.DefaultQuery("limit", "20"), 20),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
		Status:     strings.TrimSpace(c.Query("status")),
		Priority:   strings.TrimSpace(c.Query("priority")),
		CustomerID: parseOptionalInt64(c.Query("customerId")),
		ServiceID:  parseOptionalInt64(c.Query("serviceId")),
		Department: strings.TrimSpace(c.Query("departmentName")),
		AdminID:    parseOptionalInt64(c.Query("assignedAdminId")),
	}
}

func parseID(c *gin.Context) (int64, bool) {
	value, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || value <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid id")
		return 0, false
	}
	return value, true
}

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"code":      code,
		"message":   message,
		"requestId": middleware.GetRequestID(c),
	})
}

func getPortalCustomerID(c *gin.Context) int64 {
	value, exists := c.Get("customerID")
	if !exists {
		return 1
	}
	customerID, _ := value.(int64)
	return customerID
}

func getPortalCustomerName(c *gin.Context) string {
	value, exists := c.Get("customerName")
	if !exists {
		return "Demo Customer"
	}
	customerName, _ := value.(string)
	return customerName
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
		return "System Admin"
	}
	adminName, _ := value.(string)
	return adminName
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func parseOptionalInt64(value string) int64 {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || parsed <= 0 {
		return 0
	}
	return parsed
}
