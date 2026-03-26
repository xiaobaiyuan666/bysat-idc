<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import { downloadCsv } from "@/utils/download";
import {
  createProduct,
  fetchProductGroups,
  fetchProducts,
  fetchProviderAccounts,
  fetchUpstreamProductTemplate,
  fetchUpstreamProducts,
  importUpstreamProducts,
  syncProductUpstream,
  updateProduct,
  type Product,
  type ProviderAccount,
  type UpstreamCatalogGroup,
  type UpstreamProductTemplate
} from "@/api/admin";

type TabKey = "ALL" | "ACTIVE" | "INACTIVE" | "LOCAL" | "MOFANG_CLOUD" | "ZJMF_API";

const router = useRouter();
const route = useRoute();

const loading = ref(false);
const syncing = ref(false);
const statusUpdating = ref(false);
const creating = ref(false);
const importLoading = ref(false);
const importSubmitting = ref(false);
const createVisible = ref(false);
const importVisible = ref(false);
const advancedVisible = ref(false);

const products = ref<Product[]>([]);
const groups = ref<string[]>([]);
const providerAccounts = ref<ProviderAccount[]>([]);
const upstreamGroups = ref<UpstreamCatalogGroup[]>([]);
const selectedRows = ref<Product[]>([]);
const activeTab = ref<TabKey>("ALL");

const createForm = reactive({
  groupName: "云主机",
  name: "",
  description: "",
  productType: "CLOUD",
  monthlyPrice: 0,
  annualPrice: 0,
  status: "ACTIVE"
});

const importForm = reactive({
  providerAccountId: 0,
  importAll: true,
  remoteProductCodes: [] as string[],
  autoSyncPricing: true,
  autoSyncConfig: true,
  autoSyncTemplate: true
});

const filters = reactive({
  keyword: "",
  groupName: "",
  status: "",
  productType: "",
  channel: "",
  providerAccountId: ""
});

const financeAccounts = computed(() =>
  providerAccounts.value.filter(item => item.providerType === "ZJMF_API" && item.status === "ACTIVE")
);

const statusTabs = computed(() => [
  { key: "ALL", label: "全部", count: products.value.length },
  { key: "ACTIVE", label: "启用", count: products.value.filter(item => item.status === "ACTIVE").length },
  { key: "INACTIVE", label: "停用", count: products.value.filter(item => item.status === "INACTIVE").length },
  { key: "LOCAL", label: "本地", count: products.value.filter(item => item.automationConfig.channel === "LOCAL").length },
  {
    key: "MOFANG_CLOUD",
    label: "魔方云",
    count: products.value.filter(item => item.automationConfig.channel === "MOFANG_CLOUD").length
  },
  {
    key: "ZJMF_API",
    label: "上下游",
    count: products.value.filter(item => item.automationConfig.channel === "ZJMF_API").length
  }
]);

const summary = computed(() => ({
  total: products.value.length,
  active: products.value.filter(item => item.status === "ACTIVE").length,
  autoProvision: products.value.filter(item => item.automationConfig.autoProvision).length,
  upstreamMapped: products.value.filter(item => item.upstreamMapping.providerType !== "NONE").length,
  selectedSyncable: selectedRows.value.filter(canSyncProduct).length
}));

const upstreamCatalogItems = computed(() =>
  upstreamGroups.value.flatMap(group =>
    group.products.map(item => ({
      ...item,
      groupName: group.groupName,
      label: `${group.groupName} / ${item.name}`
    }))
  )
);

const importSummary = computed(() => ({
  total: upstreamCatalogItems.value.length,
  selected: importForm.importAll ? upstreamCatalogItems.value.length : importForm.remoteProductCodes.length
}));
const selectedSingleProduct = computed(() => (selectedRows.value.length === 1 ? selectedRows.value[0] : null));

const importPreview = computed(() => {
  const items = importForm.importAll
    ? upstreamCatalogItems.value
    : upstreamCatalogItems.value.filter(item => importForm.remoteProductCodes.includes(item.remoteProductCode));
  return items.slice(0, 8);
});

