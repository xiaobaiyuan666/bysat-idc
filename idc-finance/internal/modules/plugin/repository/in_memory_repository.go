package repository

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"idc-finance/internal/modules/plugin/domain"
)

type MemoryRepository struct {
	mu            sync.RWMutex
	definitions   map[string]domain.PluginDefinition
	configs       map[string]domain.PluginConfig
	kycRecords    []domain.KYCRecord
	smsRecords    []domain.SMSRecord
	paymentOrders []domain.PaymentOrder
	nextID        int64
}

func NewMemoryRepository() *MemoryRepository {
	repository := &MemoryRepository{
		definitions: map[string]domain.PluginDefinition{},
		configs:     map[string]domain.PluginConfig{},
		nextID:      1,
	}

	for _, definition := range defaultDefinitions() {
		repository.definitions[definition.Code] = cloneDefinition(definition)
		repository.configs[definition.Code] = domain.PluginConfig{
			Code:    definition.Code,
			Name:    definition.Name,
			Enabled: true,
			Config:  map[string]any{},
		}
	}

	return repository
}

func (repository *MemoryRepository) ListDefinitions() []domain.PluginDefinition {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.PluginDefinition, 0, len(repository.definitions))
	for _, item := range repository.definitions {
		items = append(items, cloneDefinition(item))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Code < items[j].Code
	})
	return items
}

func (repository *MemoryRepository) GetConfig(code string) (domain.PluginConfig, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	item, ok := repository.configs[normalizeCode(code)]
	if !ok {
		return domain.PluginConfig{}, false
	}
	return cloneConfig(item), true
}

func (repository *MemoryRepository) ListConfigs() []domain.PluginConfig {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.PluginConfig, 0, len(repository.configs))
	for _, item := range repository.configs {
		items = append(items, cloneConfig(item))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Code < items[j].Code
	})
	return items
}

func (repository *MemoryRepository) SaveConfig(code string, enabled bool, config map[string]any, updatedBy string) (domain.PluginConfig, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	normalized := normalizeCode(code)
	definition, ok := repository.definitions[normalized]
	if !ok {
		return domain.PluginConfig{}, fmt.Errorf("plugin not found: %s", code)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	item := domain.PluginConfig{
		Code:      definition.Code,
		Name:      definition.Name,
		Enabled:   enabled,
		Config:    cloneMap(config),
		UpdatedBy: strings.TrimSpace(updatedBy),
		UpdatedAt: now,
	}
	repository.configs[normalized] = item
	return cloneConfig(item), nil
}

func (repository *MemoryRepository) CreateKYCRecord(record domain.KYCRecord) (domain.KYCRecord, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	record.ID = repository.takeIDLocked()
	record.PluginCode = normalizeCode(record.PluginCode)
	record.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	record.Raw = cloneMap(record.Raw)
	repository.kycRecords = append(repository.kycRecords, record)
	return cloneKYCRecord(record), nil
}

func (repository *MemoryRepository) ListKYCRecords(customerID int64, limit int) []domain.KYCRecord {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.KYCRecord, 0)
	for index := len(repository.kycRecords) - 1; index >= 0; index-- {
		item := repository.kycRecords[index]
		if customerID > 0 && item.CustomerID != customerID {
			continue
		}
		items = append(items, cloneKYCRecord(item))
		if limit > 0 && len(items) >= limit {
			break
		}
	}
	return items
}

func (repository *MemoryRepository) CreateSMSRecord(record domain.SMSRecord) (domain.SMSRecord, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	record.ID = repository.takeIDLocked()
	record.PluginCode = normalizeCode(record.PluginCode)
	record.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	record.Raw = cloneMap(record.Raw)
	repository.smsRecords = append(repository.smsRecords, record)
	return cloneSMSRecord(record), nil
}

func (repository *MemoryRepository) ListSMSRecords(customerID int64, limit int) []domain.SMSRecord {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.SMSRecord, 0)
	for index := len(repository.smsRecords) - 1; index >= 0; index-- {
		item := repository.smsRecords[index]
		if customerID > 0 && item.CustomerID != customerID {
			continue
		}
		items = append(items, cloneSMSRecord(item))
		if limit > 0 && len(items) >= limit {
			break
		}
	}
	return items
}

func (repository *MemoryRepository) CreatePaymentOrder(order domain.PaymentOrder) (domain.PaymentOrder, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05")
	order.ID = repository.takeIDLocked()
	order.PluginCode = normalizeCode(order.PluginCode)
	order.CreatedAt = now
	order.UpdatedAt = now
	order.Raw = cloneMap(order.Raw)
	repository.paymentOrders = append(repository.paymentOrders, order)
	return clonePaymentOrder(order), nil
}

