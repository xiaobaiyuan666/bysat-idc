<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchPluginConfig,
  fetchPluginConfigs,
  fetchPluginDefinitions,
  savePluginConfig,
  type PluginConfig,
  type PluginDefinition
} from "@/api/admin";

type PluginEditorState = {
  enabled: boolean;
  config: Record<string, unknown>;
};

const loading = ref(false);
const saving = ref(false);
const pluginDefs = ref<PluginDefinition[]>([]);
const pluginConfigs = ref<PluginConfig[]>([]);
const selectedCode = ref("");
const selectedConfig = ref<PluginConfig | null>(null);

const editor = reactive<PluginEditorState>({
  enabled: false,
  config: {}
});

const selectedDef = computed(() => pluginDefs.value.find(item => item.code === selectedCode.value) ?? null);
const selectedName = computed(() => selectedConfig.value?.name || selectedDef.value?.name || "-");
const selectedDesc = computed(() => selectedDef.value?.description || "");
const totalCount = computed(() => pluginDefs.value.length || pluginConfigs.value.length);
const enabledCount = computed(() => pluginConfigs.value.filter(item => item.enabled).length);
const requiredCount = computed(() => selectedDef.value?.fields.filter(field => field.required).length ?? 0);

function normalizeConfig(source: Record<string, unknown> | undefined) {
  if (!source || typeof source !== "object") return {};
  return { ...source };
}

function toInputValue(value: unknown) {
  if (value === null || value === undefined) return "";
  return String(value);
}

function onFieldInput(key: string, value: string | number) {
  editor.config[key] = String(value ?? "");
}

function fillMissingFieldDefaults(definition: PluginDefinition | null, source: Record<string, unknown>) {
  if (!definition) return source;
  const next = { ...source };
  for (const field of definition.fields) {
    if (next[field.key] === undefined || next[field.key] === null) {
      next[field.key] = "";
    }
  }
  return next;
}

function mergeConfigs(defs: PluginDefinition[], configs: PluginConfig[]) {
  return defs.map(definition => {
    const matched = configs.find(item => item.code === definition.code);
    return (
      matched ?? {
        code: definition.code,
        name: definition.name,
        enabled: false,
        config: {}
      }
    );
  });
}

function applySelectedConfig(config: PluginConfig | null) {
  selectedConfig.value = config;
  editor.enabled = config?.enabled ?? false;
  editor.config = fillMissingFieldDefaults(selectedDef.value, normalizeConfig(config?.config));
}

async function loadList() {
  loading.value = true;
  try {
    const [defsResp, cfgResp] = await Promise.all([fetchPluginDefinitions(), fetchPluginConfigs()]);
    pluginDefs.value = defsResp.items;
    pluginConfigs.value = mergeConfigs(defsResp.items, cfgResp.items);

    if (!selectedCode.value) {
      selectedCode.value = pluginConfigs.value[0]?.code || pluginDefs.value[0]?.code || "";
    }
    if (selectedCode.value) {
      const listItem = pluginConfigs.value.find(item => item.code === selectedCode.value) ?? null;
      applySelectedConfig(listItem);
    }
  } catch (error: any) {
    ElMessage.error(error?.message ?? "插件配置加载失败");
  } finally {
    loading.value = false;
  }
}

async function changePlugin(code: string | number) {
  if (!code) return;
  selectedCode.value = String(code);
  loading.value = true;
  try {
    const detail = await fetchPluginConfig(String(code));
    const index = pluginConfigs.value.findIndex(item => item.code === detail.code);
    if (index >= 0) {
      pluginConfigs.value[index] = detail;
    } else {
      pluginConfigs.value.push(detail);
    }
    applySelectedConfig(detail);
  } catch (error: any) {
    ElMessage.error(error?.message ?? "插件详情加载失败");
  } finally {
    loading.value = false;
  }
}