const filteredProducts = computed(() =>
  products.value.filter(item => {
    const keyword = filters.keyword.trim().toLowerCase();
    const providerAccountId = Number(filters.providerAccountId);
    const keywordMatched =
      !keyword ||
      item.name.toLowerCase().includes(keyword) ||
      item.productNo.toLowerCase().includes(keyword) ||
      item.description.toLowerCase().includes(keyword) ||
      item.upstreamMapping.remoteProductCode.toLowerCase().includes(keyword);
    const groupMatched = !filters.groupName || item.groupName === filters.groupName;
    const statusMatched = !filters.status || item.status === filters.status;
    const typeMatched = !filters.productType || item.productType === filters.productType;
    const channelMatched = !filters.channel || item.automationConfig.channel === filters.channel;
    const providerAccountMatched =
      !filters.providerAccountId ||
      item.automationConfig.providerAccountId === providerAccountId ||
      item.upstreamMapping.providerAccountId === providerAccountId;
    const tabMatched =
      activeTab.value === "ALL" ||
      item.status === activeTab.value ||
      item.automationConfig.channel === activeTab.value;
    return (
      keywordMatched &&
      groupMatched &&
      statusMatched &&
      typeMatched &&
      channelMatched &&
      providerAccountMatched &&
      tabMatched
    );
  })
);

function safeAccountName(account?: ProviderAccount) {
  if (!account) return "未绑定";
  const raw = `${account.name || ""}${account.sourceName || ""}`;
  if (/[?？]{2,}/.test(raw)) {
    try {
      return new URL(account.baseUrl).hostname;
    } catch {
      return `接口账户 #${account.id}`;
    }
  }
  return account.name || account.sourceName || `接口账户 #${account.id}`;
}

function accountLabel(id?: number) {
  if (!id) return "未绑定";
  return safeAccountName(providerAccounts.value.find(item => item.id === id));
}

function formatProductType(type: string) {
  const mapping: Record<string, string> = {
    CLOUD: "云产品",
    CLOUD_HOST: "云主机",
    BANDWIDTH: "带宽",
    COLOCATION: "托管",
    CDN: "CDN",
    SSL: "SSL 证书",
    SERVER: "物理服务器",
    DCIM: "DCIM 服务器",
    DCIMCLOUD: "弹性云"
  };
  return mapping[type] ?? type;
}

function formatStatus(status: string) {
  return status === "ACTIVE" ? "启用" : "停用";
}

function formatChannel(channel: string) {
  const mapping: Record<string, string> = {
    LOCAL: "本地",
    MOFANG_CLOUD: "魔方云",
    ZJMF_API: "上下游",
    RESOURCE: "资源池",
    MANUAL: "手动资源"
  };
  return mapping[channel] ?? channel;
}

function formatAutomation(product: Product) {
  return `${formatChannel(product.automationConfig.channel)} / ${accountLabel(product.automationConfig.providerAccountId)}`;
}

function formatUpstream(product: Product) {
  if (product.upstreamMapping.providerType === "NONE") return "未配置";
  return `${accountLabel(product.upstreamMapping.providerAccountId)} / ${product.upstreamMapping.remoteProductCode || "-"}`;
}

function formatSyncStatus(product: Product) {
  return product.upstreamMapping.syncStatus || "未同步";
}

function formatPricing(product: Product) {
  return product.pricing.map(item => `${item.cycleName} ¥${item.price.toFixed(2)}`).join(" / ");
}

function statusTagType(status: string) {
  return status === "ACTIVE" ? "success" : "info";
}

function canSyncProduct(product: Product) {
  if (!product.upstreamMapping.providerAccountId) return false;
  if (product.upstreamMapping.providerType === "ZJMF_API") {
    return Boolean(product.upstreamMapping.remoteProductCode);
  }
  return product.upstreamMapping.providerType !== "NONE";
}

function buildProductPayload(product: Product, overrides: Record<string, unknown> = {}) {
  return {
    groupName: product.groupName,
    name: product.name,
    description: product.description,
    productType: product.productType,
    status: product.status,
    pricing: product.pricing.map(item => ({ ...item })),
    configOptions: product.configOptions.map(item => ({
      ...item,
      choices: item.choices.map(choice => ({ ...choice }))
    })),
    resourceTemplate: { ...product.resourceTemplate },
    automationConfig: { ...product.automationConfig },
    upstreamMapping: { ...product.upstreamMapping },
    ...overrides
  };
}

function resetCreateForm() {
  createForm.groupName = "云主机";
  createForm.name = "";
  createForm.description = "";
  createForm.productType = "CLOUD";
  createForm.monthlyPrice = 0;
  createForm.annualPrice = 0;
  createForm.status = "ACTIVE";
}

