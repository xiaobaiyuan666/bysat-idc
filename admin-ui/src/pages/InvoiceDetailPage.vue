<script setup lang="ts">
import { ArrowLeft, Money } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";

import { http } from "@/api/http";
import { useAuthStore } from "@/stores/auth";
import { formatCurrency, formatDate, formatDateTime } from "@/utils/format";
import {
  billingJobStatusMap,
  billingJobTypeMap,
  callbackStatusMap,
  getLabel,
  getStatusTagType,
  invoiceStatusMap,
  invoiceTypeMap,
  paymentMethodMap,
  paymentStatusMap,
  refundStatusMap,
  roleMap,
} from "@/utils/maps";

const authStore = useAuthStore();
const route = useRoute();
const router = useRouter();

const loading = ref(false);
const actionSubmitting = ref(false);
const paymentSubmitting = ref(false);
const refundSubmitting = ref(false);
const actionDialogVisible = ref(false);
const paymentDialogVisible = ref(false);
const refundDialogVisible = ref(false);
const detail = ref<any | null>(null);

const actionForm = reactive({
  action: "issue",
  dueDate: "",
  reason: "",
});

const paymentForm = reactive({
  method: "ALIPAY",
  amount: 0,
  transactionNo: "",
});

const refundForm = reactive({
  paymentId: "",
  paymentNo: "",
  method: "ALIPAY",
  amount: 0,
  refundableAmount: 0,
  reason: "",
  transactionNo: "",
});

const canManageInvoices = computed(() =>
  Boolean(authStore.user?.permissions.includes("invoices.manage")),
);

const canManagePayments = computed(() =>
  Boolean(authStore.user?.permissions.includes("payments.manage")),
);

const paymentRows = computed(() =>
  (detail.value?.invoice?.payments ?? []).map((payment: any) => {
    const refundedAmount = (payment.refunds ?? []).reduce(
      (total: number, refund: any) => total + (refund.amount ?? 0),
      0,
    );

    return {
      ...payment,
      refundedAmount,
      refundableAmount: Math.max((payment.amount ?? 0) - refundedAmount, 0),
    };
  }),
);

const financeHint = computed(() => {
  const invoice = detail.value?.invoice;
  if (!invoice) return "-";
  if (invoice.status === "DRAFT") return "账单仍是草稿，需要先签发后才能进入收款流程。";
  if (invoice.status === "VOID") return "账单已作废，不能再登记收款。";
  if ((detail.value?.summary?.outstandingAmount ?? 0) > 0) return "账单仍有未结清金额，可继续登记收款或核对支付回调。";
  return "账单已结清，可根据支付记录处理退款或核对回调日志。";
});

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get(`/invoices/${route.params.id}`);
    detail.value = data.data;
  } finally {
    loading.value = false;
  }
}

function openIssueDialog() {
  actionForm.action = "issue";
  actionForm.dueDate = detail.value?.invoice?.dueDate ? String(detail.value.invoice.dueDate).slice(0, 10) : "";
  actionForm.reason = "";
  actionDialogVisible.value = true;
}

function openVoidDialog() {
  actionForm.action = "void";
  actionForm.dueDate = "";
  actionForm.reason = "";
  actionDialogVisible.value = true;
}

async function submitAction() {
  actionSubmitting.value = true;
  try {
    await http.post(`/invoices/${route.params.id}/action`, actionForm);
    ElMessage.success(actionForm.action === "issue" ? "账单已签发" : "账单已作废");
    actionDialogVisible.value = false;
    await loadData();
  } finally {
    actionSubmitting.value = false;
  }
}

function openPaymentDialog() {
  paymentForm.method = "ALIPAY";
  paymentForm.amount = Number((((detail.value?.summary?.outstandingAmount ?? 0) / 100)).toFixed(2));
  paymentForm.transactionNo = "";
  paymentDialogVisible.value = true;
}

async function submitPayment() {
  paymentSubmitting.value = true;
  try {
    await http.post("/payments", {
      invoiceId: detail.value.invoice.id,
      method: paymentForm.method,
      amount: paymentForm.amount,
      transactionNo: paymentForm.transactionNo,
    });
    ElMessage.success("收款登记成功");
    paymentDialogVisible.value = false;
    await loadData();
  } finally {
    paymentSubmitting.value = false;
  }
}

