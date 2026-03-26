package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	automationService "idc-finance/internal/modules/automation/service"
	customerService "idc-finance/internal/modules/customer/service"
	orderService "idc-finance/internal/modules/order/service"
	"idc-finance/internal/modules/ticket/domain"
	"idc-finance/internal/modules/ticket/dto"
	"idc-finance/internal/modules/ticket/repository"
	"idc-finance/internal/platform/audit"
)

const (
	presetReplyConfigKey = "ticket.preset_replies"
	departmentConfigKey  = "ticket.departments"
)

type Service struct {
	repository repository.Repository
	customers  *customerService.Service
	orders     *orderService.Service
	audit      *audit.Service
	automation *automationService.Service
}

func New(
	repo repository.Repository,
	customerSvc *customerService.Service,
	orderSvc *orderService.Service,
	auditSvc *audit.Service,
	automationSvc *automationService.Service,
) *Service {
	return &Service{
		repository: repo,
		customers:  customerSvc,
		orders:     orderSvc,
		audit:      auditSvc,
		automation: automationSvc,
	}
}

func (service *Service) List(filter domain.ListFilter) ([]domain.Ticket, int) {
	items, total := service.repository.List(withDefaultPage(filter))
	return service.decorateTickets(items), total
}

func (service *Service) ListByCustomer(customerID int64, filter domain.ListFilter) ([]domain.Ticket, int, error) {
	if service.customers != nil {
		if _, exists := service.customers.GetByID(customerID); !exists {
			return nil, 0, repository.ErrCustomerNotFound
		}
	}
	items, total := service.repository.ListByCustomer(customerID, withDefaultPage(filter))
	return service.decorateTickets(items), total, nil
}

func (service *Service) GetDetail(id int64) (dto.DetailResponse, bool) {
	ticket, exists := service.repository.GetByID(id)
	if !exists {
		return dto.DetailResponse{}, false
	}
	return dto.DetailResponse{
		Ticket:    service.decorateTicket(ticket),
		Replies:   service.repository.ListReplies(id),
		AuditLogs: service.audit.ListByTarget("ticket", id),
	}, true
}

func (service *Service) GetDetailByCustomer(customerID, id int64) (dto.DetailResponse, bool) {
	ticket, exists := service.repository.GetByCustomer(customerID, id)
	if !exists {
		return dto.DetailResponse{}, false
	}
	return dto.DetailResponse{
		Ticket:    service.decorateTicket(ticket),
		Replies:   visibleReplies(service.repository.ListReplies(id)),
		AuditLogs: service.audit.ListByTarget("ticket", id),
	}, true
}

func (service *Service) ListPresetReplies(departmentName string) []domain.PresetReply {
	presets := service.loadPresetReplies()
	department := normalizeDepartmentName(departmentName)
	if department == "" {
		return presets
	}

	result := make([]domain.PresetReply, 0, len(presets))
	for _, item := range presets {
		itemDepartment := normalizeDepartmentName(item.DepartmentName)
		if itemDepartment == "" || strings.EqualFold(itemDepartment, department) {
			result = append(result, item)
		}
	}
	return result
}

func (service *Service) SavePresetReplies(items []domain.PresetReply) ([]domain.PresetReply, error) {
	normalized := normalizePresetReplies(items)
	if service.automation == nil {
		return normalized, nil
	}

	encoded, err := json.Marshal(normalized)
	if err != nil {
		return nil, err
	}
	if err := service.automation.SaveSystemConfig(presetReplyConfigKey, string(encoded), "ticket preset replies"); err != nil {
		return nil, err
	}
	return normalized, nil
}

func (service *Service) ListDepartments() []domain.Department {
	return service.loadDepartments()
}

func (service *Service) SaveDepartments(items []domain.Department) ([]domain.Department, error) {
	normalized := normalizeDepartments(items)
	if service.automation == nil {
		return normalized, nil
	}

	encoded, err := json.Marshal(normalized)
	if err != nil {
		return nil, err
	}
	if err := service.automation.SaveSystemConfig(departmentConfigKey, string(encoded), "ticket departments"); err != nil {
		return nil, err
	}
	return normalized, nil
}

