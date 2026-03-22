package audit

import (
	"slices"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	nextID int64
	items  []Entry
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 3,
		items: []Entry{
			{
				ID:          1,
				ActorType:   "ADMIN",
				ActorID:     1,
				Actor:       "系统管理员",
				Action:      "customer.bootstrap",
				TargetType:  "customer",
				TargetID:    1,
				Target:      "CUST-20260320-0001",
				Description: "初始化演示客户与联系人数据",
				CreatedAt:   "2026-03-20 18:00:00",
			},
			{
				ID:          2,
				ActorType:   "ADMIN",
				ActorID:     1,
				Actor:       "系统管理员",
				Action:      "identity.pending",
				TargetType:  "customer",
				TargetID:    2,
				Target:      "CUST-20260320-0002",
				Description: "客户实名认证已提交，等待后台审核",
				CreatedAt:   "2026-03-20 18:05:00",
			},
		},
	}
}

func (repository *MemoryRepository) Record(entry Entry) Entry {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	entry.ID = repository.nextID
	repository.nextID++
	if entry.CreatedAt == "" {
		entry.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	repository.items = append([]Entry{entry}, repository.items...)
	return entry
}

func (repository *MemoryRepository) List() []Entry {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	return slices.Clone(repository.items)
}

func (repository *MemoryRepository) ListByTarget(targetType string, targetID int64) []Entry {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	results := make([]Entry, 0)
	for _, item := range repository.items {
		if item.TargetType == targetType && item.TargetID == targetID {
			results = append(results, item)
		}
	}
	return results
}
