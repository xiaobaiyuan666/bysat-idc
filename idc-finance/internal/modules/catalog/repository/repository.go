package repository

import "idc-finance/internal/modules/catalog/domain"

type Repository interface {
	List() []domain.Product
	ListActive() []domain.Product
	ListGroups() []string
	GetByID(id int64) (domain.Product, bool)
	Create(product domain.Product) domain.Product
	Update(id int64, product domain.Product) (domain.Product, bool)
}