func (service *Service) GetStatistics() dto.StatisticsResponse {
	departments := service.loadDepartments()
	departmentStatsMap := make(map[string]*domain.DepartmentStats)
	for _, item := range departments {
		departmentStatsMap[item.Name] = &domain.DepartmentStats{
			Key:  item.Key,
			Name: item.Name,
		}
	}

	adminStatsMap := make(map[string]*domain.AdminStats)
	summary := domain.SummaryStats{}
	allTickets := service.collectAllTickets()

	for _, raw := range allTickets {
		ticket := service.decorateTicket(raw)
		summary.Total++
		if ticket.AssignedAdminID == 0 {
			summary.Unassigned++
		}
		if ticket.Status == domain.StatusWaitingCustomer {
			summary.WaitingCustomer++
		}
		if ticket.Status == domain.StatusClosed {
			summary.Closed++
		}
		if ticket.SLAStatus == "BREACHED" {
			summary.Breached++
		}

		departmentName := firstNonEmpty(ticket.DepartmentName, "未分配部门")
		dept, exists := departmentStatsMap[departmentName]
		if !exists {
			dept = &domain.DepartmentStats{
				Key:  slugifyKey(departmentName, len(departmentStatsMap)+1),
				Name: departmentName,
			}
			departmentStatsMap[departmentName] = dept
		}
		dept.Total++
		switch ticket.Status {
		case domain.StatusOpen:
			dept.Open++
		case domain.StatusProcessing:
			dept.Processing++
		case domain.StatusWaitingCustomer:
			dept.WaitingCustomer++
		case domain.StatusClosed:
			dept.Closed++
		}
		if ticket.SLAStatus == "BREACHED" {
			dept.Breached++
		}

		adminName := firstNonEmpty(ticket.AssignedAdminName, "未分配")
		adminKey := fmt.Sprintf("%d:%s", ticket.AssignedAdminID, adminName)
		admin, exists := adminStatsMap[adminKey]
		if !exists {
			admin = &domain.AdminStats{
				AdminID:   ticket.AssignedAdminID,
				AdminName: adminName,
			}
			adminStatsMap[adminKey] = admin
		}
		admin.Total++
		switch ticket.Status {
		case domain.StatusOpen:
			admin.Open++
		case domain.StatusProcessing:
			admin.Processing++
		case domain.StatusWaitingCustomer:
			admin.WaitingCustomer++
		case domain.StatusClosed:
			admin.Closed++
		}
		if ticket.SLAStatus == "BREACHED" {
			admin.Breached++
		}
	}

	departmentStats := make([]domain.DepartmentStats, 0, len(departmentStatsMap))
	for _, item := range departmentStatsMap {
		departmentStats = append(departmentStats, *item)
	}
	sort.SliceStable(departmentStats, func(i, j int) bool {
		if departmentStats[i].Total == departmentStats[j].Total {
			return departmentStats[i].Name < departmentStats[j].Name
		}
		return departmentStats[i].Total > departmentStats[j].Total
	})

	adminStats := make([]domain.AdminStats, 0, len(adminStatsMap))
	for _, item := range adminStatsMap {
		adminStats = append(adminStats, *item)
	}
	sort.SliceStable(adminStats, func(i, j int) bool {
		if adminStats[i].Total == adminStats[j].Total {
			return adminStats[i].AdminName < adminStats[j].AdminName
		}
		return adminStats[i].Total > adminStats[j].Total
	})

	return dto.StatisticsResponse{
		Summary:         summary,
		DepartmentStats: departmentStats,
		AdminStats:      adminStats,
	}
}

func (service *Service) ClaimTicket(ticketID, adminID int64, adminName, requestID string) (domain.Ticket, bool, error) {
	current, exists := service.repository.GetByID(ticketID)
	if !exists {
		return domain.Ticket{}, false, nil
	}
	if current.Status == domain.StatusClosed {
		return domain.Ticket{}, false, repository.ErrTicketClosed
	}
	if current.AssignedAdminID > 0 && current.AssignedAdminID != adminID {
		return domain.Ticket{}, false, fmt.Errorf("ticket already assigned to %s", firstNonEmpty(current.AssignedAdminName, "another admin"))
	}
	if current.AssignedAdminID == adminID && adminID > 0 {
		return service.decorateTicket(current), true, nil
	}

	status := current.Status
	if status == domain.StatusOpen {
		status = domain.StatusProcessing
	}
	departmentName := current.DepartmentName
	if strings.TrimSpace(departmentName) == "" {
		departmentName = service.resolveDepartmentName("")
	}

	updated, ok, err := service.repository.Update(ticketID, domain.UpdateInput{
		Status:            &status,
		DepartmentName:    stringPointer(departmentName),
		AssignedAdminID:   &adminID,
		AssignedAdminName: stringPointer(firstNonEmpty(adminName, "System Admin")),
	})
	if err != nil || !ok {
		return domain.Ticket{}, ok, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       firstNonEmpty(adminName, "System Admin"),
		Action:      "ticket.claim",
		TargetType:  "ticket",
		TargetID:    updated.ID,
		Target:      updated.TicketNo,
		RequestID:   requestID,
		Description: "Admin claimed ticket",
		Payload: map[string]any{
			"status":            updated.Status,
			"departmentName":    updated.DepartmentName,
			"assignedAdminId":   updated.AssignedAdminID,
			"assignedAdminName": updated.AssignedAdminName,
		},
	})
	return service.decorateTicket(updated), true, nil
}

