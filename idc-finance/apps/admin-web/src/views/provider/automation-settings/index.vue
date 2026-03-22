<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { fetchAutomationSettings, updateAutomationSettings, type AutomationSettings } from "@/api/admin";

const loading = ref(false);
const saving = ref(false);

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
  logRetentionDays: 0
});

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
    hint: "自动开通、逾期暂停、逾期删除、账单提醒、信用额提醒、工单自动关闭"
  },
  {
    title: "暂停窗口",
    value: policy.autoSuspendEnabled ? `${policy.suspendAfterDays} 天` : "未启用",
    hint: `暂停前提前 ${policy.suspendNoticeDays} 天通知`
  },
  {
    title: "删除窗口",
    value: policy.autoTerminateEnabled ? `${policy.terminateAfterDays} 天` : "未启用",
    hint: `删除前提前 ${policy.terminateNoticeDays} 天通知`
  },
  {
    title: "工单自动关闭",
    value: policy.ticketAutoCloseOn ? `${policy.ticketAutoCloseHours} 小时` : "未启用",
    hint: "按最后回复时间自动关闭等待中的工单"
  },
  {
    title: "日志保留",
    value: `${policy.logRetentionDays} 天`,
    hint: "自动化任务日志和执行记录保留周期"
  }
]);

async function loadSettings() {
  loading.value = true;
  try {
    const data = await fetchAutomationSettings();
    Object.assign(policy, data);
  } finally {
    loading.value = false;
  }
}

async function saveSettings() {
  saving.value = true;
  try {
    const data = await updateAutomationSettings({ ...policy });
    Object.assign(policy, data);
    ElMessage.success("自动化策略已保存");
  } catch (error: any) {
    ElMessage.error(error?.message ?? "自动化策略保存失败");
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  void loadSettings();
});
</script>

<template>
  <PageWorkbench
    eyebrow="资源与商店 / 自动化"
    title="自动化策略"
    subtitle="按魔方财务的 cron 配置思路，把自动开通、逾期处理、提醒和日志保留集中管理。"
  >
    <template #actions>
      <el-button :loading="loading" @click="loadSettings">重新加载</el-button>
      <el-button type="primary" :loading="saving" @click="saveSettings">保存策略</el-button>
    </template>

    <template #metrics>
      <div class="automation-settings-metrics">
        <div v-for="card in summaryCards" :key="card.title" class="automation-settings-card automation-settings-card--summary">
          <span>{{ card.title }}</span>
          <strong>{{ card.value }}</strong>
          <p>{{ card.hint }}</p>
        </div>
      </div>
    </template>

    <div class="page-card" v-loading="loading">
      <div class="page-header">
        <div>
          <h2 class="section-title">策略配置</h2>
          <p class="page-subtitle">
            策略页负责定义规则，具体执行结果请到“自动化任务中心”查看。这样结构和老魔方的自动化逻辑更接近。
          </p>
        </div>
      </div>

      <el-alert
        title="策略保存后会作为自动化任务调度与提醒规则的基础参数。"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 18px"
      />

      <div class="automation-settings-grid">
        <div class="panel-card">
          <div class="section-card__head">
            <strong>自动开通与逾期暂停</strong>
            <span class="section-card__meta">控制付款后自动开通、逾期暂停和暂停前提醒</span>
          </div>
          <el-form label-position="top" :model="policy" class="automation-settings-form">
            <el-form-item label="自动开通">
              <el-switch v-model="policy.autoProvisionEnabled" active-text="启用" inactive-text="关闭" />
            </el-form-item>
            <el-form-item label="自动暂停">
              <el-switch v-model="policy.autoSuspendEnabled" active-text="启用" inactive-text="关闭" />
            </el-form-item>
            <div class="automation-settings-columns">
              <el-form-item label="逾期暂停天数">
                <el-input-number v-model="policy.suspendAfterDays" :min="0" :precision="0" controls-position="right" />
              </el-form-item>
              <el-form-item label="暂停提醒天数">
                <el-input-number v-model="policy.suspendNoticeDays" :min="0" :precision="0" controls-position="right" />
              </el-form-item>
            </div>
          </el-form>
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>逾期删除与提醒</strong>
            <span class="section-card__meta">控制逾期删除、账单提醒和信用额提醒</span>
          </div>
          <el-form label-position="top" :model="policy" class="automation-settings-form">
            <el-form-item label="自动删除">
              <el-switch v-model="policy.autoTerminateEnabled" active-text="启用" inactive-text="关闭" />
            </el-form-item>
            <div class="automation-settings-columns">
              <el-form-item label="逾期删除天数">
                <el-input-number v-model="policy.terminateAfterDays" :min="0" :precision="0" controls-position="right" />
              </el-form-item>
              <el-form-item label="删除提醒天数">
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
            <el-form-item label="信用额提醒">
              <el-switch v-model="policy.creditReminderOn" active-text="启用" inactive-text="关闭" />
            </el-form-item>
            <el-form-item label="信用额提醒天数">
              <el-input v-model="policy.creditReminderDays" placeholder="例如 3,1" />
            </el-form-item>
          </el-form>
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>工单与日志</strong>
            <span class="section-card__meta">控制工单自动关闭和任务日志保留周期</span>
          </div>
          <el-form label-position="top" :model="policy" class="automation-settings-form">
            <el-form-item label="工单自动关闭">
              <el-switch v-model="policy.ticketAutoCloseOn" active-text="启用" inactive-text="关闭" />
            </el-form-item>
            <el-form-item label="自动关闭小时数">
              <el-input-number v-model="policy.ticketAutoCloseHours" :min="0" :precision="0" controls-position="right" />
            </el-form-item>
            <el-form-item label="日志保留天数">
              <el-input-number v-model="policy.logRetentionDays" :min="1" :precision="0" controls-position="right" />
            </el-form-item>
          </el-form>
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>执行说明</strong>
            <span class="section-card__meta">把策略定义和任务执行拆开，便于排障与重试</span>
          </div>
          <div class="automation-settings-notes">
            <p>1. 自动开通控制付款完成后的交付阶段。</p>
            <p>2. 逾期暂停和逾期删除对应不同的生命周期节点。</p>
            <p>3. 账单提醒和信用额提醒建议配置多个提前天数窗口。</p>
            <p>4. 工单自动关闭通常按最后回复时间计算。</p>
            <p>5. 日志保留周期建议和审计策略保持一致。</p>
          </div>
        </div>
      </div>
    </div>
  </PageWorkbench>
</template>
