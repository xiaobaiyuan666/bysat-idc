package repository

import (
	"errors"

	"idc-finance/internal/modules/ticket/domain"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrServiceNotFound  = errors.New("service not found")
	ErrTicketNotFound   = errors.New("ticket not found")
	ErrTicketClosed     = errors.New("ticket is closed")
)

type Repository interface {
	List(filter domain.ListFilter) ([]domain.Ticket, int)
	ListByCustomer(customerID int64, filter domain.ListFilter) ([]domain.Ticket, int)
	GetByID(id int64) (domain.Ticket, bool)
	GetByCustomer(customerID, id int64) (domain.Ticket, bool)
	ListReplies(ticketID int64) []domain.Reply
	Create(input domain.CreateInput) (domain.Ticket, error)
	Update(id int64, input domain.UpdateInput) (domain.Ticket, bool, error)
	AddReply(ticketID int64, input domain.ReplyInput) (domain.Reply, domain.Ticket, bool, error)
	Close(ticketID int64) (domain.Ticket, bool, error)
}
