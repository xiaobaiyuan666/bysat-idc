<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import { http } from "@/api/http";
import ArchitectureInfraSection from "@/components/architecture/ArchitectureInfraSection.vue";
import ArchitecturePlansSection from "@/components/architecture/ArchitecturePlansSection.vue";
import PortalUsersSection from "@/components/architecture/PortalUsersSection.vue";

const loading = ref(false);
const architecture = ref<any>({
  regions: [],
  zones: [],
  flavors: [],
  images: [],
  plans: [],
  products: [],
  customers: [],
  portalUsers: [],
});

const metrics = computed(() => [
  {
    label: "已发布地域",
    value: architecture.value.regions.length,
    hint: "用于销售地域、资源池和节点编排。",
  },
  {
    label: "可用区与规格",
    value: architecture.value.zones.length + architecture.value.flavors.length,
    hint: "可用区和规格共同决定上架能力边界。",
  },
  {
    label: "售卖方案",
    value: architecture.value.plans.length,
    hint: "前台可直接下单的云产品套餐。",
  },
  {
    label: "门户账号",
    value: architecture.value.portalUsers.length,
    hint: "客户联系人、主账号和协作成员体系。",
  },
]);

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/architecture");
    architecture.value = data.data;
  } finally {
    loading.value = false;
  }
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">云架构中心</h1>
        <div class="page-subtitle">
          统一维护地域、可用区、规格、镜像、售卖方案和门户联系人，作为公有云业务上架、交付和运营的底层配置中心。
        </div>
      </div>
      <el-button @click="loadData">刷新架构</el-button>
    </div>

    <div class="metric-grid">
      <div v-for="item in metrics" :key="item.label" class="metric-tile">
        <div class="metric-label">{{ item.label }}</div>
        <div class="metric-value">{{ item.value }}</div>
        <div class="metric-hint">{{ item.hint }}</div>
      </div>
    </div>

    <architecture-infra-section
      :loading="loading"
      :regions="architecture.regions"
      :zones="architecture.zones"
      :flavors="architecture.flavors"
      :images="architecture.images"
      @updated="loadData"
    />

    <architecture-plans-section
      :loading="loading"
      :plans="architecture.plans"
      :products="architecture.products"
      :regions="architecture.regions"
      :zones="architecture.zones"
      :flavors="architecture.flavors"
      :images="architecture.images"
      @updated="loadData"
    />

    <portal-users-section
      :loading="loading"
      :portal-users="architecture.portalUsers"
      :customers="architecture.customers"
      @updated="loadData"
    />
  </div>
</template>
