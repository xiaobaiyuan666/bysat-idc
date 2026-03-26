package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"idc-finance/internal/platform/audit"
	"idc-finance/internal/platform/rbac"
)

type AdminMember struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"displayName"`
	Email       string  `json:"email"`
	Mobile      string  `json:"mobile"`
	Status      string  `json:"status"`
	RoleIDs     []int64 `json:"roleIds"`
	Roles       []string `json:"roles"`
	RoleNames   []string `json:"roleNames"`
	LastLoginAt string  `json:"lastLoginAt"`
	LastLoginIP string  `json:"lastLoginIp"`
	Remarks     string  `json:"remarks"`
}

type Role struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	MenuIDs     []int64  `json:"menuIds"`
	Permissions []string `json:"permissions"`
	Users       int      `json:"users"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type SaveAdminInput struct {
	Username    string  `json:"username"`
	DisplayName string  `json:"displayName"`
	Email       string  `json:"email"`
	Mobile      string  `json:"mobile"`
	Status      string  `json:"status"`
	RoleIDs     []int64 `json:"roleIds"`
	Remarks     string  `json:"remarks"`
}

type SaveRoleInput struct {
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	MenuIDs     []int64  `json:"menuIds"`
	Permissions []string `json:"permissions"`
}

type state struct {
	NextAdminID int64         `json:"nextAdminId"`
	NextRoleID  int64         `json:"nextRoleId"`
	Admins      []AdminMember `json:"admins"`
	Roles       []Role        `json:"roles"`
}

type Service struct {
	mu        sync.RWMutex
	stateFile string
	audit     *audit.Service
	menus     []rbac.Menu
	state     state
}

func PermissionCodes() []string {
	values := append([]string{}, extraPermissionCodes()...)
	for _, menu := range rbac.PhaseOneMenus() {
		if menu.Permission != "" {
			values = append(values, menu.Permission)
		}
	}
	return uniqueSortedStrings(values)
}

func NewService(stateFile string, auditService *audit.Service) *Service {
	service := &Service{
		stateFile: stateFile,
		audit:     auditService,
		menus:     rbac.PhaseOneMenus(),
	}
	service.load()
	return service
}

func (service *Service) ListAdmins() []AdminMember {
	service.mu.RLock()
	defer service.mu.RUnlock()
	return service.snapshotAdminsLocked()
}

func (service *Service) CreateAdmin(input SaveAdminInput) (AdminMember, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	admin, err := service.buildAdminLocked(0, input)
	if err != nil {
		return AdminMember{}, err
	}

	admin.ID = service.state.NextAdminID
	service.state.NextAdminID++
	service.state.Admins = append(service.state.Admins, admin)
	if err := service.saveLocked(); err != nil {
		return AdminMember{}, err
	}

	service.recordAudit("rbac.admin.create", "admin", admin.ID, admin.Username, fmt.Sprintf("创建管理员 %s", admin.DisplayName), map[string]any{
		"username": admin.Username,
		"roleIds":  admin.RoleIDs,
		"status":   admin.Status,
	})
	return service.adminViewLocked(admin), nil
}

func (service *Service) UpdateAdmin(id int64, input SaveAdminInput) (AdminMember, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	index := -1
	for idx, item := range service.state.Admins {
		if item.ID == id {
			index = idx
			break
		}
	}
	if index == -1 {
		return AdminMember{}, fmt.Errorf("管理员不存在")
	}

	admin, err := service.buildAdminLocked(id, input)
	if err != nil {
		return AdminMember{}, err
	}
	admin.ID = id
	admin.LastLoginAt = service.state.Admins[index].LastLoginAt
	admin.LastLoginIP = service.state.Admins[index].LastLoginIP
	service.state.Admins[index] = admin
	if err := service.saveLocked(); err != nil {
		return AdminMember{}, err
	}

	service.recordAudit("rbac.admin.update", "admin", admin.ID, admin.Username, fmt.Sprintf("更新管理员 %s", admin.DisplayName), map[string]any{
		"username": admin.Username,
		"roleIds":  admin.RoleIDs,
		"status":   admin.Status,
	})
	return service.adminViewLocked(admin), nil
}

