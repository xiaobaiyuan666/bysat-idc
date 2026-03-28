package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"idc-finance/internal/modules/plugin/domain"
	"idc-finance/internal/modules/plugin/dto"
	"idc-finance/internal/modules/plugin/repository"
	"idc-finance/internal/platform/audit"
)

const smsBaoAPIURL = "http://api.smsbao.com/sms"

var smsBaoStatusMap = map[string]string{
	"0":  "发送成功",
	"30": "密码错误",
	"40": "账号不存在",
	"41": "余额不足",
	"43": "IP 地址受限",
	"50": "内容包含敏感词",
	"51": "手机号码格式不正确",
}

type OnPaymentSuccessFunc func(customerID int64, amount float64, orderNo string) error

type Service struct {
	repository       repository.Repository
	audit            *audit.Service
	siteURL          string
	onPaymentSuccess OnPaymentSuccessFunc
}

func New(repo repository.Repository, auditSvc *audit.Service) *Service {
	return &Service{
		repository: repo,
		audit:      auditSvc,
	}
}

func (service *Service) SetSiteURL(value string) {
	service.siteURL = strings.TrimRight(strings.TrimSpace(value), "/")
}

func (service *Service) SetOnPaymentSuccess(fn OnPaymentSuccessFunc) {
	service.onPaymentSuccess = fn
}

func (service *Service) ListDefinitions() []domain.PluginDefinition {
	return service.repository.ListDefinitions()
}

func (service *Service) ListConfigs() []domain.PluginConfig {
	return service.repository.ListConfigs()
}

func (service *Service) GetConfig(code string) (domain.PluginConfig, bool) {
	return service.repository.GetConfig(code)
}

func (service *Service) SaveConfig(code string, request dto.SavePluginConfigRequest, adminID int64, adminName, requestID string) (domain.PluginConfig, error) {
	item, err := service.repository.SaveConfig(code, request.Enabled, request.Config, adminName)
	if err != nil {
		return domain.PluginConfig{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "plugin.config.update",
		TargetType:  "plugin",
		TargetID:    0,
		Target:      item.Code,
		RequestID:   requestID,
		Description: "后台更新插件配置",
		Payload: map[string]any{
			"enabled": item.Enabled,
			"code":    item.Code,
		},
	})

	return item, nil
}

