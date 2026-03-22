package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	orderDomain "idc-finance/internal/modules/order/domain"
)

type MofangCatalogOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type MofangProvisionRequest struct {
	ServiceID        int64                                `json:"serviceId"`
	CustomerID       int64                                `json:"customerId"`
	ServiceNo        string                               `json:"serviceNo"`
	CustomerName     string                               `json:"customerName"`
	ProductName      string                               `json:"productName"`
	ProductType      string                               `json:"productType"`
	BillingCycle     string                               `json:"billingCycle"`
	ProviderNode     string                               `json:"providerNode"`
	ServerGroup      string                               `json:"serverGroup"`
	Configuration    []orderDomain.ServiceConfigSelection `json:"configuration"`
	ResourceSnapshot orderDomain.ServiceResourceSnapshot  `json:"resourceSnapshot"`
}

type MofangProvisionResult struct {
	ProviderType       string                              `json:"providerType"`
	ProviderResourceID string                              `json:"providerResourceId"`
	RegionName         string                              `json:"regionName"`
	IPAddress          string                              `json:"ipAddress"`
	Status             orderDomain.ServiceStatus           `json:"status"`
	SyncStatus         string                              `json:"syncStatus"`
	SyncMessage        string                              `json:"syncMessage"`
	LastAction         string                              `json:"lastAction"`
	ResourceSnapshot   orderDomain.ServiceResourceSnapshot `json:"resourceSnapshot"`
}

type mofangCustomerProfile struct {
	Name   string
	Email  string
	Mobile string
}

func (service *Service) ListMofangAreas() ([]MofangCatalogOption, error) {
	raw, err := service.request(http.MethodGet, "/v1/areas?sort=asc&list_type=all", nil)
	if err != nil {
		return nil, err
	}

	result := make([]MofangCatalogOption, 0)
	for _, record := range extractRecords(raw.Parsed) {
		value := pickString(record, "id")
		if value == "" {
			continue
		}
		label := strings.TrimSpace(pickString(record, "name", "short_name"))
		country := strings.TrimSpace(pickString(record, "country_code"))
		if label == "" {
			label = "未命名区域"
		}
		if country != "" {
			label = country + " / " + label
		}
		result = append(result, MofangCatalogOption{
			Value: value,
			Label: label,
		})
	}
	return result, nil
}

func (service *Service) ListMofangImages() ([]MofangCatalogOption, error) {
	raw, err := service.request(http.MethodGet, "/v1/image?per_page=200&sort=asc", nil)
	if err != nil {
		return nil, err
	}

	result := make([]MofangCatalogOption, 0)
	for _, record := range extractRecords(raw.Parsed) {
		value := pickString(record, "id")
		if value == "" {
			continue
		}
		if status := pickInt(record, "status"); status != 0 && status != 1 {
			continue
		}
		filename := strings.ToLower(pickString(record, "filename"))
		if strings.Contains(filename, "rescue") {
			continue
		}
		groupName := strings.TrimSpace(pickString(record, "group_name"))
		if groupName == "" {
			if groupRecord, ok := record["group"].(map[string]any); ok {
				groupName = strings.TrimSpace(pickString(groupRecord, "name"))
				if pickInt(groupRecord, "svg") == 10 {
					continue
				}
			}
		}
		name := strings.TrimSpace(pickString(record, "name"))
		if name == "" {
			continue
		}
		label := name
		if groupName != "" {
			label = groupName + " / " + name
		}
		result = append(result, MofangCatalogOption{
			Value: value,
			Label: label,
		})
	}
	return result, nil
}