func (service *Service) ListRoles() []Role {
	service.mu.RLock()
	defer service.mu.RUnlock()
	return service.snapshotRolesLocked()
}

func (service *Service) CreateRole(input SaveRoleInput) (Role, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	role, err := service.buildRoleLocked(0, input)
	if err != nil {
		return Role{}, err
	}

	role.ID = service.state.NextRoleID
	service.state.NextRoleID++
	service.state.Roles = append(service.state.Roles, role)
	if err := service.saveLocked(); err != nil {
		return Role{}, err
	}

	service.recordAudit("rbac.role.create", "role", role.ID, role.Code, fmt.Sprintf("创建角色 %s", role.Name), map[string]any{
		"code":        role.Code,
		"permissions": role.Permissions,
		"menuIds":     role.MenuIDs,
	})
	return service.roleViewLocked(role), nil
}

func (service *Service) UpdateRole(id int64, input SaveRoleInput) (Role, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	index := -1
	for idx, item := range service.state.Roles {
		if item.ID == id {
			index = idx
			break
		}
	}
	if index == -1 {
		return Role{}, fmt.Errorf("角色不存在")
	}

	role, err := service.buildRoleLocked(id, input)
	if err != nil {
		return Role{}, err
	}
	role.ID = id
	role.CreatedAt = service.state.Roles[index].CreatedAt
	service.state.Roles[index] = role
	if err := service.saveLocked(); err != nil {
		return Role{}, err
	}

	service.recordAudit("rbac.role.update", "role", role.ID, role.Code, fmt.Sprintf("更新角色 %s", role.Name), map[string]any{
		"code":        role.Code,
		"permissions": role.Permissions,
		"menuIds":     role.MenuIDs,
		"status":      role.Status,
	})
	return service.roleViewLocked(role), nil
}

func (service *Service) load() {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.state = defaultState(service.menus)
	if service.stateFile == "" {
		return
	}

	content, err := os.ReadFile(service.stateFile)
	if err != nil {
		_ = service.saveLocked()
		return
	}

	var stored state
	if err := json.Unmarshal(content, &stored); err != nil {
		_ = service.saveLocked()
		return
	}

	service.state = stored
	service.state = service.normalizeState(service.state)
}

func (service *Service) normalizeState(stored state) state {
	if stored.NextAdminID <= 0 {
		stored.NextAdminID = 1
	}
	if stored.NextRoleID <= 0 {
		stored.NextRoleID = 1
	}

	for index, item := range stored.Roles {
		stored.Roles[index] = service.roleViewLocked(item)
		if stored.Roles[index].ID >= stored.NextRoleID {
			stored.NextRoleID = stored.Roles[index].ID + 1
		}
	}

	for index, item := range stored.Admins {
		stored.Admins[index] = service.adminViewLocked(item)
		if stored.Admins[index].ID >= stored.NextAdminID {
			stored.NextAdminID = stored.Admins[index].ID + 1
		}
	}

	return stored
}

func (service *Service) buildAdminLocked(currentID int64, input SaveAdminInput) (AdminMember, error) {
	username := strings.TrimSpace(input.Username)
	displayName := strings.TrimSpace(input.DisplayName)
	if username == "" {
		return AdminMember{}, fmt.Errorf("登录账号不能为空")
	}
	if displayName == "" {
		return AdminMember{}, fmt.Errorf("显示名称不能为空")
	}

	for _, item := range service.state.Admins {
		if item.ID != currentID && strings.EqualFold(item.Username, username) {
			return AdminMember{}, fmt.Errorf("登录账号已存在")
		}
	}

	roleIDs := uniqueSortedInt64s(input.RoleIDs)
	if len(roleIDs) == 0 {
		return AdminMember{}, fmt.Errorf("至少选择一个角色")
	}

	roleCodes, roleNames, err := service.resolveRoleSummaryLocked(roleIDs)
	if err != nil {
		return AdminMember{}, err
	}

	status := strings.ToUpper(strings.TrimSpace(input.Status))
	if status == "" {
		status = "ACTIVE"
	}

	return AdminMember{
		Username:    username,
		DisplayName: displayName,
		Email:       strings.TrimSpace(input.Email),
		Mobile:      strings.TrimSpace(input.Mobile),
		Status:      status,
		RoleIDs:     roleIDs,
		Roles:       roleCodes,
		RoleNames:   roleNames,
		Remarks:     strings.TrimSpace(input.Remarks),
	}, nil
}

