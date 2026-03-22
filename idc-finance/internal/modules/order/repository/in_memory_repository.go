package repository

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"strings"
	"sync"
	"time"

	"idc-finance/internal/modules/order/domain"
)

type MemoryRepository struct {
	mu             sync.RWMutex
	nextOrderID    int64
	nextInvoiceID  int64
	nextServiceID  int64
	nextPaymentID  int64
	nextRefundID   int64
	orders         []domain.Order
	invoices       []domain.Invoice
	services       []domain.ServiceRecord
	payments       []domain.PaymentRecord
	refunds        []domain.RefundRecord
	serviceChanges []memoryServiceChangeLink
}

type memoryServiceChangeLink struct {
	ID           int64
	ServiceID    int64
	OrderID      int64
	InvoiceID    int64
	ActionName   string
	Title        string
	Status       string
	Reason       string
	BillingCycle string
	Amount       float64
	Payload      map[string]any
	PaidAt       string
	RefundedAt   string
	CreatedAt    string
	UpdatedAt    string
}

func NewMemoryRepository() *MemoryRepository {
	now := time.Now()

	orderOneSnapshot := domain.ServiceResourceSnapshot{
		RegionName:      "华南广州",
		ZoneName:        "广州三区",
		Hostname:        "srv-00000001.idc.local",
		OperatingSystem: "Rocky Linux 9",
		LoginUsername:   "root",
		PasswordHint:    "初始化密码已下发到站内信",
		SecurityGroup:   "default-cloud",
		CPUCores:        4,
		MemoryGB:        8,
		SystemDiskGB:    60,
		DataDiskGB:      120,
		BandwidthMbps:   20,
		PublicIPv4:      "203.0.113.11",
		PublicIPv6:      "2408:4004:10::11",
	}

	orderTwoSnapshot := domain.ServiceResourceSnapshot{
		RegionName:      "华南广州",
		ZoneName:        "骨干网络",
		Hostname:        "network-addon-0002",
		OperatingSystem: "-",
		LoginUsername:   "-",
		PasswordHint:    "-",
		SecurityGroup:   "network-edge",
		BandwidthMbps:   50,
		PublicIPv4:      "203.0.113.51",
	}

	return &MemoryRepository{
		nextOrderID:   3,
		nextInvoiceID: 3,
		nextServiceID: 2,
		nextPaymentID: 2,
		nextRefundID:  1,
		orders: []domain.Order{
			{
				ID:                1,
				OrderNo:           "ORD-20260320-0001",
				CustomerID:        1,
				CustomerName:      "演示客户",
				ProductID:         1,
				ProductName:       "弹性云主机 CN2 标准型",
				ProductType:       "CLOUD",
				AutomationType:    "MOFANG_CLOUD",
				ProviderAccountID: 1,
				BillingCycle:      "annual",
				Amount:            1999,
				Status:            domain.OrderStatusActive,
				Configuration: []domain.ServiceConfigSelection{
					{Code: "cpu", Name: "CPU 规格", Value: "4", ValueLabel: "4 核", PriceDelta: 80},
					{Code: "memory", Name: "内存规格", Value: "8", ValueLabel: "8 GB", PriceDelta: 60},
					{Code: "backup", Name: "云备份", Value: "enabled", ValueLabel: "启用", PriceDelta: 30},
				},
				ResourceSnapshot: orderOneSnapshot,
				CreatedAt:        now.Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
			},
			{
				ID:                2,
				OrderNo:           "ORD-20260320-0002",
				CustomerID:        1,
				CustomerName:      "演示客户",
				ProductID:         2,
				ProductName:       "精品带宽 50M",
				ProductType:       "BANDWIDTH",
				AutomationType:    "LOCAL",
				ProviderAccountID: 0,
				BillingCycle:      "monthly",
				Amount:            299,
				Status:            domain.OrderStatusPending,
				Configuration: []domain.ServiceConfigSelection{
					{Code: "commit", Name: "承诺带宽", Value: "50", ValueLabel: "50 Mbps", PriceDelta: 0},
				},
				ResourceSnapshot: orderTwoSnapshot,
				CreatedAt:        now.Add(-2 * time.Hour).Format("2006-01-02 15:04:05"),
			},
		},
		invoices: []domain.Invoice{
			{
				ID:           1,
				InvoiceNo:    "INV-20260320-0001",
				OrderID:      1,
				OrderNo:      "ORD-20260320-0001",
				CustomerID:   1,
				ProductName:  "弹性云主机 CN2 标准型",
				TotalAmount:  1999,
				Status:       domain.InvoiceStatusPaid,
				DueAt:        now.Add(-23 * time.Hour).Format("2006-01-02 15:04:05"),
				PaidAt:       now.Add(-23 * time.Hour).Format("2006-01-02 15:04:05"),
				BillingCycle: "annual",
			},
			{
				ID:           2,
				InvoiceNo:    "INV-20260320-0002",
				OrderID:      2,
				OrderNo:      "ORD-20260320-0002",
				CustomerID:   1,
				ProductName:  "精品带宽 50M",
				TotalAmount:  299,
				Status:       domain.InvoiceStatusUnpaid,
				DueAt:        now.Add(48 * time.Hour).Format("2006-01-02 15:04:05"),
				BillingCycle: "monthly",
			},
		},
		services: []domain.ServiceRecord{
			{
				ID:                1,
				ServiceNo:         "SRV-20260320-0001",
				OrderID:           1,
				InvoiceID:         1,
				CustomerID:        1,
				ProductName:       "弹性云主机 CN2 标准型",
				ProviderType:      "MOFANG_CLOUD",
				ProviderAccountID: 1,
				Status:            domain.ServiceStatusActive,
				NextDueAt:         now.AddDate(1, 0, 0).Format("2006-01-02"),
				LastAction:        "activate",
				UpdatedAt:         now.Add(-23 * time.Hour).Format("2006-01-02 15:04:05"),
				Configuration: []domain.ServiceConfigSelection{
					{Code: "cpu", Name: "CPU 规格", Value: "4", ValueLabel: "4 核", PriceDelta: 80},
					{Code: "memory", Name: "内存规格", Value: "8", ValueLabel: "8 GB", PriceDelta: 60},
					{Code: "backup", Name: "云备份", Value: "enabled", ValueLabel: "启用", PriceDelta: 30},
				},
				ResourceSnapshot: orderOneSnapshot,
				CreatedAt:        now.Add(-23 * time.Hour).Format("2006-01-02 15:04:05"),
			},
		},
		payments: []domain.PaymentRecord{
			{
				ID:         1,
				PaymentNo:  "PAY-20260320-0001",
				InvoiceID:  1,
				OrderID:    1,
				CustomerID: 1,
				Channel:    "ONLINE",
				TradeNo:    "TRADE-20260320-0001",
				Amount:     1999,
				Source:     "PORTAL",
				Status:     domain.PaymentStatusCompleted,
				Operator:   "演示客户",
				PaidAt:     now.Add(-23 * time.Hour).Format("2006-01-02 15:04:05"),
			},
		},
		refunds:        []domain.RefundRecord{},
		serviceChanges: []memoryServiceChangeLink{},
	}
}

