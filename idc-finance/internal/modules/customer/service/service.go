package service

import (
	"idc-finance/internal/modules/customer/domain"
	"idc-finance/internal/modules/customer/dto"
	"idc-finance/internal/modules/customer/repository"
	"idc-finance/internal/platform/audit"
)

type Service struct {
	repository repository.Repository
	audit      *audit.Service
}

func New(repo repository.Repository, auditService *audit.Service) *Service {
	return &Service{repository: repo, audit: auditService}
}

func (service *Service) List() []domain.Customer {
	return service.repository.List()
}

func (service *Service) ListGroups() []domain.CustomerGroup {
	return service.repository.ListGroups()
}

func (service *Service) ListLevels() []domain.CustomerLevel {
	return service.repository.ListLevels()
}

func (service *Service) CreateGroup(request dto.SaveCustomerGroupRequest, adminID int64, adminName, requestID string) (domain.CustomerGroup, error) {
	created, err := service.repository.CreateGroup(domain.CustomerGroup{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		return domain.CustomerGroup{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.group.create",
		TargetType:  "customer-group",
		TargetID:    created.ID,
		Target:      created.Name,
		RequestID:   requestID,
		Description: "创建客户分组",
	})
	return created, nil
}

func (service *Service) UpdateGroup(id int64, request dto.SaveCustomerGroupRequest, adminID int64, adminName, requestID string) (domain.CustomerGroup, error) {
	updated, err := service.repository.UpdateGroup(id, domain.CustomerGroup{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		return domain.CustomerGroup{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.group.update",
		TargetType:  "customer-group",
		TargetID:    updated.ID,
		Target:      updated.Name,
		RequestID:   requestID,
		Description: "更新客户分组",
	})
	return updated, nil
}

func (service *Service) DeleteGroup(id int64, adminID int64, adminName, requestID string) error {
	err := service.repository.DeleteGroup(id)
	if err != nil {
		return err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.group.delete",
		TargetType:  "customer-group",
		TargetID:    id,
		RequestID:   requestID,
		Description: "删除客户分组",
	})
	return nil
}

func (service *Service) CreateLevel(request dto.SaveCustomerLevelRequest, adminID int64, adminName, requestID string) (domain.CustomerLevel, error) {
	created, err := service.repository.CreateLevel(domain.CustomerLevel{
		Name:        request.Name,
		Priority:    request.Priority,
		Description: request.Description,
	})
	if err != nil {
		return domain.CustomerLevel{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.level.create",
		TargetType:  "customer-level",
		TargetID:    created.ID,
		Target:      created.Name,
		RequestID:   requestID,
		Description: "创建客户等级",
	})
	return created, nil
}

func (service *Service) UpdateLevel(id int64, request dto.SaveCustomerLevelRequest, adminID int64, adminName, requestID string) (domain.CustomerLevel, error) {
	updated, err := service.repository.UpdateLevel(id, domain.CustomerLevel{
		Name:        request.Name,
		Priority:    request.Priority,
		Description: request.Description,
	})
	if err != nil {
		return domain.CustomerLevel{}, err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.level.update",
		TargetType:  "customer-level",
		TargetID:    updated.ID,
		Target:      updated.Name,
		RequestID:   requestID,
		Description: "更新客户等级",
	})
	return updated, nil
}

func (service *Service) DeleteLevel(id int64, adminID int64, adminName, requestID string) error {
	err := service.repository.DeleteLevel(id)
	if err != nil {
		return err
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.level.delete",
		TargetType:  "customer-level",
		TargetID:    id,
		RequestID:   requestID,
		Description: "删除客户等级",
	})
	return nil
}

func (service *Service) GetByID(id int64) (domain.Customer, bool) {
	return service.repository.GetByID(id)
}

func (service *Service) Create(request dto.CreateCustomerRequest) domain.Customer {
	return service.repository.Create(domain.Customer{
		Name:       request.Name,
		Email:      request.Email,
		Mobile:     request.Mobile,
		Type:       request.Type,
		GroupName:  request.GroupName,
		LevelName:  request.LevelName,
		SalesOwner: request.SalesOwner,
		Remarks:    request.Remarks,
	})
}

