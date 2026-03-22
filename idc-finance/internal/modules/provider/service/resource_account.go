package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (service *Service) performResourceAction(
	remoteID,
	action string,
	request ResourceActionRequest,
	remote map[string]any,
	accountIDs ...int64,
) (map[string]any, string, string, error) {
	accountID := firstID(accountIDs...)
	if accountID == 0 {
		return service.legacyPerformResourceAction(remoteID, action, request, remote)
	}

	account, err := service.resolveMofangAccount(accountID)
	if err != nil {
		return nil, "", "", err
	}

	requester := func(method, path string, payload map[string]any) (rawResponse, error) {
		return service.requestWithAccount(account, method, path, payload)
	}

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
			return nil, "", "", fmt.Errorf("未识别到当前实例的 IPv4 段")
		}
		raw, err := requester(http.MethodPut, fmt.Sprintf("/v1/clouds/%s/ip", remoteID), map[string]any{
			"num":      targetCount,
			"ip_group": ipGroup,
		})
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), remoteID, fmt.Sprintf("已提交 IPv4 扩容，总数 %d", targetCount), nil

	case "add-ipv6":
		targetCount := pickInt(remote, "ipv6_num") + maxInt(request.Count, 1)
		raw, err := requester(http.MethodPut, fmt.Sprintf("/v1/clouds/%s/ipv6", remoteID), map[string]any{
			"num": targetCount,
		})
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), remoteID, fmt.Sprintf("已提交 IPv6 扩容，总数 %d", targetCount), nil

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
		raw, err := requester(http.MethodPost, fmt.Sprintf("/v1/clouds/%s/disks", remoteID), map[string]any{
			"size":   request.SizeGB,
			"store":  storeID,
			"driver": driver,
		})
		if err != nil {
			return nil, "", "", err
		}
		record := extractRecord(raw.Parsed)
		return record, pickString(record, "diskid", "id"), fmt.Sprintf("已提交新增 %d GB 数据盘", request.SizeGB), nil

	case "resize-disk":
		diskID := strings.TrimSpace(request.DiskID)
		if diskID == "" {
			return nil, "", "", fmt.Errorf("请选择需要扩容的磁盘")
		}
		if request.SizeGB <= 0 {
			return nil, "", "", fmt.Errorf("请填写扩容后的磁盘容量")
		}
		raw, err := requester(http.MethodPut, fmt.Sprintf("/v1/disks/%s", diskID), map[string]any{
			"size": request.SizeGB,
		})
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), diskID, fmt.Sprintf("已提交磁盘 %s 扩容到 %d GB", diskID, request.SizeGB), nil

	case "create-snapshot":
		diskID := strings.TrimSpace(request.DiskID)
		if diskID == "" {
			return nil, "", "", fmt.Errorf("请选择快照所属磁盘")
		}
		name := strings.TrimSpace(request.Name)
		if name == "" {
			name = "snapshot-" + time.Now().Format("20060102150405")
		}
		raw, err := requester(http.MethodPost, fmt.Sprintf("/v1/disks/%s/snapshots", diskID), map[string]any{
			"type": "snap",
			"name": name,
		})
		if err != nil {
			return nil, "", "", err
		}
		record := extractRecord(raw.Parsed)
		return record, firstNonEmptyString(pickString(record, "snapshotid", "id"), diskID), fmt.Sprintf("已创建快照 %s", name), nil

	case "delete-snapshot":
		snapshotID := strings.TrimSpace(request.SnapshotID)
		if snapshotID == "" {
			return nil, "", "", fmt.Errorf("请选择要删除的快照")
		}
		raw, err := requester(http.MethodDelete, fmt.Sprintf("/v1/snapshots/%s", snapshotID), nil)
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), snapshotID, fmt.Sprintf("已删除快照 %s", snapshotID), nil

	case "restore-snapshot":
		snapshotID := strings.TrimSpace(request.SnapshotID)
		if snapshotID == "" {
			return nil, "", "", fmt.Errorf("请选择要恢复的快照")
		}
		raw, err := requester(http.MethodPost, fmt.Sprintf("/v1/snapshots/%s/restore?hostid=%s", snapshotID, remoteID), nil)
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), snapshotID, fmt.Sprintf("已提交快照恢复 %s", snapshotID), nil

	case "create-backup":
		diskID := strings.TrimSpace(request.DiskID)
		if diskID == "" {
			return nil, "", "", fmt.Errorf("请选择备份所属磁盘")
		}
		name := strings.TrimSpace(request.Name)
		if name == "" {
			name = "backup-" + time.Now().Format("20060102150405")
		}
		raw, err := requester(http.MethodPost, fmt.Sprintf("/v1/disks/%s/snapshots", diskID), map[string]any{
			"type": "backup",
			"name": name,
		})
		if err != nil {
			return nil, "", "", err
		}
		record := extractRecord(raw.Parsed)
		return record, firstNonEmptyString(pickString(record, "snapshotid", "id"), diskID), fmt.Sprintf("已创建备份 %s", name), nil

	case "delete-backup":
		backupID := strings.TrimSpace(request.BackupID)
		if backupID == "" {
			return nil, "", "", fmt.Errorf("请选择要删除的备份")
		}
		raw, err := requester(http.MethodDelete, fmt.Sprintf("/v1/snapshots/%s", backupID), nil)
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), backupID, fmt.Sprintf("已删除备份 %s", backupID), nil

	case "restore-backup":
		backupID := strings.TrimSpace(request.BackupID)
		if backupID == "" {
			return nil, "", "", fmt.Errorf("请选择要恢复的备份")
		}
		raw, err := requester(http.MethodPost, fmt.Sprintf("/v1/snapshots/%s/restore?hostid=%s", backupID, remoteID), nil)
		if err != nil {
			return nil, "", "", err
		}
		return extractRecord(raw.Parsed), backupID, fmt.Sprintf("已提交备份恢复 %s", backupID), nil
	}

	return nil, "", "", fmt.Errorf("不支持的资源动作: %s", action)
}
