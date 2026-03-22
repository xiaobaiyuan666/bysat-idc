package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"idc-finance/internal/modules/catalog/domain"
	"idc-finance/internal/modules/catalog/dto"
	"idc-finance/internal/modules/catalog/repository"
	providerService "idc-finance/internal/modules/provider/service"
	"idc-finance/internal/platform/audit"
)

type Service struct {
	repository repository.Repository
	audit      *audit.Service
	provider   *providerService.Service
}

func New(repo repository.Repository, auditService *audit.Service, provider *providerService.Service) *Service {
	return &Service{repository: repo, audit: auditService, provider: provider}
}

func (service *Service) List() []domain.Product {
	return service.repository.List()
}

func (service *Service) ListActive() []domain.Product {
	return service.repository.ListActive()
}

func (service *Service) ListGroups() []string {
	return service.repository.ListGroups()
}

func (service *Service) GetByID(id int64) (domain.Product, bool) {
	return service.repository.GetByID(id)
}

func (service *Service) GetDetail(id int64) (dto.ProductDetailResponse, bool) {
	product, exists := service.repository.GetByID(id)
	if !exists {
		return dto.ProductDetailResponse{}, false
	}
	return dto.ProductDetailResponse{
		Product:   product,
		AuditLogs: service.audit.ListByTarget("product", id),
	}, true
}

func (service *Service) Create(request dto.CreateProductRequest, adminID int64, adminName, requestID string) domain.Product {
	product := service.repository.Create(buildProductPayload(
		request.GroupName,
		request.Name,
		request.Description,
		request.ProductType,
		request.Status,
		request.Pricing,
		request.ConfigOptions,
		request.ResourceTemplate,
		request.AutomationConfig,
		request.UpstreamMapping,
	))

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "catalog.product.create",
		TargetType:  "product",
		TargetID:    product.ID,
		Target:      product.ProductNo,
		RequestID:   requestID,
		Description: "创建商品并保存价格矩阵、配置项和资源模板",
		Payload: map[string]any{
			"name":        product.Name,
			"groupName":   product.GroupName,
			"pricing":     len(product.Pricing),
			"configItems": len(product.ConfigOptions),
		},
	})

	return product
}

func (service *Service) Update(id int64, request dto.UpdateProductRequest, adminID int64, adminName, requestID string) (domain.Product, bool) {
	existing, ok := service.repository.GetByID(id)
	if !ok {
		return domain.Product{}, false
	}

	payload := buildProductPayload(
		request.GroupName,
		request.Name,
		request.Description,
		request.ProductType,
		request.Status,
		request.Pricing,
		request.ConfigOptions,
		request.ResourceTemplate,
		request.AutomationConfig,
		request.UpstreamMapping,
	)
	if isEmptyAutomationConfigInput(request.AutomationConfig) {
		payload.AutomationConfig = existing.AutomationConfig
	}
	payload.UpstreamMapping = mergeProductUpstreamMapping(existing.UpstreamMapping, request.UpstreamMapping)

	product, ok := service.repository.Update(id, payload)
	if !ok {
		return domain.Product{}, false
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "catalog.product.update",
		TargetType:  "product",
		TargetID:    product.ID,
		Target:      product.ProductNo,
		RequestID:   requestID,
		Description: "更新商品基础信息、价格矩阵和资源模板",
		Payload: map[string]any{
			"name":        product.Name,
			"groupName":   product.GroupName,
			"pricing":     len(product.Pricing),
			"configItems": len(product.ConfigOptions),
		},
	})

	return product, true
}