func (service *Service) RunAutoCloseSweep(operatorID int64, operatorName, requestID string) (dto.AutoCloseSweepResponse, error) {
	response := dto.AutoCloseSweepResponse{
		AutoCloseHours: service.autoCloseHours(),
		Items:          make([]dto.AutoCloseSweepItem, 0),
	}

	actorName := firstNonEmpty(operatorName, "Automation Engine")
	taskID := int64(0)
	if service.automation != nil {
		task := service.automation.Start(automationService.StartTaskRequest{
			TaskType:     "TICKET_AUTO_CLOSE",
			Title:        "Ticket auto-close sweep",
			Channel:      "LOCAL",
			Stage:        "TICKET",
			SourceType:   "cron",
			ActionName:   "ticket-auto-close",
			OperatorType: "ADMIN",
			OperatorName: actorName,
			Message:      "Manual ticket auto-close sweep triggered",
		})
		response.TaskID = task.ID
		response.TaskNo = task.TaskNo
		taskID = task.ID
	}

	if response.AutoCloseHours <= 0 {
		response.Message = "工单自动关单未启用"
		if taskID > 0 {
			service.automation.MarkBlocked(taskID, response.Message, response)
		}
		return response, nil
	}

	response.Enabled = true
	waitingTickets := service.collectTicketsByStatus(domain.StatusWaitingCustomer)
	response.CheckedCount = len(waitingTickets)

	for _, item := range waitingTickets {
		decorated := service.decorateTicket(item)
		if decorated.AutoCloseAt == "" || decorated.AutoCloseMins > 0 {
			response.SkippedCount++
			continue
		}

		replyContent := fmt.Sprintf("系统已自动关闭工单：客户超过 %d 小时未回复，如需继续处理请重新提交工单。", response.AutoCloseHours)
		_, updated, ok, err := service.repository.AddReply(item.ID, domain.ReplyInput{
			AuthorType: domain.AuthorTypeAdmin,
			AuthorID:   0,
			AuthorName: "Automation Engine",
			Content:    replyContent,
			Status:     domain.StatusClosed,
			IsInternal: false,
		})
		if err != nil {
			response.FailedCount++
			response.Items = append(response.Items, dto.AutoCloseSweepItem{
				TicketID:     item.ID,
				TicketNo:     item.TicketNo,
				CustomerName: item.CustomerName,
				Department:   item.DepartmentName,
				AutoCloseAt:  decorated.AutoCloseAt,
				Result:       "failed",
				Message:      err.Error(),
			})
			continue
		}
		if !ok {
			response.FailedCount++
			response.Items = append(response.Items, dto.AutoCloseSweepItem{
				TicketID:     item.ID,
				TicketNo:     item.TicketNo,
				CustomerName: item.CustomerName,
				Department:   item.DepartmentName,
				AutoCloseAt:  decorated.AutoCloseAt,
				Result:       "failed",
				Message:      "ticket not found",
			})
			continue
		}

		response.ClosedCount++
		response.Items = append(response.Items, dto.AutoCloseSweepItem{
			TicketID:     updated.ID,
			TicketNo:     updated.TicketNo,
			CustomerName: updated.CustomerName,
			Department:   updated.DepartmentName,
			AutoCloseAt:  decorated.AutoCloseAt,
			ClosedAt:     updated.ClosedAt,
			Result:       "closed",
			Message:      "已自动关闭",
		})

		service.audit.Record(audit.Entry{
			ActorType:   "SYSTEM",
			ActorID:     operatorID,
			Actor:       actorName,
			Action:      "ticket.auto_close",
			TargetType:  "ticket",
			TargetID:    updated.ID,
			Target:      updated.TicketNo,
			RequestID:   requestID,
			Description: "Ticket closed automatically after customer timeout",
			Payload: map[string]any{
				"departmentName": updated.DepartmentName,
				"autoCloseAt":    decorated.AutoCloseAt,
				"autoCloseHours": response.AutoCloseHours,
			},
		})
	}

	switch {
	case response.ClosedCount > 0 && response.FailedCount == 0:
		response.Message = fmt.Sprintf("本次巡检已自动关闭 %d 个工单", response.ClosedCount)
	case response.ClosedCount > 0 && response.FailedCount > 0:
		response.Message = fmt.Sprintf("本次巡检关闭 %d 个工单，失败 %d 个", response.ClosedCount, response.FailedCount)
	case response.FailedCount > 0:
		response.Message = fmt.Sprintf("本次巡检没有成功关闭工单，失败 %d 个", response.FailedCount)
	default:
		response.Message = "本次巡检未发现需要自动关闭的工单"
	}

	if taskID > 0 {
		switch {
		case response.FailedCount > 0 && response.ClosedCount == 0:
			service.automation.MarkBlocked(taskID, response.Message, response)
		default:
			service.automation.MarkSuccess(taskID, response.Message, response)
		}
	}
	return response, nil
}