func (service *Service) legacyProvisionMofangCloud(request MofangProvisionRequest) (MofangProvisionResult, error) {
	if strings.TrimSpace(service.config.BaseURL) == "" {
		return MofangProvisionResult{}, fmt.Errorf("未配置魔方云地址")
	}

	profile := service.lookupMofangCustomerProfile(request.CustomerID, request.CustomerName)
	userID, err := service.ensureMofangCloudUser(profile, request)
	if err != nil {
		service.logProvisionResult(request, "", "FAILED", err.Error(), nil, nil)
		return MofangProvisionResult{}, err
	}

	payload, err := buildMofangProvisionPayload(request, userID)
	if err != nil {
		service.logProvisionResult(request, "", "FAILED", err.Error(), nil, nil)
		return MofangProvisionResult{}, err
	}

	raw, err := service.request(http.MethodPost, "/v1/clouds", payload)
	if err != nil {
		service.logProvisionResult(request, "", "FAILED", err.Error(), payload, nil)
		return MofangProvisionResult{}, err
	}

	record := extractRecord(raw.Parsed)
	remoteID := pickString(record, "id", "cloudid", "cloud_id", "instance_id", "hostid")
	if remoteID == "" {
		service.logProvisionResult(request, "", "FAILED", "魔方云未返回实例编号", payload, record)
		return MofangProvisionResult{}, fmt.Errorf("魔方云未返回实例编号")
	}

	detail, detailErr := service.GetInstanceDetail(remoteID)
	snapshot := request.ResourceSnapshot
	if detailErr == nil {
		if strings.TrimSpace(detail.Name) != "" {
			snapshot.Hostname = detail.Name
		}
		if strings.TrimSpace(detail.Region) != "" {
			snapshot.RegionName = detail.Region
		}
		if strings.TrimSpace(detail.IPAddress) != "" {
			snapshot.PublicIPv4 = detail.IPAddress
		}
	}

	result := MofangProvisionResult{
		ProviderType:       "MOFANG_CLOUD",
		ProviderResourceID: remoteID,
		RegionName:         snapshot.RegionName,
		IPAddress:          snapshot.PublicIPv4,
		Status:             orderDomain.ServiceStatusPending,
		SyncStatus:         "PROVISIONING",
		SyncMessage:        "魔方云实例创建任务已提交",
		LastAction:         "provision",
		ResourceSnapshot:   snapshot,
	}

	if detailErr == nil {
		result.RegionName = firstNonEmpty(detail.Region, result.RegionName)
		result.IPAddress = firstNonEmpty(detail.IPAddress, result.IPAddress)
		result.Status = mapProvisionedServiceStatus(detail.Status)
		if result.Status == orderDomain.ServiceStatusActive {
			result.SyncStatus = "SUCCESS"
			result.SyncMessage = "魔方云实例已创建并同步"
		}
	}

	service.logProvisionResult(request, remoteID, result.SyncStatus, result.SyncMessage, payload, record)
	return result, nil
}

func (service *Service) lookupMofangCustomerProfile(customerID int64, fallbackName string) mofangCustomerProfile {
	profile := mofangCustomerProfile{
		Name: strings.TrimSpace(fallbackName),
	}
	if service.db == nil || customerID == 0 {
		return profile
	}

	row := service.db.QueryRowContext(context.Background(), `
SELECT c.name, c.email, c.mobile
FROM customers c
WHERE c.id = ?
LIMIT 1`, customerID)
	_ = row.Scan(&profile.Name, &profile.Email, &profile.Mobile)
	profile.Name = firstNonEmpty(strings.TrimSpace(profile.Name), strings.TrimSpace(fallbackName))
	return profile
}

func (service *Service) ensureMofangCloudUser(profile mofangCustomerProfile, request MofangProvisionRequest) (string, error) {
	username := buildMofangUsername(profile, request.CustomerID)
	password := buildMofangTempPassword(request.ServiceNo)

	_, _ = service.request(http.MethodPost, "/v1/user", map[string]any{
		"username":  username,
		"email":     profile.Email,
		"status":    1,
		"real_name": firstNonEmpty(profile.Name, request.CustomerName),
		"password":  password,
	})

	raw, err := service.request(http.MethodPost, "/v1/user/check", map[string]any{
		"username": username,
	})
	if err != nil {
		return "", err
	}

	record := extractRecord(raw.Parsed)
	userID := pickString(record, "id", "user_id", "uid")
	if userID == "" {
		userID = pickString(extractRecord(record["data"]), "id", "user_id", "uid")
	}
	if userID == "" {
		return "", fmt.Errorf("魔方云用户创建失败")
	}
	return userID, nil
}

