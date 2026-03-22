package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	automationService "idc-finance/internal/modules/automation/service"
)

type ResourceActionRequest struct {
	Count      int    `json:"count"`
	SizeGB     int    `json:"sizeGb"`
	StoreID    int    `json:"storeId"`
	IPGroup    int    `json:"ipGroup"`
	Driver     string `json:"driver"`
	Name       string `json:"name"`
	DiskID     string `json:"diskId"`
	SnapshotID string `json:"snapshotId"`
	BackupID   string `json:"backupId"`
}

type ResourceActionResponse struct {
	OK         bool          `json:"ok"`
	Action     string        `json:"action"`
	ServiceID  int64         `json:"serviceId"`
	ServiceNo  string        `json:"serviceNo"`
	RemoteID   string        `json:"remoteId"`
	ResourceID string        `json:"resourceId,omitempty"`
	Message    string        `json:"message"`
	PoweredOff bool          `json:"poweredOff"`
	PoweredOn  bool          `json:"poweredOn"`
	SyncItem   *PullSyncItem `json:"syncItem,omitempty"`
}

type serviceBinding struct {
	ID                 int64
	ServiceNo          string
	ProviderType       string
	ProviderAccountID  int64
	ProviderResourceID string
}

type syncLogInput struct {
	Action       string
	ResourceType string
	ResourceID   string
	ServiceID    int64
	Status       string
	Message      string
	RequestBody  any
	ResponseBody any
}

func (service *Service) GetServiceResources(serviceID int64) (ServiceResourcesResponse, error) {
	if service.db == nil {
		return ServiceResourcesResponse{}, fmt.Errorf("褰撳墠杩愯瀹炰緥鏈繛鎺?MySQL")
	}

	row := service.db.QueryRow(`
SELECT id, service_no, provider_type, provider_resource_id, status, region_name, ip_address, sync_status, sync_message,
       IFNULL(DATE_FORMAT(last_sync_at, '%Y-%m-%d %H:%i:%s'), '')
FROM services
WHERE id = ?`, serviceID)

	var response ServiceResourcesResponse
	if err := row.Scan(
		&response.ServiceID,
		&response.ServiceNo,
		&response.ProviderType,
		&response.ProviderResourceID,
		&response.Status,
		&response.RegionName,
		&response.IPAddress,
		&response.SyncStatus,
		&response.SyncMessage,
		&response.LastSyncAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return ServiceResourcesResponse{}, fmt.Errorf("未找到服务")
		}
		return ServiceResourcesResponse{}, err
	}

	interfaces, err := service.listNetworkInterfaces(serviceID)
	if err != nil {
		return ServiceResourcesResponse{}, err
	}
	ipAddresses, err := service.listIPAddresses(serviceID)
	if err != nil {
		return ServiceResourcesResponse{}, err
	}
	disks, err := service.listDisks(serviceID)
	if err != nil {
		return ServiceResourcesResponse{}, err
	}
	snapshots, err := service.listSnapshots(serviceID)
	if err != nil {
		return ServiceResourcesResponse{}, err
	}
	backups, err := service.listBackups(serviceID)
	if err != nil {
		return ServiceResourcesResponse{}, err
	}
	vpcs, err := service.listVPCNetworks(serviceID)
	if err != nil {
		return ServiceResourcesResponse{}, err
	}
	securityGroups, err := service.listSecurityGroups(serviceID)
	if err != nil {
		return ServiceResourcesResponse{}, err
	}
	syncLogs, err := service.ListSyncLogs(30, serviceID)
	if err != nil {
		return ServiceResourcesResponse{}, err
	}

	response.NetworkInterfaces = interfaces
	response.IPAddresses = ipAddresses
	response.Disks = disks
	response.Snapshots = snapshots
	response.Backups = backups
	response.VPCNetworks = vpcs
	response.SecurityGroups = securityGroups
	response.SyncLogs = syncLogs
	return response, nil
}

