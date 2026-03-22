<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import ContextTabs from "@/components/workbench/ContextTabs.vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createCustomerContact,
  deleteCustomerContact,
  fetchCustomerAuditLogs,
  fetchCustomerContacts,
  fetchCustomerDetail,
  fetchCustomerIdentitiesDetail,
  fetchCustomerInvoices,
  fetchCustomerServices,
  fetchCustomerTickets,
  reviewCustomerIdentity,
  updateCustomer,
  updateCustomerContact,
  type AuditLog,
  type Contact,
  type Customer,
  type IdentityRecord,
  type RelatedItem
} from "@/api/admin";

type WorkbenchRelatedItem = RelatedItem & { entityId?: number };

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const contactDialogVisible = ref(false);
const contactSubmitting = ref(false);
const savingCustomer = ref(false);
const editingContactId = ref<number | null>(null);

const customer = reactive<Customer>({
  id: 0,
  customerNo: "",
  name: "",
  email: "",
  mobile: "",
  type: "COMPANY",
  status: "ACTIVE",
  groupName: "",
  levelName: "",
  salesOwner: "",
  remarks: ""
});

const contacts = ref<Contact[]>([]);
const identities = ref<IdentityRecord[]>([]);
const services = ref<WorkbenchRelatedItem[]>([]);
const invoices = ref<WorkbenchRelatedItem[]>([]);
const tickets = ref<WorkbenchRelatedItem[]>([]);
const auditLogs = ref<AuditLog[]>([]);

const contactForm = reactive({
  name: "",
  email: "",
  mobile: "",
  roleName: "",
  isPrimary: false
});

const customerId = computed(() => route.params.id as string);
const approvedIdentityCount = computed(
  () => identities.value.filter(item => item.verifyStatus === "APPROVED").length
);
const pendingIdentityCount = computed(
  () => identities.value.filter(item => item.verifyStatus === "PENDING").length
);
const activeServiceCount = computed(() => services.value.filter(item => item.status === "ACTIVE").length);
const suspendedServiceCount = computed(() => services.value.filter(item => item.status === "SUSPENDED").length);
const mofangServiceCount = computed(() => services.value.filter(item => item.providerType === "MOFANG_CLOUD").length);
const serviceWithPublicIpCount = computed(
  () => services.value.filter(item => item.ipAddress && item.ipAddress !== "-").length
);
const syncedCustomer = computed(() => customer.remarks.includes("魔方云同步"));
const primaryContact = computed(() => contacts.value.find(item => item.isPrimary) ?? contacts.value[0] ?? null);
const latestService = computed(() => services.value[0] ?? null);
const latestInvoice = computed(() => invoices.value[0] ?? null);
const latestTicket = computed(() => tickets.value[0] ?? null);

const contextItems = computed(() => [
  {
    key: "customer",
    label: "客户工作台",
    active: true,
    badge: customer.customerNo || "-"
  },
  {
    key: "identity",
    label: "实名审核",
    to: "/customer/identities",
    badge: identities.value.length
  },
  {
    key: "service",
    label: "业务列表",
    to: "/services/list",
    badge: services.value.length
  },
  {
    key: "invoice",
    label: "账单管理",
    to: "/billing/invoices",
    badge: invoices.value.length
  },
  {
    key: "order",
    label: "产品订单",
    to: "/orders/list"
  }
]);

