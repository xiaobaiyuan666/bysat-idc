package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	automationService "idc-finance/internal/modules/automation/service"
)

func (service *Service) PullSync(options PullSyncOptions) (PullSyncResult, error) {
	if options.ProviderAccountID == 0 {
		return service.legacyPullSync(options)
	}

	if err := service.ensureSyncReady(options.ProviderAccountID); err != nil {
		return PullSyncResult{}, err
	}

	taskID := service.startTask(automationService.StartTaskRequest{
		TaskType:     "PULL_SYNC_BATCH",
		Title:        "魔方云实例批量同步",
		Channel:      "MOFANG_CLOUD",
		Stage:        "SYNC",
		SourceType:   "provider",
		ActionName:   "pull-sync",
		OperatorType: "SYSTEM",
		OperatorName: "魔方云同步引擎",
		RequestPayload: map[string]any{
			"providerAccountId": options.ProviderAccountID,
			"includeResources":  options.IncludeResources,
			"limit":             options.Limit,
			"remoteIds":         options.RemoteIDs,
		},
		Message: "批量同步任务已启动",
	})

	instances, err := service.ListInstances(options.ProviderAccountID)
	if err != nil {
		service.failTask(taskID, err.Error(), map[string]any{
			"providerAccountId": options.ProviderAccountID,
		})
		return PullSyncResult{}, err
	}

	filtered := filterInstances(instances, options.RemoteIDs)
	remoteCount := len(filtered)
	if options.Limit > 0 && options.Limit < len(filtered) {
		filtered = filtered[:options.Limit]
	}

	counters := syncCounters{}
	items := make([]PullSyncItem, 0, len(filtered))
	for _, item := range filtered {
		syncItem, err := service.syncInstanceWithAccount(item, options.ProviderAccountID, options.IncludeResources, &counters)
		if err != nil {
			counters.failedServices++
			_ = service.insertSyncLog(context.Background(), syncLogInput{
				Action:       "pull_sync_service",
				ResourceType: "instance",
				ResourceID:   item.RemoteID,
				Status:       "FAILED",
				Message:      err.Error(),
				ResponseBody: item.Raw,
			})
			items = append(items, PullSyncItem{
				RemoteID:  item.RemoteID,
				Operation: "failed",
				Status:    "FAILED",
				Message:   err.Error(),
			})
			continue
		}
		items = append(items, syncItem)
	}

	message := fmt.Sprintf("完成 %d 台实例同步", len(filtered))
	result := PullSyncResult{
		Summary: PullSyncSummary{
			RemoteCount:             remoteCount,
			ProcessedCount:          len(filtered),
			CreatedCustomers:        counters.createdCustomers,
			CreatedServices:         counters.createdServices,
			UpdatedServices:         counters.updatedServices,
			FailedServices:          counters.failedServices,
			SyncedNetworkInterfaces: counters.syncedNetworkInterfaces,
			SyncedIps:               counters.syncedIPs,
			SyncedDisks:             counters.syncedDisks,
			SyncedSnapshots:         counters.syncedSnapshots,
			SyncedBackups:           counters.syncedBackups,
			SyncedVpcs:              counters.syncedVPCs,
			SyncedSecurityGroups:    counters.syncedSecurityGroups,
			SyncedSecurityRules:     counters.syncedSecurityRules,
			Message:                 message,
		},
		Items: items,
	}
	service.successTask(taskID, message, result.Summary)
	return result, nil
}

