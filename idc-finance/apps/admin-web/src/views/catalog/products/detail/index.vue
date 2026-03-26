<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import ContextTabs from "@/components/workbench/ContextTabs.vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createProduct,
  fetchProductDetail,
  fetchProductGroups,
  fetchProviderAccounts,
  fetchUpstreamProductTemplate,
  fetchUpstreamProducts,
  syncProductUpstream,
  updateProduct,
  type AuditLog,
  type Product,
  type ProductConfigChoice,
  type ProductConfigOption,
  type ProductPriceOption,
  type ProviderAccount,
  type UpstreamCatalogGroup,
  type UpstreamProductTemplate
} from "@/api/admin";

type ProductDetailForm = {
  id: number;
  productNo: string;
  groupName: string;
  name: string;
  description: string;
  productType: string;
  status: string;
  pricing: ProductPriceOption[];
  configOptions: ProductConfigOption[];
  resourceTemplate: Product["resourceTemplate"];
  automationConfig: Product["automationConfig"];
  upstreamMapping: Product["upstreamMapping"];
};

type RemoteCatalogOption = {
  remoteProductCode: string;
  name: string;
  groupName: string;
  label: string;
};

const t = {
  eyebrow: "\u5546\u54c1\u8bbe\u7f6e",
  title: "\u5546\u54c1\u5de5\u4f5c\u53f0",
  subtitle:
    "\u7edf\u4e00\u7ef4\u62a4\u5546\u54c1\u5b9a\u4e49\u3001\u81ea\u52a8\u5316\u6e20\u9053\u3001\u63a5\u53e3\u8d26\u6237\u3001\u4e0a\u6e38\u6620\u5c04\u548c\u4e91\u6a21\u677f\u3002",
  back: "\u8fd4\u56de\u5217\u8868",
  refresh: "\u5237\u65b0\u8be6\u60c5",
  sync: "\u4e00\u952e\u540c\u6b65\u4e0a\u6e38\u6a21\u677f",
  save: "\u4fdd\u5b58\u5546\u54c1",
  products: "\u5546\u54c1\u5217\u8868",
  workbench: "\u5546\u54c1\u5de5\u4f5c\u53f0",
  basic: "\u57fa\u7840\u4fe1\u606f",
  automation: "\u81ea\u52a8\u5316\u914d\u7f6e",
  upstream: "\u4e0a\u6e38\u6620\u5c04",
  pricing: "\u4ef7\u683c\u77e9\u9635",
  config: "\u914d\u7f6e\u9879",
  template: "\u4e91\u6a21\u677f",
  audit: "\u5ba1\u8ba1\u65e5\u5fd7"
} as const;

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const saving = ref(false);
const syncing = ref(false);
const copying = ref(false);
const statusUpdating = ref(false);
const previewLoading = ref(false);
const groupsLoading = ref(false);
const product = ref<ProductDetailForm | null>(null);
const auditLogs = ref<AuditLog[]>([]);
const providerAccounts = ref<ProviderAccount[]>([]);
const productGroups = ref<string[]>([]);
const upstreamGroups = ref<UpstreamCatalogGroup[]>([]);
const remoteTemplate = ref<UpstreamProductTemplate | null>(null);
const activeTab = ref("basic");
const showAdvancedAutomation = ref(false);

const automationPresets = [
  {
    key: "LOCAL",
    title: "\u672c\u5730\u624b\u5de5\u4ea4\u4ed8",
    description: "\u9002\u5408\u81ea\u5efa\u5e93\u5b58\u3001\u4eba\u5de5\u5ba1\u6838\u3001\u4eba\u5de5\u5f00\u901a\u7684\u5546\u54c1\u3002"
  },
  {
    key: "MOFANG_CLOUD",
    title: "\u9b54\u65b9\u4e91\u81ea\u52a8\u5f00\u901a",
    description: "\u4ed8\u6b3e\u540e\u81ea\u52a8\u5411\u9b54\u65b9\u4e91\u521b\u5efa\u5b9e\u4f8b\u5e76\u56de\u5199\u672c\u5730\u670d\u52a1\u3002"
  },
  {
    key: "ZJMF_API",
    title: "\u8d22\u52a1\u4e0a\u4e0b\u6e38\u540c\u6b65",
    description: "\u540c\u6b65\u4e0a\u6e38\u5546\u54c1\u3001\u4ef7\u683c\u3001\u914d\u7f6e\u9879\uff0c\u5e76\u6309\u4e0a\u6e38\u94fe\u8def\u4ea4\u4ed8\u3002"
  }
] as const;

const contextTabs = computed(() => [
  { key: "products", label: t.products, to: "/catalog/products" },
  { key: "detail", label: t.workbench, active: true, badge: product.value?.productNo || "-" }
]);

const financeAccounts = computed(() =>
  providerAccounts.value.filter(item => item.providerType === "ZJMF_API" && item.status === "ACTIVE")
);
const mofangAccounts = computed(() =>
  providerAccounts.value.filter(item => item.providerType === "MOFANG_CLOUD" && item.status === "ACTIVE")
);
const remoteCatalogOptions = computed<RemoteCatalogOption[]>(() =>
  upstreamGroups.value.flatMap(group =>
    group.products.map(item => ({
      remoteProductCode: item.remoteProductCode,
      name: item.name,
      groupName: group.groupName,
      label: `${group.groupName} / ${item.name}`
    }))
  )
);

const automationAccountOptions = computed(() => {
  const activeAccounts = providerAccounts.value.filter(item => item.status === "ACTIVE");
  const channel = product.value?.automationConfig.channel;
  if (channel === "MOFANG_CLOUD") {
    return activeAccounts.filter(item => item.providerType === "MOFANG_CLOUD");
  }
  if (channel === "ZJMF_API") {
    return activeAccounts.filter(item => item.providerType === "ZJMF_API");
  }
  return [];
});

const automationGuide = computed(() => {
  const channel = product.value?.automationConfig.channel || "LOCAL";
  const autoProvision = Boolean(product.value?.automationConfig.autoProvision);
  const accountName = safeAccountName(product.value?.automationConfig.providerAccountId);
  const stageLabel = formatProvisionStage(product.value?.automationConfig.provisionStage || "MANUAL_REVIEW");
  if (channel === "MOFANG_CLOUD") {
    return {
      type: "success" as const,
      activeStep: autoProvision ? 4 : 3,
      title: "当前是魔方云自动开通商品",
      description: `客户付款后，系统会通过魔方云接口账户自动创建实例，并把实例编号、IP 和同步状态回写到本地服务。当前开通阶段：${stageLabel}。`,
      accountPlaceholder: automationAccountOptions.value.length ? "选择用于开通实例的魔方云账户" : "请先到接口账户里新增魔方云账户",
      accountRequired: true,
      showAccount: true,
      tipText: `当前执行账户：${accountName}。建议先确认地域、镜像和默认模板，再开启自动开通。`
    };
  }
  if (channel === "ZJMF_API") {
    return {
      type: "warning" as const,
      activeStep: autoProvision ? 4 : 3,
      title: "当前是财务上下游同步商品",
      description: `系统会先同步上游商品、价格和配置项，再按上游财务链路下单、支付和交付。当前开通阶段：${stageLabel}。`,
      accountPlaceholder: automationAccountOptions.value.length ? "选择用于同步与交付的上游财务账户" : "请先到接口账户里新增财务上下游账户",
      accountRequired: true,
      showAccount: true,
      tipText: `当前执行账户：${accountName}。如果要完全跟随上游，建议把价格、配置项和云模板同步开关都打开。`
    };
  }
  if (channel === "RESOURCE") {
    return {
      type: "info" as const,
      activeStep: 2,
      title: "当前是资源池商品",
      description: "适合从资源池或库存里分配现成资源。系统不会自动向云平台下单，建议由运维审核后再交付。",
      accountPlaceholder: "资源池商品通常不需要接口账户",
      accountRequired: false,
      showAccount: false,
      tipText: "如果这类商品最终仍需对接云平台，请改成魔方云或财务上下游渠道。"
    };
  }
  if (channel === "MANUAL") {
    return {
      type: "info" as const,
      activeStep: 2,
      title: "当前是手工资源商品",
      description: "适合临时业务、代维业务或特殊资源。订单支付后不会自动创建实例，需要管理员手工处理。",
      accountPlaceholder: "手工资源不需要接口账户",
      accountRequired: false,
      showAccount: false,
      tipText: "建议关闭自动开通，并在服务工作台中补齐远端资源编号、IP 和到期时间。"
    };
  }
  return {
    type: "info" as const,
    activeStep: autoProvision ? 3 : 2,
    title: "当前是本地交付商品",
    description: `系统只保留本地订单、账单和服务记录，不会自动调用外部平台。当前开通阶段：${stageLabel}。`,
    accountPlaceholder: "本地交付不需要接口账户",
    accountRequired: false,
    showAccount: false,
    tipText: autoProvision ? "已开启自动开通，建议确认本地模块和资源模板已经准备好。" : "建议保留人工审核，便于客服或运维二次确认。"
  };
});

