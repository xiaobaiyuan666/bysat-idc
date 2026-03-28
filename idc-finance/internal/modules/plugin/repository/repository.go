package repository

import "idc-finance/internal/modules/plugin/domain"

type Repository interface {
	ListDefinitions() []domain.PluginDefinition
	GetConfig(code string) (domain.PluginConfig, bool)
	ListConfigs() []domain.PluginConfig
	SaveConfig(code string, enabled bool, config map[string]any, updatedBy string) (domain.PluginConfig, error)

	CreateKYCRecord(record domain.KYCRecord) (domain.KYCRecord, error)
	ListKYCRecords(customerID int64, limit int) []domain.KYCRecord

	CreateSMSRecord(record domain.SMSRecord) (domain.SMSRecord, error)
	ListSMSRecords(customerID int64, limit int) []domain.SMSRecord

	CreatePaymentOrder(order domain.PaymentOrder) (domain.PaymentOrder, error)
	GetPaymentOrderByOrderNo(orderNo string) (domain.PaymentOrder, bool)
	UpdatePaymentOrder(order domain.PaymentOrder) (domain.PaymentOrder, bool, error)
	ListPaymentOrders(customerID int64, limit int) []domain.PaymentOrder
}