const serviceRegionSummary = computed(() =>
  summarizeLabels(services.value.map(item => item.regionName), 6)
);
const providerSummary = computed(() =>
  summarizeLabels(services.value.map(item => formatProviderTypeLabel(item.providerType)), 4)
);
const serviceStatusSummary = computed(() =>
  summarizeLabels(services.value.map(item => formatServiceStatusLabel(item.status)), 4)
);
const customerSignals = computed(() => {
  const result: Array<{
    type: "success" | "info" | "warning";
    title: string;
    description: string;
  }> = [];

  if (syncedCustomer.value) {
    result.push({
      type: "info",
      title: "当前客户来自魔方云同步",
      description: "这类客户会优先带入服务、资源和实例关系，联系人、实名和账单通常需要后续按运营流程补齐。"
    });
  }

  if (!primaryContact.value) {
    result.push({
      type: "warning",
      title: "还没有主联系人",
      description: "建议至少补充 1 个主联系人，方便后续客服通知、欠费提醒和工单协作。"
    });
  }

  if (approvedIdentityCount.value === 0) {
    result.push({
      type: "warning",
      title: "未完成实名",
      description: "当前没有已通过的实名档案，如需按合规流程交付资源，建议先补充实名认证资料。"
    });
  } else if (pendingIdentityCount.value > 0) {
    result.push({
      type: "info",
      title: "存在待审核实名资料",
      description: `当前还有 ${pendingIdentityCount.value} 份实名资料待审核，可以优先处理。`
    });
  }

  if (services.value.length > 0 && invoices.value.length === 0) {
    result.push({
      type: "info",
      title: "当前没有本地账单",
      description: "说明这位客户目前以同步服务为主，财务对象还没有在本地系统完整沉淀。"
    });
  }

  if (services.value.length > 0 && activeServiceCount.value === 0) {
    result.push({
      type: "warning",
      title: "当前没有运行中的服务",
      description: "客户虽然有关联服务，但当前都不在运行态，建议核对是否已暂停、已终止或待恢复。"
    });
  }

  if (result.length === 0) {
    result.push({
      type: "success",
      title: "客户资料状态良好",
      description: "联系人、服务和关联数据已经具备继续跟进和运营处理的基础条件。"
    });
  }

  return result;
});
const recentActivities = computed(() => {
  const result: Array<{
    key: string;
    label: string;
    value: string;
    meta: string;
    actionLabel?: string;
    action?: () => void;
  }> = [];

  if (latestService.value) {
    result.push({
      key: "service",
      label: "最近服务",
      value: `${latestService.value.serviceNo} / ${latestService.value.productName}`,
      meta: `${formatServiceStatusLabel(latestService.value.status)} · ${latestService.value.regionName || "未同步地域"}`,
      actionLabel: latestService.value.entityId ? "打开服务" : undefined,
      action: latestService.value.entityId ? () => openServiceDetail(latestService.value!) : undefined
    });
  }

  if (latestInvoice.value) {
    result.push({
      key: "invoice",
      label: "最近账单",
      value: `${latestInvoice.value.invoiceNo} / ${formatAmount(latestInvoice.value.totalAmount)}`,
      meta: `${formatInvoiceStatusLabel(latestInvoice.value.status)} · ${latestInvoice.value.dueAt || "无到期时间"}`,
      actionLabel: latestInvoice.value.entityId ? "打开账单" : undefined,
      action: latestInvoice.value.entityId ? () => openInvoiceDetail(latestInvoice.value!) : undefined
    });
  }

  if (latestTicket.value) {
    result.push({
      key: "ticket",
      label: "最近工单",
      value: `${latestTicket.value.ticketNo} / ${latestTicket.value.title}`,
      meta: `${formatTicketStatusLabel(latestTicket.value.status)} · ${latestTicket.value.updatedAt || "无更新时间"}`
    });
  }

  if (auditLogs.value[0]) {
    result.push({
      key: "audit",
      label: "最近操作",
      value: auditLogs.value[0].description || auditLogs.value[0].action,
      meta: `${auditLogs.value[0].actor || "系统"} · ${auditLogs.value[0].createdAt || "未知时间"}`
    });
  }

  return result;
});

function pickEntityId(item: RelatedItem, keys: string[]) {
  for (const key of keys) {
    const value = Number((item as Record<string, unknown>)[key]);
    if (Number.isFinite(value) && value > 0) return value;
  }
  return undefined;
}

function summarizeLabels(values: Array<string | undefined>, limit = 5) {
  const counter = new Map<string, number>();
  for (const item of values) {
    const label = (item || "").trim();
    if (!label || label === "-") continue;
    counter.set(label, (counter.get(label) ?? 0) + 1);
  }
  return Array.from(counter.entries())
    .map(([label, count]) => ({ label, count }))
    .sort((left, right) => {
      if (right.count !== left.count) return right.count - left.count;
      return left.label.localeCompare(right.label, "zh-CN");
    })
    .slice(0, limit);
}