func (service *Service) CreateWithAudit(request dto.CreateCustomerRequest, adminID int64, adminName, requestID string) domain.Customer {
	customer := service.Create(request)
	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.create",
		TargetType:  "customer",
		TargetID:    customer.ID,
		Target:      customer.CustomerNo,
		RequestID:   requestID,
		Description: "创建客户并初始化实名档案",
		Payload: map[string]any{
			"name":  customer.Name,
			"email": customer.Email,
		},
	})
	return customer
}

func (service *Service) Update(id int64, request dto.UpdateCustomerRequest, adminID int64, adminName, requestID string) (domain.Customer, bool) {
	current, exists := service.repository.GetByID(id)
	if !exists {
		return domain.Customer{}, false
	}

	status := current.Status
	if request.Status != "" {
		status = domain.CustomerStatus(request.Status)
	}

	updated, ok := service.repository.Update(id, domain.Customer{
		Name:       request.Name,
		Email:      request.Email,
		Mobile:     request.Mobile,
		Type:       current.Type,
		Status:     status,
		GroupName:  request.GroupName,
		LevelName:  request.LevelName,
		SalesOwner: request.SalesOwner,
		Remarks:    request.Remarks,
	})
	if !ok {
		return domain.Customer{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.update",
		TargetType:  "customer",
		TargetID:    updated.ID,
		Target:      updated.CustomerNo,
		RequestID:   requestID,
		Description: "更新客户基础资料",
	})
	return updated, true
}

func (service *Service) UpdateSelf(customerID int64, request dto.UpdateCustomerRequest, requestID string) (domain.Customer, bool) {
	current, exists := service.repository.GetByID(customerID)
	if !exists {
		return domain.Customer{}, false
	}

	updated, ok := service.repository.Update(customerID, domain.Customer{
		Name:       request.Name,
		Email:      request.Email,
		Mobile:     request.Mobile,
		Type:       current.Type,
		Status:     current.Status,
		GroupName:  current.GroupName,
		LevelName:  current.LevelName,
		SalesOwner: current.SalesOwner,
		Remarks:    request.Remarks,
	})
	if !ok {
		return domain.Customer{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "CUSTOMER",
		ActorID:     updated.ID,
		Actor:       updated.Name,
		Action:      "customer.self.update",
		TargetType:  "customer",
		TargetID:    updated.ID,
		Target:      updated.CustomerNo,
		RequestID:   requestID,
		Description: "客户更新门户基础资料",
	})
	return updated, true
}

func (service *Service) ListIdentityOverview() []dto.IdentityOverviewItem {
	customers := service.List()
	identities := make([]dto.IdentityOverviewItem, 0, len(customers))
	for _, customer := range customers {
		if customer.Identity != nil {
			identities = append(identities, dto.IdentityOverviewItem{
				ID:           customer.Identity.ID,
				CustomerID:   customer.ID,
				CustomerNo:   customer.CustomerNo,
				CustomerName: customer.Name,
				CustomerType: customer.Type,
				IdentityType: customer.Identity.IdentityType,
				VerifyStatus: customer.Identity.VerifyStatus,
				SubjectName:  customer.Identity.SubjectName,
				CertNo:       customer.Identity.CertNo,
				CountryCode:  customer.Identity.CountryCode,
				ReviewRemark: customer.Identity.ReviewRemark,
				ReviewedAt:   customer.Identity.ReviewedAt,
				SubmittedAt:  customer.Identity.SubmittedAt,
			})
		}
	}
	return identities
}

func (service *Service) GetWorkbench(id int64) (dto.CustomerWorkbenchResponse, bool) {
	customer, exists := service.repository.GetByID(id)
	if !exists {
		return dto.CustomerWorkbenchResponse{}, false
	}

	return dto.CustomerWorkbenchResponse{
		Customer:  customer,
		Services:  service.repository.ListServiceItems(id),
		Invoices:  service.repository.ListInvoiceItems(id),
		Tickets:   service.repository.ListTicketItems(id),
		AuditLogs: service.audit.ListByTarget("customer", id),
	}, true
}