func (service *Service) ListSyncLogs(limit int, serviceID int64) ([]SyncLogItem, error) {
	if service.db == nil {
		return nil, fmt.Errorf("褰撳墠杩愯瀹炰緥鏈繛鎺?MySQL")
	}
	if limit <= 0 {
		limit = 50
	}

	query := `
SELECT id, provider_type, action, resource_type, resource_id, IFNULL(service_id, 0), status, message,
       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
FROM provider_sync_logs`
	args := make([]any, 0, 2)
	if serviceID > 0 {
		query += ` WHERE service_id = ?`
		args = append(args, serviceID)
	}
	query += ` ORDER BY id DESC LIMIT ?`
	args = append(args, limit)

	rows, err := service.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]SyncLogItem, 0)
	for rows.Next() {
		var item SyncLogItem
		if err := rows.Scan(
			&item.ID,
			&item.ProviderType,
			&item.Action,
			&item.ResourceType,
			&item.ResourceID,
			&item.ServiceID,
			&item.Status,
			&item.Message,
			&item.CreatedAt,
		); err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (service *Service) syncResourcesTx(ctx context.Context, tx *sql.Tx, serviceID int64, remote map[string]any, counters *syncCounters) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM service_security_group_rules WHERE security_group_id IN (SELECT id FROM service_security_groups WHERE service_id = ?)`, serviceID); err != nil {
		return err
	}
	for _, statement := range []string{
		`DELETE FROM service_security_groups WHERE service_id = ?`,
		`DELETE FROM service_vpc_networks WHERE service_id = ?`,
		`DELETE FROM service_snapshots WHERE service_id = ?`,
		`DELETE FROM service_backups WHERE service_id = ?`,
		`DELETE FROM service_ip_addresses WHERE service_id = ?`,
		`DELETE FROM service_network_interfaces WHERE service_id = ?`,
		`DELETE FROM service_disks WHERE service_id = ?`,
	} {
		if _, err := tx.ExecContext(ctx, statement, serviceID); err != nil {
			return err
		}
	}

	interfaceMap, err := service.insertNetworkInterfacesTx(ctx, tx, serviceID, remote, counters)
	if err != nil {
		return err
	}
	diskMap, err := service.insertDisksTx(ctx, tx, serviceID, remote, counters)
	if err != nil {
		return err
	}
	if err := service.insertIPAddressesTx(ctx, tx, serviceID, remote, interfaceMap, counters); err != nil {
		return err
	}
	if err := service.insertSnapshotsTx(ctx, tx, serviceID, remote, diskMap, counters); err != nil {
		return err
	}
	if err := service.insertBackupsTx(ctx, tx, serviceID, remote, counters); err != nil {
		return err
	}
	if err := service.insertVPCNetworksTx(ctx, tx, serviceID, remote, counters); err != nil {
		return err
	}
	if err := service.insertSecurityGroupsTx(ctx, tx, serviceID, remote, counters); err != nil {
		return err
	}
	return nil
}

func (service *Service) insertNetworkInterfacesTx(ctx context.Context, tx *sql.Tx, serviceID int64, remote map[string]any, counters *syncCounters) (map[string]int64, error) {
	result := make(map[string]int64)
	for _, item := range recordArray(remote["network"]) {
		res, err := tx.ExecContext(ctx, `
INSERT INTO service_network_interfaces (
  service_id, provider_interface_id, name, mac_address, bridge, network_name, interface_name, nic_model, inbound_mbps, outbound_mbps, vpc_name, status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'ACTIVE')`,
			serviceID,
			firstNonEmptyString(pickString(item, "id"), ""),
			firstNonEmptyString(pickString(item, "name"), firstNonEmptyString(pickString(item, "interface_name"), "eth0")),
			pickString(item, "mac"),
			pickString(item, "bridge"),
			pickString(item, "network_name", "net_name"),
			pickString(item, "interface_name"),
			pickString(item, "niccard_name"),
			pickInt(item, "in_bw"),
			pickInt(item, "out_bw"),
			pickString(item, "vpc_name"),
		)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		key := firstNonEmptyString(pickString(item, "id"), pickString(item, "name"), pickString(item, "interface_name"))
		if key != "" {
			result[key] = id
		}
		counters.syncedNetworkInterfaces += 1
	}
	return result, nil
}

func (service *Service) insertIPAddressesTx(ctx context.Context, tx *sql.Tx, serviceID int64, remote map[string]any, interfaceMap map[string]int64, counters *syncCounters) error {
	type ipRow struct {
		InterfaceID int64
		Address     string
		Version     string
		Role        string
		SubnetMask  string
		Gateway     string
		Bandwidth   int
		IsPrimary   bool
		Status      string
		ProviderID  string
	}

	ips := make([]ipRow, 0)
	seen := make(map[string]struct{})
	addIP := func(interfaceID int64, record map[string]any, primary bool, bandwidth int) {
		address := firstNonEmptyString(pickString(record, "ipaddress", "ip", "address"), "")
		if address == "" {
			return
		}
		if _, ok := seen[address]; ok {
			return
		}
		seen[address] = struct{}{}
		version := "IPv4"
		if strings.Contains(address, ":") {
			version = "IPv6"
		}
		ips = append(ips, ipRow{
			InterfaceID: interfaceID,
			Address:     address,
			Version:     version,
			Role:        "public",
			SubnetMask:  pickString(record, "subnet_mask"),
			Gateway:     pickString(record, "gateway"),
			Bandwidth:   bandwidth,
			IsPrimary:   primary,
			Status:      "ACTIVE",
			ProviderID:  pickString(record, "id", "ip_id"),
		})
	}

	mainIP := firstNonEmptyString(pickString(remote, "mainip", "main_ip"), "")
	for _, item := range recordArray(remote["ip"]) {
		addIP(0, item, firstNonEmptyString(pickString(item, "ipaddress"), "") == mainIP, 0)
	}
	for _, networkItem := range recordArray(remote["network"]) {
		interfaceKey := firstNonEmptyString(pickString(networkItem, "id"), pickString(networkItem, "name"), pickString(networkItem, "interface_name"))
		interfaceID := interfaceMap[interfaceKey]
		bandwidth := pickInt(networkItem, "out_bw", "in_bw")
		for _, ipItem := range recordArray(networkItem["ipaddress"]) {
			addIP(interfaceID, ipItem, firstNonEmptyString(pickString(ipItem, "ipaddress"), "") == mainIP, bandwidth)
		}
	}

	for _, item := range ips {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO service_ip_addresses (
  service_id, network_interface_id, provider_ip_id, address, version, ip_role, subnet_mask, gateway, bandwidth_mbps, is_primary, status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			serviceID,
			nullIfZero(item.InterfaceID),
			nullIfEmpty(item.ProviderID),
			item.Address,
			item.Version,
			item.Role,
			item.SubnetMask,
			item.Gateway,
			item.Bandwidth,
			item.IsPrimary,
			item.Status,
		); err != nil {
			return err
		}
		counters.syncedIPs += 1
	}
	return nil
}

func (service *Service) insertDisksTx(ctx context.Context, tx *sql.Tx, serviceID int64, remote map[string]any, counters *syncCounters) (map[string]int64, error) {
	result := make(map[string]int64)
	for _, item := range recordArray(remote["disk"]) {
		diskType := firstNonEmptyString(pickString(item, "type", "disk_type"), "data")
		name := firstNonEmptyString(pickString(item, "name"), "disk")
		res, err := tx.ExecContext(ctx, `
INSERT INTO service_disks (
  service_id, provider_disk_id, name, disk_type, size_gb, device_name, driver_name, fs_type, mount_point, is_system, status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			serviceID,
			nullIfEmpty(pickString(item, "id", "disk_id")),
			name,
			diskType,
			pickInt(item, "size", "size_gb"),
			pickString(item, "dev"),
			pickString(item, "driver"),
			pickString(item, "fs_type"),
			pickString(item, "mount_point", "path"),
			diskType == "system",
			normalizeDiskStatus(pickString(item, "status")),
		)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		result[name] = id
		if providerID := pickString(item, "id", "disk_id"); providerID != "" {
			result[providerID] = id
		}
		counters.syncedDisks += 1
	}
	return result, nil
}