func (repository *MemoryRepository) ListOrders(filter domain.OrderListFilter) ([]domain.Order, int) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := repository.decorateOrdersLocked(cloneOrders(repository.orders))
	items = filterMemoryOrders(items, filter)
	total := len(items)
	items = paginateMemoryOrders(items, filter.Page, filter.Limit)
	return items, total
}

func (repository *MemoryRepository) ListOrdersByCustomer(customerID int64) []domain.Order {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.Order, 0)
	for _, item := range repository.orders {
		if item.CustomerID == customerID {
			items = append(items, cloneOrder(item))
		}
	}
	return repository.decorateOrdersLocked(items)
}

func (repository *MemoryRepository) GetOrderByID(id int64) (domain.Order, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	for _, item := range repository.orders {
		if item.ID == id {
			return repository.decorateOrderLocked(cloneOrder(item)), true
		}
	}
	return domain.Order{}, false
}

func (repository *MemoryRepository) ListInvoices(filter domain.InvoiceListFilter) ([]domain.Invoice, int) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	items := slices.Clone(repository.invoices)
	items = filterMemoryInvoices(items, filter)
	total := len(items)
	items = paginateMemoryInvoices(items, filter.Page, filter.Limit)
	return items, total
}

func (repository *MemoryRepository) ListInvoicesByCustomer(customerID int64) []domain.Invoice {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.Invoice, 0)
	for _, item := range repository.invoices {
		if item.CustomerID == customerID {
			items = append(items, item)
		}
	}
	return items
}

func (repository *MemoryRepository) ListInvoicesByOrder(orderID int64) []domain.Invoice {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.Invoice, 0)
	for _, item := range repository.invoices {
		if item.OrderID == orderID {
			items = append(items, item)
		}
	}
	return items
}

func (repository *MemoryRepository) GetInvoiceByID(id int64) (domain.Invoice, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	for _, item := range repository.invoices {
		if item.ID == id {
			return item, true
		}
	}
	return domain.Invoice{}, false
}

func (repository *MemoryRepository) GetServiceChangeOrderByInvoiceID(invoiceID int64) (domain.ServiceChangeOrder, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	link, ok := repository.findServiceChangeByInvoiceLocked(invoiceID)
	if !ok {
		return domain.ServiceChangeOrder{}, false
	}

	return cloneServiceChange(link), true
}

func (repository *MemoryRepository) GetServiceChangeOrderByOrderID(orderID int64) (domain.ServiceChangeOrder, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	link, ok := repository.findServiceChangeByOrderLocked(orderID)
	if !ok {
		return domain.ServiceChangeOrder{}, false
	}

	return cloneServiceChange(link), true
}

func (repository *MemoryRepository) ListServiceChangeOrdersByService(serviceID int64) []domain.ServiceChangeOrder {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.ServiceChangeOrder, 0)
	for _, item := range repository.serviceChanges {
		if item.ServiceID != serviceID {
			continue
		}
		items = append(items, cloneServiceChange(item))
	}
	return items
}

func (repository *MemoryRepository) UpdatePendingOrder(
	orderID int64,
	update domain.PendingOrderUpdate,
) (domain.Order, *domain.Invoice, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for orderIndex, order := range repository.orders {
		if order.ID != orderID {
			continue
		}

		var linkedInvoice *domain.Invoice
		hasPaidInvoice := false
		_, isServiceChange := repository.findServiceChangeByOrderLocked(order.ID)
		nextStatus := string(order.Status)
		if strings.TrimSpace(update.Status) != "" {
			nextStatus = strings.ToUpper(strings.TrimSpace(update.Status))
			switch nextStatus {
			case string(domain.OrderStatusPending), string(domain.OrderStatusActive), string(domain.OrderStatusCompleted), string(domain.OrderStatusCancelled):
			default:
				return domain.Order{}, nil, false, fmt.Errorf("订单状态不合法")
			}
		}
		for invoiceIndex, invoice := range repository.invoices {
			if invoice.OrderID != order.ID {
				continue
			}
			if strings.TrimSpace(update.ProductName) != "" {
				order.ProductName = strings.TrimSpace(update.ProductName)
				invoice.ProductName = order.ProductName
			}
			if strings.TrimSpace(update.BillingCycle) != "" {
				order.BillingCycle = strings.TrimSpace(update.BillingCycle)
				invoice.BillingCycle = order.BillingCycle
			}
			if update.Amount != nil {
				if *update.Amount < 0 {
					return domain.Order{}, nil, false, fmt.Errorf("订单金额不能小于 0")
				}
				order.Amount = *update.Amount
				invoice.TotalAmount = *update.Amount
			}
			if strings.TrimSpace(update.DueAt) != "" {
				parsed, ok := parseOrderTime(strings.TrimSpace(update.DueAt))
				if !ok {
					return domain.Order{}, nil, false, fmt.Errorf("账单到期时间格式不正确")
				}
				invoice.DueAt = parsed.Format("2006-01-02 15:04:05")
			}
			repository.invoices[invoiceIndex] = invoice
			invoiceCopy := invoice
			linkedInvoice = &invoiceCopy
			if invoice.Status == domain.InvoiceStatusPaid {
				hasPaidInvoice = true
			}
			break
		}
		order.Status = domain.OrderStatus(nextStatus)
		now := time.Now().Format("2006-01-02 15:04:05")
		if !isServiceChange {
			for serviceIndex, item := range repository.services {
				if item.OrderID != order.ID {
					continue
				}
				switch order.Status {
				case domain.OrderStatusCancelled:
					item.Status = domain.ServiceStatusTerminated
					item.SyncStatus = "SUCCESS"
					item.SyncMessage = "后台人工改签订单为已取消"
				case domain.OrderStatusPending:
					if item.Status != domain.ServiceStatusTerminated {
						item.Status = domain.ServiceStatusSuspended
					}
					item.SyncStatus = "SUCCESS"
					item.SyncMessage = "后台人工改签订单为待处理"
				case domain.OrderStatusActive, domain.OrderStatusCompleted:
					if hasPaidInvoice && item.Status != domain.ServiceStatusTerminated {
						item.Status = domain.ServiceStatusActive
						item.SyncStatus = "SUCCESS"
						item.SyncMessage = "后台人工改签订单为生效"
					}
				}
				item.LastAction = "manual-adjust"
				item.UpdatedAt = now
				repository.services[serviceIndex] = item
			}
		}

		repository.orders[orderIndex] = order
		orderCopy := repository.decorateOrderLocked(cloneOrder(order))
		return orderCopy, linkedInvoice, true, nil
	}

	return domain.Order{}, nil, false, nil
}

