package service

import (
	"fmt"
	"strconv"
	"strings"

	automationDomain "idc-finance/internal/modules/automation/domain"
	automationService "idc-finance/internal/modules/automation/service"
	catalogDomain "idc-finance/internal/modules/catalog/domain"
	catalogService "idc-finance/internal/modules/catalog/service"
	"idc-finance/internal/modules/order/domain"
	"idc-finance/internal/modules/order/dto"
	"idc-finance/internal/modules/order/repository"
	providerService "idc-finance/internal/modules/provider/service"
	"idc-finance/internal/platform/audit"
)

type Service struct {
	repository repository.Repository
	catalog    *catalogService.Service
	audit      *audit.Service
	provider   *providerService.Service
	tasks      *automationService.Service
}

func New(
	repo repository.Repository,
	catalogSvc *catalogService.Service,
	auditService *audit.Service,
	provider *providerService.Service,
	taskService *automationService.Service,
) *Service {
	return &Service{
		repository: repo,
		catalog:    catalogSvc,
		audit:      auditService,
		provider:   provider,
		tasks:      taskService,
	}
}

func (service *Service) ListOrders(filter domain.OrderListFilter) ([]domain.Order, int) {
	return service.repository.ListOrders(filter)
}

func (service *Service) ListOrdersByCustomer(customerID int64) []domain.Order {
	return service.repository.ListOrdersByCustomer(customerID)
}

func (service *Service) GetOrderDetail(orderID int64) (dto.OrderDetailResponse, bool) {
	order, exists := service.repository.GetOrderByID(orderID)
	if !exists {
		return dto.OrderDetailResponse{}, false
	}
	changeOrder, hasChangeOrder := service.repository.GetServiceChangeOrderByOrderID(orderID)
	return dto.OrderDetailResponse{
		Order:       order,
		Invoices:    service.repository.ListInvoicesByOrder(orderID),
		Services:    service.repository.ListServicesByOrder(orderID),
		Requests:    service.repository.ListOrderRequestsByOrder(orderID),
		ChangeOrder: service.buildChangeOrderRecord(changeOrder, hasChangeOrder),
		AuditLogs:   service.audit.ListByTarget("order", orderID),
	}, true
}

func (service *Service) UpdatePendingOrder(
	orderID int64,
	request dto.UpdatePendingOrderRequest,
	adminID int64,
	adminName,
	requestID string,
) (dto.OrderDetailResponse, bool, error) {
	beforeOrder, exists := service.repository.GetOrderByID(orderID)
	if !exists {
		return dto.OrderDetailResponse{}, false, nil
	}
	beforeInvoices := service.repository.ListInvoicesByOrder(orderID)
	beforeServices := service.repository.ListServicesByOrder(orderID)

	order, invoice, ok, err := service.repository.UpdatePendingOrder(orderID, domain.PendingOrderUpdate{
		ProductName:  request.ProductName,
		BillingCycle: request.BillingCycle,
		Amount:       request.Amount,
		DueAt:        request.DueAt,
		Status:       request.Status,
	})
	if err != nil {
		return dto.OrderDetailResponse{}, false, err
	}
	if !ok {
		return dto.OrderDetailResponse{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "order.manual_adjust",
		TargetType:  "order",
		TargetID:    order.ID,
		Target:      order.OrderNo,
		RequestID:   requestID,
		Description: "管理员人工调整订单内容，并同步回写关联账单",
		Payload: map[string]any{
			"reason": firstNonEmptyString(strings.TrimSpace(request.Reason), "后台人工修单"),
			"before": map[string]any{
				"status":          beforeOrder.Status,
				"productName":     beforeOrder.ProductName,
				"billingCycle":    beforeOrder.BillingCycle,
				"amount":          beforeOrder.Amount,
				"dueAt":           firstInvoiceDueAt(beforeInvoices),
				"invoiceStatuses": summarizeInvoices(beforeInvoices),
				"serviceStatuses": summarizeServices(beforeServices),
			},
			"after": map[string]any{
				"status":          order.Status,
				"productName":     order.ProductName,
				"billingCycle":    order.BillingCycle,
				"amount":          order.Amount,
				"dueAt":           valueOrZeroInvoiceDueAt(invoice),
				"invoiceId":       valueOrZeroInvoice(invoice),
				"invoiceStatuses": summarizeInvoices(service.repository.ListInvoicesByOrder(order.ID)),
				"serviceStatuses": summarizeServices(service.repository.ListServicesByOrder(order.ID)),
			},
		},
	})

	return dto.OrderDetailResponse{
		Order:     order,
		Invoices:  service.repository.ListInvoicesByOrder(order.ID),
		Services:  service.repository.ListServicesByOrder(order.ID),
		AuditLogs: service.audit.ListByTarget("order", order.ID),
	}, true, nil
}

func (service *Service) ListInvoices(filter domain.InvoiceListFilter) ([]domain.Invoice, int) {
	return service.repository.ListInvoices(filter)
}

func (service *Service) ListInvoicesByCustomer(customerID int64) []domain.Invoice {
	return service.repository.ListInvoicesByCustomer(customerID)
}

func (service *Service) GetCustomerWallet(customerID int64) (domain.CustomerWallet, bool) {
	return service.repository.GetCustomerWallet(customerID)
}

func (service *Service) ListAccountTransactions(filter domain.AccountTransactionFilter) ([]domain.AccountTransaction, int) {
	return service.repository.ListAccountTransactions(filter)
}

func (service *Service) ListAccountTransactionsByCustomer(customerID int64, limit int) []domain.AccountTransaction {
	return service.repository.ListAccountTransactionsByCustomer(customerID, limit)
}

func (service *Service) ListPayments(filter domain.PaymentListFilter) ([]domain.PaymentRecord, int) {
	return service.repository.ListPayments(filter)
}

func (service *Service) ListRefunds(filter domain.RefundListFilter) ([]domain.RefundRecord, int) {
	return service.repository.ListRefunds(filter)
}

func (service *Service) ListOrderRequests(filter domain.OrderRequestListFilter) ([]domain.OrderRequest, int) {
	return service.repository.ListOrderRequests(filter)
}

func (service *Service) CreateOrderRequest(
	orderID int64,
	request dto.CreateOrderRequestRequest,
	adminID int64,
	adminName,
	requestID string,
) (dto.OrderDetailResponse, bool, error) {
	order, exists := service.repository.GetOrderByID(orderID)
	if !exists {
		return dto.OrderDetailResponse{}, false, nil
	}

	requestType := normalizeOrderRequestType(request.Type)
	if requestType == "" {
		return dto.OrderDetailResponse{}, false, fmt.Errorf("订单申请类型不正确")
	}
	if request.RequestedAmount != nil && *request.RequestedAmount < 0 {
		return dto.OrderDetailResponse{}, false, fmt.Errorf("申请金额不能小于 0")
	}
	reason := strings.TrimSpace(request.Reason)
	if reason == "" {
		return dto.OrderDetailResponse{}, false, fmt.Errorf("申请原因不能为空")
	}

	summary := strings.TrimSpace(request.Summary)
	if summary == "" {
		summary = defaultOrderRequestSummary(requestType, order.ProductName)
	}

	created, ok, err := service.repository.CreateOrderRequest(orderID, domain.OrderRequestCreateInput{
		Type:                  requestType,
		Summary:               summary,
		Reason:                reason,
		RequestedAmount:       request.RequestedAmount,
		RequestedBillingCycle: strings.TrimSpace(request.RequestedBillingCycle),
		SourceType:            "ADMIN",
		SourceID:              adminID,
		SourceName:            adminName,
		Payload:               request.Payload,
	})
	if err != nil {
		return dto.OrderDetailResponse{}, false, err
	}
	if !ok {
		return dto.OrderDetailResponse{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "order.request.create",
		TargetType:  "order",
		TargetID:    order.ID,
		Target:      order.OrderNo,
		RequestID:   requestID,
		Description: "后台创建订单申请，进入后续审批或跟进流程",
		Payload: map[string]any{
			"requestNo":             created.RequestNo,
			"type":                  created.Type,
			"summary":               created.Summary,
			"reason":                created.Reason,
			"currentAmount":         created.CurrentAmount,
			"requestedAmount":       created.RequestedAmount,
			"currentBillingCycle":   created.CurrentBillingCycle,
			"requestedBillingCycle": created.RequestedBillingCycle,
			"status":                created.Status,
			"sourceType":            created.SourceType,
			"sourceName":            created.SourceName,
		},
	})

	detail, ok := service.GetOrderDetail(orderID)
	return detail, ok, nil
}