func (service *Service) CreatePortalTicket(
	customerID int64,
	customerName string,
	request dto.CreateTicketRequest,
	requestID string,
) (domain.Ticket, error) {
	customer, exists := service.customers.GetByID(customerID)
	if !exists {
		return domain.Ticket{}, repository.ErrCustomerNotFound
	}

	serviceID := int64Value(request.ServiceID)
	if serviceID > 0 && service.orders != nil {
		detail, ok := service.orders.GetServiceDetail(serviceID)
		if !ok || detail.Service.CustomerID != customerID {
			return domain.Ticket{}, repository.ErrServiceNotFound
		}
	}

	priority, err := normalizePriority(request.Priority)
	if err != nil {
		return domain.Ticket{}, err
	}

	departmentName := service.resolveDepartmentName(request.DepartmentName)
	assignedAdminID := int64Value(request.AssignedAdminID)
	if err := service.validateDepartmentAssignment(departmentName, assignedAdminID); err != nil {
		return domain.Ticket{}, err
	}

	created, err := service.repository.Create(domain.CreateInput{
		CustomerID:      customerID,
		CustomerNo:      customer.CustomerNo,
		CustomerName:    firstNonEmpty(customerName, customer.Name),
		ServiceID:       serviceID,
		Title:           strings.TrimSpace(request.Title),
		Content:         strings.TrimSpace(request.Content),
		Status:          domain.StatusOpen,
		Priority:        priority,
		Source:          "PORTAL",
		DepartmentName:  departmentName,
		AssignedAdminID: assignedAdminID,
	})
	if err != nil {
		return domain.Ticket{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "CUSTOMER",
		ActorID:     customerID,
		Actor:       firstNonEmpty(customerName, customer.Name, created.CustomerName),
		Action:      "ticket.create",
		TargetType:  "ticket",
		TargetID:    created.ID,
		Target:      created.TicketNo,
		RequestID:   requestID,
		Description: "Customer created a support ticket",
		Payload: map[string]any{
			"title":          created.Title,
			"serviceId":      created.ServiceID,
			"priority":       created.Priority,
			"departmentName": created.DepartmentName,
		},
	})
	return created, nil
}

func (service *Service) CreateAdminTicket(
	request dto.CreateAdminTicketRequest,
	adminID int64,
	adminName string,
	requestID string,
) (domain.Ticket, error) {
	customer, exists := service.customers.GetByID(request.CustomerID)
	if !exists {
		return domain.Ticket{}, repository.ErrCustomerNotFound
	}

	serviceID := int64Value(request.ServiceID)
	if serviceID > 0 && service.orders != nil {
		detail, ok := service.orders.GetServiceDetail(serviceID)
		if !ok || detail.Service.CustomerID != request.CustomerID {
			return domain.Ticket{}, repository.ErrServiceNotFound
		}
	}

	priority, err := normalizePriority(request.Priority)
	if err != nil {
		return domain.Ticket{}, err
	}

	departmentName := service.resolveDepartmentName(request.DepartmentName)
	assignedAdminID := int64Value(request.AssignedAdminID)
	if err := service.validateDepartmentAssignment(departmentName, assignedAdminID); err != nil {
		return domain.Ticket{}, err
	}

	assignedAdminName := strings.TrimSpace(stringValue(request.AssignedAdminName))
	status := domain.StatusOpen
	if assignedAdminID > 0 {
		status = domain.StatusProcessing
	}

	created, err := service.repository.Create(domain.CreateInput{
		CustomerID:      request.CustomerID,
		CustomerNo:      customer.CustomerNo,
		CustomerName:    customer.Name,
		ServiceID:       serviceID,
		Title:           strings.TrimSpace(request.Title),
		Content:         strings.TrimSpace(request.Content),
		Status:          status,
		Priority:        priority,
		Source:          "ADMIN",
		DepartmentName:  departmentName,
		AssignedAdminID: assignedAdminID,
		AssignedAdmin:   assignedAdminName,
	})
	if err != nil {
		return domain.Ticket{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       firstNonEmpty(adminName, "System Admin"),
		Action:      "ticket.create_admin",
		TargetType:  "ticket",
		TargetID:    created.ID,
		Target:      created.TicketNo,
		RequestID:   requestID,
		Description: "Admin created a support ticket",
		Payload: map[string]any{
			"customerId":        created.CustomerID,
			"serviceId":         created.ServiceID,
			"priority":          created.Priority,
			"departmentName":    created.DepartmentName,
			"assignedAdminId":   created.AssignedAdminID,
			"assignedAdminName": created.AssignedAdminName,
		},
	})
	return created, nil
}

