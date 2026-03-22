package dto

import (
	customerDomain "idc-finance/internal/modules/customer/domain"
	orderDomain "idc-finance/internal/modules/order/domain"
)

type Stat struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Ticket struct {
	No        string `json:"no"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updatedAt"`
}

type WalletSummary struct {
	Balance     string `json:"balance"`
	CreditLimit string `json:"creditLimit"`
}

type DashboardResponse struct {
	Stats    []Stat                      `json:"stats"`
	Services []orderDomain.ServiceRecord `json:"services"`
	Orders   []orderDomain.Order         `json:"orders"`
	Invoices []orderDomain.Invoice       `json:"invoices"`
	Tickets  []Ticket                    `json:"tickets"`
	Wallet   WalletSummary               `json:"wallet"`
}

type AccountResponse struct {
	Customer       customerDomain.Customer `json:"customer"`
	PrimaryContact *customerDomain.Contact `json:"primaryContact,omitempty"`
	Wallet         WalletSummary           `json:"wallet"`
}
