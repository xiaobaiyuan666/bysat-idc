package repository

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"idc-finance/internal/modules/customer/domain"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	items  []domain.Customer
	groups []domain.CustomerGroup
	levels []domain.CustomerLevel
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: []domain.Customer{
			{
				ID:         1,
				CustomerNo: "CUST-20260320-0001",
				Name:       "星环网络",
				Email:      "ops@starring.cn",
				Mobile:     "13800000001",
				Type:       "COMPANY",
				Status:     domain.CustomerStatusActive,
				GroupName:  "旗舰客户",
				LevelName:  "VIP",
				SalesOwner: "华东销售一组",
				Remarks:    "重点云业务客户",
				Contacts: []domain.Contact{
					{ID: 1, Name: "林浩", Email: "ops@starring.cn", Mobile: "13800000001", RoleName: "主联系人", IsPrimary: true},
					{ID: 2, Name: "陈会", Email: "finance@starring.cn", Mobile: "13800000002", RoleName: "财务联系人", IsPrimary: false},
				},
				Identity: &domain.Identity{
					ID:           1,
					IdentityType: "COMPANY",
					VerifyStatus: domain.IdentityStatusApproved,
					SubjectName:  "上海星环网络科技有限公司",
					CertNo:       "91310000MA1K123456",
					CountryCode:  "CN",
					ReviewRemark: "历史实名数据已完成核验导入",
					ReviewedAt:   "2026-03-20 18:00:00",
					SubmittedAt:  "2026-03-20 10:00:00",
				},
			},
			{
				ID:         2,
				CustomerNo: "CUST-20260320-0002",
				Name:       "云脉互联",
				Email:      "admin@cloudpulse.cn",
				Mobile:     "13800000003",
				Type:       "COMPANY",
				Status:     domain.CustomerStatusActive,
				GroupName:  "云业务客户",
				LevelName:  "标准",
				SalesOwner: "华南销售二组",
				Remarks:    "以弹性云和带宽为主",
				Contacts: []domain.Contact{
					{ID: 3, Name: "周越", Email: "admin@cloudpulse.cn", Mobile: "13800000003", RoleName: "主联系人", IsPrimary: true},
				},
				Identity: &domain.Identity{
					ID:           2,
					IdentityType: "COMPANY",
					VerifyStatus: domain.IdentityStatusPending,
					SubjectName:  "深圳云脉互联科技有限公司",
					CertNo:       "91440300MA5P654321",
					CountryCode:  "CN",
					ReviewRemark: "",
					ReviewedAt:   "",
					SubmittedAt:  "2026-03-21 09:30:00",
				},
			},
		},
		groups: []domain.CustomerGroup{
			{ID: 1, Name: "旗舰客户", Description: "高价值 IDC 与云资源客户"},
			{ID: 2, Name: "云业务客户", Description: "以弹性云、带宽和高防产品为主"},
		},
		levels: []domain.CustomerLevel{
			{ID: 1, Name: "VIP", Priority: 100, Description: "重点维护客户"},
			{ID: 2, Name: "标准", Priority: 50, Description: "默认标准等级"},
		},
	}
}

func (repository *MemoryRepository) List() []domain.Customer {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	return slices.Clone(repository.items)
}

func (repository *MemoryRepository) ListGroups() []domain.CustomerGroup {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	result := make([]domain.CustomerGroup, 0, len(repository.groups))
	for _, item := range repository.groups {
		item.CustomerCount = repository.countCustomersByGroupLocked(item.Name)
		result = append(result, item)
	}
	return result
}

func (repository *MemoryRepository) ListLevels() []domain.CustomerLevel {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	result := make([]domain.CustomerLevel, 0, len(repository.levels))
	for _, item := range repository.levels {
		item.CustomerCount = repository.countCustomersByLevelLocked(item.Name)
		result = append(result, item)
	}
	return result
}

func (repository *MemoryRepository) CreateGroup(group domain.CustomerGroup) (domain.CustomerGroup, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	group.ID = repository.nextGroupIDLocked()
	repository.groups = append(repository.groups, group)
	group.CustomerCount = 0
	return group, nil
}

