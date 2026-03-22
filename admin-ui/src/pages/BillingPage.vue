<script setup lang="ts">
import { ElMessage } from "element-plus";
import { onMounted, reactive, ref } from "vue";

import { http } from "@/api/http";
import { formatCurrency, formatDate, formatDateTime } from "@/utils/format";
import {
  billingJobStatusMap,
  billingJobTypeMap,
  getLabel,
  getStatusTagType,
  invoiceStatusMap,
  paymentMethodMap,
  serviceStatusMap,
} from "@/utils/maps";

const paymentSignTypeMap: Record<string, string> = {
  HEADER_SECRET: "请求头密钥",
  HMAC_SHA256: "HMAC-SHA256",
};

const loading = ref(false);
const running = ref(false);
const saving = ref(false);
const taxSubmitting = ref(false);
const gatewaySubmitting = ref(false);
const lastRunSummary = ref<any | null>(null);
const taxDialogVisible = ref(false);
const gatewayDialogVisible = ref(false);

const dataSource = ref<any>({
  jobs: [],
  dueServices: [],
  openInvoices: [],
  taxProfiles: [],
  paymentGateways: [],
});

const settings = reactive({
  invoiceLeadDays: 7,
  graceDays: 3,
  autoSuspendDays: 7,
  autoTerminateDays: 30,
  autoRenewByBalance: true,
  allowNegativeBalance: false,
  invoicePrefix: "INV",
  invoiceIssuerName: "IDC云业务管理系统",
  invoiceTaxNo: "",
  financeEmail: "",
  defaultTaxRate: 13,
});

const taxForm = reactive({
  id: "",
  code: "",
  name: "",
  taxRate: 13,
  description: "",
  isDefault: false,
  isActive: true,
});

const gatewayForm = reactive({
  id: "",
  method: "ALIPAY",
  name: "",
  merchantId: "",
  appId: "",
  apiBaseUrl: "",
  signType: "HEADER_SECRET",
  callbackSecret: "",
  callbackHeader: "",
  notifyUrl: "",
  returnUrl: "",
  isEnabled: true,
  remark: "",
});

function assignSettings(value: any) {
  settings.invoiceLeadDays = value.invoiceLeadDays;
  settings.graceDays = value.graceDays;
  settings.autoSuspendDays = value.autoSuspendDays;
  settings.autoTerminateDays = value.autoTerminateDays;
  settings.autoRenewByBalance = value.autoRenewByBalance;
  settings.allowNegativeBalance = value.allowNegativeBalance;
  settings.invoicePrefix = value.invoicePrefix || "INV";
  settings.invoiceIssuerName = value.invoiceIssuerName || "IDC云业务管理系统";
  settings.invoiceTaxNo = value.invoiceTaxNo || "";
  settings.financeEmail = value.financeEmail || "";
  settings.defaultTaxRate = value.defaultTaxRate ?? 13;
}

function resetTaxForm() {
  Object.assign(taxForm, {
    id: "",
    code: "",
    name: "",
    taxRate: settings.defaultTaxRate,
    description: "",
    isDefault: false,
    isActive: true,
  });
}

function resetGatewayForm() {
  Object.assign(gatewayForm, {
    id: "",
    method: "ALIPAY",
    name: "",
    merchantId: "",
    appId: "",
    apiBaseUrl: "",
    signType: "HEADER_SECRET",
    callbackSecret: "",
    callbackHeader: "",
    notifyUrl: "",
    returnUrl: "",
    isEnabled: true,
    remark: "",
  });
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/billing");
    dataSource.value = data.data;
    assignSettings(data.data.settings);
  } finally {
    loading.value = false;
  }
}

async function saveSettings() {
  saving.value = true;
  try {
    await http.put("/billing/settings", settings);
    ElMessage.success("财务配置已保存");
    await loadData();
  } finally {
    saving.value = false;
  }
}

async function runBilling() {
  running.value = true;
  try {
    const { data } = await http.post("/billing/run");
    lastRunSummary.value = data.data.summary;
    ElMessage.success("计费引擎执行完成");
    await loadData();
  } finally {
    running.value = false;
  }
}