function resetContactForm() {
  editingContactId.value = null;
  contactForm.name = "";
  contactForm.email = "";
  contactForm.mobile = "";
  contactForm.roleName = "";
  contactForm.isPrimary = false;
}

function normalizeServices(items: RelatedItem[]) {
  return items.map(item => ({
    ...item,
    entityId: pickEntityId(item, ["serviceId", "id"]),
    serviceNo: item.serviceNo ?? item.no ?? "-",
    productName: item.productName ?? item.name ?? "-",
    nextDueAt: item.nextDueAt ?? item.dueAt ?? "-",
    providerType: item.providerType ?? "-",
    regionName: item.regionName ?? "-",
    ipAddress: item.ipAddress ?? "-",
    description: item.description ?? "-"
  }));
}

function normalizeInvoices(items: RelatedItem[]) {
  return items.map(item => ({
    ...item,
    entityId: pickEntityId(item, ["invoiceId", "id"]),
    invoiceNo: item.invoiceNo ?? item.no ?? "-",
    totalAmount: item.totalAmount ?? item.amount ?? "-",
    dueAt: item.dueAt ?? "-"
  }));
}

function normalizeTickets(items: RelatedItem[]) {
  return items.map(item => ({
    ...item,
    entityId: pickEntityId(item, ["ticketId", "id"]),
    ticketNo: item.ticketNo ?? item.no ?? "-",
    title: item.title ?? item.name ?? "-",
    updatedAt: item.updatedAt ?? "-"
  }));
}

function formatCustomerType(type: string) {
  const mapping: Record<string, string> = {
    COMPANY: "企业客户",
    PERSONAL: "个人客户"
  };
  return mapping[type] ?? type;
}

function formatCustomerStatus(status: string) {
  const mapping: Record<string, string> = {
    ACTIVE: "正常",
    SUSPENDED: "暂停",
    DISABLED: "禁用"
  };
  return mapping[status] ?? status;
}

function formatProviderTypeLabel(type?: string) {
  const mapping: Record<string, string> = {
    MOFANG_CLOUD: "魔方云",
    ZJMF_API: "上下游财务",
    RESOURCE: "资源池",
    MANUAL: "手动资源",
    LOCAL: "本地模块"
  };
  return mapping[type ?? ""] ?? type ?? "-";
}

function formatServiceStatusLabel(status?: string) {
  const mapping: Record<string, string> = {
    ACTIVE: "运行中",
    PENDING: "待开通",
    SUSPENDED: "已暂停",
    TERMINATED: "已终止",
    PROVISIONING: "处理中",
    FAILED: "异常"
  };
  return mapping[status ?? ""] ?? status ?? "-";
}

function serviceStatusTag(status?: string) {
  const mapping: Record<string, string> = {
    ACTIVE: "success",
    PENDING: "warning",
    SUSPENDED: "info",
    TERMINATED: "danger",
    PROVISIONING: "primary",
    FAILED: "danger"
  };
  return mapping[status ?? ""] ?? "info";
}

function formatInvoiceStatusLabel(status?: string) {
  const mapping: Record<string, string> = {
    PAID: "已支付",
    UNPAID: "未支付",
    OVERDUE: "已逾期",
    CANCELLED: "已取消",
    REFUNDED: "已退款",
    DRAFT: "草稿"
  };
  return mapping[status ?? ""] ?? status ?? "-";
}

function invoiceStatusTag(status?: string) {
  const mapping: Record<string, string> = {
    PAID: "success",
    UNPAID: "warning",
    OVERDUE: "danger",
    CANCELLED: "info",
    REFUNDED: "info",
    DRAFT: "info"
  };
  return mapping[status ?? ""] ?? "info";
}