func (repository *MemoryRepository) UpdateGroup(id int64, group domain.CustomerGroup) (domain.CustomerGroup, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.groups {
		if item.ID != id {
			continue
		}
		oldName := item.Name
		item.Name = group.Name
		item.Description = group.Description
		repository.groups[index] = item

		if oldName != group.Name {
			for customerIndex := range repository.items {
				if repository.items[customerIndex].GroupName == oldName {
					repository.items[customerIndex].GroupName = group.Name
				}
			}
		}

		item.CustomerCount = repository.countCustomersByGroupLocked(item.Name)
		return item, nil
	}
	return domain.CustomerGroup{}, ErrGroupNotFound
}

func (repository *MemoryRepository) DeleteGroup(id int64) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.groups {
		if item.ID != id {
			continue
		}
		if repository.countCustomersByGroupLocked(item.Name) > 0 {
			return ErrGroupInUse
		}
		repository.groups = append(repository.groups[:index], repository.groups[index+1:]...)
		return nil
	}
	return ErrGroupNotFound
}

func (repository *MemoryRepository) CreateLevel(level domain.CustomerLevel) (domain.CustomerLevel, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	level.ID = repository.nextLevelIDLocked()
	repository.levels = append(repository.levels, level)
	level.CustomerCount = 0
	return level, nil
}

func (repository *MemoryRepository) UpdateLevel(id int64, level domain.CustomerLevel) (domain.CustomerLevel, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.levels {
		if item.ID != id {
			continue
		}
		oldName := item.Name
		item.Name = level.Name
		item.Priority = level.Priority
		item.Description = level.Description
		repository.levels[index] = item

		if oldName != level.Name {
			for customerIndex := range repository.items {
				if repository.items[customerIndex].LevelName == oldName {
					repository.items[customerIndex].LevelName = level.Name
				}
			}
		}

		item.CustomerCount = repository.countCustomersByLevelLocked(item.Name)
		return item, nil
	}
	return domain.CustomerLevel{}, ErrLevelNotFound
}

func (repository *MemoryRepository) DeleteLevel(id int64) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.levels {
		if item.ID != id {
			continue
		}
		if repository.countCustomersByLevelLocked(item.Name) > 0 {
			return ErrLevelInUse
		}
		repository.levels = append(repository.levels[:index], repository.levels[index+1:]...)
		return nil
	}
	return ErrLevelNotFound
}

func (repository *MemoryRepository) GetByID(id int64) (domain.Customer, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	for _, item := range repository.items {
		if item.ID == id {
			return item, true
		}
	}
	return domain.Customer{}, false
}

func (repository *MemoryRepository) Create(customer domain.Customer) domain.Customer {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	customer.ID = int64(len(repository.items) + 1)
	customer.CustomerNo = fmt.Sprintf("CUST-20260320-%04d", customer.ID+2)
	customer.Status = domain.CustomerStatusActive
	customer.Contacts = []domain.Contact{}
	customer.Identity = &domain.Identity{
		ID:           customer.ID,
		IdentityType: customer.Type,
		VerifyStatus: domain.IdentityStatusPending,
		SubjectName:  customer.Name,
		CountryCode:  "CN",
		SubmittedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}
	repository.items = append(repository.items, customer)
	return customer
}

func (repository *MemoryRepository) Update(id int64, customer domain.Customer) (domain.Customer, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for index, item := range repository.items {
		if item.ID == id {
			customer.ID = item.ID
			customer.CustomerNo = item.CustomerNo
			customer.Type = item.Type
			customer.Contacts = item.Contacts
			customer.Identity = item.Identity
			repository.items[index] = customer
			return customer, true
		}
	}
	return domain.Customer{}, false
}

func (repository *MemoryRepository) AddContact(customerID int64, contact domain.Contact) (domain.Customer, domain.Contact, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for customerIndex, item := range repository.items {
		if item.ID != customerID {
			continue
		}

		if contact.IsPrimary {
			for index := range item.Contacts {
				item.Contacts[index].IsPrimary = false
			}
		}
		contact.ID = int64(len(item.Contacts) + 1)
		item.Contacts = append(item.Contacts, contact)
		repository.items[customerIndex] = item
		return item, contact, true
	}
	return domain.Customer{}, domain.Contact{}, false
}