func (repository *MemoryRepository) UpdateUnpaidInvoice(
	invoiceID int64,
	update domain.UnpaidInvoiceUpdate,
) (domain.Invoice, *domain.Order, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for invoiceIndex, invoice := range repository.invoices {
		if invoice.ID != invoiceID {
			continue
		}

		var linkedOrder *domain.Order
		changeLink, isServiceChange := repository.findServiceChangeByInvoiceLocked(invoice.ID)
		previousStatus := invoice.Status
		nextStatus := string(invoice.Status)
		if strings.TrimSpace(update.Status) != "" {
			nextStatus = strings.ToUpper(strings.TrimSpace(update.Status))
			switch nextStatus {
			case string(domain.InvoiceStatusUnpaid), string(domain.InvoiceStatusPaid), string(domain.InvoiceStatusRefunded):
			default:
				return domain.Invoice{}, nil, false, fmt.Errorf("账单状态不合法")
			}
		}
		if strings.TrimSpace(update.ProductName) != "" {
			invoice.ProductName = strings.TrimSpace(update.ProductName)
		}
		if strings.TrimSpace(update.BillingCycle) != "" {
			invoice.BillingCycle = strings.TrimSpace(update.BillingCycle)
		}
		if update.Amount != nil {
			if *update.Amount < 0 {
				return domain.Invoice{}, nil, false, fmt.Errorf("账单金额不能小于 0")
			}
			invoice.TotalAmount = *update.Amount
		}
		if strings.TrimSpace(update.DueAt) != "" {
			parsed, ok := parseOrderTime(strings.TrimSpace(update.DueAt))
			if !ok {
				return domain.Invoice{}, nil, false, fmt.Errorf("账单到期时间格式不正确")
			}
			invoice.DueAt = parsed.Format("2006-01-02 15:04:05")
		}
		invoice.Status = domain.InvoiceStatus(nextStatus)
		if invoice.Status == domain.InvoiceStatusUnpaid {
			invoice.PaidAt = ""
		} else if strings.TrimSpace(invoice.PaidAt) == "" {
			invoice.PaidAt = time.Now().Format("2006-01-02 15:04:05")
		}
		repository.invoices[invoiceIndex] = invoice

		for orderIndex, order := range repository.orders {
			if order.ID != invoice.OrderID {
				continue
			}
			order.ProductName = invoice.ProductName
			order.BillingCycle = invoice.BillingCycle
			order.Amount = invoice.TotalAmount
			switch invoice.Status {
			case domain.InvoiceStatusPaid:
				order.Status = domain.OrderStatusActive
			case domain.InvoiceStatusUnpaid, domain.InvoiceStatusRefunded:
				order.Status = domain.OrderStatusPending
			}
			repository.orders[orderIndex] = order
			orderCopy := repository.decorateOrderLocked(cloneOrder(order))
			linkedOrder = &orderCopy
			break
		}

		now := time.Now().Format("2006-01-02 15:04:05")
		if invoice.Status == domain.InvoiceStatusPaid && previousStatus != domain.InvoiceStatusPaid {
			if _, ok := repository.findLatestPaymentLocked(invoice.ID); !ok {
				payment := domain.PaymentRecord{
					ID:         repository.nextPaymentID,
					PaymentNo:  fmt.Sprintf("PAY-%08d", repository.nextPaymentID),
					InvoiceID:  invoice.ID,
					OrderID:    invoice.OrderID,
					CustomerID: invoice.CustomerID,
					Channel:    "MANUAL",
					TradeNo:    fmt.Sprintf("TRADE-%08d", repository.nextPaymentID),
					Amount:     invoice.TotalAmount,
					Source:     "ADMIN",
					Status:     domain.PaymentStatusCompleted,
					Operator:   "后台人工改签",
					PaidAt:     now,
				}
				repository.nextPaymentID++
				repository.payments = append([]domain.PaymentRecord{payment}, repository.payments...)
			}
		}
		if invoice.Status == domain.InvoiceStatusRefunded && previousStatus != domain.InvoiceStatusRefunded {
			refund := domain.RefundRecord{
				ID:         repository.nextRefundID,
				RefundNo:   fmt.Sprintf("RFD-%08d", repository.nextRefundID),
				InvoiceID:  invoice.ID,
				OrderID:    invoice.OrderID,
				CustomerID: invoice.CustomerID,
				Amount:     invoice.TotalAmount,
				Reason:     "后台人工改签为已退款",
				Status:     domain.RefundStatusCompleted,
				CreatedAt:  now,
			}
			repository.nextRefundID++
			repository.refunds = append([]domain.RefundRecord{refund}, repository.refunds...)
		}
		if isServiceChange {
			changeLink.Status = string(invoice.Status)
			repository.replaceServiceChangeLocked(changeLink)
		} else {
			for serviceIndex, item := range repository.services {
				if item.InvoiceID != invoice.ID {
					continue
				}
				switch invoice.Status {
				case domain.InvoiceStatusPaid:
					item.Status = domain.ServiceStatusActive
					item.SyncMessage = "后台人工改签账单为已支付"
				case domain.InvoiceStatusUnpaid:
					if item.Status != domain.ServiceStatusTerminated {
						item.Status = domain.ServiceStatusSuspended
					}
					item.SyncMessage = "后台人工改签账单为未支付"
				case domain.InvoiceStatusRefunded:
					item.Status = domain.ServiceStatusTerminated
					item.SyncMessage = "后台人工改签账单为已退款"
				}
				item.SyncStatus = "SUCCESS"
				item.LastAction = "manual-adjust"
				item.UpdatedAt = now
				repository.services[serviceIndex] = item
			}
		}

		invoiceCopy := invoice
		return invoiceCopy, linkedOrder, true, nil
	}

	return domain.Invoice{}, nil, false, nil
}

func (repository *MemoryRepository) ListServices(filter domain.ServiceListFilter) ([]domain.ServiceRecord, int) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	items := cloneServices(repository.services)
	items = filterMemoryServices(items, filter)
	total := len(items)
	items = paginateMemoryServices(items, filter.Page, filter.Limit)
	return items, total
}