func (service *Service) CreatePortalOrderRequest(
	customerID int64,
	customerName string,
	orderID int64,
	request dto.CreateOrderRequestRequest,
	requestID string,
) (dto.OrderDetailResponse, bool, error) {
	order, exists := service.repository.GetOrderByID(orderID)
	if !exists || order.CustomerID != customerID {
		return dto.OrderDetailResponse{}, false, nil
	}

	requestType := normalizeOrderRequestType(request.Type)
	if requestType == "" {
		return dto.OrderDetailResponse{}, false, fmt.Errorf("订单申请类型不正确")
	}
	if request.RequestedAmount != nil && *request.RequestedAmount < 0 {
		return dto.OrderDetailResponse{}, false, fmt.Errorf("申请金额不能小于 0")
	}
	reason := strings.TrimSpace(request.Reason)
	if reason == "" {
		return dto.OrderDetailResponse{}, false, fmt.Errorf("申请原因不能为空")
	}

	summary := strings.TrimSpace(request.Summary)
	if summary == "" {
		summary = defaultOrderRequestSummary(requestType, order.ProductName)
	}

	created, ok, err := service.repository.CreateOrderRequest(orderID, domain.OrderRequestCreateInput{
		Type:                  requestType,
		Summary:               summary,
		Reason:                reason,
		RequestedAmount:       request.RequestedAmount,
		RequestedBillingCycle: strings.TrimSpace(request.RequestedBillingCycle),
		SourceType:            "CUSTOMER",
		SourceID:              customerID,
		SourceName:            customerName,
		Payload:               request.Payload,
	})
	if err != nil {
		return dto.OrderDetailResponse{}, false, err
	}
	if !ok {
		return dto.OrderDetailResponse{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "CUSTOMER",
		ActorID:     customerID,
		Actor:       customerName,
		Action:      "order.request.portal_create",
		TargetType:  "order",
		TargetID:    order.ID,
		Target:      order.OrderNo,
		RequestID:   requestID,
		Description: "客户在门户提交订单申请",
		Payload: map[string]any{
			"requestNo":             created.RequestNo,
			"type":                  created.Type,
			"summary":               created.Summary,
			"reason":                created.Reason,
			"requestedAmount":       created.RequestedAmount,
			"requestedBillingCycle": created.RequestedBillingCycle,
		},
	})

	detail, detailOK := service.GetOrderDetail(orderID)
	return detail, detailOK, nil
}

func (service *Service) ProcessOrderRequest(
	requestIDValue int64,
	request dto.ProcessOrderRequestRequest,
	adminID int64,
	adminName,
	requestTraceID string,
) (dto.OrderDetailResponse, bool, error) {
	status := normalizeOrderRequestStatus(request.Status)
	if status == "" || status == string(domain.OrderRequestStatusPending) {
		return dto.OrderDetailResponse{}, false, fmt.Errorf("订单申请处理状态不正确")
	}

	updated, ok, err := service.repository.ProcessOrderRequest(requestIDValue, domain.OrderRequestProcessInput{
		Status:        status,
		ProcessNote:   strings.TrimSpace(request.ProcessNote),
		ProcessorType: "ADMIN",
		ProcessorID:   adminID,
		ProcessorName: adminName,
	})
	if err != nil {
		return dto.OrderDetailResponse{}, false, err
	}
	if !ok {
		return dto.OrderDetailResponse{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "order.request.process",
		TargetType:  "order",
		TargetID:    updated.OrderID,
		Target:      updated.OrderNo,
		RequestID:   requestTraceID,
		Description: "后台处理订单申请，并记录审批结果",
		Payload: map[string]any{
			"requestNo":     updated.RequestNo,
			"type":          updated.Type,
			"status":        updated.Status,
			"processNote":   updated.ProcessNote,
			"processorName": updated.ProcessorName,
			"processedAt":   updated.ProcessedAt,
		},
	})

	detail, exists := service.GetOrderDetail(updated.OrderID)
	return detail, exists, nil
}

func (service *Service) ListServiceChangeOrders(filter domain.ServiceChangeOrderListFilter) ([]dto.ServiceChangeOrderRecord, int) {
	items, total := service.repository.ListServiceChangeOrders(filter)
	return service.buildChangeOrderRecords(items), total
}

func (service *Service) GetCustomerWalletOverview(customerID int64, limit int) (dto.CustomerWalletResponse, bool) {
	wallet, ok := service.repository.GetCustomerWallet(customerID)
	if !ok {
		return dto.CustomerWalletResponse{}, false
	}
	return dto.CustomerWalletResponse{
		Wallet:       wallet,
		Transactions: service.repository.ListAccountTransactionsByCustomer(customerID, limit),
	}, true
}

func (service *Service) AdjustCustomerWallet(
	customerID int64,
	request dto.AdjustWalletRequest,
	adminID int64,
	adminName,
	requestID string,
) (dto.AdjustWalletResponse, bool, error) {
	if request.Amount < 0 {
		return dto.AdjustWalletResponse{}, false, fmt.Errorf("adjust amount cannot be negative")
	}

	wallet, transaction, ok, err := service.repository.AdjustCustomerWallet(customerID, domain.WalletAdjustment{
		Target:       request.Target,
		Operation:    request.Operation,
		Amount:       request.Amount,
		Summary:      request.Summary,
		Remark:       request.Remark,
		OperatorType: "ADMIN",
		OperatorID:   adminID,
		OperatorName: adminName,
	})
	if err != nil {
		return dto.AdjustWalletResponse{}, false, err
	}
	if !ok {
		return dto.AdjustWalletResponse{}, false, nil
	}

	action := "finance.balance_adjust"
	if strings.EqualFold(strings.TrimSpace(request.Target), "CREDIT_LIMIT") {
		action = "finance.credit_limit_adjust"
	}
	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      action,
		TargetType:  "customer",
		TargetID:    customerID,
		Target:      wallet.CustomerNo,
		RequestID:   requestID,
		Description: "后台调整客户钱包额度",
		Payload: map[string]any{
			"target":        request.Target,
			"operation":     request.Operation,
			"amount":        request.Amount,
			"transactionNo": transaction.TransactionNo,
			"summary":       transaction.Summary,
			"remark":        transaction.Remark,
			"balanceAfter":  wallet.Balance,
			"creditAfter":   wallet.CreditLimit,
		},
	})

	return dto.AdjustWalletResponse{
		Wallet:      wallet,
		Transaction: transaction,
	}, true, nil
}

func (service *Service) GetInvoiceDetail(invoiceID int64) (dto.InvoiceDetailResponse, bool) {
	invoice, exists := service.repository.GetInvoiceByID(invoiceID)
	if !exists {
		return dto.InvoiceDetailResponse{}, false
	}

	var order *domain.Order
	if linkedOrder, ok := service.repository.GetOrderByID(invoice.OrderID); ok {
		order = &linkedOrder
	}
	changeOrder, hasChangeOrder := service.repository.GetServiceChangeOrderByInvoiceID(invoiceID)

	return dto.InvoiceDetailResponse{
		Invoice:     invoice,
		Order:       order,
		Services:    service.repository.ListServicesByOrder(invoice.OrderID),
		ChangeOrder: service.buildChangeOrderRecord(changeOrder, hasChangeOrder),
		Payments:    service.repository.ListPaymentsByInvoice(invoiceID),
		Refunds:     service.repository.ListRefundsByInvoice(invoiceID),
		AuditLogs:   service.audit.ListByTarget("invoice", invoiceID),
	}, true
}