func (service *Service) UpdateAdminTicket(
	ticketID int64,
	request dto.UpdateTicketRequest,
	adminID int64,
	adminName,
	requestID string,
) (domain.Ticket, bool, error) {
	current, exists := service.repository.GetByID(ticketID)
	if !exists {
		return domain.Ticket{}, false, nil
	}

	status, err := normalizeOptionalStatus(request.Status)
	if err != nil {
		return domain.Ticket{}, false, err
	}
	priority, err := normalizeOptionalPriority(request.Priority)
	if err != nil {
		return domain.Ticket{}, false, err
	}

	var departmentName *string
	effectiveDepartment := current.DepartmentName
	if request.DepartmentName != nil {
		resolved := service.resolveDepartmentName(*request.DepartmentName)
		departmentName = &resolved
		effectiveDepartment = resolved
	}

	effectiveAdminID := current.AssignedAdminID
	if request.AssignedAdminID != nil {
		effectiveAdminID = *request.AssignedAdminID
	}
	if err := service.validateDepartmentAssignment(effectiveDepartment, effectiveAdminID); err != nil {
		return domain.Ticket{}, false, err
	}

	updated, ok, err := service.repository.Update(ticketID, domain.UpdateInput{
		Title:             trimOptionalString(request.Title),
		Status:            status,
		Priority:          priority,
		DepartmentName:    departmentName,
		AssignedAdminID:   request.AssignedAdminID,
		AssignedAdminName: resolveAssignedAdminName(request.AssignedAdminID, request.AssignedAdminName, adminName),
	})
	if err != nil {
		return domain.Ticket{}, false, err
	}
	if !ok {
		return domain.Ticket{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "ticket.update",
		TargetType:  "ticket",
		TargetID:    updated.ID,
		Target:      updated.TicketNo,
		RequestID:   requestID,
		Description: "Admin updated ticket fields",
		Payload: map[string]any{
			"status":            updated.Status,
			"priority":          updated.Priority,
			"departmentName":    updated.DepartmentName,
			"assignedAdminId":   updated.AssignedAdminID,
			"assignedAdminName": updated.AssignedAdminName,
		},
	})
	return updated, true, nil
}

func (service *Service) ReplyAsAdmin(
	ticketID int64,
	request dto.ReplyRequest,
	adminID int64,
	adminName,
	requestID string,
) (dto.DetailResponse, bool, error) {
	status, err := normalizeReplyStatus(request.Status, domain.StatusWaitingCustomer)
	if err != nil {
		return dto.DetailResponse{}, false, err
	}

	_, ticket, ok, err := service.repository.AddReply(ticketID, domain.ReplyInput{
		AuthorType: domain.AuthorTypeAdmin,
		AuthorID:   adminID,
		AuthorName: adminName,
		Content:    strings.TrimSpace(request.Content),
		Status:     status,
		IsInternal: request.IsInternal,
	})
	if err != nil {
		return dto.DetailResponse{}, false, err
	}
	if !ok {
		return dto.DetailResponse{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "ticket.reply",
		TargetType:  "ticket",
		TargetID:    ticket.ID,
		Target:      ticket.TicketNo,
		RequestID:   requestID,
		Description: "Admin replied to ticket",
		Payload: map[string]any{
			"status":     ticket.Status,
			"isInternal": request.IsInternal,
		},
	})

	result, _ := service.GetDetail(ticketID)
	return result, true, nil
}

func (service *Service) ReplyAsCustomer(
	customerID int64,
	customerName string,
	ticketID int64,
	request dto.ReplyRequest,
	requestID string,
) (dto.DetailResponse, bool, error) {
	ticket, exists := service.repository.GetByCustomer(customerID, ticketID)
	if !exists {
		return dto.DetailResponse{}, false, nil
	}
	if ticket.Status == domain.StatusClosed {
		return dto.DetailResponse{}, false, repository.ErrTicketClosed
	}

	status, err := normalizeReplyStatus(request.Status, domain.StatusOpen)
	if err != nil {
		return dto.DetailResponse{}, false, err
	}

	_, updated, ok, err := service.repository.AddReply(ticketID, domain.ReplyInput{
		AuthorType: domain.AuthorTypeCustomer,
		AuthorID:   customerID,
		AuthorName: firstNonEmpty(customerName, ticket.CustomerName),
		Content:    strings.TrimSpace(request.Content),
		Status:     status,
		IsInternal: false,
	})
	if err != nil {
		return dto.DetailResponse{}, false, err
	}
	if !ok {
		return dto.DetailResponse{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "CUSTOMER",
		ActorID:     customerID,
		Actor:       firstNonEmpty(customerName, ticket.CustomerName),
		Action:      "ticket.reply",
		TargetType:  "ticket",
		TargetID:    updated.ID,
		Target:      updated.TicketNo,
		RequestID:   requestID,
		Description: "Customer replied to ticket",
		Payload: map[string]any{
			"status": updated.Status,
		},
	})

	result, _ := service.GetDetailByCustomer(customerID, ticketID)
	return result, true, nil
}

func (service *Service) CloseByCustomer(
	customerID int64,
	customerName string,
	ticketID int64,
	requestID string,
) (domain.Ticket, bool, error) {
	ticket, exists := service.repository.GetByCustomer(customerID, ticketID)
	if !exists {
		return domain.Ticket{}, false, nil
	}

	updated, ok, err := service.repository.Close(ticketID)
	if err != nil {
		return domain.Ticket{}, false, err
	}
	if !ok {
		return domain.Ticket{}, false, nil
	}

	service.audit.Record(audit.Entry{
		ActorType:   "CUSTOMER",
		ActorID:     customerID,
		Actor:       firstNonEmpty(customerName, ticket.CustomerName),
		Action:      "ticket.close",
		TargetType:  "ticket",
		TargetID:    updated.ID,
		Target:      updated.TicketNo,
		RequestID:   requestID,
		Description: "Customer closed ticket",
	})
	return updated, true, nil
}

