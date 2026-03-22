<script setup lang="ts">
import { ElMessage } from "element-plus";
import { onMounted, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";

import { http } from "@/api/http";
import { formatCurrency, formatDateTime } from "@/utils/format";
import {
  callbackStatusMap,
  getLabel,
  getStatusTagType,
  paymentMethodMap,
  refundStatusMap,
} from "@/utils/maps";

const router = useRouter();
const loading = ref(false);
const dialogVisible = ref(false);
const refundDialogVisible = ref(false);
const payments = ref<any[]>([]);
const invoices = ref<any[]>([]);
const callbackLogs = ref<any[]>([]);
const refunds = ref<any[]>([]);
const summary = ref<any>({
  paymentCount: 0,
  callbackCount: 0,
  pendingCallbacks: 0,
  totalSuccess: 0,
  refundCount: 0,
  totalRefunded: 0,
  pendingRefunds: 0,
});

const form = reactive({
  invoiceId: "",
  method: "ALIPAY",
  amount: 0,
  transactionNo: "",
});

const refundForm = reactive({
  paymentId: "",
  amount: 0,
  method: "BALANCE",
  reason: "",
  transactionNo: "",
});

function refundableAmount(row: any) {
  const refundedAmount = (row.refunds || [])
    .filter((item: any) => item.status === "SUCCESS")
    .reduce((total: number, item: any) => total + item.amount, 0);

  return Math.max(row.amount - refundedAmount, 0);
}

function syncAmountFromInvoice(invoiceId: string) {
  const invoice = invoices.value.find((item) => item.id === invoiceId);

  if (!invoice) {
    return;
  }

  form.amount = Math.max(invoice.totalAmount - invoice.paidAmount, 0) / 100;
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/payments");
    payments.value = data.data.payments;
    invoices.value = data.data.invoices;
    callbackLogs.value = data.data.callbackLogs;
    refunds.value = data.data.refunds;
    summary.value = data.data.summary;

    if (!form.invoiceId && invoices.value[0]) {
      form.invoiceId = invoices.value[0].id;
    }

    syncAmountFromInvoice(form.invoiceId);
  } finally {
    loading.value = false;
  }
}

watch(
  () => form.invoiceId,
  (value) => {
    syncAmountFromInvoice(value);
  },
);

async function submitForm() {
  await http.post("/payments", form);
  ElMessage.success("收款登记成功");
  dialogVisible.value = false;
  form.transactionNo = "";
  await loadData();
}

function openRefund(row: any) {
  refundForm.paymentId = row.id;
  refundForm.amount = refundableAmount(row) / 100;
  refundForm.method = "BALANCE";
  refundForm.reason = "";
  refundForm.transactionNo = "";
  refundDialogVisible.value = true;
}

async function submitRefund() {
  await http.post(`/payments/${refundForm.paymentId}/refund`, refundForm);
  ElMessage.success("退款已处理");
  refundDialogVisible.value = false;
  await loadData();
}

