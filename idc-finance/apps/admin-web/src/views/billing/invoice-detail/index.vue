<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import AutomationTaskPanel from "@/components/workbench/AutomationTaskPanel.vue";
import AuditTrailTable from "@/components/workbench/AuditTrailTable.vue";
import ContextTabs from "@/components/workbench/ContextTabs.vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { useLocaleStore } from "@/store";
import {
  fetchInvoiceDetail,
  fetchProviderAccounts,
  receiveInvoicePayment,
  refundInvoice,
  updateUnpaidInvoice,
  type InvoiceDetailResponse,
  type ProviderAccount
} from "@/api/admin";
import {
  billingCycleOptions as createBillingCycleOptions,
  formatBillingCycle,
  formatChangeOrderAction,
  formatChangeOrderExecution,
  formatInvoiceStatus as formatInvoiceStatusLabel,
  formatMoney,
  formatOrderStatus,
  formatPaymentChannel,
  formatProviderType,
  formatServiceStatus,
  invoiceStatusOptions as createInvoiceStatusOptions
} from "@/utils/business";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();

const loading = ref(false);
const receiving = ref(false);
const refunding = ref(false);
const savingEdit = ref(false);
const receiveDialogVisible = ref(false);
const editDialogVisible = ref(false);
const detail = ref<InvoiceDetailResponse | null>(null);
const providerAccounts = ref<ProviderAccount[]>([]);
const billingCycles = computed(() => createBillingCycleOptions(localeStore.locale));
const invoiceStatusOptions = computed(() => createInvoiceStatusOptions(localeStore.locale));
const _legacyBillingCycles = [
  { label: "月付", value: "monthly" },
  { label: "季付", value: "quarterly" },
  { label: "半年付", value: "semiannual" },
  { label: "年付", value: "annual" },
  { label: "两年付", value: "biennially" },
  { label: "三年付", value: "triennially" },
  { label: "一次性", value: "onetime" }
];
const _legacyInvoiceStatusOptions = [
  { label: "未支付", value: "UNPAID" },
  { label: "已支付", value: "PAID" },
  { label: "已退款", value: "REFUNDED" }
];

const receiveForm = reactive({
  channel: "OFFLINE",
  tradeNo: ""
});
const editForm = reactive({
  productName: "",
  billingCycle: "monthly",
  amount: 0,
  dueAt: "",
  status: "UNPAID",
  reason: ""
});

const customerId = computed(() => detail.value?.invoice.customerId ?? detail.value?.order?.customerId ?? 0);
const invoiceAccount = computed(() => {
  const accountId = detail.value?.order?.providerAccountId || detail.value?.services[0]?.providerAccountId || 0;
  if (!accountId) return null;
  return providerAccounts.value.find(item => item.id === accountId) ?? null;
});
const invoiceEditImpact = computed(() => {
  switch (editForm.status) {
    case "PAID":
      return {
        type: "success" as const,
        title: "将账单改为已支付",
        description: "系统会回写订单为生效状态，并尝试恢复或激活关联服务。"
      };
    case "REFUNDED":
      return {
        type: "warning" as const,
        title: "将账单改为已退款",
        description: "系统会写入退款记录、回退订单，并终止关联服务。"
      };
    default:
      return {
        type: "info" as const,
        title: "将账单改为未支付",
        description: "系统会把订单回退为待处理，并将关联服务改为挂起状态。"
      };
  }
});

const contextTabs = computed(() => [
  { key: "customer", label: "客户", to: customerId.value ? `/customer/detail/${customerId.value}` : undefined },
  { key: "invoice", label: "账单工作台", active: true, badge: detail.value?.invoice.invoiceNo },
  {
    key: "order",
    label: "订单",
    to: detail.value?.order ? `/orders/detail/${detail.value.order.id}` : undefined,
    badge: detail.value?.order?.orderNo
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
    onetime: "一次性"
  };
  return mapping[cycle] ?? cycle;
}

function formatInvoiceStatus(status: string) {
  return formatInvoiceStatusLabel(localeStore.locale, status);
  const mapping: Record<string, string> = {
    UNPAID: "未支付",
    PAID: "已支付",
    REFUNDED: "已退款"
  };
  return mapping[status] ?? status;
}

function invoiceStatusType(status: string) {
  const mapping: Record<string, string> = {
    UNPAID: "warning",
    PAID: "success",
    REFUNDED: "info"
  };
  return mapping[status] ?? "info";
}

function paymentChannelLabel(channel: string) {
  return formatPaymentChannel(localeStore.locale, channel);
  const mapping: Record<string, string> = {
    OFFLINE: "线下汇款",
    ALIPAY: "支付宝",
    WECHAT: "微信支付",
    BALANCE: "余额抵扣"
  };
  return mapping[channel] ?? channel;
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
    detail.value = await fetchInvoiceDetail(route.params.id as string);
  } finally {
    loading.value = false;
  }
}