function resetImportForm() {
  importForm.importAll = true;
  importForm.remoteProductCodes = [];
  importForm.autoSyncPricing = true;
  importForm.autoSyncConfig = true;
  importForm.autoSyncTemplate = true;
  ensureImportAccount();
}

function resetFilters() {
  activeTab.value = "ALL";
  filters.keyword = "";
  filters.groupName = "";
  filters.status = "";
  filters.productType = "";
  filters.channel = "";
  filters.providerAccountId = "";
}

function readRouteQueryValue(value: unknown) {
  if (Array.isArray(value)) return String(value[0] ?? "");
  if (value === undefined || value === null) return "";
  return String(value);
}

function syncFiltersFromRoute() {
  const providerAccountId = readRouteQueryValue(route.query.providerAccountId);
  const providerType = readRouteQueryValue(route.query.providerType).toUpperCase();
  const keywordFromRoute = readRouteQueryValue(route.query.keyword);

  filters.providerAccountId = providerAccountId;
  filters.keyword = keywordFromRoute;
  filters.channel = ["LOCAL", "MOFANG_CLOUD", "ZJMF_API", "RESOURCE", "MANUAL"].includes(providerType)
    ? providerType
    : "";
  advancedVisible.value = Boolean(providerAccountId || filters.channel);
}

function ensureImportAccount() {
  if (!importForm.providerAccountId && financeAccounts.value.length > 0) {
    importForm.providerAccountId = financeAccounts.value[0].id;
  }
}

function handleSelectionChange(rows: Product[]) {
  selectedRows.value = rows;
}

async function loadProducts() {
  loading.value = true;
  try {
    const [productData, groupData, accountData] = await Promise.all([
      fetchProducts(),
      fetchProductGroups(),
      fetchProviderAccounts()
    ]);
    products.value = productData.items;
    groups.value = groupData;
    providerAccounts.value = accountData;
    ensureImportAccount();
  } finally {
    loading.value = false;
  }
}

async function loadUpstreamCatalog() {
  ensureImportAccount();
  if (!importForm.providerAccountId) {
    upstreamGroups.value = [];
    return;
  }
  importLoading.value = true;
  try {
    const result = await fetchUpstreamProducts(importForm.providerAccountId);
    upstreamGroups.value = result.groups;
    if (importForm.importAll) {
      importForm.remoteProductCodes = [];
    }
  } finally {
    importLoading.value = false;
  }
}

function openImportDialog() {
  ensureImportAccount();
  importVisible.value = true;
  void loadUpstreamCatalog();
}

async function handleCreateProduct() {
  creating.value = true;
  try {
    await createProduct({
      groupName: createForm.groupName,
      name: createForm.name,
      description: createForm.description,
      productType: createForm.productType,
      status: createForm.status,
      pricing: [
        { cycleCode: "monthly", cycleName: "月付", price: createForm.monthlyPrice, setupFee: 0 },
        { cycleCode: "annual", cycleName: "年付", price: createForm.annualPrice, setupFee: 0 }
      ]
    });
    ElMessage.success("商品已创建");
    createVisible.value = false;
    resetCreateForm();
    await loadProducts();
  } finally {
    creating.value = false;
  }
}

async function syncSingleProduct(product: Product) {
  if (!canSyncProduct(product)) {
    ElMessage.warning(`商品 ${product.productNo} 还没有可用的上游映射`);
    return;
  }
  syncing.value = true;
  try {
    await syncProductUpstream(product.id, {
      providerAccountId: product.upstreamMapping.providerAccountId,
      providerType: product.upstreamMapping.providerType,
      sourceName: product.upstreamMapping.sourceName,
      remoteProductCode: product.upstreamMapping.remoteProductCode,
      remoteProductName: product.upstreamMapping.remoteProductName,
      pricePolicy: product.upstreamMapping.pricePolicy,
      autoSyncPricing: product.upstreamMapping.autoSyncPricing,
      autoSyncConfig: product.upstreamMapping.autoSyncConfig,
      autoSyncTemplate: product.upstreamMapping.autoSyncTemplate
    });
    ElMessage.success(`商品 ${product.productNo} 已同步`);
    await loadProducts();
  } finally {
    syncing.value = false;
  }
}