func (service *Service) SyncUpstream(
	id int64,
	request dto.SyncProductUpstreamRequest,
	adminID int64,
	adminName,
	requestID string,
) (domain.Product, bool, error) {
	product, exists := service.repository.GetByID(id)
	if !exists {
		return domain.Product{}, false, nil
	}

	mapping := normalizeUpstreamMapping(request)
	if mapping.ProviderType == "" || mapping.ProviderType == "NONE" {
		mapping.ProviderType = firstNonEmpty(product.UpstreamMapping.ProviderType, "MOFANG_CLOUD")
	}
	if mapping.SourceName == "" {
		mapping.SourceName = firstNonEmpty(product.UpstreamMapping.SourceName, defaultSourceName(mapping.ProviderType))
	}
	if mapping.ProviderAccountID == 0 {
		mapping.ProviderAccountID = product.UpstreamMapping.ProviderAccountID
	}
	if mapping.RemoteProductCode == "" {
		mapping.RemoteProductCode = product.UpstreamMapping.RemoteProductCode
	}
	if mapping.RemoteProductName == "" {
		mapping.RemoteProductName = product.UpstreamMapping.RemoteProductName
	}

	now := time.Now().Format(time.RFC3339)
	messages := make([]string, 0, 6)
	messages = append(messages, fmt.Sprintf("已保存上游映射：%s/%s", mapping.ProviderType, fallbackString(mapping.RemoteProductCode, "未指定远端商品编码")))

	switch mapping.ProviderType {
	case "ZJMF_API":
		if service.provider == nil {
			return domain.Product{}, false, fmt.Errorf("未初始化上游提供方服务")
		}
		health := service.provider.CheckFinanceHealth(mapping.ProviderAccountID)
		if !health.Enabled {
			return domain.Product{}, false, fmt.Errorf("未配置上游财务地址")
		}
		if !health.Connected {
			return domain.Product{}, false, fmt.Errorf("%s", health.Message)
		}
		if mapping.RemoteProductCode == "" {
			return domain.Product{}, false, fmt.Errorf("请先填写远端商品编码")
		}

		template, err := service.provider.GetUpstreamProductTemplate(mapping.RemoteProductCode, mapping.ProviderAccountID)
		if err != nil {
			return domain.Product{}, false, err
		}

		if strings.TrimSpace(template.GroupName) != "" {
			product.GroupName = strings.TrimSpace(template.GroupName)
		}
		if strings.TrimSpace(template.Name) != "" {
			product.Name = strings.TrimSpace(template.Name)
			mapping.RemoteProductName = template.Name
		}
		if strings.TrimSpace(template.Description) != "" {
			product.Description = strings.TrimSpace(template.Description)
		}
		if strings.TrimSpace(template.ProductType) != "" {
			product.ProductType = normalizeRemoteProductType(template.ProductType)
		}
		product.AutomationConfig = deriveAutomationConfig(product.AutomationConfig, "ZJMF_API", product.ProductType)
		product.AutomationConfig.ProviderAccountID = mapping.ProviderAccountID

		messages = append(messages, "已从上游财务拉取商品主数据")

		if mapping.AutoSyncPricing {
			product.Pricing = mapUpstreamPricing(template.Pricing)
			messages = append(messages, fmt.Sprintf("已同步 %d 条价格周期", len(product.Pricing)))
		}
		if mapping.AutoSyncConfig {
			product.ConfigOptions = mapUpstreamConfigOptions(template.ConfigOptions)
			messages = append(messages, fmt.Sprintf("已同步 %d 个配置项", len(product.ConfigOptions)))
		}
		if mapping.AutoSyncTemplate {
			product.ResourceTemplate = deriveResourceTemplate(product.ResourceTemplate, template.ConfigOptions)
			messages = append(messages, "已按远端配置推导云资源模板")
		}

	case "MOFANG_CLOUD":
		product.AutomationConfig = deriveAutomationConfig(product.AutomationConfig, "MOFANG_CLOUD", product.ProductType)
		product.AutomationConfig.ProviderAccountID = mapping.ProviderAccountID
		if service.provider != nil {
			health := service.provider.CheckHealth(mapping.ProviderAccountID)
			if health.Enabled && health.Connected {
				messages = append(messages, "魔方云连接已验证")
			} else if health.Enabled && !health.Connected {
				messages = append(messages, "魔方云未连接，已按本地模板完成映射")
			}
		}
		if mapping.AutoSyncConfig {
			product.ConfigOptions = service.buildMofangCloudConfigOptions(product.ConfigOptions)
			messages = append(messages, "已同步标准云配置项")
		}
		if mapping.AutoSyncTemplate {
			product.ResourceTemplate = ensureMofangCloudTemplate(product.ResourceTemplate)
			messages = append(messages, "已同步云资源模板")
		}
		if mapping.AutoSyncPricing {
			product.Pricing = ensureCloudPricing(product.Pricing)
			messages = append(messages, "已同步基础价格矩阵")
		}

	default:
		if mapping.AutoSyncPricing || mapping.AutoSyncConfig || mapping.AutoSyncTemplate {
			messages = append(messages, "当前上游类型暂未定义自动同步规则")
		}
	}

	mapping.SyncStatus = "SUCCESS"
	mapping.SyncMessage = strings.Join(messages, "；")
	mapping.LastSyncedAt = now
	product.UpstreamMapping = mapping

	updated, ok := service.repository.Update(id, product)
	if !ok {
		return domain.Product{}, false, fmt.Errorf("商品更新失败")
	}

	service.audit.Record(audit.Entry{
		ActorType:   "ADMIN",
		ActorID:     adminID,
		Actor:       adminName,
		Action:      "catalog.product.sync_upstream",
		TargetType:  "product",
		TargetID:    updated.ID,
		Target:      updated.ProductNo,
		RequestID:   requestID,
		Description: "同步商品上游映射、价格矩阵、配置项和云模板",
		Payload: map[string]any{
			"providerType":      updated.UpstreamMapping.ProviderType,
			"remoteProductCode": updated.UpstreamMapping.RemoteProductCode,
			"autoSyncPricing":   updated.UpstreamMapping.AutoSyncPricing,
			"autoSyncConfig":    updated.UpstreamMapping.AutoSyncConfig,
			"autoSyncTemplate":  updated.UpstreamMapping.AutoSyncTemplate,
		},
	})

	return updated, true, nil
}