const pricingDiffRows = computed(() => {
  if (!product.value || !remoteTemplate.value) return [];
  const remoteMap = new Map(remoteTemplate.value.pricing.map(item => [item.cycleCode, item]));
  const localMap = new Map(product.value.pricing.map(item => [item.cycleCode, item]));
  const cycles = new Set([...remoteMap.keys(), ...localMap.keys()]);
  return Array.from(cycles).map(cycleCode => {
    const local = localMap.get(cycleCode);
    const remote = remoteMap.get(cycleCode);
    return {
      cycleCode,
      cycleName: remote?.cycleName || local?.cycleName || cycleCode,
      localPrice: local?.price ?? null,
      remotePrice: remote?.price ?? null
    };
  });
});

const configDiffRows = computed(() => {
  if (!product.value || !remoteTemplate.value) return [];
  const remoteMap = new Map(remoteTemplate.value.configOptions.map(item => [item.code, item]));
  const localMap = new Map(product.value.configOptions.map(item => [item.code, item]));
  const codes = new Set([...remoteMap.keys(), ...localMap.keys()]);
  return Array.from(codes).map(code => {
    const local = localMap.get(code);
    const remote = remoteMap.get(code);
    return {
      code,
      name: remote?.name || local?.name || code,
      localDefault: local?.defaultValue || "-",
      remoteDefault: remote?.defaultValue || "-"
    };
  });
});

const templateDiffRows = computed(() => {
  if (!product.value || !remoteTemplate.value) return [];
  const derived = buildResourceTemplateFromRemote(remoteTemplate.value);
  return [
    { label: "\u5730\u57df", local: product.value.resourceTemplate.regionName, remote: derived.regionName },
    { label: "\u53ef\u7528\u533a", local: product.value.resourceTemplate.zoneName, remote: derived.zoneName },
    { label: "\u64cd\u4f5c\u7cfb\u7edf", local: product.value.resourceTemplate.operatingSystem, remote: derived.operatingSystem },
    { label: "CPU", local: `${product.value.resourceTemplate.cpuCores}\u6838`, remote: `${derived.cpuCores}\u6838` },
    { label: "\u5185\u5b58", local: `${product.value.resourceTemplate.memoryGB} GB`, remote: `${derived.memoryGB} GB` },
    { label: "\u7cfb\u7edf\u76d8", local: `${product.value.resourceTemplate.systemDiskGB} GB`, remote: `${derived.systemDiskGB} GB` }
  ];
});

const upstreamDiffSummary = computed(() => ({
  pricingChanged: pricingDiffRows.value.filter(row => row.localPrice !== row.remotePrice).length,
  configChanged: configDiffRows.value.filter(row => row.localDefault !== row.remoteDefault).length,
  templateChanged: templateDiffRows.value.filter(row => row.local !== row.remote).length,
  remotePricingCount: remoteTemplate.value?.pricing.length ?? 0,
  remoteConfigCount: remoteTemplate.value?.configOptions.length ?? 0
}));

const pricingSummary = computed(() => {
  const items = product.value?.pricing ?? [];
  return {
    total: items.length,
    chargeable: items.filter(item => Number(item.price || 0) > 0).length,
    free: items.filter(item => Number(item.price || 0) <= 0).length,
    highestPrice: items.reduce((max, item) => Math.max(max, Number(item.price || 0)), 0)
  };
});

const configSummary = computed(() => {
  const items = product.value?.configOptions ?? [];
  return {
    total: items.length,
    required: items.filter(item => item.required).length,
    selectLike: items.filter(item => item.inputType === "select" || item.inputType === "radio").length,
    totalChoices: items.reduce((sum, item) => sum + item.choices.length, 0)
  };
});

const commonPricingPresets = [
  { cycleCode: "monthly", cycleName: "月付" },
  { cycleCode: "quarterly", cycleName: "季付" },
  { cycleCode: "semiannual", cycleName: "半年付" },
  { cycleCode: "annual", cycleName: "年付" },
  { cycleCode: "biennially", cycleName: "两年付" },
  { cycleCode: "triennially", cycleName: "三年付" },
  { cycleCode: "onetime", cycleName: "一次性" }
] as const;

function cloneProduct(input: Product): ProductDetailForm {
  return {
    id: input.id,
    productNo: input.productNo,
    groupName: input.groupName,
    name: input.name,
    description: input.description,
    productType: input.productType,
    status: input.status,
    pricing: input.pricing.map(item => ({ ...item })),
    configOptions: input.configOptions.map(item => ({
      ...item,
      choices: item.choices.map(choice => ({ ...choice }))
    })),
    resourceTemplate: { ...input.resourceTemplate },
    automationConfig: { ...input.automationConfig },
    upstreamMapping: { ...input.upstreamMapping }
  };
}

function buildProductPayload(source: ProductDetailForm, overrides: Record<string, unknown> = {}) {
  return {
    groupName: source.groupName,
    name: source.name,
    description: source.description,
    productType: source.productType,
    status: source.status,
    pricing: source.pricing.map(item => ({ ...item })),
    configOptions: cloneConfigOptions(source.configOptions),
    resourceTemplate: { ...source.resourceTemplate },
    automationConfig: { ...source.automationConfig },
    upstreamMapping: { ...source.upstreamMapping },
    ...overrides
  };
}

function safeAccountName(accountId?: number) {
  if (!accountId) return "\u672a\u7ed1\u5b9a";
  const account = providerAccounts.value.find(item => item.id === accountId);
  if (!account) return "\u672a\u7ed1\u5b9a";
  return account.name || account.sourceName || `#${account.id}`;
}

function formatProductType(type: string) {
  const mapping: Record<string, string> = {
    CLOUD: "\u4e91\u4ea7\u54c1",
    CLOUD_HOST: "\u4e91\u4e3b\u673a",
    BANDWIDTH: "\u5e26\u5bbd",
    COLOCATION: "\u6258\u7ba1",
    CDN: "CDN",
    SSL: "SSL",
    SERVER: "\u7269\u7406\u670d\u52a1\u5668",
    DCIM: "DCIM",
    DCIMCLOUD: "\u5f39\u6027\u4e91"
  };
  return mapping[type] ?? type;
}