async function batchSyncSelected() {
  const syncable = selectedRows.value.filter(canSyncProduct);
  if (syncable.length === 0) {
    ElMessage.warning("请至少选择一个可同步的商品");
    return;
  }

  syncing.value = true;
  let success = 0;
  let failed = 0;
  try {
    for (const item of syncable) {
      try {
        await syncProductUpstream(item.id, {
          providerAccountId: item.upstreamMapping.providerAccountId,
          providerType: item.upstreamMapping.providerType,
          sourceName: item.upstreamMapping.sourceName,
          remoteProductCode: item.upstreamMapping.remoteProductCode,
          remoteProductName: item.upstreamMapping.remoteProductName,
          pricePolicy: item.upstreamMapping.pricePolicy,
          autoSyncPricing: item.upstreamMapping.autoSyncPricing,
          autoSyncConfig: item.upstreamMapping.autoSyncConfig,
          autoSyncTemplate: item.upstreamMapping.autoSyncTemplate
        });
        success += 1;
      } catch {
        failed += 1;
      }
    }
    ElMessage.success(`批量同步完成，成功 ${success}，失败 ${failed}`);
    await loadProducts();
  } finally {
    syncing.value = false;
  }
}

async function batchUpdateStatus(nextStatus: "ACTIVE" | "INACTIVE", rows?: Product[]) {
  const targetRows = rows && rows.length > 0 ? rows : selectedRows.value;
  if (targetRows.length === 0) {
    ElMessage.info("请先选择商品");
    return;
  }

  statusUpdating.value = true;
  let success = 0;
  let failed = 0;
  try {
    for (const item of targetRows) {
      if (item.status === nextStatus) continue;
      try {
        await updateProduct(
          item.id,
          buildProductPayload(item, {
            status: nextStatus
          })
        );
        success += 1;
      } catch {
        failed += 1;
      }
    }
    ElMessage.success(
      failed > 0 ? `批量状态更新完成，成功 ${success}，失败 ${failed}` : `已更新 ${success} 个商品状态`
    );
    await loadProducts();
  } finally {
    statusUpdating.value = false;
  }
}

async function duplicateSelectedProduct() {
  const source = selectedSingleProduct.value;
  if (!source) {
    ElMessage.info("请只选择一个商品进行复制");
    return;
  }

  creating.value = true;
  try {
    await createProduct(
      buildProductPayload(source, {
        name: `${source.name} - 副本`,
        description: source.description
          ? `${source.description}\n\n复制自 ${source.productNo}`
          : `复制自 ${source.productNo}`,
        upstreamMapping: {
          providerAccountId: 0,
          providerType: "NONE",
          sourceName: "",
          remoteProductCode: "",
          remoteProductName: "",
          pricePolicy: source.upstreamMapping.pricePolicy,
          autoSyncPricing: false,
          autoSyncConfig: false,
          autoSyncTemplate: false
        }
      })
    );
    ElMessage.success(`已复制商品 ${source.productNo}`);
    await loadProducts();
  } finally {
    creating.value = false;
  }
}

async function copySelectedProductNos() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择商品");
    return;
  }
  await navigator.clipboard.writeText(selectedRows.value.map(item => item.productNo).join("\n"));
  ElMessage.success(`已复制 ${selectedRows.value.length} 个商品编号`);
}

function exportProducts(rows: Product[], filename: string) {
  downloadCsv(
    filename,
    ["商品编号", "商品名称", "分组", "类型", "自动化渠道", "接口账户", "上游编码", "同步状态", "价格矩阵", "状态"],
    rows.map(item => [
      item.productNo,
      item.name,
      item.groupName,
      formatProductType(item.productType),
      formatChannel(item.automationConfig.channel),
      accountLabel(item.automationConfig.providerAccountId),
      item.upstreamMapping.remoteProductCode,
      formatSyncStatus(item),
      formatPricing(item),
      formatStatus(item.status)
    ])
  );
}

function exportCurrent() {
  exportProducts(filteredProducts.value, "products-current.csv");
  ElMessage.success("已导出当前列表");
}

function exportSelected() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择商品");
    return;
  }
  exportProducts(selectedRows.value, "products-selected.csv");
  ElMessage.success("已导出选中商品");
}

function openDetail(product: Product) {
  void router.push(`/catalog/products/${product.id}`);
}

function openOrderWorkbench(product: Product) {
  void router.push({
    path: "/orders/list",
    query: {
      productName: product.name
    }
  });
}

