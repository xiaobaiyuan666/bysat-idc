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
	router.POST("/orders/:id/requests", handler.adminCreateOrderRequest)
	router.GET("/order-requests", handler.listAdminOrderRequests)
	router.PATCH("/order-requests/:id", handler.adminProcessOrderRequest)
	router.GET("/invoices", handler.listAdminInvoices)
	router.GET("/invoices/:id", handler.adminInvoiceDetail)
	router.PATCH("/invoices/:id", handler.adminUpdateUnpaidInvoice)
	router.POST("/invoices/:id/receive-payment", handler.adminReceivePayment)
	router.POST("/invoices/:id/refund", handler.adminRefundInvoice)
	router.GET("/payments", handler.listAdminPayments)
	router.GET("/refunds", handler.listAdminRefunds)
	router.GET("/service-change-orders", handler.listAdminServiceChangeOrders)
	router.GET("/accounts/transactions", handler.listAdminAccountTransactions)
	router.GET("/accounts/customers/:id/wallet", handler.adminCustomerWallet)
	router.POST("/accounts/customers/:id/adjustments", handler.adminAdjustCustomerWallet)
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

func (handler *Handler) adminCreateOrderRequest(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "订单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.CreateOrderRequestRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "订单申请参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok, createErr := handler.service.CreateOrderRequest(
		orderID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if createErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "ORDER_REQUEST_CREATE_FAILED",
			"message":   createErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "ORDER_NOT_FOUND",
			"message":   "订单不存在，无法创建申请",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) listAdminOrderRequests(c *gin.Context) {
	filter := parseOrderRequestListFilter(c)
	items, total := handler.service.ListOrderRequests(filter)
	c.JSON(http.StatusOK, appErrors.Ok(dto.OrderRequestListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) adminProcessOrderRequest(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "申请编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.ProcessOrderRequestRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "申请处理参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok, processErr := handler.service.ProcessOrderRequest(
		requestID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if processErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "ORDER_REQUEST_PROCESS_FAILED",
			"message":   processErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "ORDER_REQUEST_NOT_FOUND",
			"message":   "订单申请不存在",
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
	router.GET("/orders/:id", handler.portalOrderDetail)
	router.POST("/orders/:id/requests", handler.portalCreateOrderRequest)
	router.GET("/invoices", handler.listPortalInvoices)
	router.GET("/invoices/:id", handler.portalInvoiceDetail)
	router.GET("/payments", handler.listPortalPayments)
	router.GET("/refunds", handler.listPortalRefunds)
	router.GET("/services", handler.listPortalServices)
	router.GET("/services/:id", handler.portalServiceDetail)
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

	invoice, serviceRecord, payment, ok, receiveErr := handler.service.ReceiveInvoicePayment(
		invoiceID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if receiveErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "RECEIVE_PAYMENT_FAILED",
			"message":   receiveErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
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

func (handler *Handler) listAdminPayments(c *gin.Context) {
	filter := parsePaymentListFilter(c)
	items, total := handler.service.ListPayments(filter)
	c.JSON(http.StatusOK, appErrors.Ok(dto.PaymentListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listAdminRefunds(c *gin.Context) {
	filter := parseRefundListFilter(c)
	items, total := handler.service.ListRefunds(filter)
	c.JSON(http.StatusOK, appErrors.Ok(dto.RefundListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listAdminAccountTransactions(c *gin.Context) {
	filter := parseAccountTransactionFilter(c)
	items, total := handler.service.ListAccountTransactions(filter)
	c.JSON(http.StatusOK, appErrors.Ok(dto.AccountTransactionListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) adminCustomerWallet(c *gin.Context) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "客户编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok := handler.service.GetCustomerWalletOverview(customerID, 20)
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

func (handler *Handler) adminAdjustCustomerWallet(c *gin.Context) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "客户编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.AdjustWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "钱包调整参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok, adjustErr := handler.service.AdjustCustomerWallet(
		customerID,
		request,
		getAdminID(c),
		getAdminName(c),
		middleware.GetRequestID(c),
	)
	if adjustErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "WALLET_ADJUST_FAILED",
			"message":   adjustErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
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

func (handler *Handler) listAdminServices(c *gin.Context) {
	filter := parseServiceListFilter(c)
	items, total := handler.service.ListServices(filter)
	c.JSON(http.StatusOK, appErrors.Ok(dto.ServiceListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listAdminServiceChangeOrders(c *gin.Context) {
	filter := parseServiceChangeOrderListFilter(c)
	items, total := handler.service.ListServiceChangeOrders(filter)
	c.JSON(http.StatusOK, appErrors.Ok(dto.ServiceChangeOrderListResponse{
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
			"message":   "服务编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.UpdateServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "服务调整参数不正确",
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
			"message":   "服务不存在，无法人工调整",
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

func (handler *Handler) portalOrderDetail(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "订单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	detail, ok := handler.service.GetOrderDetail(orderID)
	if !ok || detail.Order.CustomerID != getPortalCustomerID(c) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "ORDER_NOT_FOUND",
			"message":   "订单不存在或无权查看",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(detail, middleware.GetRequestID(c)))
}

func (handler *Handler) portalCreateOrderRequest(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "订单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	var request dto.CreateOrderRequestRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "订单申请参数不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	result, ok, createErr := handler.service.CreatePortalOrderRequest(
		getPortalCustomerID(c),
		getPortalCustomerName(c),
		orderID,
		request,
		middleware.GetRequestID(c),
	)
	if createErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "ORDER_REQUEST_CREATE_FAILED",
			"message":   createErr.Error(),
			"requestId": middleware.GetRequestID(c),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "ORDER_NOT_FOUND",
			"message":   "订单不存在或无权提交申请",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(result, middleware.GetRequestID(c)))
}

func (handler *Handler) listPortalInvoices(c *gin.Context) {
	customerID := getPortalCustomerID(c)
	items := handler.service.ListInvoicesByCustomer(customerID)
	c.JSON(http.StatusOK, appErrors.Ok(dto.InvoiceListResponse{
		Items: items,
		Total: len(items),
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) portalInvoiceDetail(c *gin.Context) {
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "账单编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	detail, ok := handler.service.GetInvoiceDetail(invoiceID)
	if !ok || detail.Invoice.CustomerID != getPortalCustomerID(c) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "INVOICE_NOT_FOUND",
			"message":   "账单不存在或无权查看",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(detail, middleware.GetRequestID(c)))
}

func (handler *Handler) listPortalPayments(c *gin.Context) {
	items, total := handler.service.ListPayments(domain.PaymentListFilter{
		Page:       1,
		Limit:      200,
		Sort:       "paid_at",
		Order:      "desc",
		CustomerID: getPortalCustomerID(c),
		InvoiceID:  parseOptionalInt64(c.Query("invoiceId")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
		Channel:    strings.TrimSpace(c.Query("channel")),
		Status:     strings.TrimSpace(c.Query("status")),
	})
	c.JSON(http.StatusOK, appErrors.Ok(dto.PaymentListResponse{
		Items: items,
		Total: total,
	}, middleware.GetRequestID(c)))
}

func (handler *Handler) listPortalRefunds(c *gin.Context) {
	items, total := handler.service.ListRefunds(domain.RefundListFilter{
		Page:       1,
		Limit:      200,
		Sort:       "created_at",
		Order:      "desc",
		CustomerID: getPortalCustomerID(c),
		InvoiceID:  parseOptionalInt64(c.Query("invoiceId")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
		Status:     strings.TrimSpace(c.Query("status")),
	})
	c.JSON(http.StatusOK, appErrors.Ok(dto.RefundListResponse{
		Items: items,
		Total: total,
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

func (handler *Handler) portalServiceDetail(c *gin.Context) {
	serviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      "INVALID_ARGUMENT",
			"message":   "服务编号格式不正确",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	detail, ok := handler.service.GetServiceDetail(serviceID)
	if !ok || detail.Service.CustomerID != getPortalCustomerID(c) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":      "SERVICE_NOT_FOUND",
			"message":   "服务不存在或无权查看",
			"requestId": middleware.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, appErrors.Ok(gin.H{
		"service":      detail.Service,
		"order":        detail.Order,
		"invoice":      detail.Invoice,
		"changeOrders": detail.ChangeOrders,
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
		CustomerID:        parseOptionalInt64(firstNonEmptyQueryValue(c, "customer_id", "customerId", "uid")),
		OrderID:           parseOptionalInt64(firstNonEmptyQueryValue(c, "order_id", "orderId")),
		ProviderType:      strings.TrimSpace(c.Query("provider_type")),
		ProviderAccountID: parseOptionalInt64(c.Query("provider_account_id")),
		SyncStatus:        strings.TrimSpace(c.Query("sync_status")),
		Keyword:           strings.TrimSpace(c.Query("keyword")),
	}
}

func parseServiceChangeOrderListFilter(c *gin.Context) domain.ServiceChangeOrderListFilter {
	return domain.ServiceChangeOrderListFilter{
		Page:      parsePositiveInt(c.DefaultQuery("page", "1"), 1),
		Limit:     parsePositiveInt(c.DefaultQuery("limit", "20"), 20),
		Sort:      strings.TrimSpace(c.DefaultQuery("sort", "created_at")),
		Order:     strings.TrimSpace(c.DefaultQuery("order", "desc")),
		Status:    strings.TrimSpace(c.Query("status")),
		ServiceID: parseOptionalInt64(c.Query("serviceId")),
		OrderID:   parseOptionalInt64(c.Query("orderId")),
		InvoiceID: parseOptionalInt64(c.Query("invoiceId")),
		Action:    strings.TrimSpace(c.Query("action")),
		Keyword:   strings.TrimSpace(c.Query("keyword")),
	}
}

func parseOrderRequestListFilter(c *gin.Context) domain.OrderRequestListFilter {
	return domain.OrderRequestListFilter{
		Page:       parsePositiveInt(c.DefaultQuery("page", "1"), 1),
		Limit:      parsePositiveInt(c.DefaultQuery("limit", "20"), 20),
		Sort:       strings.TrimSpace(c.DefaultQuery("sort", "created_at")),
		Order:      strings.TrimSpace(c.DefaultQuery("order", "desc")),
		OrderID:    parseOptionalInt64(c.Query("orderId")),
		CustomerID: parseOptionalInt64(firstNonEmptyQueryValue(c, "customerId", "uid")),
		Type:       strings.TrimSpace(c.Query("type")),
		Status:     strings.TrimSpace(c.Query("status")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
	}
}

func parseAccountTransactionFilter(c *gin.Context) domain.AccountTransactionFilter {
	return domain.AccountTransactionFilter{
		Page:            parsePositiveInt(c.DefaultQuery("page", "1"), 1),
		Limit:           parsePositiveInt(c.DefaultQuery("limit", "20"), 20),
		Sort:            strings.TrimSpace(c.DefaultQuery("sort", "occurred_at")),
		Order:           strings.TrimSpace(c.DefaultQuery("order", "desc")),
		CustomerID:      parseOptionalInt64(c.Query("customerId")),
		TransactionType: strings.TrimSpace(c.Query("transactionType")),
		Direction:       strings.TrimSpace(c.Query("direction")),
		Channel:         strings.TrimSpace(c.Query("channel")),
		Keyword:         strings.TrimSpace(c.Query("keyword")),
		StartTime:       strings.TrimSpace(c.Query("startTime")),
		EndTime:         strings.TrimSpace(c.Query("endTime")),
	}
}

func parsePaymentListFilter(c *gin.Context) domain.PaymentListFilter {
	return domain.PaymentListFilter{
		Page:       parsePositiveInt(c.DefaultQuery("page", "1"), 1),
		Limit:      parsePositiveInt(c.DefaultQuery("limit", "20"), 20),
		Sort:       strings.TrimSpace(c.DefaultQuery("sort", "paid_at")),
		Order:      strings.TrimSpace(c.DefaultQuery("order", "desc")),
		CustomerID: parseOptionalInt64(firstNonEmptyQueryValue(c, "customerId", "uid")),
		InvoiceID:  parseOptionalInt64(c.Query("invoiceId")),
		Keyword:    strings.TrimSpace(firstNonEmptyQueryValue(c, "keyword", "paymentNo", "tradeNo")),
		Channel:    strings.TrimSpace(c.Query("channel")),
		Status:     strings.TrimSpace(c.Query("status")),
	}
}

func parseRefundListFilter(c *gin.Context) domain.RefundListFilter {
	return domain.RefundListFilter{
		Page:       parsePositiveInt(c.DefaultQuery("page", "1"), 1),
		Limit:      parsePositiveInt(c.DefaultQuery("limit", "20"), 20),
		Sort:       strings.TrimSpace(c.DefaultQuery("sort", "created_at")),
		Order:      strings.TrimSpace(c.DefaultQuery("order", "desc")),
		CustomerID: parseOptionalInt64(firstNonEmptyQueryValue(c, "customerId", "uid")),
		InvoiceID:  parseOptionalInt64(c.Query("invoiceId")),
		Keyword:    strings.TrimSpace(firstNonEmptyQueryValue(c, "keyword", "refundNo")),
		Status:     strings.TrimSpace(c.Query("status")),
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

func firstNonEmptyQueryValue(c *gin.Context, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(c.Query(key)); value != "" {
			return value
		}
	}
	return ""
}
