package service

import (
	"fmt"
	"strings"

	customerDomain "idc-finance/internal/modules/customer/domain"
	customerDTO "idc-finance/internal/modules/customer/dto"
	customerService "idc-finance/internal/modules/customer/service"
	orderDomain "idc-finance/internal/modules/order/domain"
	orderDTO "idc-finance/internal/modules/order/dto"
	orderService "idc-finance/internal/modules/order/service"
	portalDTO "idc-finance/internal/modules/portal/dto"
)

type Service struct {
	customer *customerService.Service
	order    *orderService.Service
}

func New(customerSvc *customerService.Service, orderSvc *orderService.Service) *Service {
	return &Service{customer: customerSvc, order: orderSvc}
}

func (service *Service) GetDashboard(customerID int64) (portalDTO.DashboardResponse, bool) {
	_, exists := service.customer.GetByID(customerID)
	if !exists {
		return portalDTO.DashboardResponse{}, false
	}

	orders := service.order.ListOrdersByCustomer(customerID)
	invoices := service.order.ListInvoicesByCustomer(customerID)
	services := service.order.ListServicesByCustomer(customerID)
	tickets, _ := service.customer.ListTickets(customerID)
	wallet := service.walletSummaryForCustomer(customerID)

	return portalDTO.DashboardResponse{
		Stats: []portalDTO.Stat{
			{Label: "服务数量", Value: fmt.Sprintf("%d", len(services))},
			{Label: "未付账单", Value: fmt.Sprintf("%d", countUnpaidInvoices(invoices))},
			{Label: "处理中工单", Value: fmt.Sprintf("%d", countOpenTickets(tickets))},
			{Label: "账户余额", Value: formatCurrency(wallet.Balance)},
		},
		Services: services,
		Orders:   orders,
		Invoices: invoices,
		Tickets:  mapTickets(tickets),
		Wallet:   wallet,
	}, true
}

func (service *Service) ListTickets(customerID int64) ([]portalDTO.Ticket, bool) {
	if _, exists := service.customer.GetByID(customerID); !exists {
		return nil, false
	}

	items, _ := service.customer.ListTickets(customerID)
	return mapTickets(items), true
}

func (service *Service) GetAccount(customerID int64) (portalDTO.AccountResponse, bool) {
	customer, exists := service.customer.GetByID(customerID)
	if !exists {
		return portalDTO.AccountResponse{}, false
	}

	account := portalDTO.AccountResponse{
		Customer: customer,
		Wallet:   service.walletSummaryForCustomer(customerID),
	}
	for index := range customer.Contacts {
		if customer.Contacts[index].IsPrimary {
			account.PrimaryContact = &customer.Contacts[index]
			break
		}
	}
	if account.PrimaryContact == nil && len(customer.Contacts) > 0 {
		account.PrimaryContact = &customer.Contacts[0]
	}

	return account, true
}

func (service *Service) UpdateProfile(customerID int64, request customerDTO.UpdateCustomerRequest, requestID string) (portalDTO.AccountResponse, bool) {
	updated, ok := service.customer.UpdateSelf(customerID, request, requestID)
	if !ok {
		return portalDTO.AccountResponse{}, false
	}
	return portalDTO.AccountResponse{
		Customer: updated,
		Wallet:   service.walletSummaryForCustomer(customerID),
	}, true
}

func (service *Service) CreateContact(customerID int64, request customerDTO.CreateContactRequest, requestID string) (customerDomain.Contact, bool) {
	return service.customer.AddContactByCustomer(customerID, request, requestID)
}

func (service *Service) UpdateContact(customerID, contactID int64, request customerDTO.UpdateContactRequest, requestID string) (customerDomain.Contact, bool) {
	return service.customer.UpdateContactByCustomer(customerID, contactID, request, requestID)
}

func (service *Service) DeleteContact(customerID, contactID int64, requestID string) bool {
	return service.customer.DeleteContactByCustomer(customerID, contactID, requestID)
}

func (service *Service) GetIdentity(customerID int64) (*customerDomain.Identity, bool) {
	account, ok := service.GetAccount(customerID)
	if !ok {
		return nil, false
	}
	return account.Customer.Identity, true
}

func (service *Service) SubmitIdentity(customerID int64, request customerDTO.SubmitIdentityRequest, requestID string) (customerDomain.Identity, bool) {
	return service.customer.SubmitIdentityByCustomer(customerID, request, requestID)
}

func (service *Service) GetWallet(customerID int64) (portalDTO.WalletOverviewResponse, bool) {
	if _, exists := service.customer.GetByID(customerID); !exists {
		return portalDTO.WalletOverviewResponse{}, false
	}

	return portalDTO.WalletOverviewResponse{
		Wallet:       service.walletSummaryForCustomer(customerID),
		Transactions: service.order.ListAccountTransactionsByCustomer(customerID, 50),
	}, true
}