func (service *Service) buildRoleLocked(currentID int64, input SaveRoleInput) (Role, error) {
	name := strings.TrimSpace(input.Name)
	code := strings.TrimSpace(input.Code)
	if name == "" {
		return Role{}, fmt.Errorf("角色名称不能为空")
	}
	if code == "" {
		return Role{}, fmt.Errorf("角色编码不能为空")
	}

	for _, item := range service.state.Roles {
		if item.ID != currentID && strings.EqualFold(item.Code, code) {
			return Role{}, fmt.Errorf("角色编码已存在")
		}
	}

	menuIDs := service.validMenuIDs(input.MenuIDs)
	permissions := service.mergeRolePermissionsLocked(menuIDs, input.Permissions)
	status := strings.ToUpper(strings.TrimSpace(input.Status))
	if status == "" {
		status = "ACTIVE"
	}
	now := nowString()
	createdAt := now
	if currentID > 0 {
		for _, item := range service.state.Roles {
			if item.ID == currentID {
				createdAt = item.CreatedAt
				break
			}
		}
	}

	return Role{
		Name:        name,
		Code:        code,
		Description: strings.TrimSpace(input.Description),
		Status:      status,
		MenuIDs:     menuIDs,
		Permissions: permissions,
		CreatedAt:   createdAt,
		UpdatedAt:   now,
	}, nil
}

func (service *Service) resolveRoleSummaryLocked(roleIDs []int64) ([]string, []string, error) {
	codeByID := make(map[int64]string, len(service.state.Roles))
	nameByID := make(map[int64]string, len(service.state.Roles))
	for _, item := range service.state.Roles {
		codeByID[item.ID] = item.Code
		nameByID[item.ID] = item.Name
	}

	roleCodes := make([]string, 0, len(roleIDs))
	roleNames := make([]string, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		code, ok := codeByID[roleID]
		if !ok {
			return nil, nil, fmt.Errorf("角色 #%d 不存在", roleID)
		}
		roleCodes = append(roleCodes, code)
		roleNames = append(roleNames, nameByID[roleID])
	}
	return roleCodes, roleNames, nil
}