func (service *Service) VerifyKYC(pluginCode string, customerID int64, request dto.KYCVerifyRequest, requestID string) (domain.KYCRecord, error) {
	plugin, err := service.requireEnabledPlugin(pluginCode)
	if err != nil {
		return domain.KYCRecord{}, err
	}

	kycGateway := stringConfig(plugin.Config, "gateway", "gatewayURL")
	if kycGateway == "" {
		kycGateway = "https://openapi.alipay.com/gateway.do"
	}

	record, err := service.repository.CreateKYCRecord(domain.KYCRecord{
		PluginCode: plugin.Code,
		CustomerID: customerID,
		RealName:   strings.TrimSpace(request.RealName),
		IDCardNo:   strings.TrimSpace(request.IDCardNo),
		Mobile:     strings.TrimSpace(request.Mobile),
		Status:     "SUCCESS",
		Message:    "实名核验已完成，当前为本地模拟验证模式",
		Raw: map[string]any{
			"mock":    true,
			"plugin":  plugin.Code,
			"gateway": kycGateway,
			"ts":      time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		return domain.KYCRecord{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "PORTAL",
		ActorID:     customerID,
		Actor:       "客户",
		Action:      "plugin.kyc.verify",
		TargetType:  "plugin",
		TargetID:    0,
		Target:      plugin.Code,
		RequestID:   requestID,
		Description: "客户触发实名认证插件核验",
		Payload: map[string]any{
			"recordId": record.ID,
			"status":   record.Status,
		},
	})

	return record, nil
}

func (service *Service) DispatchSMS(pluginCode string, customerID int64, request dto.SMSDispatchRequest, requestID string) (domain.SMSRecord, error) {
	plugin, err := service.requireEnabledPlugin(pluginCode)
	if err != nil {
		return domain.SMSRecord{}, err
	}

	mobile := strings.TrimSpace(request.Mobile)
	content := strings.TrimSpace(request.Content)
	sign := stringConfig(plugin.Config, "sign")
	if sign != "" && !strings.HasPrefix(content, "【") {
		content = sign + content
	}

	username := stringConfig(plugin.Config, "username")
	password := stringConfig(plugin.Config, "password")

	status := "SUCCESS"
	message := ""
	raw := map[string]any{
		"plugin": plugin.Code,
		"ts":     time.Now().Format(time.RFC3339),
	}

	if username == "" || password == "" {
		message = "短信宝未配置账号密码，已按本地模拟发送成功处理"
		raw["mock"] = true
	} else {
		params := url.Values{}
		params.Set("u", username)
		params.Set("p", password)
		params.Set("m", mobile)
		params.Set("c", content)

		resp, httpErr := http.Get(smsBaoAPIURL + "?" + params.Encode())
		if httpErr != nil {
			status = "FAILED"
			message = fmt.Sprintf("短信宝请求失败: %v", httpErr)
			raw["error"] = httpErr.Error()
		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			code := strings.TrimSpace(string(body))
			raw["responseCode"] = code
			raw["httpStatus"] = resp.StatusCode

			if code == "0" {
				message = "短信发送成功"
			} else {
				status = "FAILED"
				if desc, ok := smsBaoStatusMap[code]; ok {
					message = fmt.Sprintf("短信宝返回错误: %s (code=%s)", desc, code)
				} else {
					message = fmt.Sprintf("短信宝返回未知状态: %s", code)
				}
			}
		}
	}

	record, err := service.repository.CreateSMSRecord(domain.SMSRecord{
		PluginCode:   plugin.Code,
		CustomerID:   customerID,
		Mobile:       mobile,
		TemplateCode: strings.TrimSpace(request.TemplateCode),
		Content:      content,
		Status:       status,
		Message:      message,
		BizNo:        fmt.Sprintf("SMS-%d", time.Now().UnixNano()),
		Raw:          raw,
	})
	if err != nil {
		return domain.SMSRecord{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "PORTAL",
		ActorID:     customerID,
		Actor:       "客户",
		Action:      "plugin.sms.dispatch",
		TargetType:  "plugin",
		TargetID:    0,
		Target:      plugin.Code,
		RequestID:   requestID,
		Description: "客户触发短信插件发送",
		Payload: map[string]any{
			"recordId": record.ID,
			"status":   record.Status,
			"mobile":   record.Mobile,
		},
	})

	return record, nil
}

func (service *Service) CreatePaymentOrder(pluginCode string, customerID int64, request dto.CreatePaymentOrderRequest, requestID string) (domain.PaymentOrder, error) {
	plugin, err := service.requireEnabledPlugin(pluginCode)
	if err != nil {
		return domain.PaymentOrder{}, err
	}
	if request.Amount <= 0 {
		return domain.PaymentOrder{}, fmt.Errorf("支付金额必须大于 0")
	}

	notifyURL := strings.TrimSpace(request.NotifyURL)
	if notifyURL == "" {
		notifyURL = stringConfig(plugin.Config, "notifyURL", "notify_url")
	}
	if notifyURL == "" && service.siteURL != "" {
		notifyURL = service.siteURL + "/api/v1/open/plugins/payment/" + plugin.Code + "/notify"
	}

	returnURL := strings.TrimSpace(request.ReturnURL)
	if returnURL == "" {
		returnURL = stringConfig(plugin.Config, "returnURL", "return_url")
	}
	if returnURL == "" && service.siteURL != "" {
		returnURL = service.siteURL + "/portal/"
	}

	orderNo := strings.TrimSpace(request.OrderNo)
	if orderNo == "" {
		orderNo = fmt.Sprintf("EPAY-%d", time.Now().UnixNano())
	}

	subject := firstNonEmpty(strings.TrimSpace(request.Subject), "钱包在线充值")
	money := fmt.Sprintf("%.2f", request.Amount)
	tradeNo := fmt.Sprintf("SIM-%d", time.Now().UnixNano())

	gateway := stringConfig(plugin.Config, "gateway", "gatewayURL")
	if gateway == "" {
		apiURL := stringConfig(plugin.Config, "apiURL", "api_url")
		if apiURL != "" {
			gateway = strings.TrimRight(apiURL, "/") + "/submit.php"
		}
	}

	merchantID := stringConfig(plugin.Config, "merchantId", "pid")
	merchantKey := stringConfig(plugin.Config, "merchantKey", "key")

	payURL := ""
	mockGateway := false
	realConfigReady := merchantID != "" && merchantKey != "" && gateway != ""
	if !realConfigReady {
		if !isLocalSiteURL(service.siteURL) {
			return domain.PaymentOrder{}, fmt.Errorf("易支付插件缺少生产配置，请先填写商户号、密钥和支付地址")
		}
		mockGateway = true
		query := url.Values{}
		query.Set("out_trade_no", orderNo)
		query.Set("trade_no", tradeNo)
		query.Set("money", money)
		query.Set("type", "alipay")
		query.Set("sign", "mock")
		payURL = notifyURL + "?" + query.Encode()
	} else {
		if gateway == "" {
			gateway = "https://epay.example.local/submit.php"
		}
		signParams := map[string]string{
			"pid":          merchantID,
			"type":         "alipay",
			"out_trade_no": orderNo,
			"notify_url":   notifyURL,
			"return_url":   returnURL,
			"name":         subject,
			"money":        money,
		}
		sign := epaySign(signParams, merchantKey)

		params := url.Values{}
		params.Set("pid", merchantID)
		params.Set("type", "alipay")
		params.Set("out_trade_no", orderNo)
		params.Set("notify_url", notifyURL)
		params.Set("return_url", returnURL)
		params.Set("name", subject)
		params.Set("money", money)
		params.Set("sign", sign)
		params.Set("sign_type", "MD5")
		payURL = strings.TrimRight(gateway, "/") + "?" + params.Encode()
	}

	order, err := service.repository.CreatePaymentOrder(domain.PaymentOrder{
		PluginCode: plugin.Code,
		OrderNo:    orderNo,
		CustomerID: customerID,
		Amount:     request.Amount,
		Currency:   firstNonEmpty(strings.TrimSpace(request.Currency), "CNY"),
		Subject:    subject,
		Status:     "UNPAID",
		PayURL:     payURL,
		NotifyURL:  notifyURL,
		ReturnURL:  returnURL,
		Raw: map[string]any{
			"mock":        mockGateway,
			"mockTradeNo": tradeNo,
			"plugin":      plugin.Code,
			"ts":          time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		return domain.PaymentOrder{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "PORTAL",
		ActorID:     customerID,
		Actor:       "客户",
		Action:      "plugin.epay.create_order",
		TargetType:  "plugin",
		TargetID:    0,
		Target:      plugin.Code,
		RequestID:   requestID,
		Description: "客户创建在线充值订单",
		Payload: map[string]any{
			"orderNo": order.OrderNo,
			"amount":  order.Amount,
			"status":  order.Status,
			"mock":    mockGateway,
		},
	})

	return order, nil
}

func (service *Service) NotifyPayment(pluginCode string, request dto.PaymentNotifyRequest, requestID string) (domain.PaymentOrder, bool, error) {
	plugin, err := service.requireEnabledPlugin(pluginCode)
	if err != nil {
		return domain.PaymentOrder{}, false, err
	}

	order, exists := service.repository.GetPaymentOrderByOrderNo(strings.TrimSpace(request.OrderNo))
	if !exists {
		return domain.PaymentOrder{}, false, nil
	}
	if strings.TrimSpace(request.Sign) == "mock" && !isLocalSiteURL(service.siteURL) {
		return domain.PaymentOrder{}, false, fmt.Errorf("production site does not allow mock payment notify")
	}

	order.PluginCode = plugin.Code
	order.Status = strings.ToUpper(strings.TrimSpace(request.Status))
	order.TradeNo = strings.TrimSpace(request.TradeNo)
	if order.Status == "PAID" || order.Status == "SUCCESS" {
		order.Status = "PAID"
		order.PaidAt = time.Now().Format("2006-01-02 15:04:05")
	}
	order.Raw = map[string]any{
		"notify": map[string]any{
			"status":  order.Status,
			"tradeNo": order.TradeNo,
			"sign":    strings.TrimSpace(request.Sign),
		},
		"mock": true,
	}

	updated, ok, err := service.repository.UpdatePaymentOrder(order)
	if err != nil {
		return domain.PaymentOrder{}, false, err
	}
	if !ok {
		return domain.PaymentOrder{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "SYSTEM",
		ActorID:     0,
		Actor:       "EPAY 回调",
		Action:      "plugin.epay.notify",
		TargetType:  "plugin",
		TargetID:    0,
		Target:      plugin.Code,
		RequestID:   requestID,
		Description: "在线充值回调处理完成",
		Payload: map[string]any{
			"orderNo": updated.OrderNo,
			"tradeNo": updated.TradeNo,
			"status":  updated.Status,
		},
	})

	if updated.Status == "PAID" && service.onPaymentSuccess != nil && updated.CustomerID > 0 {
		if cbErr := service.onPaymentSuccess(updated.CustomerID, updated.Amount, updated.OrderNo); cbErr != nil {
			service.audit.Record(audit.Entry{
				ActorType:   "SYSTEM",
				ActorID:     0,
				Actor:       "EPAY 回调",
				Action:      "plugin.epay.recharge_failed",
				TargetType:  "customer",
				TargetID:    updated.CustomerID,
				Target:      updated.OrderNo,
				RequestID:   requestID,
				Description: fmt.Sprintf("支付成功但钱包入账失败: %v", cbErr),
				Payload: map[string]any{
					"orderNo": updated.OrderNo,
					"amount":  updated.Amount,
					"error":   cbErr.Error(),
				},
			})
		}
	}

	return updated, true, nil
}

func (service *Service) ListKYCRecords(customerID int64, limit int) []domain.KYCRecord {
	if limit <= 0 {
		limit = 20
	}
	return service.repository.ListKYCRecords(customerID, limit)
}

func (service *Service) ListSMSRecords(customerID int64, limit int) []domain.SMSRecord {
	if limit <= 0 {
		limit = 20
	}
	return service.repository.ListSMSRecords(customerID, limit)
}

func (service *Service) ListPaymentOrders(customerID int64, limit int) []domain.PaymentOrder {
	if limit <= 0 {
		limit = 20
	}
	return service.repository.ListPaymentOrders(customerID, limit)
}

func (service *Service) requireEnabledPlugin(code string) (domain.PluginConfig, error) {
	item, exists := service.repository.GetConfig(code)
	if !exists {
		return domain.PluginConfig{}, fmt.Errorf("插件不存在: %s", strings.TrimSpace(code))
	}
	if !item.Enabled {
		return domain.PluginConfig{}, fmt.Errorf("插件未启用: %s", item.Code)
	}
	return item, nil
}

func stringConfig(config map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := config[key]
		if !ok {
			continue
		}
		if text := strings.TrimSpace(fmt.Sprint(value)); text != "" {
			return text
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func epaySign(params map[string]string, merchantKey string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		if key == "sign" || key == "sign_type" || params[key] == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for index, key := range keys {
		if index > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(params[key])
	}
	builder.WriteString(merchantKey)

	hash := md5.Sum([]byte(builder.String()))
	return hex.EncodeToString(hash[:])
}

func isLocalSiteURL(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	return lower == "" ||
		strings.Contains(lower, "127.0.0.1") ||
		strings.Contains(lower, "localhost")
}