func (service *Service) ImportUpstreamProducts(
	request dto.ImportUpstreamProductsRequest,
	adminID int64,
	adminName,
	requestID string,
) (dto.ImportUpstreamProductsResponse, error) {
	if service.provider == nil {
		return dto.ImportUpstreamProductsResponse{}, fmt.Errorf("未初始化上游接口服务")
	}
	if request.ProviderAccountID == 0 {
		return dto.ImportUpstreamProductsResponse{}, fmt.Errorf("请选择上游接口账户")
	}

	health := service.provider.CheckFinanceHealth(request.ProviderAccountID)
	if !health.Enabled {
		return dto.ImportUpstreamProductsResponse{}, fmt.Errorf("当前接口账户未启用上游财务")
	}
	if !health.Connected {
		return dto.ImportUpstreamProductsResponse{}, fmt.Errorf(
			"%s",
			firstNonEmpty(strings.TrimSpace(health.Message), "上游财务连接失败"),
		)
	}

	groups, err := service.provider.ListUpstreamProducts(request.ProviderAccountID)
	if err != nil {
		return dto.ImportUpstreamProductsResponse{}, err
	}

	selectedCodes := make(map[string]struct{}, len(request.RemoteProductCodes))
	for _, code := range request.RemoteProductCodes {
		trimmed := strings.TrimSpace(code)
		if trimmed == "" {
			continue
		}
		selectedCodes[trimmed] = struct{}{}
	}
	if !request.ImportAll && len(selectedCodes) == 0 {
		return dto.ImportUpstreamProductsResponse{}, fmt.Errorf("请至少选择一个上游商品")
	}

	autoSyncPricing := request.AutoSyncPricing || (!request.AutoSyncConfig && !request.AutoSyncTemplate && !request.AutoSyncPricing)
	autoSyncConfig := request.AutoSyncConfig || (!request.AutoSyncConfig && !request.AutoSyncTemplate && !request.AutoSyncPricing)
	autoSyncTemplate := request.AutoSyncTemplate || (!request.AutoSyncConfig && !request.AutoSyncTemplate && !request.AutoSyncPricing)
	now := time.Now().Format(time.RFC3339)

	products := service.repository.List()
	items := make([]dto.ImportUpstreamProductItem, 0)
	importedCount := 0
	updatedCount := 0
	failedCount := 0

	for _, group := range groups {
		for _, remoteProduct := range group.Products {
			if !request.ImportAll {
				if _, ok := selectedCodes[strings.TrimSpace(remoteProduct.RemoteProductCode)]; !ok {
					continue
				}
			}

			template, templateErr := service.provider.GetUpstreamProductTemplate(remoteProduct.RemoteProductCode, request.ProviderAccountID)
			if templateErr != nil {
				failedCount++
				items = append(items, dto.ImportUpstreamProductItem{
					RemoteProductCode: remoteProduct.RemoteProductCode,
					RemoteProductName: firstNonEmpty(template.Name, remoteProduct.Name),
					GroupName:         firstNonEmpty(template.GroupName, remoteProduct.GroupName, group.GroupName),
					Status:            "FAILED",
					Operation:         "skip",
					Message:           templateErr.Error(),
				})
				continue
			}

			existing, exists := findProductByUpstreamMapping(products, request.ProviderAccountID, "ZJMF_API", remoteProduct.RemoteProductCode)
			product := service.buildImportedProduct(existing, exists, remoteProduct, template, request.ProviderAccountID, autoSyncPricing, autoSyncConfig, autoSyncTemplate, now)

			var (
				saved   domain.Product
				ok      bool
				action  string
				message string
			)

			if exists {
				saved, ok = service.repository.Update(existing.ID, product)
				action = "update"
				message = "已按上游模板更新本地商品"
			} else {
				saved = service.repository.Create(product)
				ok = saved.ID > 0
				action = "create"
				message = "已从上游创建本地商品"
			}

			if !ok {
				failedCount++
				items = append(items, dto.ImportUpstreamProductItem{
					RemoteProductCode: remoteProduct.RemoteProductCode,
					RemoteProductName: firstNonEmpty(saved.Name, product.Name, remoteProduct.Name),
					GroupName:         firstNonEmpty(saved.GroupName, product.GroupName, group.GroupName),
					Status:            "FAILED",
					Operation:         action,
					Message:           "本地商品保存失败",
				})
				continue
			}

			if exists {
				updatedCount++
			} else {
				importedCount++
				products = append(products, saved)
			}

			service.audit.Record(audit.Entry{
				ActorType:   "ADMIN",
				ActorID:     adminID,
				Actor:       adminName,
				Action:      "catalog.product.import_upstream",
				TargetType:  "product",
				TargetID:    saved.ID,
				Target:      saved.ProductNo,
				RequestID:   requestID,
				Description: "从上游财务导入商品模板、价格矩阵、配置项和云模板",
				Payload: map[string]any{
					"providerAccountId": request.ProviderAccountID,
					"providerType":      "ZJMF_API",
					"remoteProductCode": remoteProduct.RemoteProductCode,
					"operation":         action,
				},
			})

			items = append(items, dto.ImportUpstreamProductItem{
				RemoteProductCode: remoteProduct.RemoteProductCode,
				RemoteProductName: saved.Name,
				GroupName:         saved.GroupName,
				ProductID:         saved.ID,
				ProductNo:         saved.ProductNo,
				Status:            "SUCCESS",
				Operation:         action,
				Message:           message,
			})
		}
	}

	return dto.ImportUpstreamProductsResponse{
		ProviderAccountID: request.ProviderAccountID,
		ImportedCount:     importedCount,
		UpdatedCount:      updatedCount,
		FailedCount:       failedCount,
		Items:             items,
		Total:             importedCount + updatedCount,
		Created:           importedCount,
		Updated:           updatedCount,
		Skipped:           0,
		Failed:            failedCount,
		Message:           fmt.Sprintf("上游商品导入完成，新增 %d，更新 %d，失败 %d", importedCount, updatedCount, failedCount),
	}, nil
}

func buildProductPayload(
	groupName, name, description, productType, status string,
	pricing []dto.PriceOptionInput,
	configOptions []dto.ConfigOptionInput,
	resource dto.ResourceTemplateInput,
	automation dto.AutomationConfigInput,
	upstream dto.UpstreamMappingInput,
) domain.Product {
	productStatus := domain.ProductStatusActive
	if status != "" {
		productStatus = domain.ProductStatus(status)
	}

	return domain.Product{
		GroupName:        strings.TrimSpace(groupName),
		Name:             strings.TrimSpace(name),
		Description:      strings.TrimSpace(description),
		ProductType:      strings.TrimSpace(productType),
		Status:           productStatus,
		Pricing:          normalizePricing(pricing),
		ConfigOptions:    normalizeConfigOptions(configOptions),
		ResourceTemplate: normalizeResourceTemplate(resource),
		AutomationConfig: normalizeAutomationConfig(automation, upstream.ProviderType, productType),
		UpstreamMapping:  normalizeUpstreamMapping(dto.SyncProductUpstreamRequest(upstream)),
	}
}

func findProductByUpstreamMapping(
	products []domain.Product,
	providerAccountID int64,
	providerType,
	remoteProductCode string,
) (domain.Product, bool) {
	for _, item := range products {
		if item.UpstreamMapping.ProviderAccountID != providerAccountID {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(item.UpstreamMapping.ProviderType), strings.TrimSpace(providerType)) {
			continue
		}
		if strings.TrimSpace(item.UpstreamMapping.RemoteProductCode) != strings.TrimSpace(remoteProductCode) {
			continue
		}
		return item, true
	}
	return domain.Product{}, false
}