function formatChannel(channel: string) {
  const mapping: Record<string, string> = {
    LOCAL: "\u672c\u5730\u4ea4\u4ed8",
    MOFANG_CLOUD: "\u9b54\u65b9\u4e91",
    ZJMF_API: "\u8d22\u52a1\u4e0a\u4e0b\u6e38",
    RESOURCE: "\u8d44\u6e90\u6c60",
    MANUAL: "\u624b\u5de5\u8d44\u6e90"
  };
  return mapping[channel] ?? channel;
}

function formatModuleType(type: string) {
  const mapping: Record<string, string> = {
    NORMAL: "\u666e\u901a\u6a21\u5757",
    CLOUD: "\u4e91\u4ea7\u54c1\u6a21\u5757",
    DCIM: "DCIM",
    DCIMCLOUD: "\u9b54\u65b9\u4e91",
    RESOURCE: "\u8d44\u6e90\u6c60\u6a21\u5757"
  };
  return mapping[type] ?? type;
}

function formatProvisionStage(stage: string) {
  const mapping: Record<string, string> = {
    AFTER_PAYMENT: "\u4ed8\u6b3e\u540e\u81ea\u52a8\u5f00\u901a",
    MANUAL_REVIEW: "\u4eba\u5de5\u5ba1\u6838\u540e\u5f00\u901a",
    AFTER_CONFIRM: "\u5ba2\u6237\u786e\u8ba4\u540e\u5f00\u901a"
  };
  return mapping[stage] ?? stage;
}

function formatSyncStatus(status: string) {
  const mapping: Record<string, string> = {
    SUCCESS: "\u5df2\u540c\u6b65",
    FAILED: "\u540c\u6b65\u5931\u8d25",
    PENDING: "\u5f85\u540c\u6b65",
    NONE: "\u672a\u540c\u6b65"
  };
  return mapping[status] ?? (status || "\u672a\u540c\u6b65");
}

function pricePolicyLabel(policy: string) {
  const mapping: Record<string, string> = {
    FOLLOW_UPSTREAM: "\u8ddf\u968f\u4e0a\u6e38",
    LOCAL_OVERRIDE: "\u672c\u5730\u8986\u76d6",
    PERCENTAGE: "\u6309\u6bd4\u4f8b\u52a0\u4ef7"
  };
  return mapping[policy] ?? (policy || "\u8ddf\u968f\u4e0a\u6e38");
}

function buildResourceTemplateFromRemote(template: UpstreamProductTemplate) {
  const findOption = (code: string) => template.configOptions.find(item => item.code === code);
  const findLabel = (code: string) => {
    const option = findOption(code);
    if (!option) return "";
    const matched = option.choices.find(choice => choice.value === option.defaultValue);
    return matched?.label || option.defaultValue || option.choices[0]?.label || "";
  };
  const extractNumber = (value: string, fallback = 0) => {
    const match = value.match(/\d+/);
    return match ? Number(match[0]) : fallback;
  };
  return {
    regionName: findLabel("area") || findLabel("node"),
    zoneName: findLabel("node") || findLabel("area"),
    operatingSystem: findLabel("os"),
    cpuCores: extractNumber(findLabel("cpu"), 2),
    memoryGB: extractNumber(findLabel("memory"), 4),
    systemDiskGB: extractNumber(findLabel("system_disk_size"), 40),
    dataDiskGB: extractNumber(findLabel("data_disk_size"), 0),
    bandwidthMbps: extractNumber(findLabel("bw"), 20),
    publicIpCount: extractNumber(findLabel("ip_num"), 1)
  };
}

function cloneConfigOptions(options: ProductConfigOption[]) {
  return options.map(item => ({
    ...item,
    choices: item.choices.map(choice => ({ ...choice }))
  }));
}

function applyRemotePricing() {
  if (!product.value || !remoteTemplate.value) return;
  product.value.pricing = remoteTemplate.value.pricing.map(item => ({ ...item }));
  ElMessage.success("已将上游价格差异回填到本地价格矩阵");
}

function applyRemoteConfigOptions() {
  if (!product.value || !remoteTemplate.value) return;
  product.value.configOptions = cloneConfigOptions(remoteTemplate.value.configOptions);
  ElMessage.success("已将上游配置项回填到本地配置草稿");
}

function applyRemoteResourceTemplate() {
  if (!product.value || !remoteTemplate.value) return;
  product.value.resourceTemplate = {
    ...product.value.resourceTemplate,
    ...buildResourceTemplateFromRemote(remoteTemplate.value)
  };
  ElMessage.success("已将上游模板推导结果回填到本地云模板");
}

function applyRemoteAll() {
  applyRemotePricing();
  applyRemoteConfigOptions();
  applyRemoteResourceTemplate();
  ElMessage.success("已将上游价格、配置项和云模板全部回填到本地草稿");
}

function applyPreset(channel: "LOCAL" | "MOFANG_CLOUD" | "ZJMF_API") {
  if (!product.value) return;
  if (channel === "LOCAL") {
    product.value.automationConfig.channel = "LOCAL";
    product.value.automationConfig.moduleType = "NORMAL";
    product.value.automationConfig.provisionStage = "MANUAL_REVIEW";
    product.value.automationConfig.autoProvision = false;
    product.value.automationConfig.providerAccountId = 0;
    return;
  }
  if (channel === "MOFANG_CLOUD") {
    product.value.automationConfig.channel = "MOFANG_CLOUD";
    product.value.automationConfig.moduleType = "DCIMCLOUD";
    product.value.automationConfig.provisionStage = "AFTER_PAYMENT";
    product.value.automationConfig.autoProvision = true;
    product.value.automationConfig.providerAccountId ||= mofangAccounts.value[0]?.id || 0;
    product.value.upstreamMapping.providerType = "MOFANG_CLOUD";
    product.value.upstreamMapping.providerAccountId = product.value.automationConfig.providerAccountId;
    product.value.upstreamMapping.sourceName ||= "\u9b54\u65b9\u4e91";
    return;
  }
  product.value.automationConfig.channel = "ZJMF_API";
  product.value.automationConfig.moduleType = "CLOUD";
  product.value.automationConfig.provisionStage = "AFTER_PAYMENT";
  product.value.automationConfig.autoProvision = true;
  product.value.automationConfig.providerAccountId ||= financeAccounts.value[0]?.id || 0;
  product.value.upstreamMapping.providerType = "ZJMF_API";
  product.value.upstreamMapping.providerAccountId = product.value.automationConfig.providerAccountId;
  product.value.upstreamMapping.sourceName ||= "\u4e0a\u6e38\u8d22\u52a1";
}

function addPricingRow() {
  product.value?.pricing.push({
    cycleCode: `custom-${product.value.pricing.length + 1}`,
    cycleName: "\u81ea\u5b9a\u4e49\u5468\u671f",
    price: 0,
    setupFee: 0
  });
}

function addCommonPricingRows() {
  if (!product.value) return;
  const existing = new Set(product.value.pricing.map(item => item.cycleCode));
  const missing = commonPricingPresets.filter(item => !existing.has(item.cycleCode));
  product.value.pricing.push(
    ...missing.map(item => ({
      cycleCode: item.cycleCode,
      cycleName: item.cycleName,
      price: 0,
      setupFee: 0
    }))
  );
  if (missing.length > 0) {
    ElMessage.success(`已补齐 ${missing.length} 个常用计费周期`);
  }
}