func (service *Service) ListWalletTransactions(customerID int64, transactionType, keyword, channel string, limit int) ([]orderDomain.AccountTransaction, bool) {
	if _, exists := service.customer.GetByID(customerID); !exists {
		return nil, false
	}

	items, _ := service.order.ListAccountTransactions(orderDomain.AccountTransactionFilter{
		Page:            1,
		Limit:           limit,
		Sort:            "occurred_at",
		Order:           "desc",
		CustomerID:      customerID,
		TransactionType: strings.TrimSpace(transactionType),
		Keyword:         strings.TrimSpace(keyword),
		Channel:         strings.TrimSpace(channel),
	})
	return items, true
}

func (service *Service) GetOrderDetail(customerID, orderID int64) (orderDTO.OrderDetailResponse, bool) {
	detail, ok := service.order.GetOrderDetail(orderID)
	if !ok || detail.Order.CustomerID != customerID {
		return orderDTO.OrderDetailResponse{}, false
	}
	return detail, true
}

func (service *Service) CreateOrderRequest(customerID int64, customerName string, orderID int64, request orderDTO.CreateOrderRequestRequest, requestID string) (orderDTO.OrderDetailResponse, bool, error) {
	return service.order.CreatePortalOrderRequest(customerID, customerName, orderID, request, requestID)
}

func (service *Service) GetInvoiceDetail(customerID, invoiceID int64) (orderDTO.InvoiceDetailResponse, bool) {
	detail, ok := service.order.GetInvoiceDetail(invoiceID)
	if !ok || detail.Invoice.CustomerID != customerID {
		return orderDTO.InvoiceDetailResponse{}, false
	}
	return detail, true
}

func (service *Service) ListPayments(customerID int64, invoiceID int64, keyword, channel, status string) ([]orderDomain.PaymentRecord, bool) {
	if _, exists := service.customer.GetByID(customerID); !exists {
		return nil, false
	}
	items, _ := service.order.ListPayments(orderDomain.PaymentListFilter{
		Page:       1,
		Limit:      200,
		Sort:       "paid_at",
		Order:      "desc",
		CustomerID: customerID,
		InvoiceID:  invoiceID,
		Keyword:    strings.TrimSpace(keyword),
		Channel:    strings.TrimSpace(channel),
		Status:     strings.TrimSpace(status),
	})
	return items, true
}

func (service *Service) ListRefunds(customerID int64, invoiceID int64, keyword, status string) ([]orderDomain.RefundRecord, bool) {
	if _, exists := service.customer.GetByID(customerID); !exists {
		return nil, false
	}
	items, _ := service.order.ListRefunds(orderDomain.RefundListFilter{
		Page:       1,
		Limit:      200,
		Sort:       "created_at",
		Order:      "desc",
		CustomerID: customerID,
		InvoiceID:  invoiceID,
		Keyword:    strings.TrimSpace(keyword),
		Status:     strings.TrimSpace(status),
	})
	return items, true
}

func mapTickets(items []customerDTO.RelatedItem) []portalDTO.Ticket {
	tickets := make([]portalDTO.Ticket, 0, len(items))
	for _, item := range items {
		tickets = append(tickets, portalDTO.Ticket{
			No:        item.No,
			Title:     item.Name,
			Status:    item.Status,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return tickets
}

func countUnpaidInvoices(items []orderDomain.Invoice) int {
	total := 0
	for _, item := range items {
		if item.Status == orderDomain.InvoiceStatusUnpaid {
			total++
		}
	}
	return total
}

func countOpenTickets(items []customerDTO.RelatedItem) int {
	total := 0
	for _, item := range items {
		if item.Status != "CLOSED" {
			total++
		}
	}
	return total
}

func (service *Service) walletSummaryForCustomer(customerID int64) portalDTO.WalletSummary {
	if wallet, ok := service.order.GetCustomerWallet(customerID); ok {
		return portalDTO.WalletSummary{
			Balance:         formatWalletValue(wallet.Balance),
			CreditLimit:     formatWalletValue(wallet.CreditLimit),
			CreditUsed:      formatWalletValue(wallet.CreditUsed),
			AvailableCredit: formatWalletValue(wallet.AvailableCredit),
		}
	}
	return portalDTO.WalletSummary{
		Balance:         "0.00",
		CreditLimit:     "0.00",
		CreditUsed:      "0.00",
		AvailableCredit: "0.00",
	}
}

func formatWalletValue(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func formatCurrency(value string) string {
	return fmt.Sprintf("¥%s", value)
}
