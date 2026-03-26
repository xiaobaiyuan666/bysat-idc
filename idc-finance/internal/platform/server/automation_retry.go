package server

import (
	"encoding/json"
	"fmt"
	"strings"

	automationDomain "idc-finance/internal/modules/automation/domain"
	orderDTO "idc-finance/internal/modules/order/dto"
	orderService "idc-finance/internal/modules/order/service"
	providerService "idc-finance/internal/modules/provider/service"
)

type automationRetryExecutor struct {
	provider *providerService.Service
	order    *orderService.Service
}

func newAutomationRetryExecutor(provider *providerService.Service, order *orderService.Service) *automationRetryExecutor {
	return &automationRetryExecutor{
		provider: provider,
		order:    order,
	}
}

func (executor *automationRetryExecutor) Execute(task automationDomain.Task) error {
	switch strings.ToUpper(strings.TrimSpace(task.TaskType)) {
	case "PULL_SYNC_SERVICE":
		if executor.provider == nil {
			return fmt.Errorf("魔方云提供方服务未初始化")
		}
		includeResources := true
		var payload struct {
			IncludeResources bool `json:"includeResources"`
		}
		if task.ServiceID <= 0 {
			return fmt.Errorf("任务缺少服务编号")
		}
		if strings.TrimSpace(task.RequestPayload) != "" {
			_ = parseTaskPayload(task.RequestPayload, &payload)
			includeResources = payload.IncludeResources
		}
		_, err := executor.provider.SyncServiceByID(task.ServiceID, includeResources)
		return err
	case "PULL_SYNC_BATCH":
		if executor.provider == nil {
			return fmt.Errorf("魔方云提供方服务未初始化")
		}
		var payload providerService.PullSyncOptions
		_ = parseTaskPayload(task.RequestPayload, &payload)
		if payload.Limit == 0 {
			payload.Limit = 50
		}
		_, err := executor.provider.PullSync(payload)
		return err
	case "SERVICE_ACTION":
		if executor.order == nil {
			return fmt.Errorf("订单服务未初始化")
		}
		if task.ServiceID <= 0 {
			return fmt.Errorf("任务缺少服务编号")
		}
		var payload orderDTO.ServiceActionRequest
		_ = parseTaskPayload(task.RequestPayload, &payload)
		_, ok := executor.order.ExecuteServiceAction(
			task.ServiceID,
			task.ActionName,
			payload,
			0,
			firstNonEmpty(task.OperatorName, "自动化任务重试"),
			"automation-retry",
		)
		if !ok {
			return fmt.Errorf("服务动作重试失败")
		}
		return nil
	case "RESOURCE_ACTION":
		if executor.provider == nil {
			return fmt.Errorf("魔方云提供方服务未初始化")
		}
		if task.ServiceID <= 0 {
			return fmt.Errorf("任务缺少服务编号")
		}
		var payload providerService.ResourceActionRequest
		_ = parseTaskPayload(task.RequestPayload, &payload)
		_, err := executor.provider.ExecuteServiceResourceAction(task.ServiceID, task.ActionName, payload)
		return err
	case "INVOICE_ACTION":
		if executor.order == nil {
			return fmt.Errorf("订单服务未初始化")
		}
		if task.InvoiceID <= 0 {
			return fmt.Errorf("任务缺少账单编号")
		}
		switch strings.TrimSpace(task.ActionName) {
		case "receive-payment":
			var payload orderDTO.ReceivePaymentRequest
			_ = parseTaskPayload(task.RequestPayload, &payload)
			_, _, _, ok, err := executor.order.ReceiveInvoicePayment(
				task.InvoiceID,
				payload,
				0,
				firstNonEmpty(task.OperatorName, "自动化任务重试"),
				"automation-retry",
			)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("账单收款重试失败")
			}
			return nil
		case "refund":
			var payload struct {
				Reason string `json:"reason"`
			}
			_ = parseTaskPayload(task.RequestPayload, &payload)
			_, _, _, ok := executor.order.RefundInvoice(
				task.InvoiceID,
				firstNonEmpty(payload.Reason, "自动化任务重试"),
				0,
				firstNonEmpty(task.OperatorName, "自动化任务重试"),
				"automation-retry",
			)
			if !ok {
				return fmt.Errorf("账单退款重试失败")
			}
			return nil
		default:
			return fmt.Errorf("暂不支持重试该账单任务动作")
		}
	default:
		return fmt.Errorf("暂不支持重试该类型任务")
	}
}

func parseTaskPayload(raw string, target any) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	return json.Unmarshal([]byte(raw), target)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
