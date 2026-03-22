<script setup lang="ts">
import { ArrowLeft, CircleCheck, CloseBold } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";

import { http } from "@/api/http";
import { useAuthStore } from "@/stores/auth";
import { formatCurrency, formatDate, formatDateTime } from "@/utils/format";

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const loading = ref(false);
const submitting = ref(false);
const detail = ref<any | null>(null);
const actionDialogVisible = ref(false);

const actionForm = reactive({
  action: "approve",
  reason: "",
});

const canManageOrders = computed(() =>
  Boolean(authStore.user?.permissions.includes("orders.manage")),
);

const textMap = {
  source: { portal: "前台门户", sales: "销售录入", admin: "后台录入", "billing-engine": "计费引擎", "mofang-sync": "魔方同步" },
  orderStatus: { PENDING: "待支付", PAID: "已付款", PROVISIONING: "开通中", ACTIVE: "已生效", CANCELLED: "已取消", REFUNDED: "已退款" },
  serviceStatus: { PENDING: "待开通", PROVISIONING: "开通中", ACTIVE: "运行中", SUSPENDED: "已暂停", OVERDUE: "已逾期", TERMINATED: "已终止", EXPIRED: "已到期", FAILED: "异常" },
  invoiceStatus: { DRAFT: "草稿", ISSUED: "待支付", PARTIAL: "部分支付", PAID: "已支付", OVERDUE: "已逾期", VOID: "已作废" },
  cycle: { MONTHLY: "月付", QUARTERLY: "季付", SEMI_ANNUALLY: "半年付", ANNUALLY: "年付", BIENNIALLY: "两年付", TRIENNIALLY: "三年付", ONETIME: "一次性" },
  paymentMethod: { BALANCE: "余额", ALIPAY: "支付宝", WECHAT: "微信支付", BANK_TRANSFER: "银行转账", OFFLINE: "线下收款", PAYPAL: "PayPal", STRIPE: "Stripe", OTHER: "其他" },
  paymentStatus: { PENDING: "待确认", SUCCESS: "成功", FAILED: "失败", REFUNDED: "已退款" },
  refundStatus: { PENDING: "待处理", SUCCESS: "已退款", FAILED: "失败", CANCELED: "已取消" },
} as const;

function labelOf(map: Record<string, string>, value?: string | null) {
  if (!value) return "-";
  return map[value] ?? value;
}

function tagTypeForStatus(status?: string | null) {
  switch (status) {
    case "ACTIVE":
    case "SUCCESS":
    case "PAID":
      return "success";
    case "PENDING":
    case "PROVISIONING":
    case "PARTIAL":
      return "warning";
    case "FAILED":
    case "VOID":
    case "CANCELLED":
    case "REFUNDED":
    case "TERMINATED":
      return "danger";
    default:
      return "info";
  }
}

const timeline = computed(() => {
  if (!detail.value) return [];

  const rows = [
    { time: detail.value.order.createdAt, title: "订单创建", description: `订单 ${detail.value.order.orderNo} 已建立` },
    detail.value.order.paidAt
      ? { time: detail.value.order.paidAt, title: "订单收款", description: "订单已完成付款确认" }
      : null,
    ...detail.value.order.invoices.map((invoice: any) => ({
      time: invoice.createdAt,
      title: "账单生成",
      description: `${invoice.invoiceNo} / ${labelOf(textMap.invoiceStatus, invoice.status)}`,
    })),
    ...detail.value.order.payments.map((payment: any) => ({
      time: payment.createdAt,
      title: "收款记录",
      description: `${payment.paymentNo} / ${labelOf(textMap.paymentMethod, payment.method)} / ${formatCurrency(payment.amount)}`,
    })),
    ...detail.value.order.services.map((service: any) => ({
      time: service.updatedAt,
      title: "服务状态",
      description: `${service.serviceNo} / ${labelOf(textMap.serviceStatus, service.status)}`,
    })),
  ]
    .filter(Boolean)
    .sort((a: any, b: any) => Number(new Date(b.time)) - Number(new Date(a.time)));

  return rows;
});

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get(`/orders/${route.params.id}`);
    detail.value = data.data;
  } finally {
    loading.value = false;
  }
}

function openActionDialog(action: "approve" | "cancel") {
  actionForm.action = action;
  actionForm.reason = "";
  actionDialogVisible.value = true;
}

