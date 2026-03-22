<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import AutomationTaskPanel from "@/components/workbench/AutomationTaskPanel.vue";
import AuditTrailTable from "@/components/workbench/AuditTrailTable.vue";
import ContextTabs from "@/components/workbench/ContextTabs.vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { useLocaleStore } from "@/store";
import {
  fetchOrderDetail,
  fetchProviderAccounts,
  updatePendingOrder,
  type OrderDetailResponse,
  type ProviderAccount
} from "@/api/admin";
import {
  billingCycleOptions as createBillingCycleOptions,
  formatBillingCycle,
  formatChangeOrderAction,
  formatChangeOrderExecution,
  formatInvoiceStatus,
  formatMoney,
  formatOrderStatus,
  formatProductType,
  formatProviderType,
  formatServiceStatus,
  orderStatusOptions as createOrderStatusOptions
} from "@/utils/business";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();
const billingCycles = computed(() => createBillingCycleOptions(localeStore.locale));
const orderStatusOptions = computed(() => createOrderStatusOptions(localeStore.locale));
const _legacyBillingCycles = [
  { label: "月付", value: "monthly" },
  { label: "季付", value: "quarterly" },
  { label: "半年付", value: "semiannual" },
  { label: "年付", value: "annual" },
  { label: "两年付", value: "biennially" },
  { label: "三年付", value: "triennially" },
  { label: "一次性", value: "onetime" }
];
const _legacyOrderStatusOptions = [
  { label: "待支付", value: "PENDING" },
  { label: "待交付", value: "ACTIVE" },
  { label: "已完成", value: "COMPLETED" },
  { label: "已取消", value: "CANCELLED" }
];

const loading = ref(false);
const savingEdit = ref(false);
const editDialogVisible = ref(false);
const detail = ref<OrderDetailResponse | null>(null);
const providerAccounts = ref<ProviderAccount[]>([]);
const editForm = reactive({
  productName: "",
  billingCycle: "monthly",
  amount: 0,
  dueAt: "",
  status: "PENDING",
  reason: ""
});

const customerId = computed(() => detail.value?.order.customerId ?? 0);
const orderAccount = computed(() => {
  const accountId = detail.value?.order.providerAccountId || detail.value?.services[0]?.providerAccountId || 0;
  if (!accountId) return null;
  return providerAccounts.value.find(item => item.id === accountId) ?? null;
});
const orderEditImpact = computed(() => {
  switch (editForm.status) {
    case "ACTIVE":
      return {
        type: "success" as const,
        title: "将订单改为待交付",
        description: "系统会把关联服务拉回运行中，并保持账单链路为可交付状态。"
      };
    case "COMPLETED":
      return {
        type: "success" as const,
        title: "将订单改为已完成",
        description: "适合人工确认业务已完成的场景，已关联服务会保持生效状态。"
      };
    case "CANCELLED":
      return {
        type: "warning" as const,
        title: "将订单改为已取消",
        description: "系统会把关联服务终止，适合订单关闭、作废或撤销业务的场景。"
      };
    default:
      return {
        type: "info" as const,
        title: "将订单改为待支付",
        description: "系统会把关联服务回退为挂起状态，等待后续收款或再次人工确认。"
      };
  }
});

const contextTabs = computed(() => [
  { key: "customer", label: "客户", to: customerId.value ? `/customer/detail/${customerId.value}` : undefined },
  { key: "order", label: "订单工作台", active: true, badge: detail.value?.order.orderNo },
  {
    key: "invoice",
    label: "账单",
    to: detail.value?.invoices[0] ? `/billing/invoices/${detail.value.invoices[0].id}` : undefined,
    badge: detail.value?.invoices.length ?? 0
  },
  {
    key: "service",
    label: "服务",
    to: detail.value?.services[0] ? `/services/detail/${detail.value.services[0].id}` : undefined,
    badge: detail.value?.services.length ?? 0
  }
]);