function openTaxDialog(row?: any) {
  if (!row) {
    resetTaxForm();
  } else {
    Object.assign(taxForm, {
      id: row.id,
      code: row.code,
      name: row.name,
      taxRate: row.taxRate,
      description: row.description || "",
      isDefault: row.isDefault,
      isActive: row.isActive,
    });
  }

  taxDialogVisible.value = true;
}

async function submitTaxForm() {
  taxSubmitting.value = true;
  try {
    if (taxForm.id) {
      await http.put(`/billing/tax-profiles/${taxForm.id}`, taxForm);
      ElMessage.success("税率档案已更新");
    } else {
      await http.post("/billing/tax-profiles", taxForm);
      ElMessage.success("税率档案已创建");
    }

    taxDialogVisible.value = false;
    resetTaxForm();
    await loadData();
  } finally {
    taxSubmitting.value = false;
  }
}

function openGatewayDialog(row?: any) {
  if (!row) {
    resetGatewayForm();
  } else {
    Object.assign(gatewayForm, {
      id: row.id,
      method: row.method,
      name: row.name,
      merchantId: row.merchantId || "",
      appId: row.appId || "",
      apiBaseUrl: row.apiBaseUrl || "",
      signType: row.signType,
      callbackSecret: row.callbackSecret,
      callbackHeader: row.callbackHeader || "",
      notifyUrl: row.notifyUrl || "",
      returnUrl: row.returnUrl || "",
      isEnabled: row.isEnabled,
      remark: row.remark || "",
    });
  }

  gatewayDialogVisible.value = true;
}