func (service *Service) snapshotAdminsLocked() []AdminMember {
	items := make([]AdminMember, 0, len(service.state.Admins))
	for _, item := range service.state.Admins {
		items = append(items, service.adminViewLocked(item))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
	return items
}

func (service *Service) snapshotRolesLocked() []Role {
	counts := make(map[int64]int)
	for _, admin := range service.state.Admins {
		for _, roleID := range admin.RoleIDs {
			counts[roleID]++
		}
	}

	items := make([]Role, 0, len(service.state.Roles))
	for _, item := range service.state.Roles {
		view := service.roleViewLocked(item)
		view.Users = counts[item.ID]
		items = append(items, view)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
	return items
}

func (service *Service) adminViewLocked(item AdminMember) AdminMember {
	roleCodes, roleNames, _ := service.resolveRoleSummaryLocked(uniqueSortedInt64s(item.RoleIDs))
	item.RoleIDs = uniqueSortedInt64s(item.RoleIDs)
	item.Roles = roleCodes
	item.RoleNames = roleNames
	return item
}

func (service *Service) roleViewLocked(item Role) Role {
	item.MenuIDs = service.validMenuIDs(item.MenuIDs)
	item.Permissions = service.mergeRolePermissionsLocked(item.MenuIDs, item.Permissions)
	return item
}

func (service *Service) validMenuIDs(menuIDs []int64) []int64 {
	valid := make(map[int64]struct{}, len(service.menus))
	for _, item := range service.menus {
		valid[item.ID] = struct{}{}
	}
	results := make([]int64, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		if _, ok := valid[menuID]; ok {
			results = append(results, menuID)
		}
	}
	return uniqueSortedInt64s(results)
}

func (service *Service) mergeRolePermissionsLocked(menuIDs []int64, explicit []string) []string {
	permissionByMenuID := make(map[int64]string, len(service.menus))
	for _, item := range service.menus {
		permissionByMenuID[item.ID] = item.Permission
	}

	values := make([]string, 0, len(menuIDs)+len(explicit))
	for _, menuID := range menuIDs {
		if permission := strings.TrimSpace(permissionByMenuID[menuID]); permission != "" {
			values = append(values, permission)
		}
	}

	allowed := make(map[string]struct{}, len(PermissionCodes()))
	for _, item := range PermissionCodes() {
		allowed[item] = struct{}{}
	}
	for _, permission := range explicit {
		value := strings.TrimSpace(permission)
		if value == "" {
			continue
		}
		if _, ok := allowed[value]; ok {
			values = append(values, value)
		}
	}
	return uniqueSortedStrings(values)
}

func (service *Service) saveLocked() error {
	if service.stateFile == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(service.stateFile), 0o755); err != nil {
		return err
	}

	content, err := json.MarshalIndent(service.state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(service.stateFile, content, 0o644)
}

func (service *Service) recordAudit(action, targetType string, targetID int64, target, description string, payload map[string]any) {
	if service.audit == nil {
		return
	}
	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     1,
		Actor:       "系统管理员",
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Target:      target,
		Description: description,
		Payload:     payload,
	})
}

func defaultState(menus []rbac.Menu) state {
	roles := defaultRoles(menus)
	admins := []AdminMember{
		{
			ID:          1,
			Username:    "admin",
			DisplayName: "系统管理员",
			Email:       "admin@example.com",
			Mobile:      "13800000001",
			Status:      "ACTIVE",
			RoleIDs:     []int64{1},
			LastLoginAt: "2026-03-26 09:30:00",
			LastLoginIP: "127.0.0.1",
			Remarks:     "平台超级管理员",
		},
		{
			ID:          2,
			Username:    "support01",
			DisplayName: "客服主管",
			Email:       "support@example.com",
			Mobile:      "13800000002",
			Status:      "ACTIVE",
			RoleIDs:     []int64{2},
			LastLoginAt: "2026-03-25 18:20:00",
			LastLoginIP: "10.0.0.12",
			Remarks:     "负责客户、工单和订单跟进",
		},
		{
			ID:          3,
			Username:    "finance01",
			DisplayName: "财务专员",
			Email:       "finance@example.com",
			Mobile:      "13800000003",
			Status:      "ACTIVE",
			RoleIDs:     []int64{3},
			LastLoginAt: "2026-03-25 17:10:00",
			LastLoginIP: "10.0.0.18",
			Remarks:     "负责账单、支付、退款和资金台账",
		},
		{
			ID:          4,
			Username:    "ops01",
			DisplayName: "运维工程师",
			Email:       "ops@example.com",
			Mobile:      "13800000004",
			Status:      "ACTIVE",
			RoleIDs:     []int64{4},
			LastLoginAt: "2026-03-26 08:55:00",
			LastLoginIP: "10.0.0.24",
			Remarks:     "负责服务、资源、魔方云和自动化任务",
		},
	}

	return state{
		NextAdminID: 5,
		NextRoleID:  5,
		Admins:      admins,
		Roles:       roles,
	}
}