func (service *Service) buildImportedProduct(
	existing domain.Product,
	exists bool,
	remoteProduct providerService.UpstreamCatalogProduct,
	template providerService.UpstreamProductTemplate,
	providerAccountID int64,
	autoSyncPricing,
	autoSyncConfig,
	autoSyncTemplate bool,
	syncedAt string,
) domain.Product {
	product := existing
	product.GroupName = firstNonEmpty(strings.TrimSpace(template.GroupName), strings.TrimSpace(remoteProduct.GroupName), strings.TrimSpace(existing.GroupName), "上游云产品")
	product.Name = firstNonEmpty(strings.TrimSpace(template.Name), strings.TrimSpace(remoteProduct.Name), strings.TrimSpace(existing.Name), "未命名上游商品")
	product.Description = firstNonEmpty(strings.TrimSpace(template.Description), strings.TrimSpace(remoteProduct.Description), strings.TrimSpace(existing.Description))
	product.ProductType = normalizeRemoteProductType(firstNonEmpty(template.ProductType, remoteProduct.ProductType, existing.ProductType))
	product.Status = domain.ProductStatusActive

	if !exists || autoSyncPricing {
		product.Pricing = mapUpstreamPricing(template.Pricing)
	}
	if !exists || autoSyncConfig {
		product.ConfigOptions = mapUpstreamConfigOptions(template.ConfigOptions)
	}
	if !exists || autoSyncTemplate {
		product.ResourceTemplate = deriveResourceTemplate(product.ResourceTemplate, template.ConfigOptions)
	}

	product.AutomationConfig = deriveAutomationConfig(product.AutomationConfig, "ZJMF_API", product.ProductType)
	product.AutomationConfig.ProviderAccountID = providerAccountID
	product.AutomationConfig.Channel = "ZJMF_API"
	product.AutomationConfig.AutoProvision = true
	if strings.TrimSpace(product.AutomationConfig.ProvisionStage) == "" {
		product.AutomationConfig.ProvisionStage = "AFTER_PAYMENT"
	}

	product.UpstreamMapping = domain.UpstreamMapping{
		ProviderAccountID: providerAccountID,
		ProviderType:      "ZJMF_API",
		SourceName:        firstNonEmpty(strings.TrimSpace(template.GroupName), strings.TrimSpace(remoteProduct.GroupName), "上游财务"),
		RemoteProductCode: strings.TrimSpace(remoteProduct.RemoteProductCode),
		RemoteProductName: firstNonEmpty(strings.TrimSpace(template.Name), strings.TrimSpace(remoteProduct.Name)),
		PricePolicy:       "FOLLOW_UPSTREAM",
		AutoSyncPricing:   autoSyncPricing,
		AutoSyncConfig:    autoSyncConfig,
		AutoSyncTemplate:  autoSyncTemplate,
		SyncStatus:        "SUCCESS",
		SyncMessage:       "已从上游财务导入模板",
		LastSyncedAt:      syncedAt,
	}

	return product
}

func normalizePricing(items []dto.PriceOptionInput) []domain.PriceOption {
	if len(items) == 0 {
		return []domain.PriceOption{
			{CycleCode: "monthly", CycleName: "月付", Price: 0, SetupFee: 0},
		}
	}

	result := make([]domain.PriceOption, 0, len(items))
	for _, item := range items {
		cycleCode := strings.TrimSpace(item.CycleCode)
		cycleName := strings.TrimSpace(item.CycleName)
		if cycleCode == "" {
			continue
		}
		if cycleName == "" {
			cycleName = cycleLabelByCode(cycleCode)
		}
		result = append(result, domain.PriceOption{
			CycleCode: cycleCode,
			CycleName: cycleName,
			Price:     item.Price,
			SetupFee:  item.SetupFee,
		})
	}
	if len(result) == 0 {
		return []domain.PriceOption{{CycleCode: "monthly", CycleName: "月付", Price: 0, SetupFee: 0}}
	}
	return result
}

func normalizeConfigOptions(items []dto.ConfigOptionInput) []domain.ConfigOption {
	result := make([]domain.ConfigOption, 0, len(items))
	for _, item := range items {
		code := strings.TrimSpace(item.Code)
		name := strings.TrimSpace(item.Name)
		if code == "" || name == "" {
			continue
		}

		choices := make([]domain.ConfigOptionChoice, 0, len(item.Choices))
		for _, choice := range item.Choices {
			value := strings.TrimSpace(choice.Value)
			label := strings.TrimSpace(choice.Label)
			if value == "" && label == "" {
				continue
			}
			if value == "" {
				value = label
			}
			if label == "" {
				label = value
			}
			choices = append(choices, domain.ConfigOptionChoice{
				Value:      value,
				Label:      label,
				PriceDelta: choice.PriceDelta,
			})
		}

		inputType := strings.TrimSpace(item.InputType)
		if inputType == "" {
			inputType = "select"
		}

		result = append(result, domain.ConfigOption{
			Code:         code,
			Name:         name,
			InputType:    inputType,
			Required:     item.Required,
			DefaultValue: strings.TrimSpace(item.DefaultValue),
			Description:  strings.TrimSpace(item.Description),
			Choices:      choices,
		})
	}
	return result
}

func normalizeResourceTemplate(item dto.ResourceTemplateInput) domain.ResourceTemplate {
	return domain.ResourceTemplate{
		RegionName:      strings.TrimSpace(item.RegionName),
		ZoneName:        strings.TrimSpace(item.ZoneName),
		OperatingSystem: strings.TrimSpace(item.OperatingSystem),
		LoginUsername:   strings.TrimSpace(item.LoginUsername),
		SecurityGroup:   strings.TrimSpace(item.SecurityGroup),
		CPUCores:        item.CPUCores,
		MemoryGB:        item.MemoryGB,
		SystemDiskGB:    item.SystemDiskGB,
		DataDiskGB:      item.DataDiskGB,
		BandwidthMbps:   item.BandwidthMbps,
		PublicIPCount:   item.PublicIPCount,
	}
}