func (repository *MemoryRepository) ListServicesByCustomer(customerID int64) []domain.ServiceRecord {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.ServiceRecord, 0)
	for _, item := range repository.services {
		if item.CustomerID == customerID {
			items = append(items, cloneService(item))
		}
	}
	return items
}

func (repository *MemoryRepository) ListServicesByOrder(orderID int64) []domain.ServiceRecord {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.ServiceRecord, 0)
	seen := make(map[int64]struct{})
	for _, item := range repository.services {
		if item.OrderID == orderID {
			items = append(items, cloneService(item))
			seen[item.ID] = struct{}{}
		}
	}
	for _, link := range repository.serviceChanges {
		if link.OrderID != orderID {
			continue
		}
		for _, item := range repository.services {
			if item.ID != link.ServiceID {
				continue
			}
			if _, exists := seen[item.ID]; exists {
				break
			}
			items = append(items, cloneService(item))
			seen[item.ID] = struct{}{}
			break
		}
	}
	return items
}

func (repository *MemoryRepository) GetServiceByID(id int64) (domain.ServiceRecord, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	for _, item := range repository.services {
		if item.ID == id {
			return cloneService(item), true
		}
	}
	return domain.ServiceRecord{}, false
}

func (repository *MemoryRepository) UpdateServiceRecord(
	serviceID int64,
	update domain.ManualServiceUpdate,
) (domain.ServiceRecord, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.services {
		if item.ID != serviceID {
			continue
		}

		nextProviderType, err := normalizeProviderTypeValue(update.ProviderType, item.ProviderType)
		if err != nil {
			return domain.ServiceRecord{}, false, err
		}
		nextStatus, err := normalizeServiceStatusValue(update.Status, string(item.Status))
		if err != nil {
			return domain.ServiceRecord{}, false, err
		}
		nextSyncStatus, err := normalizeSyncStatusValue(update.SyncStatus, item.SyncStatus)
		if err != nil {
			return domain.ServiceRecord{}, false, err
		}

		if update.ProviderAccountID != nil {
			item.ProviderAccountID = *update.ProviderAccountID
		}
		item.ProviderType = nextProviderType
		item.ProviderResourceID = strings.TrimSpace(update.ProviderResourceID)
		item.RegionName = strings.TrimSpace(update.RegionName)
		item.IPAddress = strings.TrimSpace(update.IPAddress)
		item.Status = domain.ServiceStatus(nextStatus)
		item.SyncStatus = nextSyncStatus
		item.SyncMessage = strings.TrimSpace(update.SyncMessage)
		if strings.TrimSpace(update.NextDueAt) != "" {
			parsed, ok := parseOrderTime(strings.TrimSpace(update.NextDueAt))
			if !ok {
				return domain.ServiceRecord{}, false, fmt.Errorf("服务到期时间格式不正确")
			}
			item.NextDueAt = parsed.Format("2006-01-02")
		}
		item.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		repository.services[index] = item
		return cloneService(item), true, nil
	}

	return domain.ServiceRecord{}, false, nil
}

func (repository *MemoryRepository) ListPaymentsByInvoice(invoiceID int64) []domain.PaymentRecord {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.PaymentRecord, 0)
	for _, item := range repository.payments {
		if item.InvoiceID == invoiceID {
			items = append(items, item)
		}
	}
	return items
}

func (repository *MemoryRepository) ListRefundsByInvoice(invoiceID int64) []domain.RefundRecord {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.RefundRecord, 0)
	for _, item := range repository.refunds {
		if item.InvoiceID == invoiceID {
			items = append(items, item)
		}
	}
	return items
}

func (repository *MemoryRepository) Checkout(
	customerID int64,
	customerName string,
	productID int64,
	productName, productType, automationType string,
	providerAccountID int64,
	cycleCode string,
	amount float64,
	configuration []domain.ServiceConfigSelection,
	resourceSnapshot domain.ServiceResourceSnapshot,
) (domain.Order, domain.Invoice) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	now := time.Now()
	order := domain.Order{
		ID:                repository.nextOrderID,
		OrderNo:           fmt.Sprintf("ORD-%08d", repository.nextOrderID),
		CustomerID:        customerID,
		CustomerName:      customerName,
		ProductID:         productID,
		ProductName:       productName,
		ProductType:       productType,
		AutomationType:    normalizeAutomationType(automationType),
		ProviderAccountID: providerAccountID,
		BillingCycle:      cycleCode,
		Amount:            amount,
		Status:            domain.OrderStatusPending,
		Configuration:     cloneConfigSelections(configuration),
		ResourceSnapshot:  cloneSnapshot(resourceSnapshot),
		CreatedAt:         now.Format("2006-01-02 15:04:05"),
	}
	repository.nextOrderID++

	invoice := domain.Invoice{
		ID:           repository.nextInvoiceID,
		InvoiceNo:    fmt.Sprintf("INV-%08d", repository.nextInvoiceID),
		OrderID:      order.ID,
		OrderNo:      order.OrderNo,
		CustomerID:   customerID,
		ProductName:  productName,
		TotalAmount:  amount,
		Status:       domain.InvoiceStatusUnpaid,
		DueAt:        now.Add(72 * time.Hour).Format("2006-01-02 15:04:05"),
		BillingCycle: cycleCode,
	}
	repository.nextInvoiceID++

	repository.orders = append([]domain.Order{order}, repository.orders...)
	repository.invoices = append([]domain.Invoice{invoice}, repository.invoices...)
	return cloneOrder(order), invoice
}