function formatCurrency(value: number) {
  return formatMoney(localeStore.locale, value);
  return `¥${value.toFixed(2)}`;
}

function formatCycle(cycle: string) {
  return formatBillingCycle(localeStore.locale, cycle);
  const mapping: Record<string, string> = {
    monthly: "月付",
    quarterly: "季付",
    semiannual: "半年付",
    semiannually: "半年付",
    annual: "年付",
    annually: "年付",
    biennially: "两年付",
    triennially: "三年付",
    onetime: "一次性"
  };
  return mapping[cycle] ?? cycle;
}

function formatStatus(status: string) {
  return formatOrderStatus(localeStore.locale, status);
  const mapping: Record<string, string> = {
    PENDING: "待支付",
    ACTIVE: "待交付",
    PAID: "已支付",
    COMPLETED: "已完成",
    CANCELLED: "已取消"
  };
  return mapping[status] ?? status;
}

function formatAutomationType(type: string) {
  return formatProviderType(localeStore.locale, type);
  const mapping: Record<string, string> = {
    LOCAL: "本地模块",
    MOFANG_CLOUD: "魔方云",
    ZJMF_API: "上下游财务",
    RESOURCE: "资源池",
    MANUAL: "手动资源"
  };
  return mapping[type] ?? type;
}

function formatProductTypeLabel(type: string) {
  return formatProductType(localeStore.locale, type);
}

function statusTagType(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "warning",
    ACTIVE: "primary",
    PAID: "success",
    COMPLETED: "success",
    CANCELLED: "info"
  };
  return mapping[status] ?? "info";
}

function changeOrderActionLabel(action: string) {
  return formatChangeOrderAction(localeStore.locale, action);
  const mapping: Record<string, string> = {
    "add-ipv4": "新增 IPv4",
    "add-ipv6": "新增 IPv6",
    "add-disk": "新增数据盘",
    "resize-disk": "扩容数据盘",
    "create-snapshot": "创建快照",
    "create-backup": "创建备份"
  };
  return mapping[action] ?? action;
}

function changeOrderExecutionLabel(status: string) {
  return formatChangeOrderExecution(localeStore.locale, status);
  const mapping: Record<string, string> = {
    WAITING_PAYMENT: "待支付",
    PAID: "待回写",
    EXECUTING: "执行中",
    EXECUTED: "已执行",
    EXECUTE_FAILED: "执行失败",
    EXECUTE_BLOCKED: "执行阻塞",
    REFUNDED: "已退款"
  };
  return mapping[status] ?? status;
}

function changeOrderExecutionType(status: string) {
  const mapping: Record<string, string> = {
    WAITING_PAYMENT: "warning",
    PAID: "info",
    EXECUTING: "primary",
    EXECUTED: "success",
    EXECUTE_FAILED: "danger",
    EXECUTE_BLOCKED: "warning",
    REFUNDED: "info"
  };
  return mapping[status] ?? "info";
}

async function loadProviderAccounts() {
  providerAccounts.value = await fetchProviderAccounts();
}

async function loadDetail() {
  loading.value = true;
  try {
    detail.value = await fetchOrderDetail(route.params.id as string);
  } finally {
    loading.value = false;
  }
}

function openEditDialog() {
  if (!detail.value) return;
  editForm.productName = detail.value.order.productName;
  editForm.billingCycle = detail.value.order.billingCycle;
  editForm.amount = detail.value.order.amount;
  editForm.dueAt = detail.value.invoices[0]?.dueAt ?? "";
  editForm.status = detail.value.order.status;
  editForm.reason = "";
  editDialogVisible.value = true;
}

async function submitEditOrder() {
  if (!detail.value) return;
  savingEdit.value = true;
  try {
    const result = await updatePendingOrder(detail.value.order.id, {
      productName: editForm.productName,
      billingCycle: editForm.billingCycle,
      amount: editForm.amount,
      dueAt: editForm.dueAt,
      status: editForm.status,
      reason: editForm.reason
    });
    detail.value = result;
    editDialogVisible.value = false;
    ElMessage.success("订单和关联账单已更新，变更记录已写入审计");
  } finally {
    savingEdit.value = false;
  }
}