func (service *Service) UpdateUnpaidInvoice(
	invoiceID int64,
	request dto.UpdateUnpaidInvoiceRequest,
	adminID int64,
	adminName,
	requestID string,
) (dto.InvoiceDetailResponse, bool, error) {
	beforeInvoice, exists := service.repository.GetInvoiceByID(invoiceID)
	if !exists {
		return dto.InvoiceDetailResponse{}, false, nil
	}
	beforePayments := service.repository.ListPaymentsByInvoice(invoiceID)
	beforeRefunds := service.repository.ListRefundsByInvoice(invoiceID)
	beforeServices := make([]domain.ServiceRecord, 0)
	var beforeOrder *domain.Order
	if beforeInvoice.OrderID > 0 {
		if item, ok := service.repository.GetOrderByID(beforeInvoice.OrderID); ok {
			beforeOrder = &item
		}
		beforeServices = service.repository.ListServicesByOrder(beforeInvoice.OrderID)
	}

	invoice, order, ok, err := service.repository.UpdateUnpaidInvoice(invoiceID, domain.UnpaidInvoiceUpdate{
		ProductName:  request.ProductName,
		BillingCycle: request.BillingCycle,
		Amount:       request.Amount,
		DueAt:        request.DueAt,
		Status:       request.Status,
	})
	if err != nil {
		return dto.InvoiceDetailResponse{}, false, err
	}
	if !ok {
		return dto.InvoiceDetailResponse{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "invoice.manual_adjust",
		TargetType:  "invoice",
		TargetID:    invoice.ID,
		Target:      invoice.InvoiceNo,
		RequestID:   requestID,
		Description: "管理员人工调整账单内容，并同步回写关联订单",
		Payload: map[string]any{
			"reason": firstNonEmptyString(strings.TrimSpace(request.Reason), "后台人工修单"),
			"before": map[string]any{
				"status":          beforeInvoice.Status,
				"productName":     beforeInvoice.ProductName,
				"billingCycle":    beforeInvoice.BillingCycle,
				"amount":          beforeInvoice.TotalAmount,
				"dueAt":           beforeInvoice.DueAt,
				"orderStatus":     valueOrZeroOrderStatus(beforeOrder),
				"paymentCount":    len(beforePayments),
				"refundCount":     len(beforeRefunds),
				"serviceStatuses": summarizeServices(beforeServices),
			},
			"after": map[string]any{
				"status":          invoice.Status,
				"productName":     invoice.ProductName,
				"billingCycle":    invoice.BillingCycle,
				"amount":          invoice.TotalAmount,
				"dueAt":           invoice.DueAt,
				"orderId":         valueOrZeroOrder(order),
				"orderStatus":     valueOrZeroOrderStatus(order),
				"paymentCount":    len(service.repository.ListPaymentsByInvoice(invoice.ID)),
				"refundCount":     len(service.repository.ListRefundsByInvoice(invoice.ID)),
				"serviceStatuses": summarizeServices(service.repository.ListServicesByOrder(invoice.OrderID)),
			},
		},
	})

	changeOrder, isServiceChange := service.repository.GetServiceChangeOrderByInvoiceID(invoice.ID)
	if isServiceChange && beforeInvoice.Status != domain.InvoiceStatusPaid && invoice.Status == domain.InvoiceStatusPaid {
		if _, _, err := service.executeServiceChangeAfterPayment(changeOrder, "ADMIN", adminID, adminName, requestID); err != nil {
			service.audit.Record(audit.Entry{
				ActorType:   "ADMIN",
				ActorID:     adminID,
				Actor:       adminName,
				Action:      "service.change_order.execute_failed",
				TargetType:  "invoice",
				TargetID:    invoice.ID,
				Target:      invoice.InvoiceNo,
				RequestID:   requestID,
				Description: "管理员改签账单为已支付后，自动执行改配动作失败",
				Payload: map[string]any{
					"changeOrderId":  changeOrder.ID,
					"changeAction":   changeOrder.ActionName,
					"executionError": err.Error(),
				},
			})
		}
	}

	return dto.InvoiceDetailResponse{
		Invoice:   invoice,
		Order:     order,
		Services:  service.repository.ListServicesByOrder(invoice.OrderID),
		Payments:  service.repository.ListPaymentsByInvoice(invoice.ID),
		Refunds:   service.repository.ListRefundsByInvoice(invoice.ID),
		AuditLogs: service.audit.ListByTarget("invoice", invoice.ID),
	}, true, nil
}

func (service *Service) ListServices(filter domain.ServiceListFilter) ([]domain.ServiceRecord, int) {
	return service.repository.ListServices(filter)
}

func (service *Service) ListServicesByCustomer(customerID int64) []domain.ServiceRecord {
	return service.repository.ListServicesByCustomer(customerID)
}

func (service *Service) GetServiceDetail(serviceID int64) (dto.ServiceDetailResponse, bool) {
	serviceRecord, exists := service.repository.GetServiceByID(serviceID)
	if !exists {
		return dto.ServiceDetailResponse{}, false
	}

	var order *domain.Order
	if linkedOrder, ok := service.repository.GetOrderByID(serviceRecord.OrderID); ok {
		order = &linkedOrder
	}

	var invoice *domain.Invoice
	if linkedInvoice, ok := service.repository.GetInvoiceByID(serviceRecord.InvoiceID); ok {
		invoice = &linkedInvoice
	}
	changeOrders := service.buildChangeOrderRecords(service.repository.ListServiceChangeOrdersByService(serviceID))

	return dto.ServiceDetailResponse{
		Service:      serviceRecord,
		Order:        order,
		Invoice:      invoice,
		ChangeOrders: changeOrders,
		AuditLogs:    service.audit.ListByTarget("service", serviceID),
	}, true
}

func (service *Service) UpdateServiceRecord(
	serviceID int64,
	request dto.UpdateServiceRequest,
	adminID int64,
	adminName,
	requestID string,
) (dto.ServiceDetailResponse, bool, error) {
	beforeService, exists := service.repository.GetServiceByID(serviceID)
	if !exists {
		return dto.ServiceDetailResponse{}, false, nil
	}
	var beforeOrder *domain.Order
	if beforeService.OrderID > 0 {
		if item, ok := service.repository.GetOrderByID(beforeService.OrderID); ok {
			beforeOrder = &item
		}
	}
	var beforeInvoice *domain.Invoice
	if beforeService.InvoiceID > 0 {
		if item, ok := service.repository.GetInvoiceByID(beforeService.InvoiceID); ok {
			beforeInvoice = &item
		}
	}

	updated, ok, err := service.repository.UpdateServiceRecord(serviceID, domain.ManualServiceUpdate{
		ProviderType:       request.ProviderType,
		ProviderAccountID:  request.ProviderAccountID,
		ProviderResourceID: request.ProviderResourceID,
		RegionName:         request.RegionName,
		IPAddress:          request.IPAddress,
		NextDueAt:          request.NextDueAt,
		Status:             request.Status,
		SyncStatus:         request.SyncStatus,
		SyncMessage:        request.SyncMessage,
	})
	if err != nil {
		return dto.ServiceDetailResponse{}, false, err
	}
	if !ok {
		return dto.ServiceDetailResponse{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "service.manual_adjust",
		TargetType:  "service",
		TargetID:    updated.ID,
		Target:      updated.ServiceNo,
		RequestID:   requestID,
		Description: "管理员人工调整服务状态、接口绑定与同步信息",
		Payload: map[string]any{
			"reason": firstNonEmptyString(strings.TrimSpace(request.Reason), "后台人工调整服务"),
			"before": map[string]any{
				"status":             beforeService.Status,
				"providerType":       beforeService.ProviderType,
				"providerAccountId":  beforeService.ProviderAccountID,
				"providerResourceId": beforeService.ProviderResourceID,
				"regionName":         beforeService.RegionName,
				"ipAddress":          beforeService.IPAddress,
				"nextDueAt":          beforeService.NextDueAt,
				"syncStatus":         beforeService.SyncStatus,
				"syncMessage":        beforeService.SyncMessage,
				"orderStatus":        valueOrZeroOrderStatus(beforeOrder),
				"invoiceStatus":      valueOrZeroInvoiceStatus(beforeInvoice),
			},
			"after": map[string]any{
				"status":             updated.Status,
				"providerType":       updated.ProviderType,
				"providerAccountId":  updated.ProviderAccountID,
				"providerResourceId": updated.ProviderResourceID,
				"regionName":         updated.RegionName,
				"ipAddress":          updated.IPAddress,
				"nextDueAt":          updated.NextDueAt,
				"syncStatus":         updated.SyncStatus,
				"syncMessage":        updated.SyncMessage,
				"orderStatus":        valueOrZeroOrderStatus(loadOrderIfExists(service.repository, updated.OrderID)),
				"invoiceStatus":      valueOrZeroInvoiceStatus(loadInvoiceIfExists(service.repository, updated.InvoiceID)),
			},
		},
	})

	detail, exists := service.GetServiceDetail(updated.ID)
	return detail, exists, nil
}

func (service *Service) CreateServiceChangeOrder(
	serviceID int64,
	request dto.CreateServiceChangeOrderRequest,
	adminID int64,
	adminName,
	requestID string,
) (dto.CreateServiceChangeOrderResponse, bool, error) {
	beforeService, exists := service.repository.GetServiceByID(serviceID)
	if !exists {
		return dto.CreateServiceChangeOrderResponse{}, false, nil
	}
	if request.Amount < 0 {
		return dto.CreateServiceChangeOrderResponse{}, false, fmt.Errorf("改配单金额不能小于 0")
	}

	order, invoice, serviceRecord, ok, err := service.repository.CreateServiceChangeOrder(serviceID, domain.ServiceChangeOrderInput{
		ActionName:   request.ActionName,
		Title:        request.Title,
		BillingCycle: request.BillingCycle,
		Amount:       request.Amount,
		Reason:       request.Reason,
		Payload:      request.Payload,
	})
	if err != nil {
		return dto.CreateServiceChangeOrderResponse{}, false, err
	}
	if !ok || serviceRecord == nil {
		return dto.CreateServiceChangeOrderResponse{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "service.change_order.create",
		TargetType:  "service",
		TargetID:    serviceRecord.ID,
		Target:      serviceRecord.ServiceNo,
		RequestID:   requestID,
		Description: "管理员为服务生成待支付改配单",
		Payload: map[string]any{
			"reason": firstNonEmptyString(strings.TrimSpace(request.Reason), "后台生成改配单"),
			"before": map[string]any{
				"serviceStatus":      beforeService.Status,
				"providerType":       beforeService.ProviderType,
				"providerAccountId":  beforeService.ProviderAccountID,
				"providerResourceId": beforeService.ProviderResourceID,
				"nextDueAt":          beforeService.NextDueAt,
			},
			"after": map[string]any{
				"actionName":   strings.TrimSpace(request.ActionName),
				"orderId":      order.ID,
				"orderNo":      order.OrderNo,
				"invoiceId":    invoice.ID,
				"invoiceNo":    invoice.InvoiceNo,
				"billingCycle": invoice.BillingCycle,
				"amount":       invoice.TotalAmount,
			},
		},
	})

	return dto.CreateServiceChangeOrderResponse{
		Service: *serviceRecord,
		Order:   order,
		Invoice: invoice,
	}, true, nil
}

