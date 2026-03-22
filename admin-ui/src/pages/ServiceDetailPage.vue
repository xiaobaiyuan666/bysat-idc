<script setup lang="ts">
import {
  ArrowLeft,
  Connection,
  Edit,
  Key,
  Refresh,
  SwitchButton,
  VideoPause,
  VideoPlay,
} from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";

import { http } from "@/api/http";
import { useAuthStore } from "@/stores/auth";
import { formatCurrency, formatDate, formatDateTime } from "@/utils/format";
import {
  billingJobStatusMap,
  billingJobTypeMap,
  getLabel,
  getStatusTagType,
  invoiceStatusMap,
  providerActionMap,
  providerSyncStatusMap,
  providerTypeMap,
  resourceStatusMap,
  serviceStatusMap,
  ticketStatusMap,
} from "@/utils/maps";

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const loading = ref(false);
const saving = ref(false);
const actionSubmitting = ref(false);
const detail = ref<any | null>(null);
const editDialogVisible = ref(false);
const actionDialogVisible = ref(false);

const serviceForm = reactive({
  name: "",
  hostname: "",
  region: "",
  status: "ACTIVE",
  monthlyCost: 0,
  nextDueDate: "",
});

const actionForm = reactive({
  action: "resetPassword",
  password: "",
  imageId: "",
  rescueType: "linux",
  notes: "",
});

const canManageServices = computed(() =>
  Boolean(authStore.user?.permissions.includes("services.manage")),
);

const operationHint = computed(() => {
  const service = detail.value?.service;
  if (!service) return "-";
  if (service.status === "SUSPENDED") return "实例已暂停，建议先解除暂停再继续其他运维操作。";
  if (service.status === "FAILED") return "实例处于异常状态，优先同步状态或重装系统。";
  if ((detail.value?.summary?.outstandingAmount ?? 0) > 0) return "当前存在未结清账单，建议先处理续费与财务状态。";
  return "实例运行正常，可执行开关机、重启、改密、重装等运维动作。";
});

function fillServiceForm() {
  const service = detail.value?.service;
  if (!service) return;

  Object.assign(serviceForm, {
    name: service.name,
    hostname: service.hostname || "",
    region: service.region || "",
    status: service.status,
    monthlyCost: Number(((service.monthlyCost ?? 0) / 100).toFixed(2)),
    nextDueDate: service.nextDueDate ? String(service.nextDueDate).slice(0, 10) : "",
  });
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get(`/services/${route.params.id}`);
    detail.value = data.data;
    fillServiceForm();
  } finally {
    loading.value = false;
  }
}

async function showProviderResult(providerResult: any, action: string) {
  if (!providerResult) return;

  if (action === "getVnc" && providerResult.consoleUrl) {
    window.open(providerResult.consoleUrl, "_blank", "noopener,noreferrer");
  }

  const lines = [
    providerResult.message ? `平台消息：${providerResult.message}` : "",
    providerResult.taskId ? `任务 ID：${providerResult.taskId}` : "",
    providerResult.consoleUrl
      ? `控制台入口：<a href="${providerResult.consoleUrl}" target="_blank" rel="noreferrer">打开 VNC 控制台</a>`
      : "",
  ].filter(Boolean);

  if (!lines.length) return;

  await ElMessageBox.alert(lines.join("<br/>"), "云平台回执", {
    dangerouslyUseHTMLString: true,
    confirmButtonText: "知道了",
  });
}

async function runAction(action: string, payload?: Record<string, unknown>) {
  if (!detail.value) return;

  try {
    const { data } = await http.post(`/services/${route.params.id}/action`, { action, payload });
    ElMessage.success(`${getLabel(providerActionMap, action)}已提交`);
    await showProviderResult(data.providerResult, action);
    await loadData();
  } catch (error: any) {
    const providerResult = error?.response?.data?.providerResult;
    const message = error?.response?.data?.message || `${getLabel(providerActionMap, action)}执行失败`;
    ElMessage.error(message);
    if (providerResult) await showProviderResult(providerResult, action);
  }
}

async function confirmAction(action: string) {
  await ElMessageBox.confirm(`确认执行“${getLabel(providerActionMap, action)}”吗？`, "实例操作确认", { type: "warning" });
  await runAction(action);
}

function openEditDialog() {
  fillServiceForm();
  editDialogVisible.value = true;
}

async function submitService() {
  saving.value = true;
  try {
    await http.put(`/services/${route.params.id}`, serviceForm);
    ElMessage.success("服务信息已更新");
    editDialogVisible.value = false;
    await loadData();
  } finally {
    saving.value = false;
  }
}

