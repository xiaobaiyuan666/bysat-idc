package handler

import (
	"net/http"
	"strconv"
	"strings"

	"idc-finance/internal/modules/order/domain"
	"idc-finance/internal/modules/order/dto"
	"idc-finance/internal/modules/order/service"
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
	router.GET("/orders", handler.listAdminOrders)
	router.GET("/orders/:id", handler.adminOrderDetail)
	router.PATCH("/orders/:id", handler.adminUpdatePendingOrder)
	router.GET("/invoices", handler.listAdminInvoices)
	router.GET("/invoices/:id", handler.adminInvoiceDetail)
	router.PATCH("/invoices/:id", handler.adminUpdateUnpaidInvoice)
	router.POST("/invoices/:id/receive-payment", handler.adminReceivePayment)
	router.POST("/invoices/:id/refund", handler.adminRefundInvoice)
	router.GET("/services", handler.listAdminServices)
	router.GET("/services/:id", handler.adminServiceDetail)
	router.PATCH("/services/:id", handler.adminUpdateService)
	router.POST("/services/:id/change-orders", handler.adminCreateServiceChangeOrder)
	router.POST("/services/:id/actions/:action", handler.adminServiceAction)
}

func (handler *Handler) adminUpdatePendingOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "订单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.UpdatePendingOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "订单调整参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok, updateErr := handler.service.UpdatePendingOrder(
		orderID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "ORDER_UPDATE_FAILED",
			"message":   updateErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "ORDER_NOT_FOUND",
			"message":   "订单不存在，无法人工调整",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminCreateServiceChangeOrder(c *gin.Context) {
	serviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "服务编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.CreateServiceChangeOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "改配单参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok, createErr := handler.service.CreateServiceChangeOrder(
		serviceID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if createErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "CHANGE_ORDER_CREATE_FAILED",
			"message":   createErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "SERVICE_NOT_FOUND",
			"message":   "服务不存在，无法生成改配单",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) RegisterPortalRoutes(router *gin.RouterGroup) {
	router.GET("/orders", handler.listPortalOrders)
	router.GET("/invoices", handler.listPortalInvoices)
	router.GET("/services", handler.listPortalServices)
	router.POST("/orders/checkout", handler.checkout)
	router.POST("/invoices/:id/pay", handler.payInvoice)
}

func (handler *Handler) listAdminOrders(c *gin.Context) {
	filter, err := parseOrderListFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	items, total := handler.service.ListOrders(filter)
	c.JSON(http.StatusOK, appErrors.Ok(dto.OrderListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) adminOrderDetail(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "订单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok := handler.service.GetOrderDetail(orderID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "ORDER_NOT_FOUND",
			"message":   "订单不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) listAdminInvoices(c *gin.Context) {
	filter := parseInvoiceListFilter(c)
	items, total := handler.service.ListInvoices(filter)
	c.JSON(http.StatusOK, appErrors.Ok(dto.InvoiceListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) adminInvoiceDetail(c *gin.Context) {
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "账单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok := handler.service.GetInvoiceDetail(invoiceID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "INVOICE_NOT_FOUND",
			"message":   "账单不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminUpdateUnpaidInvoice(c *gin.Context) {
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "账单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.UpdateUnpaidInvoiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "账单调整参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok, updateErr := handler.service.UpdateUnpaidInvoice(
		invoiceID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVOICE_UPDATE_FAILED",
			"message":   updateErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVOICE_NOT_FOUND",
			"message":   "账单不存在，无法人工调整",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminReceivePayment(c *gin.Context) {
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "账单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.ReceivePaymentRequest
	_ = c.ShouldBindJSON(&request)

	invoice, serviceRecord, payment, ok := handler.service.ReceiveInvoicePayment(
		invoiceID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVOICE_NOT_FOUND",
			"message":   "账单不存在或当前状态不可收款",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"invoice": invoice,
		"service": serviceRecord,
		"payment": payment,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) adminRefundInvoice(c *gin.Context) {
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "账单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.RefundRequest
	_ = c.ShouldBindJSON(&request)
	if strings.TrimSpace(request.Reason) == "" {
		request.Reason = "后台手工退款"
	}

	invoice, serviceRecord, refund, ok := handler.service.RefundInvoice(
		invoiceID,
		request.Reason,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "REFUND_FAILED",
			"message":   "账单当前不可退款",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"invoice": invoice,
		"service": serviceRecord,
		"refund":  refund,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listAdminServices(c *gin.Context) {
	filter := parseServiceListFilter(c)
	items, total := handler.service.ListServices(filter)
	c.JSON(http.StatusOK, appErrors.Ok(dto.ServiceListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) adminServiceDetail(c *gin.Context) {
	serviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "服务编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok := handler.service.GetServiceDetail(serviceID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "SERVICE_NOT_FOUND",
			"message":   "服务不存在",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminUpdateService(c *gin.Context) {
	serviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "鏈嶅姟缂栧彿鏍煎紡涓嶆纭?",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.UpdateServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "鏈嶅姟璋冩暣鍙傛暟涓嶆纭?",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok, updateErr := handler.service.UpdateServiceRecord(
		serviceID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "SERVICE_UPDATE_FAILED",
			"message":   updateErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "SERVICE_NOT_FOUND",
			"message":   "鏈嶅姟涓嶅瓨鍦紝鏃犳硶浜哄伐璋冩暣",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) adminServiceAction(c *gin.Context) {
	serviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "服务编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	action := c.Param("action")
	var request dto.ServiceActionRequest
	_ = c.ShouldBindJSON(&request)

	result, ok := handler.service.ExecuteServiceAction(
		serviceID,
		action,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "SERVICE_ACTION_FAILED",
			"message":   "服务动作执行失败",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) listPortalOrders(c *gin.Context) {
	customerID := getPortalCustomerID(c)
	items := handler.service.ListOrdersByCustomer(customerID)
	c.JSON(http.StatusOK, appErrors.Ok(dto.OrderListResponse{
		Items: items,
		Total: len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listPortalInvoices(c *gin.Context) {
	customerID := getPortalCustomerID(c)
	items := handler.service.ListInvoicesByCustomer(customerID)
	c.JSON(http.StatusOK, appErrors.Ok(dto.InvoiceListResponse{
		Items: items,
		Total: len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listPortalServices(c *gin.Context) {
	customerID := getPortalCustomerID(c)
	items := handler.service.ListServicesByCustomer(customerID)
	c.JSON(http.StatusOK, appErrors.Ok(dto.ServiceListResponse{
		Items: items,
		Total: len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) checkout(c *gin.Context) {
	var request dto.CheckoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "下单参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, err := handler.service.Checkout(
		getPortalCustomerID(c),
		getPortalCustomerName(c),
		request,
		middleware.GetRequestID(c),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "CHECKOUT_FAILED",
			"message":   err.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusCreated, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) payInvoice(c *gin.Context) {
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "账单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	invoice, serviceRecord, payment, ok := handler.service.PayInvoice(
		getPortalCustomerID(c),
		getPortalCustomerName(c),
		invoiceID,
		middleware.GetRequestID(c),
	)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "INVOICE_NOT_FOUND",
			"message":   "账单不存在或当前不可支付",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"invoice": invoice,
		"service": serviceRecord,
		"payment": payment,
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

func getPortalCustomerName(c *gin.Context) string {
	value, exists := c.Get("customerName")
	if !exists {
		return "演示客户"
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
		return "系统管理员"
	}
	adminName, _ := value.(string)
	return adminName
}

func parseOrderListFilter(c *gin.Context) (domain.OrderListFilter, error) {
	filter := domain.OrderListFilter{
		Page:        parsePositiveInt(c.DefaultQuery("page", "1"), 1),
		Limit:       parsePositiveInt(c.DefaultQuery("limit", "20"), 20),
		Sort:        strings.TrimSpace(c.DefaultQuery("sort", "create_time")),
		Order:       strings.TrimSpace(c.DefaultQuery("order", "desc")),
		Status:      strings.TrimSpace(c.Query("status")),
		OrderNo:     strings.TrimSpace(c.Query("ordernum")),
		StartTime:   strings.TrimSpace(c.Query("start_time")),
		EndTime:     strings.TrimSpace(c.Query("end_time")),
		Payment:     strings.TrimSpace(c.Query("payment")),
		PayStatus:   strings.TrimSpace(c.Query("pay_status")),
		ProductName: strings.TrimSpace(c.Query("product_name")),
	}

	if rawAmount := strings.TrimSpace(c.Query("amount")); rawAmount != "" {
		amount, err := strconv.ParseFloat(rawAmount, 64)
		if err != nil {
			return domain.OrderListFilter{}, err
		}
		filter.Amount = amount
		filter.HasAmount = true
	}

	if rawCustomerID := strings.TrimSpace(c.Query("uid")); rawCustomerID != "" {
		customerID, err := strconv.ParseInt(rawCustomerID, 10, 64)
		if err != nil {
			return domain.OrderListFilter{}, err
		}
		filter.CustomerID = customerID
	}

	return filter, nil
}

func parseInvoiceListFilter(c *gin.Context) domain.InvoiceListFilter {
	return domain.InvoiceListFilter{
		Page:         parsePositiveInt(c.DefaultQuery("page", "1"), 1),
		Limit:        parsePositiveInt(c.DefaultQuery("limit", "20"), 20),
		Sort:         strings.TrimSpace(c.DefaultQuery("sort", "created_at")),
		Order:        strings.TrimSpace(c.DefaultQuery("order", "desc")),
		Status:       strings.TrimSpace(c.Query("status")),
		InvoiceNo:    strings.TrimSpace(c.Query("invoice_no")),
		OrderNo:      strings.TrimSpace(c.Query("order_no")),
		ProductName:  strings.TrimSpace(c.Query("product_name")),
		BillingCycle: strings.TrimSpace(c.Query("billing_cycle")),
		CustomerID:   parseOptionalInt64(c.Query("uid")),
	}
}

func parseServiceListFilter(c *gin.Context) domain.ServiceListFilter {
	return domain.ServiceListFilter{
		Page:              parsePositiveInt(c.DefaultQuery("page", "1"), 1),
		Limit:             parsePositiveInt(c.DefaultQuery("limit", "20"), 20),
		Sort:              strings.TrimSpace(c.DefaultQuery("sort", "created_at")),
		Order:             strings.TrimSpace(c.DefaultQuery("order", "desc")),
		Status:            strings.TrimSpace(c.Query("status")),
		ProviderType:      strings.TrimSpace(c.Query("provider_type")),
		ProviderAccountID: parseOptionalInt64(c.Query("provider_account_id")),
		SyncStatus:        strings.TrimSpace(c.Query("sync_status")),
		Keyword:           strings.TrimSpace(c.Query("keyword")),
	}
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
