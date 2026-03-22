package repository

import (
	"fmt"
	"slices"
	"sync"

	"idc-finance/internal/modules/catalog/domain"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	items  []domain.Product
	nextID int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 4,
		items:  seedProducts(),
	}
}

func (repository *MemoryRepository) List() []domain.Product {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	return cloneProducts(repository.items)
}

func (repository *MemoryRepository) ListActive() []domain.Product {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	result := make([]domain.Product, 0)
	for _, item := range repository.items {
		if item.Status == domain.ProductStatusActive {
			result = append(result, cloneProduct(item))
		}
	}
	return result
}

func (repository *MemoryRepository) ListGroups() []string {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	set := make(map[string]struct{})
	for _, item := range repository.items {
		set[item.GroupName] = struct{}{}
	}

	groups := make([]string, 0, len(set))
	for group := range set {
		groups = append(groups, group)
	}
	slices.Sort(groups)
	return groups
}

func (repository *MemoryRepository) GetByID(id int64) (domain.Product, bool) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	for _, item := range repository.items {
		if item.ID == id {
			return cloneProduct(item), true
		}
	}
	return domain.Product{}, false
}

func (repository *MemoryRepository) Create(product domain.Product) domain.Product {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	product.ID = repository.nextID
	repository.nextID++
	product.ProductNo = fmt.Sprintf("PROD-%06d", product.ID)
	product = cloneProduct(product)
	repository.items = append(repository.items, product)
	return cloneProduct(product)
}

func (repository *MemoryRepository) Update(id int64, product domain.Product) (domain.Product, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	for index, item := range repository.items {
		if item.ID != id {
			continue
		}
		product.ID = item.ID
		product.ProductNo = item.ProductNo
		product = cloneProduct(product)
		repository.items[index] = product
		return cloneProduct(product), true
	}
	return domain.Product{}, false
}

