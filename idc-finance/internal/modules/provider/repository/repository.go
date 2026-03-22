package repository

import "idc-finance/internal/modules/provider/domain"

type Repository interface {
	ListAccounts(providerType string) []domain.Account
	GetAccountByID(id int64) (domain.Account, bool)
	CreateAccount(account domain.Account) domain.Account
	UpdateAccount(id int64, account domain.Account) (domain.Account, bool)
	DeleteAccount(id int64) bool
	FindByProviderAndBaseURL(providerType, baseURL string) (domain.Account, bool)
}