func (service *Service) ExecuteServiceAction(serviceID int64, action string, params dto.ServiceActionRequest, adminID int64, adminName, requestID string) (domain.ServiceRecord, bool) {
	serviceRecord, exists := service.repository.GetServiceByID(serviceID)
	if !exists {
		return domain.ServiceRecord{}, false
	}

	taskID := service.startTask(automationService.StartTaskRequest{
		TaskType:           "SERVICE_ACTION",
		Title:              "鏈嶅姟鍔ㄤ綔鎵ц",
		Channel:            firstNonEmptyString(strings.TrimSpace(serviceRecord.ProviderType), "LOCAL"),
		Stage:              "MANUAL",
		SourceType:         "service",
		SourceID:           serviceRecord.ID,
		CustomerID:         serviceRecord.CustomerID,
		ProductName:        serviceRecord.ProductName,
		ServiceID:          serviceRecord.ID,
		ServiceNo:          serviceRecord.ServiceNo,
		ProviderType:       serviceRecord.ProviderType,
		ProviderResourceID: serviceRecord.ProviderResourceID,
		ActionName:         action,
		OperatorType:       "ADMIN",
		OperatorName:       adminName,
		RequestPayload:     params,
		Message:            "任务已进入执行阶段",
	})

	if serviceRecord.ProviderType == "MOFANG_CLOUD" && service.provider != nil {
		updated, ok := service.executeMofangCloudServiceAction(serviceRecord, action, params)
		if !ok {
			service.failTask(taskID, "鏈嶅姟鍔ㄤ綔鎵ц澶辫触", map[string]any{
				"action":      action,
				"serviceId":   serviceRecord.ID,
				"provider":    serviceRecord.ProviderType,
				"requestId":   requestID,
				"serviceNo":   serviceRecord.ServiceNo,
				"resourceId":  serviceRecord.ProviderResourceID,
				"targetState": serviceRecord.Status,
			})
			return domain.ServiceRecord{}, false
		}
		serviceRecord = updated
	} else {
		var ok bool
		serviceRecord, ok = service.repository.ExecuteServiceAction(serviceID, action, domain.ServiceActionParams{
			ImageName: params.ImageName,
			Password:  params.Password,
		})
		if !ok {
			service.failTask(taskID, "鏈湴鏈嶅姟鍔ㄤ綔鎵ц澶辫触", map[string]any{
				"action":    action,
				"serviceId": serviceID,
			})
			return domain.ServiceRecord{}, false
		}
	}

	service.successTask(taskID, fmt.Sprintf("服务动作 %s 已完成", action), map[string]any{
		"serviceId":    serviceRecord.ID,
		"serviceNo":    serviceRecord.ServiceNo,
		"status":       serviceRecord.Status,
		"lastAction":   serviceRecord.LastAction,
		"providerType": serviceRecord.ProviderType,
	})

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "service." + action,
		TargetType:  "service",
		TargetID:    serviceRecord.ID,
		Target:      serviceRecord.ServiceNo,
		RequestID:   requestID,
		Description: serviceActionDescription(action),
		Payload: map[string]any{
			"imageName":    params.ImageName,
			"hasSecret":    strings.TrimSpace(params.Password) != "",
			"providerType": serviceRecord.ProviderType,
		},
	})

	return serviceRecord, true
}

func (service *Service) executeMofangCloudServiceAction(
	serviceRecord domain.ServiceRecord,
	action string,
	params dto.ServiceActionRequest,
) (domain.ServiceRecord, bool) {
	remoteID := strings.TrimSpace(serviceRecord.ProviderResourceID)
	if remoteID == "" {
		return domain.ServiceRecord{}, false
	}

	providerAction, providerRequest, shouldPullSync, ok := mapMofangAction(action, params)
	if !ok {
		return domain.ServiceRecord{}, false
	}

	result, err := service.provider.ExecuteAction(remoteID, providerAction, providerRequest, serviceRecord.ProviderAccountID)
	if err != nil || !result.OK {
		return domain.ServiceRecord{}, false
	}

	localRecord, ok := service.repository.ExecuteServiceAction(serviceRecord.ID, action, domain.ServiceActionParams{
		ImageName: params.ImageName,
		Password:  params.Password,
	})
	if !ok {
		return domain.ServiceRecord{}, false
	}

	if shouldPullSync {
		if _, err := service.provider.SyncServiceByID(serviceRecord.ID, true); err == nil {
			if synced, exists := service.repository.GetServiceByID(serviceRecord.ID); exists {
				return synced, true
			}
		}
	}

	return localRecord, true
}

func mapMofangAction(
	action string,
	params dto.ServiceActionRequest,
) (string, providerService.ActionRequest, bool, bool) {
	request := providerService.ActionRequest{
		ImageName: params.ImageName,
		Password:  params.Password,
	}

	switch action {
	case "activate":
		return "activate", request, true, true
	case "suspend":
		return "suspend", request, true, true
	case "terminate":
		return "terminate", request, false, true
	case "reboot":
		return "reboot", request, true, true
	case "reset-password":
		return "reset-password", request, false, true
	case "reinstall":
		return "reinstall", request, true, true
	default:
		return "", providerService.ActionRequest{}, false, false
	}
}

func (service *Service) Checkout(customerID int64, customerName string, request dto.CheckoutRequest, requestID string) (dto.CheckoutResponse, error) {
	product, exists := service.catalog.GetByID(request.ProductID)
	if !exists {
		return dto.CheckoutResponse{}, fmt.Errorf("商品不存在")
	}

	price, err := pickPrice(product, request.CycleCode)
	if err != nil {
		return dto.CheckoutResponse{}, err
	}

	selectedConfiguration, err := buildSelectedConfiguration(product, request.SelectedOptions)
	if err != nil {
		return dto.CheckoutResponse{}, err
	}
	totalAmount := price.Price + price.SetupFee + configurationDelta(selectedConfiguration)
	automationType := resolveProductAutomationType(product)
	providerAccountID := resolveProductProviderAccountID(product)
	order, invoice := service.repository.Checkout(
		customerID,
		customerName,
		product.ID,
		product.Name,
		product.ProductType,
		automationType,
		providerAccountID,
		price.CycleCode,
		totalAmount,
		selectedConfiguration,
		buildResourceSnapshot(product, customerID, selectedConfiguration),
	)

	service.audit.Record(audit.Entry{
		ActorType:   "PORTAL",
		ActorID:     customerID,
		Actor:       customerName,
		Action:      "order.checkout",
		TargetType:  "order",
		TargetID:    order.ID,
		Target:      order.OrderNo,
		RequestID:   requestID,
		Description: "客户从门户提交新购订单",
		Payload: map[string]any{
			"productId":      product.ID,
			"cycle":          price.CycleCode,
			"amount":         totalAmount,
			"selectedConfig": len(selectedConfiguration),
			"automationType": automationType,
		},
	})

	return dto.CheckoutResponse{Order: order, Invoice: invoice}, nil
}