func (repository *MemoryRepository) PayInvoice(invoiceID int64, channel, source, operator, tradeNo string) (domain.Invoice, *domain.ServiceRecord, domain.PaymentRecord, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for invoiceIndex, invoice := range repository.invoices {
		if invoice.ID != invoiceID {
			continue
		}

		if invoice.Status == domain.InvoiceStatusRefunded {
			return domain.Invoice{}, nil, domain.PaymentRecord{}, false
		}

		if invoice.Status == domain.InvoiceStatusPaid {
			var existingService *domain.ServiceRecord
			for _, serviceItem := range repository.services {
				if serviceItem.InvoiceID == invoice.ID {
					serviceCopy := cloneService(serviceItem)
					existingService = &serviceCopy
					break
				}
			}
			payment, _ := repository.findLatestPaymentLocked(invoice.ID)
			return invoice, existingService, payment, true
		}

		now := time.Now()
		invoice.Status = domain.InvoiceStatusPaid
		invoice.PaidAt = now.Format("2006-01-02 15:04:05")
		repository.invoices[invoiceIndex] = invoice

		payment := domain.PaymentRecord{
			ID:         repository.nextPaymentID,
			PaymentNo:  fmt.Sprintf("PAY-%08d", repository.nextPaymentID),
			InvoiceID:  invoice.ID,
			OrderID:    invoice.OrderID,
			CustomerID: invoice.CustomerID,
			Channel:    firstNonEmpty(channel, "ONLINE"),
			TradeNo:    firstNonEmpty(tradeNo, fmt.Sprintf("TRADE-%08d", repository.nextPaymentID)),
			Amount:     invoice.TotalAmount,
			Source:     firstNonEmpty(source, "PORTAL"),
			Status:     domain.PaymentStatusCompleted,
			Operator:   firstNonEmpty(operator, "系统"),
			PaidAt:     now.Format("2006-01-02 15:04:05"),
		}
		repository.nextPaymentID++
		repository.payments = append([]domain.PaymentRecord{payment}, repository.payments...)
		changeLink, isServiceChange := repository.findServiceChangeByInvoiceLocked(invoice.ID)

		for orderIndex, order := range repository.orders {
			if order.ID != invoice.OrderID {
				continue
			}

			order.Status = domain.OrderStatusActive
			repository.orders[orderIndex] = order
			if isServiceChange {
				changeLink.Status = string(domain.InvoiceStatusPaid)
				changeLink.PaidAt = now.Format("2006-01-02 15:04:05")
				changeLink.UpdatedAt = now.Format("2006-01-02 15:04:05")
				repository.replaceServiceChangeLocked(changeLink)
				for _, serviceItem := range repository.services {
					if serviceItem.ID != changeLink.ServiceID {
						continue
					}
					serviceCopy := cloneService(serviceItem)
					return invoice, &serviceCopy, payment, true
				}
				return invoice, nil, payment, true
			}

			for serviceIndex, existing := range repository.services {
				if existing.InvoiceID != invoice.ID {
					continue
				}
				existing.Status = domain.ServiceStatusActive
				existing.LastAction = "activate"
				existing.UpdatedAt = now.Format("2006-01-02 15:04:05")
				repository.services[serviceIndex] = existing
				serviceCopy := cloneService(existing)
				return invoice, &serviceCopy, payment, true
			}

			snapshot := cloneSnapshot(order.ResourceSnapshot)
			if snapshot.Hostname == "" {
				snapshot.Hostname = fmt.Sprintf("srv-%08d.idc.local", repository.nextServiceID)
			}
			if snapshot.PublicIPv4 == "" {
				snapshot.PublicIPv4 = fmt.Sprintf("203.0.113.%d", 100+repository.nextServiceID)
			}
			if snapshot.PasswordHint == "" {
				snapshot.PasswordHint = "初始化密码已下发到站内信"
			}

			serviceStatus := domain.ServiceStatusActive
			syncStatus := "SUCCESS"
			syncMessage := ""
			lastAction := "activate"
			providerType := normalizeAutomationType(order.AutomationType)
			if providerType != "LOCAL" {
				serviceStatus = domain.ServiceStatusPending
				syncStatus = "PENDING"
				syncMessage = fmt.Sprintf("待按自动化渠道 %s 开通", providerType)
				lastAction = "pending-provision"
			}

			service := domain.ServiceRecord{
				ID:                repository.nextServiceID,
				ServiceNo:         fmt.Sprintf("SRV-%08d", repository.nextServiceID),
				OrderID:           order.ID,
				InvoiceID:         invoice.ID,
				CustomerID:        order.CustomerID,
				ProductName:       order.ProductName,
				ProviderType:      providerType,
				ProviderAccountID: order.ProviderAccountID,
				Status:            serviceStatus,
				SyncStatus:        syncStatus,
				SyncMessage:       syncMessage,
				NextDueAt:         nextDueDate(order.BillingCycle, now),
				LastAction:        lastAction,
				UpdatedAt:         now.Format("2006-01-02 15:04:05"),
				Configuration:     cloneConfigSelections(order.Configuration),
				ResourceSnapshot:  snapshot,
				CreatedAt:         now.Format("2006-01-02 15:04:05"),
			}
			repository.nextServiceID++
			repository.services = append([]domain.ServiceRecord{service}, repository.services...)
			serviceCopy := cloneService(service)
			return invoice, &serviceCopy, payment, true
		}

		return invoice, nil, payment, true
	}

	return domain.Invoice{}, nil, domain.PaymentRecord{}, false
}

func (repository *MemoryRepository) ExecuteServiceAction(serviceID int64, action string, params domain.ServiceActionParams) (domain.ServiceRecord, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05")
	for index, item := range repository.services {
		if item.ID != serviceID {
			continue
		}

		switch action {
		case "activate":
			item.Status = domain.ServiceStatusActive
		case "suspend":
			item.Status = domain.ServiceStatusSuspended
		case "terminate":
			item.Status = domain.ServiceStatusTerminated
		case "reboot":
			if item.Status == domain.ServiceStatusSuspended || item.Status == domain.ServiceStatusTerminated {
				return domain.ServiceRecord{}, false
			}
		case "reset-password":
			if strings.TrimSpace(params.Password) != "" {
				item.ResourceSnapshot.PasswordHint = fmt.Sprintf("最近一次密码重置：%s", now)
			} else {
				item.ResourceSnapshot.PasswordHint = fmt.Sprintf("最近一次密码重置：%s（系统生成密码）", now)
			}
		case "reinstall":
			imageName := strings.TrimSpace(params.ImageName)
			if imageName == "" {
				imageName = "Rocky Linux 9"
			}
			item.ResourceSnapshot.OperatingSystem = imageName
			item.Status = domain.ServiceStatusActive
		default:
			return domain.ServiceRecord{}, false
		}

		item.LastAction = action
		item.UpdatedAt = now
		repository.services[index] = item
		return cloneService(item), true
	}

	return domain.ServiceRecord{}, false
}