function sortPricingRows() {
  if (!product.value) return;
  const order = new Map<string, number>(commonPricingPresets.map((item, index) => [item.cycleCode, index]));
  product.value.pricing = [...product.value.pricing].sort((a, b) => {
    const aIndex = order.get(a.cycleCode) ?? 999;
    const bIndex = order.get(b.cycleCode) ?? 999;
    if (aIndex !== bIndex) return aIndex - bIndex;
    return a.cycleCode.localeCompare(b.cycleCode);
  });
}

function duplicatePricingRow(index: number) {
  if (!product.value) return;
  const target = product.value.pricing[index];
  if (!target) return;
  product.value.pricing.splice(index + 1, 0, {
    ...target,
    cycleCode: `${target.cycleCode}_copy`,
    cycleName: `${target.cycleName} 副本`
  });
}

function removePricingRow(index: number) {
  product.value?.pricing.splice(index, 1);
}

function addConfigOption() {
  product.value?.configOptions.push({
    code: `config_${product.value.configOptions.length + 1}`,
    name: "\u65b0\u914d\u7f6e\u9879",
    inputType: "select",
    required: false,
    defaultValue: "",
    description: "",
    choices: []
  });
}

function addSelectConfigOption() {
  addConfigOption();
  const option = product.value?.configOptions.at(-1);
  if (!option) return;
  option.name = "新下拉配置项";
  option.inputType = "select";
  addChoice(option);
}

function addTextConfigOption() {
  addConfigOption();
  const option = product.value?.configOptions.at(-1);
  if (!option) return;
  option.name = "新文本配置项";
  option.inputType = "text";
  option.choices = [];
}

function duplicateConfigOption(index: number) {
  if (!product.value) return;
  const target = product.value.configOptions[index];
  if (!target) return;
  product.value.configOptions.splice(index + 1, 0, {
    ...target,
    code: `${target.code}_copy`,
    name: `${target.name} 副本`,
    choices: target.choices.map(choice => ({ ...choice }))
  });
}

function removeConfigOption(index: number) {
  product.value?.configOptions.splice(index, 1);
}

function addChoice(option: ProductConfigOption) {
  option.choices.push({
    value: `value_${option.choices.length + 1}`,
    label: "\u65b0\u9009\u9879",
    priceDelta: 0
  });
  if (!option.defaultValue) option.defaultValue = option.choices[0]?.value || "";
}

function removeChoice(option: ProductConfigOption, index: number) {
  option.choices.splice(index, 1);
  if (option.defaultValue && !option.choices.some(choice => choice.value === option.defaultValue)) {
    option.defaultValue = option.choices[0]?.value || "";
  }
}

async function loadProviderAccountList() {
  providerAccounts.value = await fetchProviderAccounts();
}

async function loadGroupList() {
  groupsLoading.value = true;
  try {
    productGroups.value = await fetchProductGroups();
  } finally {
    groupsLoading.value = false;
  }
}

async function loadUpstreamCatalog(accountId?: number) {
  const targetId = accountId || product.value?.upstreamMapping.providerAccountId;
  if (!targetId) {
    upstreamGroups.value = [];
    return;
  }
  upstreamGroups.value = (await fetchUpstreamProducts(targetId)).groups;
}

function applyRemoteSelection(code: string) {
  if (!product.value) return;
  const matched = remoteCatalogOptions.value.find(item => item.remoteProductCode === code);
  if (!matched) return;
  product.value.upstreamMapping.remoteProductCode = matched.remoteProductCode;
  product.value.upstreamMapping.remoteProductName = matched.name;
  product.value.upstreamMapping.sourceName = matched.groupName;
}

async function loadRemotePreview() {
  if (!product.value) return;
  if (product.value.upstreamMapping.providerType !== "ZJMF_API") {
    remoteTemplate.value = null;
    return;
  }
  if (!product.value.upstreamMapping.providerAccountId || !product.value.upstreamMapping.remoteProductCode) {
    remoteTemplate.value = null;
    return;
  }
  previewLoading.value = true;
  try {
    remoteTemplate.value = await fetchUpstreamProductTemplate(
      product.value.upstreamMapping.remoteProductCode,
      product.value.upstreamMapping.providerAccountId
    );
  } catch {
    remoteTemplate.value = null;
  } finally {
    previewLoading.value = false;
  }
}

async function loadProductDetail() {
  loading.value = true;
  try {
    const detail = await fetchProductDetail(route.params.id as string);
    product.value = cloneProduct(detail.product);
    auditLogs.value = detail.auditLogs;
    if (product.value.upstreamMapping.providerType === "ZJMF_API" && product.value.upstreamMapping.providerAccountId) {
      await loadUpstreamCatalog(product.value.upstreamMapping.providerAccountId);
    } else {
      upstreamGroups.value = [];
    }
    await loadRemotePreview();
  } finally {
    loading.value = false;
  }
}

async function handleSave() {
  if (!product.value) return;
  saving.value = true;
  try {
    const updated = await updateProduct(product.value.id, buildProductPayload(product.value));
    product.value = cloneProduct(updated);
    ElMessage.success(`\u5546\u54c1 ${updated.productNo} \u5df2\u4fdd\u5b58`);
  } finally {
    saving.value = false;
  }
}

async function handleSync() {
  if (!product.value) return;
  syncing.value = true;
  try {
    const updated = await syncProductUpstream(product.value.id, product.value.upstreamMapping);
    product.value = cloneProduct(updated);
    ElMessage.success(`\u5546\u54c1 ${updated.productNo} \u5df2\u540c\u6b65`);
    await loadRemotePreview();
  } finally {
    syncing.value = false;
  }
}

async function duplicateCurrentProduct() {
  if (!product.value) return;
  copying.value = true;
  try {
    const created = await createProduct(
      buildProductPayload(product.value, {
        name: `${product.value.name} - 副本`,
        description: product.value.description
          ? `${product.value.description}\n\n复制自 ${product.value.productNo}`
          : `复制自 ${product.value.productNo}`,
        upstreamMapping: {
          providerAccountId: 0,
          providerType: "NONE",
          sourceName: "",
          remoteProductCode: "",
          remoteProductName: "",
          pricePolicy: product.value.upstreamMapping.pricePolicy,
          autoSyncPricing: false,
          autoSyncConfig: false,
          autoSyncTemplate: false
        }
      })
    );
    ElMessage.success(`已复制商品 ${created.productNo}`);
    await router.push(`/catalog/products/${created.id}`);
  } finally {
    copying.value = false;
  }
}

async function toggleCurrentProductStatus() {
  if (!product.value) return;
  statusUpdating.value = true;
  const nextStatus = product.value.status === "ACTIVE" ? "INACTIVE" : "ACTIVE";
  try {
    const updated = await updateProduct(
      product.value.id,
      buildProductPayload(product.value, {
        status: nextStatus
      })
    );
    product.value = cloneProduct(updated);
    ElMessage.success(`商品状态已切换为${nextStatus === "ACTIVE" ? "启用" : "停用"}`);
  } finally {
    statusUpdating.value = false;
  }
}

function openAutomationAccountWorkbench() {
  if (!product.value) return;
  const providerType = product.value.automationConfig.channel;
  if (providerType === "LOCAL" || providerType === "MANUAL") return;
  void router.push({
    path: "/providers/accounts",
    query: {
      providerType,
      accountId: product.value.automationConfig.providerAccountId
        ? String(product.value.automationConfig.providerAccountId)
        : undefined
    }
  });
}