func (service *Service) PayInvoice(customerID int64, customerName string, invoiceID int64, requestID string) (domain.Invoice, *domain.ServiceRecord, domain.PaymentRecord, bool) {
	invoice, exists := service.repository.GetInvoiceByID(invoiceID)
	if !exists || invoice.CustomerID != customerID {
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}
	beforeStatus := invoice.Status
	paymentChannel := "ONLINE"
	if wallet, ok := service.repository.GetCustomerWallet(customerID); ok && wallet.Balance+0.00001 >= invoice.TotalAmount {
		paymentChannel = "BALANCE"
	}

	invoice, serviceRecord, payment, ok := service.repository.PayInvoice(invoiceID, paymentChannel, "PORTAL", customerName, "")
	if !ok {
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "PORTAL",
		ActorID:     customerID,
		Actor:       customerName,
		Action:      "invoice.pay",
		TargetType:  "invoice",
		TargetID:    invoice.ID,
		Target:      invoice.InvoiceNo,
		RequestID:   requestID,
		Description: "客户支付账单并激活订单",
		Payload: map[string]any{
			"amount":    invoice.TotalAmount,
			"paymentNo": payment.PaymentNo,
			"channel":   payment.Channel,
		},
	})

	if serviceRecord != nil {
		changeOrder, isServiceChange := service.repository.GetServiceChangeOrderByInvoiceID(invoice.ID)
		if isServiceChange && beforeStatus != domain.InvoiceStatusPaid && invoice.Status == domain.InvoiceStatusPaid {
			if updatedService, _, err := service.executeServiceChangeAfterPayment(changeOrder, "PORTAL", customerID, customerName, requestID); err == nil && updatedService != nil {
				serviceRecord = updatedService
			} else if err != nil {
				service.audit.Record(audit.Entry{
					ActorType:   "PORTAL",
					ActorID:     customerID,
					Actor:       customerName,
					Action:      "service.change_order.execute_failed",
					TargetType:  "service",
					TargetID:    serviceRecord.ID,
					Target:      serviceRecord.ServiceNo,
					RequestID:   requestID,
					Description: "改配单支付成功，但自动执行资源动作失败",
					Payload: map[string]any{
						"invoiceId":   invoice.ID,
						"invoiceNo":   invoice.InvoiceNo,
						"actionName":  changeOrder.ActionName,
						"changeOrder": changeOrder.ID,
						"error":       err.Error(),
					},
				})
			}
			return invoice, serviceRecord, payment, true
		}
		if order, exists := service.repository.GetOrderByID(invoice.OrderID); exists {
			serviceRecord = service.finalizeProvisionAfterPayment(order, serviceRecord, "PORTAL", customerID, customerName, requestID)
		}
		action := "service.activate"
		if serviceRecord.Status == domain.ServiceStatusPending {
			action = "service.provision_dispatch"
		}
		service.audit.Record(audit.Entry{
			ActorType:   "PORTAL",
			ActorID:     customerID,
			Actor:       customerName,
			Action:      action,
			TargetType:  "service",
			TargetID:    serviceRecord.ID,
			Target:      serviceRecord.ServiceNo,
			RequestID:   requestID,
			Description: "账单支付后自动激活服务",
		})
	}

	return invoice, serviceRecord, payment, true
}

func (service *Service) ReceiveInvoicePayment(invoiceID int64, request dto.ReceivePaymentRequest, adminID int64, adminName, requestID string) (domain.Invoice, *domain.ServiceRecord, domain.PaymentRecord, bool, error) {
	invoice, exists := service.repository.GetInvoiceByID(invoiceID)
	if !exists {
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false, nil
	}
	beforeStatus := invoice.Status
	channel := strings.ToUpper(strings.TrimSpace(firstNonEmptyString(request.Channel, "OFFLINE")))
	if channel == "BALANCE" {
		wallet, ok := service.repository.GetCustomerWallet(invoice.CustomerID)
		if !ok {
			return domain.Invoice{}, nil, domain.PaymentRecord{}, false, fmt.Errorf("客户钱包不存在，无法使用余额抵扣")
		}
		if wallet.Balance+0.00001 < invoice.TotalAmount {
			return domain.Invoice{}, nil, domain.PaymentRecord{}, false, fmt.Errorf("客户余额不足，无法使用余额抵扣")
		}
	}

	taskID := service.startTask(automationService.StartTaskRequest{
		TaskType:     "INVOICE_ACTION",
		Title:        "账单收款处理",
		Channel:      "LOCAL",
		Stage:        "FINANCE",
		SourceType:   "invoice",
		SourceID:     invoice.ID,
		CustomerID:   invoice.CustomerID,
		OrderID:      invoice.OrderID,
		InvoiceID:    invoice.ID,
		ActionName:   "receive-payment",
		OperatorType: "ADMIN",
		OperatorName: adminName,
		RequestPayload: map[string]any{
			"channel": channel,
			"tradeNo": request.TradeNo,
		},
		Message: "账单收款任务已启动",
	})

	invoice, serviceRecord, payment, ok := service.repository.PayInvoice(
		invoiceID,
		channel,
		"ADMIN",
		adminName,
		request.TradeNo,
	)
	if !ok {
		service.failTask(taskID, "账单收款失败", map[string]any{
			"invoiceId": invoiceID,
			"channel":   channel,
		})
		return domain.Invoice{}, nil, domain.PaymentRecord{}, false, fmt.Errorf("账单不存在或当前状态不可收款")
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "invoice.receive_payment",
		TargetType:  "invoice",
		TargetID:    invoice.ID,
		Target:      invoice.InvoiceNo,
		RequestID:   requestID,
		Description: "后台登记账单收款",
		Payload: map[string]any{
			"amount":     invoice.TotalAmount,
			"customerId": invoice.CustomerID,
			"paymentNo":  payment.PaymentNo,
			"channel":    payment.Channel,
		},
	})

	changeOrder, isServiceChange := service.repository.GetServiceChangeOrderByInvoiceID(invoice.ID)
	if serviceRecord != nil {
		if isServiceChange && beforeStatus != domain.InvoiceStatusPaid && invoice.Status == domain.InvoiceStatusPaid {
			if updatedService, resourceAction, err := service.executeServiceChangeAfterPayment(changeOrder, "ADMIN", adminID, adminName, requestID); err == nil {
				if updatedService != nil {
					serviceRecord = updatedService
				}
				service.successTask(taskID, "账单已收款，改配动作已下发", map[string]any{
					"invoiceId":      invoice.ID,
					"invoiceNo":      invoice.InvoiceNo,
					"paymentNo":      payment.PaymentNo,
					"serviceId":      valueOrZero(serviceRecord),
					"serviceNo":      serviceNoOrEmpty(serviceRecord),
					"invoiceStatus":  invoice.Status,
					"changeOrderId":  changeOrder.ID,
					"changeAction":   changeOrder.ActionName,
					"resourceAction": resourceAction,
				})
				return invoice, serviceRecord, payment, true, nil
			} else {
				service.audit.Record(audit.Entry{
					ActorType:   "ADMIN",
					ActorID:     adminID,
					Actor:       adminName,
					Action:      "service.change_order.execute_failed",
					TargetType:  "service",
					TargetID:    serviceRecord.ID,
					Target:      serviceRecord.ServiceNo,
					RequestID:   requestID,
					Description: "改配单已收款，但自动执行资源动作失败",
					Payload: map[string]any{
						"invoiceId":      invoice.ID,
						"invoiceNo":      invoice.InvoiceNo,
						"changeOrderId":  changeOrder.ID,
						"changeAction":   changeOrder.ActionName,
						"providerType":   serviceRecord.ProviderType,
						"providerId":     serviceRecord.ProviderResourceID,
						"executionError": err.Error(),
					},
				})
				service.successTask(taskID, "账单已收款，但改配动作执行失败", map[string]any{
					"invoiceId":      invoice.ID,
					"invoiceNo":      invoice.InvoiceNo,
					"paymentNo":      payment.PaymentNo,
					"serviceId":      valueOrZero(serviceRecord),
					"serviceNo":      serviceNoOrEmpty(serviceRecord),
					"invoiceStatus":  invoice.Status,
					"changeOrderId":  changeOrder.ID,
					"changeAction":   changeOrder.ActionName,
					"executionError": err.Error(),
				})
				return invoice, serviceRecord, payment, true, nil
			}
		}
		if order, exists := service.repository.GetOrderByID(invoice.OrderID); exists {
			serviceRecord = service.finalizeProvisionAfterPayment(order, serviceRecord, "ADMIN", adminID, adminName, requestID)
		}
		action := "service.activate"
		if serviceRecord.Status == domain.ServiceStatusPending {
			action = "service.provision_dispatch"
		}
		service.audit.Record(audit.Entry{
			ActorType:   "ADMIN",
			ActorID:     adminID,
			Actor:       adminName,
			Action:      action,
			TargetType:  "service",
			TargetID:    serviceRecord.ID,
			Target:      serviceRecord.ServiceNo,
			RequestID:   requestID,
			Description: "后台收款后激活服务",
		})
	}

	service.successTask(taskID, "账单收款已完成", map[string]any{
		"invoiceId":     invoice.ID,
		"invoiceNo":     invoice.InvoiceNo,
		"paymentNo":     payment.PaymentNo,
		"serviceId":     valueOrZero(serviceRecord),
		"serviceNo":     serviceNoOrEmpty(serviceRecord),
		"invoiceStatus": invoice.Status,
	})

	return invoice, serviceRecord, payment, true, nil
}