func (repository *MemoryRepository) UpdateServiceProvisioning(serviceID int64, update domain.ServiceProvisioningUpdate) (domain.ServiceRecord, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05")
	for index, item := range repository.services {
		if item.ID != serviceID {
			continue
		}

		if strings.TrimSpace(update.ProviderType) != "" {
			item.ProviderType = update.ProviderType
		}
		if update.ProviderAccountID > 0 {
			item.ProviderAccountID = update.ProviderAccountID
		}
		item.ProviderResourceID = strings.TrimSpace(update.ProviderResourceID)
		item.RegionName = strings.TrimSpace(update.RegionName)
		item.IPAddress = strings.TrimSpace(update.IPAddress)
		item.Status = update.Status
		item.SyncStatus = strings.TrimSpace(update.SyncStatus)
		item.SyncMessage = strings.TrimSpace(update.SyncMessage)
		item.LastAction = strings.TrimSpace(update.LastAction)
		item.LastSyncAt = now
		item.UpdatedAt = now
		item.Configuration = cloneConfigSelections(update.Configuration)
		item.ResourceSnapshot = cloneSnapshot(update.ResourceSnapshot)
		repository.services[index] = item
		return cloneService(item), true
	}

	return domain.ServiceRecord{}, false
}

func (repository *MemoryRepository) CreateServiceChangeOrder(
	serviceID int64,
	input domain.ServiceChangeOrderInput,
) (domain.Order, domain.Invoice, *domain.ServiceRecord, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	var (
		serviceRecord domain.ServiceRecord
		foundService  bool
		baseOrder     domain.Order
		foundOrder    bool
	)
	for _, item := range repository.services {
		if item.ID == serviceID {
			serviceRecord = cloneService(item)
			foundService = true
			break
		}
	}
	if !foundService {
		return domain.Order{}, domain.Invoice{}, nil, false, nil
	}
	for _, item := range repository.orders {
		if item.ID == serviceRecord.OrderID {
			baseOrder = cloneOrder(item)
			foundOrder = true
			break
		}
	}
	if !foundOrder {
		baseOrder = domain.Order{
			CustomerID:        serviceRecord.CustomerID,
			ProductType:       "SERVICE_CHANGE",
			ProviderAccountID: serviceRecord.ProviderAccountID,
			BillingCycle:      "monthly",
		}
	}

	billingCycle := firstNonEmpty(strings.TrimSpace(input.BillingCycle), firstNonEmpty(strings.TrimSpace(baseOrder.BillingCycle), "monthly"))
	title := firstNonEmpty(strings.TrimSpace(input.Title), serviceRecord.ProductName+" 改配单")
	now := time.Now()

	order := domain.Order{
		ID:                repository.nextOrderID,
		OrderNo:           fmt.Sprintf("ORD-%08d", repository.nextOrderID),
		CustomerID:        serviceRecord.CustomerID,
		CustomerName:      baseOrder.CustomerName,
		ProductID:         baseOrder.ProductID,
		ProductName:       title,
		ProductType:       firstNonEmpty(baseOrder.ProductType, "SERVICE_CHANGE"),
		AutomationType:    "SERVICE_CHANGE",
		ProviderAccountID: serviceRecord.ProviderAccountID,
		BillingCycle:      billingCycle,
		Amount:            input.Amount,
		Status:            domain.OrderStatusPending,
		Configuration:     cloneConfigSelections(serviceRecord.Configuration),
		ResourceSnapshot:  cloneSnapshot(serviceRecord.ResourceSnapshot),
		CreatedAt:         now.Format("2006-01-02 15:04:05"),
	}
	repository.nextOrderID++

	invoice := domain.Invoice{
		ID:           repository.nextInvoiceID,
		InvoiceNo:    fmt.Sprintf("INV-%08d", repository.nextInvoiceID),
		OrderID:      order.ID,
		OrderNo:      order.OrderNo,
		CustomerID:   serviceRecord.CustomerID,
		ProductName:  title,
		TotalAmount:  input.Amount,
		Status:       domain.InvoiceStatusUnpaid,
		DueAt:        now.Add(72 * time.Hour).Format("2006-01-02 15:04:05"),
		BillingCycle: billingCycle,
	}
	repository.nextInvoiceID++

	repository.orders = append([]domain.Order{order}, repository.orders...)
	repository.invoices = append([]domain.Invoice{invoice}, repository.invoices...)
	repository.serviceChanges = append([]memoryServiceChangeLink{{
		ID:           int64(len(repository.serviceChanges) + 1),
		ServiceID:    serviceID,
		OrderID:      order.ID,
		InvoiceID:    invoice.ID,
		ActionName:   strings.TrimSpace(input.ActionName),
		Title:        title,
		Status:       string(domain.InvoiceStatusUnpaid),
		Reason:       strings.TrimSpace(input.Reason),
		BillingCycle: billingCycle,
		Amount:       input.Amount,
		Payload:      maps.Clone(input.Payload),
		CreatedAt:    now.Format("2006-01-02 15:04:05"),
		UpdatedAt:    now.Format("2006-01-02 15:04:05"),
	}}, repository.serviceChanges...)

	order = repository.decorateOrderLocked(cloneOrder(order))
	serviceCopy := cloneService(serviceRecord)
	return order, invoice, &serviceCopy, true, nil
}

func (repository *MemoryRepository) RefundInvoice(invoiceID int64, reason string) (domain.Invoice, *domain.ServiceRecord, domain.RefundRecord, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for invoiceIndex, invoice := range repository.invoices {
		if invoice.ID != invoiceID || invoice.Status != domain.InvoiceStatusPaid {
			continue
		}

		now := time.Now()
		invoice.Status = domain.InvoiceStatusRefunded
		repository.invoices[invoiceIndex] = invoice

		refund := domain.RefundRecord{
			ID:         repository.nextRefundID,
			RefundNo:   fmt.Sprintf("RFD-%08d", repository.nextRefundID),
			InvoiceID:  invoice.ID,
			OrderID:    invoice.OrderID,
			CustomerID: invoice.CustomerID,
			Amount:     invoice.TotalAmount,
			Reason:     reason,
			Status:     domain.RefundStatusCompleted,
			CreatedAt:  now.Format("2006-01-02 15:04:05"),
		}
		repository.nextRefundID++
		repository.refunds = append([]domain.RefundRecord{refund}, repository.refunds...)

		var updatedService *domain.ServiceRecord
		changeLink, isServiceChange := repository.findServiceChangeByInvoiceLocked(invoice.ID)
		for orderIndex, order := range repository.orders {
			if order.ID != invoice.OrderID {
				continue
			}
			order.Status = domain.OrderStatusPending
			repository.orders[orderIndex] = order
			break
		}

		if isServiceChange {
			changeLink.Status = string(domain.InvoiceStatusRefunded)
			changeLink.RefundedAt = now.Format("2006-01-02 15:04:05")
			changeLink.UpdatedAt = now.Format("2006-01-02 15:04:05")
			repository.replaceServiceChangeLocked(changeLink)
			for _, item := range repository.services {
				if item.ID != changeLink.ServiceID {
					continue
				}
				serviceCopy := cloneService(item)
				updatedService = &serviceCopy
				break
			}
		} else {
			for serviceIndex, item := range repository.services {
				if item.InvoiceID != invoice.ID {
					continue
				}
				item.Status = domain.ServiceStatusTerminated
				item.LastAction = "refund-terminate"
				item.UpdatedAt = now.Format("2006-01-02 15:04:05")
				repository.services[serviceIndex] = item
				serviceCopy := cloneService(item)
				updatedService = &serviceCopy
				break
			}
		}

		return invoice, updatedService, refund, true
	}

	return domain.Invoice{}, nil, domain.RefundRecord{}, false
}