function openUpstreamAccountWorkbench() {
  if (!product.value || product.value.upstreamMapping.providerType === "NONE") return;
  void router.push({
    path: "/providers/accounts",
    query: {
      providerType: product.value.upstreamMapping.providerType,
      accountId: product.value.upstreamMapping.providerAccountId
        ? String(product.value.upstreamMapping.providerAccountId)
        : undefined
    }
  });
}

function openProductOrdersWorkbench() {
  if (!product.value) return;
  void router.push({
    path: "/orders/list",
    query: {
      productName: product.value.name
    }
  });
}

function openProductServicesWorkbench() {
  if (!product.value) return;
  void router.push({
    path: "/services/list",
    query: {
      keyword: product.value.name,
      providerType: product.value.automationConfig.channel !== "LOCAL" ? product.value.automationConfig.channel : undefined,
      providerAccountId: product.value.automationConfig.providerAccountId
        ? String(product.value.automationConfig.providerAccountId)
        : undefined
    }
  });
}

function openProductAutomationWorkbench() {
  if (!product.value) return;
  void router.push({
    path: "/providers/automation",
    query: {
      sourceType: "product",
      sourceId: String(product.value.id),
      channel: product.value.automationConfig.channel !== "LOCAL" ? product.value.automationConfig.channel : undefined
    }
  });
}

function openProductResourcesWorkbench() {
  if (!product.value) return;
  const accountId = product.value.automationConfig.providerAccountId || product.value.upstreamMapping.providerAccountId;
  const providerType =
    product.value.automationConfig.channel !== "LOCAL"
      ? product.value.automationConfig.channel
      : product.value.upstreamMapping.providerType !== "NONE"
        ? product.value.upstreamMapping.providerType
        : undefined;
  void router.push({
    path: "/providers/resources",
    query: {
      providerType,
      accountId: accountId ? String(accountId) : undefined
    }
  });
}

function openUpstreamSyncHistoryWorkbench() {
  if (!product.value?.upstreamMapping.providerAccountId) return;
  void router.push({
    path: "/providers/upstream-sync-history",
    query: {
      providerAccountId: String(product.value.upstreamMapping.providerAccountId),
      providerType: product.value.upstreamMapping.providerType || undefined,
      keyword: product.value.upstreamMapping.remoteProductCode || product.value.productNo
    }
  });
}

watch(
  () => route.params.id,
  () => {
    void loadProductDetail();
  }
);

watch(
  () => product.value?.upstreamMapping.providerAccountId,
  async accountId => {
    if (!product.value) return;
    if (product.value.upstreamMapping.providerType === "ZJMF_API") {
      await loadUpstreamCatalog(accountId);
      await loadRemotePreview();
    }
  }
);

watch(
  () => product.value?.upstreamMapping.remoteProductCode,
  async code => {
    if (!code || !product.value) return;
    applyRemoteSelection(code);
    if (product.value.upstreamMapping.providerType === "ZJMF_API") {
      await loadRemotePreview();
    }
  }
);

watch(
  () => product.value?.automationConfig.channel,
  channel => {
    if (!product.value || !channel) return;
    const automation = product.value.automationConfig;
    const upstream = product.value.upstreamMapping;
    const availableIds = new Set(automationAccountOptions.value.map(item => item.id));

    if (channel === "MOFANG_CLOUD") {
      automation.moduleType = "DCIMCLOUD";
      automation.provisionStage ||= "AFTER_PAYMENT";
      automation.providerAccountId = availableIds.has(automation.providerAccountId)
        ? automation.providerAccountId
        : automationAccountOptions.value[0]?.id || 0;
      upstream.providerType = "MOFANG_CLOUD";
      upstream.providerAccountId = automation.providerAccountId;
      upstream.sourceName ||= "\u9b54\u65b9\u4e91";
      return;
    }

    if (channel === "ZJMF_API") {
      automation.moduleType = automation.moduleType === "DCIMCLOUD" ? "CLOUD" : automation.moduleType || "CLOUD";
      automation.provisionStage ||= "AFTER_PAYMENT";
      automation.providerAccountId = availableIds.has(automation.providerAccountId)
        ? automation.providerAccountId
        : automationAccountOptions.value[0]?.id || 0;
      upstream.providerType = "ZJMF_API";
      upstream.providerAccountId = automation.providerAccountId;
      upstream.sourceName ||= "\u4e0a\u6e38\u8d22\u52a1";
      return;
    }

    automation.providerAccountId = 0;
    if (channel === "LOCAL") {
      automation.moduleType = "NORMAL";
    }
  }
);

watch(
  () => product.value?.automationConfig.providerAccountId,
  accountId => {
    if (!product.value) return;
    if (product.value.automationConfig.channel === "MOFANG_CLOUD" && product.value.upstreamMapping.providerType === "MOFANG_CLOUD") {
      product.value.upstreamMapping.providerAccountId = accountId || 0;
    }
    if (product.value.automationConfig.channel === "ZJMF_API" && product.value.upstreamMapping.providerType === "ZJMF_API") {
      product.value.upstreamMapping.providerAccountId = accountId || 0;
    }
  }
);