func withDefaultPage(filter domain.ListFilter) domain.ListFilter {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return filter
}

func normalizePriority(value string) (domain.Priority, error) {
	if strings.TrimSpace(value) == "" {
		return domain.PriorityNormal, nil
	}
	switch domain.Priority(strings.ToUpper(strings.TrimSpace(value))) {
	case domain.PriorityLow, domain.PriorityNormal, domain.PriorityHigh, domain.PriorityUrgent:
		return domain.Priority(strings.ToUpper(strings.TrimSpace(value))), nil
	default:
		return "", fmt.Errorf("invalid ticket priority")
	}
}

func normalizeOptionalPriority(value *string) (*domain.Priority, error) {
	if value == nil {
		return nil, nil
	}
	normalized, err := normalizePriority(*value)
	if err != nil {
		return nil, err
	}
	return &normalized, nil
}

func normalizeOptionalStatus(value *string) (*domain.Status, error) {
	if value == nil {
		return nil, nil
	}
	normalized, err := normalizeStatus(*value)
	if err != nil {
		return nil, err
	}
	return &normalized, nil
}

func normalizeReplyStatus(value *string, fallback domain.Status) (domain.Status, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback, nil
	}
	return normalizeStatus(*value)
}

func normalizeStatus(value string) (domain.Status, error) {
	switch domain.Status(strings.ToUpper(strings.TrimSpace(value))) {
	case domain.StatusOpen, domain.StatusProcessing, domain.StatusWaitingCustomer, domain.StatusClosed:
		return domain.Status(strings.ToUpper(strings.TrimSpace(value))), nil
	default:
		return "", fmt.Errorf("invalid ticket status")
	}
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	return &trimmed
}

func resolveAssignedAdminName(adminID *int64, explicit *string, fallback string) *string {
	if adminID == nil && explicit == nil {
		return nil
	}
	if explicit != nil {
		trimmed := strings.TrimSpace(*explicit)
		return &trimmed
	}
	if adminID != nil && *adminID > 0 {
		name := strings.TrimSpace(fallback)
		return &name
	}
	empty := ""
	return &empty
}

func int64Value(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func stringPointer(value string) *string {
	trimmed := strings.TrimSpace(value)
	return &trimmed
}

func firstNonEmpty(values ...string) string {
	for _, item := range values {
		item = strings.TrimSpace(item)
		if item != "" {
			return item
		}
	}
	return ""
}

func (service *Service) collectAllTickets() []domain.Ticket {
	result := make([]domain.Ticket, 0)
	for page := 1; ; page++ {
		items, total := service.repository.List(domain.ListFilter{
			Page:  page,
			Limit: 200,
		})
		if len(items) == 0 {
			break
		}
		result = append(result, items...)
		if len(result) >= total || len(items) < 200 {
			break
		}
	}
	return result
}

func (service *Service) collectTicketsByStatus(status domain.Status) []domain.Ticket {
	result := make([]domain.Ticket, 0)
	for page := 1; ; page++ {
		items, total := service.repository.List(domain.ListFilter{
			Page:   page,
			Limit:  200,
			Status: string(status),
		})
		if len(items) == 0 {
			break
		}
		result = append(result, items...)
		if len(result) >= total || len(items) < 200 {
			break
		}
	}
	return result
}

func (service *Service) decorateTickets(items []domain.Ticket) []domain.Ticket {
	result := make([]domain.Ticket, 0, len(items))
	for _, item := range items {
		result = append(result, service.decorateTicket(item))
	}
	return result
}

func (service *Service) decorateTicket(item domain.Ticket) domain.Ticket {
	now := time.Now()
	settings := service.autoCloseHours()
	baseTime := parseTicketTime(firstNonEmpty(item.LastReplyAt, item.CreatedAt))

	switch item.Status {
	case domain.StatusClosed:
		item.SLAStatus = "MET"
		item.SLADeadlineAt = ""
		item.SLARemainingMins = 0
		item.SLAPaused = false
		item.AutoCloseAt = ""
		item.AutoCloseMins = 0
		return item
	case domain.StatusWaitingCustomer:
		item.SLAStatus = "PAUSED"
		item.SLAPaused = true
		item.SLADeadlineAt = ""
		item.SLARemainingMins = 0
		if !baseTime.IsZero() && settings > 0 {
			autoCloseAt := baseTime.Add(time.Duration(settings) * time.Hour)
			item.AutoCloseAt = autoCloseAt.Format("2006-01-02 15:04:05")
			item.AutoCloseMins = int(autoCloseAt.Sub(now).Minutes())
		}
		return item
	default:
		item.SLAPaused = false
		item.AutoCloseAt = ""
		item.AutoCloseMins = 0
	}

	target := prioritySLA(item.Priority)
	if target <= 0 || baseTime.IsZero() {
		item.SLAStatus = "UNKNOWN"
		item.SLADeadlineAt = ""
		item.SLARemainingMins = 0
		return item
	}

	deadline := baseTime.Add(target)
	remaining := int(deadline.Sub(now).Minutes())
	item.SLADeadlineAt = deadline.Format("2006-01-02 15:04:05")
	item.SLARemainingMins = remaining

	switch {
	case remaining <= 0:
		item.SLAStatus = "BREACHED"
	case remaining <= int(target.Minutes()/4):
		item.SLAStatus = "AT_RISK"
	default:
		item.SLAStatus = "ON_TRACK"
	}
	return item
}

func (service *Service) autoCloseHours() int {
	if service.automation == nil {
		return 72
	}
	settings := service.automation.GetSettings()
	if !settings.TicketAutoCloseOn || settings.TicketAutoCloseHours <= 0 {
		return 0
	}
	return settings.TicketAutoCloseHours
}

func (service *Service) loadPresetReplies() []domain.PresetReply {
	if service.automation == nil {
		return defaultPresetReplies()
	}

	raw := service.automation.GetSystemConfig(presetReplyConfigKey, "")
	if strings.TrimSpace(raw) == "" {
		return defaultPresetReplies()
	}

	var items []domain.PresetReply
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return defaultPresetReplies()
	}
	return normalizePresetReplies(items)
}