function formatTicketStatusLabel(status?: string) {
  const mapping: Record<string, string> = {
    OPEN: "待处理",
    IN_PROGRESS: "处理中",
    REPLY: "待回复",
    CLOSED: "已关闭",
    RESOLVED: "已解决"
  };
  return mapping[status ?? ""] ?? status ?? "-";
}

function ticketStatusTag(status?: string) {
  const mapping: Record<string, string> = {
    OPEN: "warning",
    IN_PROGRESS: "primary",
    REPLY: "warning",
    CLOSED: "info",
    RESOLVED: "success"
  };
  return mapping[status ?? ""] ?? "info";
}

function customerStatusTag(status: string) {
  const mapping: Record<string, string> = {
    ACTIVE: "success",
    SUSPENDED: "warning",
    DISABLED: "danger"
  };
  return mapping[status] ?? "info";
}

function identityStatusLabel(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "待审核",
    APPROVED: "已实名",
    REJECTED: "已驳回"
  };
  return mapping[status] ?? status;
}

function identityTagType(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "warning",
    APPROVED: "success",
    REJECTED: "danger"
  };
  return mapping[status] ?? "info";
}

function formatAmount(value: string | number | undefined) {
  if (value === undefined || value === null || value === "") return "-";
  const amount = typeof value === "number" ? value : Number(value);
  if (!Number.isFinite(amount)) return String(value);
  return `¥${amount.toFixed(2)}`;
}

function openServiceDetail(row: WorkbenchRelatedItem) {
  if (!row.entityId) return;
  void router.push(`/services/detail/${row.entityId}`);
}

function openInvoiceDetail(row: WorkbenchRelatedItem) {
  if (!row.entityId) return;
  void router.push(`/billing/invoices/${row.entityId}`);
}

function openLatestService() {
  if (!latestService.value) return;
  openServiceDetail(latestService.value);
}

async function loadDetail() {
  loading.value = true;
  try {
    const id = customerId.value;
    const [customerData, contactData, identityData, serviceData, invoiceData, ticketData, auditData] =
      await Promise.all([
        fetchCustomerDetail(id),
        fetchCustomerContacts(id),
        fetchCustomerIdentitiesDetail(id),
        fetchCustomerServices(id),
        fetchCustomerInvoices(id),
        fetchCustomerTickets(id),
        fetchCustomerAuditLogs(id)
      ]);

    Object.assign(customer, customerData);
    contacts.value = contactData;
    identities.value = identityData;
    services.value = normalizeServices(serviceData);
    invoices.value = normalizeInvoices(invoiceData);
    tickets.value = normalizeTickets(ticketData);
    auditLogs.value = auditData;
  } finally {
    loading.value = false;
  }
}

async function saveCustomer() {
  savingCustomer.value = true;
  try {
    const updated = await updateCustomer(customer.id, customer);
    Object.assign(customer, updated);
    ElMessage.success("客户资料已更新");
    await loadDetail();
  } finally {
    savingCustomer.value = false;
  }
}

function openCreateContact() {
  resetContactForm();
  contactDialogVisible.value = true;
}

function openEditContact(contact: Contact) {
  editingContactId.value = contact.id;
  contactForm.name = contact.name;
  contactForm.email = contact.email;
  contactForm.mobile = contact.mobile;
  contactForm.roleName = contact.roleName;
  contactForm.isPrimary = contact.isPrimary;
  contactDialogVisible.value = true;
}

async function submitContact() {
  contactSubmitting.value = true;
  try {
    if (editingContactId.value) {
      await updateCustomerContact(customer.id, editingContactId.value, contactForm);
      ElMessage.success("联系人已更新");
    } else {
      await createCustomerContact(customer.id, contactForm);
      ElMessage.success("联系人已新增");
    }
    contactDialogVisible.value = false;
    resetContactForm();
    await loadDetail();
  } finally {
    contactSubmitting.value = false;
  }
}

async function removeContact(contact: Contact) {
  await ElMessageBox.confirm(`确认删除联系人 ${contact.name} 吗？`, "删除联系人", {
    type: "warning"
  });
  await deleteCustomerContact(customer.id, contact.id);
  ElMessage.success("联系人已删除");
  await loadDetail();
}

