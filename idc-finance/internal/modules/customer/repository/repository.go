package repository

import (
	"errors"

	"idc-finance/internal/modules/customer/domain"
)

var (
	ErrGroupNotFound = errors.New("customer group not found")
	ErrGroupInUse    = errors.New("customer group in use")
	ErrLevelNotFound = errors.New("customer level not found")
	ErrLevelInUse    = errors.New("customer level in use")
)

type Repository interface {
	List() []domain.Customer
	ListGroups() []domain.CustomerGroup
	ListLevels() []domain.CustomerLevel
	CreateGroup(group domain.CustomerGroup) (domain.CustomerGroup, error)
	UpdateGroup(id int64, group domain.CustomerGroup) (domain.CustomerGroup, error)
	DeleteGroup(id int64) error
	CreateLevel(level domain.CustomerLevel) (domain.CustomerLevel, error)
	UpdateLevel(id int64, level domain.CustomerLevel) (domain.CustomerLevel, error)
	DeleteLevel(id int64) error
	GetByID(id int64) (domain.Customer, bool)
	Create(customer domain.Customer) domain.Customer
	Update(id int64, customer domain.Customer) (domain.Customer, bool)
	AddContact(customerID int64, contact domain.Contact) (domain.Customer, domain.Contact, bool)
	UpdateContact(customerID, contactID int64, contact domain.Contact) (domain.Customer, domain.Contact, bool)
	DeleteContact(customerID, contactID int64) (domain.Customer, bool)
	SubmitIdentity(customerID int64, identity domain.Identity) (domain.Customer, bool)
	ReviewIdentity(customerID int64, status domain.IdentityStatus, reviewRemark string) (domain.Customer, bool)
	ListServiceItems(customerID int64) []domain.RelatedItem
	ListInvoiceItems(customerID int64) []domain.RelatedItem
	ListTicketItems(customerID int64) []domain.RelatedItem
}