func (service *Service) loadDepartments() []domain.Department {
	if service.automation == nil {
		return defaultDepartments()
	}

	raw := service.automation.GetSystemConfig(departmentConfigKey, "")
	if strings.TrimSpace(raw) == "" {
		return defaultDepartments()
	}

	var items []domain.Department
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return defaultDepartments()
	}
	return normalizeDepartments(items)
}

func (service *Service) validateDepartmentAssignment(departmentName string, adminID int64) error {
	if adminID <= 0 {
		return nil
	}
	department, exists := service.findDepartment(departmentName)
	if !exists || len(department.AdminIDs) == 0 {
		return nil
	}
	for _, item := range department.AdminIDs {
		if item == adminID {
			return nil
		}
	}
	return fmt.Errorf("admin %d is not allowed in department %s", adminID, department.Name)
}

func (service *Service) findDepartment(value string) (domain.Department, bool) {
	normalized := normalizeDepartmentName(value)
	for _, item := range service.loadDepartments() {
		if strings.EqualFold(item.Name, normalized) {
			return item, true
		}
	}
	return domain.Department{}, false
}

func (service *Service) resolveDepartmentName(value string) string {
	departments := service.loadDepartments()
	normalized := normalizeDepartmentName(value)

	if normalized != "" {
		for _, item := range departments {
			if item.Enabled && strings.EqualFold(item.Name, normalized) {
				return item.Name
			}
		}
		return normalized
	}

	for _, item := range departments {
		if item.Enabled && item.IsDefault {
			return item.Name
		}
	}
	for _, item := range departments {
		if item.Enabled {
			return item.Name
		}
	}
	return "Technical Support"
}

func defaultPresetReplies() []domain.PresetReply {
	return normalizePresetReplies([]domain.PresetReply{
		{
			Key:     "need-more-info",
			Title:   "补充信息",
			Content: "您好，为了继续处理，请补充实例编号、操作时间和报错截图，我们收到后会优先跟进。",
			Status:  string(domain.StatusWaitingCustomer),
		},
		{
			Key:            "processing",
			Title:          "处理中",
			Content:        "您好，工单已接单，正在安排技术同事排查，处理进度会持续同步到本工单。",
			DepartmentName: "Technical Support",
			Status:         string(domain.StatusProcessing),
		},
		{
			Key:            "resolved-confirm",
			Title:          "已处理待确认",
			Content:        "您好，当前问题已经处理完成，请您核实结果。如无其他问题可直接关闭工单。",
			DepartmentName: "Technical Support",
			Status:         string(domain.StatusWaitingCustomer),
		},
		{
			Key:            "finance-adjusting",
			Title:          "财务核对中",
			Content:        "您好，财务同事已收到您的申请，正在核对账单和金额，预计稍后给您反馈。",
			DepartmentName: "Finance",
			Status:         string(domain.StatusProcessing),
		},
		{
			Key:            "finance-confirm",
			Title:          "已调整待确认",
			Content:        "您好，相关账单信息已经调整，请您刷新后查看，如确认无误可继续支付。",
			DepartmentName: "Finance",
			Status:         string(domain.StatusWaitingCustomer),
		},
	})
}

func defaultDepartments() []domain.Department {
	return normalizeDepartments([]domain.Department{
		{
			Key:         "technical-support",
			Name:        "Technical Support",
			Description: "技术问题、网络故障、实例异常、远程连接等",
			Enabled:     true,
			IsDefault:   true,
			Sort:        10,
		},
		{
			Key:         "finance",
			Name:        "Finance",
			Description: "账单、退款、改价、余额与开票相关问题",
			Enabled:     true,
			IsDefault:   false,
			Sort:        20,
		},
		{
			Key:         "sales",
			Name:        "Sales",
			Description: "售前咨询、方案推荐、产品选型与商务沟通",
			Enabled:     true,
			IsDefault:   false,
			Sort:        30,
		},
	})
}