func (service *Service) RefundInvoice(invoiceID int64, reason string, adminID int64, adminName, requestID string) (domain.Invoice, *domain.ServiceRecord, domain.RefundRecord, bool) {
	taskID := service.startTask(automationService.StartTaskRequest{
		TaskType:     "INVOICE_ACTION",
		Title:        "账单退款处理",
		Channel:      "LOCAL",
		Stage:        "FINANCE",
		SourceType:   "invoice",
		SourceID:     invoiceID,
		InvoiceID:    invoiceID,
		ActionName:   "refund",
		OperatorType: "ADMIN",
		OperatorName: adminName,
		RequestPayload: map[string]any{
			"reason": reason,
		},
		Message: "账单退款任务已启动",
	})

	invoice, serviceRecord, refund, ok := service.repository.RefundInvoice(invoiceID, reason)
	if !ok {
		service.failTask(taskID, "账单退款失败", map[string]any{
			"invoiceId": invoiceID,
			"reason":    reason,
		})
		return domain.Invoice{}, nil, domain.RefundRecord{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "invoice.refund",
		TargetType:  "invoice",
		TargetID:    invoice.ID,
		Target:      invoice.InvoiceNo,
		RequestID:   requestID,
		Description: "后台执行账单退款",
		Payload: map[string]any{
			"refundNo": refund.RefundNo,
			"amount":   refund.Amount,
			"reason":   refund.Reason,
		},
	})

	if serviceRecord != nil {
		service.audit.Record(audit.Entry{
			ActorType:   "ADMIN",
			ActorID:     adminID,
			Actor:       adminName,
			Action:      "service.terminate",
			TargetType:  "service",
			TargetID:    serviceRecord.ID,
			Target:      serviceRecord.ServiceNo,
			RequestID:   requestID,
			Description: "閫€娆惧悗缁堟鍏宠仈鏈嶅姟",
		})
	}

	service.successTask(taskID, "账单退款已完成", map[string]any{
		"invoiceId":     invoice.ID,
		"invoiceNo":     invoice.InvoiceNo,
		"orderId":       invoice.OrderID,
		"customerId":    invoice.CustomerID,
		"refundNo":      refund.RefundNo,
		"refundAmount":  refund.Amount,
		"serviceId":     valueOrZero(serviceRecord),
		"serviceNo":     serviceNoOrEmpty(serviceRecord),
		"invoiceStatus": invoice.Status,
	})

	return invoice, serviceRecord, refund, true
}

func (service *Service) executeServiceChangeAfterPayment(
	changeOrder domain.ServiceChangeOrder,
	actorType string,
	actorID int64,
	actorName,
	requestID string,
) (*domain.ServiceRecord, *providerService.ResourceActionResponse, error) {
	if service.provider == nil {
		return nil, nil, fmt.Errorf("资源通道未初始化")
	}
	if changeOrder.ServiceID <= 0 {
		return nil, nil, fmt.Errorf("改配单缺少关联服务")
	}
	actionName := strings.TrimSpace(changeOrder.ActionName)
	if actionName == "" {
		return nil, nil, fmt.Errorf("改配单缺少资源动作")
	}

	request := buildResourceActionRequest(changeOrder.Payload)
	result, err := service.provider.ExecuteServiceResourceAction(changeOrder.ServiceID, actionName, request)
	if err != nil {
		return nil, nil, err
	}

	var updatedService *domain.ServiceRecord
	if item, ok := service.repository.GetServiceByID(changeOrder.ServiceID); ok {
		updatedService = &item
	}

	targetID := changeOrder.ServiceID
	targetNo := fmt.Sprintf("service-%d", changeOrder.ServiceID)
	if updatedService != nil {
		targetID = updatedService.ID
		targetNo = updatedService.ServiceNo
	}

	service.audit.Record(audit.Entry{
		ActorType:   actorType,
		ActorID:     actorID,
		Actor:       actorName,
		Action:      "service.change_order.execute",
		TargetType:  "service",
		TargetID:    targetID,
		Target:      targetNo,
		RequestID:   requestID,
		Description: "改配单支付后自动执行资源动作",
		Payload: map[string]any{
			"changeOrderId": changeOrder.ID,
			"invoiceId":     changeOrder.InvoiceID,
			"orderId":       changeOrder.OrderID,
			"actionName":    actionName,
			"request":       request,
			"result":        result,
		},
	})

	return updatedService, &result, nil
}

func (service *Service) buildChangeOrderRecords(items []domain.ServiceChangeOrder) []dto.ServiceChangeOrderRecord {
	result := make([]dto.ServiceChangeOrderRecord, 0, len(items))
	for _, item := range items {
		record := service.buildChangeOrderRecord(item, true)
		if record == nil {
			continue
		}
		result = append(result, *record)
	}
	return result
}

func (service *Service) buildChangeOrderRecord(item domain.ServiceChangeOrder, ok bool) *dto.ServiceChangeOrderRecord {
	if !ok {
		return nil
	}

	record := dto.ServiceChangeOrderRecord{
		ID:               item.ID,
		ServiceID:        item.ServiceID,
		OrderID:          item.OrderID,
		InvoiceID:        item.InvoiceID,
		ActionName:       item.ActionName,
		Title:            item.Title,
		Amount:           item.Amount,
		Status:           item.Status,
		Reason:           item.Reason,
		BillingCycle:     item.BillingCycle,
		Payload:          item.Payload,
		PaidAt:           item.PaidAt,
		RefundedAt:       item.RefundedAt,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
		ExecutionStatus:  "WAITING_PAYMENT",
		ExecutionMessage: "改配单尚未支付，资源动作不会自动下发。",
	}

	if serviceRecord, exists := service.repository.GetServiceByID(item.ServiceID); exists {
		record.ServiceNo = serviceRecord.ServiceNo
		record.ProductName = serviceRecord.ProductName
		record.ProviderType = serviceRecord.ProviderType
	}
	if order, exists := service.repository.GetOrderByID(item.OrderID); exists {
		record.OrderNo = order.OrderNo
		if strings.TrimSpace(record.ProductName) == "" {
			record.ProductName = order.ProductName
		}
	}
	if invoice, exists := service.repository.GetInvoiceByID(item.InvoiceID); exists {
		record.InvoiceNo = invoice.InvoiceNo
		if strings.TrimSpace(record.BillingCycle) == "" {
			record.BillingCycle = invoice.BillingCycle
		}
	}

	if item.Status == string(domain.InvoiceStatusRefunded) {
		record.ExecutionStatus = "REFUNDED"
		record.ExecutionMessage = "改配单已退款，已停止继续自动执行。"
		return &record
	}

	if item.Status != string(domain.InvoiceStatusPaid) {
		return &record
	}

	record.ExecutionStatus = "PAID"
	record.ExecutionMessage = "改配单已支付，等待系统回写执行结果。"
	if service.tasks == nil {
		return &record
	}

	tasks, _ := service.tasks.List(automationDomain.TaskFilter{
		InvoiceID: item.InvoiceID,
		Limit:     6,
	})
	for _, task := range tasks {
		if task.InvoiceID != item.InvoiceID {
			continue
		}
		taskCopy := task
		record.LatestTask = &taskCopy
		switch task.Status {
		case automationDomain.TaskStatusSuccess:
			record.ExecutionStatus = "EXECUTED"
		case automationDomain.TaskStatusFailed:
			record.ExecutionStatus = "EXECUTE_FAILED"
		case automationDomain.TaskStatusBlocked:
			record.ExecutionStatus = "EXECUTE_BLOCKED"
		case automationDomain.TaskStatusRunning:
			record.ExecutionStatus = "EXECUTING"
		default:
			record.ExecutionStatus = "PAID"
		}
		if strings.TrimSpace(task.Message) != "" {
			record.ExecutionMessage = task.Message
		}
		break
	}

	return &record
}

func buildResourceActionRequest(payload map[string]any) providerService.ResourceActionRequest {
	return providerService.ResourceActionRequest{
		Count:      payloadInt(payload, "count"),
		SizeGB:     payloadInt(payload, "sizeGb", "sizeGB"),
		StoreID:    payloadInt(payload, "storeId"),
		IPGroup:    payloadInt(payload, "ipGroup"),
		Driver:     payloadString(payload, "driver"),
		Name:       payloadString(payload, "name"),
		DiskID:     payloadString(payload, "diskId"),
		SnapshotID: payloadString(payload, "snapshotId"),
		BackupID:   payloadString(payload, "backupId"),
	}
}

func payloadString(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := payload[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case string:
			trimmed := strings.TrimSpace(typed)
			if trimmed != "" {
				return trimmed
			}
		default:
			text := strings.TrimSpace(fmt.Sprint(typed))
			if text != "" && text != "<nil>" {
				return text
			}
		}
	}
	return ""
}

func payloadInt(payload map[string]any, keys ...string) int {
	for _, key := range keys {
		value, ok := payload[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case int:
			return typed
		case int8:
			return int(typed)
		case int16:
			return int(typed)
		case int32:
			return int(typed)
		case int64:
			return int(typed)
		case float32:
			return int(typed)
		case float64:
			return int(typed)
		case string:
			if parsed, err := strconv.Atoi(strings.TrimSpace(typed)); err == nil {
				return parsed
			}
		}
	}
	return 0
}

func valueOrZero(serviceRecord *domain.ServiceRecord) int64 {
	if serviceRecord == nil {
		return 0
	}
	return serviceRecord.ID
}

func valueOrZeroInvoice(invoice *domain.Invoice) int64 {
	if invoice == nil {
		return 0
	}
	return invoice.ID
}

func valueOrZeroOrder(order *domain.Order) int64 {
	if order == nil {
		return 0
	}
	return order.ID
}

func valueOrZeroInvoiceDueAt(invoice *domain.Invoice) string {
	if invoice == nil {
		return ""
	}
	return invoice.DueAt
}