func (service *Service) ListContacts(customerID int64) ([]domain.Contact, bool) {
	customer, exists := service.repository.GetByID(customerID)
	if !exists {
		return nil, false
	}
	return customer.Contacts, true
}

func (service *Service) ListIdentityDetails(customerID int64) ([]domain.Identity, bool) {
	customer, exists := service.repository.GetByID(customerID)
	if !exists {
		return nil, false
	}
	if customer.Identity == nil {
		return []domain.Identity{}, true
	}
	return []domain.Identity{*customer.Identity}, true
}

func (service *Service) ListServices(customerID int64) ([]domain.RelatedItem, bool) {
	if _, exists := service.repository.GetByID(customerID); !exists {
		return nil, false
	}
	return service.repository.ListServiceItems(customerID), true
}

func (service *Service) ListInvoices(customerID int64) ([]domain.RelatedItem, bool) {
	if _, exists := service.repository.GetByID(customerID); !exists {
		return nil, false
	}
	return service.repository.ListInvoiceItems(customerID), true
}

func (service *Service) ListTickets(customerID int64) ([]domain.RelatedItem, bool) {
	if _, exists := service.repository.GetByID(customerID); !exists {
		return nil, false
	}
	return service.repository.ListTicketItems(customerID), true
}

func (service *Service) ListCustomerAuditLogs(customerID int64) []audit.Entry {
	return service.audit.ListByTarget("customer", customerID)
}

func (service *Service) AddContactByCustomer(customerID int64, request dto.CreateContactRequest, requestID string) (domain.Contact, bool) {
	customer, contact, ok := service.repository.AddContact(customerID, domain.Contact{
		Name:      request.Name,
		Email:     request.Email,
		Mobile:    request.Mobile,
		RoleName:  request.RoleName,
		IsPrimary: request.IsPrimary,
	})
	if !ok {
		return domain.Contact{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "CUSTOMER",
		ActorID:     customer.ID,
		Actor:       customer.Name,
		Action:      "customer.contact.self_create",
		TargetType:  "customer",
		TargetID:    customer.ID,
		Target:      customer.CustomerNo,
		RequestID:   requestID,
		Description: "客户在门户新增联系人",
		Payload: map[string]any{
			"contactName": contact.Name,
			"roleName":    contact.RoleName,
		},
	})
	return contact, true
}

func (service *Service) AddContact(customerID int64, request dto.CreateContactRequest, adminID int64, adminName, requestID string) (domain.Contact, bool) {
	customer, contact, ok := service.repository.AddContact(customerID, domain.Contact{
		Name:      request.Name,
		Email:     request.Email,
		Mobile:    request.Mobile,
		RoleName:  request.RoleName,
		IsPrimary: request.IsPrimary,
	})
	if !ok {
		return domain.Contact{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.contact.create",
		TargetType:  "customer",
		TargetID:    customer.ID,
		Target:      customer.CustomerNo,
		RequestID:   requestID,
		Description: "新增客户联系人",
		Payload: map[string]any{
			"contactName": contact.Name,
			"roleName":    contact.RoleName,
		},
	})
	return contact, true
}

func (service *Service) UpdateContactByCustomer(customerID, contactID int64, request dto.UpdateContactRequest, requestID string) (domain.Contact, bool) {
	customer, contact, ok := service.repository.UpdateContact(customerID, contactID, domain.Contact{
		Name:      request.Name,
		Email:     request.Email,
		Mobile:    request.Mobile,
		RoleName:  request.RoleName,
		IsPrimary: request.IsPrimary,
	})
	if !ok {
		return domain.Contact{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "CUSTOMER",
		ActorID:     customer.ID,
		Actor:       customer.Name,
		Action:      "customer.contact.self_update",
		TargetType:  "customer",
		TargetID:    customer.ID,
		Target:      customer.CustomerNo,
		RequestID:   requestID,
		Description: "客户在门户更新联系人",
		Payload: map[string]any{
			"contactId":   contact.ID,
			"contactName": contact.Name,
		},
	})
	return contact, true
}