func (repository *MemoryRepository) findLatestPaymentLocked(invoiceID int64) (domain.PaymentRecord, bool) {
	for _, payment := range repository.payments {
		if payment.InvoiceID == invoiceID {
			return payment, true
		}
	}
	return domain.PaymentRecord{}, false
}

func (repository *MemoryRepository) findServiceChangeByInvoiceLocked(invoiceID int64) (memoryServiceChangeLink, bool) {
	for _, item := range repository.serviceChanges {
		if item.InvoiceID == invoiceID {
			return item, true
		}
	}
	return memoryServiceChangeLink{}, false
}

func (repository *MemoryRepository) findServiceChangeByOrderLocked(orderID int64) (memoryServiceChangeLink, bool) {
	for _, item := range repository.serviceChanges {
		if item.OrderID == orderID {
			return item, true
		}
	}
	return memoryServiceChangeLink{}, false
}

func (repository *MemoryRepository) replaceServiceChangeLocked(link memoryServiceChangeLink) {
	for index, item := range repository.serviceChanges {
		if item.OrderID == link.OrderID {
			repository.serviceChanges[index] = link
			return
		}
	}
}

func cloneServiceChange(link memoryServiceChangeLink) domain.ServiceChangeOrder {
	return domain.ServiceChangeOrder{
		ID:           link.ID,
		ServiceID:    link.ServiceID,
		OrderID:      link.OrderID,
		InvoiceID:    link.InvoiceID,
		ActionName:   link.ActionName,
		Title:        link.Title,
		Amount:       link.Amount,
		Status:       link.Status,
		Reason:       link.Reason,
		BillingCycle: link.BillingCycle,
		Payload:      maps.Clone(link.Payload),
		PaidAt:       link.PaidAt,
		RefundedAt:   link.RefundedAt,
		CreatedAt:    link.CreatedAt,
		UpdatedAt:    link.UpdatedAt,
	}
}

func firstNonEmpty(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func nextDueDate(cycle string, now time.Time) string {
	switch cycle {
	case "annual":
		return now.AddDate(1, 0, 0).Format("2006-01-02")
	case "quarterly":
		return now.AddDate(0, 3, 0).Format("2006-01-02")
	default:
		return now.AddDate(0, 1, 0).Format("2006-01-02")
	}
}

func normalizeAutomationType(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "MOFANG_CLOUD":
		return "MOFANG_CLOUD"
	case "ZJMF_API":
		return "ZJMF_API"
	case "SERVICE_CHANGE":
		return "SERVICE_CHANGE"
	case "RESOURCE":
		return "RESOURCE"
	case "MANUAL":
		return "MANUAL"
	default:
		return "LOCAL"
	}
}

func cloneOrders(items []domain.Order) []domain.Order {
	result := make([]domain.Order, 0, len(items))
	for _, item := range items {
		result = append(result, cloneOrder(item))
	}
	return result
}

func cloneOrder(item domain.Order) domain.Order {
	item.Configuration = cloneConfigSelections(item.Configuration)
	item.ResourceSnapshot = cloneSnapshot(item.ResourceSnapshot)
	return item
}

func cloneServices(items []domain.ServiceRecord) []domain.ServiceRecord {
	result := make([]domain.ServiceRecord, 0, len(items))
	for _, item := range items {
		result = append(result, cloneService(item))
	}
	return result
}

func cloneService(item domain.ServiceRecord) domain.ServiceRecord {
	item.Configuration = cloneConfigSelections(item.Configuration)
	item.ResourceSnapshot = cloneSnapshot(item.ResourceSnapshot)
	return item
}

func cloneConfigSelections(items []domain.ServiceConfigSelection) []domain.ServiceConfigSelection {
	return slices.Clone(items)
}

func cloneSnapshot(item domain.ServiceResourceSnapshot) domain.ServiceResourceSnapshot {
	return item
}

func (repository *MemoryRepository) decorateOrdersLocked(items []domain.Order) []domain.Order {
	for index := range items {
		items[index] = repository.decorateOrderLocked(items[index])
	}
	return items
}

func (repository *MemoryRepository) decorateOrderLocked(item domain.Order) domain.Order {
	item.Payment = ""
	item.PayStatus = "UNPAID"

	var latestInvoice *domain.Invoice
	for _, invoice := range repository.invoices {
		if invoice.OrderID != item.ID {
			continue
		}
		if latestInvoice == nil || invoice.ID > latestInvoice.ID {
			current := invoice
			latestInvoice = &current
		}
	}
	if latestInvoice != nil {
		switch latestInvoice.Status {
		case domain.InvoiceStatusRefunded:
			item.PayStatus = "REFUNDED"
		case domain.InvoiceStatusPaid:
			item.PayStatus = "PAID"
		default:
			item.PayStatus = "UNPAID"
		}
	}

	var latestPayment *domain.PaymentRecord
	for _, payment := range repository.payments {
		if payment.OrderID != item.ID {
			continue
		}
		if latestPayment == nil || payment.ID > latestPayment.ID {
			current := payment
			latestPayment = &current
		}
	}
	if latestPayment != nil {
		item.Payment = latestPayment.Channel
	}

	return item
}