function openServiceWorkbench(product: Product) {
  void router.push({
    path: "/services/list",
    query: {
      keyword: product.name
    }
  });
}

function openProviderWorkbench(product: Product) {
  const providerAccountId = product.automationConfig.providerAccountId || product.upstreamMapping.providerAccountId;
  const providerType =
    product.automationConfig.channel && product.automationConfig.channel !== "LOCAL"
      ? product.automationConfig.channel
      : product.upstreamMapping.providerType !== "NONE"
        ? product.upstreamMapping.providerType
        : "";
  void router.push({
    path: "/providers/accounts",
    query: {
      providerType: providerType || undefined,
      accountId: providerAccountId ? String(providerAccountId) : undefined
    }
  });
}

function buildResourceTemplate(template: UpstreamProductTemplate) {
  const lookup = (code: string) => {
    const option = template.configOptions.find(item => item.code === code);
    if (!option) return "";
    const matched = option.choices.find(choice => choice.value === option.defaultValue);
    return matched?.label || option.defaultValue || option.choices[0]?.label || "";
  };

  const pickNumber = (value: string) => {
    const matched = value.match(/-?\d+/);
    return matched ? Number(matched[0]) : 0;
  };

  return {
    regionName: lookup("area") || lookup("node"),
    zoneName: lookup("node"),
    operatingSystem: lookup("os"),
    loginUsername: "root",
    securityGroup: "",
    cpuCores: pickNumber(lookup("cpu")),
    memoryGB: pickNumber(lookup("memory")),
    systemDiskGB: pickNumber(lookup("system_disk_size")),
    dataDiskGB: pickNumber(lookup("data_disk_size")),
    bandwidthMbps: pickNumber(lookup("bw")),
    publicIpCount: pickNumber(lookup("ip_num"))
  };
}

async function createLocalProductsFromUpstream() {
  const selected = importForm.importAll
    ? upstreamCatalogItems.value
    : upstreamCatalogItems.value.filter(item => importForm.remoteProductCodes.includes(item.remoteProductCode));

  let created = 0;
  let failed = 0;
  for (const item of selected) {
    try {
      const template = await fetchUpstreamProductTemplate(item.remoteProductCode, importForm.providerAccountId);
      await createProduct({
        groupName: template.groupName || item.groupName || "上游导入",
        name: template.name || item.name,
        description: template.description || item.description || "",
        productType: template.productType || item.productType || "CLOUD_HOST",
        status: "ACTIVE",
        pricing: template.pricing,
        configOptions: template.configOptions,
        resourceTemplate: buildResourceTemplate(template),
        automationConfig: {
          providerAccountId: importForm.providerAccountId,
          channel: "ZJMF_API",
          moduleType: "CLOUD",
          provisionStage: "AFTER_PAYMENT",
          autoProvision: true,
          serverGroup: "",
          providerNode: ""
        },
        upstreamMapping: {
          providerAccountId: importForm.providerAccountId,
          providerType: "ZJMF_API",
          sourceName: template.groupName || item.groupName,
          remoteProductCode: item.remoteProductCode,
          remoteProductName: template.name || item.name,
          pricePolicy: "FOLLOW_UPSTREAM",
          autoSyncPricing: importForm.autoSyncPricing,
          autoSyncConfig: importForm.autoSyncConfig,
          autoSyncTemplate: importForm.autoSyncTemplate
        }
      });
      created += 1;
    } catch {
      failed += 1;
    }
  }

  return { created, failed };
}

async function handleImportUpstream() {
  if (!importForm.providerAccountId) {
    ElMessage.warning("请先选择一个上下游接口账户");
    return;
  }
  if (!importForm.importAll && importForm.remoteProductCodes.length === 0) {
    ElMessage.warning("请至少勾选一个上游商品");
    return;
  }

  importSubmitting.value = true;
  try {
    try {
      const result = await importUpstreamProducts({
        providerAccountId: importForm.providerAccountId,
        providerType: "ZJMF_API",
        importAll: importForm.importAll,
        remoteProductCodes: importForm.remoteProductCodes,
        autoSyncPricing: importForm.autoSyncPricing,
        autoSyncConfig: importForm.autoSyncConfig,
        autoSyncTemplate: importForm.autoSyncTemplate
      });
      ElMessage.success(result.message || `导入完成，新增 ${result.created}，更新 ${result.updated}`);
    } catch {
      const fallback = await createLocalProductsFromUpstream();
      ElMessage.success(`后端批量导入不可用，已走前端兜底导入 ${fallback.created} 个商品`);
      if (fallback.failed > 0) {
        ElMessage.warning(`另外有 ${fallback.failed} 个商品导入失败`);
      }
    }

    importVisible.value = false;
    await loadProducts();
  } finally {
    importSubmitting.value = false;
  }
}