function openRefundDialog(payment: any) {
  refundForm.paymentId = payment.id;
  refundForm.paymentNo = payment.paymentNo;
  refundForm.method = payment.method;
  refundForm.refundableAmount = Number(((payment.refundableAmount ?? 0) / 100).toFixed(2));
  refundForm.amount = refundForm.refundableAmount;
  refundForm.reason = "";
  refundForm.transactionNo = "";
  refundDialogVisible.value = true;
}

async function submitRefund() {
  refundSubmitting.value = true;
  try {
    await http.post(`/payments/${refundForm.paymentId}/refund`, {
      amount: refundForm.amount,
      reason: refundForm.reason,
      method: refundForm.method,
      transactionNo: refundForm.transactionNo,
    });
    ElMessage.success("退款已创建");
    refundDialogVisible.value = false;
    await loadData();
  } finally {
    refundSubmitting.value = false;
  }
}

onMounted(async () => {
  await authStore.ensureUser();
  await loadData();
});

watch(() => route.params.id, () => {
  void loadData();
});
</script>

<template>
  <div class="page-body" v-loading="loading">
    <div class="toolbar-row">
      <div class="toolbar-actions">
        <el-button :icon="ArrowLeft" @click="router.push('/invoices')">返回账单中心</el-button>
        <el-button v-if="canManageInvoices && detail?.invoice?.status === 'DRAFT'" type="primary" @click="openIssueDialog">签发账单</el-button>
        <el-button v-if="canManageInvoices && detail?.invoice?.status !== 'VOID' && detail?.invoice?.paidAmount === 0" type="danger" plain @click="openVoidDialog">作废账单</el-button>
        <el-button v-if="canManagePayments && detail?.invoice?.status !== 'DRAFT' && detail?.invoice?.status !== 'VOID' && detail?.summary?.outstandingAmount > 0" :icon="Money" @click="openPaymentDialog">登记收款</el-button>
      </div>
    </div>

    <div v-if="detail" class="page-body">
      <div class="detail-grid">
        <div class="detail-hero">
          <div class="toolbar-row" style="margin-bottom: 12px">
            <div>
              <div class="detail-hero-title">{{ detail.invoice.invoiceNo }}</div>
              <div class="detail-hero-copy">
                {{ detail.invoice.customer.name }} / {{ getLabel(invoiceTypeMap, detail.invoice.type) }}
              </div>
            </div>
            <el-tag :type="getStatusTagType(detail.invoice.status)" size="large">
              {{ getLabel(invoiceStatusMap, detail.invoice.status) }}
            </el-tag>
          </div>

          <div class="summary-strip">
            <div class="summary-card"><div class="metric-label">账单总额</div><h3>{{ formatCurrency(detail.invoice.totalAmount) }}</h3><p>含税后的账单总金额</p></div>
            <div class="summary-card"><div class="metric-label">已收金额</div><h3>{{ formatCurrency(detail.invoice.paidAmount) }}</h3><p>已入账的支付金额</p></div>
            <div class="summary-card"><div class="metric-label">待收金额</div><h3>{{ formatCurrency(detail.summary.outstandingAmount) }}</h3><p>当前账单未结清金额</p></div>
            <div class="summary-card"><div class="metric-label">退款金额</div><h3>{{ formatCurrency(detail.summary.refundedAmount) }}</h3><p>关联退款累计金额</p></div>
          </div>
        </div>

        <div class="detail-side">
          <el-card class="page-card">
            <div class="section-heading"><h2>账单信息</h2><el-tag effect="plain">{{ detail.summary.paymentCount }} 笔支付</el-tag></div>
            <div class="detail-meta-list">
              <div class="detail-meta-item"><label>客户</label><strong>{{ detail.invoice.customer.name }}</strong></div>
              <div class="detail-meta-item"><label>关联订单</label><strong><el-button v-if="detail.invoice.order" text type="primary" @click="router.push(`/orders/${detail.invoice.order.id}`)">{{ detail.invoice.order.orderNo }}</el-button><span v-else>-</span></strong></div>
              <div class="detail-meta-item"><label>关联服务</label><strong><el-button v-if="detail.invoice.service" text type="primary" @click="router.push(`/services/${detail.invoice.service.id}`)">{{ detail.invoice.service.serviceNo }}</el-button><span v-else>-</span></strong></div>
              <div class="detail-meta-item"><label>税率</label><strong>{{ detail.invoice.taxProfileName || '自定义税率' }} / {{ detail.invoice.taxRate }}%</strong></div>
              <div class="detail-meta-item"><label>到期日</label><strong>{{ formatDate(detail.invoice.dueDate) }}</strong></div>
              <div class="detail-meta-item"><label>备注</label><strong>{{ detail.invoice.remark || '-' }}</strong></div>
            </div>
          </el-card>

          <el-card class="page-card">
            <div class="detail-meta-list">
              <div class="detail-meta-item"><label>处理建议</label><strong>{{ financeHint }}</strong></div>
              <div class="detail-meta-item"><label>支付记录</label><strong>{{ detail.summary.paymentCount }}</strong></div>
              <div class="detail-meta-item"><label>退款记录</label><strong>{{ detail.summary.refundCount }}</strong></div>
              <div class="detail-meta-item"><label>回调日志</label><strong>{{ detail.summary.callbackCount }}</strong></div>
              <div class="detail-meta-item"><label>财务任务</label><strong>{{ detail.invoice.billingJobs.length }}</strong></div>
            </div>
          </el-card>
        </div>
      </div>

      <el-card class="page-card">
        <el-tabs class="content-tabs">
          <el-tab-pane label="支付记录">
            <el-table :data="paymentRows" stripe>
              <el-table-column label="支付单号" min-width="180"><template #default="{ row }"><div class="primary-line">{{ row.paymentNo }}</div><div class="muted-line">{{ formatDateTime(row.createdAt) }}</div></template></el-table-column>
              <el-table-column label="渠道" width="140"><template #default="{ row }">{{ getLabel(paymentMethodMap, row.method) }}</template></el-table-column>
              <el-table-column label="金额" width="120"><template #default="{ row }">{{ formatCurrency(row.amount) }}</template></el-table-column>
              <el-table-column label="可退余额" width="120"><template #default="{ row }">{{ formatCurrency(row.refundableAmount) }}</template></el-table-column>
              <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(paymentStatusMap, row.status) }}</el-tag></template></el-table-column>
              <el-table-column label="交易号" min-width="180" prop="transactionNo" />
              <el-table-column v-if="canManagePayments" label="操作" width="120"><template #default="{ row }"><el-button v-if="row.refundableAmount > 0" text type="danger" @click="openRefundDialog(row)">发起退款</el-button><span v-else>-</span></template></el-table-column>
            </el-table>
          </el-tab-pane>
          <el-tab-pane label="退款与回调">
            <el-row :gutter="16">
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>退款记录</h2></div>
                <el-table :data="detail.invoice.refunds" stripe>
                  <el-table-column label="退款单号" min-width="180" prop="refundNo" />
                  <el-table-column label="金额" width="120"><template #default="{ row }">{{ formatCurrency(row.amount) }}</template></el-table-column>
                  <el-table-column label="渠道" width="140"><template #default="{ row }">{{ getLabel(paymentMethodMap, row.method) }}</template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(refundStatusMap, row.status) }}</el-tag></template></el-table-column>
                  <el-table-column label="处理人" width="160"><template #default="{ row }">{{ row.processedBy?.name || '-' }}<div class="muted-line">{{ getLabel(roleMap, row.processedBy?.role) }}</div></template></el-table-column>
                  <el-table-column prop="reason" label="原因" min-width="200" />
                </el-table>
              </el-col>
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>回调日志</h2></div>
                <el-table :data="detail.callbackLogs" stripe>
                  <el-table-column label="回调时间" width="180"><template #default="{ row }">{{ formatDateTime(row.createdAt) }}</template></el-table-column>
                  <el-table-column label="支付单号" min-width="160" prop="paymentNo" />
                  <el-table-column label="交易号" min-width="180" prop="transactionNo" />
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.callbackStatus)">{{ getLabel(callbackStatusMap, row.callbackStatus) }}</el-tag></template></el-table-column>
                  <el-table-column prop="message" label="说明" min-width="220" />
                </el-table>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane label="财务任务">
            <el-table :data="detail.invoice.billingJobs" stripe>
              <el-table-column label="任务类型" min-width="180"><template #default="{ row }">{{ getLabel(billingJobTypeMap, row.jobType) }}</template></el-table-column>
              <el-table-column prop="message" label="说明" min-width="220" />
              <el-table-column label="执行时间" width="180"><template #default="{ row }">{{ formatDateTime(row.executedAt) }}</template></el-table-column>
              <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(billingJobStatusMap, row.status) }}</el-tag></template></el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="审计日志">
            <el-table :data="detail.auditLogs" stripe>
              <el-table-column label="时间" width="180"><template #default="{ row }">{{ formatDateTime(row.createdAt) }}</template></el-table-column>
              <el-table-column label="动作" width="120" prop="action" />
              <el-table-column label="摘要" min-width="220" prop="summary" />
              <el-table-column label="详情" min-width="260" prop="detail" show-overflow-tooltip />
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>

    <el-dialog v-model="actionDialogVisible" :title="actionForm.action === 'issue' ? '签发账单' : '作废账单'" width="520px">
      <el-form label-position="top">
        <el-form-item v-if="actionForm.action === 'issue'" label="到期日期"><el-date-picker v-model="actionForm.dueDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" style="width: 100%" /></el-form-item>
        <el-form-item v-else label="作废原因"><el-input v-model="actionForm.reason" type="textarea" :rows="4" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="actionDialogVisible = false">取消</el-button><el-button type="primary" :loading="actionSubmitting" @click="submitAction">确认提交</el-button></template>
    </el-dialog>

    <el-dialog v-model="paymentDialogVisible" title="登记收款" width="520px">
      <el-form label-position="top">
        <el-form-item label="支付渠道"><el-select v-model="paymentForm.method" style="width: 100%"><el-option label="支付宝" value="ALIPAY" /><el-option label="微信支付" value="WECHAT" /><el-option label="余额" value="BALANCE" /><el-option label="银行转账" value="BANK_TRANSFER" /><el-option label="线下收款" value="OFFLINE" /><el-option label="PayPal" value="PAYPAL" /><el-option label="Stripe" value="STRIPE" /><el-option label="其他" value="OTHER" /></el-select></el-form-item>
        <el-form-item label="收款金额"><el-input-number v-model="paymentForm.amount" :min="0.01" :precision="2" style="width: 100%" /></el-form-item>
        <el-form-item label="交易流水号"><el-input v-model="paymentForm.transactionNo" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="paymentDialogVisible = false">取消</el-button><el-button type="primary" :loading="paymentSubmitting" @click="submitPayment">确认入账</el-button></template>
    </el-dialog>
    <el-dialog v-model="refundDialogVisible" title="发起退款" width="520px">
      <el-form label-position="top">
        <el-form-item label="原支付单号"><el-input :model-value="refundForm.paymentNo" disabled /></el-form-item>
        <el-form-item label="退款渠道"><el-select v-model="refundForm.method" style="width: 100%"><el-option label="支付宝" value="ALIPAY" /><el-option label="微信支付" value="WECHAT" /><el-option label="余额" value="BALANCE" /><el-option label="银行转账" value="BANK_TRANSFER" /><el-option label="线下收款" value="OFFLINE" /><el-option label="PayPal" value="PAYPAL" /><el-option label="Stripe" value="STRIPE" /><el-option label="其他" value="OTHER" /></el-select></el-form-item>
        <el-form-item label="可退余额"><el-input :model-value="refundForm.refundableAmount.toFixed(2)" disabled /></el-form-item>
        <el-form-item label="退款金额"><el-input-number v-model="refundForm.amount" :min="0.01" :max="refundForm.refundableAmount" :precision="2" style="width: 100%" /></el-form-item>
        <el-form-item label="退款原因"><el-input v-model="refundForm.reason" type="textarea" :rows="4" /></el-form-item>
        <el-form-item label="退款流水号"><el-input v-model="refundForm.transactionNo" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="refundDialogVisible = false">取消</el-button><el-button type="primary" :loading="refundSubmitting" @click="submitRefund">确认退款</el-button></template>
    </el-dialog>
  </div>
</template>
