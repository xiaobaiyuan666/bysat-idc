package repository

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"idc-finance/internal/modules/ticket/domain"
)

type MemoryRepository struct {
	mu           sync.RWMutex
	tickets      []domain.Ticket
	replies      []domain.Reply
	nextTicketID int64
	nextReplyID  int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tickets: []domain.Ticket{
			{
				ID:                 1,
				TicketNo:           "TIC-00000001",
				CustomerID:         1,
				CustomerNo:         "CUST-20260320-0001",
				CustomerName:       "Demo Customer",
				ServiceID:          1,
				ServiceNo:          "SRV-00000001",
				ProductName:        "CN2 Cloud Server",
				Title:              "Console login timeout",
				Content:            "The instance became unreachable after a maintenance window.",
				Status:             domain.StatusProcessing,
				Priority:           domain.PriorityHigh,
				Source:             "PORTAL",
				DepartmentName:     "Technical Support",
				AssignedAdminID:    1,
				AssignedAdminName:  "System Admin",
				LatestReplyExcerpt: "Please verify whether remote port 3389 is blocked.",
				LastReplyAt:        "2026-03-25 14:20:00",
				CreatedAt:          "2026-03-25 11:10:00",
				UpdatedAt:          "2026-03-25 14:20:00",
			},
			{
				ID:                 2,
				TicketNo:           "TIC-00000002",
				CustomerID:         1,
				CustomerNo:         "CUST-20260320-0001",
				CustomerName:       "Demo Customer",
				Title:              "Invoice amount confirmation",
				Content:            "Need a manual amount adjustment before payment.",
				Status:             domain.StatusOpen,
				Priority:           domain.PriorityNormal,
				Source:             "PORTAL",
				DepartmentName:     "Finance",
				LatestReplyExcerpt: "Need a manual amount adjustment before payment.",
				LastReplyAt:        "2026-03-25 09:05:00",
				CreatedAt:          "2026-03-25 09:05:00",
				UpdatedAt:          "2026-03-25 09:05:00",
			},
		},
		replies: []domain.Reply{
			{
				ID:         1,
				TicketID:   1,
				AuthorType: domain.AuthorTypeCustomer,
				AuthorID:   1,
				AuthorName: "Demo Customer",
				Content:    "The instance became unreachable after a maintenance window.",
				CreatedAt:  "2026-03-25 11:10:00",
				UpdatedAt:  "2026-03-25 11:10:00",
			},
			{
				ID:         2,
				TicketID:   1,
				AuthorType: domain.AuthorTypeAdmin,
				AuthorID:   1,
				AuthorName: "System Admin",
				Content:    "Please verify whether remote port 3389 is blocked.",
				CreatedAt:  "2026-03-25 14:20:00",
				UpdatedAt:  "2026-03-25 14:20:00",
			},
		},
		nextTicketID: 3,
		nextReplyID:  3,
	}
}

func (repository *MemoryRepository) List(filter domain.ListFilter) ([]domain.Ticket, int) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := repository.filterLocked(filter, 0)
	total := len(items)
	return paginateTickets(items, filter.Page, filter.Limit), total
}

func (repository *MemoryRepository) ListByCustomer(customerID int64, filter domain.ListFilter) ([]domain.Ticket, int) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := repository.filterLocked(filter, customerID)
	total := len(items)
	return paginateTickets(items, filter.Page, filter.Limit), total
}

func (repository *MemoryRepository) GetByID(id int64) (domain.Ticket, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	for _, item := range repository.tickets {
		if item.ID == id {
			return item, true
		}
	}
	return domain.Ticket{}, false
}

func (repository *MemoryRepository) GetByCustomer(customerID, id int64) (domain.Ticket, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	for _, item := range repository.tickets {
		if item.ID == id && item.CustomerID == customerID {
			return item, true
		}
	}
	return domain.Ticket{}, false
}

func (repository *MemoryRepository) ListReplies(ticketID int64) []domain.Reply {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	result := make([]domain.Reply, 0)
	for _, item := range repository.replies {
		if item.TicketID == ticketID {
			result = append(result, item)
		}
	}
	return result
}

