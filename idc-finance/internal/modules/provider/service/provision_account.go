package service

import (
	"fmt"
	"net/http"
	"strings"

	orderDomain "idc-finance/internal/modules/order/domain"
)

func (service *Service) ProvisionMofangCloud(request MofangProvisionRequest, accountIDs ...int64) (MofangProvisionResult, error) {
	accountID := firstID(accountIDs...)
	if accountID == 0 {
		return service.legacyProvisionMofangCloud(request)
	}

	account, err := service.resolveMofangAccount(accountID)
	if err != nil {
		return MofangProvisionResult{}, err
	}

	profile := service.lookupMofangCustomerProfile(request.CustomerID, request.CustomerName)
	userID, err := service.ensureMofangCloudUserWithAccount(account, profile, request)
	if err != nil {
		service.logProvisionResult(request, "", "FAILED", err.Error(), nil, nil)
		return MofangProvisionResult{}, err
	}

	payload, err := buildMofangProvisionPayload(request, userID)
	if err != nil {
		service.logProvisionResult(request, "", "FAILED", err.Error(), nil, nil)
		return MofangProvisionResult{}, err
	}

	raw, err := service.requestWithAccount(account, http.MethodPost, "/v1/clouds", payload)
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

	detail, detailErr := service.GetInstanceDetail(remoteID, account.AccountID)
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

func (service *Service) ensureMofangCloudUserWithAccount(
	account resolvedProviderAccount,
	profile mofangCustomerProfile,
	request MofangProvisionRequest,
) (string, error) {
	username := buildMofangUsername(profile, request.CustomerID)
	password := buildMofangTempPassword(request.ServiceNo)

	_, _ = service.requestWithAccount(account, http.MethodPost, "/v1/user", map[string]any{
		"username":  username,
		"email":     profile.Email,
		"status":    1,
		"real_name": firstNonEmpty(profile.Name, request.CustomerName),
		"password":  password,
	})

	raw, err := service.requestWithAccount(account, http.MethodPost, "/v1/user/check", map[string]any{
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