func firstInvoiceDueAt(items []domain.Invoice) string {
	if len(items) == 0 {
		return ""
	}
	return items[0].DueAt
}

func valueOrZeroOrderStatus(order *domain.Order) string {
	if order == nil {
		return ""
	}
	return string(order.Status)
}

func valueOrZeroInvoiceStatus(invoice *domain.Invoice) string {
	if invoice == nil {
		return ""
	}
	return string(invoice.Status)
}

func serviceNoOrEmpty(serviceRecord *domain.ServiceRecord) string {
	if serviceRecord == nil {
		return ""
	}
	return serviceRecord.ServiceNo
}

func pickPrice(product catalogDomain.Product, cycleCode string) (catalogDomain.PriceOption, error) {
	for _, item := range product.Pricing {
		if item.CycleCode == cycleCode {
			return item, nil
		}
	}
	return catalogDomain.PriceOption{}, fmt.Errorf("未找到对应计费周期")
}

func buildSelectedConfiguration(product catalogDomain.Product, selectedOptions []dto.CheckoutOptionRequest) ([]domain.ServiceConfigSelection, error) {
	selectedMap := make(map[string]string, len(selectedOptions))
	for _, item := range selectedOptions {
		code := strings.TrimSpace(item.Code)
		value := strings.TrimSpace(item.Value)
		if code == "" || value == "" {
			continue
		}
		selectedMap[code] = value
	}

	result := make([]domain.ServiceConfigSelection, 0, len(product.ConfigOptions))
	for _, item := range product.ConfigOptions {
		value := item.DefaultValue
		if selectedValue, ok := selectedMap[item.Code]; ok {
			value = selectedValue
		}
		if item.Required && value == "" {
			return nil, fmt.Errorf("配置项 %s 为必选项", item.Name)
		}

		label := value
		priceDelta := 0.0
		matchedChoice := false
		for _, choice := range item.Choices {
			if choice.Value == value {
				label = choice.Label
				priceDelta = choice.PriceDelta
				matchedChoice = true
				break
			}
		}
		if len(item.Choices) > 0 && item.InputType != "text" && value != "" && !matchedChoice {
			return nil, fmt.Errorf("配置项 %s 的值无效", item.Name)
		}
		result = append(result, domain.ServiceConfigSelection{
			Code:       item.Code,
			Name:       item.Name,
			Value:      value,
			ValueLabel: label,
			PriceDelta: priceDelta,
		})
	}
	return result, nil
}

func configurationDelta(items []domain.ServiceConfigSelection) float64 {
	total := 0.0
	for _, item := range items {
		total += item.PriceDelta
	}
	return total
}

func buildResourceSnapshot(product catalogDomain.Product, customerID int64, configuration []domain.ServiceConfigSelection) domain.ServiceResourceSnapshot {
	template := product.ResourceTemplate
	snapshot := domain.ServiceResourceSnapshot{
		RegionName:      template.RegionName,
		ZoneName:        template.ZoneName,
		Hostname:        fmt.Sprintf("c%d-%s", customerID, slugify(product.Name)),
		OperatingSystem: template.OperatingSystem,
		LoginUsername:   template.LoginUsername,
		PasswordHint:    "初始密码已下发到站内信",
		SecurityGroup:   template.SecurityGroup,
		CPUCores:        template.CPUCores,
		MemoryGB:        template.MemoryGB,
		SystemDiskGB:    template.SystemDiskGB,
		DataDiskGB:      template.DataDiskGB,
		BandwidthMbps:   template.BandwidthMbps,
		PublicIPv4:      "",
		PublicIPv6:      "",
	}

	for _, item := range configuration {
		switch item.Code {
		case "area":
			region := strings.TrimSpace(item.ValueLabel)
			if region == "" {
				region = strings.TrimSpace(item.Value)
			}
			if region != "" {
				snapshot.RegionName = region
				if strings.TrimSpace(snapshot.ZoneName) == "" {
					snapshot.ZoneName = region
				}
			}
		case "os":
			osName := strings.TrimSpace(item.ValueLabel)
			if osName == "" {
				osName = strings.TrimSpace(item.Value)
			}
			if osName != "" {
				snapshot.OperatingSystem = osName
			}
		case "cpu":
			if value, err := strconv.Atoi(item.Value); err == nil {
				snapshot.CPUCores = value
			}
		case "memory":
			if value, err := strconv.Atoi(item.Value); err == nil {
				snapshot.MemoryGB = value
			}
		case "system_disk_size":
			if value, err := strconv.Atoi(item.Value); err == nil {
				snapshot.SystemDiskGB = value
			}
		case "data_disk_size":
			if value, err := strconv.Atoi(item.Value); err == nil {
				snapshot.DataDiskGB = value
			}
		case "bw", "in_bw", "out_bw":
			if value, err := strconv.Atoi(item.Value); err == nil {
				snapshot.BandwidthMbps = value
			}
		}
	}

	return snapshot
}

func summarizeInvoices(items []domain.Invoice) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, fmt.Sprintf("%s:%s", item.InvoiceNo, item.Status))
	}
	return result
}

func summarizeServices(items []domain.ServiceRecord) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, fmt.Sprintf("%s:%s", item.ServiceNo, item.Status))
	}
	return result
}

func loadOrderIfExists(repo repository.Repository, orderID int64) *domain.Order {
	if orderID <= 0 {
		return nil
	}
	if item, ok := repo.GetOrderByID(orderID); ok {
		return &item
	}
	return nil
}

func loadInvoiceIfExists(repo repository.Repository, invoiceID int64) *domain.Invoice {
	if invoiceID <= 0 {
		return nil
	}
	if item, ok := repo.GetInvoiceByID(invoiceID); ok {
		return &item
	}
	return nil
}

func resolveProductAutomationType(product catalogDomain.Product) string {
	switch strings.ToUpper(strings.TrimSpace(product.AutomationConfig.Channel)) {
	case "MOFANG_CLOUD":
		return "MOFANG_CLOUD"
	case "ZJMF_API":
		return "ZJMF_API"
	case "RESOURCE":
		return "RESOURCE"
	case "MANUAL":
		return "MANUAL"
	}

	switch strings.ToUpper(strings.TrimSpace(product.UpstreamMapping.ProviderType)) {
	case "MOFANG_CLOUD":
		return "MOFANG_CLOUD"
	case "ZJMF_API":
		return "ZJMF_API"
	case "RESOURCE":
		return "RESOURCE"
	case "MANUAL":
		return "MANUAL"
	default:
		return "LOCAL"
	}
}

func resolveProductProviderAccountID(product catalogDomain.Product) int64 {
	if product.AutomationConfig.ProviderAccountID > 0 {
		return product.AutomationConfig.ProviderAccountID
	}
	if product.UpstreamMapping.ProviderAccountID > 0 {
		return product.UpstreamMapping.ProviderAccountID
	}
	return 0
}