func buildMofangProvisionPayload(request MofangProvisionRequest, userID string) (map[string]any, error) {
	configMap := make(map[string]string, len(request.Configuration))
	for _, item := range request.Configuration {
		if strings.TrimSpace(item.Code) == "" {
			continue
		}
		configMap[strings.TrimSpace(item.Code)] = strings.TrimSpace(item.Value)
	}

	area := firstNonEmpty(strings.TrimSpace(configMap["area"]), strings.TrimSpace(request.ResourceSnapshot.RegionName))
	if area == "" {
		return nil, fmt.Errorf("缺少魔方云区域配置")
	}
	image := strings.TrimSpace(configMap["os"])
	if image == "" {
		return nil, fmt.Errorf("缺少魔方云镜像配置")
	}

	cpu := pickConfigInt(configMap, "cpu", request.ResourceSnapshot.CPUCores, 2)
	memory := pickConfigInt(configMap, "memory", request.ResourceSnapshot.MemoryGB, 2)
	if memory > 0 && memory < 128 {
		memory *= 1024
	}
	systemDisk := pickConfigInt(configMap, "system_disk_size", request.ResourceSnapshot.SystemDiskGB, 40)
	dataDisk := pickConfigInt(configMap, "data_disk_size", request.ResourceSnapshot.DataDiskGB, 0)
	bandwidth := pickConfigInt(configMap, "bw", request.ResourceSnapshot.BandwidthMbps, 5)
	inboundBandwidth := pickConfigInt(configMap, "in_bw", bandwidth, bandwidth)
	outboundBandwidth := pickConfigInt(configMap, "out_bw", bandwidth, bandwidth)
	ipNum := pickConfigInt(configMap, "ip_num", 1, 1)
	backupNum := pickConfigInt(configMap, "backup_num", -1, -1)
	snapNum := pickConfigInt(configMap, "snap_num", -1, -1)
	flowLimit := pickConfigInt(configMap, "flow_limit", 0, 0)

	networkType := strings.TrimSpace(configMap["network_type"])
	if networkType == "" {
		networkType = "normal"
	}
	cloudType := strings.TrimSpace(configMap["type"])
	if cloudType == "" {
		cloudType = "host"
	}
	hostname := resolveMofangProvisionHostname(configMap, request)

	payload := map[string]any{
		"area":             area,
		"os":               image,
		"cpu":              cpu,
		"memory":           memory,
		"network_type":     networkType,
		"in_bw":            inboundBandwidth,
		"out_bw":           outboundBandwidth,
		"ip_num":           ipNum,
		"backup_num":       backupNum,
		"snap_num":         snapNum,
		"client":           userID,
		"hostname":         hostname,
		"rootpass":         buildMofangRootPassword(request.ServiceNo),
		"traffic_quota":    flowLimit,
		"traffic_type":     mapTrafficType(configMap["flow_way"]),
		"type":             cloudType,
		"system_disk_size": systemDisk,
	}

	for _, key := range []string{"node", "node_group", "ip_group", "node_priority", "advanced_cpu", "advanced_bw", "link_clone"} {
		if value := strings.TrimSpace(configMap[key]); value != "" {
			payload[key] = value
		}
	}
	if strings.TrimSpace(request.ProviderNode) != "" {
		if _, exists := payload["node"]; !exists {
			payload["node"] = strings.TrimSpace(request.ProviderNode)
		}
	}
	if strings.TrimSpace(request.ServerGroup) != "" {
		if _, exists := payload["node_group"]; !exists {
			payload["node_group"] = strings.TrimSpace(request.ServerGroup)
		}
	}
	for _, key := range []string{"nat_acl_limit", "nat_web_limit"} {
		if value := strings.TrimSpace(configMap[key]); value != "" {
			payload[key] = value
		}
	}

	if dataDisk > 0 {
		payload["other_data_disk[0][size]"] = dataDisk
	}
	if networkType == "vpc" {
		payload["vpc_name"] = buildMofangVPCName(request.ServiceNo)
	}
	if strings.EqualFold(strings.TrimSpace(configMap["traffic_bill_type"]), "last_30days") {
		payload["reset_flow_day"] = time.Now().Day()
	}

	return payload, nil
}