async function saveCurrentConfig() {
  if (!selectedCode.value) {
    ElMessage.warning("请先选择插件");
    return;
  }

  saving.value = true;
  try {
    const updated = await savePluginConfig(selectedCode.value, {
      enabled: editor.enabled,
      config: { ...editor.config }
    });
    applySelectedConfig(updated);

    const index = pluginConfigs.value.findIndex(item => item.code === updated.code);
    if (index >= 0) {
      pluginConfigs.value[index] = updated;
    } else {
      pluginConfigs.value.unshift(updated);
    }

    ElMessage.success("插件配置已保存");
  } catch (error: any) {
    ElMessage.error(error?.message ?? "插件配置保存失败");
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  void loadList();
});
</script>

<template>
  <PageWorkbench
    eyebrow="接口与上游"
    title="插件配置"
    subtitle="统一维护实名认证、短信、支付等插件开关与参数，让门户充值、短信验证码和插件实名都能在后台集中管理。"
  >
    <template #actions>
      <el-button :loading="loading" @click="loadList">重新加载</el-button>
      <el-button type="primary" :loading="saving" @click="saveCurrentConfig">保存配置</el-button>
    </template>

    <template #metrics>
      <div class="summary-grid">
        <div class="summary-card">
          <span>插件总数</span>
          <strong>{{ totalCount }}</strong>
          <p>当前可维护的插件定义数量</p>
        </div>
        <div class="summary-card">
          <span>已启用</span>
          <strong>{{ enabledCount }}</strong>
          <p>在业务链路中实际可调用的插件</p>
        </div>
        <div class="summary-card">
          <span>当前插件</span>
          <strong>{{ selectedName }}</strong>
          <p>{{ selectedDesc || "选择插件后查看配置说明" }}</p>
        </div>
        <div class="summary-card">
          <span>必填字段</span>
          <strong>{{ requiredCount }}</strong>
          <p>当前插件定义中的必填字段数量</p>
        </div>
      </div>
    </template>

    <div class="plugin-layout" v-loading="loading">
      <section class="page-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">插件选择</h2>
            <p class="page-subtitle">先选择要维护的插件，再在右侧保存启用状态和参数。</p>
          </div>
        </div>

        <el-radio-group
          v-if="pluginDefs.length"
          :model-value="selectedCode"
          class="plugin-radio-grid"
          @update:model-value="changePlugin"
        >
          <el-radio-button v-for="item in pluginDefs" :key="item.code" :label="item.code">
            {{ item.name }}
          </el-radio-button>
        </el-radio-group>
        <el-empty v-else description="暂无插件定义" />
      </section>

      <section class="page-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">参数配置</h2>
            <p class="page-subtitle">字段来自插件定义，保存后写入后台配置存储。</p>
          </div>
        </div>

        <template v-if="selectedCode && selectedDef">
          <el-form label-position="top" class="config-form">
            <el-form-item label="启用状态">
              <el-switch v-model="editor.enabled" active-text="启用" inactive-text="停用" />
            </el-form-item>

            <el-form-item
              v-for="field in selectedDef.fields"
              :key="field.key"
              :label="`${field.label}${field.required ? ' *' : ''}`"
            >
              <el-input
                :model-value="toInputValue(editor.config[field.key])"
                :placeholder="field.placeholder || `请输入${field.label}`"
                :type="field.secret ? 'password' : 'text'"
                show-password
                @update:model-value="onFieldInput(field.key, $event)"
              />
            </el-form-item>
          </el-form>
        </template>
        <el-empty v-else description="请选择插件后进行配置" />
      </section>
    </div>
  </PageWorkbench>
</template>

<style scoped>
.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.summary-card {
  border: 1px solid var(--el-border-color-light);
  border-radius: 12px;
  padding: 12px 14px;
  background: var(--el-bg-color-page);
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.summary-card span {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.summary-card strong {
  font-size: 22px;
  line-height: 1.2;
}

.summary-card p {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.plugin-layout {
  display: grid;
  grid-template-columns: minmax(300px, 360px) minmax(0, 1fr);
  gap: 16px;
}

.plugin-radio-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.config-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.config-form :deep(.el-form-item:first-child) {
  grid-column: 1 / -1;
}

@media (max-width: 1100px) {
  .plugin-layout {
    grid-template-columns: 1fr;
  }

  .config-form {
    grid-template-columns: 1fr;
  }
}
</style>
