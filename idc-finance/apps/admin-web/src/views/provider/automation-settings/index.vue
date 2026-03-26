<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchAutomationSettings,
  fetchTicketPresets,
  runTicketAutoCloseSweep,
  updateAutomationSettings,
  updateTicketPresets,
  type AutomationSettings,
  type TicketAutoCloseSweepResponse,
  type TicketPresetReply
} from "@/api/admin";

const loading = ref(false);
const savingPolicy = ref(false);
const savingPresets = ref(false);
const sweeping = ref(false);
const sweepResult = ref<TicketAutoCloseSweepResponse | null>(null);

const policy = reactive<AutomationSettings>({
  autoProvisionEnabled: false,
  autoSuspendEnabled: false,
  suspendAfterDays: 0,
  suspendNoticeDays: 0,
  autoTerminateEnabled: false,
  terminateAfterDays: 0,
  terminateNoticeDays: 0,
  invoiceReminderOn: false,
  invoiceReminderDays: "",
  creditReminderOn: false,
  creditReminderDays: "",
  ticketAutoCloseOn: false,
  ticketAutoCloseHours: 0,
  logRetentionDays: 30
});

const presetRows = ref<TicketPresetReply[]>([]);

const departmentOptions = [
  { label: "通用", value: "" },
  { label: "技术支持", value: "Technical Support" },
  { label: "财务", value: "Finance" }
];

const statusOptions = [
  { label: "处理中", value: "PROCESSING" },
  { label: "等待客户回复", value: "WAITING_CUSTOMER" },
  { label: "打开", value: "OPEN" },
  { label: "已关闭", value: "CLOSED" }
];

const enabledCount = computed(
  () =>
    [
      policy.autoProvisionEnabled,
      policy.autoSuspendEnabled,
      policy.autoTerminateEnabled,
      policy.invoiceReminderOn,
      policy.creditReminderOn,
      policy.ticketAutoCloseOn
    ].filter(Boolean).length
);

const summaryCards = computed(() => [
  {
    title: "已启用策略",
    value: `${enabledCount.value}/6`,
    hint: "自动开通、到期暂停、到期删除、账单提醒、余额提醒、工单自动关单"
  },
  {
    title: "工单自动关单",
    value: policy.ticketAutoCloseOn ? `${policy.ticketAutoCloseHours} 小时` : "未启用",
    hint: "按最后一次回复时间计算，超时后系统自动关闭工单"
  },
  {
    title: "快捷回复模板",
    value: `${presetRows.value.length} 条`,
    hint: "后台详情页可一键套用，减少重复录入"
  },
  {
    title: "日志保留",
    value: `${policy.logRetentionDays} 天`,
    hint: "自动化任务日志与执行记录的默认保留周期"
  }
]);

function departmentLabel(value: string) {
  return departmentOptions.find(item => item.value === value)?.label || (value ? value : "通用");
}

function statusLabel(value: string) {
  return statusOptions.find(item => item.value === value)?.label || value;
}

async function loadData() {
  loading.value = true;
  try {
    const [settings, presets] = await Promise.all([fetchAutomationSettings(), fetchTicketPresets()]);
    Object.assign(policy, settings);
    presetRows.value = presets.items.map(item => ({ ...item }));
  } finally {
    loading.value = false;
  }
}

async function savePolicy() {
  savingPolicy.value = true;
  try {
    const data = await updateAutomationSettings({ ...policy });
    Object.assign(policy, data);
    ElMessage.success("自动化策略已保存");
  } catch (error: any) {
    ElMessage.error(error?.message ?? "自动化策略保存失败");
  } finally {
    savingPolicy.value = false;
  }
}

async function savePresets() {
  savingPresets.value = true;
  try {
    const data = await updateTicketPresets(presetRows.value);
    presetRows.value = data.items.map(item => ({ ...item }));
    ElMessage.success("快捷回复模板已保存");
  } catch (error: any) {
    ElMessage.error(error?.message ?? "快捷回复模板保存失败");
  } finally {
    savingPresets.value = false;
  }
}

async function runSweep() {
  sweeping.value = true;
  try {
    sweepResult.value = await runTicketAutoCloseSweep();
    ElMessage.success(sweepResult.value.message || "巡检已执行");
  } catch (error: any) {
    ElMessage.error(error?.message ?? "自动关单巡检执行失败");
  } finally {
    sweeping.value = false;
  }
}

function addPreset() {
  presetRows.value.push({
    key: `preset-${Date.now()}`,
    title: "",
    content: "",
    departmentName: "",
    status: "WAITING_CUSTOMER"
  });
}

function removePreset(index: number) {
  presetRows.value.splice(index, 1);
}

onMounted(() => {
  void loadData();
});
</script>