func (repository *MemoryRepository) Create(input domain.CreateInput) (domain.Ticket, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05")
	item := domain.Ticket{
		ID:                 repository.nextTicketID,
		TicketNo:           fmt.Sprintf("TIC-%08d", repository.nextTicketID),
		CustomerID:         input.CustomerID,
		CustomerNo:         input.CustomerNo,
		CustomerName:       input.CustomerName,
		ServiceID:          input.ServiceID,
		ServiceNo:          serviceNoByID(input.ServiceID),
		ProductName:        productNameByServiceID(input.ServiceID),
		Title:              input.Title,
		Content:            input.Content,
		Status:             input.Status,
		Priority:           input.Priority,
		Source:             input.Source,
		DepartmentName:     input.DepartmentName,
		AssignedAdminID:    input.AssignedAdminID,
		AssignedAdminName:  input.AssignedAdmin,
		LatestReplyExcerpt: excerpt(input.Content),
		LastReplyAt:        now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	repository.nextTicketID++
	repository.tickets = append(repository.tickets, item)
	return item, nil
}

func (repository *MemoryRepository) Update(id int64, input domain.UpdateInput) (domain.Ticket, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.tickets {
		if item.ID != id {
			continue
		}
		if input.Title != nil {
			item.Title = strings.TrimSpace(*input.Title)
		}
		if input.Status != nil {
			item.Status = *input.Status
			if item.Status == domain.StatusClosed {
				item.ClosedAt = time.Now().Format("2006-01-02 15:04:05")
			} else {
				item.ClosedAt = ""
			}
		}
		if input.Priority != nil {
			item.Priority = *input.Priority
		}
		if input.DepartmentName != nil {
			item.DepartmentName = strings.TrimSpace(*input.DepartmentName)
		}
		if input.AssignedAdminID != nil {
			item.AssignedAdminID = *input.AssignedAdminID
			if *input.AssignedAdminID == 0 {
				item.AssignedAdminName = ""
			}
		}
		if input.AssignedAdminName != nil {
			item.AssignedAdminName = strings.TrimSpace(*input.AssignedAdminName)
		}
		item.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		repository.tickets[index] = item
		return item, true, nil
	}
	return domain.Ticket{}, false, nil
}

func (repository *MemoryRepository) AddReply(ticketID int64, input domain.ReplyInput) (domain.Reply, domain.Ticket, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.tickets {
		if item.ID != ticketID {
			continue
		}
		if item.Status == domain.StatusClosed && input.Status == "" {
			return domain.Reply{}, domain.Ticket{}, false, ErrTicketClosed
		}

		now := time.Now().Format("2006-01-02 15:04:05")
		reply := domain.Reply{
			ID:         repository.nextReplyID,
			TicketID:   ticketID,
			AuthorType: input.AuthorType,
			AuthorID:   input.AuthorID,
			AuthorName: input.AuthorName,
			Content:    input.Content,
			IsInternal: input.IsInternal,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		repository.nextReplyID++
		repository.replies = append(repository.replies, reply)

		if input.Status != "" {
			item.Status = input.Status
		}
		if item.Status == domain.StatusClosed {
			item.ClosedAt = now
		} else if input.Status != "" {
			item.ClosedAt = ""
		}
		item.LatestReplyExcerpt = excerpt(input.Content)
		item.LastReplyAt = now
		item.UpdatedAt = now
		repository.tickets[index] = item
		return reply, item, true, nil
	}
	return domain.Reply{}, domain.Ticket{}, false, nil
}

func (repository *MemoryRepository) Close(ticketID int64) (domain.Ticket, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.tickets {
		if item.ID != ticketID {
			continue
		}
		now := time.Now().Format("2006-01-02 15:04:05")
		item.Status = domain.StatusClosed
		item.ClosedAt = now
		item.UpdatedAt = now
		repository.tickets[index] = item
		return item, true, nil
	}
	return domain.Ticket{}, false, nil
}

func (repository *MemoryRepository) filterLocked(filter domain.ListFilter, onlyCustomerID int64) []domain.Ticket {
	keyword := strings.ToLower(strings.TrimSpace(filter.Keyword))
	status := strings.TrimSpace(filter.Status)
	priority := strings.TrimSpace(filter.Priority)
	department := strings.TrimSpace(filter.Department)
	result := make([]domain.Ticket, 0)

	for _, item := range repository.tickets {
		if onlyCustomerID > 0 && item.CustomerID != onlyCustomerID {
			continue
		}
		if filter.CustomerID > 0 && item.CustomerID != filter.CustomerID {
			continue
		}
		if filter.ServiceID > 0 && item.ServiceID != filter.ServiceID {
			continue
		}
		if status != "" && string(item.Status) != status {
			continue
		}
		if priority != "" && string(item.Priority) != priority {
			continue
		}
		if department != "" && !strings.EqualFold(item.DepartmentName, department) {
			continue
		}
		if filter.AdminID > 0 && item.AssignedAdminID != filter.AdminID {
			continue
		}
		if keyword != "" {
			text := strings.ToLower(strings.Join([]string{
				item.TicketNo,
				item.Title,
				item.CustomerName,
				item.ServiceNo,
				item.ProductName,
			}, " "))
			if !strings.Contains(text, keyword) {
				continue
			}
		}
		result = append(result, item)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].UpdatedAt > result[j].UpdatedAt
	})
	return result
}

func paginateTickets(items []domain.Ticket, page, limit int) []domain.Ticket {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	if start >= len(items) {
		return []domain.Ticket{}
	}
	end := start + limit
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func serviceNoByID(id int64) string {
	if id <= 0 {
		return ""
	}
	return fmt.Sprintf("SRV-%08d", id)
}

func productNameByServiceID(id int64) string {
	if id <= 0 {
		return ""
	}
	return "Cloud Service"
}

func excerpt(content string) string {
	content = strings.TrimSpace(content)
	if len(content) <= 120 {
		return content
	}
	return content[:120]
}