func (service *Service) finalizeProvisionAfterPayment(
	order domain.Order,
	serviceRecord *domain.ServiceRecord,
	actorType string,
	actorID int64,
	actorName,
	requestID string,
) *domain.ServiceRecord {
	if serviceRecord == nil {
		return nil
	}
	if serviceRecord.LastAction != "pending-provision" {
		return serviceRecord
	}

	taskID := service.startTask(automationService.StartTaskRequest{
		TaskType:           "AUTO_PROVISION",
		Title:              "商品自动开通",
		Channel:            firstNonEmptyString(serviceRecord.ProviderType, order.AutomationType, "LOCAL"),
		Stage:              "AFTER_PAYMENT",
		SourceType:         "order",
		SourceID:           order.ID,
		CustomerID:         order.CustomerID,
		CustomerName:       order.CustomerName,
		ProductName:        order.ProductName,
		OrderID:            order.ID,
		InvoiceID:          serviceRecord.InvoiceID,
		ServiceID:          serviceRecord.ID,
		ServiceNo:          serviceRecord.ServiceNo,
		ProviderType:       serviceRecord.ProviderType,
		ProviderResourceID: serviceRecord.ProviderResourceID,
		ActionName:         "provision",
		OperatorType:       actorType,
		OperatorName:       actorName,
		RequestPayload: map[string]any{
			"billingCycle":     order.BillingCycle,
			"automationType":   order.AutomationType,
			"providerType":     serviceRecord.ProviderType,
			"providerResource": serviceRecord.ProviderResourceID,
		},
		Message: "自动开通任务已进入执行阶段",
	})

	product, exists := service.catalog.GetByID(order.ProductID)
	if !exists {
		service.blockTask(taskID, "缺少商品定义，无法执行自动开通", map[string]any{
			"orderId":   order.ID,
			"serviceId": serviceRecord.ID,
		})
		return service.markProvisionBlocked(serviceRecord, "缺少商品定义，无法执行自动开通", actorType, actorID, actorName, requestID)
	}
	if !product.AutomationConfig.AutoProvision || strings.ToUpper(strings.TrimSpace(product.AutomationConfig.ProvisionStage)) != "AFTER_PAYMENT" {
		service.blockTask(taskID, "商品未启用付款后自动开通", map[string]any{
			"productId":  product.ID,
			"serviceId":  serviceRecord.ID,
			"autoEnable": product.AutomationConfig.AutoProvision,
			"stage":      product.AutomationConfig.ProvisionStage,
		})
		return service.markProvisionBlocked(serviceRecord, "商品未启用付款后自动开通", actorType, actorID, actorName, requestID)
	}

	switch resolveProductAutomationType(product) {
	case "MOFANG_CLOUD":
		if service.provider == nil {
			service.blockTask(taskID, "未初始化魔方云提供方服务", map[string]any{
				"serviceId": serviceRecord.ID,
				"productId": product.ID,
			})
			return service.markProvisionBlocked(serviceRecord, "未初始化魔方云提供方服务", actorType, actorID, actorName, requestID)
		}

		result, err := service.provider.ProvisionMofangCloud(providerService.MofangProvisionRequest{
			ServiceID:        serviceRecord.ID,
			CustomerID:       order.CustomerID,
			ServiceNo:        serviceRecord.ServiceNo,
			CustomerName:     order.CustomerName,
			ProductName:      order.ProductName,
			ProductType:      order.ProductType,
			BillingCycle:     order.BillingCycle,
			ProviderNode:     product.AutomationConfig.ProviderNode,
			ServerGroup:      product.AutomationConfig.ServerGroup,
			Configuration:    serviceRecord.Configuration,
			ResourceSnapshot: serviceRecord.ResourceSnapshot,
		}, order.ProviderAccountID)
		if err != nil {
			service.failTask(taskID, err.Error(), map[string]any{
				"serviceId": serviceRecord.ID,
				"productId": product.ID,
			})
			return service.markProvisionBlocked(serviceRecord, err.Error(), actorType, actorID, actorName, requestID)
		}

		updated, ok := service.repository.UpdateServiceProvisioning(serviceRecord.ID, domain.ServiceProvisioningUpdate{
			ProviderAccountID:  order.ProviderAccountID,
			ProviderType:       result.ProviderType,
			ProviderResourceID: result.ProviderResourceID,
			RegionName:         result.RegionName,
			IPAddress:          result.IPAddress,
			Status:             result.Status,
			SyncStatus:         result.SyncStatus,
			SyncMessage:        result.SyncMessage,
			LastAction:         result.LastAction,
			ResourceSnapshot:   result.ResourceSnapshot,
			Configuration:      serviceRecord.Configuration,
		})
		if !ok {
			service.failTask(taskID, "本地服务写回失败", map[string]any{
				"serviceId":          serviceRecord.ID,
				"providerResourceId": result.ProviderResourceID,
			})
			return service.markProvisionBlocked(serviceRecord, "本地服务写回失败", actorType, actorID, actorName, requestID)
		}

		if strings.TrimSpace(updated.ProviderResourceID) != "" && service.provider != nil {
			if _, err := service.provider.SyncServiceByID(updated.ID, true); err == nil {
				if synced, exists := service.repository.GetServiceByID(updated.ID); exists {
					updated = synced
				}
			}
		}

		service.successTask(taskID, "魔方云自动开通已完成", map[string]any{
			"providerType":       updated.ProviderType,
			"providerResourceId": updated.ProviderResourceID,
			"syncStatus":         updated.SyncStatus,
			"status":             updated.Status,
		})

		service.audit.Record(audit.Entry{
			ActorType:   actorType,
			ActorID:     actorID,
			Actor:       actorName,
			Action:      "service.provision.success",
			TargetType:  "service",
			TargetID:    updated.ID,
			Target:      updated.ServiceNo,
			RequestID:   requestID,
			Description: "按魔方云自动化渠道完成实例开通回写",
			Payload: map[string]any{
				"providerType":       updated.ProviderType,
				"providerResourceId": updated.ProviderResourceID,
				"syncStatus":         updated.SyncStatus,
				"status":             updated.Status,
			},
		})
		return &updated

	case "ZJMF_API":
		service.blockTask(taskID, "上下游财务自动开通链路待接通", map[string]any{
			"serviceId": serviceRecord.ID,
			"productId": product.ID,
		})
		return service.markProvisionBlocked(serviceRecord, "上下游财务自动开通链路待接通", actorType, actorID, actorName, requestID)
	case "RESOURCE", "MANUAL":
		service.blockTask(taskID, "当前商品走资源池或人工交付，暂未接入自动开通", map[string]any{
			"serviceId": serviceRecord.ID,
			"productId": product.ID,
		})
		return service.markProvisionBlocked(serviceRecord, "当前商品走资源池或人工交付，暂未接入自动开通", actorType, actorID, actorName, requestID)
	default:
		service.successTask(taskID, "当前商品无需自动开通任务", map[string]any{
			"serviceId": serviceRecord.ID,
			"status":    serviceRecord.Status,
		})
		return serviceRecord
	}
}

func (service *Service) markProvisionBlocked(
	serviceRecord *domain.ServiceRecord,
	message,
	actorType string,
	actorID int64,
	actorName,
	requestID string,
) *domain.ServiceRecord {
	updated, ok := service.repository.UpdateServiceProvisioning(serviceRecord.ID, domain.ServiceProvisioningUpdate{
		ProviderType:       serviceRecord.ProviderType,
		ProviderResourceID: serviceRecord.ProviderResourceID,
		RegionName:         serviceRecord.RegionName,
		IPAddress:          serviceRecord.IPAddress,
		Status:             domain.ServiceStatusPending,
		SyncStatus:         "FAILED",
		SyncMessage:        message,
		LastAction:         "provision-blocked",
		ResourceSnapshot:   serviceRecord.ResourceSnapshot,
		Configuration:      serviceRecord.Configuration,
	})
	target := serviceRecord
	if ok {
		target = &updated
	}

	service.audit.Record(audit.Entry{
		ActorType:   actorType,
		ActorID:     actorID,
		Actor:       actorName,
		Action:      "service.provision.failed",
		TargetType:  "service",
		TargetID:    target.ID,
		Target:      target.ServiceNo,
		RequestID:   requestID,
		Description: "鑷姩寮€閫氭湭鎵ц鎴愬姛",
		Payload: map[string]any{
			"message": message,
		},
	})
	return target
}

func normalizeOrderRequestType(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case string(domain.OrderRequestTypeCancel):
		return string(domain.OrderRequestTypeCancel)
	case string(domain.OrderRequestTypeRenew):
		return string(domain.OrderRequestTypeRenew)
	case string(domain.OrderRequestTypePriceAdjust):
		return string(domain.OrderRequestTypePriceAdjust)
	default:
		return ""
	}
}

func normalizeOrderRequestStatus(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case string(domain.OrderRequestStatusPending):
		return string(domain.OrderRequestStatusPending)
	case string(domain.OrderRequestStatusApproved):
		return string(domain.OrderRequestStatusApproved)
	case string(domain.OrderRequestStatusRejected):
		return string(domain.OrderRequestStatusRejected)
	case string(domain.OrderRequestStatusCompleted):
		return string(domain.OrderRequestStatusCompleted)
	case string(domain.OrderRequestStatusCancelled):
		return string(domain.OrderRequestStatusCancelled)
	default:
		return ""
	}
}

func defaultOrderRequestSummary(requestType, productName string) string {
	productName = strings.TrimSpace(productName)
	if productName == "" {
		productName = "当前订单"
	}
	switch requestType {
	case string(domain.OrderRequestTypeCancel):
		return fmt.Sprintf("%s取消申请", productName)
	case string(domain.OrderRequestTypeRenew):
		return fmt.Sprintf("%s续费请求", productName)
	case string(domain.OrderRequestTypePriceAdjust):
		return fmt.Sprintf("%s改价申请", productName)
	default:
		return fmt.Sprintf("%s业务申请", productName)
	}
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(" ", "-", "。", "-", "_", "-", "/", "-", ".", "-", ":", "-", "\\", "-")
	value = replacer.Replace(value)

	var builder strings.Builder
	lastDash := false
	for _, char := range value {
		switch {
		case char >= 'a' && char <= 'z':
			builder.WriteRune(char)
			lastDash = false
		case char >= '0' && char <= '9':
			builder.WriteRune(char)
			lastDash = false
		case char == '-':
			if !lastDash {
				builder.WriteRune(char)
				lastDash = true
			}
		}
	}

	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return "service"
	}
	return result
}

func serviceActionDescription(action string) string {
	switch action {
	case "activate":
		return "后台恢复服务运行"
	case "suspend":
		return "后台暂停服务"
	case "terminate":
		return "后台终止服务"
	case "reboot":
		return "后台发起实例重启"
	case "reset-password":
		return "后台重置实例密码"
	case "reinstall":
		return "后台发起实例重装"
	default:
		return "后台更新服务状态"
	}
}