function openActionDialog(action: "reinstall" | "resetPassword" | "rescueStart") {
  Object.assign(actionForm, {
    action,
    password: "",
    imageId: detail.value?.service?.plan?.image?.id || "",
    rescueType: "linux",
    notes: "",
  });
  actionDialogVisible.value = true;
}

async function submitPayloadAction() {
  if (actionForm.action !== "rescueStart" && actionForm.password && actionForm.password.length < 8) {
    ElMessage.warning("密码长度至少 8 位");
    return;
  }

  actionSubmitting.value = true;
  try {
    const payload: Record<string, unknown> = {};
    if (actionForm.password) payload.password = actionForm.password;
    if (actionForm.imageId) payload.imageId = actionForm.imageId;
    if (actionForm.rescueType) payload.rescueType = actionForm.rescueType;
    if (actionForm.notes) payload.notes = actionForm.notes;
    await runAction(actionForm.action, payload);
    actionDialogVisible.value = false;
  } finally {
    actionSubmitting.value = false;
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
        <el-button :icon="ArrowLeft" @click="router.push('/services')">返回业务中心</el-button>
        <el-button v-if="canManageServices" :icon="Edit" @click="openEditDialog">编辑服务</el-button>
        <el-button :icon="Refresh" @click="confirmAction('sync')">同步状态</el-button>
        <el-button type="primary" :icon="VideoPlay" @click="confirmAction('powerOn')">开机</el-button>
        <el-button :icon="VideoPause" @click="confirmAction('powerOff')">关机</el-button>
        <el-button :icon="Connection" @click="confirmAction('reboot')">重启</el-button>
        <el-button type="warning" :icon="SwitchButton" @click="confirmAction('suspend')">暂停</el-button>
        <el-button type="success" :icon="SwitchButton" @click="confirmAction('unsuspend')">解除暂停</el-button>
        <el-button type="info" :icon="Key" @click="confirmAction('getVnc')">打开 VNC</el-button>
      </div>
    </div>

    <div v-if="detail" class="page-body">
      <div class="detail-grid">
        <div class="detail-hero">
          <div class="toolbar-row" style="margin-bottom: 12px">
            <div>
              <div class="detail-hero-title">{{ detail.service.name }}</div>
              <div class="detail-hero-copy">
                {{ detail.service.serviceNo }} / {{ detail.service.customer.name }} / {{ detail.service.product.name }}
              </div>
            </div>
            <el-tag :type="getStatusTagType(detail.service.status)" size="large">
              {{ getLabel(serviceStatusMap, detail.service.status) }}
            </el-tag>
          </div>

          <div class="summary-strip">
            <div class="summary-card">
              <div class="metric-label">月成本</div>
              <h3>{{ formatCurrency(detail.service.monthlyCost) }}</h3>
              <p>当前服务计费成本</p>
            </div>
            <div class="summary-card">
              <div class="metric-label">待收金额</div>
              <h3>{{ formatCurrency(detail.summary.outstandingAmount) }}</h3>
              <p>关联账单剩余未结金额</p>
            </div>
            <div class="summary-card">
              <div class="metric-label">资源数量</div>
              <h3>{{ detail.summary.diskCount + detail.summary.ipCount }}</h3>
              <p>当前挂载磁盘与 IP 合计</p>
            </div>
            <div class="summary-card">
              <div class="metric-label">工单数量</div>
              <h3>{{ detail.summary.ticketCount }}</h3>
              <p>当前服务关联的工单总数</p>
            </div>
          </div>
        </div>

        <div class="detail-side">
          <el-card class="page-card">
            <div class="section-heading">
              <h2>服务信息</h2>
              <el-tag effect="plain">{{ getLabel(providerTypeMap, detail.service.providerType) }}</el-tag>
            </div>
            <div class="detail-meta-list">
              <div class="detail-meta-item"><label>主机名</label><strong>{{ detail.service.hostname || '-' }}</strong></div>
              <div class="detail-meta-item"><label>公网 IP</label><strong>{{ detail.service.ipAddress || '-' }}</strong></div>
              <div class="detail-meta-item"><label>远端资源 ID</label><strong>{{ detail.service.providerResourceId || '-' }}</strong></div>
              <div class="detail-meta-item"><label>下次到期</label><strong>{{ formatDate(detail.service.nextDueDate) }}</strong></div>
              <div class="detail-meta-item"><label>最近同步</label><strong>{{ formatDateTime(detail.service.lastSyncAt) }}</strong></div>
            </div>
          </el-card>

          <el-card class="page-card">
            <div class="detail-meta-list">
              <div class="detail-meta-item"><label>处理建议</label><strong>{{ operationHint }}</strong></div>
              <div class="detail-meta-item"><label>所属订单</label><strong><el-button v-if="detail.service.order" text type="primary" @click="router.push(`/orders/${detail.service.order.id}`)">{{ detail.service.order.orderNo }}</el-button><span v-else>-</span></strong></div>
              <div class="detail-meta-item"><label>VPC 网络</label><strong>{{ detail.service.vpcNetwork?.name || '-' }}</strong></div>
              <div class="detail-meta-item"><label>售卖方案</label><strong>{{ detail.service.plan?.name || '-' }}</strong></div>
              <div class="detail-meta-item"><label>地域 / 可用区</label><strong>{{ detail.service.plan?.region?.name || detail.service.region || '-' }} / {{ detail.service.plan?.zone?.name || '-' }}</strong></div>
              <div class="detail-meta-item"><label>规格 / 镜像</label><strong>{{ detail.service.plan?.flavor?.name || '-' }} / {{ detail.service.plan?.image?.name || '-' }}</strong></div>
              <div class="detail-meta-item"><label>更多动作</label><strong class="actions-wrap"><el-button text type="primary" @click="confirmAction('renew')">续费</el-button><el-button text type="primary" @click="openActionDialog('resetPassword')">改密</el-button><el-button text type="primary" @click="openActionDialog('reinstall')">重装</el-button><el-button text type="primary" @click="openActionDialog('rescueStart')">进入救援</el-button><el-button text type="primary" @click="confirmAction('rescueStop')">退出救援</el-button><el-button text type="danger" @click="confirmAction('terminate')">终止实例</el-button></strong></div>
            </div>
          </el-card>
        </div>
      </div>

      <el-card class="page-card">
        <el-tabs class="content-tabs">
          <el-tab-pane label="资源概览">
            <el-row :gutter="16">
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>IP 地址</h2><el-tag effect="plain">{{ detail.service.ipAddresses.length }} 条</el-tag></div>
                <el-table :data="detail.service.ipAddresses" stripe>
                  <el-table-column prop="address" label="地址" min-width="160" />
                  <el-table-column prop="version" label="版本" width="90" />
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(resourceStatusMap, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>磁盘</h2><el-tag effect="plain">{{ detail.service.disks.length }} 块</el-tag></div>
                <el-table :data="detail.service.disks" stripe>
                  <el-table-column prop="name" label="磁盘" min-width="160" />
                  <el-table-column label="容量" width="110"><template #default="{ row }">{{ row.sizeGb }} GB</template></el-table-column>
                  <el-table-column prop="mountPoint" label="挂载点" min-width="120" />
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(resourceStatusMap, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
            </el-row>
            <el-row :gutter="16" style="margin-top: 16px">
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>快照 / 备份</h2></div>
                <div class="list-stack">
                  <div class="mini-card"><div class="mini-title">快照数量</div><div class="mini-meta"><span>{{ detail.service.snapshots.length }} 个</span><span>最近更新 {{ formatDateTime(detail.service.snapshots[0]?.createdAt) }}</span></div></div>
                  <div class="mini-card"><div class="mini-title">备份数量</div><div class="mini-meta"><span>{{ detail.service.backups.length }} 个</span><span>最近更新 {{ formatDateTime(detail.service.backups[0]?.createdAt) }}</span></div></div>
                </div>
              </el-col>
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>安全组</h2></div>
                <el-table :data="detail.service.securityGroups" stripe>
                  <el-table-column prop="name" label="安全组" min-width="160" />
                  <el-table-column label="规则数" width="120"><template #default="{ row }">{{ row.rules?.length ?? 0 }}</template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(resourceStatusMap, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane label="账单与工单">
            <el-row :gutter="16">
              <el-col :xl="14" :span="24">
                <div class="section-heading"><h2>关联账单</h2></div>
                <el-table :data="detail.service.invoices" stripe>
                  <el-table-column label="账单号" min-width="180"><template #default="{ row }"><el-button text type="primary" @click="router.push(`/invoices/${row.id}`)">{{ row.invoiceNo }}</el-button></template></el-table-column>
                  <el-table-column label="总额" width="120"><template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template></el-table-column>
                  <el-table-column label="已收" width="120"><template #default="{ row }">{{ formatCurrency(row.paidAmount) }}</template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(invoiceStatusMap, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
              <el-col :xl="10" :span="24">
                <div class="section-heading"><h2>关联工单</h2></div>
                <el-table :data="detail.service.tickets" stripe>
                  <el-table-column label="工单" min-width="180"><template #default="{ row }"><div class="primary-line">{{ row.ticketNo }}</div><div class="muted-line">{{ row.subject }}</div></template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(ticketStatusMap, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane label="同步日志">
            <el-row :gutter="16">
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>云平台同步日志</h2></div>
                <el-table :data="detail.providerLogs" stripe>
                  <el-table-column label="动作" width="140"><template #default="{ row }">{{ getLabel(providerActionMap, row.action) }}</template></el-table-column>
                  <el-table-column prop="message" label="说明" min-width="220" />
                  <el-table-column label="时间" width="180"><template #default="{ row }">{{ formatDateTime(row.syncedAt) }}</template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(providerSyncStatusMap, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
              <el-col :xl="12" :span="24">
                <div class="section-heading"><h2>计费任务</h2></div>
                <el-table :data="detail.service.billingJobs" stripe>
                  <el-table-column label="任务类型" min-width="190"><template #default="{ row }">{{ getLabel(billingJobTypeMap, row.jobType) }}</template></el-table-column>
                  <el-table-column prop="message" label="说明" min-width="180" />
                  <el-table-column label="执行时间" width="180"><template #default="{ row }">{{ formatDateTime(row.executedAt) }}</template></el-table-column>
                  <el-table-column label="状态" width="120"><template #default="{ row }"><el-tag :type="getStatusTagType(row.status)">{{ getLabel(billingJobStatusMap, row.status) }}</el-tag></template></el-table-column>
                </el-table>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane label="审计日志">
            <el-table :data="detail.auditLogs" stripe>
              <el-table-column label="时间" width="180"><template #default="{ row }">{{ formatDateTime(row.createdAt) }}</template></el-table-column>
              <el-table-column label="动作" width="120" prop="action" />
              <el-table-column prop="summary" label="摘要" min-width="220" />
              <el-table-column prop="detail" label="详情" min-width="260" show-overflow-tooltip />
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>

    <el-dialog v-model="editDialogVisible" title="编辑服务" width="720px">
      <el-form label-position="top">
        <el-row :gutter="16">
          <el-col :md="12" :span="24"><el-form-item label="服务名称"><el-input v-model="serviceForm.name" /></el-form-item></el-col>
          <el-col :md="12" :span="24"><el-form-item label="主机名"><el-input v-model="serviceForm.hostname" /></el-form-item></el-col>
          <el-col :md="12" :span="24"><el-form-item label="地域"><el-input v-model="serviceForm.region" /></el-form-item></el-col>
          <el-col :md="12" :span="24"><el-form-item label="状态"><el-select v-model="serviceForm.status" style="width: 100%"><el-option label="待开通" value="PENDING" /><el-option label="开通中" value="PROVISIONING" /><el-option label="运行中" value="ACTIVE" /><el-option label="已暂停" value="SUSPENDED" /><el-option label="已逾期" value="OVERDUE" /><el-option label="已到期" value="EXPIRED" /><el-option label="已终止" value="TERMINATED" /><el-option label="异常" value="FAILED" /></el-select></el-form-item></el-col>
          <el-col :md="12" :span="24"><el-form-item label="月成本"><el-input-number v-model="serviceForm.monthlyCost" :min="0" :precision="2" style="width: 100%" /></el-form-item></el-col>
          <el-col :md="12" :span="24"><el-form-item label="下次到期"><el-date-picker v-model="serviceForm.nextDueDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" style="width: 100%" /></el-form-item></el-col>
        </el-row>
      </el-form>
      <template #footer><el-button @click="editDialogVisible = false">取消</el-button><el-button type="primary" :loading="saving" @click="submitService">保存服务</el-button></template>
    </el-dialog>
    <el-dialog v-model="actionDialogVisible" :title="getLabel(providerActionMap, actionForm.action)" width="560px">
      <el-form label-position="top">
        <el-form-item v-if="actionForm.action === 'reinstall'" label="镜像 ID">
          <el-input v-model="actionForm.imageId" placeholder="可选，填写目标镜像 ID" />
        </el-form-item>
        <el-form-item v-if="actionForm.action === 'reinstall' || actionForm.action === 'resetPassword'" label="新密码">
          <el-input v-model="actionForm.password" type="password" show-password placeholder="请输入至少 8 位密码" />
        </el-form-item>
        <el-form-item v-if="actionForm.action === 'rescueStart'" label="救援类型">
          <el-select v-model="actionForm.rescueType" style="width: 100%">
            <el-option label="Linux 救援模式" value="linux" />
            <el-option label="Windows 救援模式" value="windows" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作备注">
          <el-input v-model="actionForm.notes" type="textarea" :rows="4" placeholder="可选，记录本次运维操作意图" />
        </el-form-item>
      </el-form>
      <template #footer><el-button @click="actionDialogVisible = false">取消</el-button><el-button type="primary" :loading="actionSubmitting" @click="submitPayloadAction">确认提交</el-button></template>
    </el-dialog>
  </div>
</template>