func (service *Service) logProvisionResult(
	request MofangProvisionRequest,
	remoteID,
	status,
	message string,
	requestBody any,
	responseBody any,
) {
	if service.db == nil || request.ServiceID == 0 {
		return
	}
	_ = service.insertSyncLog(context.Background(), syncLogInput{
		Action:       "provision_create",
		ResourceType: "instance",
		ResourceID:   firstNonEmpty(strings.TrimSpace(remoteID), request.ServiceNo),
		ServiceID:    request.ServiceID,
		Status:       firstNonEmpty(strings.TrimSpace(status), "SUCCESS"),
		Message:      message,
		RequestBody:  requestBody,
		ResponseBody: responseBody,
	})
}

func buildMofangUsername(profile mofangCustomerProfile, customerID int64) string {
	for _, candidate := range []string{
		strings.ToLower(strings.TrimSpace(profile.Email)),
		strings.TrimSpace(profile.Mobile),
	} {
		if candidate != "" {
			return candidate
		}
	}
	return fmt.Sprintf("idc-customer-%d", customerID)
}

func buildMofangTempPassword(serviceNo string) string {
	suffix := sanitizeProvisionToken(serviceNo)
	if len(suffix) > 8 {
		suffix = suffix[len(suffix)-8:]
	}
	if suffix == "" {
		suffix = "portal"
	}
	return "Portal!" + suffix + "1"
}

func buildMofangRootPassword(serviceNo string) string {
	suffix := sanitizeProvisionToken(serviceNo)
	if len(suffix) > 10 {
		suffix = suffix[len(suffix)-10:]
	}
	if suffix == "" {
		suffix = "rootpass"
	}
	return "Root!" + suffix + "1"
}

func buildMofangVPCName(serviceNo string) string {
	suffix := sanitizeProvisionToken(serviceNo)
	if len(suffix) > 10 {
		suffix = suffix[len(suffix)-10:]
	}
	if suffix == "" {
		suffix = "auto"
	}
	return "VPC-" + suffix
}

func resolveMofangProvisionHostname(configMap map[string]string, request MofangProvisionRequest) string {
	for _, key := range []string{"hostname", "host_name"} {
		if value := strings.TrimSpace(configMap[key]); value != "" {
			return value
		}
	}
	serviceHostname := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(request.ServiceNo), "_", "-"))
	if serviceHostname != "" {
		return serviceHostname
	}
	if value := strings.TrimSpace(request.ResourceSnapshot.Hostname); value != "" {
		return value
	}
	return fmt.Sprintf("srv-%d", request.ServiceID)
}

func sanitizeProvisionToken(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	for _, char := range value {
		switch {
		case char >= 'a' && char <= 'z':
			builder.WriteRune(char)
		case char >= '0' && char <= '9':
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

func pickConfigInt(config map[string]string, key string, fallback int, defaultValue int) int {
	if value := strings.TrimSpace(config[key]); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	if fallback > 0 || defaultValue < 0 {
		if fallback != 0 {
			return fallback
		}
	}
	return defaultValue
}

func mapTrafficType(value string) int {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "in":
		return 1
	case "out":
		return 2
	default:
		return 3
	}
}

func mapProvisionedServiceStatus(status string) orderDomain.ServiceStatus {
	switch normalizeStatus(status) {
	case "ACTIVE":
		return orderDomain.ServiceStatusActive
	case "SUSPENDED":
		return orderDomain.ServiceStatusSuspended
	case "TERMINATED":
		return orderDomain.ServiceStatusTerminated
	default:
		return orderDomain.ServiceStatusPending
	}
}