func normalizePresetReplies(items []domain.PresetReply) []domain.PresetReply {
	result := make([]domain.PresetReply, 0, len(items))
	seenKeys := make(map[string]int)

	for index, item := range items {
		title := strings.TrimSpace(item.Title)
		content := strings.TrimSpace(item.Content)
		if title == "" || content == "" {
			continue
		}

		key := strings.TrimSpace(item.Key)
		if key == "" {
			key = fmt.Sprintf("preset-%02d", index+1)
		}
		seenKeys[key]++
		if seenKeys[key] > 1 {
			key = fmt.Sprintf("%s-%d", key, seenKeys[key])
		}

		status := strings.ToUpper(strings.TrimSpace(item.Status))
		switch domain.Status(status) {
		case domain.StatusOpen, domain.StatusProcessing, domain.StatusWaitingCustomer, domain.StatusClosed:
		default:
			status = string(domain.StatusWaitingCustomer)
		}

		result = append(result, domain.PresetReply{
			Key:            key,
			Title:          title,
			Content:        content,
			DepartmentName: normalizeDepartmentName(item.DepartmentName),
			Status:         status,
		})
	}
	return result
}

func normalizeDepartments(items []domain.Department) []domain.Department {
	result := make([]domain.Department, 0, len(items))
	seenKeys := make(map[string]int)

	for index, item := range items {
		name := normalizeDepartmentName(item.Name)
		if name == "" {
			continue
		}

		key := strings.TrimSpace(item.Key)
		if key == "" {
			key = slugifyKey(name, index+1)
		}
		seenKeys[key]++
		if seenKeys[key] > 1 {
			key = fmt.Sprintf("%s-%d", key, seenKeys[key])
		}

		sortIndex := item.Sort
		if sortIndex <= 0 {
			sortIndex = (index + 1) * 10
		}

		result = append(result, domain.Department{
			Key:         key,
			Name:        name,
			Description: strings.TrimSpace(item.Description),
			Enabled:     item.Enabled,
			IsDefault:   item.IsDefault,
			Sort:        sortIndex,
			AdminIDs:    normalizeAdminIDs(item.AdminIDs),
		})
	}

	if len(result) == 0 {
		return defaultDepartments()
	}

	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Sort == result[j].Sort {
			return result[i].Name < result[j].Name
		}
		return result[i].Sort < result[j].Sort
	})

	hasEnabled := false
	for index := range result {
		if result[index].Enabled {
			hasEnabled = true
			break
		}
	}
	if !hasEnabled {
		result[0].Enabled = true
	}

	defaultAssigned := false
	for index := range result {
		if result[index].Enabled && result[index].IsDefault && !defaultAssigned {
			defaultAssigned = true
			continue
		}
		result[index].IsDefault = false
	}
	if !defaultAssigned {
		for index := range result {
			if result[index].Enabled {
				result[index].IsDefault = true
				break
			}
		}
	}

	return result
}

func normalizeAdminIDs(items []int64) []int64 {
	if len(items) == 0 {
		return nil
	}

	result := make([]int64, 0, len(items))
	seen := make(map[int64]struct{}, len(items))
	for _, item := range items {
		if item <= 0 {
			continue
		}
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func normalizeDepartmentName(value string) string {
	trimmed := strings.TrimSpace(value)
	switch strings.ToLower(trimmed) {
	case "", "all", "common", "global", "通用":
		return ""
	case "technical", "technical support", "support", "技术支持", "售后支持":
		return "Technical Support"
	case "finance", "财务":
		return "Finance"
	default:
		return trimmed
	}
}

func slugifyKey(value string, fallbackIndex int) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return fmt.Sprintf("department-%02d", fallbackIndex)
	}

	var builder strings.Builder
	lastDash := false
	for _, char := range value {
		switch {
		case (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9'):
			builder.WriteRune(char)
			lastDash = false
		case char == ' ' || char == '-' || char == '_' || char == '/':
			if !lastDash {
				builder.WriteRune('-')
				lastDash = true
			}
		}
	}

	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return fmt.Sprintf("department-%02d", fallbackIndex)
	}
	return result
}

func prioritySLA(priority domain.Priority) time.Duration {
	switch priority {
	case domain.PriorityUrgent:
		return 30 * time.Minute
	case domain.PriorityHigh:
		return 2 * time.Hour
	case domain.PriorityLow:
		return 24 * time.Hour
	default:
		return 8 * time.Hour
	}
}

func parseTicketTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
	if err == nil {
		return parsed
	}
	parsed, err = time.ParseInLocation("2006-01-02", value, time.Local)
	if err == nil {
		return parsed
	}
	return time.Time{}
}

func visibleReplies(items []domain.Reply) []domain.Reply {
	result := make([]domain.Reply, 0, len(items))
	for _, item := range items {
		if item.IsInternal {
			continue
		}
		result = append(result, item)
	}
	return result
}
