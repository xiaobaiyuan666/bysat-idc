package repository

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"idc-finance/internal/modules/provider/domain"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	nextID int64
	items  []domain.Account
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 1,
		items:  []domain.Account{},
	}
}

func (repository *MemoryRepository) ListAccounts(providerType string) []domain.Account {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	result := make([]domain.Account, 0, len(repository.items))
	for _, item := range repository.items {
		if providerType != "" && !strings.EqualFold(item.ProviderType, providerType) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func (repository *MemoryRepository) GetAccountByID(id int64) (domain.Account, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	for _, item := range repository.items {
		if item.ID == id {
			return item, true
		}
	}
	return domain.Account{}, false
}

func (repository *MemoryRepository) CreateAccount(account domain.Account) domain.Account {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05")
	account.ID = repository.nextID
	repository.nextID++
	account.CreatedAt = now
	account.UpdatedAt = now
	if account.Status == "" {
		account.Status = domain.AccountStatusActive
	}
	repository.items = append(repository.items, account)
	return account
}

func (repository *MemoryRepository) UpdateAccount(id int64, account domain.Account) (domain.Account, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.items {
		if item.ID != id {
			continue
		}
		account.ID = id
		account.CreatedAt = item.CreatedAt
		account.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		repository.items[index] = account
		return account, true
	}
	return domain.Account{}, false
}

func (repository *MemoryRepository) DeleteAccount(id int64) bool {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.items {
		if item.ID != id {
			continue
		}
		repository.items = append(repository.items[:index], repository.items[index+1:]...)
		return true
	}
	return false
}

func (repository *MemoryRepository) FindByProviderAndBaseURL(providerType, baseURL string) (domain.Account, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	for _, item := range repository.items {
		if strings.EqualFold(item.ProviderType, providerType) && normalizeBaseURL(item.BaseURL) == normalizeBaseURL(baseURL) {
			return item, true
		}
	}
	return domain.Account{}, false
}

func normalizeBaseURL(value string) string {
	return strings.TrimRight(strings.TrimSpace(value), "/")
}

func seedAccountName(providerType string, id int64) string {
	return fmt.Sprintf("%s-%d", providerType, id)
}