func normalizeUpstreamMapping(item dto.SyncProductUpstreamRequest) domain.UpstreamMapping {
	providerType := strings.ToUpper(strings.TrimSpace(item.ProviderType))
	if providerType == "" {
		providerType = "NONE"
	}

	pricePolicy := strings.ToUpper(strings.TrimSpace(item.PricePolicy))
	if pricePolicy == "" {
		pricePolicy = "CUSTOM"
	}

	return domain.UpstreamMapping{
		ProviderAccountID: item.ProviderAccountID,
		ProviderType:      providerType,
		SourceName:        strings.TrimSpace(item.SourceName),
		RemoteProductCode: strings.TrimSpace(item.RemoteProductCode),
		RemoteProductName: strings.TrimSpace(item.RemoteProductName),
		PricePolicy:       pricePolicy,
		AutoSyncPricing:   item.AutoSyncPricing,
		AutoSyncConfig:    item.AutoSyncConfig,
		AutoSyncTemplate:  item.AutoSyncTemplate,
	}
}

func normalizeAutomationConfig(item dto.AutomationConfigInput, providerType, productType string) domain.AutomationConfig {
	channel := strings.ToUpper(strings.TrimSpace(item.Channel))
	if channel == "" {
		channel = strings.ToUpper(strings.TrimSpace(providerType))
	}
	switch channel {
	case "MOFANG_CLOUD", "ZJMF_API", "MANUAL", "RESOURCE":
	default:
		channel = "LOCAL"
	}

	moduleType := strings.ToUpper(strings.TrimSpace(item.ModuleType))
	if moduleType == "" {
		moduleType = defaultModuleType(productType, channel)
	}

	provisionStage := strings.ToUpper(strings.TrimSpace(item.ProvisionStage))
	if provisionStage == "" {
		provisionStage = "AFTER_PAYMENT"
	}

	autoProvision := item.AutoProvision
	if !item.AutoProvision && item.Channel == "" && item.ProvisionStage == "" {
		autoProvision = provisionStage == "AFTER_PAYMENT"
	}

	return domain.AutomationConfig{
		ProviderAccountID: item.ProviderAccountID,
		Channel:           channel,
		ModuleType:        moduleType,
		ProvisionStage:    provisionStage,
		AutoProvision:     autoProvision,
		ServerGroup:       strings.TrimSpace(item.ServerGroup),
		ProviderNode:      strings.TrimSpace(item.ProviderNode),
	}
}

func deriveAutomationConfig(current domain.AutomationConfig, providerType, productType string) domain.AutomationConfig {
	return normalizeAutomationConfig(dto.AutomationConfigInput{
		ProviderAccountID: current.ProviderAccountID,
		Channel:           current.Channel,
		ModuleType:        current.ModuleType,
		ProvisionStage:    current.ProvisionStage,
		AutoProvision:     current.AutoProvision,
		ServerGroup:       current.ServerGroup,
		ProviderNode:      current.ProviderNode,
	}, providerType, productType)
}

func isEmptyAutomationConfigInput(item dto.AutomationConfigInput) bool {
	return item.ProviderAccountID == 0 &&
		strings.TrimSpace(item.Channel) == "" &&
		strings.TrimSpace(item.ModuleType) == "" &&
		strings.TrimSpace(item.ProvisionStage) == "" &&
		!item.AutoProvision &&
		strings.TrimSpace(item.ServerGroup) == "" &&
		strings.TrimSpace(item.ProviderNode) == ""
}

func mergeProductUpstreamMapping(existing domain.UpstreamMapping, input dto.UpstreamMappingInput) domain.UpstreamMapping {
	if input.ProviderAccountID == 0 &&
		strings.TrimSpace(input.ProviderType) == "" &&
		strings.TrimSpace(input.SourceName) == "" &&
		strings.TrimSpace(input.RemoteProductCode) == "" &&
		strings.TrimSpace(input.RemoteProductName) == "" &&
		strings.TrimSpace(input.PricePolicy) == "" &&
		!input.AutoSyncPricing &&
		!input.AutoSyncConfig &&
		!input.AutoSyncTemplate {
		return existing
	}

	merged := existing
	if input.ProviderAccountID != 0 {
		merged.ProviderAccountID = input.ProviderAccountID
	}
	if value := strings.ToUpper(strings.TrimSpace(input.ProviderType)); value != "" {
		merged.ProviderType = value
	}
	if value := strings.TrimSpace(input.SourceName); value != "" {
		merged.SourceName = value
	}
	if value := strings.TrimSpace(input.RemoteProductCode); value != "" {
		merged.RemoteProductCode = value
	}
	if value := strings.TrimSpace(input.RemoteProductName); value != "" {
		merged.RemoteProductName = value
	}
	if value := strings.ToUpper(strings.TrimSpace(input.PricePolicy)); value != "" {
		merged.PricePolicy = value
	}
	merged.AutoSyncPricing = input.AutoSyncPricing
	merged.AutoSyncConfig = input.AutoSyncConfig
	merged.AutoSyncTemplate = input.AutoSyncTemplate
	return merged
}

func defaultModuleType(productType, channel string) string {
	switch {
	case channel == "MOFANG_CLOUD":
		return "DCIMCLOUD"
	case strings.EqualFold(productType, "COLOCATION"):
		return "DCIM"
	case strings.EqualFold(productType, "CLOUD"):
		return "CLOUD"
	default:
		return "NORMAL"
	}
}

func defaultSourceName(providerType string) string {
	switch providerType {
	case "MOFANG_CLOUD":
		return "魔方云"
	case "ZJMF_API":
		return "上下游财务"
	case "MANUAL":
		return "手动资源池"
	default:
		return "本地映射"
	}
}