func (repository *MemoryRepository) GetPaymentOrderByOrderNo(orderNo string) (domain.PaymentOrder, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	target := strings.TrimSpace(orderNo)
	for _, item := range repository.paymentOrders {
		if item.OrderNo == target {
			return clonePaymentOrder(item), true
		}
	}
	return domain.PaymentOrder{}, false
}

func (repository *MemoryRepository) UpdatePaymentOrder(order domain.PaymentOrder) (domain.PaymentOrder, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	target := strings.TrimSpace(order.OrderNo)
	for index := range repository.paymentOrders {
		if repository.paymentOrders[index].OrderNo != target {
			continue
		}
		order.PluginCode = normalizeCode(order.PluginCode)
		order.Raw = cloneMap(order.Raw)
		order.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		repository.paymentOrders[index] = order
		return clonePaymentOrder(order), true, nil
	}
	return domain.PaymentOrder{}, false, nil
}

func (repository *MemoryRepository) ListPaymentOrders(customerID int64, limit int) []domain.PaymentOrder {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.PaymentOrder, 0)
	for index := len(repository.paymentOrders) - 1; index >= 0; index-- {
		item := repository.paymentOrders[index]
		if customerID > 0 && item.CustomerID != customerID {
			continue
		}
		items = append(items, clonePaymentOrder(item))
		if limit > 0 && len(items) >= limit {
			break
		}
	}
	return items
}

func (repository *MemoryRepository) takeIDLocked() int64 {
	value := repository.nextID
	repository.nextID++
	return value
}

func normalizeCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func cloneMap(source map[string]any) map[string]any {
	if len(source) == 0 {
		return map[string]any{}
	}
	result := make(map[string]any, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

func cloneDefinition(source domain.PluginDefinition) domain.PluginDefinition {
	item := source
	item.Fields = append([]domain.PluginField(nil), source.Fields...)
	return item
}

func cloneConfig(source domain.PluginConfig) domain.PluginConfig {
	item := source
	item.Config = cloneMap(source.Config)
	return item
}

func cloneKYCRecord(source domain.KYCRecord) domain.KYCRecord {
	item := source
	item.Raw = cloneMap(source.Raw)
	return item
}

func cloneSMSRecord(source domain.SMSRecord) domain.SMSRecord {
	item := source
	item.Raw = cloneMap(source.Raw)
	return item
}

func clonePaymentOrder(source domain.PaymentOrder) domain.PaymentOrder {
	item := source
	item.Raw = cloneMap(source.Raw)
	return item
}

func defaultDefinitions() []domain.PluginDefinition {
	return []domain.PluginDefinition{
		{
			Code:        string(domain.PluginCodeZhimaKYC),
			Name:        "芝麻实名认证",
			Description: "芝麻实名认证插件骨架，当前提供可演示的本地模拟校验链路。",
			Fields: []domain.PluginField{
				{Key: "appId", Label: "应用 ID", Required: true, Placeholder: "支付宝 / 芝麻应用 ID", Secret: false},
				{Key: "privateKey", Label: "商户私钥", Required: true, Placeholder: "PKCS8 私钥", Secret: true},
				{Key: "publicKey", Label: "平台公钥", Required: true, Placeholder: "芝麻平台公钥", Secret: true},
				{Key: "gateway", Label: "网关地址", Required: false, Placeholder: "默认 https://openapi.alipay.com/gateway.do", Secret: false},
			},
		},
		{
			Code:        string(domain.PluginCodeSMSBao),
			Name:        "短信宝短信",
			Description: "短信宝短信插件骨架，可用于手机号验证码和通知短信。",
			Fields: []domain.PluginField{
				{Key: "username", Label: "账号", Required: true, Placeholder: "短信宝账号", Secret: false},
				{Key: "password", Label: "密码 / MD5", Required: true, Placeholder: "短信宝密码或 MD5", Secret: true},
				{Key: "sign", Label: "短信签名", Required: false, Placeholder: "例如： 【白猿科技】", Secret: false},
			},
		},
		{
			Code:        string(domain.PluginCodeEPay),
			Name:        "易支付",
			Description: "易支付插件骨架，可用于钱包在线充值和支付回调联动。",
			Fields: []domain.PluginField{
				{Key: "merchantId", Label: "商户号", Required: true, Placeholder: "PID", Secret: false},
				{Key: "merchantKey", Label: "商户密钥", Required: true, Placeholder: "KEY", Secret: true},
				{Key: "apiURL", Label: "易支付地址", Required: true, Placeholder: "https://pay.example.com", Secret: false},
				{Key: "gateway", Label: "支付网关", Required: false, Placeholder: "留空则自动拼接为 submit.php", Secret: false},
				{Key: "notifyURL", Label: "异步回调", Required: false, Placeholder: "留空则按 SITE_URL 自动生成", Secret: false},
				{Key: "returnURL", Label: "同步回跳", Required: false, Placeholder: "留空则按 SITE_URL 自动生成", Secret: false},
			},
		},
	}
}