func (repository *MemoryRepository) UpdateContact(customerID, contactID int64, contact domain.Contact) (domain.Customer, domain.Contact, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for customerIndex, item := range repository.items {
		if item.ID != customerID {
			continue
		}
		for contactIndex, existing := range item.Contacts {
			if existing.ID != contactID {
				continue
			}
			contact.ID = contactID
			if contact.IsPrimary {
				for resetIndex := range item.Contacts {
					item.Contacts[resetIndex].IsPrimary = false
				}
			} else {
				contact.IsPrimary = existing.IsPrimary
			}
			item.Contacts[contactIndex] = contact
			repository.items[customerIndex] = item
			return item, contact, true
		}
	}
	return domain.Customer{}, domain.Contact{}, false
}

func (repository *MemoryRepository) DeleteContact(customerID, contactID int64) (domain.Customer, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for customerIndex, item := range repository.items {
		if item.ID != customerID {
			continue
		}
		for contactIndex, contact := range item.Contacts {
			if contact.ID != contactID {
				continue
			}
			item.Contacts = append(item.Contacts[:contactIndex], item.Contacts[contactIndex+1:]...)
			if contact.IsPrimary && len(item.Contacts) > 0 {
				item.Contacts[0].IsPrimary = true
			}
			repository.items[customerIndex] = item
			return item, true
		}
	}
	return domain.Customer{}, false
}

func (repository *MemoryRepository) SubmitIdentity(customerID int64, identity domain.Identity) (domain.Customer, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for customerIndex, item := range repository.items {
		if item.ID != customerID {
			continue
		}
		if item.Identity == nil {
			identity.ID = customerID
		} else {
			identity.ID = item.Identity.ID
		}
		identity.VerifyStatus = domain.IdentityStatusPending
		identity.ReviewRemark = ""
		identity.ReviewedAt = ""
		identity.SubmittedAt = time.Now().Format("2006-01-02 15:04:05")
		item.Identity = &identity
		repository.items[customerIndex] = item
		return item, true
	}
	return domain.Customer{}, false
}

func (repository *MemoryRepository) ReviewIdentity(customerID int64, status domain.IdentityStatus, reviewRemark string) (domain.Customer, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for customerIndex, item := range repository.items {
		if item.ID != customerID || item.Identity == nil {
			continue
		}
		item.Identity.VerifyStatus = status
		item.Identity.ReviewRemark = reviewRemark
		item.Identity.ReviewedAt = time.Now().Format("2006-01-02 15:04:05")
		repository.items[customerIndex] = item
		return item, true
	}
	return domain.Customer{}, false
}

func (repository *MemoryRepository) ListServiceItems(customerID int64) []domain.RelatedItem {
	if _, exists := repository.GetByID(customerID); !exists {
		return nil
	}
	return []domain.RelatedItem{
		{
			No:          "SRV-DEMO-0001",
			Name:        "弹性云主机 CN2 标准型",
			Status:      "ACTIVE",
			DueAt:       "2026-04-20",
			Description: "后续接入真实服务模块",
		},
	}
}

func (repository *MemoryRepository) ListInvoiceItems(customerID int64) []domain.RelatedItem {
	if _, exists := repository.GetByID(customerID); !exists {
		return nil
	}
	return []domain.RelatedItem{
		{
			No:     "INV-DEMO-0001",
			Name:   "云主机新购账单",
			Status: "UNPAID",
			Amount: "199.00",
			DueAt:  "2026-03-30",
		},
	}
}

func (repository *MemoryRepository) ListTicketItems(customerID int64) []domain.RelatedItem {
	if _, exists := repository.GetByID(customerID); !exists {
		return nil
	}
	return []domain.RelatedItem{
		{
			No:        "TIC-DEMO-0001",
			Name:      "实例网络异常排查",
			Status:    "PROCESSING",
			UpdatedAt: "2026-03-20 18:20:00",
		},
	}
}

func (repository *MemoryRepository) countCustomersByGroupLocked(groupName string) int {
	total := 0
	for _, item := range repository.items {
		if item.GroupName == groupName {
			total++
		}
	}
	return total
}

func (repository *MemoryRepository) countCustomersByLevelLocked(levelName string) int {
	total := 0
	for _, item := range repository.items {
		if item.LevelName == levelName {
			total++
		}
	}
	return total
}

func (repository *MemoryRepository) nextGroupIDLocked() int64 {
	var maxID int64
	for _, item := range repository.groups {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}

func (repository *MemoryRepository) nextLevelIDLocked() int64 {
	var maxID int64
	for _, item := range repository.levels {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}