<template>
  <PageWorkbench
    eyebrow="资源与渠道 / 自动化"
    title="自动化策略"
    subtitle="把工单自动关单、提醒策略和快捷回复模板放到一个页面统一维护，后台和客户侧使用同一套配置。"
  >
    <template #actions>
      <el-button :loading="loading" @click="loadData">重新加载</el-button>
      <el-button type="primary" plain :loading="sweeping" @click="runSweep">立即巡检</el-button>
      <el-button type="primary" :loading="savingPolicy" @click="savePolicy">保存策略</el-button>
    </template>

    <template #metrics>
      <div class="summary-grid">
        <div v-for="card in summaryCards" :key="card.title" class="summary-card">
          <span>{{ card.title }}</span>
          <strong>{{ card.value }}</strong>
          <p>{{ card.hint }}</p>
        </div>
      </div>
    </template>

    <div class="settings-layout" v-loading="loading">
      <div class="page-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">基础自动化策略</h2>
            <p class="page-subtitle">控制自动开通、到期处理、提醒窗口与日志保留周期。</p>
          </div>
        </div>

        <div class="policy-grid">
          <div class="panel-card">
            <div class="panel-head">
              <strong>交付与停机</strong>
              <span>围绕付款后开通与到期暂停配置规则</span>
            </div>
            <el-form label-position="top" class="policy-form">
              <el-form-item label="自动开通">
                <el-switch v-model="policy.autoProvisionEnabled" active-text="启用" inactive-text="关闭" />
              </el-form-item>
              <el-form-item label="自动暂停">
                <el-switch v-model="policy.autoSuspendEnabled" active-text="启用" inactive-text="关闭" />
              </el-form-item>
              <div class="policy-columns">
                <el-form-item label="到期后多少天暂停">
                  <el-input-number v-model="policy.suspendAfterDays" :min="0" :precision="0" controls-position="right" />
                </el-form-item>
                <el-form-item label="暂停前提前几天提醒">
                  <el-input-number v-model="policy.suspendNoticeDays" :min="0" :precision="0" controls-position="right" />
                </el-form-item>
              </div>
            </el-form>
          </div>

          <div class="panel-card">
            <div class="panel-head">
              <strong>续费与删除</strong>
              <span>账单提醒、余额提醒和删除窗口集中管理</span>
            </div>
            <el-form label-position="top" class="policy-form">
              <el-form-item label="自动删除">
                <el-switch v-model="policy.autoTerminateEnabled" active-text="启用" inactive-text="关闭" />
              </el-form-item>
              <div class="policy-columns">
                <el-form-item label="到期后多少天删除">
                  <el-input-number v-model="policy.terminateAfterDays" :min="0" :precision="0" controls-position="right" />
                </el-form-item>
                <el-form-item label="删除前提前几天提醒">
                  <el-input-number v-model="policy.terminateNoticeDays" :min="0" :precision="0" controls-position="right" />
                </el-form-item>
              </div>
              <el-divider />
              <el-form-item label="账单提醒">
                <el-switch v-model="policy.invoiceReminderOn" active-text="启用" inactive-text="关闭" />
              </el-form-item>
              <el-form-item label="账单提醒天数">
                <el-input v-model="policy.invoiceReminderDays" placeholder="例如 7,3,1" />
              </el-form-item>
              <el-form-item label="余额提醒">
                <el-switch v-model="policy.creditReminderOn" active-text="启用" inactive-text="关闭" />
              </el-form-item>
              <el-form-item label="余额提醒天数">
                <el-input v-model="policy.creditReminderDays" placeholder="例如 3,1" />
              </el-form-item>
            </el-form>
          </div>

          <div class="panel-card">
            <div class="panel-head">
              <strong>工单与日志</strong>
              <span>真正影响工单自动关单的执行参数</span>
            </div>
            <el-form label-position="top" class="policy-form">
              <el-form-item label="启用工单自动关单">
                <el-switch v-model="policy.ticketAutoCloseOn" active-text="启用" inactive-text="关闭" />
              </el-form-item>
              <el-form-item label="等待客户多少小时后自动关闭">
                <el-input-number
                  v-model="policy.ticketAutoCloseHours"
                  :min="0"
                  :precision="0"
                  controls-position="right"
                />
              </el-form-item>
              <el-form-item label="日志保留天数">
                <el-input-number v-model="policy.logRetentionDays" :min="1" :precision="0" controls-position="right" />
              </el-form-item>
              <el-alert
                title="手动点击“立即巡检”会立即执行一次自动关单，并写入自动化任务中心。"
                type="info"
                :closable="false"
                show-icon
              />
            </el-form>
          </div>
        </div>
      </div>

      <div class="page-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">快捷回复模板</h2>
            <p class="page-subtitle">后台工单详情页会按部门拉取模板，技术、财务和通用模板可分别维护。</p>
          </div>
          <div class="inline-actions">
            <el-button @click="addPreset">新增模板</el-button>
            <el-button type="primary" :loading="savingPresets" @click="savePresets">保存模板</el-button>
          </div>
        </div>

        <el-table :data="presetRows" border stripe empty-text="暂无快捷回复模板">
          <el-table-column label="模板标题" min-width="180">
            <template #default="{ row }">
              <el-input v-model="row.title" placeholder="例如：补充信息" />
            </template>
          </el-table-column>
          <el-table-column label="适用部门" width="160">
            <template #default="{ row }">
              <el-select v-model="row.departmentName" style="width: 100%">
                <el-option v-for="item in departmentOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="回复后状态" width="180">
            <template #default="{ row }">
              <el-select v-model="row.status" style="width: 100%">
                <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="回复内容" min-width="360">
            <template #default="{ row }">
              <el-input v-model="row.content" type="textarea" :rows="3" resize="vertical" placeholder="请输入模板内容" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ $index }">
              <el-button type="danger" link @click="removePreset($index)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="page-card" v-if="sweepResult">
        <div class="page-header">
          <div>
            <h2 class="section-title">最近一次自动关单巡检</h2>
            <p class="page-subtitle">{{ sweepResult.message }}</p>
          </div>
          <div class="inline-tags">
            <el-tag type="success">关闭 {{ sweepResult.closedCount }}</el-tag>
            <el-tag type="info">跳过 {{ sweepResult.skippedCount }}</el-tag>
            <el-tag :type="sweepResult.failedCount ? 'danger' : 'success'">失败 {{ sweepResult.failedCount }}</el-tag>
          </div>
        </div>

        <div class="result-summary">
          <div class="result-card">
            <span>巡检任务</span>
            <strong>{{ sweepResult.taskNo || "-" }}</strong>
            <p>自动化任务中心可查看完整执行记录</p>
          </div>
          <div class="result-card">
            <span>扫描工单数</span>
            <strong>{{ sweepResult.checkedCount }}</strong>
            <p>仅扫描等待客户回复状态的工单</p>
          </div>
          <div class="result-card">
            <span>自动关单窗口</span>
            <strong>{{ sweepResult.autoCloseHours }} 小时</strong>
            <p>{{ sweepResult.enabled ? "当前策略已启用" : "当前策略未启用" }}</p>
          </div>
        </div>

        <el-table :data="sweepResult.items" border stripe empty-text="本次巡检没有关闭任何工单">
          <el-table-column prop="ticketNo" label="工单号" min-width="160" />
          <el-table-column prop="customerName" label="客户" min-width="140" />
          <el-table-column label="部门" min-width="120">
            <template #default="{ row }">{{ departmentLabel(row.departmentName) }}</template>
          </el-table-column>
          <el-table-column prop="autoCloseAt" label="应关单时间" min-width="170" />
          <el-table-column prop="closedAt" label="实际关闭时间" min-width="170" />
          <el-table-column label="结果" width="120">
            <template #default="{ row }">
              <el-tag :type="row.result === 'closed' ? 'success' : 'danger'">
                {{ row.result === "closed" ? "已关闭" : "失败" }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="message" label="说明" min-width="220" show-overflow-tooltip />
        </el-table>
      </div>

      <div class="page-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">配置说明</h2>
            <p class="page-subtitle">这套规则会同时影响后台工单工作台、工单详情页和客户中心工单详情。</p>
          </div>
        </div>
        <div class="notes-grid">
          <div class="note-item">
            <strong>1. 自动关单触发点</strong>
            <p>只处理“等待客户回复”的工单，按最后一次回复时间加关单小时数计算。</p>
          </div>
          <div class="note-item">
            <strong>2. 快捷回复模板</strong>
            <p>通用模板对所有部门可见，技术和财务模板仅在对应部门的工单详情页显示。</p>
          </div>
          <div class="note-item">
            <strong>3. 可回溯</strong>
            <p>自动关单会写入自动化任务中心和审计日志，方便追踪“谁触发、关了哪些单”。</p>
          </div>
        </div>
      </div>
    </div>
  </PageWorkbench>
</template>

<style scoped>
.settings-layout {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.summary-card,
.result-card {
  border-radius: 18px;
  padding: 18px 20px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(246, 248, 252, 0.98));
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 14px 40px rgba(15, 23, 42, 0.06);
}

.summary-card span,
.result-card span {
  display: block;
  font-size: 13px;
  color: #64748b;
}

.summary-card strong,
.result-card strong {
  display: block;
  margin-top: 10px;
  font-size: 28px;
  color: #0f172a;
}

.summary-card p,
.result-card p {
  margin: 10px 0 0;
  font-size: 13px;
  line-height: 1.6;
  color: #475569;
}

.policy-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.panel-card {
  border-radius: 18px;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: #fff;
  padding: 18px;
}

.panel-head {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 14px;
}

.panel-head strong {
  font-size: 16px;
  color: #0f172a;
}

.panel-head span {
  font-size: 13px;
  color: #64748b;
}

.policy-form :deep(.el-form-item) {
  margin-bottom: 16px;
}

.policy-columns {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.result-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 16px;
}

.notes-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.note-item {
  padding: 18px;
  border-radius: 16px;
  background: #f8fafc;
  border: 1px solid rgba(148, 163, 184, 0.2);
}

.note-item strong {
  display: block;
  color: #0f172a;
  margin-bottom: 8px;
}

.note-item p {
  margin: 0;
  line-height: 1.7;
  color: #475569;
}

.inline-tags,
.inline-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

@media (max-width: 1280px) {
  .summary-grid,
  .policy-grid,
  .result-summary,
  .notes-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .summary-grid,
  .policy-grid,
  .result-summary,
  .notes-grid,
  .policy-columns {
    grid-template-columns: 1fr;
  }
}
</style>