async function submitGatewayForm() {
  gatewaySubmitting.value = true;
  try {
    if (gatewayForm.id) {
      await http.put(`/billing/payment-gateways/${gatewayForm.id}`, gatewayForm);
      ElMessage.success("支付渠道已更新");
    } else {
      await http.post("/billing/payment-gateways", gatewayForm);
      ElMessage.success("支付渠道已创建");
    }

    gatewayDialogVisible.value = false;
    resetGatewayForm();
    await loadData();
  } finally {
    gatewaySubmitting.value = false;
  }
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">计费与财务中心</h1>
        <div class="page-subtitle">
          统一配置续费规则、税务档案、支付渠道和回调签名，并查看待续费服务、未结清账单与任务执行结果。
        </div>
      </div>
      <div class="toolbar-actions">
        <el-button :loading="running" type="primary" @click="runBilling">立即执行</el-button>
        <el-button :loading="saving" @click="saveSettings">保存配置</el-button>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xl="12" :span="24">
        <el-card class="page-card">
          <template #header>
            <span>计费与开票基础设置</span>
          </template>
          <el-form label-position="top">
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="提前生成续费账单天数">
                  <el-input-number v-model="settings.invoiceLeadDays" :min="0" :max="90" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="逾期宽限天数">
                  <el-input-number v-model="settings.graceDays" :min="0" :max="30" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="逾期后自动暂停天数">
                  <el-input-number v-model="settings.autoSuspendDays" :min="0" :max="90" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="逾期后自动终止天数">
                  <el-input-number v-model="settings.autoTerminateDays" :min="1" :max="180" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="发票前缀">
                  <el-input v-model="settings.invoicePrefix" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="默认税率 (%)">
                  <el-input-number v-model="settings.defaultTaxRate" :min="0" :max="100" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="开票主体">
              <el-input v-model="settings.invoiceIssuerName" />
            </el-form-item>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="纳税人识别号">
                  <el-input v-model="settings.invoiceTaxNo" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="财务邮箱">
                  <el-input v-model="settings.financeEmail" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="允许余额自动续费">
                  <el-switch v-model="settings.autoRenewByBalance" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="允许客户余额为负">
                  <el-switch v-model="settings.allowNegativeBalance" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xl="12" :span="24">
        <el-card class="page-card">
          <template #header>
            <span>最近一次执行摘要</span>
          </template>
          <div v-if="lastRunSummary" class="summary-grid">
            <div class="summary-tile">
              <strong>{{ lastRunSummary.createdInvoices }}</strong>
              <span>生成账单</span>
            </div>
            <div class="summary-tile">
              <strong>{{ lastRunSummary.autoRenewed }}</strong>
              <span>自动续费</span>
            </div>
            <div class="summary-tile">
              <strong>{{ lastRunSummary.markedOverdue }}</strong>
              <span>标记逾期</span>
            </div>
            <div class="summary-tile">
              <strong>{{ lastRunSummary.suspended }}</strong>
              <span>自动暂停</span>
            </div>
            <div class="summary-tile">
              <strong>{{ lastRunSummary.terminated }}</strong>
              <span>自动终止</span>
            </div>
            <div class="summary-tile">
              <strong>{{ lastRunSummary.jobs }}</strong>
              <span>写入任务</span>
            </div>
          </div>
          <el-empty v-else description="执行一次计费引擎后，这里会展示最新处理摘要" />
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      <el-col :xl="10" :span="24">
        <el-card class="page-card">
          <template #header>
            <div class="card-header-with-action">
              <span>税率档案</span>
              <el-button type="primary" text @click="openTaxDialog()">新增档案</el-button>
            </div>
          </template>
          <el-table v-loading="loading" :data="dataSource.taxProfiles" stripe>
            <el-table-column label="档案" min-width="180">
              <template #default="{ row }">
                <div style="font-weight: 700">{{ row.name }}</div>
                <div class="muted-line">{{ row.code }}</div>
              </template>
            </el-table-column>
            <el-table-column label="税率" width="120">
              <template #default="{ row }">
                {{ row.taxRate }}%
              </template>
            </el-table-column>
            <el-table-column prop="description" label="说明" min-width="220" show-overflow-tooltip />
            <el-table-column label="状态" width="140">
              <template #default="{ row }">
                <div class="tag-group">
                  <el-tag :type="row.isActive ? 'success' : 'info'">
                    {{ row.isActive ? "启用" : "停用" }}
                  </el-tag>
                  <el-tag v-if="row.isDefault" type="warning">默认</el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button size="small" text @click="openTaxDialog(row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :xl="14" :span="24">
        <el-card class="page-card">
          <template #header>
            <div class="card-header-with-action">
              <span>支付渠道与签名策略</span>
              <el-button type="primary" text @click="openGatewayDialog()">新增渠道</el-button>
            </div>
          </template>
          <el-table v-loading="loading" :data="dataSource.paymentGateways" stripe>
            <el-table-column label="渠道" min-width="190">
              <template #default="{ row }">
                <div style="font-weight: 700">{{ row.name }}</div>
                <div class="muted-line">{{ getLabel(paymentMethodMap, row.method) }}</div>
              </template>
            </el-table-column>
            <el-table-column label="商户 / 应用" min-width="180">
              <template #default="{ row }">
                <div>{{ row.merchantId || "-" }}</div>
                <div class="muted-line">{{ row.appId || "-" }}</div>
              </template>
            </el-table-column>
            <el-table-column label="签名策略" width="150">
              <template #default="{ row }">
                {{ paymentSignTypeMap[row.signType] || row.signType }}
              </template>
            </el-table-column>
            <el-table-column label="签名头" min-width="150">
              <template #default="{ row }">
                {{ row.callbackHeader || "-" }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.isEnabled ? 'success' : 'info'">
                  {{ row.isEnabled ? "启用" : "停用" }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button size="small" text @click="openGatewayDialog(row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      <el-col :xl="12" :span="24">
        <el-card class="page-card">
          <template #header>
            <span>近期到期服务</span>
          </template>
          <el-table v-loading="loading" :data="dataSource.dueServices" stripe>
            <el-table-column prop="name" label="服务" min-width="180" />
            <el-table-column label="客户" min-width="160">
              <template #default="{ row }">
                {{ row.customer.name }}
              </template>
            </el-table-column>
            <el-table-column label="月费" width="140">
              <template #default="{ row }">
                {{ formatCurrency(row.monthlyCost) }}
              </template>
            </el-table-column>
            <el-table-column label="到期日" width="140">
              <template #default="{ row }">
                {{ formatDate(row.nextDueDate) }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">
                  {{ getLabel(serviceStatusMap, row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :xl="12" :span="24">
        <el-card class="page-card">
          <template #header>
            <span>未结清账单</span>
          </template>
          <el-table v-loading="loading" :data="dataSource.openInvoices" stripe>
            <el-table-column prop="invoiceNo" label="账单号" min-width="180" />
            <el-table-column label="客户" min-width="160">
              <template #default="{ row }">
                {{ row.customer.name }}
              </template>
            </el-table-column>
            <el-table-column label="未收金额" width="140">
              <template #default="{ row }">
                {{ formatCurrency(row.totalAmount - row.paidAmount) }}
              </template>
            </el-table-column>
            <el-table-column label="到期日" width="140">
              <template #default="{ row }">
                {{ formatDate(row.dueDate) }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">
                  {{ getLabel(invoiceStatusMap, row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="page-card">
      <template #header>
        <span>最近任务记录</span>
      </template>
      <el-table v-loading="loading" :data="dataSource.jobs" stripe>
        <el-table-column label="任务类型" width="180">
          <template #default="{ row }">
            {{ getLabel(billingJobTypeMap, row.jobType) }}
          </template>
        </el-table-column>
        <el-table-column label="客户" min-width="150">
          <template #default="{ row }">
            {{ row.customer?.name || "-" }}
          </template>
        </el-table-column>
        <el-table-column label="服务" min-width="170">
          <template #default="{ row }">
            {{ row.service?.serviceNo || row.service?.name || "-" }}
          </template>
        </el-table-column>
        <el-table-column prop="message" label="执行说明" min-width="260" />
        <el-table-column label="执行时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.executedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getLabel(billingJobStatusMap, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="taxDialogVisible" :title="taxForm.id ? '编辑税率档案' : '新增税率档案'" width="560px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="税率编码">
              <el-input v-model="taxForm.code" placeholder="例如 VAT13" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="税率名称">
              <el-input v-model="taxForm.name" placeholder="例如 增值税 13%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="税率 (%)">
              <el-input-number v-model="taxForm.taxRate" :min="0" :max="100" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="默认档案">
              <el-switch v-model="taxForm.isDefault" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="说明">
          <el-input v-model="taxForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="taxForm.isActive" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="taxDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="taxSubmitting" @click="submitTaxForm">
          保存
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="gatewayDialogVisible" :title="gatewayForm.id ? '编辑支付渠道' : '新增支付渠道'" width="720px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="支付方式">
              <el-select v-model="gatewayForm.method" style="width: 100%">
                <el-option
                  v-for="(label, value) in paymentMethodMap"
                  :key="value"
                  :label="label"
                  :value="value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="渠道名称">
              <el-input v-model="gatewayForm.name" placeholder="例如 微信支付 V3" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="商户号">
              <el-input v-model="gatewayForm.merchantId" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="应用标识">
              <el-input v-model="gatewayForm.appId" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="接口地址">
              <el-input v-model="gatewayForm.apiBaseUrl" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="签名策略">
              <el-select v-model="gatewayForm.signType" style="width: 100%">
                <el-option
                  v-for="(label, value) in paymentSignTypeMap"
                  :key="value"
                  :label="label"
                  :value="value"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="回调密钥">
              <el-input v-model="gatewayForm.callbackSecret" show-password />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="签名头">
              <el-input v-model="gatewayForm.callbackHeader" placeholder="如 x-payment-signature" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通知地址">
              <el-input v-model="gatewayForm.notifyUrl" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="返回地址">
              <el-input v-model="gatewayForm.returnUrl" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注">
          <el-input v-model="gatewayForm.remark" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="gatewayForm.isEnabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="gatewayDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="gatewaySubmitting" @click="submitGatewayForm">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.toolbar-actions {
  display: flex;
  gap: 12px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 12px;
}

.summary-tile {
  display: grid;
  gap: 6px;
  padding: 18px;
  border-radius: 16px;
  background: linear-gradient(180deg, #ffffff 0%, #f7fbff 100%);
  border: 1px solid var(--border-color);
}

.summary-tile strong {
  font-size: 28px;
  letter-spacing: -0.04em;
}

.summary-tile span,
.muted-line {
  color: var(--text-secondary);
  font-size: 13px;
}

.card-header-with-action {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.tag-group {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
</style>