func (service *Service) insertSnapshotsTx(ctx context.Context, tx *sql.Tx, serviceID int64, remote map[string]any, diskMap map[string]int64, counters *syncCounters) error {
	for _, item := range recordArray(remote["snapshot"]) {
		diskKey := firstNonEmptyString(pickString(item, "disk_id", "source_disk_id", "disk", "disk_name"), "")
		if _, err := tx.ExecContext(ctx, `
INSERT INTO service_snapshots (service_id, disk_id, provider_snapshot_id, name, size_gb, status)
VALUES (?, ?, ?, ?, ?, ?)`,
			serviceID,
			nullIfZero(diskMap[diskKey]),
			nullIfEmpty(pickString(item, "id", "snapshot_id")),
			firstNonEmptyString(pickString(item, "name"), "snapshot"),
			pickInt(item, "size", "size_gb"),
			firstNonEmptyString(pickString(item, "status"), "READY"),
		); err != nil {
			return err
		}
		counters.syncedSnapshots += 1
	}
	return nil
}

func (service *Service) insertBackupsTx(ctx context.Context, tx *sql.Tx, serviceID int64, remote map[string]any, counters *syncCounters) error {
	for _, item := range recordArray(remote["backup"]) {
		expiresAt := nullableTimeString(
			firstNonEmptyString(pickString(item, "expires_at", "expire_time", "end_time"), ""),
		)
		if _, err := tx.ExecContext(ctx, `
INSERT INTO service_backups (service_id, provider_backup_id, name, size_gb, status, expires_at)
VALUES (?, ?, ?, ?, ?, ?)`,
			serviceID,
			nullIfEmpty(pickString(item, "id", "backup_id")),
			firstNonEmptyString(pickString(item, "name"), "backup"),
			pickInt(item, "size", "size_gb"),
			firstNonEmptyString(pickString(item, "status"), "READY"),
			expiresAt,
		); err != nil {
			return err
		}
		counters.syncedBackups += 1
	}
	return nil
}

