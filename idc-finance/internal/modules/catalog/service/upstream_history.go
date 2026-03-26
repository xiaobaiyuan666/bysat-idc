package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"idc-finance/internal/modules/catalog/dto"
)

type upstreamImportHistoryState struct {
	Items []dto.UpstreamImportHistoryDetail `json:"items"`
}

func (service *Service) loadUpstreamImportHistory() {
	service.historyMu.Lock()
	defer service.historyMu.Unlock()

	service.history = upstreamImportHistoryState{Items: []dto.UpstreamImportHistoryDetail{}}
	if strings.TrimSpace(service.historyFile) == "" {
		return
	}

	content, err := os.ReadFile(service.historyFile)
	if err != nil {
		return
	}

	_ = json.Unmarshal(content, &service.history)
	if service.history.Items == nil {
		service.history.Items = []dto.UpstreamImportHistoryDetail{}
	}
}

func (service *Service) ListUpstreamImportHistory(query dto.UpstreamImportHistoryQuery) dto.UpstreamImportHistoryListResponse {
	service.historyMu.RLock()
	defer service.historyMu.RUnlock()

	limit := query.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	providerType := strings.ToUpper(strings.TrimSpace(query.ProviderType))
	status := strings.ToUpper(strings.TrimSpace(query.Status))
	keyword := strings.ToLower(strings.TrimSpace(query.Keyword))

	items := make([]dto.UpstreamImportHistoryRecord, 0, len(service.history.Items))
	for _, item := range service.history.Items {
		record := item.Record
		if query.ProviderAccountID > 0 && record.ProviderAccountID != query.ProviderAccountID {
			continue
		}
		if providerType != "" && !strings.EqualFold(record.ProviderType, providerType) {
			continue
		}
		if status != "" && !strings.EqualFold(record.Status, status) {
			continue
		}
		if keyword != "" {
			searchBlob := strings.ToLower(strings.Join([]string{
				record.AccountName,
				record.SourceName,
				record.Message,
				record.ProviderType,
				strings.Join(record.RequestedCodes, ","),
			}, " "))
			if !strings.Contains(searchBlob, keyword) {
				continue
			}
		}
		items = append(items, record)
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	if len(items) > limit {
		items = items[:limit]
	}

	return dto.UpstreamImportHistoryListResponse{
		Items: items,
		Total: len(items),
	}
}

func (service *Service) GetUpstreamImportHistory(historyID string) (dto.UpstreamImportHistoryDetail, bool) {
	service.historyMu.RLock()
	defer service.historyMu.RUnlock()

	for _, item := range service.history.Items {
		if item.Record.HistoryID == historyID {
			return item, true
		}
	}
	return dto.UpstreamImportHistoryDetail{}, false
}

func (service *Service) recordUpstreamImportHistory(
	request dto.ImportUpstreamProductsRequest,
	response dto.ImportUpstreamProductsResponse,
	runErr error,
	startedAt time.Time,
	finishedAt time.Time,
) string {
	service.historyMu.Lock()
	defer service.historyMu.Unlock()

	record := dto.UpstreamImportHistoryRecord{
		HistoryID:         buildUpstreamImportHistoryID(finishedAt),
		ProviderAccountID: request.ProviderAccountID,
		ProviderType:      strings.ToUpper(strings.TrimSpace(request.ProviderType)),
		Status:            deriveUpstreamImportRunStatus(response, runErr),
		Message:           deriveUpstreamImportRunMessage(response, runErr),
		ImportAll:         request.ImportAll,
		AutoSyncPricing:   request.AutoSyncPricing,
		AutoSyncConfig:    request.AutoSyncConfig,
		AutoSyncTemplate:  request.AutoSyncTemplate,
		DeactivateMissing: request.DeactivateMissing,
		RequestedCodes:    append([]string{}, request.RemoteProductCodes...),
		Created:           response.Created,
		Updated:           response.Updated,
		Deactivated:       response.DeactivatedCount,
		Failed:            response.Failed,
		Total:             len(response.Items),
		CreatedAt:         startedAt.Format("2006-01-02 15:04:05"),
		FinishedAt:        finishedAt.Format("2006-01-02 15:04:05"),
		DurationMs:        finishedAt.Sub(startedAt).Milliseconds(),
	}

	if record.ProviderType == "" {
		record.ProviderType = "ZJMF_API"
	}
	if service.provider != nil {
		if account, ok := service.provider.GetAccount(request.ProviderAccountID); ok {
			record.AccountName = strings.TrimSpace(account.Name)
			record.SourceName = strings.TrimSpace(account.SourceName)
			if record.ProviderType == "" {
				record.ProviderType = strings.ToUpper(strings.TrimSpace(account.ProviderType))
			}
		}
	}

	service.history.Items = append([]dto.UpstreamImportHistoryDetail{{
		Record: record,
		Items:  append([]dto.ImportUpstreamProductItem{}, response.Items...),
	}}, service.history.Items...)
	if len(service.history.Items) > 200 {
		service.history.Items = service.history.Items[:200]
	}
	_ = service.saveUpstreamImportHistoryLocked()
	return record.HistoryID
}

func (service *Service) saveUpstreamImportHistoryLocked() error {
	if strings.TrimSpace(service.historyFile) == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(service.historyFile), 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(service.history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(service.historyFile, content, 0o644)
}

func buildUpstreamImportHistoryID(finishedAt time.Time) string {
	return "UPSYNC-" + finishedAt.Format("20060102150405.000")
}

func deriveUpstreamImportRunStatus(response dto.ImportUpstreamProductsResponse, runErr error) string {
	if runErr != nil {
		return "FAILED"
	}
	successCount := response.Created + response.Updated + response.DeactivatedCount
	if response.Failed > 0 && successCount == 0 {
		return "FAILED"
	}
	if response.Failed > 0 {
		return "PARTIAL"
	}
	return "SUCCESS"
}

func deriveUpstreamImportRunMessage(response dto.ImportUpstreamProductsResponse, runErr error) string {
	if runErr != nil {
		return runErr.Error()
	}
	if response.Created > 0 || response.Updated > 0 || response.DeactivatedCount > 0 || response.Failed > 0 {
		return "上游商品同步完成"
	}
	if strings.TrimSpace(response.Message) != "" {
		return response.Message
	}
	return "本次同步没有产出变更"
}