function exportPayments(kind: "payments" | "refunds" | "callbacks") {
  window.open(`/api/payments/export?kind=${kind}`, "_blank");
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">收款管理</h1>
        <div class="page-subtitle">
          统一处理手工收款、退款、支付回调和财务导出。
        </div>
      </div>
      <div class="toolbar-actions">
        <el-button @click="router.push('/reconciliation')">财务对账</el-button>
        <el-button @click="exportPayments('payments')">导出收款</el-button>
        <el-button @click="exportPayments('refunds')">导出退款</el-button>
        <el-button @click="exportPayments('callbacks')">导出回调</el-button>
        <el-button type="primary" @click="dialogVisible = true">登记收款</el-button>
      </div>
    </div>

    <div class="metric-grid">
      <div class="metric-tile">
        <div class="metric-label">累计实收</div>
        <div class="metric-value">{{ formatCurrency(summary.totalSuccess) }}</div>
        <div class="metric-hint">当前财务流水中已成功入账的收款总额。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">收款记录</div>
        <div class="metric-value">{{ summary.paymentCount }}</div>
        <div class="metric-hint">系统内已登记的收款单总数。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">回调日志</div>
        <div class="metric-value">{{ summary.callbackCount }}</div>
        <div class="metric-hint">近期渠道异步回调日志总数。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">待处理回调</div>
        <div class="metric-value">{{ summary.pendingCallbacks }}</div>
        <div class="metric-hint">仍未完成确认或入账的回调数量。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">累计退款</div>
        <div class="metric-value">{{ formatCurrency(summary.totalRefunded) }}</div>
        <div class="metric-hint">已成功完成回冲的退款总额。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">退款记录</div>
        <div class="metric-value">{{ summary.refundCount }}</div>
        <div class="metric-hint">系统内已创建的退款单数量。</div>
      </div>
    </div>

    <el-card class="page-card">
      <template #header>
        <span>收款流水</span>
      </template>
      <el-table v-loading="loading" :data="payments" stripe>
        <el-table-column prop="paymentNo" label="收款单号" min-width="180" />
        <el-table-column label="客户" min-width="160">
          <template #default="{ row }">
            {{ row.customer.name }}
          </template>
        </el-table-column>
        <el-table-column label="关联账单" min-width="160">
          <template #default="{ row }">
            {{ row.invoice?.invoiceNo || "-" }}
          </template>
        </el-table-column>
        <el-table-column label="支付方式" width="140">
          <template #default="{ row }">
            {{ getLabel(paymentMethodMap, row.method) }}
          </template>
        </el-table-column>
        <el-table-column label="金额" width="140">
          <template #default="{ row }">
            {{ formatCurrency(row.amount) }}
          </template>
        </el-table-column>
        <el-table-column label="渠道流水号" min-width="170">
          <template #default="{ row }">
            {{ row.transactionNo || "-" }}
          </template>
        </el-table-column>
        <el-table-column label="入账时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.paidAt || row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="退款" width="160" fixed="right">
          <template #default="{ row }">
            <el-button
              size="small"
              type="danger"
              plain
              :disabled="row.status === 'FAILED' || refundableAmount(row) <= 0"
              @click="openRefund(row)"
            >
              发起退款
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="page-card">
      <template #header>
        <span>退款记录</span>
      </template>
      <el-table v-loading="loading" :data="refunds" stripe>
        <el-table-column label="退款单号" min-width="180">
          <template #default="{ row }">
            <div style="font-weight: 700">{{ row.refundNo }}</div>
            <div class="muted-line">{{ row.payment?.paymentNo || "-" }}</div>
          </template>
        </el-table-column>
        <el-table-column label="客户" min-width="160">
          <template #default="{ row }">
            {{ row.customer.name }}
          </template>
        </el-table-column>
        <el-table-column label="退款方式" width="130">
          <template #default="{ row }">
            {{ getLabel(paymentMethodMap, row.method) }}
          </template>
        </el-table-column>
        <el-table-column label="金额" width="140">
          <template #default="{ row }">
            {{ formatCurrency(row.amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="退款原因" min-width="220" show-overflow-tooltip />
        <el-table-column label="处理时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.processedAt || row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getLabel(refundStatusMap, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="page-card">
      <template #header>
        <span>渠道回调日志</span>
      </template>
      <el-table v-loading="loading" :data="callbackLogs" stripe>
        <el-table-column label="接收时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="支付方式" width="130">
          <template #default="{ row }">
            {{ getLabel(paymentMethodMap, row.method) }}
          </template>
        </el-table-column>
        <el-table-column prop="invoiceNo" label="账单号" min-width="160" />
        <el-table-column prop="transactionNo" label="渠道流水号" min-width="180" />
        <el-table-column prop="message" label="处理说明" min-width="260" />
        <el-table-column label="状态" width="140">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.isHandled ? row.callbackStatus : 'PENDING')">
              {{ getLabel(callbackStatusMap, row.isHandled ? row.callbackStatus : "PENDING") }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="登记收款" width="540px">
      <el-form label-position="top">
        <el-form-item label="选择账单">
          <el-select v-model="form.invoiceId" style="width: 100%">
            <el-option
              v-for="item in invoices"
              :key="item.id"
              :label="`${item.invoiceNo} / ${item.customer.name} / ${formatCurrency(item.totalAmount - item.paidAmount)}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="支付方式">
          <el-select v-model="form.method" style="width: 100%">
            <el-option
              v-for="(label, value) in paymentMethodMap"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="收款金额">
          <el-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="渠道流水号">
          <el-input v-model="form.transactionNo" placeholder="可选，用于对账或回调关联" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确认登记</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="refundDialogVisible" title="发起退款" width="540px">
      <el-form label-position="top">
        <el-form-item label="退款金额">
          <el-input-number
            v-model="refundForm.amount"
            :min="0.01"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="退款方式">
          <el-select v-model="refundForm.method" style="width: 100%">
            <el-option
              v-for="(label, value) in paymentMethodMap"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="退款原因">
          <el-input v-model="refundForm.reason" placeholder="例如 多收回退 / 客户取消订单" />
        </el-form-item>
        <el-form-item label="退款流水号">
          <el-input v-model="refundForm.transactionNo" placeholder="可选，用于渠道对账" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="refundDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRefund">确认退款</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.toolbar-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.muted-line {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 12px;
}
</style>