func (service *Service) insertVPCNetworksTx(ctx context.Context, tx *sql.Tx, serviceID int64, remote map[string]any, counters *syncCounters) error {
	for _, item := range recordArray(remote["network"]) {
		vpcName := firstNonEmptyString(pickString(item, "vpc_name"), "")
		if vpcName == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO service_vpc_networks (service_id, provider_vpc_id, name, cidr, gateway, interface_name, status)
VALUES (?, ?, ?, ?, ?, ?, 'ACTIVE')`,
			serviceID,
			nullIfEmpty(pickString(item, "vpc", "vpc_id")),
			vpcName,
			firstNonEmptyString(pickString(item, "cidr", "subnet"), ""),
			firstNonEmptyString(pickString(item, "gateway"), ""),
			pickString(item, "interface_name", "name"),
		); err != nil {
			return err
		}
		counters.syncedVPCs += 1
	}
	return nil
}

func (service *Service) insertSecurityGroupsTx(ctx context.Context, tx *sql.Tx, serviceID int64, remote map[string]any, counters *syncCounters) error {
	groupName := firstNonEmptyString(pickString(remote, "security_name"), "")
	if groupName == "" && pickInt(remote, "security") == 0 {
		return nil
	}
	res, err := tx.ExecContext(ctx, `
INSERT INTO service_security_groups (service_id, provider_security_group_id, name, status)
VALUES (?, ?, ?, 'ACTIVE')`,
		serviceID,
		nullIfEmpty(pickString(remote, "security", "security_group_id")),
		firstNonEmptyString(groupName, "默认安全组"),
	)
	if err != nil {
		return err
	}
	groupID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	counters.syncedSecurityGroups += 1

	for _, item := range recordArray(remote["security_rules"]) {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO service_security_group_rules (security_group_id, direction, protocol, port_range, source_cidr, action, description)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
			groupID,
			firstNonEmptyString(pickString(item, "direction"), "ingress"),
			firstNonEmptyString(pickString(item, "protocol"), "tcp"),
			firstNonEmptyString(pickString(item, "port_range", "port"), "all"),
			firstNonEmptyString(pickString(item, "source", "cidr"), "0.0.0.0/0"),
			firstNonEmptyString(pickString(item, "action"), "ALLOW"),
			firstNonEmptyString(pickString(item, "description"), ""),
		); err != nil {
			return err
		}
		counters.syncedSecurityRules += 1
	}
	return nil
}

func (service *Service) insertSyncLog(ctx context.Context, input syncLogInput) error {
	return service.insertSyncLogTx(ctx, service.db, input)
}

func (service *Service) insertSyncLogTx(ctx context.Context, exec sqlExecutor, input syncLogInput) error {
	var requestJSON any
	var responseJSON any
	if input.RequestBody != nil {
		encoded, _ := json.Marshal(input.RequestBody)
		requestJSON = string(encoded)
	}
	if input.ResponseBody != nil {
		encoded, _ := json.Marshal(input.ResponseBody)
		responseJSON = string(encoded)
	}
	_, err := exec.ExecContext(ctx, `
INSERT INTO provider_sync_logs (provider_type, action, resource_type, resource_id, service_id, status, message, request_body, response_body)
VALUES ('MOFANG_CLOUD', ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.Action,
		firstNonEmptyString(input.ResourceType, "instance"),
		input.ResourceID,
		nullIfZero(input.ServiceID),
		firstNonEmptyString(input.Status, "SUCCESS"),
		input.Message,
		requestJSON,
		responseJSON,
	)
	return err
}

func (service *Service) listNetworkInterfaces(serviceID int64) ([]NetworkInterface, error) {
	rows, err := service.db.Query(`
SELECT id, IFNULL(provider_interface_id, ''), name, mac_address, bridge, network_name, interface_name, nic_model, inbound_mbps, outbound_mbps, vpc_name, status
FROM service_network_interfaces
WHERE service_id = ?
ORDER BY id`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]NetworkInterface, 0)
	for rows.Next() {
		var item NetworkInterface
		if err := rows.Scan(
			&item.ID,
			&item.ProviderInterfaceID,
			&item.Name,
			&item.MACAddress,
			&item.Bridge,
			&item.NetworkName,
			&item.InterfaceName,
			&item.NICModel,
			&item.InboundMbps,
			&item.OutboundMbps,
			&item.VPCName,
			&item.Status,
		); err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (service *Service) listIPAddresses(serviceID int64) ([]IPAddress, error) {
	rows, err := service.db.Query(`
SELECT id, IFNULL(network_interface_id, 0), IFNULL(provider_ip_id, ''), address, version, ip_role, subnet_mask, gateway, bandwidth_mbps, is_primary, status
FROM service_ip_addresses
WHERE service_id = ?
ORDER BY is_primary DESC, id`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]IPAddress, 0)
	for rows.Next() {
		var item IPAddress
		if err := rows.Scan(
			&item.ID,
			&item.NetworkInterfaceID,
			&item.ProviderIPID,
			&item.Address,
			&item.Version,
			&item.IPRole,
			&item.SubnetMask,
			&item.Gateway,
			&item.BandwidthMbps,
			&item.IsPrimary,
			&item.Status,
		); err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (service *Service) listDisks(serviceID int64) ([]Disk, error) {
	rows, err := service.db.Query(`
SELECT id, IFNULL(provider_disk_id, ''), name, disk_type, size_gb, device_name, driver_name, fs_type, mount_point, is_system, status
FROM service_disks
WHERE service_id = ?
ORDER BY is_system DESC, id`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Disk, 0)
	for rows.Next() {
		var item Disk
		if err := rows.Scan(
			&item.ID,
			&item.ProviderDiskID,
			&item.Name,
			&item.DiskType,
			&item.SizeGB,
			&item.DeviceName,
			&item.DriverName,
			&item.FSType,
			&item.MountPoint,
			&item.IsSystem,
			&item.Status,
		); err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (service *Service) listSnapshots(serviceID int64) ([]Snapshot, error) {
	rows, err := service.db.Query(`
SELECT id, IFNULL(disk_id, 0), IFNULL(provider_snapshot_id, ''), name, size_gb, status
FROM service_snapshots
WHERE service_id = ?
ORDER BY id DESC`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Snapshot, 0)
	for rows.Next() {
		var item Snapshot
		if err := rows.Scan(
			&item.ID,
			&item.DiskID,
			&item.ProviderSnapshotID,
			&item.Name,
			&item.SizeGB,
			&item.Status,
		); err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (service *Service) listBackups(serviceID int64) ([]Backup, error) {
	rows, err := service.db.Query(`
SELECT id, IFNULL(provider_backup_id, ''), name, size_gb, status, IFNULL(DATE_FORMAT(expires_at, '%Y-%m-%d %H:%i:%s'), '')
FROM service_backups
WHERE service_id = ?
ORDER BY id DESC`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Backup, 0)
	for rows.Next() {
		var item Backup
		if err := rows.Scan(
			&item.ID,
			&item.ProviderBackupID,
			&item.Name,
			&item.SizeGB,
			&item.Status,
			&item.ExpiresAt,
		); err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (service *Service) listVPCNetworks(serviceID int64) ([]VPCNetwork, error) {
	rows, err := service.db.Query(`
SELECT id, IFNULL(provider_vpc_id, ''), name, cidr, gateway, interface_name, status
FROM service_vpc_networks
WHERE service_id = ?
ORDER BY id`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]VPCNetwork, 0)
	for rows.Next() {
		var item VPCNetwork
		if err := rows.Scan(
			&item.ID,
			&item.ProviderVPCID,
			&item.Name,
			&item.CIDR,
			&item.Gateway,
			&item.InterfaceName,
			&item.Status,
		); err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (service *Service) listSecurityGroups(serviceID int64) ([]SecurityGroup, error) {
	rows, err := service.db.Query(`
SELECT id, IFNULL(provider_security_group_id, ''), name, status
FROM service_security_groups
WHERE service_id = ?
ORDER BY id`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]SecurityGroup, 0)
	for rows.Next() {
		var item SecurityGroup
		if err := rows.Scan(
			&item.ID,
			&item.ProviderSecurityGroupID,
			&item.Name,
			&item.Status,
		); err != nil {
			return result, err
		}
		rules, err := service.listSecurityGroupRules(item.ID)
		if err != nil {
			return result, err
		}
		item.Rules = rules
		result = append(result, item)
	}
	return result, nil
}

func (service *Service) listSecurityGroupRules(groupID int64) ([]SecurityGroupRule, error) {
	rows, err := service.db.Query(`
SELECT id, direction, protocol, port_range, source_cidr, action, description
FROM service_security_group_rules
WHERE security_group_id = ?
ORDER BY id`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]SecurityGroupRule, 0)
	for rows.Next() {
		var item SecurityGroupRule
		if err := rows.Scan(
			&item.ID,
			&item.Direction,
			&item.Protocol,
			&item.PortRange,
			&item.SourceCIDR,
			&item.Action,
			&item.Description,
		); err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func nullIfZero(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

func nullableTimeString(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	if parsed, ok := parseRemoteTime(value); ok {
		return parsed
	}
	return nil
}

func normalizeDiskStatus(value string) string {
	if value == "" {
		return "ACTIVE"
	}
	switch value {
	case "1":
		return "ACTIVE"
	case "0":
		return "INACTIVE"
	default:
		return strings.ToUpper(value)
	}
}

func (service *Service) ExecuteServiceResourceAction(serviceID int64, action string, request ResourceActionRequest) (ResourceActionResponse, error) {
	binding, err := service.loadServiceBinding(serviceID)
	if err != nil {
		return ResourceActionResponse{}, err
	}
	if err := service.ensureSyncReady(binding.ProviderAccountID); err != nil {
		return ResourceActionResponse{}, err
	}
	if strings.TrimSpace(binding.ProviderType) != "MOFANG_CLOUD" {
		return ResourceActionResponse{}, fmt.Errorf("当前服务不是魔方云实例")
	}

	taskID := service.startTask(automationService.StartTaskRequest{
		TaskType:           "RESOURCE_ACTION",
		Title:              "魔方云资源动作",
		Channel:            "MOFANG_CLOUD",
		Stage:              "MANUAL",
		SourceType:         "service",
		SourceID:           binding.ID,
		ServiceID:          binding.ID,
		ServiceNo:          binding.ServiceNo,
		ProviderType:       binding.ProviderType,
		ProviderResourceID: binding.ProviderResourceID,
		ActionName:         action,
		OperatorType:       "ADMIN",
		OperatorName:       "资源工作台",
		RequestPayload:     request,
		Message:            "资源动作任务已启动",
	})

	detail, err := service.GetInstanceDetail(binding.ProviderResourceID, binding.ProviderAccountID)
	if err != nil {
		service.failTask(taskID, err.Error(), map[string]any{
			"serviceId": binding.ID,
			"remoteId":  binding.ProviderResourceID,
			"action":    action,
		})
		return ResourceActionResponse{}, err
	}
	remote := mergeRecords(detail.Raw)

	poweredOff := false
	poweredOn := false
	requiresPowerCycle := resourceActionRequiresPowerCycle(action, remote)
	wasRunning := isRemoteRunning(remote)
	if requiresPowerCycle && wasRunning {
		if _, err := service.ExecuteAction(binding.ProviderResourceID, "power-off", ActionRequest{}, binding.ProviderAccountID); err != nil {
			service.failTask(taskID, err.Error(), map[string]any{
				"serviceId": binding.ID,
				"remoteId":  binding.ProviderResourceID,
				"action":    action,
				"phase":     "power-off",
			})
			return ResourceActionResponse{}, err
		}
		poweredOff = true
		if err := service.waitForRemoteState(binding.ProviderResourceID, func(status string) bool {
			return status != "ACTIVE" && status != "PROVISIONING"
		}, 2*time.Minute, binding.ProviderAccountID); err != nil {
			service.failTask(taskID, err.Error(), map[string]any{
				"serviceId": binding.ID,
				"remoteId":  binding.ProviderResourceID,
				"action":    action,
				"phase":     "wait-power-off",
			})
			return ResourceActionResponse{}, err
		}
		if refreshed, refreshErr := service.GetInstanceDetail(binding.ProviderResourceID, binding.ProviderAccountID); refreshErr == nil {
			remote = mergeRecords(refreshed.Raw)
		}
	}

	responseBody, resourceID, message, err := service.performResourceAction(binding.ProviderResourceID, action, request, remote, binding.ProviderAccountID)
	if err != nil {
		if poweredOff && wasRunning {
			if _, powerErr := service.ExecuteAction(binding.ProviderResourceID, "power-on", ActionRequest{}, binding.ProviderAccountID); powerErr == nil {
				poweredOn = true
				_ = service.waitForRemoteState(binding.ProviderResourceID, func(status string) bool {
					return status == "ACTIVE"
				}, 2*time.Minute, binding.ProviderAccountID)
			}
		}
		_ = service.insertSyncLog(context.Background(), syncLogInput{
			Action:       "resource_" + action,
			ResourceType: resourceActionResourceType(action),
			ResourceID:   firstNonEmptyString(resourceID, binding.ProviderResourceID),
			ServiceID:    binding.ID,
			Status:       "FAILED",
			Message:      err.Error(),
			RequestBody:  request,
		})
		service.failTask(taskID, err.Error(), map[string]any{
			"serviceId":  binding.ID,
			"serviceNo":  binding.ServiceNo,
			"remoteId":   binding.ProviderResourceID,
			"resourceId": resourceID,
			"action":     action,
		})
		return ResourceActionResponse{}, err
	}

	if poweredOff && wasRunning {
		if _, err := service.ExecuteAction(binding.ProviderResourceID, "power-on", ActionRequest{}, binding.ProviderAccountID); err == nil {
			poweredOn = true
			_ = service.waitForRemoteState(binding.ProviderResourceID, func(status string) bool {
				return status == "ACTIVE"
			}, 2*time.Minute, binding.ProviderAccountID)
		}
	}

	syncItem, syncErr := service.SyncServiceByID(binding.ID, true)
	if syncErr != nil {
		message = message + "，但本地同步回写失败：" + syncErr.Error()
	}

	_ = service.insertSyncLog(context.Background(), syncLogInput{
		Action:       "resource_" + action,
		ResourceType: resourceActionResourceType(action),
		ResourceID:   firstNonEmptyString(resourceID, binding.ProviderResourceID),
		ServiceID:    binding.ID,
		Status:       "SUCCESS",
		Message:      message,
		RequestBody:  request,
		ResponseBody: responseBody,
	})

	result := ResourceActionResponse{
		OK:         true,
		Action:     action,
		ServiceID:  binding.ID,
		ServiceNo:  binding.ServiceNo,
		RemoteID:   binding.ProviderResourceID,
		ResourceID: resourceID,
		Message:    message,
		PoweredOff: poweredOff,
		PoweredOn:  poweredOn,
	}
	if syncErr == nil {
		result.SyncItem = &syncItem
	}
	service.successTask(taskID, "资源动作执行完成", result)
	return result, nil
}

func (service *Service) loadServiceBinding(serviceID int64) (serviceBinding, error) {
	if service.db == nil {
		return serviceBinding{}, fmt.Errorf("褰撳墠杩愯瀹炰緥鏈繛鎺?MySQL")
	}

row := service.db.QueryRow(`
SELECT id, service_no, provider_type, provider_account_id, IFNULL(provider_resource_id, '')
FROM services
WHERE id = ?`, serviceID)

	var binding serviceBinding
	if err := row.Scan(
		&binding.ID,
		&binding.ServiceNo,
		&binding.ProviderType,
		&binding.ProviderAccountID,
		&binding.ProviderResourceID,
	); err != nil {
		if err == sql.ErrNoRows {
			return serviceBinding{}, fmt.Errorf("未找到服务")
		}
		return serviceBinding{}, err
	}
	if strings.TrimSpace(binding.ProviderResourceID) == "" {
		return serviceBinding{}, fmt.Errorf("该服务缺少远程实例标识")
	}
	return binding, nil
}

func (service *Service) legacyPerformResourceAction(remoteID, action string, request ResourceActionRequest, remote map[string]any) (map[string]any, string, string, error) {
	switch action {
	case "add-ipv4":
		targetCount := pickInt(remote, "ip_num") + maxInt(request.Count, 1)
		ipGroup := request.IPGroup
		if ipGroup == 0 {
			for _, item := range recordArray(remote["ip"]) {
				ipGroup = pickInt(item, "ip_segment_id")
				if ipGroup > 0 {
					break
				}
			}
		}
		if ipGroup == 0 {
			return nil, "", "", fmt.Errorf("未识别到当前实例的 IP 段")
		}
		raw, err := service.request(http.MethodPut, fmt.Sprintf("/v1/clouds/%s/ip", remoteID), map[string]any{
			"num":      targetCount,
			"ip_group": ipGroup,
		})
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), remoteID, fmt.Sprintf("已向魔方云提交 IPv4 扩容，请求总数 %d", targetCount), nil

	case "add-ipv6":
		targetCount := pickInt(remote, "ipv6_num") + maxInt(request.Count, 1)
		raw, err := service.request(http.MethodPut, fmt.Sprintf("/v1/clouds/%s/ipv6", remoteID), map[string]any{
			"num": targetCount,
		})
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), remoteID, fmt.Sprintf("已向魔方云提交 IPv6 扩容，请求总数 %d", targetCount), nil

	case "add-disk":
		if request.SizeGB <= 0 {
			return nil, "", "", fmt.Errorf("请填写有效的数据盘容量")
		}
		storeID := request.StoreID
		if storeID == 0 {
			storeID = preferredStoreID(remote)
		}
		if storeID == 0 {
			return nil, "", "", fmt.Errorf("未识别到可用存储池")
		}
		driver := strings.TrimSpace(request.Driver)
		if driver == "" {
			driver = "virtio"
		}
		raw, err := service.request(http.MethodPost, fmt.Sprintf("/v1/clouds/%s/disks", remoteID), map[string]any{
			"size":   request.SizeGB,
			"store":  storeID,
			"driver": driver,
		})
		if err != nil {
			return nil, "", "", err
		}
		record := extractRecord(raw.Parsed)
		return record, pickString(record, "diskid", "id"), fmt.Sprintf("已向魔方云提交新增 %d GB 数据盘", request.SizeGB), nil

	case "resize-disk":
		diskID := strings.TrimSpace(request.DiskID)
		if diskID == "" {
			return nil, "", "", fmt.Errorf("璇烽€夋嫨瑕佹墿瀹圭殑纾佺洏")
		}
		if request.SizeGB <= 0 {
			return nil, "", "", fmt.Errorf("请填写扩容后的磁盘容量")
		}
		raw, err := service.request(http.MethodPut, fmt.Sprintf("/v1/disks/%s", diskID), map[string]any{
			"size": request.SizeGB,
		})
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), diskID, fmt.Sprintf("已向魔方云提交磁盘 %s 扩容到 %d GB", diskID, request.SizeGB), nil

	case "create-snapshot":
		diskID := strings.TrimSpace(request.DiskID)
		if diskID == "" {
			return nil, "", "", fmt.Errorf("璇烽€夋嫨瑕佸垱寤哄揩鐓х殑纾佺洏")
		}
		name := strings.TrimSpace(request.Name)
		if name == "" {
			name = "snapshot-" + time.Now().Format("20060102150405")
		}
		raw, err := service.request(http.MethodPost, fmt.Sprintf("/v1/disks/%s/snapshots", diskID), map[string]any{
			"type": "snap",
			"name": name,
		})
		if err != nil {
			return nil, "", "", err
		}
		record := extractRecord(raw.Parsed)
		return record, firstNonEmptyString(pickString(record, "snapshotid", "id"), diskID), fmt.Sprintf("宸插悜榄旀柟浜戞彁浜ゅ揩鐓у垱寤猴細%s", name), nil

	case "delete-snapshot":
		snapshotID := strings.TrimSpace(request.SnapshotID)
		if snapshotID == "" {
			return nil, "", "", fmt.Errorf("璇烽€夋嫨瑕佸垹闄ょ殑蹇収")
		}
		raw, err := service.request(http.MethodDelete, fmt.Sprintf("/v1/snapshots/%s", snapshotID), nil)
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), snapshotID, fmt.Sprintf("宸插悜榄旀柟浜戞彁浜ゅ揩鐓у垹闄わ細%s", snapshotID), nil

	case "restore-snapshot":
		snapshotID := strings.TrimSpace(request.SnapshotID)
		if snapshotID == "" {
			return nil, "", "", fmt.Errorf("璇烽€夋嫨瑕佹仮澶嶇殑蹇収")
		}
		raw, err := service.request(http.MethodPost, fmt.Sprintf("/v1/snapshots/%s/restore?hostid=%s", snapshotID, remoteID), nil)
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), snapshotID, fmt.Sprintf("宸插悜榄旀柟浜戞彁浜ゅ揩鐓ф仮澶嶏細%s", snapshotID), nil

	case "create-backup":
		diskID := strings.TrimSpace(request.DiskID)
		if diskID == "" {
			return nil, "", "", fmt.Errorf("璇烽€夋嫨瑕佸垱寤哄浠界殑纾佺洏")
		}
		name := strings.TrimSpace(request.Name)
		if name == "" {
			name = "backup-" + time.Now().Format("20060102150405")
		}
		raw, err := service.request(http.MethodPost, fmt.Sprintf("/v1/disks/%s/snapshots", diskID), map[string]any{
			"type": "backup",
			"name": name,
		})
		if err != nil {
			return nil, "", "", err
		}
		record := extractRecord(raw.Parsed)
		return record, firstNonEmptyString(pickString(record, "snapshotid", "id"), diskID), fmt.Sprintf("宸插悜榄旀柟浜戞彁浜ゅ浠藉垱寤猴細%s", name), nil

	case "delete-backup":
		backupID := strings.TrimSpace(request.BackupID)
		if backupID == "" {
			return nil, "", "", fmt.Errorf("璇烽€夋嫨瑕佸垹闄ょ殑澶囦唤")
		}
		raw, err := service.request(http.MethodDelete, fmt.Sprintf("/v1/snapshots/%s", backupID), nil)
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), backupID, fmt.Sprintf("宸插悜榄旀柟浜戞彁浜ゅ浠藉垹闄わ細%s", backupID), nil

	case "restore-backup":
		backupID := strings.TrimSpace(request.BackupID)
		if backupID == "" {
			return nil, "", "", fmt.Errorf("璇烽€夋嫨瑕佹仮澶嶇殑澶囦唤")
		}
		raw, err := service.request(http.MethodPost, fmt.Sprintf("/v1/snapshots/%s/restore?hostid=%s", backupID, remoteID), nil)
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), backupID, fmt.Sprintf("宸插悜榄旀柟浜戞彁浜ゅ浠芥仮澶嶏細%s", backupID), nil
	}

	return nil, "", "", fmt.Errorf("鏆備笉鏀寔鐨勮祫婧愬姩浣? %s", action)
}

func (service *Service) waitForRemoteState(remoteID string, matcher func(string) bool, timeout time.Duration, accountIDs ...int64) error {
	deadline := time.Now().Add(timeout)
	for {
		detail, err := service.GetInstanceDetail(remoteID, accountIDs...)
		if err == nil && matcher(strings.ToUpper(strings.TrimSpace(detail.Status))) {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("等待魔方云实例状态刷新超时")
		}
		time.Sleep(5 * time.Second)
	}
}

func resourceActionRequiresPowerCycle(action string, remote map[string]any) bool {
	switch action {
	case "add-ipv4":
		return strings.EqualFold(pickString(remote, "network_type"), "normal")
	case "add-ipv6", "resize-disk":
		return true
	default:
		return false
	}
}

func resourceActionResourceType(action string) string {
	switch action {
	case "add-ipv4", "add-ipv6":
		return "ip"
	case "add-disk", "resize-disk":
		return "disk"
	case "create-snapshot", "delete-snapshot", "restore-snapshot":
		return "snapshot"
	case "create-backup", "delete-backup", "restore-backup":
		return "backup"
	default:
		return "instance"
	}
}

func isRemoteRunning(remote map[string]any) bool {
	return normalizeStatus(pickString(remote, "status")) == "ACTIVE"
}

func preferredStoreID(remote map[string]any) int {
	for _, item := range recordArray(remote["disk"]) {
		if storeID := pickInt(item, "store_id"); storeID > 0 {
			return storeID
		}
	}
	return 0
}

func maxInt(current, fallback int) int {
	if current > fallback {
		return current
	}
	return fallback
}
