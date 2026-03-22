<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { loadServices, type PortalService } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalConfiguration,
  formatPortalServiceStatus,
  portalTagTypeByStatus
} from "@/utils/business";

const localeStore = useLocaleStore();
const loading = ref(true);
const error = ref("");
const services = ref<PortalService[]>([]);
const keyword = ref("");
const status = ref("");

const filteredServices = computed(() =>
  services.value.filter(item => {
    const matchedKeyword =
      !keyword.value ||
      item.serviceNo.includes(keyword.value) ||
      item.productName.includes(keyword.value);
    const matchedStatus = !status.value || item.status === status.value;
    return matchedKeyword && matchedStatus;
  })
);

const activeCount = computed(() => services.value.filter(item => item.status === "ACTIVE").length);
const suspendedCount = computed(() => services.value.filter(item => item.status === "SUSPENDED").length);
const expiringSoon = computed(
  () =>
    services.value.filter(item => {
      const nextDue = item.nextDueAt ? new Date(item.nextDueAt) : null;
      if (!nextDue || Number.isNaN(nextDue.getTime())) {
        return false;
      }

      const diff = nextDue.getTime() - Date.now();
      return diff <= 7 * 24 * 60 * 60 * 1000 && diff >= 0;
    }).length
);

const copy = computed(() => ({
  title: pickLabel(localeStore.locale, "服务中心", "Services"),
  subtitle: pickLabel(
    localeStore.locale,
    "汇总当前客户已开通的全部服务，用于查看到期时间、运行状态和配置快照。",
    "Review due dates, running status, and configuration snapshots for all services."
  ),
  total: pickLabel(localeStore.locale, "服务总数", "Total Services"),
  active: pickLabel(localeStore.locale, "运行中", "Active"),
  suspended: pickLabel(localeStore.locale, "已暂停", "Suspended"),
  expiringSoon: pickLabel(localeStore.locale, "7 天内到期", "Due in 7 days"),
  listTitle: pickLabel(localeStore.locale, "服务列表", "Service List"),
  listDesc: pickLabel(
    localeStore.locale,
    "按服务编号、产品名称和状态快速筛选当前服务。",
    "Filter services by service number, product, and status."
  ),
  keywordPlaceholder: pickLabel(localeStore.locale, "搜索服务编号或产品名称", "Search service no. or product"),
  statusPlaceholder: pickLabel(localeStore.locale, "服务状态", "Service Status"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  serviceNo: pickLabel(localeStore.locale, "服务编号", "Service No."),
  productName: pickLabel(localeStore.locale, "产品名称", "Product"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  createdAt: pickLabel(localeStore.locale, "开通时间", "Created At"),
  nextDueAt: pickLabel(localeStore.locale, "下次到期", "Next Due"),
  configuration: pickLabel(localeStore.locale, "配置摘要", "Configuration"),
  empty: pickLabel(localeStore.locale, "暂无匹配的服务记录。", "No matching services.")
}));

async function fetchServices() {
  loading.value = true;
  error.value = "";

  try {
    services.value = await loadServices();
  } catch (err) {
    error.value = err instanceof Error ? err.message : pickLabel(localeStore.locale, "服务列表加载失败", "Failed to load services");
  } finally {
    loading.value = false;
  }
}

onMounted(fetchServices);
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />

    <section class="portal-card">
      <div class="portal-card-head">
        <div>
          <h1 class="portal-title">{{ copy.title }}</h1>
          <p class="portal-subtitle">{{ copy.subtitle }}</p>
        </div>
      </div>

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <div class="portal-stat">
          <h3>{{ copy.total }}</h3>
          <strong>{{ services.length }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.active }}</h3>
          <strong>{{ activeCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.suspended }}</h3>
          <strong>{{ suspendedCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.expiringSoon }}</h3>
          <strong>{{ expiringSoon }}</strong>
        </div>
      </div>
    </section>

    <section class="portal-card">
      <div class="portal-card-head">
        <div>
          <h2 class="portal-panel__title">{{ copy.listTitle }}</h2>
          <div class="portal-panel__meta">{{ copy.listDesc }}</div>
        </div>
      </div>

      <div class="portal-toolbar" style="margin: 18px 20px 0">
        <el-input v-model="keyword" :placeholder="copy.keywordPlaceholder" clearable />
        <el-select v-model="status" :placeholder="copy.statusPlaceholder" clearable>
          <el-option :label="formatPortalServiceStatus(localeStore.locale, 'ACTIVE')" value="ACTIVE" />
          <el-option :label="formatPortalServiceStatus(localeStore.locale, 'PENDING')" value="PENDING" />
          <el-option :label="formatPortalServiceStatus(localeStore.locale, 'SUSPENDED')" value="SUSPENDED" />
          <el-option :label="formatPortalServiceStatus(localeStore.locale, 'TERMINATED')" value="TERMINATED" />
        </el-select>
        <div class="portal-summary" style="margin: 0; padding: 12px 16px">
          <div class="portal-summary-row" style="font-size: 12px">
            <span>{{ copy.visible }}</span>
            <strong>{{ filteredServices.length }}</strong>
          </div>
        </div>
      </div>

      <el-table :data="filteredServices" border>
        <el-table-column prop="serviceNo" :label="copy.serviceNo" min-width="170" />
        <el-table-column prop="productName" :label="copy.productName" min-width="220" />
        <el-table-column :label="copy.status" min-width="120">
          <template #default="{ row }">
            <el-tag :type="portalTagTypeByStatus(row.status)" effect="light">
              {{ formatPortalServiceStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="copy.createdAt" min-width="180" />
        <el-table-column prop="nextDueAt" :label="copy.nextDueAt" min-width="160" />
        <el-table-column :label="copy.configuration" min-width="240">
          <template #default="{ row }">
            {{ formatPortalConfiguration(row.configuration, localeStore.locale) }}
          </template>
        </el-table-column>
      </el-table>

      <div v-if="!filteredServices.length" class="portal-empty-state" style="margin: 20px">
        {{ copy.empty }}
      </div>
    </section>
  </div>
</template>