func (service *Service) buildMofangCloudConfigOptions(existing []domain.ConfigOption) []domain.ConfigOption {
	areaChoices := make([]domain.ConfigOptionChoice, 0)
	imageChoices := make([]domain.ConfigOptionChoice, 0)

	if service.provider != nil {
		if areas, err := service.provider.ListMofangAreas(); err == nil {
			for _, item := range areas {
				areaChoices = append(areaChoices, domain.ConfigOptionChoice{
					Value:      item.Value,
					Label:      item.Label,
					PriceDelta: 0,
				})
			}
		}
		if images, err := service.provider.ListMofangImages(); err == nil {
			for _, item := range images {
				imageChoices = append(imageChoices, domain.ConfigOptionChoice{
					Value:      item.Value,
					Label:      item.Label,
					PriceDelta: 0,
				})
			}
		}
	}

	options := ensureMofangCloudConfigOptions(existing, areaChoices, imageChoices)
	for index := range options {
		switch options[index].Code {
		case "area":
			if len(areaChoices) > 0 {
				defaultValue := preserveMofangDefaultValue(options[index].DefaultValue, areaChoices)
				options[index].Choices = areaChoices
				options[index].DefaultValue = defaultValue
			}
		case "os":
			if len(imageChoices) > 0 {
				defaultValue := preserveMofangDefaultValue(options[index].DefaultValue, imageChoices)
				options[index].Choices = imageChoices
				options[index].DefaultValue = defaultValue
			}
		}
	}
	return options
}

func fallbackString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func ensureCloudPricing(pricing []domain.PriceOption) []domain.PriceOption {
	if len(pricing) != 0 {
		return pricing
	}
	return []domain.PriceOption{
		{CycleCode: "monthly", CycleName: "月付", Price: 199, SetupFee: 0},
		{CycleCode: "quarterly", CycleName: "季付", Price: 569, SetupFee: 0},
		{CycleCode: "annual", CycleName: "年付", Price: 1999, SetupFee: 0},
	}
}

func ensureMofangCloudTemplate(template domain.ResourceTemplate) domain.ResourceTemplate {
	if strings.TrimSpace(template.RegionName) == "" {
		template.RegionName = "华南广州"
	}
	if strings.TrimSpace(template.ZoneName) == "" {
		template.ZoneName = "广州三区"
	}
	if strings.TrimSpace(template.OperatingSystem) == "" {
		template.OperatingSystem = "Rocky Linux 9"
	}
	if strings.TrimSpace(template.LoginUsername) == "" {
		template.LoginUsername = "root"
	}
	if strings.TrimSpace(template.SecurityGroup) == "" {
		template.SecurityGroup = "default-cloud"
	}
	if template.CPUCores == 0 {
		template.CPUCores = 4
	}
	if template.MemoryGB == 0 {
		template.MemoryGB = 8
	}
	if template.SystemDiskGB == 0 {
		template.SystemDiskGB = 60
	}
	if template.DataDiskGB == 0 {
		template.DataDiskGB = 120
	}
	if template.BandwidthMbps == 0 {
		template.BandwidthMbps = 20
	}
	if template.PublicIPCount == 0 {
		template.PublicIPCount = 1
	}
	return template
}