async function handleReview(identity: IdentityRecord, status: "APPROVED" | "REJECTED") {
  const remark =
    status === "APPROVED" ? "资料齐全，审核通过" : "资料不完整，请补充后重新提交";
  await reviewCustomerIdentity(customer.id, identity.id, { status, remark });
  ElMessage.success(status === "APPROVED" ? "已通过实名认证" : "已驳回实名认证");
  await loadDetail();
}

watch(
  () => route.params.id,
  () => {
    void loadDetail();
  }
);

onMounted(() => {
  void loadDetail();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="客户"
      title="客户工作台"
      subtitle="把客户资料、联系人、实名、服务、账单、工单和审计记录集中到同一页处理，作为后台客户主工作台。"
    >
      <template #actions>
        <el-button @click="router.push('/customer/list')">返回列表</el-button>
        <el-button @click="loadDetail">刷新详情</el-button>
        <el-button type="primary" :loading="savingCustomer" @click="saveCustomer">
          保存客户资料
        </el-button>
      </template>

      <template #context>
        <ContextTabs :items="contextItems" />
      </template>

      <template #metrics>
        <div class="detail-kpi-grid">
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户编号</span>
            <strong class="detail-kpi-card__value">{{ customer.customerNo || "N/A" }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户状态</span>
            <el-tag :type="customerStatusTag(customer.status)" effect="light">
              {{ formatCustomerStatus(customer.status) }}
            </el-tag>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户分组</span>
            <strong class="detail-kpi-card__value">{{ customer.groupName || "未设置" }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户等级</span>
            <strong class="detail-kpi-card__value">{{ customer.levelName || "未设置" }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">已实名档案</span>
            <strong class="detail-kpi-card__value">{{ approvedIdentityCount }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">关联服务</span>
            <strong class="detail-kpi-card__value">{{ services.length }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">公网 IP 服务</span>
            <strong class="detail-kpi-card__value">{{ serviceWithPublicIpCount }}</strong>
          </div>
        </div>
      </template>

      <el-tabs>
        <el-tab-pane label="客户概览">
          <div class="signal-grid">
            <el-alert
              v-for="item in customerSignals"
              :key="item.title"
              :type="item.type"
              :closable="false"
              show-icon
              :title="item.title"
              :description="item.description"
            />
          </div>

          <div class="detail-kpi-grid" style="margin-bottom: 16px">
            <div class="detail-kpi-card">
              <span class="detail-kpi-card__label">关联服务</span>
              <strong class="detail-kpi-card__value">{{ services.length }}</strong>
            </div>
            <div class="detail-kpi-card">
              <span class="detail-kpi-card__label">运行中</span>
              <strong class="detail-kpi-card__value">{{ activeServiceCount }}</strong>
            </div>
            <div class="detail-kpi-card">
              <span class="detail-kpi-card__label">已暂停</span>
              <strong class="detail-kpi-card__value">{{ suspendedServiceCount }}</strong>
            </div>
            <div class="detail-kpi-card">
              <span class="detail-kpi-card__label">魔方云服务</span>
              <strong class="detail-kpi-card__value">{{ mofangServiceCount }}</strong>
            </div>
            <div class="detail-kpi-card">
              <span class="detail-kpi-card__label">公网 IP 覆盖</span>
              <strong class="detail-kpi-card__value">{{ serviceWithPublicIpCount }}</strong>
            </div>
          </div>

          <div class="customer-grid">
            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>基础资料</strong>
                <span class="section-card__meta">{{ formatCustomerType(customer.type) }}</span>
              </div>
              <el-form label-position="top" class="customer-form-grid">
                <el-form-item label="客户名称">
                  <el-input v-model="customer.name" />
                </el-form-item>
                <el-form-item label="邮箱">
                  <el-input v-model="customer.email" />
                </el-form-item>
                <el-form-item label="手机号">
                  <el-input v-model="customer.mobile" />
                </el-form-item>
                <el-form-item label="销售归属">
                  <el-input v-model="customer.salesOwner" />
                </el-form-item>
                <el-form-item label="客户分组">
                  <el-input v-model="customer.groupName" />
                </el-form-item>
                <el-form-item label="客户等级">
                  <el-input v-model="customer.levelName" />
                </el-form-item>
                <el-form-item label="备注" class="customer-form-grid__full">
                  <el-input v-model="customer.remarks" type="textarea" :rows="4" />
                </el-form-item>
              </el-form>
            </div>

            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>客户摘要</strong>
                <span class="section-card__meta">对象关系总览</span>
              </div>
              <el-descriptions :column="1" border>
                <el-descriptions-item label="主联系人">
                  {{
                    primaryContact
                      ? `${primaryContact.name} / ${primaryContact.mobile || primaryContact.email || "未填写联系方式"}`
                      : "未设置"
                  }}
                </el-descriptions-item>
                <el-descriptions-item label="客户类型">
                  {{ formatCustomerType(customer.type) }}
                </el-descriptions-item>
                <el-descriptions-item label="联系人数量">
                  {{ contacts.length }}
                </el-descriptions-item>
                <el-descriptions-item label="实名档案">
                  {{ identities.length }}
                </el-descriptions-item>
                <el-descriptions-item label="关联账单">
                  {{ invoices.length }}
                </el-descriptions-item>
                <el-descriptions-item label="关联工单">
                  {{ tickets.length }}
                </el-descriptions-item>
                <el-descriptions-item label="审计记录">
                  {{ auditLogs.length }}
                </el-descriptions-item>
              </el-descriptions>

              <div class="quick-actions-grid">
                <el-button type="primary" :disabled="!latestService?.entityId" @click="openLatestService">
                  打开最近服务
                </el-button>
                <el-button @click="openCreateContact">补充联系人</el-button>
                <el-button @click="router.push('/services/list')">查看服务列表</el-button>
                <el-button :disabled="!latestInvoice?.entityId" @click="latestInvoice && openInvoiceDetail(latestInvoice)">
                  打开最近账单
                </el-button>
              </div>
            </div>
          </div>

          <div class="overview-grid">
            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>服务画像</strong>
                <span class="section-card__meta">地域、渠道和状态分布</span>
              </div>

              <div class="overview-section">
                <span class="overview-section__label">主要地域</span>
                <div class="summary-chip-list">
                  <span v-for="item in serviceRegionSummary" :key="`region-${item.label}`" class="summary-chip">
                    {{ item.label }} · {{ item.count }}
                  </span>
                  <span v-if="!serviceRegionSummary.length" class="summary-chip summary-chip--muted">暂无地域数据</span>
                </div>
              </div>

              <div class="overview-section">
                <span class="overview-section__label">渠道分布</span>
                <div class="summary-chip-list">
                  <span v-for="item in providerSummary" :key="`provider-${item.label}`" class="summary-chip">
                    {{ item.label }} · {{ item.count }}
                  </span>
                  <span v-if="!providerSummary.length" class="summary-chip summary-chip--muted">暂无渠道数据</span>
                </div>
              </div>

              <div class="overview-section">
                <span class="overview-section__label">服务状态</span>
                <div class="summary-chip-list">
                  <span v-for="item in serviceStatusSummary" :key="`status-${item.label}`" class="summary-chip">
                    {{ item.label }} · {{ item.count }}
                  </span>
                  <span v-if="!serviceStatusSummary.length" class="summary-chip summary-chip--muted">暂无状态数据</span>
                </div>
              </div>

              <el-table :data="services.slice(0, 5)" border stripe empty-text="暂无服务记录">
                <el-table-column prop="serviceNo" label="服务编号" min-width="160" />
                <el-table-column prop="productName" label="产品" min-width="180" />
                <el-table-column label="状态" min-width="110">
                  <template #default="{ row }">
                    <el-tag :type="serviceStatusTag(row.status)" effect="light">
                      {{ formatServiceStatusLabel(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="ipAddress" label="IP" min-width="150" />
                <el-table-column label="操作" min-width="120" fixed="right">
                  <template #default="{ row }">
                    <el-button type="primary" link :disabled="!row.entityId" @click="openServiceDetail(row)">
                      打开
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>

            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>最近业务</strong>
                <span class="section-card__meta">最近服务、账单、工单和操作</span>
              </div>

              <div v-if="recentActivities.length" class="activity-list">
                <div v-for="item in recentActivities" :key="item.key" class="activity-card">
                  <div class="activity-card__head">
                    <strong>{{ item.label }}</strong>
                    <el-button v-if="item.action && item.actionLabel" type="primary" link @click="item.action()">
                      {{ item.actionLabel }}
                    </el-button>
                  </div>
                  <div class="activity-card__value">{{ item.value }}</div>
                  <div class="activity-card__meta">{{ item.meta }}</div>
                </div>
              </div>
              <el-empty v-else description="当前没有最近业务记录" :image-size="88" />
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="联系人">
          <div class="page-header">
            <div>
              <h2 class="section-title">联系人管理</h2>
              <p class="page-subtitle">支持新增、编辑、删除和设置主联系人。</p>
            </div>
            <el-button type="primary" @click="openCreateContact">新增联系人</el-button>
          </div>
          <el-table :data="contacts" border stripe empty-text="暂无联系人">
            <el-table-column prop="name" label="姓名" min-width="140" />
            <el-table-column prop="email" label="邮箱" min-width="180" />
            <el-table-column prop="mobile" label="手机号" min-width="140" />
            <el-table-column prop="roleName" label="角色" min-width="160" />
            <el-table-column label="主联系人" min-width="120">
              <template #default="{ row }">
                <el-tag :type="row.isPrimary ? 'primary' : 'info'" effect="plain">
                  {{ row.isPrimary ? "是" : "否" }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" min-width="160" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openEditContact(row)">编辑</el-button>
                <el-button link type="danger" @click="removeContact(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="实名信息">
          <el-table :data="identities" border stripe empty-text="暂无实名档案">
            <el-table-column prop="subjectName" label="主体名称" min-width="220" />
            <el-table-column prop="certNo" label="证件号" min-width="220" />
            <el-table-column prop="identityType" label="证件类型" min-width="120" />
            <el-table-column prop="countryCode" label="国家/地区" min-width="120" />
            <el-table-column prop="submittedAt" label="提交时间" min-width="180" />
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <el-tag :type="identityTagType(row.verifyStatus)" effect="light">
                  {{ identityStatusLabel(row.verifyStatus) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="reviewRemark" label="审核备注" min-width="220" />
            <el-table-column label="操作" min-width="180" fixed="right">
              <template #default="{ row }">
                <el-button
                  link
                  type="primary"
                  :disabled="row.verifyStatus === 'APPROVED'"
                  @click="handleReview(row, 'APPROVED')"
                >
                  通过
                </el-button>
                <el-button
                  link
                  type="danger"
                  :disabled="row.verifyStatus === 'REJECTED'"
                  @click="handleReview(row, 'REJECTED')"
                >
                  驳回
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="关联服务">
          <el-table :data="services" border stripe empty-text="暂无关联服务">
            <el-table-column prop="serviceNo" label="服务编号" min-width="180" />
            <el-table-column prop="productName" label="产品名称" min-width="220" />
            <el-table-column label="渠道" min-width="130">
              <template #default="{ row }">{{ formatProviderTypeLabel(row.providerType) }}</template>
            </el-table-column>
            <el-table-column prop="regionName" label="地域" min-width="150" />
            <el-table-column prop="ipAddress" label="公网 IP" min-width="160" />
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <el-tag :type="serviceStatusTag(row.status)" effect="light">
                  {{ formatServiceStatusLabel(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="nextDueAt" label="到期时间" min-width="180" />
            <el-table-column prop="description" label="同步说明" min-width="220" show-overflow-tooltip />
            <el-table-column label="操作" min-width="140" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link :disabled="!row.entityId" @click="openServiceDetail(row)">
                  打开服务
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="关联账单">
          <el-table :data="invoices" border stripe empty-text="暂无关联账单">
            <el-table-column prop="invoiceNo" label="账单编号" min-width="180" />
            <el-table-column label="金额" min-width="120">
              <template #default="{ row }">{{ formatAmount(row.totalAmount) }}</template>
            </el-table-column>
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <el-tag :type="invoiceStatusTag(row.status)" effect="light">
                  {{ formatInvoiceStatusLabel(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="dueAt" label="到期时间" min-width="180" />
            <el-table-column label="操作" min-width="140" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link :disabled="!row.entityId" @click="openInvoiceDetail(row)">
                  打开账单
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="关联工单">
          <el-table :data="tickets" border stripe empty-text="暂无关联工单">
            <el-table-column prop="ticketNo" label="工单编号" min-width="180" />
            <el-table-column prop="title" label="工单标题" min-width="260" />
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <el-tag :type="ticketStatusTag(row.status)" effect="light">
                  {{ formatTicketStatusLabel(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="updatedAt" label="更新时间" min-width="180" />
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="审计日志">
          <el-table :data="auditLogs" border stripe empty-text="暂无审计记录">
            <el-table-column prop="createdAt" label="时间" min-width="180" />
            <el-table-column prop="actor" label="操作人" min-width="140" />
            <el-table-column prop="action" label="动作" min-width="180" />
            <el-table-column prop="description" label="说明" min-width="280" />
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <el-dialog
        v-model="contactDialogVisible"
        :title="editingContactId ? '编辑联系人' : '新增联系人'"
        width="520px"
      >
        <el-form label-position="top">
          <el-form-item label="姓名">
            <el-input v-model="contactForm.name" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="contactForm.email" />
          </el-form-item>
          <el-form-item label="手机号">
            <el-input v-model="contactForm.mobile" />
          </el-form-item>
          <el-form-item label="角色">
            <el-input v-model="contactForm.roleName" placeholder="例如：技术联系人 / 财务联系人" />
          </el-form-item>
          <el-form-item label="主联系人">
            <el-switch v-model="contactForm.isPrimary" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="contactDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="contactSubmitting" @click="submitContact">
            保存联系人
          </el-button>
        </template>
      </el-dialog>
    </PageWorkbench>
  </div>
</template>

<style scoped>
.customer-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(320px, 0.65fr);
  gap: 16px;
}

.overview-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(320px, 0.9fr);
  gap: 16px;
  margin-top: 16px;
}

.inner-card {
  padding: 18px;
}

.customer-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px 16px;
}

.customer-form-grid__full {
  grid-column: 1 / -1;
}

.section-card__head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
}

.section-card__meta {
  color: #5e7093;
  font-size: 12px;
}

.signal-grid {
  display: grid;
  gap: 12px;
  margin-bottom: 16px;
}

.quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 16px;
}

.overview-section + .overview-section {
  margin-top: 16px;
}

.overview-section__label {
  display: block;
  margin-bottom: 10px;
  color: #41506b;
  font-size: 13px;
  font-weight: 600;
}

.summary-chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.summary-chip {
  display: inline-flex;
  align-items: center;
  min-height: 30px;
  padding: 0 10px;
  border-radius: 999px;
  background: #eef3ff;
  color: #3056a8;
  font-size: 12px;
  font-weight: 600;
}

.summary-chip--muted {
  background: #f5f7fa;
  color: #7b8798;
}

.activity-list {
  display: grid;
  gap: 12px;
}

.activity-card {
  padding: 14px;
  border: 1px solid #e5e9f2;
  border-radius: 14px;
  background: #fbfcff;
}

.activity-card__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.activity-card__value {
  color: #1f2a3d;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.5;
}

.activity-card__meta {
  margin-top: 6px;
  color: #6b778c;
  font-size: 12px;
}

@media (max-width: 1280px) {
  .customer-grid {
    grid-template-columns: 1fr;
  }

  .overview-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .customer-form-grid {
    grid-template-columns: 1fr;
  }

  .quick-actions-grid {
    grid-template-columns: 1fr;
  }
}
</style>