watch(() => route.params.id, () => void loadDetail());

onMounted(async () => {
  await Promise.all([loadProviderAccounts(), loadDetail()]);
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="业务 / 订单"
      title="订单工作台"
      subtitle="围绕订单对象集中处理客户、账单、服务交付、接口账户与自动化任务。"
    >
      <template #actions>
        <el-button @click="router.push('/orders/list')">返回列表</el-button>
        <el-button @click="loadDetail">刷新详情</el-button>
        <el-button v-if="detail" type="primary" plain @click="openEditDialog">
          人工调整订单
        </el-button>
      </template>

      <template #context>
        <ContextTabs v-if="detail" :items="contextTabs" />
      </template>

      <template #metrics>
        <div v-if="detail" class="detail-kpi-grid">
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">订单编号</span>
            <strong class="detail-kpi-card__value">{{ detail.order.orderNo }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户</span>
            <strong class="detail-kpi-card__value">{{ detail.order.customerName }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">订单金额</span>
            <strong class="detail-kpi-card__value">{{ formatCurrency(detail.order.amount) }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">订单状态</span>
            <el-tag :type="statusTagType(detail.order.status)" effect="light">
              {{ formatStatus(detail.order.status) }}
            </el-tag>
          </div>
        </div>
      </template>

      <template v-if="detail">
        <el-tabs>
          <el-tab-pane label="订单概览">
            <div class="portal-grid portal-grid--two">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>基础信息</strong>
                  <span class="section-card__meta">{{ formatProductTypeLabel(detail.order.productType) }}</span>
                </div>
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="商品名称">{{ detail.order.productName }}</el-descriptions-item>
                  <el-descriptions-item label="计费周期">{{ formatCycle(detail.order.billingCycle) }}</el-descriptions-item>
                  <el-descriptions-item label="自动化渠道">
                    {{ formatAutomationType(detail.order.automationType) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="接口账户">
                    {{ orderAccount ? `${orderAccount.name} / ${orderAccount.baseUrl}` : "未绑定" }}
                  </el-descriptions-item>
                  <el-descriptions-item label="创建时间">{{ detail.order.createdAt }}</el-descriptions-item>
                  <el-descriptions-item label="订单编号">{{ detail.order.orderNo }}</el-descriptions-item>
                </el-descriptions>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>交付链路</strong>
                  <span class="section-card__meta">订单 -> 账单 -> 服务 -> 自动化任务</span>
                </div>
                <div class="summary-strip">
                  <div class="summary-pill">
                    <span>关联账单</span>
                    <strong>{{ detail.invoices.length }}</strong>
                  </div>
                  <div class="summary-pill">
                    <span>已开通服务</span>
                    <strong>{{ detail.services.length }}</strong>
                  </div>
                  <div class="summary-pill">
                    <span>审计记录</span>
                    <strong>{{ detail.auditLogs.length }}</strong>
                  </div>
                </div>
                <el-alert
                  title="订单工作台负责串联支付、交付与同步状态。当前接口账户会直接决定后续自动化执行落到哪套魔方或上下游环境。"
                  type="info"
                  :closable="false"
                  show-icon
                />
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="关联账单">
            <el-table :data="detail.invoices" border stripe empty-text="暂无关联账单">
              <el-table-column prop="invoiceNo" label="账单编号" min-width="180" />
              <el-table-column label="金额" min-width="120">
                <template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template>
              </el-table-column>
              <el-table-column label="状态" min-width="120">
                <template #default="{ row }">{{ formatInvoiceStatus(localeStore.locale, row.status) }}</template>
              </el-table-column>
              <el-table-column prop="dueAt" label="到期时间" min-width="180" />
              <el-table-column label="操作" min-width="140" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="router.push(`/billing/invoices/${row.id}`)">打开账单</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="已开通服务">
            <el-table :data="detail.services" border stripe empty-text="暂无关联服务">
              <el-table-column prop="serviceNo" label="服务编号" min-width="170" />
              <el-table-column prop="productName" label="产品名称" min-width="220" />
              <el-table-column label="渠道" min-width="140">
                <template #default="{ row }">{{ formatAutomationType(row.providerType) }}</template>
              </el-table-column>
              <el-table-column label="接口账户" min-width="220">
                <template #default="{ row }">
                  {{
                    providerAccounts.find(item => item.id === row.providerAccountId)
                      ? `${providerAccounts.find(item => item.id === row.providerAccountId)?.name} / ${providerAccounts.find(item => item.id === row.providerAccountId)?.baseUrl}`
                      : "未绑定"
                  }}
                </template>
              </el-table-column>
              <el-table-column label="状态" min-width="120">
                <template #default="{ row }">{{ formatServiceStatus(localeStore.locale, row.status) }}</template>
              </el-table-column>
              <el-table-column prop="createdAt" label="开通时间" min-width="180" />
              <el-table-column prop="nextDueAt" label="下次到期" min-width="160" />
              <el-table-column label="操作" min-width="140" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="router.push(`/services/detail/${row.id}`)">打开服务</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane v-if="detail.changeOrder" label="改配记录">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="改配订单">{{ detail.changeOrder.orderNo }}</el-descriptions-item>
              <el-descriptions-item label="改配账单">{{ detail.changeOrder.invoiceNo }}</el-descriptions-item>
              <el-descriptions-item label="资源动作">
                {{ changeOrderActionLabel(detail.changeOrder.actionName) }}
              </el-descriptions-item>
              <el-descriptions-item label="支付状态">{{ formatInvoiceStatus(localeStore.locale, detail.changeOrder.status) }}</el-descriptions-item>
              <el-descriptions-item label="执行状态">
                <el-tag :type="changeOrderExecutionType(detail.changeOrder.executionStatus)" effect="light">
                  {{ changeOrderExecutionLabel(detail.changeOrder.executionStatus) }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="改配金额">{{ formatCurrency(detail.changeOrder.amount) }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ detail.changeOrder.createdAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="已支付时间">{{ detail.changeOrder.paidAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="执行结果" :span="2">
                {{ detail.changeOrder.executionMessage || "-" }}
              </el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane label="自动化任务">
            <AutomationTaskPanel title="订单自动化任务" :order-id="detail.order.id" />
          </el-tab-pane>

          <el-tab-pane label="审计日志">
            <AuditTrailTable :items="detail.auditLogs" empty-text="暂无审计记录" />
          </el-tab-pane>
        </el-tabs>
      </template>
    </PageWorkbench>

    <el-dialog v-model="editDialogVisible" title="人工调整订单" width="560px">
      <el-form label-position="top">
        <el-alert
          :title="orderEditImpact.title"
          :description="orderEditImpact.description"
          :type="orderEditImpact.type"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-form-item label="商品名称">
          <el-input v-model="editForm.productName" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="计费周期">
          <el-select v-model="editForm.billingCycle" style="width: 100%">
            <el-option v-for="item in billingCycles" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="待支付金额">
          <el-input-number v-model="editForm.amount" :min="0" :precision="2" :step="10" style="width: 100%" />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="editForm.status" style="width: 100%">
            <el-option v-for="item in orderStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="账单到期时间">
          <el-date-picker
            v-model="editForm.dueAt"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="请选择到期时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="变更原因">
          <el-input
            v-model="editForm.reason"
            type="textarea"
            :rows="3"
            placeholder="例如：客户申请改周期、财务复核改价、销售补贴调整"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingEdit" @click="submitEditOrder">保存调整</el-button>
      </template>
    </el-dialog>
  </div>
</template>