func seedProducts() []domain.Product {
	return []domain.Product{
		{
			ID:          1,
			ProductNo:   "PROD-ELASTIC-001",
			GroupName:   "云主机",
			Name:        "弹性云主机 CN2 标准型",
			Description: "适合中小型 IDC 客户业务部署，含 1 个公网 IPv4，可按月或按年售卖。",
			ProductType: "CLOUD",
			Status:      domain.ProductStatusActive,
			Pricing: []domain.PriceOption{
				{CycleCode: "monthly", CycleName: "月付", Price: 199, SetupFee: 0},
				{CycleCode: "quarterly", CycleName: "季付", Price: 569, SetupFee: 0},
				{CycleCode: "annual", CycleName: "年付", Price: 1999, SetupFee: 0},
			},
			ConfigOptions: []domain.ConfigOption{
				{
					Code:         "cpu",
					Name:         "CPU 规格",
					InputType:    "select",
					Required:     true,
					DefaultValue: "4",
					Description:  "影响实例可分配的 vCPU 核数。",
					Choices: []domain.ConfigOptionChoice{
						{Value: "2", Label: "2 核", PriceDelta: 0},
						{Value: "4", Label: "4 核", PriceDelta: 80},
						{Value: "8", Label: "8 核", PriceDelta: 220},
					},
				},
				{
					Code:         "memory",
					Name:         "内存规格",
					InputType:    "select",
					Required:     true,
					DefaultValue: "8",
					Description:  "按 GB 计费，可和 CPU 规格组合售卖。",
					Choices: []domain.ConfigOptionChoice{
						{Value: "4", Label: "4 GB", PriceDelta: 0},
						{Value: "8", Label: "8 GB", PriceDelta: 60},
						{Value: "16", Label: "16 GB", PriceDelta: 180},
					},
				},
				{
					Code:         "backup",
					Name:         "云备份",
					InputType:    "radio",
					Required:     false,
					DefaultValue: "enabled",
					Description:  "是否附带每周快照备份服务。",
					Choices: []domain.ConfigOptionChoice{
						{Value: "enabled", Label: "启用", PriceDelta: 30},
						{Value: "disabled", Label: "关闭", PriceDelta: 0},
					},
				},
			},
			ResourceTemplate: domain.ResourceTemplate{
				RegionName:      "华南广州",
				ZoneName:        "广州三区",
				OperatingSystem: "Rocky Linux 9",
				LoginUsername:   "root",
				SecurityGroup:   "default-cloud",
				CPUCores:        4,
				MemoryGB:        8,
				SystemDiskGB:    60,
				DataDiskGB:      120,
				BandwidthMbps:   20,
				PublicIPCount:   1,
			},
			AutomationConfig: domain.AutomationConfig{
				Channel:        "MOFANG_CLOUD",
				ModuleType:     "DCIMCLOUD",
				ProvisionStage: "AFTER_PAYMENT",
				AutoProvision:  true,
				ServerGroup:    "mf-cloud-default",
			},
			UpstreamMapping: domain.UpstreamMapping{
				ProviderType:      "MOFANG_CLOUD",
				SourceName:        "魔方云",
				RemoteProductCode: "mf-cloud-standard",
				RemoteProductName: "弹性云主机标准型",
				PricePolicy:       "CUSTOM",
				AutoSyncPricing:   true,
				AutoSyncConfig:    true,
				AutoSyncTemplate:  true,
				SyncStatus:        "SUCCESS",
				SyncMessage:       "已预置示例映射",
				LastSyncedAt:      "2026-03-21T10:00:00+08:00",
			},
		},
		{
			ID:          2,
			ProductNo:   "PROD-BW-001",
			GroupName:   "网络",
			Name:        "精品带宽 50M",
			Description: "适合高并发 Web 业务，支持后续扩容和带宽包升级。",
			ProductType: "BANDWIDTH",
			Status:      domain.ProductStatusActive,
			Pricing: []domain.PriceOption{
				{CycleCode: "monthly", CycleName: "月付", Price: 299, SetupFee: 50},
				{CycleCode: "annual", CycleName: "年付", Price: 2990, SetupFee: 0},
			},
			ConfigOptions: []domain.ConfigOption{
				{
					Code:         "commit",
					Name:         "承诺带宽",
					InputType:    "select",
					Required:     true,
					DefaultValue: "50",
					Description:  "固定承诺带宽。",
					Choices: []domain.ConfigOptionChoice{
						{Value: "50", Label: "50 Mbps", PriceDelta: 0},
						{Value: "100", Label: "100 Mbps", PriceDelta: 260},
						{Value: "200", Label: "200 Mbps", PriceDelta: 680},
					},
				},
			},
			ResourceTemplate: domain.ResourceTemplate{
				RegionName:      "华南广州",
				ZoneName:        "骨干网络",
				OperatingSystem: "-",
				LoginUsername:   "-",
				SecurityGroup:   "network-edge",
				BandwidthMbps:   50,
				PublicIPCount:   1,
			},
			AutomationConfig: domain.AutomationConfig{
				Channel:        "LOCAL",
				ModuleType:     "NORMAL",
				ProvisionStage: "AFTER_PAYMENT",
				AutoProvision:  true,
			},
		},
		{
			ID:          3,
			ProductNo:   "PROD-COLO-001",
			GroupName:   "机柜托管",
			Name:        "1U 标准托管",
			Description: "适合物理服务器托管业务，含基础运维巡检与上架服务。",
			ProductType: "COLOCATION",
			Status:      domain.ProductStatusActive,
			Pricing: []domain.PriceOption{
				{CycleCode: "monthly", CycleName: "月付", Price: 499, SetupFee: 200},
				{CycleCode: "annual", CycleName: "年付", Price: 4990, SetupFee: 0},
			},
			ConfigOptions: []domain.ConfigOption{
				{
					Code:         "power",
					Name:         "电力额度",
					InputType:    "select",
					Required:     true,
					DefaultValue: "300",
					Description:  "默认含 A/B 路冗余电力。",
					Choices: []domain.ConfigOptionChoice{
						{Value: "300", Label: "300W", PriceDelta: 0},
						{Value: "500", Label: "500W", PriceDelta: 120},
						{Value: "800", Label: "800W", PriceDelta: 280},
					},
				},
			},
			ResourceTemplate: domain.ResourceTemplate{
				RegionName:    "华东上海",
				ZoneName:      "A1 机房",
				SecurityGroup: "cabinet-standard",
				BandwidthMbps: 20,
				PublicIPCount: 1,
			},
			AutomationConfig: domain.AutomationConfig{
				Channel:        "LOCAL",
				ModuleType:     "DCIM",
				ProvisionStage: "AFTER_PAYMENT",
				AutoProvision:  true,
			},
		},
	}
}

func cloneProducts(items []domain.Product) []domain.Product {
	result := make([]domain.Product, 0, len(items))
	for _, item := range items {
		result = append(result, cloneProduct(item))
	}
	return result
}

func cloneProduct(item domain.Product) domain.Product {
	item.Pricing = slices.Clone(item.Pricing)
	item.ConfigOptions = cloneConfigOptions(item.ConfigOptions)
	return item
}

func cloneConfigOptions(items []domain.ConfigOption) []domain.ConfigOption {
	result := make([]domain.ConfigOption, 0, len(items))
	for _, item := range items {
		copyItem := item
		copyItem.Choices = slices.Clone(item.Choices)
		result = append(result, copyItem)
	}
	return result
}
