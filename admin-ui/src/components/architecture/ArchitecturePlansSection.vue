<script setup lang="ts">
import { Delete, Edit, Plus } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { computed, reactive, ref } from "vue";

import { http } from "@/api/http";
import { formatCurrency } from "@/utils/format";
import { cycleMap, getLabel } from "@/utils/maps";

type PlanConfig = {
  node?: string;
  os?: string;
  cpu?: number;
  memory?: number;
  system_disk_size?: number;
  data_disk_size?: number;
  network_type?: string;
  bw?: number;
  in_bw?: number;
  flow_limit?: number;
  flow_way?: string;
  ip_num?: number;
  peak_defence?: number;
};

const props = defineProps<{
  loading: boolean;
  plans: any[];
  products: any[];
  regions: any[];
  zones: any[];
  flavors: any[];
  images: any[];
}>();

const emit = defineEmits<{
  updated: [];
}>();

const dialogVisible = ref(false);
const editingId = ref("");
const form = reactive({
  productId: "",
  regionId: "",
  zoneId: "",
  flavorId: "",
  imageId: "",
  code: "",
  name: "",
  billingCycle: "MONTHLY",
  salePrice: 99,
  marketPrice: 129,
  setupFee: 0,
  stock: 100,
  isPublic: true,
  isActive: true,
  providerNodeId: "",
  providerOsCode: "",
  configCpu: 0,
  configMemoryGb: 0,
  systemDiskSize: 0,
  dataDiskSize: 0,
  networkType: "public",
  bandwidthMbps: 0,
  inboundBandwidthMbps: 0,
  flowLimitGb: 0,
  flowBillingMode: "month",
  ipCount: 1,
  peakDefenceGbps: 0,
  description: "",
});

const availableZones = computed(() =>
  props.zones.filter((item) => (!form.regionId ? true : item.regionId === form.regionId)),
);

function parseConfigOptions(raw?: string | null): PlanConfig {
  if (!raw) {
    return {};
  }

  try {
    const parsed = JSON.parse(raw) as PlanConfig;
    return parsed && typeof parsed === "object" ? parsed : {};
  } catch {
    return {};
  }
}

function numericValue(value: unknown) {
  const normalized = Number(value ?? 0);
  return Number.isFinite(normalized) ? normalized : 0;
}

function buildConfigSummary(row: any) {
  const config = parseConfigOptions(row.configOptions);
  const cpu = numericValue(config.cpu ?? row.flavor?.cpu);
  const memory = numericValue(config.memory ?? row.flavor?.memoryGb);
  const systemDisk = numericValue(config.system_disk_size ?? row.flavor?.storageGb);
  const bandwidth = numericValue(config.bw ?? row.flavor?.bandwidthMbps);
  const ipCount = numericValue(config.ip_num);
  const dataDisk = numericValue(config.data_disk_size);

  const parts = [
    cpu > 0 && memory > 0 ? `${cpu} 核 / ${memory}G` : "",
    systemDisk > 0 ? `系统盘 ${systemDisk}G` : "",
    dataDisk > 0 ? `数据盘 ${dataDisk}G` : "",
    bandwidth > 0 ? `${bandwidth} Mbps` : "",
    ipCount > 0 ? `${ipCount} 个 IP` : "",
    config.network_type ? String(config.network_type) : "",
  ].filter(Boolean);

  return parts.length > 0 ? parts.join(" / ") : "未配置云平台参数";
}

function resetForm() {
  Object.assign(form, {
    productId: props.products[0]?.id ?? "",
    regionId: props.regions[0]?.id ?? "",
    zoneId: "",
    flavorId: props.flavors[0]?.id ?? "",
    imageId: props.images[0]?.id ?? "",
    code: "",
    name: "",
    billingCycle: "MONTHLY",
    salePrice: 99,
    marketPrice: 129,
    setupFee: 0,
    stock: 100,
    isPublic: true,
    isActive: true,
    providerNodeId: "",
    providerOsCode: "",
    configCpu: 0,
    configMemoryGb: 0,
    systemDiskSize: 0,
    dataDiskSize: 0,
    networkType: "public",
    bandwidthMbps: 0,
    inboundBandwidthMbps: 0,
    flowLimitGb: 0,
    flowBillingMode: "month",
    ipCount: 1,
    peakDefenceGbps: 0,
    description: "",
  });
}

