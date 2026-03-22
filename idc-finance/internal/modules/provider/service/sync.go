package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	automationService "idc-finance/internal/modules/automation/service"
	"idc-finance/internal/platform/audit"
)

type syncCounters struct {
	createdCustomers        int
	createdServices         int
	updatedServices         int
	failedServices          int
	syncedNetworkInterfaces int
	syncedIPs               int
	syncedDisks             int
	syncedSnapshots         int
	syncedBackups           int
	syncedVPCs              int
	syncedSecurityGroups    int
	syncedSecurityRules     int
}

type localServiceRow struct {
	ID                 int64
	ServiceNo          string
	CustomerID         int64
	ProductName        string
	Status             string
	ProviderType       string
	ProviderAccountID  int64
	ProviderResourceID string
	RegionName         string
	IPAddress          string
	SyncStatus         string
	SyncMessage        string
}

type customerRow struct {
	ID         int64
	CustomerNo string
	Name       string
	Email      string
}

func (service *Service) legacyPullSync(options PullSyncOptions) (PullSyncResult, error) {
	if err := service.legacyEnsureSyncReady(); err != nil {
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
			"includeResources": options.IncludeResources,
			"limit":            options.Limit,
			"remoteIds":        options.RemoteIDs,
		},
		Message: "批量同步任务已启动",
	})

	instances, err := service.ListInstances()
	if err != nil {
		service.failTask(taskID, err.Error(), map[string]any{
			"includeResources": options.IncludeResources,
			"limit":            options.Limit,
			"remoteIds":        options.RemoteIDs,
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
		syncItem, err := service.legacySyncInstance(item, options.IncludeResources, &counters)
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

func (service *Service) legacySyncServiceByID(serviceID int64, includeResources bool) (PullSyncItem, error) {
	if err := service.legacyEnsureSyncReady(); err != nil {
		return PullSyncItem{}, err
	}

	row := service.db.QueryRow(`
SELECT id, service_no, provider_resource_id
FROM services
WHERE id = ? AND provider_type = 'MOFANG_CLOUD'`, serviceID)

	var (
		serviceNo string
		remoteID  sql.NullString
	)
	if err := row.Scan(&serviceID, &serviceNo, &remoteID); err != nil {
		if err == sql.ErrNoRows {
			return PullSyncItem{}, fmt.Errorf("未找到魔方云服务")
		}
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
			"includeResources": includeResources,
		},
		Message: "单服务同步任务已启动",
	})

	instance, err := service.GetInstanceDetail(remoteID.String)
	if err != nil {
		service.failTask(taskID, err.Error(), map[string]any{
			"serviceId": serviceID,
			"remoteId":  remoteID.String,
		})
		return PullSyncItem{}, err
	}
	instance.RemoteID = remoteID.String

	item, err := service.legacySyncInstance(instance, includeResources, &syncCounters{})
	if err != nil {
		service.failTask(taskID, err.Error(), map[string]any{
			"serviceId": serviceID,
			"remoteId":  remoteID.String,
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

func (service *Service) legacySyncInstance(listItem InstanceSummary, includeResources bool, counters *syncCounters) (PullSyncItem, error) {
	detail, err := service.GetInstanceDetail(listItem.RemoteID)
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

	serviceRow, operation, err := service.upsertServiceTx(ctx, tx, customer, remote, listItem.RemoteID, 0)
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

func (service *Service) legacyEnsureSyncReady() error {
	if service.db == nil {
		return fmt.Errorf("当前运行实例未连接 MySQL，无法执行同步")
	}
	if strings.TrimSpace(service.config.BaseURL) == "" {
		return fmt.Errorf("未配置魔方云地址")
	}
	return nil
}

func (service *Service) ensureCustomerTx(ctx context.Context, tx *sql.Tx, remote map[string]any, remoteID string) (customerRow, bool, error) {
	userRecord := firstRecord(recordMap(remote["user"]), remote)
	customerKey := firstNonEmptyString(
		pickString(userRecord, "id", "user_id", "uid"),
		pickString(remote, "user_id", "uid"),
		remoteID,
	)
	customerNo := fmt.Sprintf("CUS-MF-%s", sanitizeToken(customerKey, 16))
	email := strings.ToLower(strings.TrimSpace(firstNonEmptyString(
		pickString(userRecord, "email"),
		pickString(remote, "email"),
	)))
	name := firstNonEmptyString(
		pickString(userRecord, "username", "name"),
		pickString(remote, "username", "hostname"),
		"魔方云用户 "+customerKey,
	)

	query := `
SELECT id, customer_no, name, email
FROM customers
WHERE customer_no = ?`
	args := []any{customerNo}
	if email != "" {
		query += ` OR email = ?`
		args = append(args, email)
	}
	query += ` ORDER BY id LIMIT 1`

	row := tx.QueryRowContext(ctx, query, args...)
	var existing customerRow
	if err := row.Scan(&existing.ID, &existing.CustomerNo, &existing.Name, &existing.Email); err == nil {
		nextEmail := existing.Email
		if email != "" {
			nextEmail = email
		}
		if _, err := tx.ExecContext(ctx, `
UPDATE customers
SET name = ?, email = ?, mobile = ?, remarks = ?, updated_at = NOW()
WHERE id = ?`,
			name,
			nextEmail,
			firstNonEmptyString(pickString(userRecord, "phone"), pickString(remote, "phone"), ""),
			"由魔方云同步引擎自动维护",
			existing.ID,
		); err != nil {
			return customerRow{}, false, err
		}
		existing.Name = name
		existing.Email = nextEmail
		return existing, false, nil
	} else if err != sql.ErrNoRows {
		return customerRow{}, false, err
	}

	if email == "" {
		uniqueEmail, ensureErr := service.buildUniqueCustomerEmailTx(ctx, tx, customerKey)
		if ensureErr != nil {
			return customerRow{}, false, ensureErr
		}
		email = uniqueEmail
	} else {
		uniqueEmail, ensureErr := service.ensureUniqueEmailTx(ctx, tx, email, 0)
		if ensureErr != nil {
			return customerRow{}, false, ensureErr
		}
		email = uniqueEmail
	}

	result, err := tx.ExecContext(ctx, `
INSERT INTO customers (customer_no, customer_type, name, email, mobile, status, sales_owner, remarks)
VALUES (?, 'PERSONAL', ?, ?, ?, 'ACTIVE', ?, ?)`,
		customerNo,
		name,
		email,
		firstNonEmptyString(pickString(userRecord, "phone"), pickString(remote, "phone"), ""),
		"魔方云同步",
		"由魔方云同步引擎自动创建",
	)
	if err != nil {
		return customerRow{}, false, err
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		return customerRow{}, false, err
	}
	return customerRow{
		ID:         customerID,
		CustomerNo: customerNo,
		Name:       name,
		Email:      email,
	}, true, nil
}

func (service *Service) buildUniqueCustomerEmailTx(ctx context.Context, tx *sql.Tx, customerKey string) (string, error) {
	base := fmt.Sprintf("mf-%s@sync.local", strings.ToLower(sanitizeToken(customerKey, 24)))
	return service.ensureUniqueEmailTx(ctx, tx, base, 0)
}

func (service *Service) ensureUniqueEmailTx(ctx context.Context, tx *sql.Tx, email string, excludeID int64) (string, error) {
	localPart := email
	domainPart := "sync.local"
	if parts := strings.SplitN(email, "@", 2); len(parts) == 2 {
		localPart = parts[0]
		domainPart = parts[1]
	}

	candidate := email
	suffix := 1
	for {
		row := tx.QueryRowContext(ctx, `SELECT id FROM customers WHERE email = ? LIMIT 1`, candidate)
		var id int64
		err := row.Scan(&id)
		if err == sql.ErrNoRows {
			return candidate, nil
		}
		if err != nil {
			return "", err
		}
		if excludeID != 0 && id == excludeID {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d@%s", localPart, suffix, domainPart)
		suffix++
	}
}

func (service *Service) upsertServiceTx(ctx context.Context, tx *sql.Tx, customer customerRow, remote map[string]any, remoteID string, providerAccountID int64) (localServiceRow, string, error) {
	existing, err := service.findServiceTx(ctx, tx, remoteID, providerAccountID)
	if err != nil {
		return localServiceRow{}, "", err
	}

	serviceName := firstNonEmptyString(
		pickString(remote, "hostname", "name", "remark"),
		fmt.Sprintf("魔方云实例 %s", remoteID),
	)
	productName := firstNonEmptyString(
		pickString(remote, "product_name", "product"),
		fmt.Sprintf("魔方云实例 / %s", serviceName),
	)
	status := normalizeStatus(firstNonEmptyString(pickString(remote, "status"), "ACTIVE"))
	regionName := firstNonEmptyString(
		pickString(remote, "area_name", "region_name", "node_name"),
		pickRegion(remote),
	)
	ipAddress := firstNonEmptyString(
		pickString(remote, "mainip", "main_ip", "ip", "ip_address"),
		firstIPFromRemote(remote),
	)

	configuration := buildConfigurationSelections(remote)
	resourceSnapshot := buildResourceSnapshot(remote)
	configurationJSON, _ := json.Marshal(configuration)
	resourceJSON, _ := json.Marshal(resourceSnapshot)
	nextDueAt := nullableRemoteTime(remote)
	syncMessage := "最近一次同步成功"
	lastAction := "provider-sync"

	if existing != nil {
		if _, err := tx.ExecContext(ctx, `
UPDATE services
SET customer_id = ?, product_name = ?, provider_account_id = ?, provider_resource_id = ?, region_name = ?, ip_address = ?, status = ?,
    sync_status = ?, sync_message = ?, next_due_at = ?, last_action = ?, configuration_snapshot = ?,
    resource_snapshot = ?, last_sync_at = NOW(), updated_at = NOW()
WHERE id = ?`,
			customer.ID,
			productName,
			providerAccountID,
			remoteID,
			regionName,
			ipAddress,
			status,
			"SUCCESS",
			syncMessage,
			nextDueAt,
			lastAction,
			configurationJSON,
			resourceJSON,
			existing.ID,
		); err != nil {
			return localServiceRow{}, "", err
		}
		existing.CustomerID = customer.ID
		existing.ProductName = productName
		existing.ProviderAccountID = providerAccountID
		existing.ProviderResourceID = remoteID
		existing.RegionName = regionName
		existing.IPAddress = ipAddress
		existing.Status = status
		existing.SyncStatus = "SUCCESS"
		existing.SyncMessage = syncMessage
		return *existing, "updated", nil
	}

	serviceNo, err := service.buildUniqueServiceNoTx(ctx, tx, fmt.Sprintf("SRV-MF-%s", sanitizeToken(remoteID, 18)))
	if err != nil {
		return localServiceRow{}, "", err
	}
	result, err := tx.ExecContext(ctx, `
INSERT INTO services (
  customer_id, service_no, product_name, provider_type, provider_account_id, provider_resource_id, region_name, ip_address, status,
  sync_status, sync_message, next_due_at, last_action, configuration_snapshot, resource_snapshot, last_sync_at
) VALUES (?, ?, ?, 'MOFANG_CLOUD', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		customer.ID,
		serviceNo,
		productName,
		providerAccountID,
		remoteID,
		regionName,
		ipAddress,
		status,
		"SUCCESS",
		syncMessage,
		nextDueAt,
		lastAction,
		configurationJSON,
		resourceJSON,
	)
	if err != nil {
		return localServiceRow{}, "", err
	}
	serviceID, err := result.LastInsertId()
	if err != nil {
		return localServiceRow{}, "", err
	}

	return localServiceRow{
		ID:                 serviceID,
		ServiceNo:          serviceNo,
		CustomerID:         customer.ID,
		ProductName:        productName,
		Status:             status,
		ProviderType:       "MOFANG_CLOUD",
		ProviderAccountID:  providerAccountID,
		ProviderResourceID: remoteID,
		RegionName:         regionName,
		IPAddress:          ipAddress,
		SyncStatus:         "SUCCESS",
		SyncMessage:        syncMessage,
	}, "created", nil
}

func (service *Service) findServiceTx(ctx context.Context, tx *sql.Tx, remoteID string, providerAccountID int64) (*localServiceRow, error) {
	query := `
SELECT id, service_no, customer_id, product_name, status, provider_type, provider_account_id, provider_resource_id, region_name, ip_address, sync_status, sync_message
FROM services
WHERE provider_type = 'MOFANG_CLOUD' AND provider_resource_id = ?`
	args := []any{remoteID}
	if providerAccountID > 0 {
		query += ` AND provider_account_id = ?`
		args = append(args, providerAccountID)
	}
	query += `
LIMIT 1`
	row := tx.QueryRowContext(ctx, query, args...)

	var item localServiceRow
	if err := row.Scan(
		&item.ID,
		&item.ServiceNo,
		&item.CustomerID,
		&item.ProductName,
		&item.Status,
		&item.ProviderType,
		&item.ProviderAccountID,
		&item.ProviderResourceID,
		&item.RegionName,
		&item.IPAddress,
		&item.SyncStatus,
		&item.SyncMessage,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (service *Service) buildUniqueServiceNoTx(ctx context.Context, tx *sql.Tx, base string) (string, error) {
	candidate := base
	sequence := 1
	for {
		row := tx.QueryRowContext(ctx, `SELECT id FROM services WHERE service_no = ? LIMIT 1`, candidate)
		var id int64
		err := row.Scan(&id)
		if err == sql.ErrNoRows {
			return candidate, nil
		}
		if err != nil {
			return "", err
		}
		candidate = fmt.Sprintf("%s-%d", base, sequence)
		sequence++
	}
}

func (service *Service) recordAudit(serviceRow localServiceRow, operation, customerName string) {
	if service.audit == nil {
		return
	}
	service.audit.Record(audit.Entry{
		ActorType:   "SYSTEM",
		ActorID:     0,
		Actor:       "魔方云同步引擎",
		Action:      "service.pull_sync",
		TargetType:  "service",
		TargetID:    serviceRow.ID,
		Target:      serviceRow.ServiceNo,
		Description: fmt.Sprintf("%s服务并同步到客户 %s", mapOperationText(operation), customerName),
		Payload: map[string]any{
			"providerType":       serviceRow.ProviderType,
			"providerResourceId": serviceRow.ProviderResourceID,
			"status":             serviceRow.Status,
			"operation":          operation,
		},
	})
}

func mapOperationText(operation string) string {
	if operation == "created" {
		return "创建"
	}
	return "更新"
}

func firstRecord(candidates ...map[string]any) map[string]any {
	for _, candidate := range candidates {
		if candidate != nil {
			return candidate
		}
	}
	return map[string]any{}
}

func mergeRecords(records ...map[string]any) map[string]any {
	merged := make(map[string]any)
	for _, record := range records {
		for key, value := range record {
			merged[key] = value
		}
	}
	return merged
}

func filterInstances(instances []InstanceSummary, remoteIDs []string) []InstanceSummary {
	if len(remoteIDs) == 0 {
		return instances
	}

	allowed := make(map[string]struct{}, len(remoteIDs))
	for _, remoteID := range remoteIDs {
		normalized := strings.TrimSpace(remoteID)
		if normalized == "" {
			continue
		}
		allowed[normalized] = struct{}{}
	}

	result := make([]InstanceSummary, 0, len(instances))
	for _, item := range instances {
		if _, ok := allowed[item.RemoteID]; ok {
			result = append(result, item)
		}
	}
	return result
}

func recordMap(value any) map[string]any {
	if record, ok := value.(map[string]any); ok {
		return record
	}
	return nil
}

func recordArray(value any) []map[string]any {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if record, ok := item.(map[string]any); ok {
			result = append(result, record)
		}
	}
	return result
}

func sanitizeToken(input string, limit int) string {
	var builder strings.Builder
	for _, char := range strings.ToUpper(strings.TrimSpace(input)) {
		switch {
		case char >= 'A' && char <= 'Z':
			builder.WriteRune(char)
		case char >= '0' && char <= '9':
			builder.WriteRune(char)
		case char == '-':
			builder.WriteRune(char)
		}
		if builder.Len() >= limit {
			break
		}
	}
	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return "SYNC"
	}
	return result
}

func nullableRemoteTime(remote map[string]any) any {
	for _, key := range []string{"next_due_date", "due_date", "renew_time", "expire_time", "expired_at", "end_time"} {
		value := strings.TrimSpace(pickString(remote, key))
		if value == "" || value == "0000-00-00 00:00:00" {
			continue
		}
		if parsed, ok := parseRemoteTime(value); ok {
			return parsed
		}
	}
	return nil
}

func parseRemoteTime(value string) (time.Time, bool) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		time.RFC3339,
	}
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func rollbackQuietly(tx *sql.Tx) {
	_ = tx.Rollback()
}