onMounted(async () => {
  await Promise.all([loadProviderAccountList(), loadGroupList()]);
  await loadProductDetail();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench :eyebrow="t.eyebrow" :title="t.title" :subtitle="t.subtitle">
      <template #actions>
        <el-button @click="router.push('/catalog/products')">{{ t.back }}</el-button>
        <el-button @click="loadProductDetail">{{ t.refresh }}</el-button>
        <el-button plain :loading="copying" @click="duplicateCurrentProduct">复制商品</el-button>
        <el-button plain :loading="statusUpdating" @click="toggleCurrentProductStatus">
          {{ product?.status === "ACTIVE" ? "停用商品" : "启用商品" }}
        </el-button>
        <el-button
          v-if="product && product.upstreamMapping.providerType !== 'NONE'"
          type="success"
          plain
          :loading="syncing"
          @click="handleSync"
        >
          {{ t.sync }}
        </el-button>
        <el-button
          v-if="product && product.upstreamMapping.providerAccountId"
          plain
          @click="openUpstreamSyncHistoryWorkbench"
        >
          上游同步记录
        </el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">{{ t.save }}</el-button>
      </template>

      <template #context>
        <ContextTabs :items="contextTabs" />
      </template>

      <template v-if="product" #metrics>
        <div class="detail-kpi-grid">
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">编号</span>
            <strong class="detail-kpi-card__value">{{ product.productNo }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">渠道</span>
            <strong class="detail-kpi-card__value">{{ formatChannel(product.automationConfig.channel) }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">接口账户</span>
            <strong class="detail-kpi-card__value">{{ safeAccountName(product.automationConfig.providerAccountId) }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">同步状态</span>
            <strong class="detail-kpi-card__value">{{ formatSyncStatus(product.upstreamMapping.syncStatus) }}</strong>
          </div>
        </div>
      </template>

      <template v-if="product">
        <el-tabs v-model="activeTab">
          <el-tab-pane :label="t.basic" name="basic">
            <div class="portal-grid portal-grid--two">
              <div class="panel-card">
                <div class="section-card__head"><strong>商品基础资料</strong></div>
                <div class="filter-bar filter-bar--compact">
                  <el-input v-model="product.name" placeholder="商品名称" />
                  <el-select v-model="product.groupName" filterable allow-create default-first-option :loading="groupsLoading">
                    <el-option v-for="group in productGroups" :key="group" :label="group" :value="group" />
                  </el-select>
                </div>
                <div class="filter-bar filter-bar--compact">
                  <el-select v-model="product.productType">
                    <el-option label="云产品" value="CLOUD" />
                    <el-option label="云主机" value="CLOUD_HOST" />
                    <el-option label="带宽" value="BANDWIDTH" />
                    <el-option label="托管" value="COLOCATION" />
                    <el-option label="CDN" value="CDN" />
                  </el-select>
                  <el-select v-model="product.status">
                    <el-option label="启用" value="ACTIVE" />
                    <el-option label="停用" value="INACTIVE" />
                  </el-select>
                </div>
                <el-input v-model="product.description" type="textarea" :rows="4" placeholder="商品说明" />
              </div>

              <div class="panel-card">
                <div class="section-card__head"><strong>当前交付识别</strong></div>
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="商品类型">{{ formatProductType(product.productType) }}</el-descriptions-item>
                  <el-descriptions-item label="自动化渠道">{{ formatChannel(product.automationConfig.channel) }}</el-descriptions-item>
                  <el-descriptions-item label="模块类型">{{ formatModuleType(product.automationConfig.moduleType) }}</el-descriptions-item>
                  <el-descriptions-item label="接口账户">{{ safeAccountName(product.automationConfig.providerAccountId) }}</el-descriptions-item>
                  <el-descriptions-item label="上游编码">{{ product.upstreamMapping.remoteProductCode || "-" }}</el-descriptions-item>
                  <el-descriptions-item label="同步状态">{{ formatSyncStatus(product.upstreamMapping.syncStatus) }}</el-descriptions-item>
                </el-descriptions>
              </div>
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>业务联动</strong>
                  <span class="section-card__meta">从商品直接跳到接口账户、订单、服务和自动化中心。</span>
                </div>
                <div class="summary-strip">
                  <div class="summary-pill"><span>自动化渠道</span><strong>{{ formatChannel(product.automationConfig.channel) }}</strong></div>
                  <div class="summary-pill"><span>自动开通</span><strong>{{ product.automationConfig.autoProvision ? "已开启" : "人工处理" }}</strong></div>
                  <div class="summary-pill"><span>上游策略</span><strong>{{ pricePolicyLabel(product.upstreamMapping.pricePolicy) }}</strong></div>
                </div>
                <div class="inline-actions" style="margin-top: 16px">
                  <el-button
                    plain
                    :disabled="product.automationConfig.channel === 'LOCAL' || product.automationConfig.channel === 'MANUAL'"
                    @click="openAutomationAccountWorkbench"
                  >
                    自动化账户
                  </el-button>
                  <el-button plain :disabled="product.upstreamMapping.providerType === 'NONE'" @click="openUpstreamAccountWorkbench">
                    上游账户
                  </el-button>
                  <el-button plain @click="openProductOrdersWorkbench">订单列表</el-button>
                  <el-button plain @click="openProductServicesWorkbench">服务列表</el-button>
                  <el-button plain @click="openProductAutomationWorkbench">自动化任务</el-button>
                  <el-button plain @click="openProductResourcesWorkbench">渠道资源</el-button>
                  <el-button plain :disabled="!product.upstreamMapping.providerAccountId" @click="openUpstreamSyncHistoryWorkbench">
                    上游同步记录
                  </el-button>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t.automation" name="automation">
            <el-alert
              :title="automationGuide.title"
              :description="automationGuide.description"
              :type="automationGuide.type"
              :closable="false"
              show-icon
            />
            <div class="panel-card" style="margin-top: 16px">
              <el-steps :active="automationGuide.activeStep" simple finish-status="success">
                <el-step title="选择交付渠道" description="先决定是本地、魔方云还是财务上下游" />
                <el-step title="选择接口账户" description="自动化商品建议绑定专用账户" />
                <el-step title="选择开通阶段" description="决定是付款后自动开通还是人工审核" />
                <el-step title="确认自动执行" description="开启后订单支付将直接进入任务中心" />
              </el-steps>
            </div>
            <div class="portal-grid portal-grid--three">
              <div v-for="preset in automationPresets" :key="preset.key" class="panel-card">
                <div class="section-card__head">
                  <strong>{{ preset.title }}</strong>
                  <el-button size="small" type="primary" plain @click="applyPreset(preset.key)">使用这套</el-button>
                </div>
                <div class="section-card__meta">{{ preset.description }}</div>
              </div>
            </div>
            <div class="panel-card" style="margin-top: 16px">
              <div class="filter-bar filter-bar--compact">
                <el-select v-model="product.automationConfig.channel">
                  <el-option label="本地交付" value="LOCAL" />
                  <el-option label="魔方云" value="MOFANG_CLOUD" />
                  <el-option label="财务上下游" value="ZJMF_API" />
                  <el-option label="资源池" value="RESOURCE" />
                  <el-option label="手工资源" value="MANUAL" />
                </el-select>
                <el-select v-model="product.automationConfig.moduleType">
                  <el-option label="普通模块" value="NORMAL" />
                  <el-option label="云产品模块" value="CLOUD" />
                  <el-option label="DCIM" value="DCIM" />
                  <el-option label="魔方云" value="DCIMCLOUD" />
                </el-select>
              </div>
              <div class="filter-bar filter-bar--compact">
                <el-select
                  v-model="product.automationConfig.providerAccountId"
                  clearable
                  :disabled="!automationGuide.showAccount"
                  :placeholder="automationGuide.accountPlaceholder"
                >
                  <el-option
                    v-for="account in automationAccountOptions"
                    :key="account.id"
                    :label="`${safeAccountName(account.id)} / ${account.baseUrl}`"
                    :value="account.id"
                  />
                </el-select>
                <el-select v-model="product.automationConfig.provisionStage">
                  <el-option label="付款后自动开通" value="AFTER_PAYMENT" />
                  <el-option label="人工审核后开通" value="MANUAL_REVIEW" />
                  <el-option label="客户确认后开通" value="AFTER_CONFIRM" />
                </el-select>
              </div>
              <div class="filter-bar filter-bar--compact">
                <div class="summary-pill">
                  <span>自动开通</span>
                  <el-switch v-model="product.automationConfig.autoProvision" />
                </div>
                <el-button plain @click="showAdvancedAutomation = !showAdvancedAutomation">
                  {{ showAdvancedAutomation ? "收起高级项" : "展开高级项" }}
                </el-button>
              </div>
              <div class="summary-strip" style="margin-top: 12px">
                <div class="summary-pill"><span>当前阶段</span><strong>{{ formatProvisionStage(product.automationConfig.provisionStage) }}</strong></div>
                <div class="summary-pill"><span>账户要求</span><strong>{{ automationGuide.accountRequired ? "必须绑定接口账户" : "无需绑定接口账户" }}</strong></div>
                <div class="summary-pill"><span>执行方式</span><strong>{{ product.automationConfig.autoProvision ? "系统自动执行" : "等待人工处理" }}</strong></div>
              </div>
              <el-alert
                style="margin-top: 12px"
                type="info"
                :closable="false"
                show-icon
                title="当前配置提示"
                :description="automationGuide.tipText"
              />
              <div v-if="showAdvancedAutomation" class="filter-bar filter-bar--compact">
                <el-input v-model="product.automationConfig.serverGroup" placeholder="节点组 / 服务器组" />
                <el-input v-model="product.automationConfig.providerNode" placeholder="指定节点" />
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t.upstream" name="upstream">
            <div class="portal-grid portal-grid--two">
              <div class="panel-card">
                <div class="section-card__head"><strong>上游接入配置</strong></div>
                <div class="filter-bar filter-bar--compact">
                  <el-select v-model="product.upstreamMapping.providerType">
                    <el-option label="不使用上游" value="NONE" />
                    <el-option label="财务上下游" value="ZJMF_API" />
                    <el-option label="魔方云" value="MOFANG_CLOUD" />
                  </el-select>
                  <el-select v-model="product.upstreamMapping.providerAccountId" clearable placeholder="选择接口账户">
                    <el-option
                      v-for="account in product.upstreamMapping.providerType === 'MOFANG_CLOUD' ? mofangAccounts : financeAccounts"
                      :key="account.id"
                      :label="`${safeAccountName(account.id)} / ${account.baseUrl}`"
                      :value="account.id"
                    />
                  </el-select>
                </div>
                <div class="filter-bar filter-bar--compact">
                  <el-input v-model="product.upstreamMapping.sourceName" placeholder="来源名称" />
                  <el-select v-model="product.upstreamMapping.pricePolicy">
                    <el-option label="跟随上游" value="FOLLOW_UPSTREAM" />
                    <el-option label="本地覆盖" value="LOCAL_OVERRIDE" />
                    <el-option label="按比例加价" value="PERCENTAGE" />
                  </el-select>
                </div>
                <div v-if="product.upstreamMapping.providerType === 'ZJMF_API'" class="filter-bar filter-bar--compact">
                  <el-select v-model="product.upstreamMapping.remoteProductCode" filterable clearable placeholder="上游商品目录">
                    <el-option
                      v-for="item in remoteCatalogOptions"
                      :key="item.remoteProductCode"
                      :label="item.label"
                      :value="item.remoteProductCode"
                    />
                  </el-select>
                  <el-input v-model="product.upstreamMapping.remoteProductName" placeholder="远端商品名称" />
                </div>
                <div v-else class="filter-bar filter-bar--compact">
                  <el-input v-model="product.upstreamMapping.remoteProductCode" placeholder="远端商品编码" />
                  <el-input v-model="product.upstreamMapping.remoteProductName" placeholder="远端商品名称" />
                </div>
                <div class="filter-bar filter-bar--compact">
                  <div class="summary-pill"><span>同步价格</span><el-switch v-model="product.upstreamMapping.autoSyncPricing" /></div>
                  <div class="summary-pill"><span>同步配置项</span><el-switch v-model="product.upstreamMapping.autoSyncConfig" /></div>
                  <div class="summary-pill"><span>同步云模板</span><el-switch v-model="product.upstreamMapping.autoSyncTemplate" /></div>
                </div>
                <el-alert
                  type="info"
                  :closable="false"
                  show-icon
                  :title="`当前策略：${pricePolicyLabel(product.upstreamMapping.pricePolicy)}`"
                  :description="product.upstreamMapping.syncMessage || '保存后会记录映射关系，点击一键同步会按当前开关拉取远端模板。'"
                />
              </div>

              <div class="panel-card">
                <div class="section-card__head"><strong>远端模板预览</strong></div>
                <div v-loading="previewLoading">
                  <el-empty v-if="!remoteTemplate" description="当前没有可预览的远端模板" />
                  <template v-else>
                    <div class="summary-strip">
                      <div class="summary-pill"><span>账户</span><strong>{{ safeAccountName(product.upstreamMapping.providerAccountId) }}</strong></div>
                      <div class="summary-pill"><span>上游商品</span><strong>{{ remoteTemplate.name }}</strong></div>
                    </div>
                    <el-descriptions :column="2" border>
                      <el-descriptions-item label="上游分组">{{ remoteTemplate.groupName }}</el-descriptions-item>
                      <el-descriptions-item label="远端编码">{{ remoteTemplate.remoteProductCode }}</el-descriptions-item>
                      <el-descriptions-item label="价格周期">{{ remoteTemplate.pricing.length }}</el-descriptions-item>
                      <el-descriptions-item label="配置项">{{ remoteTemplate.configOptions.length }}</el-descriptions-item>
                    </el-descriptions>
                  </template>
                </div>
              </div>
            </div>

            <div v-if="remoteTemplate" class="panel-card" style="margin-top: 16px">
              <div class="section-card__head">
                <strong>差异处理动作</strong>
                <span class="section-card__meta">先把上游模板回填到当前草稿，再决定是否保存到正式商品。</span>
              </div>
              <div class="summary-strip">
                <div class="summary-pill"><span>价格差异</span><strong>{{ upstreamDiffSummary.pricingChanged }}</strong></div>
                <div class="summary-pill"><span>配置差异</span><strong>{{ upstreamDiffSummary.configChanged }}</strong></div>
                <div class="summary-pill"><span>模板差异</span><strong>{{ upstreamDiffSummary.templateChanged }}</strong></div>
                <div class="summary-pill"><span>远端价格周期</span><strong>{{ upstreamDiffSummary.remotePricingCount }}</strong></div>
                <div class="summary-pill"><span>远端配置项</span><strong>{{ upstreamDiffSummary.remoteConfigCount }}</strong></div>
              </div>
              <div class="inline-actions" style="margin-top: 16px">
                <el-button plain :loading="previewLoading" @click="loadRemotePreview">刷新远端模板</el-button>
                <el-button plain :disabled="product.upstreamMapping.providerType === 'NONE'" @click="openUpstreamAccountWorkbench">
                  上游账户
                </el-button>
                <el-button plain :disabled="!product.upstreamMapping.providerAccountId" @click="openUpstreamSyncHistoryWorkbench">
                  同步记录
                </el-button>
                <el-button plain @click="applyRemotePricing">回填价格</el-button>
                <el-button plain @click="applyRemoteConfigOptions">回填配置项</el-button>
                <el-button plain @click="applyRemoteResourceTemplate">回填云模板</el-button>
                <el-button type="primary" plain @click="applyRemoteAll">全部回填</el-button>
                <el-button type="primary" :loading="saving" @click="handleSave">保存商品</el-button>
              </div>
            </div>
            <div class="portal-grid portal-grid--three" style="margin-top: 16px">
              <div class="panel-card">
                <div class="section-card__head"><strong>价格差异</strong></div>
                <el-table :data="pricingDiffRows" border stripe empty-text="暂无差异数据">
                  <el-table-column prop="cycleName" label="周期" min-width="120" />
                  <el-table-column label="本地价格" min-width="110">
                    <template #default="{ row }">{{ row.localPrice === null ? "-" : `¥${Number(row.localPrice).toFixed(2)}` }}</template>
                  </el-table-column>
                  <el-table-column label="远端价格" min-width="110">
                    <template #default="{ row }">{{ row.remotePrice === null ? "-" : `¥${Number(row.remotePrice).toFixed(2)}` }}</template>
                  </el-table-column>
                </el-table>
              </div>
              <div class="panel-card">
                <div class="section-card__head"><strong>配置差异</strong></div>
                <el-table :data="configDiffRows" border stripe empty-text="暂无差异数据">
                  <el-table-column prop="name" label="配置项" min-width="120" />
                  <el-table-column prop="localDefault" label="本地默认值" min-width="120" />
                  <el-table-column prop="remoteDefault" label="远端默认值" min-width="120" />
                </el-table>
              </div>
              <div class="panel-card">
                <div class="section-card__head"><strong>云模板差异</strong></div>
                <el-table :data="templateDiffRows" border stripe empty-text="暂无差异数据">
                  <el-table-column prop="label" label="字段" min-width="120" />
                  <el-table-column prop="local" label="本地" min-width="120" />
                  <el-table-column prop="remote" label="远端推导" min-width="120" />
                </el-table>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t.pricing" name="pricing">
            <div class="panel-card">
              <div class="section-card__head">
                <strong>价格矩阵</strong>
                <div class="inline-actions">
                  <el-button size="small" plain @click="addCommonPricingRows">补齐常用周期</el-button>
                  <el-button size="small" plain @click="sortPricingRows">排序周期</el-button>
                  <el-button type="primary" plain size="small" @click="addPricingRow">新增周期</el-button>
                </div>
              </div>
              <div class="summary-strip" style="margin-bottom: 16px">
                <div class="summary-pill"><span>周期总数</span><strong>{{ pricingSummary.total }}</strong></div>
                <div class="summary-pill"><span>可收费周期</span><strong>{{ pricingSummary.chargeable }}</strong></div>
                <div class="summary-pill"><span>免费周期</span><strong>{{ pricingSummary.free }}</strong></div>
                <div class="summary-pill"><span>最高售价</span><strong>¥{{ pricingSummary.highestPrice.toFixed(2) }}</strong></div>
              </div>
              <el-table :data="product.pricing" border stripe>
                <el-table-column label="周期编码" min-width="140">
                  <template #default="{ row }"><el-input v-model="row.cycleCode" /></template>
                </el-table-column>
                <el-table-column label="周期名称" min-width="140">
                  <template #default="{ row }"><el-input v-model="row.cycleName" /></template>
                </el-table-column>
                <el-table-column label="销售价" min-width="120">
                  <template #default="{ row }"><el-input-number v-model="row.price" :min="0" style="width: 100%" /></template>
                </el-table-column>
                <el-table-column label="开户费" min-width="120">
                  <template #default="{ row }"><el-input-number v-model="row.setupFee" :min="0" style="width: 100%" /></template>
                </el-table-column>
                <el-table-column label="操作" min-width="140" fixed="right">
                  <template #default="{ $index }">
                    <div class="inline-actions">
                      <el-button link type="primary" @click="duplicatePricingRow($index)">复制</el-button>
                      <el-button link type="danger" @click="removePricingRow($index)">删除</el-button>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t.config" name="config">
            <div class="summary-strip" style="margin-bottom: 16px">
              <div class="summary-pill"><span>配置项总数</span><strong>{{ configSummary.total }}</strong></div>
              <div class="summary-pill"><span>必填项</span><strong>{{ configSummary.required }}</strong></div>
              <div class="summary-pill"><span>选择类</span><strong>{{ configSummary.selectLike }}</strong></div>
              <div class="summary-pill"><span>选项总数</span><strong>{{ configSummary.totalChoices }}</strong></div>
            </div>
            <div class="table-toolbar">
              <div class="table-toolbar__meta">
                <strong>配置项编辑台</strong>
                <span>快速维护下拉、单选、文本类配置，并控制默认值和加价策略</span>
              </div>
              <div class="inline-actions">
                <el-button plain @click="addSelectConfigOption">新增下拉项</el-button>
                <el-button plain @click="addTextConfigOption">新增文本项</el-button>
                <el-button type="primary" plain @click="addConfigOption">新增配置项</el-button>
              </div>
            </div>
            <div class="portal-grid" style="gap: 16px">
              <div v-for="(option, optionIndex) in product.configOptions" :key="`${option.code}-${optionIndex}`" class="panel-card">
                <div class="section-card__head">
                  <strong>{{ option.name || `配置项 ${optionIndex + 1}` }}</strong>
                  <div class="inline-actions">
                    <el-button size="small" plain @click="duplicateConfigOption(optionIndex)">复制配置项</el-button>
                    <el-button size="small" type="primary" plain @click="addChoice(option)">新增选项</el-button>
                    <el-button size="small" link type="danger" @click="removeConfigOption(optionIndex)">删除配置项</el-button>
                  </div>
                </div>
                <div class="summary-strip" style="margin-bottom: 12px">
                  <div class="summary-pill"><span>编码</span><strong>{{ option.code || "-" }}</strong></div>
                  <div class="summary-pill"><span>类型</span><strong>{{ option.inputType }}</strong></div>
                  <div class="summary-pill"><span>默认值</span><strong>{{ option.defaultValue || "-" }}</strong></div>
                  <div class="summary-pill"><span>选项数</span><strong>{{ option.choices.length }}</strong></div>
                </div>
                <div class="filter-bar filter-bar--compact">
                  <el-input v-model="option.code" placeholder="编码" />
                  <el-input v-model="option.name" placeholder="名称" />
                  <el-select v-model="option.inputType">
                    <el-option label="下拉选择" value="select" />
                    <el-option label="单选" value="radio" />
                    <el-option label="文本" value="text" />
                  </el-select>
                </div>
                <div class="filter-bar filter-bar--compact">
                  <el-input v-model="option.defaultValue" placeholder="默认值" />
                  <div class="summary-pill"><span>必填</span><el-switch v-model="option.required" /></div>
                </div>
                <el-input v-model="option.description" type="textarea" :rows="2" placeholder="说明" />
                <el-table :data="option.choices" border stripe empty-text="当前还没有选项" style="margin-top: 12px">
                  <el-table-column label="值" min-width="120">
                    <template #default="{ row }"><el-input v-model="(row as ProductConfigChoice).value" /></template>
                  </el-table-column>
                  <el-table-column label="显示名称" min-width="140">
                    <template #default="{ row }"><el-input v-model="(row as ProductConfigChoice).label" /></template>
                  </el-table-column>
                  <el-table-column label="加价" min-width="120">
                    <template #default="{ row }"><el-input-number v-model="(row as ProductConfigChoice).priceDelta" :min="0" style="width: 100%" /></template>
                  </el-table-column>
                  <el-table-column label="操作" min-width="90" fixed="right">
                    <template #default="{ $index }"><el-button link type="danger" @click="removeChoice(option, $index)">删除</el-button></template>
                  </el-table-column>
                </el-table>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t.template" name="template">
            <div class="panel-card">
              <div class="section-card__head"><strong>云资源模板</strong></div>
              <div class="filter-bar filter-bar--compact">
                <el-input v-model="product.resourceTemplate.regionName" placeholder="地域" />
                <el-input v-model="product.resourceTemplate.zoneName" placeholder="可用区" />
                <el-input v-model="product.resourceTemplate.operatingSystem" placeholder="系统镜像" />
              </div>
              <div class="filter-bar filter-bar--compact">
                <el-input-number v-model="product.resourceTemplate.cpuCores" :min="1" style="width: 100%" />
                <el-input-number v-model="product.resourceTemplate.memoryGB" :min="1" style="width: 100%" />
                <el-input-number v-model="product.resourceTemplate.bandwidthMbps" :min="1" style="width: 100%" />
              </div>
              <div class="filter-bar filter-bar--compact">
                <el-input-number v-model="product.resourceTemplate.systemDiskGB" :min="20" style="width: 100%" />
                <el-input-number v-model="product.resourceTemplate.dataDiskGB" :min="0" style="width: 100%" />
                <el-input-number v-model="product.resourceTemplate.publicIpCount" :min="0" style="width: 100%" />
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t.audit" name="audit">
            <div class="panel-card">
              <el-table :data="auditLogs" border stripe empty-text="当前商品暂无审计记录">
                <el-table-column prop="createdAt" label="时间" min-width="180" />
                <el-table-column prop="actor" label="操作人" min-width="140" />
                <el-table-column prop="action" label="动作" min-width="180" />
                <el-table-column prop="description" label="说明" min-width="260" />
              </el-table>
            </div>
          </el-tab-pane>
        </el-tabs>
      </template>
    </PageWorkbench>
  </div>
</template>