function openCreate() {
  editingId.value = "";
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: any) {
  const config = parseConfigOptions(row.configOptions);

  editingId.value = row.id;
  Object.assign(form, {
    productId: row.productId,
    regionId: row.regionId,
    zoneId: row.zoneId || "",
    flavorId: row.flavorId || "",
    imageId: row.imageId || "",
    code: row.code,
    name: row.name,
    billingCycle: row.billingCycle,
    salePrice: row.salePrice / 100,
    marketPrice: row.marketPrice ? row.marketPrice / 100 : 0,
    setupFee: row.setupFee / 100,
    stock: row.stock,
    isPublic: row.isPublic,
    isActive: row.isActive,
    providerNodeId: String(config.node ?? ""),
    providerOsCode: String(config.os ?? ""),
    configCpu: numericValue(config.cpu),
    configMemoryGb: numericValue(config.memory),
    systemDiskSize: numericValue(config.system_disk_size),
    dataDiskSize: numericValue(config.data_disk_size),
    networkType: String(config.network_type ?? "public"),
    bandwidthMbps: numericValue(config.bw),
    inboundBandwidthMbps: numericValue(config.in_bw),
    flowLimitGb: numericValue(config.flow_limit),
    flowBillingMode: String(config.flow_way ?? "month"),
    ipCount: numericValue(config.ip_num) || 1,
    peakDefenceGbps: numericValue(config.peak_defence),
    description: row.description || "",
  });
  dialogVisible.value = true;
}

async function submitForm() {
  const payload = {
    kind: "plan",
    ...form,
    ...(editingId.value ? { id: editingId.value } : {}),
  };

  if (editingId.value) {
    await http.put("/architecture", payload);
    ElMessage.success("售卖方案已更新");
  } else {
    await http.post("/architecture", payload);
    ElMessage.success("售卖方案已创建");
  }

  dialogVisible.value = false;
  emit("updated");
}

async function removeItem(row: any) {
  await ElMessageBox.confirm(`确认删除方案 ${row.name} 吗？`, "删除确认", {
    type: "warning",
  });

  await http.delete("/architecture", {
    data: {
      kind: "plan",
      id: row.id,
    },
  });

  ElMessage.success("售卖方案已删除");
  emit("updated");
}

resetForm();
</script>