function openEditDialog() {
  if (!detail.value) return;
  editForm.productName = detail.value.invoice.productName;
  editForm.billingCycle = detail.value.invoice.billingCycle;
  editForm.amount = detail.value.invoice.totalAmount;
  editForm.dueAt = detail.value.invoice.dueAt;
  editForm.status = detail.value.invoice.status;
  editForm.reason = "";
  editDialogVisible.value = true;
}

async function submitInvoiceEdit() {
  if (!detail.value) return;
  savingEdit.value = true;
  try {
    const result = await updateUnpaidInvoice(detail.value.invoice.id, {
      productName: editForm.productName,
      billingCycle: editForm.billingCycle,
      amount: editForm.amount,
      dueAt: editForm.dueAt,
      status: editForm.status,
      reason: editForm.reason
    });
    detail.value = result;
    editDialogVisible.value = false;
    ElMessage.success("账单和关联订单已更新，变更记录已写入审计");
  } finally {
    savingEdit.value = false;
  }
}

async function submitReceivePayment() {
  if (!detail.value) return;
  receiving.value = true;
  try {
    const result = await receiveInvoicePayment(detail.value.invoice.id, {
      channel: receiveForm.channel,
      tradeNo: receiveForm.tradeNo
    });
    receiveDialogVisible.value = false;
    receiveForm.tradeNo = "";
    ElMessage.success(
      result.service
        ? `收款成功，支付单 ${result.payment.paymentNo} 已生成，服务 ${result.service.serviceNo} 已联动处理`
        : `收款成功，支付单 ${result.payment.paymentNo} 已生成`
    );
    await loadDetail();
  } finally {
    receiving.value = false;
  }
}