func ensureMofangCloudConfigOptions(options []domain.ConfigOption, areaChoices, imageChoices []domain.ConfigOptionChoice) []domain.ConfigOption {
	if len(areaChoices) == 0 {
		areaChoices = []domain.ConfigOptionChoice{
			{Value: "gz-3", Label: "广州三区", PriceDelta: 0},
			{Value: "sh-1", Label: "上海一区", PriceDelta: 0},
		}
	}
	if len(imageChoices) == 0 {
		imageChoices = []domain.ConfigOptionChoice{
			{Value: "rocky-9", Label: "Rocky Linux 9", PriceDelta: 0},
			{Value: "ubuntu-24.04", Label: "Ubuntu 24.04 LTS", PriceDelta: 0},
			{Value: "debian-12", Label: "Debian 12", PriceDelta: 0},
		}
	}
	standard := []domain.ConfigOption{
		{
			Code:         "area",
			Name:         "地域",
			InputType:    "select",
			Required:     true,
			DefaultValue: areaChoices[0].Value,
			Description:  "选择实例所属地域和可用区",
			Choices:      cloneConfigChoices(areaChoices),
		},
		{
			Code:         "os",
			Name:         "操作系统",
			InputType:    "select",
			Required:     true,
			DefaultValue: imageChoices[0].Value,
			Description:  "同步默认云镜像模板",
			Choices:      cloneConfigChoices(imageChoices),
		},
		{
			Code:         "cpu",
			Name:         "CPU",
			InputType:    "select",
			Required:     true,
			DefaultValue: "4",
			Description:  "实例 CPU 核数",
			Choices: []domain.ConfigOptionChoice{
				{Value: "2", Label: "2 核", PriceDelta: 0},
				{Value: "4", Label: "4 核", PriceDelta: 80},
				{Value: "8", Label: "8 核", PriceDelta: 220},
			},
		},
		{
			Code:         "memory",
			Name:         "内存",
			InputType:    "select",
			Required:     true,
			DefaultValue: "8",
			Description:  "实例内存容量",
			Choices: []domain.ConfigOptionChoice{
				{Value: "4", Label: "4 GB", PriceDelta: 0},
				{Value: "8", Label: "8 GB", PriceDelta: 60},
				{Value: "16", Label: "16 GB", PriceDelta: 180},
			},
		},
		{
			Code:         "system_disk_size",
			Name:         "系统盘",
			InputType:    "select",
			Required:     true,
			DefaultValue: "60",
			Description:  "默认系统盘容量",
			Choices: []domain.ConfigOptionChoice{
				{Value: "60", Label: "60 GB", PriceDelta: 0},
				{Value: "100", Label: "100 GB", PriceDelta: 40},
			},
		},
		{
			Code:         "data_disk_size",
			Name:         "数据盘",
			InputType:    "select",
			Required:     false,
			DefaultValue: "120",
			Description:  "默认数据盘容量",
			Choices: []domain.ConfigOptionChoice{
				{Value: "120", Label: "120 GB", PriceDelta: 0},
				{Value: "240", Label: "240 GB", PriceDelta: 90},
			},
		},
		{
			Code:         "bw",
			Name:         "带宽",
			InputType:    "select",
			Required:     true,
			DefaultValue: "20",
			Description:  "公网带宽模板",
			Choices: []domain.ConfigOptionChoice{
				{Value: "20", Label: "20 Mbps", PriceDelta: 0},
				{Value: "50", Label: "50 Mbps", PriceDelta: 120},
			},
		},
		{
			Code:         "ip_num",
			Name:         "IP 数量",
			InputType:    "select",
			Required:     true,
			DefaultValue: "1",
			Description:  "公网 IPv4 数量",
			Choices: []domain.ConfigOptionChoice{
				{Value: "1", Label: "1 个", PriceDelta: 0},
				{Value: "2", Label: "2 个", PriceDelta: 30},
			},
		},
		{
			Code:         "network_type",
			Name:         "网络类型",
			InputType:    "select",
			Required:     true,
			DefaultValue: "normal",
			Description:  "实例网络类型",
			Choices: []domain.ConfigOptionChoice{
				{Value: "normal", Label: "经典网络", PriceDelta: 0},
				{Value: "vpc", Label: "VPC 专有网络", PriceDelta: 0},
			},
		},
		{
			Code:         "backup_num",
			Name:         "备份份数",
			InputType:    "select",
			Required:     false,
			DefaultValue: "0",
			Description:  "自动备份保留份数",
			Choices: []domain.ConfigOptionChoice{
				{Value: "0", Label: "不启用", PriceDelta: 0},
				{Value: "1", Label: "1 份", PriceDelta: 15},
				{Value: "3", Label: "3 份", PriceDelta: 45},
			},
		},
		{
			Code:         "snap_num",
			Name:         "快照份数",
			InputType:    "select",
			Required:     false,
			DefaultValue: "0",
			Description:  "快照保留份数",
			Choices: []domain.ConfigOptionChoice{
				{Value: "0", Label: "不启用", PriceDelta: 0},
				{Value: "1", Label: "1 份", PriceDelta: 10},
				{Value: "3", Label: "3 份", PriceDelta: 30},
			},
		},
	}

	merged := make([]domain.ConfigOption, 0, len(standard)+len(options))
	used := make(map[string]struct{}, len(standard))
	currentByCode := make(map[string]domain.ConfigOption, len(options))
	for _, item := range options {
		code := strings.TrimSpace(item.Code)
		if code == "" {
			continue
		}
		currentByCode[code] = item
	}

	for _, item := range standard {
		if current, ok := currentByCode[item.Code]; ok {
			item = mergeMofangConfigOption(item, current)
		}
		merged = append(merged, item)
		used[item.Code] = struct{}{}
	}

	for _, item := range options {
		code := strings.TrimSpace(item.Code)
		if code == "" {
			continue
		}
		if _, ok := used[code]; ok {
			continue
		}
		if _, skip := legacyMofangAlias(code); skip {
			continue
		}
		merged = append(merged, item)
	}

	return merged
}

func mergeMofangConfigOption(base, current domain.ConfigOption) domain.ConfigOption {
	if strings.TrimSpace(current.DefaultValue) != "" && optionChoiceExists(current.DefaultValue, base.Choices) {
		base.DefaultValue = current.DefaultValue
	}
	if strings.TrimSpace(current.Description) != "" {
		base.Description = current.Description
	}
	if len(base.Choices) != 0 && len(current.Choices) != 0 {
		currentChoices := make(map[string]domain.ConfigOptionChoice, len(current.Choices))
		for _, choice := range current.Choices {
			currentChoices[strings.TrimSpace(choice.Value)] = choice
		}
		for index := range base.Choices {
			if currentChoice, ok := currentChoices[strings.TrimSpace(base.Choices[index].Value)]; ok {
				base.Choices[index].PriceDelta = currentChoice.PriceDelta
				if strings.TrimSpace(currentChoice.Label) != "" {
					base.Choices[index].Label = currentChoice.Label
				}
			}
		}
	}
	return base
}

func optionChoiceExists(value string, choices []domain.ConfigOptionChoice) bool {
	for _, choice := range choices {
		if strings.TrimSpace(choice.Value) == strings.TrimSpace(value) {
			return true
		}
	}
	return false
}

func cloneConfigChoices(items []domain.ConfigOptionChoice) []domain.ConfigOptionChoice {
	result := make([]domain.ConfigOptionChoice, len(items))
	copy(result, items)
	return result
}

func preserveMofangDefaultValue(current string, choices []domain.ConfigOptionChoice) string {
	current = strings.TrimSpace(current)
	if current != "" {
		for _, choice := range choices {
			if strings.TrimSpace(choice.Value) == current {
				return current
			}
		}
	}
	if len(choices) == 0 {
		return current
	}
	return choices[0].Value
}

func legacyMofangAlias(code string) (string, bool) {
	switch strings.TrimSpace(code) {
	case "backup":
		return "backup_num", true
	case "snapshot":
		return "snap_num", true
	default:
		return "", false
	}
}

func normalizeRemoteProductType(input string) string {
	switch strings.ToUpper(strings.TrimSpace(input)) {
	case "DCIMCLOUD", "CLOUD":
		return "CLOUD"
	case "BANDWIDTH", "CDN":
		return "BANDWIDTH"
	case "DCIM", "SERVER", "BAREMETAL", "BAREMETALSERVER":
		return "COLOCATION"
	default:
		return "CLOUD"
	}
}