<template>
  <el-card class="page-card">
    <template #header>
      <div class="card-header">
        <div>
          <div class="card-title">售卖方案</div>
          <div class="card-subtitle">
            公有云最终面向客户售卖的标准 SKU，统一沉淀价格、库存和云平台下发参数。
          </div>
        </div>
        <el-button type="primary" :icon="Plus" @click="openCreate">新增方案</el-button>
      </div>
    </template>

    <el-table v-loading="loading" :data="plans" stripe>
      <el-table-column label="方案信息" min-width="220">
        <template #default="{ row }">
          <div style="font-weight: 700">{{ row.name }}</div>
          <div class="muted-line">{{ row.code }} / {{ row.product?.name }}</div>
        </template>
      </el-table-column>
      <el-table-column label="地域 / 可用区" min-width="160">
        <template #default="{ row }">
          {{ row.region?.name }} {{ row.zone ? `/ ${row.zone.name}` : "" }}
        </template>
      </el-table-column>
      <el-table-column label="规格 / 镜像" min-width="180">
        <template #default="{ row }">
          <div>{{ row.flavor?.name || "自定义规格" }}</div>
          <div class="muted-line">{{ row.image?.name || "默认镜像" }}</div>
        </template>
      </el-table-column>
      <el-table-column label="云平台参数" min-width="240">
        <template #default="{ row }">
          <div>{{ buildConfigSummary(row) }}</div>
          <div class="muted-line">
            节点 {{ parseConfigOptions(row.configOptions).node || "-" }} / 系统
            {{ parseConfigOptions(row.configOptions).os || "-" }}
          </div>
        </template>
      </el-table-column>
      <el-table-column label="周期" width="110">
        <template #default="{ row }">
          {{ getLabel(cycleMap, row.billingCycle) }}
        </template>
      </el-table-column>
      <el-table-column label="售价" width="130">
        <template #default="{ row }">
          {{ formatCurrency(row.salePrice) }}
        </template>
      </el-table-column>
      <el-table-column prop="stock" label="库存" width="90" />
      <el-table-column label="公开" width="90">
        <template #default="{ row }">
          <el-tag :type="row.isPublic ? 'success' : 'info'">{{ row.isPublic ? "公开" : "内部" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.isActive ? 'success' : 'info'">{{ row.isActive ? "启用" : "停用" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button text type="primary" :icon="Edit" @click="openEdit(row)">编辑</el-button>
            <el-button text type="danger" :icon="Delete" @click="removeItem(row)">删除</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑售卖方案' : '新增售卖方案'" width="960px">
      <el-form label-position="top">
        <div class="section-title">商品与定价</div>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="关联产品">
              <el-select v-model="form.productId" style="width: 100%">
                <el-option
                  v-for="item in products"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="计费周期">
              <el-select v-model="form.billingCycle" style="width: 100%">
                <el-option
                  v-for="(label, value) in cycleMap"
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
            <el-form-item label="所属地域">
              <el-select v-model="form.regionId" style="width: 100%">
                <el-option
                  v-for="item in regions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="可用区">
              <el-select v-model="form.zoneId" clearable style="width: 100%">
                <el-option
                  v-for="item in availableZones"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="绑定规格">
              <el-select v-model="form.flavorId" clearable style="width: 100%">
                <el-option
                  v-for="item in flavors"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="默认镜像">
              <el-select v-model="form.imageId" clearable style="width: 100%">
                <el-option
                  v-for="item in images"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="方案编码">
              <el-input v-model="form.code" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="方案名称">
              <el-input v-model="form.name" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="销售价">
              <el-input-number v-model="form.salePrice" :precision="2" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="划线价">
              <el-input-number v-model="form.marketPrice" :precision="2" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="开通费">
              <el-input-number v-model="form.setupFee" :precision="2" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="库存">
              <el-input-number v-model="form.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="方案说明">
              <el-input v-model="form.description" />
            </el-form-item>
          </el-col>
        </el-row>

        <div class="section-title">云平台下发参数</div>
        <div class="section-copy">
          这些参数会写入售卖方案配置，并在后续服务开通、同步和前台展示中复用。
        </div>
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="平台节点标识">
              <el-input v-model="form.providerNodeId" placeholder="例如 hk-node-01" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="平台系统标识">
              <el-input v-model="form.providerOsCode" placeholder="例如 ubuntu-2204" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="网络类型">
              <el-input v-model="form.networkType" placeholder="例如 public / private" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="6">
            <el-form-item label="CPU(核)">
              <el-input-number v-model="form.configCpu" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="内存(GB)">
              <el-input-number v-model="form.configMemoryGb" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="系统盘(GB)">
              <el-input-number v-model="form.systemDiskSize" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="数据盘(GB)">
              <el-input-number v-model="form.dataDiskSize" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="6">
            <el-form-item label="出带宽(Mbps)">
              <el-input-number v-model="form.bandwidthMbps" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="入带宽(Mbps)">
              <el-input-number v-model="form.inboundBandwidthMbps" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="流量限制(GB)">
              <el-input-number v-model="form.flowLimitGb" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="IP 数量">
              <el-input-number v-model="form.ipCount" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="流量计费方式">
              <el-input v-model="form.flowBillingMode" placeholder="例如 month / peak95 / flow" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="防御峰值(Gbps)">
              <el-input-number v-model="form.peakDefenceGbps" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <div class="switch-row">
          <el-switch v-model="form.isPublic" active-text="门户公开" inactive-text="内部方案" />
          <el-switch v-model="form.isActive" active-text="启用" inactive-text="停用" />
        </div>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">{{ editingId ? "保存修改" : "确认创建" }}</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.card-title {
  font-size: 18px;
  font-weight: 700;
}

.card-subtitle,
.muted-line,
.section-copy {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 12px;
}

.row-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.switch-row {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.section-title {
  margin: 8px 0 4px;
  font-size: 14px;
  font-weight: 700;
  color: var(--text-primary);
}
</style>