watch(
  () => importForm.providerAccountId,
  () => {
    if (importVisible.value) {
      void loadUpstreamCatalog();
    }
  }
);

watch(
  () => importForm.importAll,
  value => {
    if (value) {
      importForm.remoteProductCodes = [];
    }
  }
);

onMounted(() => {
  syncFiltersFromRoute();
  void loadProducts();
});

watch(
  () => route.fullPath,
  () => {
    syncFiltersFromRoute();
  }
);
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="商品设置"
      title="商品管理"
      subtitle="这里统一管理本地商品、自动化渠道、上下游映射和一键导入。小白用户优先用预设和导入，后续再细化配置。"
    >
      <template #actions>
        <el-button @click="loadProducts">刷新列表</el-button>
        <el-button type="success" plain @click="openImportDialog">一键导入上游商品</el-button>
        <el-button type="primary" @click="createVisible = true">新增商品</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>商品总数</span><strong>{{ summary.total }}</strong></div>
          <div class="summary-pill"><span>启用商品</span><strong>{{ summary.active }}</strong></div>
          <div class="summary-pill"><span>自动开通</span><strong>{{ summary.autoProvision }}</strong></div>
          <div class="summary-pill"><span>已绑定上游</span><strong>{{ summary.upstreamMapped }}</strong></div>
          <div class="summary-pill"><span>选中可同步</span><strong>{{ summary.selectedSyncable }}</strong></div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />
        <div class="filter-bar">
          <el-input v-model="filters.keyword" placeholder="搜索商品编号、名称、说明或上游编码" clearable />
          <el-select v-model="filters.groupName" placeholder="商品分组" clearable>
            <el-option v-for="group in groups" :key="group" :label="group" :value="group" />
          </el-select>
          <el-select v-model="filters.status" placeholder="状态" clearable>
            <el-option label="启用" value="ACTIVE" />
            <el-option label="停用" value="INACTIVE" />
          </el-select>
          <div class="action-group">
            <el-button type="primary" @click="loadProducts">应用筛选</el-button>
            <el-button plain @click="advancedVisible = !advancedVisible">
              {{ advancedVisible ? "收起高级筛选" : "展开高级筛选" }}
            </el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>

        <div v-if="advancedVisible" class="filter-bar filter-bar--compact">
          <el-select v-model="filters.productType" placeholder="商品类型" clearable>
            <el-option label="云产品" value="CLOUD" />
            <el-option label="云主机" value="CLOUD_HOST" />
            <el-option label="带宽" value="BANDWIDTH" />
            <el-option label="托管" value="COLOCATION" />
            <el-option label="CDN" value="CDN" />
            <el-option label="SSL 证书" value="SSL" />
          </el-select>
          <el-select v-model="filters.channel" placeholder="自动化渠道" clearable>
            <el-option label="本地" value="LOCAL" />
            <el-option label="魔方云" value="MOFANG_CLOUD" />
            <el-option label="上下游" value="ZJMF_API" />
            <el-option label="资源池" value="RESOURCE" />
            <el-option label="手动资源" value="MANUAL" />
          </el-select>
          <el-select v-model="filters.providerAccountId" placeholder="接口账户" clearable>
            <el-option
              v-for="item in providerAccounts"
              :key="item.id"
              :label="accountLabel(item.id)"
              :value="String(item.id)"
            />
          </el-select>
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>商品列表</strong>
          <span>显示 {{ filteredProducts.length }} 条</span>
          <span>已选 {{ selectedRows.length }} 条</span>
        </div>
        <div class="action-group">
          <el-button plain @click="copySelectedProductNos">复制商品编号</el-button>
          <el-button plain :loading="statusUpdating" :disabled="selectedRows.length === 0" @click="batchUpdateStatus('ACTIVE')">
            启用选中
          </el-button>
          <el-button plain :loading="statusUpdating" :disabled="selectedRows.length === 0" @click="batchUpdateStatus('INACTIVE')">
            停用选中
          </el-button>
          <el-button plain :disabled="selectedRows.length !== 1" @click="duplicateSelectedProduct">复制商品</el-button>
          <el-button plain :disabled="!selectedSingleProduct" @click="selectedSingleProduct && openOrderWorkbench(selectedSingleProduct)">
            查看订单
          </el-button>
          <el-button plain :disabled="!selectedSingleProduct" @click="selectedSingleProduct && openServiceWorkbench(selectedSingleProduct)">
            查看服务
          </el-button>
          <el-button plain @click="exportSelected">导出选中商品</el-button>
          <el-button plain @click="exportCurrent">导出当前列表</el-button>
          <el-button type="primary" plain :loading="syncing" @click="batchSyncSelected">批量同步上游</el-button>
        </div>
      </div>

      <el-table :data="filteredProducts" border stripe @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column prop="productNo" label="商品编号" min-width="160" />
        <el-table-column prop="name" label="商品名称" min-width="240" />
        <el-table-column prop="groupName" label="分组" min-width="140" />
        <el-table-column label="类型" min-width="120">
          <template #default="{ row }">{{ formatProductType(row.productType) }}</template>
        </el-table-column>
        <el-table-column label="自动化 / 接口账户" min-width="240">
          <template #default="{ row }">{{ formatAutomation(row) }}</template>
        </el-table-column>
        <el-table-column label="上游映射" min-width="260">
          <template #default="{ row }">{{ formatUpstream(row) }}</template>
        </el-table-column>
        <el-table-column label="同步状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="row.upstreamMapping.syncStatus === 'SUCCESS' ? 'success' : 'info'" effect="light">
              {{ formatSyncStatus(row) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="价格矩阵" min-width="260">
          <template #default="{ row }">{{ formatPricing(row) }}</template>
        </el-table-column>
        <el-table-column label="配置项" min-width="90" align="center">
          <template #default="{ row }">{{ row.configOptions.length }}</template>
        </el-table-column>
        <el-table-column label="状态" min-width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" effect="light">{{ formatStatus(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="300" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button type="primary" link @click="openDetail(row)">进入工作台</el-button>
              <el-button type="primary" link @click="openOrderWorkbench(row)">订单</el-button>
              <el-button type="primary" link @click="openServiceWorkbench(row)">服务</el-button>
              <el-button v-if="canSyncProduct(row)" type="success" link :disabled="syncing" @click="syncSingleProduct(row)">
                同步
              </el-button>
              <el-dropdown>
                <el-button type="primary" link>更多</el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="openProviderWorkbench(row)">接口账户</el-dropdown-item>
                    <el-dropdown-item @click="batchUpdateStatus(row.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE', [row])">
                      {{ row.status === "ACTIVE" ? "停用当前商品" : "启用当前商品" }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div style="display: flex; justify-content: flex-end; margin: 12px 0;">
        <el-button type="success" plain @click="openImportDialog">一键导入上游商品</el-button>
      </div>

      <div class="page-table-count">共 {{ filteredProducts.length }} 条商品记录</div>

      <el-dialog v-model="createVisible" title="新增商品" width="680px" @closed="resetCreateForm">
        <el-form label-position="top">
          <div class="filter-bar filter-bar--compact">
            <el-form-item label="商品名称" style="flex: 1">
              <el-input v-model="createForm.name" />
            </el-form-item>
            <el-form-item label="商品分组" style="flex: 1">
              <el-select v-model="createForm.groupName" filterable allow-create default-first-option>
                <el-option v-for="group in groups" :key="group" :label="group" :value="group" />
              </el-select>
            </el-form-item>
          </div>

          <div class="filter-bar filter-bar--compact">
            <el-form-item label="商品类型" style="flex: 1">
              <el-select v-model="createForm.productType">
                <el-option label="云产品" value="CLOUD" />
                <el-option label="带宽" value="BANDWIDTH" />
                <el-option label="托管" value="COLOCATION" />
              </el-select>
            </el-form-item>
            <el-form-item label="状态" style="flex: 1">
              <el-select v-model="createForm.status">
                <el-option label="启用" value="ACTIVE" />
                <el-option label="停用" value="INACTIVE" />
              </el-select>
            </el-form-item>
          </div>

          <div class="filter-bar filter-bar--compact">
            <el-form-item label="月付价格" style="flex: 1">
              <el-input-number v-model="createForm.monthlyPrice" :min="0" style="width: 100%" />
            </el-form-item>
            <el-form-item label="年付价格" style="flex: 1">
              <el-input-number v-model="createForm.annualPrice" :min="0" style="width: 100%" />
            </el-form-item>
          </div>

          <el-form-item label="商品说明">
            <el-input v-model="createForm.description" type="textarea" :rows="3" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="createVisible = false">取消</el-button>
          <el-button type="primary" :loading="creating" @click="handleCreateProduct">保存商品</el-button>
        </template>
      </el-dialog>

      <el-dialog v-model="importVisible" title="一键导入上游商品" width="980px" @closed="resetImportForm">
        <div class="summary-strip" style="margin-bottom: 16px">
          <div class="summary-pill"><span>上游商品总数</span><strong>{{ importSummary.total }}</strong></div>
          <div class="summary-pill"><span>本次准备导入</span><strong>{{ importSummary.selected }}</strong></div>
          <div class="summary-pill"><span>同步价格</span><strong>{{ importForm.autoSyncPricing ? "开启" : "关闭" }}</strong></div>
          <div class="summary-pill"><span>同步配置</span><strong>{{ importForm.autoSyncConfig ? "开启" : "关闭" }}</strong></div>
          <div class="summary-pill"><span>同步云模板</span><strong>{{ importForm.autoSyncTemplate ? "开启" : "关闭" }}</strong></div>
        </div>

        <el-form label-position="top" v-loading="importLoading">
          <div class="filter-bar filter-bar--compact">
            <el-form-item label="上下游接口账户" style="flex: 1">
              <el-select v-model="importForm.providerAccountId" placeholder="请选择一个上下游接口账户">
                <el-option
                  v-for="item in financeAccounts"
                  :key="item.id"
                  :label="`${safeAccountName(item)} / ${item.baseUrl}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="导入方式" style="width: 180px">
              <el-switch v-model="importForm.importAll" inline-prompt active-text="全部导入" inactive-text="手动挑选" />
            </el-form-item>
          </div>

          <div class="filter-bar filter-bar--compact">
            <el-form-item label="同步价格" style="width: 160px">
              <el-switch v-model="importForm.autoSyncPricing" />
            </el-form-item>
            <el-form-item label="同步配置项" style="width: 160px">
              <el-switch v-model="importForm.autoSyncConfig" />
            </el-form-item>
            <el-form-item label="同步云模板" style="width: 160px">
              <el-switch v-model="importForm.autoSyncTemplate" />
            </el-form-item>
          </div>

          <el-alert
            :closable="false"
            type="info"
            title="先选上下游接口账户，再导入商品。后端批量导入不可用时，前端会自动走兜底导入，保证这页始终能把演示商品拉进来。"
            style="margin-bottom: 16px"
          />

          <el-form-item label="上游商品目录">
            <el-select
              v-model="importForm.remoteProductCodes"
              multiple
              filterable
              clearable
              :disabled="importForm.importAll"
              placeholder="关闭“全部导入”后，可以手动勾选要导入的商品"
              style="width: 100%"
            >
              <el-option
                v-for="item in upstreamCatalogItems"
                :key="item.remoteProductCode"
                :label="item.label"
                :value="item.remoteProductCode"
              />
            </el-select>
          </el-form-item>

          <div class="panel-card">
            <div class="section-card__head">
              <strong>导入预览</strong>
              <span class="section-card__meta">先看这批商品，再执行一键导入，避免一次导太多不好辨认。</span>
            </div>
            <el-empty v-if="importSummary.selected === 0" description="当前没有可导入的上游商品" />
            <div v-else style="display: grid; gap: 10px">
              <div v-for="item in importPreview" :key="item.remoteProductCode" class="summary-pill">
                <span>{{ item.label }}</span>
              </div>
              <div v-if="importSummary.selected > importPreview.length" class="section-card__meta">
                其余 {{ importSummary.selected - importPreview.length }} 项会在导入时一起处理。
              </div>
            </div>
          </div>
        </el-form>

        <template #footer>
          <el-button @click="importVisible = false">取消</el-button>
          <el-button type="primary" :loading="importSubmitting" @click="handleImportUpstream">开始导入</el-button>
        </template>
      </el-dialog>
    </PageWorkbench>
  </div>
</template>