async function submitAction() {
  submitting.value = true;
  try {
    await http.post(`/orders/${route.params.id}/action`, actionForm);
    ElMessage.success(actionForm.action === "approve" ? "订单已核验通过" : "订单已取消");
    actionDialogVisible.value = false;
    await loadData();
  } finally {
    submitting.value = false;
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
        <el-button :icon="ArrowLeft" @click="router.push('/orders')">返回订单中心</el-button>
        <el-button
          v-if="canManageOrders && detail?.order?.status !== 'CANCELLED' && detail?.order?.status !== 'ACTIVE'"
          type="primary"
          :icon="CircleCheck"
          @click="openActionDialog('approve')"
        >
          核验通过
        </el-button>
        <el-button
          v-if="canManageOrders && detail?.order?.status !== 'CANCELLED'"
          type="danger"
          plain
          :icon="CloseBold"
          @click="openActionDialog('cancel')"
        >
          取消订单
        </el-button>
      </div>
    </div>

    <div v-if="detail" class="page-body">
      <div class="detail-grid">
        <div class="detail-hero">
          <div class="toolbar-row" style="margin-bottom: 12px">
            <div>
              <div class="detail-hero-title">{{ detail.order.orderNo }}</div>
              <div class="detail-hero-copy">
                {{ detail.order.customer.name }} / {{ labelOf(textMap.source, detail.order.source) }} / {{ detail.order.orderType }}
              </div>
            </div>
            <el-tag :type="tagTypeForStatus(detail.order.status)" size="large">
              {{ labelOf(textMap.orderStatus, detail.order.status) }}
            </el-tag>
          </div>

          <div class="summary-strip">
            <div class="summary-card"><div class="metric-label">订单金额</div><h3>{{ formatCurrency(detail.order.totalAmount) }}</h3><p>订单应收总金额</p></div>
            <div class="summary-card"><div class="metric-label">已收金额</div><h3>{{ formatCurrency(detail.order.paidAmount) }}</h3><p>订单已入账金额</p></div>
            <div class="summary-card"><div class="metric-label">待收金额</div><h3>{{ formatCurrency(detail.summary.outstandingAmount) }}</h3><p>尚未结清的订单金额</p></div>
            <div class="summary-card"><div class="metric-label">交付服务</div><h3>{{ detail.summary.serviceCount }}</h3><p>订单创建出的服务实例数</p></div>
          </div>
        </div>

        <div class="detail-side">
          <el-card class="page-card">
            <div class="section-heading"><h2>订单信息</h2><el-tag effect="plain">{{ detail.summary.itemCount }} 个项目</el-tag></div>
            <div class="detail-meta-list">
              <div class="detail-meta-item"><label>客户</label><strong>{{ detail.order.customer.name }}</strong></div>
              <div class="detail-meta-item"><label>创建时间</label><strong>{{ formatDateTime(detail.order.createdAt) }}</strong></div>
              <div class="detail-meta-item"><label>到期时间</label><strong>{{ formatDate(detail.order.dueDate) }}</strong></div>
              <div class="detail-meta-item"><label>付款时间</label><strong>{{ formatDateTime(detail.order.paidAt) }}</strong></div>
              <div class="detail-meta-item"><label>备注</label><strong>{{ detail.order.notes || "-" }}</strong></div>
            </div>
          </el-card>

          <el-card class="page-card">
            <div class="section-heading"><h2>审核工作区</h2></div>
            <div class="detail-meta-list">
              <div class="detail-meta-item"><label>账单数</label><strong>{{ detail.summary.invoiceCount }}</strong></div>
              <div class="detail-meta-item"><label>收款数</label><strong>{{ detail.summary.paymentCount }}</strong></div>
              <div class="detail-meta-item"><label>累计退款</label><strong>{{ formatCurrency(detail.summary.totalRefunded) }}</strong></div>
              <div class="detail-meta-item"><label>处理建议</label><strong>{{ detail.summary.outstandingAmount > 0 ? "待收款完成后再核验通过" : "可执行核验通过或进入开通流程" }}</strong></div>
            </div>
          </el-card>
        </div>
      </div>

      <el-card class="page-card">
        <el-tabs class="content-tabs">
          <el-tab-pane label="订单项目">
            <el-table :data="detail.order.items" stripe>
              <el-table-column label="项目" min-width="260">
                <template #default="{ row }">
                  <div class="primary-line">{{ row.title }}</div>
                  <div class="muted-line">{{ row.product.name }}</div>
                </template>
              </el-table-column>
              <el-table-column label="周期" width="120"><template #default="{ row }">{{ labelOf(textMap.cycle, row.cycle) }}</template></el-table-column>
              <el-table-column label="数量" width="100" prop="quantity" />
              <el-table-column label="单价" width="120"><template #default="{ row }">{{ formatCurrency(row.unitPrice) }}</template></el-table-column>
              <el-table-column label="金额" width="120"><template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template></el-table-column>
              <el-table-column label="服务" min-width="160">
                <template #default="{ row }">
                  <el-button v-if="row.service" text type="primary" @click="router.push(`/services/${row.service.id}`)">{{ row.service.serviceNo }}</el-button>
                  <span v-else>-</span>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="账单与支付">
            <el-row :gutter="16">
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>关联账单</h2></div>
                <el-table :data="detail.order.invoices" stripe>
                  <el-table-column label="账单号" min-width="170">
                    <template #default="{ row }">
                      <el-button text type="primary" @click="router.push(`/invoices/${row.id}`)">{{ row.invoiceNo }}</el-button>
                      <div class="muted-line">到期 {{ formatDate(row.dueDate) }}</div>
                    </template>
                  </el-table-column>
                  <el-table-column label="总额" width="120"><template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="tagTypeForStatus(row.status)">{{ labelOf(textMap.invoiceStatus, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>支付记录</h2></div>
                <el-table :data="detail.order.payments" stripe>
                  <el-table-column label="支付单号" min-width="170">
                    <template #default="{ row }">
                      <div class="primary-line">{{ row.paymentNo }}</div>
                      <div class="muted-line">{{ formatDateTime(row.createdAt) }}</div>
                    </template>
                  </el-table-column>
                  <el-table-column label="渠道" width="120"><template #default="{ row }">{{ labelOf(textMap.paymentMethod, row.method) }}</template></el-table-column>
                  <el-table-column label="金额" width="120"><template #default="{ row }">{{ formatCurrency(row.amount) }}</template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="tagTypeForStatus(row.status)">{{ labelOf(textMap.paymentStatus, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane label="交付服务">
            <el-table :data="detail.order.services" stripe>
              <el-table-column label="服务" min-width="220">
                <template #default="{ row }">
                  <el-button text type="primary" @click="router.push(`/services/${row.id}`)">{{ row.serviceNo }}</el-button>
                  <div class="muted-line">{{ row.name }} / {{ row.product.name }}</div>
                </template>
              </el-table-column>
              <el-table-column label="区域" min-width="180">
                <template #default="{ row }">{{ row.plan?.region?.name || row.region || "-" }} / {{ row.plan?.zone?.name || "-" }}</template>
              </el-table-column>
              <el-table-column label="规格" min-width="160"><template #default="{ row }">{{ row.plan?.flavor?.name || "-" }}</template></el-table-column>
              <el-table-column label="到期" width="120"><template #default="{ row }">{{ formatDate(row.nextDueDate) }}</template></el-table-column>
              <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="tagTypeForStatus(row.status)">{{ labelOf(textMap.serviceStatus, row.status) }}</el-tag></template></el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="订单时间线">
            <div class="list-stack">
              <div v-for="item in timeline" :key="`${item.title}-${item.time}`" class="mini-card">
                <div class="mini-title">{{ item.title }}</div>
                <div class="mini-copy">{{ item.description }}</div>
                <div class="mini-meta"><span>{{ formatDateTime(item.time) }}</span></div>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="审计日志">
            <el-table :data="detail.auditLogs" stripe>
              <el-table-column label="时间" width="180"><template #default="{ row }">{{ formatDateTime(row.createdAt) }}</template></el-table-column>
              <el-table-column label="动作" width="160" prop="action" />
              <el-table-column label="摘要" min-width="220" prop="summary" />
              <el-table-column label="详情" min-width="260" prop="detail" show-overflow-tooltip />
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>

    <el-dialog v-model="actionDialogVisible" :title="actionForm.action === 'approve' ? '核验通过订单' : '取消订单'" width="520px">
      <el-form label-position="top">
        <el-form-item :label="actionForm.action === 'approve' ? '审核说明' : '取消原因'">
          <el-input v-model="actionForm.reason" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="actionDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitAction">确认提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>
