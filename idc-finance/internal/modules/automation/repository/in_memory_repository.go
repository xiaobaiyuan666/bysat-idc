package repository

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"idc-finance/internal/modules/automation/domain"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	items []domain.Task
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: []domain.Task{
			{
				ID:           1,
				TaskNo:       "TASK-20260321-0001",
				TaskType:     "PULL_SYNC",
				Title:        "魔方云实例拉取同步",
				Channel:      "MOFANG_CLOUD",
				Stage:        "SYNC",
				Status:       domain.TaskStatusSuccess,
				SourceType:   "service",
				SourceID:     72,
				CustomerID:   2,
				CustomerName: "云脉互联",
				ProductName:  "弹性云主机标准型",
				ServiceID:    72,
				ServiceNo:    "SRV-00000072",
				ProviderType: "MOFANG_CLOUD",
				Message:      "最近一次服务同步已完成",
				CreatedAt:    "2026-03-21 10:00:00",
				StartedAt:    "2026-03-21 10:00:01",
				FinishedAt:   "2026-03-21 10:00:08",
				UpdatedAt:    "2026-03-21 10:00:08",
			},
		},
	}
}

func (repository *MemoryRepository) List(filter domain.TaskFilter) []domain.Task {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]domain.Task, 0, len(repository.items))
	for _, item := range repository.items {
		if filter.Status != "" && string(item.Status) != filter.Status {
			continue
		}
		if filter.TaskType != "" && item.TaskType != filter.TaskType {
			continue
		}
		if filter.Channel != "" && item.Channel != filter.Channel {
			continue
		}
		if filter.Stage != "" && item.Stage != filter.Stage {
			continue
		}
		if filter.SourceType != "" && item.SourceType != filter.SourceType {
			continue
		}
		if filter.SourceID > 0 && item.SourceID != filter.SourceID {
			continue
		}
		if filter.OrderID > 0 && item.OrderID != filter.OrderID {
			continue
		}
		if filter.InvoiceID > 0 && item.InvoiceID != filter.InvoiceID {
			continue
		}
		if filter.ServiceID > 0 && item.ServiceID != filter.ServiceID {
			continue
		}
		if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
			if !strings.Contains(item.TaskNo, keyword) &&
				!strings.Contains(item.Title, keyword) &&
				!strings.Contains(item.ServiceNo, keyword) &&
				!strings.Contains(item.CustomerName, keyword) {
				continue
			}
		}
		items = append(items, item)
	}
	slices.Reverse(items)
	if filter.Limit > 0 && filter.Limit < len(items) {
		return items[:filter.Limit]
	}
	return items
}

func (repository *MemoryRepository) GetByID(id int64) (domain.Task, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	for _, item := range repository.items {
		if item.ID == id {
			return item, true
		}
	}
	return domain.Task{}, false
}

func (repository *MemoryRepository) Create(task domain.Task) domain.Task {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	task.ID = int64(len(repository.items) + 1)
	if strings.TrimSpace(task.TaskNo) == "" {
		task.TaskNo = fmt.Sprintf("TASK-%s-%04d", time.Now().Format("20060102"), task.ID)
	}
	repository.items = append(repository.items, task)
	return task
}

func (repository *MemoryRepository) Update(task domain.Task) (domain.Task, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for index, item := range repository.items {
		if item.ID != task.ID {
			continue
		}
		repository.items[index] = task
		return task, true
	}
	return domain.Task{}, false
}