func mapUpstreamPricing(items []providerService.UpstreamCyclePrice) []domain.PriceOption {
	if len(items) == 0 {
		return []domain.PriceOption{{CycleCode: "monthly", CycleName: "月付", Price: 0, SetupFee: 0}}
	}
	result := make([]domain.PriceOption, 0, len(items))
	for _, item := range items {
		cycleCode := strings.TrimSpace(item.CycleCode)
		if cycleCode == "" {
			continue
		}
		cycleName := strings.TrimSpace(item.CycleName)
		if cycleName == "" {
			cycleName = cycleLabelByCode(cycleCode)
		}
		result = append(result, domain.PriceOption{
			CycleCode: cycleCode,
			CycleName: cycleName,
			Price:     item.Price,
			SetupFee:  item.SetupFee,
		})
	}
	if len(result) == 0 {
		return []domain.PriceOption{{CycleCode: "monthly", CycleName: "月付", Price: 0, SetupFee: 0}}
	}
	return result
}

func mapUpstreamConfigOptions(items []providerService.UpstreamConfigOption) []domain.ConfigOption {
	result := make([]domain.ConfigOption, 0, len(items))
	for _, item := range items {
		code := strings.TrimSpace(item.Code)
		name := strings.TrimSpace(item.Name)
		if code == "" || name == "" {
			continue
		}
		choices := make([]domain.ConfigOptionChoice, 0, len(item.Choices))
		for _, choice := range item.Choices {
			value := strings.TrimSpace(choice.Value)
			label := strings.TrimSpace(choice.Label)
			if value == "" && label == "" {
				continue
			}
			if value == "" {
				value = label
			}
			if label == "" {
				label = value
			}
			choices = append(choices, domain.ConfigOptionChoice{
				Value:      value,
				Label:      label,
				PriceDelta: choice.PriceDelta,
			})
		}
		defaultValue := strings.TrimSpace(item.DefaultValue)
		if defaultValue == "" && len(choices) > 0 {
			defaultValue = choices[0].Value
		}
		result = append(result, domain.ConfigOption{
			Code:         code,
			Name:         name,
			InputType:    fallbackString(strings.TrimSpace(item.InputType), "select"),
			Required:     item.Required,
			DefaultValue: defaultValue,
			Description:  strings.TrimSpace(item.Description),
			Choices:      choices,
		})
	}
	return result
}

func deriveResourceTemplate(base domain.ResourceTemplate, options []providerService.UpstreamConfigOption) domain.ResourceTemplate {
	template := base
	optionMap := make(map[string]providerService.UpstreamConfigOption, len(options))
	for _, item := range options {
		optionMap[item.Code] = item
	}

	if template.RegionName == "" {
		template.RegionName = preferredOptionLabel(optionMap["area"], optionMap["node"])
	}
	if template.ZoneName == "" {
		template.ZoneName = template.RegionName
	}
	if template.OperatingSystem == "" {
		template.OperatingSystem = preferredOptionLabel(optionMap["os"])
	}
	if template.LoginUsername == "" {
		template.LoginUsername = "root"
	}
	if template.SecurityGroup == "" {
		template.SecurityGroup = "default"
	}
	if template.CPUCores == 0 {
		template.CPUCores = preferredOptionNumber(optionMap["cpu"], 2)
	}
	if template.MemoryGB == 0 {
		template.MemoryGB = preferredOptionNumber(optionMap["memory"], 4)
	}
	if template.SystemDiskGB == 0 {
		template.SystemDiskGB = preferredOptionNumber(optionMap["system_disk_size"], 40)
	}
	if template.DataDiskGB == 0 {
		template.DataDiskGB = preferredOptionNumber(optionMap["data_disk_size"], 0)
	}
	if template.BandwidthMbps == 0 {
		template.BandwidthMbps = preferredOptionNumber(optionMap["bw"], 20)
	}
	if template.PublicIPCount == 0 {
		template.PublicIPCount = preferredOptionNumber(optionMap["ip_num"], 1)
	}

	return template
}

func preferredOptionLabel(items ...providerService.UpstreamConfigOption) string {
	for _, item := range items {
		if strings.TrimSpace(item.DefaultValue) != "" {
			for _, choice := range item.Choices {
				if choice.Value == item.DefaultValue && strings.TrimSpace(choice.Label) != "" {
					return choice.Label
				}
			}
		}
		if len(item.Choices) > 0 && strings.TrimSpace(item.Choices[0].Label) != "" {
			return item.Choices[0].Label
		}
	}
	return ""
}

func preferredOptionNumber(item providerService.UpstreamConfigOption, fallback int) int {
	if strings.TrimSpace(item.DefaultValue) != "" {
		if parsed := extractLeadingInt(item.DefaultValue); parsed > 0 {
			return parsed
		}
		for _, choice := range item.Choices {
			if choice.Value == item.DefaultValue {
				if parsed := extractLeadingInt(choice.Label); parsed > 0 {
					return parsed
				}
			}
		}
	}
	for _, choice := range item.Choices {
		if parsed := extractLeadingInt(choice.Label); parsed > 0 {
			return parsed
		}
		if parsed := extractLeadingInt(choice.Value); parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func extractLeadingInt(input string) int {
	builder := strings.Builder{}
	for _, r := range input {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
			continue
		}
		if builder.Len() > 0 {
			break
		}
	}
	if builder.Len() == 0 {
		return 0
	}
	value, err := strconv.Atoi(builder.String())
	if err != nil {
		return 0
	}
	return value
}

func cycleLabelByCode(cycleCode string) string {
	switch cycleCode {
	case "monthly":
		return "月付"
	case "quarterly":
		return "季付"
	case "semiannually", "semiannual":
		return "半年付"
	case "annually", "annual":
		return "年付"
	case "biennially":
		return "两年付"
	case "triennially":
		return "三年付"
	case "onetime":
		return "一次性"
	default:
		return cycleCode
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