async function handleRefund() {
  if (!detail.value) return;
  let reason = "后台人工退款";
  try {
    const result = await ElMessageBox.prompt("请输入退款原因", "执行退款", {
      confirmButtonText: "确认退款",
      cancelButtonText: "取消",
      inputValue: reason
    });
    reason = result.value;
  } catch {
    return;
  }

  refunding.value = true;
  try {
    const result = await refundInvoice(detail.value.invoice.id, reason);
    ElMessage.success(
      result.service
        ? `退款完成，退款单 ${result.refund.refundNo} 已生成，服务 ${result.service.serviceNo} 已联动处理`
        : `退款完成，退款单 ${result.refund.refundNo} 已生成`
    );
    await loadDetail();
  } finally {
    refunding.value = false;
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
      eyebrow="财务 / 账单"
      title="账单工作台"
      subtitle="集中处理收款、退款、支付流水、接口账户与关联服务联动。"
    >
      <template #actions>
        <el-button @click="router.push('/billing/invoices')">返回列表</el-button>
        <el-button @click="loadDetail">刷新详情</el-button>
        <el-button v-if="detail" type="primary" plain @click="openEditDialog">
          人工调整账单
        </el-button>
        <el-button v-if="detail?.invoice.status === 'UNPAID'" type="primary" @click="receiveDialogVisible = true">
          登记收款
        </el-button>
        <el-button v-if="detail?.invoice.status === 'PAID'" type="danger" :loading="refunding" @click="handleRefund">
          执行退款
        </el-button>
      </template>

      <template #context>
        <ContextTabs v-if="detail" :items="contextTabs" />
      </template>

      <template #metrics>
        <div v-if="detail" class="detail-kpi-grid">
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">账单编号</span>
            <strong class="detail-kpi-card__value">{{ detail.invoice.invoiceNo }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">账单金额</span>
            <strong class="detail-kpi-card__value">{{ formatCurrency(detail.invoice.totalAmount) }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">账单状态</span>
            <el-tag :type="invoiceStatusType(detail.invoice.status)" effect="light">
              {{ formatInvoiceStatus(detail.invoice.status) }}
            </el-tag>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">到期时间</span>
            <strong class="detail-kpi-card__value">{{ detail.invoice.dueAt }}</strong>
          </div>
        </div>
      </template>

      <template v-if="detail">
        <el-tabs>
          <el-tab-pane label="账单概览">
            <div class="portal-grid portal-grid--two">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>账单信息</strong>
                  <span class="section-card__meta">{{ formatCycle(detail.invoice.billingCycle) }}</span>
                </div>
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="订单编号">{{ detail.invoice.orderNo }}</el-descriptions-item>
                  <el-descriptions-item label="商品名称">{{ detail.invoice.productName }}</el-descriptions-item>
                  <el-descriptions-item label="支付时间">{{ detail.invoice.paidAt || "-" }}</el-descriptions-item>
                  <el-descriptions-item label="计费周期">{{ formatCycle(detail.invoice.billingCycle) }}</el-descriptions-item>
                  <el-descriptions-item label="接口账户" :span="2">
                    {{ invoiceAccount ? `${invoiceAccount.name} / ${invoiceAccount.baseUrl}` : "未绑定" }}
                  </el-descriptions-item>
                </el-descriptions>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>处理摘要</strong>
                  <span class="section-card__meta">收款、退款与服务联动</span>
                </div>
                <div class="summary-strip">
                  <div class="summary-pill">
                    <span>支付记录</span>
                    <strong>{{ detail.payments.length }}</strong>
                  </div>
                  <div class="summary-pill">
                    <span>退款记录</span>
                    <strong>{{ detail.refunds.length }}</strong>
                  </div>
                  <div class="summary-pill">
                    <span>关联服务</span>
                    <strong>{{ detail.services.length }}</strong>
                  </div>
                </div>
                <el-alert
                  title="账单工作台负责承接支付、退款和交付回写。绑定的接口账户决定支付完成后服务会落到哪套自动化渠道。"
                  type="info"
                  :closable="false"
                  show-icon
                />
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="关联订单">
            <el-descriptions v-if="detail.order" :column="2" border>
              <el-descriptions-item label="订单编号">{{ detail.order.orderNo }}</el-descriptions-item>
              <el-descriptions-item label="客户名称">{{ detail.order.customerName }}</el-descriptions-item>
              <el-descriptions-item label="商品名称">{{ detail.order.productName }}</el-descriptions-item>
              <el-descriptions-item label="订单状态">{{ formatOrderStatus(localeStore.locale, detail.order.status) }}</el-descriptions-item>
              <el-descriptions-item label="自动化渠道">{{ formatProviderType(localeStore.locale, detail.order.automationType) }}</el-descriptions-item>
              <el-descriptions-item label="接口账户">
                {{ invoiceAccount ? `${invoiceAccount.name} / ${invoiceAccount.baseUrl}` : "未绑定" }}
              </el-descriptions-item>
              <el-descriptions-item label="订单金额">{{ formatCurrency(detail.order.amount) }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ detail.order.createdAt }}</el-descriptions-item>
            </el-descriptions>
            <div class="inline-actions" style="margin-top: 16px">
              <el-button v-if="detail.order" type="primary" link @click="router.push(`/orders/detail/${detail.order.id}`)">
                打开订单详情
              </el-button>
            </div>
          </el-tab-pane>

          <el-tab-pane v-if="detail.changeOrder" label="改配记录">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="改配订单">{{ detail.changeOrder.orderNo }}</el-descriptions-item>
              <el-descriptions-item label="改配账单">{{ detail.changeOrder.invoiceNo }}</el-descriptions-item>
              <el-descriptions-item label="资源动作">
                {{ changeOrderActionLabel(detail.changeOrder.actionName) }}
              </el-descriptions-item>
              <el-descriptions-item label="支付状态">{{ formatInvoiceStatus(detail.changeOrder.status) }}</el-descriptions-item>
              <el-descriptions-item label="执行状态">
                <el-tag :type="changeOrderExecutionType(detail.changeOrder.executionStatus)" effect="light">
                  {{ changeOrderExecutionLabel(detail.changeOrder.executionStatus) }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="改配金额">{{ formatCurrency(detail.changeOrder.amount) }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ detail.changeOrder.createdAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="已支付时间">{{ detail.changeOrder.paidAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="已退款时间">{{ detail.changeOrder.refundedAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="变更原因" :span="2">{{ detail.changeOrder.reason || "-" }}</el-descriptions-item>
              <el-descriptions-item label="执行结果" :span="2">
                {{ detail.changeOrder.executionMessage || "-" }}
              </el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane label="支付记录">
            <el-table :data="detail.payments" border stripe empty-text="暂无支付记录">
              <el-table-column prop="paymentNo" label="支付单号" min-width="160" />
              <el-table-column label="渠道" min-width="120">
                <template #default="{ row }">{{ paymentChannelLabel(row.channel) }}</template>
              </el-table-column>
              <el-table-column prop="tradeNo" label="交易号" min-width="180" />
              <el-table-column label="金额" min-width="120">
                <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="operator" label="操作人" min-width="120" />
              <el-table-column prop="paidAt" label="支付时间" min-width="180" />
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="退款记录">
            <el-table :data="detail.refunds" border stripe empty-text="暂无退款记录">
              <el-table-column prop="refundNo" label="退款单号" min-width="160" />
              <el-table-column label="金额" min-width="120">
                <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="reason" label="原因" min-width="220" />
              <el-table-column label="状态" min-width="120">
                <template #default="{ row }">{{ formatInvoiceStatus(row.status) }}</template>
              </el-table-column>
              <el-table-column prop="createdAt" label="创建时间" min-width="180" />
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="关联服务">
            <el-table :data="detail.services" border stripe empty-text="暂无关联服务">
              <el-table-column prop="serviceNo" label="服务编号" min-width="180" />
              <el-table-column prop="productName" label="产品名称" min-width="220" />
              <el-table-column label="渠道" min-width="140">
                <template #default="{ row }">{{ formatProviderType(localeStore.locale, row.providerType) }}</template>
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
              <el-table-column label="操作" min-width="140" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="router.push(`/services/detail/${row.id}`)">打开服务</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="自动化任务">
            <AutomationTaskPanel title="账单自动化任务" :invoice-id="detail.invoice.id" />
          </el-tab-pane>

          <el-tab-pane label="审计日志">
            <AuditTrailTable :items="detail.auditLogs" empty-text="暂无审计记录" />
          </el-tab-pane>
        </el-tabs>
      </template>
    </PageWorkbench>

    <el-dialog v-model="receiveDialogVisible" title="登记收款" width="420px">
      <el-form label-position="top">
        <el-form-item label="收款渠道">
          <el-select v-model="receiveForm.channel">
            <el-option label="线下汇款" value="OFFLINE" />
            <el-option label="支付宝" value="ALIPAY" />
            <el-option label="微信支付" value="WECHAT" />
            <el-option label="余额抵扣" value="BALANCE" />
          </el-select>
        </el-form-item>
        <el-form-item label="交易号">
          <el-input v-model="receiveForm.tradeNo" placeholder="可留空" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="receiveDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="receiving" @click="submitReceivePayment">确认收款</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="editDialogVisible" title="人工调整账单" width="560px">
      <el-form label-position="top">
        <el-alert
          :title="invoiceEditImpact.title"
          :description="invoiceEditImpact.description"
          :type="invoiceEditImpact.type"
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
        <el-form-item label="账单金额">
          <el-input-number v-model="editForm.amount" :min="0" :precision="2" :step="10" style="width: 100%" />
        </el-form-item>
        <el-form-item label="账单状态">
          <el-select v-model="editForm.status" style="width: 100%">
            <el-option v-for="item in invoiceStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="到期时间">
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
            placeholder="例如：财务复核、历史账单修正、销售方案调整"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingEdit" @click="submitInvoiceEdit">保存调整</el-button>
      </template>
    </el-dialog>
  </div>
</template>