func filterMemoryOrders(items []domain.Order, filter domain.OrderListFilter) []domain.Order {
	filtered := make([]domain.Order, 0, len(items))
	for _, item := range items {
		if filter.Status != "" && !strings.EqualFold(string(item.Status), filter.Status) {
			continue
		}
		if filter.OrderNo != "" && !strings.Contains(strings.ToLower(item.OrderNo), strings.ToLower(filter.OrderNo)) {
			continue
		}
		if filter.ProductName != "" && !strings.Contains(strings.ToLower(item.ProductName), strings.ToLower(filter.ProductName)) {
			continue
		}
		if filter.CustomerID > 0 && item.CustomerID != filter.CustomerID {
			continue
		}
		if filter.Payment != "" && !strings.EqualFold(item.Payment, filter.Payment) {
			continue
		}
		if filter.PayStatus != "" && !strings.EqualFold(item.PayStatus, filter.PayStatus) {
			continue
		}
		if filter.HasAmount && math.Abs(item.Amount-filter.Amount) > 0.00001 {
			continue
		}
		if filter.StartTime != "" && !timeMatchesRange(item.CreatedAt, filter.StartTime, "") {
			continue
		}
		if filter.EndTime != "" && !timeMatchesRange(item.CreatedAt, "", filter.EndTime) {
			continue
		}
		filtered = append(filtered, item)
	}

	slices.SortFunc(filtered, func(left, right domain.Order) int {
		desc := strings.EqualFold(filter.Order, "desc")
		switch strings.ToLower(filter.Sort) {
		case "amount":
			if left.Amount == right.Amount {
				return compareByID(left.ID, right.ID, desc)
			}
			if desc {
				if left.Amount > right.Amount {
					return -1
				}
				return 1
			}
			if left.Amount < right.Amount {
				return -1
			}
			return 1
		case "ordernum":
			return compareString(left.OrderNo, right.OrderNo, desc)
		default:
			return compareString(left.CreatedAt, right.CreatedAt, desc)
		}
	})

	return filtered
}

func paginateMemoryOrders(items []domain.Order, page, limit int) []domain.Order {
	if limit <= 0 {
		return items
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	if start >= len(items) {
		return []domain.Order{}
	}
	end := start + limit
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func compareByID(left, right int64, desc bool) int {
	if left == right {
		return 0
	}
	if desc {
		if left > right {
			return -1
		}
		return 1
	}
	if left < right {
		return -1
	}
	return 1
}

func compareString(left, right string, desc bool) int {
	if left == right {
		return 0
	}
	if desc {
		if left > right {
			return -1
		}
		return 1
	}
	if left < right {
		return -1
	}
	return 1
}

func timeMatchesRange(createdAt, start, end string) bool {
	createdTime, ok := parseOrderTime(createdAt)
	if !ok {
		return true
	}
	if start != "" {
		startTime, ok := parseFilterTime(start, false)
		if ok && createdTime.Before(startTime) {
			return false
		}
	}
	if end != "" {
		endTime, ok := parseFilterTime(end, true)
		if ok && createdTime.After(endTime) {
			return false
		}
	}
	return true
}

func parseOrderTime(value string) (time.Time, bool) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func parseFilterTime(value string, endOfDay bool) (time.Time, bool) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err != nil {
			continue
		}
		if layout == "2006-01-02" && endOfDay {
			parsed = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		return parsed, true
	}
	return time.Time{}, false
}

func filterMemoryInvoices(items []domain.Invoice, filter domain.InvoiceListFilter) []domain.Invoice {
	filtered := make([]domain.Invoice, 0, len(items))
	for _, item := range items {
		if filter.Status != "" && !strings.EqualFold(string(item.Status), filter.Status) {
			continue
		}
		if filter.InvoiceNo != "" && !strings.Contains(strings.ToLower(item.InvoiceNo), strings.ToLower(filter.InvoiceNo)) {
			continue
		}
		if filter.OrderNo != "" && !strings.Contains(strings.ToLower(item.OrderNo), strings.ToLower(filter.OrderNo)) {
			continue
		}
		if filter.ProductName != "" && !strings.Contains(strings.ToLower(item.ProductName), strings.ToLower(filter.ProductName)) {
			continue
		}
		if filter.BillingCycle != "" && !strings.EqualFold(item.BillingCycle, filter.BillingCycle) {
			continue
		}
		if filter.CustomerID > 0 && item.CustomerID != filter.CustomerID {
			continue
		}
		filtered = append(filtered, item)
	}

	slices.SortFunc(filtered, func(left, right domain.Invoice) int {
		desc := strings.EqualFold(filter.Order, "desc")
		switch strings.ToLower(filter.Sort) {
		case "invoice_no":
			return compareString(left.InvoiceNo, right.InvoiceNo, desc)
		case "amount", "total_amount":
			if left.TotalAmount == right.TotalAmount {
				return compareByID(left.ID, right.ID, desc)
			}
			if desc {
				if left.TotalAmount > right.TotalAmount {
					return -1
				}
				return 1
			}
			if left.TotalAmount < right.TotalAmount {
				return -1
			}
			return 1
		case "due_at":
			return compareString(left.DueAt, right.DueAt, desc)
		default:
			return compareByID(left.ID, right.ID, desc)
		}
	})

	return filtered
}

func paginateMemoryInvoices(items []domain.Invoice, page, limit int) []domain.Invoice {
	if limit <= 0 {
		return items
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	if start >= len(items) {
		return []domain.Invoice{}
	}
	end := start + limit
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func filterMemoryServices(items []domain.ServiceRecord, filter domain.ServiceListFilter) []domain.ServiceRecord {
	filtered := make([]domain.ServiceRecord, 0, len(items))
	for _, item := range items {
		if filter.Status != "" && !strings.EqualFold(string(item.Status), filter.Status) {
			continue
		}
		if filter.ProviderType != "" && !strings.EqualFold(item.ProviderType, filter.ProviderType) {
			continue
		}
		if filter.ProviderAccountID > 0 && item.ProviderAccountID != filter.ProviderAccountID {
			continue
		}
		if filter.SyncStatus != "" && !strings.EqualFold(item.SyncStatus, filter.SyncStatus) {
			continue
		}
		if filter.Keyword != "" {
			keyword := strings.ToLower(filter.Keyword)
			if !strings.Contains(strings.ToLower(item.ServiceNo), keyword) &&
				!strings.Contains(strings.ToLower(item.ProductName), keyword) &&
				!strings.Contains(strings.ToLower(item.IPAddress), keyword) &&
				!strings.Contains(strings.ToLower(item.ProviderResourceID), keyword) {
				continue
			}
		}
		filtered = append(filtered, item)
	}

	slices.SortFunc(filtered, func(left, right domain.ServiceRecord) int {
		desc := strings.EqualFold(filter.Order, "desc")
		switch strings.ToLower(filter.Sort) {
		case "service_no":
			return compareString(left.ServiceNo, right.ServiceNo, desc)
		case "next_due_at":
			return compareString(left.NextDueAt, right.NextDueAt, desc)
		default:
			return compareByID(left.ID, right.ID, desc)
		}
	})

	return filtered
}

func paginateMemoryServices(items []domain.ServiceRecord, page, limit int) []domain.ServiceRecord {
	if limit <= 0 {
		return items
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	if start >= len(items) {
		return []domain.ServiceRecord{}
	}
	end := start + limit
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}