func (service *Service) UpdateContact(customerID, contactID int64, request dto.UpdateContactRequest, adminID int64, adminName, requestID string) (domain.Contact, bool) {
	customer, contact, ok := service.repository.UpdateContact(customerID, contactID, domain.Contact{
		Name:      request.Name,
		Email:     request.Email,
		Mobile:    request.Mobile,
		RoleName:  request.RoleName,
		IsPrimary: request.IsPrimary,
	})
	if !ok {
		return domain.Contact{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.contact.update",
		TargetType:  "customer",
		TargetID:    customer.ID,
		Target:      customer.CustomerNo,
		RequestID:   requestID,
		Description: "更新客户联系人",
		Payload: map[string]any{
			"contactId":   contact.ID,
			"contactName": contact.Name,
		},
	})
	return contact, true
}

func (service *Service) DeleteContactByCustomer(customerID, contactID int64, requestID string) bool {
	customer, ok := service.repository.DeleteContact(customerID, contactID)
	if !ok {
		return false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "CUSTOMER",
		ActorID:     customer.ID,
		Actor:       customer.Name,
		Action:      "customer.contact.self_delete",
		TargetType:  "customer",
		TargetID:    customer.ID,
		Target:      customer.CustomerNo,
		RequestID:   requestID,
		Description: "客户在门户删除联系人",
		Payload: map[string]any{
			"contactId": contactID,
		},
	})
	return true
}

func (service *Service) DeleteContact(customerID, contactID int64, adminID int64, adminName, requestID string) bool {
	customer, ok := service.repository.DeleteContact(customerID, contactID)
	if !ok {
		return false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.contact.delete",
		TargetType:  "customer",
		TargetID:    customer.ID,
		Target:      customer.CustomerNo,
		RequestID:   requestID,
		Description: "删除客户联系人",
		Payload: map[string]any{
			"contactId": contactID,
		},
	})
	return true
}

func (service *Service) SubmitIdentityByCustomer(customerID int64, request dto.SubmitIdentityRequest, requestID string) (domain.Identity, bool) {
	updatedCustomer, ok := service.repository.SubmitIdentity(customerID, domain.Identity{
		IdentityType: request.IdentityType,
		SubjectName:  request.SubjectName,
		CertNo:       request.CertNo,
		CountryCode:  request.CountryCode,
	})
	if !ok || updatedCustomer.Identity == nil {
		return domain.Identity{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "CUSTOMER",
		ActorID:     updatedCustomer.ID,
		Actor:       updatedCustomer.Name,
		Action:      "customer.identity.submit",
		TargetType:  "customer",
		TargetID:    updatedCustomer.ID,
		Target:      updatedCustomer.CustomerNo,
		RequestID:   requestID,
		Description: "客户在门户提交实名认证资料",
		Payload: map[string]any{
			"identityType": updatedCustomer.Identity.IdentityType,
			"subjectName":  updatedCustomer.Identity.SubjectName,
		},
	})
	return *updatedCustomer.Identity, true
}

func (service *Service) ReviewIdentity(customerID int64, request dto.ReviewIdentityRequest, adminID int64, adminName, requestID string) (domain.Identity, bool) {
	reviewRemark := request.ReviewRemark
	if reviewRemark == "" {
		reviewRemark = request.Remark
	}

	updatedCustomer, ok := service.repository.ReviewIdentity(customerID, domain.IdentityStatus(request.Status), reviewRemark)
	if !ok || updatedCustomer.Identity == nil {
		return domain.Identity{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "customer.identity.review",
		TargetType:  "customer",
		TargetID:    updatedCustomer.ID,
		Target:      updatedCustomer.CustomerNo,
		RequestID:   requestID,
		Description: "审核客户实名信息",
		Payload: map[string]any{
			"status": updatedCustomer.Identity.VerifyStatus,
			"remark": reviewRemark,
		},
	})
	return *updatedCustomer.Identity, true
}