func (service *Service) SyncServiceByID(serviceID int64, includeResources bool) (PullSyncItem, error) {
	row := service.db.QueryRow(`
SELECT id, service_no, provider_account_id, provider_resource_id
FROM services
WHERE id = ? AND provider_type = 'MOFANG_CLOUD'`, serviceID)

	var (
		serviceNo         string
		providerAccountID int64
		remoteID          sql.NullString
	)
	if err := row.Scan(&serviceID, &serviceNo, &providerAccountID, &remoteID); err != nil {
		if err == sql.ErrNoRows {
			return PullSyncItem{}, fmt.Errorf("未找到魔方云服务")
		}
		return PullSyncItem{}, err
	}

	if providerAccountID == 0 {
		return service.legacySyncServiceByID(serviceID, includeResources)
	}
	if err := service.ensureSyncReady(providerAccountID); err != nil {
		return PullSyncItem{}, err
	}
	if !remoteID.Valid || strings.TrimSpace(remoteID.String) == "" {
		return PullSyncItem{}, fmt.Errorf("该服务缺少远程实例标识")
	}

	taskID := service.startTask(automationService.StartTaskRequest{
		TaskType:           "PULL_SYNC_SERVICE",
		Title:              "魔方云单服务同步",
		Channel:            "MOFANG_CLOUD",
		Stage:              "SYNC",
		SourceType:         "service",
		SourceID:           serviceID,
		ServiceID:          serviceID,
		ServiceNo:          serviceNo,
		ProviderType:       "MOFANG_CLOUD",
		ProviderResourceID: strings.TrimSpace(remoteID.String),
		ActionName:         "pull-sync",
		OperatorType:       "SYSTEM",
		OperatorName:       "魔方云同步引擎",
		RequestPayload: map[string]any{
			"providerAccountId": providerAccountID,
			"includeResources":  includeResources,
		},
		Message: "单服务同步任务已启动",
	})

	instance, err := service.GetInstanceDetail(remoteID.String, providerAccountID)
	if err != nil {
		service.failTask(taskID, err.Error(), map[string]any{
			"serviceId":         serviceID,
			"providerAccountId": providerAccountID,
			"remoteId":          remoteID.String,
		})
		return PullSyncItem{}, err
	}
	instance.RemoteID = remoteID.String

	item, err := service.syncInstanceWithAccount(instance, providerAccountID, includeResources, &syncCounters{})
	if err != nil {
		service.failTask(taskID, err.Error(), map[string]any{
			"serviceId":         serviceID,
			"providerAccountId": providerAccountID,
			"remoteId":          remoteID.String,
		})
		return PullSyncItem{}, err
	}

	item.ServiceID = serviceID
	if item.ServiceNo == "" {
		item.ServiceNo = serviceNo
	}
	service.successTask(taskID, "单服务同步已完成", item)
	return item, nil
}

func (service *Service) syncInstanceWithAccount(listItem InstanceSummary, providerAccountID int64, includeResources bool, counters *syncCounters) (PullSyncItem, error) {
	detail, err := service.GetInstanceDetail(listItem.RemoteID, providerAccountID)
	if err != nil {
		return PullSyncItem{}, err
	}

	remote := mergeRecords(listItem.Raw, detail.Raw)
	ctx := context.Background()
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return PullSyncItem{}, err
	}
	defer rollbackQuietly(tx)

	customer, createdCustomer, err := service.ensureCustomerTx(ctx, tx, remote, listItem.RemoteID)
	if err != nil {
		return PullSyncItem{}, err
	}
	if createdCustomer {
		counters.createdCustomers++
	}

	serviceRow, operation, err := service.upsertServiceTx(ctx, tx, customer, remote, listItem.RemoteID, providerAccountID)
	if err != nil {
		return PullSyncItem{}, err
	}
	if operation == "created" {
		counters.createdServices++
	} else {
		counters.updatedServices++
	}

	if includeResources {
		if err := service.syncResourcesTx(ctx, tx, serviceRow.ID, remote, counters); err != nil {
			return PullSyncItem{}, err
		}
	}

	if err := service.insertSyncLogTx(ctx, tx, syncLogInput{
		Action:       "pull_sync_service",
		ResourceType: "instance",
		ResourceID:   listItem.RemoteID,
		ServiceID:    serviceRow.ID,
		Status:       "SUCCESS",
		Message:      fmt.Sprintf("%s服务 %s", mapOperationText(operation), serviceRow.ServiceNo),
		ResponseBody: remote,
	}); err != nil {
		return PullSyncItem{}, err
	}

	if err := tx.Commit(); err != nil {
		return PullSyncItem{}, err
	}

	service.recordAudit(serviceRow, operation, customer.Name)

	return PullSyncItem{
		RemoteID:     listItem.RemoteID,
		ServiceID:    serviceRow.ID,
		ServiceNo:    serviceRow.ServiceNo,
		CustomerID:   customer.ID,
		CustomerName: customer.Name,
		Operation:    operation,
		Status:       serviceRow.Status,
		Message:      fmt.Sprintf("%s到本地服务池", mapOperationText(operation)),
	}, nil
}

func (service *Service) ensureSyncReady(providerAccountID int64) error {
	if service.db == nil {
		return fmt.Errorf("当前运行实例未连接 MySQL，无法执行同步")
	}
	_, err := service.resolveMofangAccount(providerAccountID)
	return err
}
