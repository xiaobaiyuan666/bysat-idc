package service

import (
	"fmt"

	customerDTO "idc-finance/internal/modules/customer/dto"
	customerService "idc-finance/internal/modules/customer/service"
	orderDomain "idc-finance/internal/modules/order/domain"
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

func (service *Service) GetWallet(customerID int64) (portalDTO.WalletOverviewResponse, bool) {
	if _, exists := service.customer.GetByID(customerID); !exists {
		return portalDTO.WalletOverviewResponse{}, false
	}

	return portalDTO.WalletOverviewResponse{
		Wallet:       service.walletSummaryForCustomer(customerID),
		Transactions: service.order.ListAccountTransactionsByCustomer(customerID, 50),
	}, true
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
	return fmt.Sprintf("¥ %s", value)
}
