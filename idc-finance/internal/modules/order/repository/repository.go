package repository

import "idc-finance/internal/modules/order/domain"

type Repository interface {
	ListOrders(filter domain.OrderListFilter) ([]domain.Order, int)
	ListOrdersByCustomer(customerID int64) []domain.Order
	GetOrderByID(id int64) (domain.Order, bool)
	ListInvoices(filter domain.InvoiceListFilter) ([]domain.Invoice, int)
	ListInvoicesByCustomer(customerID int64) []domain.Invoice
	ListInvoicesByOrder(orderID int64) []domain.Invoice
	GetInvoiceByID(id int64) (domain.Invoice, bool)
	GetCustomerWallet(customerID int64) (domain.CustomerWallet, bool)
	ListAccountTransactions(filter domain.AccountTransactionFilter) ([]domain.AccountTransaction, int)
	ListAccountTransactionsByCustomer(customerID int64, limit int) []domain.AccountTransaction
	ListPayments(filter domain.PaymentListFilter) ([]domain.PaymentRecord, int)
	ListRefunds(filter domain.RefundListFilter) ([]domain.RefundRecord, int)
	ListServiceChangeOrders(filter domain.ServiceChangeOrderListFilter) ([]domain.ServiceChangeOrder, int)
	AdjustCustomerWallet(customerID int64, input domain.WalletAdjustment) (domain.CustomerWallet, domain.AccountTransaction, bool, error)
	GetServiceChangeOrderByInvoiceID(invoiceID int64) (domain.ServiceChangeOrder, bool)
	GetServiceChangeOrderByOrderID(orderID int64) (domain.ServiceChangeOrder, bool)
	ListServiceChangeOrdersByService(serviceID int64) []domain.ServiceChangeOrder
	UpdatePendingOrder(orderID int64, update domain.PendingOrderUpdate) (domain.Order, *domain.Invoice, bool, error)
	UpdateUnpaidInvoice(invoiceID int64, update domain.UnpaidInvoiceUpdate) (domain.Invoice, *domain.Order, bool, error)
	ListServices(filter domain.ServiceListFilter) ([]domain.ServiceRecord, int)
	ListServicesByCustomer(customerID int64) []domain.ServiceRecord
	ListServicesByOrder(orderID int64) []domain.ServiceRecord
	GetServiceByID(id int64) (domain.ServiceRecord, bool)
	UpdateServiceRecord(serviceID int64, update domain.ManualServiceUpdate) (domain.ServiceRecord, bool, error)
	CreateServiceChangeOrder(serviceID int64, input domain.ServiceChangeOrderInput) (domain.Order, domain.Invoice, *domain.ServiceRecord, bool, error)
	ListPaymentsByInvoice(invoiceID int64) []domain.PaymentRecord
	ListRefundsByInvoice(invoiceID int64) []domain.RefundRecord
	Checkout(
		customerID int64,
		customerName string,
		productID int64,
		productName, productType, automationType string,
		providerAccountID int64,
		cycleCode string,
		amount float64,
		configuration []domain.ServiceConfigSelection,
		resourceSnapshot domain.ServiceResourceSnapshot,
	) (domain.Order, domain.Invoice)
	PayInvoice(invoiceID int64, channel, source, operator, tradeNo string) (domain.Invoice, *domain.ServiceRecord, domain.PaymentRecord, bool)
	ExecuteServiceAction(serviceID int64, action string, params domain.ServiceActionParams) (domain.ServiceRecord, bool)
	UpdateServiceProvisioning(serviceID int64, update domain.ServiceProvisioningUpdate) (domain.ServiceRecord, bool)
	RefundInvoice(invoiceID int64, reason string) (domain.Invoice, *domain.ServiceRecord, domain.RefundRecord, bool)
}