func defaultRoles(menus []rbac.Menu) []Role {
	allMenuIDs := make([]int64, 0, len(menus))
	for _, item := range menus {
		allMenuIDs = append(allMenuIDs, item.ID)
	}
	allPermissions := PermissionCodes()

	supportMenus := []int64{1, 10, 11, 12, 13, 14, 20, 21, 22, 30, 32, 33, 40, 41, 50, 52, 60, 61, 70, 74}
	financeMenus := []int64{1, 10, 11, 20, 21, 30, 31, 34, 35, 36, 37, 38, 60, 61, 70, 74}
	opsMenus := []int64{1, 20, 21, 22, 40, 41, 50, 51, 52, 53, 60, 61, 70, 74}

	return []Role{
		{
			ID:          1,
			Name:        "系统管理员",
			Code:        "super-admin",
			Description: "拥有平台全部菜单与权限，可维护系统与业务全局配置。",
			Status:      "ACTIVE",
			MenuIDs:     uniqueSortedInt64s(allMenuIDs),
			Permissions: allPermissions,
			Users:       1,
			CreatedAt:   "2026-03-20 18:00:00",
			UpdatedAt:   "2026-03-26 09:30:00",
		},
		{
			ID:          2,
			Name:        "客服主管",
			Code:        "support-manager",
			Description: "负责客户、工单、订单跟进与基础业务处理。",
			Status:      "ACTIVE",
			MenuIDs:     uniqueSortedInt64s(supportMenus),
			Permissions: uniqueSortedStrings(append(menuPermissionsByIDs(menus, supportMenus), "customer:view", "order:view", "service:view", "ticket:update")),
			Users:       1,
			CreatedAt:   "2026-03-20 18:10:00",
			UpdatedAt:   "2026-03-26 09:00:00",
		},
		{
			ID:          3,
			Name:        "财务专员",
			Code:        "finance-operator",
			Description: "负责账单、支付、退款、余额和财务报表处理。",
			Status:      "ACTIVE",
			MenuIDs:     uniqueSortedInt64s(financeMenus),
			Permissions: uniqueSortedStrings(append(menuPermissionsByIDs(menus, financeMenus), "invoice:receive", "invoice:refund")),
			Users:       1,
			CreatedAt:   "2026-03-20 18:12:00",
			UpdatedAt:   "2026-03-26 08:40:00",
		},
		{
			ID:          4,
			Name:        "运维工程师",
			Code:        "ops-engineer",
			Description: "负责服务、资源、接口账户、魔方云和自动化任务处理。",
			Status:      "ACTIVE",
			MenuIDs:     uniqueSortedInt64s(opsMenus),
			Permissions: uniqueSortedStrings(append(menuPermissionsByIDs(menus, opsMenus), "provider:action", "service:update")),
			Users:       1,
			CreatedAt:   "2026-03-20 18:15:00",
			UpdatedAt:   "2026-03-26 08:50:00",
		},
	}
}

func menuPermissionsByIDs(menus []rbac.Menu, ids []int64) []string {
	targets := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		targets[id] = struct{}{}
	}
	values := make([]string, 0, len(ids))
	for _, item := range menus {
		if _, ok := targets[item.ID]; ok && item.Permission != "" {
			values = append(values, item.Permission)
		}
	}
	return values
}

func extraPermissionCodes() []string {
	return []string{
		"customer:view",
		"customer:create",
		"customer:contact:list",
		"customer:contact:create",
		"customer:contact:update",
		"customer:contact:delete",
		"customer:identity:update",
		"catalog:product:create",
		"invoice:receive",
		"invoice:refund",
		"provider:action",
	}
}

func uniqueSortedStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	results := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		results = append(results, value)
	}
	sort.Strings(results)
	return results
}

func uniqueSortedInt64s(values []int64) []int64 {
	seen := make(map[int64]struct{}, len(values))
	results := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		results = append(results, value)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})
	return results
}

func nowString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
