<script setup lang="ts">
import { Delete, Edit, Plus } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { reactive, ref } from "vue";

import { http } from "@/api/http";
import { getLabel, productCategoryMap, providerTypeMap } from "@/utils/maps";

type InfraKind = "region" | "zone" | "flavor" | "image";

const props = defineProps<{
  loading: boolean;
  regions: any[];
  zones: any[];
  flavors: any[];
  images: any[];
}>();

const emit = defineEmits<{
  updated: [];
}>();

const activeTab = ref<InfraKind>("region");
const dialogKind = ref<InfraKind>("region");
const dialogVisible = ref(false);
const editingId = ref("");

const regionForm = reactive({
  code: "",
  name: "",
  location: "",
  providerType: "MOFANG_CLOUD",
  providerRegionId: "",
  description: "",
  isActive: true,
  sortOrder: 0,
});

const zoneForm = reactive({
  regionId: "",
  code: "",
  name: "",
  providerZoneId: "",
  isActive: true,
  sortOrder: 0,
});

const flavorForm = reactive({
  code: "",
  name: "",
  category: "CLOUD_SERVER",
  cpu: 2,
  memoryGb: 4,
  storageGb: 40,
  bandwidthMbps: 10,
  gpuCount: 0,
  description: "",
  isActive: true,
  sortOrder: 0,
});

const imageForm = reactive({
  code: "",
  name: "",
  osType: "Linux",
  version: "Ubuntu 22.04",
  architecture: "x86_64",
  regionId: "",
  description: "",
  isPublic: true,
  isActive: true,
});

function resetForm(kind: InfraKind) {
  if (kind === "region") {
    Object.assign(regionForm, {
      code: "",
      name: "",
      location: "",
      providerType: "MOFANG_CLOUD",
      providerRegionId: "",
      description: "",
      isActive: true,
      sortOrder: 0,
    });
  }

  if (kind === "zone") {
    Object.assign(zoneForm, {
      regionId: props.regions[0]?.id ?? "",
      code: "",
      name: "",
      providerZoneId: "",
      isActive: true,
      sortOrder: 0,
    });
  }

  if (kind === "flavor") {
    Object.assign(flavorForm, {
      code: "",
      name: "",
      category: "CLOUD_SERVER",
      cpu: 2,
      memoryGb: 4,
      storageGb: 40,
      bandwidthMbps: 10,
      gpuCount: 0,
      description: "",
      isActive: true,
      sortOrder: 0,
    });
  }

  if (kind === "image") {
    Object.assign(imageForm, {
      code: "",
      name: "",
      osType: "Linux",
      version: "Ubuntu 22.04",
      architecture: "x86_64",
      regionId: "",
      description: "",
      isPublic: true,
      isActive: true,
    });
  }
}

function openCreate(kind: InfraKind) {
  dialogKind.value = kind;
  editingId.value = "";
  resetForm(kind);
  dialogVisible.value = true;
}

function openEdit(kind: InfraKind, row: any) {
  dialogKind.value = kind;
  editingId.value = row.id;
  resetForm(kind);

  if (kind === "region") {
    Object.assign(regionForm, {
      code: row.code,
      name: row.name,
      location: row.location || "",
      providerType: row.providerType,
      providerRegionId: row.providerRegionId || "",
      description: row.description || "",
      isActive: row.isActive,
      sortOrder: row.sortOrder,
    });
  }

  if (kind === "zone") {
    Object.assign(zoneForm, {
      regionId: row.regionId,
      code: row.code,
      name: row.name,
      providerZoneId: row.providerZoneId || "",
      isActive: row.isActive,
      sortOrder: row.sortOrder,
    });
  }

  if (kind === "flavor") {
    Object.assign(flavorForm, {
      code: row.code,
      name: row.name,
      category: row.category,
      cpu: row.cpu,
      memoryGb: row.memoryGb,
      storageGb: row.storageGb || 0,
      bandwidthMbps: row.bandwidthMbps || 0,
      gpuCount: row.gpuCount || 0,
      description: row.description || "",
      isActive: row.isActive,
      sortOrder: row.sortOrder,
    });
  }

  if (kind === "image") {
    Object.assign(imageForm, {
      code: row.code,
      name: row.name,
      osType: row.osType,
      version: row.version || "",
      architecture: row.architecture,
      regionId: row.regionId || "",
      description: row.description || "",
      isPublic: row.isPublic,
      isActive: row.isActive,
    });
  }

  dialogVisible.value = true;
}

function currentPayload(kind: InfraKind) {
  if (kind === "region") return { ...regionForm };
  if (kind === "zone") return { ...zoneForm };
  if (kind === "flavor") return { ...flavorForm };
  return { ...imageForm };
}

async function submitForm() {
  const payload = {
    kind: dialogKind.value,
    ...currentPayload(dialogKind.value),
    ...(editingId.value ? { id: editingId.value } : {}),
  };

  if (editingId.value) {
    await http.put("/architecture", payload);
    ElMessage.success("基础架构已更新");
  } else {
    await http.post("/architecture", payload);
    ElMessage.success("基础架构已创建");
  }

  dialogVisible.value = false;
  emit("updated");
}

async function removeItem(kind: InfraKind, row: any) {
  await ElMessageBox.confirm(`确认删除 ${row.name || row.code} 吗？`, "删除确认", {
    type: "warning",
  });

  await http.delete("/architecture", {
    data: {
      kind,
      id: row.id,
    },
  });

  ElMessage.success("基础架构已删除");
  emit("updated");
}
</script>

<template>
  <el-card class="page-card">
    <template #header>
      <div class="card-header">
        <div>
          <div class="card-title">基础架构</div>
          <div class="card-subtitle">地域、可用区、规格和镜像是公有云售卖与调度的底层模型。</div>
        </div>
      </div>
    </template>

    <el-tabs v-model="activeTab">
      <el-tab-pane label="地域" name="region">
        <div class="toolbar-inline">
          <span class="helper-copy">定义公有云销售地域与资源池入口。</span>
          <el-button type="primary" :icon="Plus" @click="openCreate('region')">新增地域</el-button>
        </div>
        <el-table v-loading="loading" :data="regions" stripe>
          <el-table-column prop="name" label="名称" min-width="140" />
          <el-table-column prop="code" label="编码" width="140" />
          <el-table-column prop="location" label="物理位置" min-width="150" />
          <el-table-column label="云平台" width="120">
            <template #default="{ row }">
              {{ getLabel(providerTypeMap, row.providerType) }}
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
                <el-button text type="primary" :icon="Edit" @click="openEdit('region', row)">编辑</el-button>
                <el-button text type="danger" :icon="Delete" @click="removeItem('region', row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="可用区" name="zone">
        <div class="toolbar-inline">
          <span class="helper-copy">用于容量隔离和可用性规划。</span>
          <el-button type="primary" :icon="Plus" @click="openCreate('zone')">新增可用区</el-button>
        </div>
        <el-table v-loading="loading" :data="zones" stripe>
          <el-table-column prop="name" label="名称" min-width="140" />
          <el-table-column prop="code" label="编码" width="140" />
          <el-table-column label="所属地域" min-width="150">
            <template #default="{ row }">
              {{ row.region?.name }}
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
                <el-button text type="primary" :icon="Edit" @click="openEdit('zone', row)">编辑</el-button>
                <el-button text type="danger" :icon="Delete" @click="removeItem('zone', row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="规格" name="flavor">
        <div class="toolbar-inline">
          <span class="helper-copy">沉淀标准规格，避免商品配置散落在多个页面。</span>
          <el-button type="primary" :icon="Plus" @click="openCreate('flavor')">新增规格</el-button>
        </div>
        <el-table v-loading="loading" :data="flavors" stripe>
          <el-table-column prop="name" label="规格名称" min-width="160" />
          <el-table-column prop="code" label="编码" width="140" />
          <el-table-column label="分类" width="120">
            <template #default="{ row }">
              {{ getLabel(productCategoryMap, row.category) }}
            </template>
          </el-table-column>
          <el-table-column label="核心参数" min-width="180">
            <template #default="{ row }">
              {{ row.cpu }}C / {{ row.memoryGb }}G / {{ row.storageGb || 0 }}G
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <div class="row-actions">
                <el-button text type="primary" :icon="Edit" @click="openEdit('flavor', row)">编辑</el-button>
                <el-button text type="danger" :icon="Delete" @click="removeItem('flavor', row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="镜像" name="image">
        <div class="toolbar-inline">
          <span class="helper-copy">统一管理公共镜像和地域镜像。</span>
          <el-button type="primary" :icon="Plus" @click="openCreate('image')">新增镜像</el-button>
        </div>
        <el-table v-loading="loading" :data="images" stripe>
          <el-table-column prop="name" label="镜像名称" min-width="180" />
          <el-table-column prop="code" label="编码" width="140" />
          <el-table-column label="系统" width="160">
            <template #default="{ row }">
              {{ row.osType }} {{ row.version || "" }}
            </template>
          </el-table-column>
          <el-table-column label="地域" min-width="140">
            <template #default="{ row }">
              {{ row.region?.name || "全局镜像" }}
            </template>
          </el-table-column>
          <el-table-column label="公开" width="90">
            <template #default="{ row }">
              <el-tag :type="row.isPublic ? 'success' : 'info'">{{ row.isPublic ? "公开" : "内部" }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <div class="row-actions">
                <el-button text type="primary" :icon="Edit" @click="openEdit('image', row)">编辑</el-button>
                <el-button text type="danger" :icon="Delete" @click="removeItem('image', row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑基础架构' : '新增基础架构'" width="680px">
      <el-form label-position="top">
        <template v-if="dialogKind === 'region'">
          <el-row :gutter="12">
            <el-col :span="12"><el-form-item label="地域编码"><el-input v-model="regionForm.code" /></el-form-item></el-col>
            <el-col :span="12"><el-form-item label="地域名称"><el-input v-model="regionForm.name" /></el-form-item></el-col>
          </el-row>
          <el-row :gutter="12">
            <el-col :span="12"><el-form-item label="物理位置"><el-input v-model="regionForm.location" /></el-form-item></el-col>
            <el-col :span="12">
              <el-form-item label="云平台类型">
                <el-select v-model="regionForm.providerType" style="width: 100%">
                  <el-option
                    v-for="(label, value) in providerTypeMap"
                    :key="value"
                    :label="label"
                    :value="value"
                  />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="平台地域 ID"><el-input v-model="regionForm.providerRegionId" /></el-form-item>
          <div class="switch-row"><el-switch v-model="regionForm.isActive" active-text="启用" inactive-text="停用" /></div>
        </template>

        <template v-else-if="dialogKind === 'zone'">
          <el-row :gutter="12">
            <el-col :span="12">
              <el-form-item label="所属地域">
                <el-select v-model="zoneForm.regionId" style="width: 100%">
                  <el-option v-for="item in regions" :key="item.id" :label="item.name" :value="item.id" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12"><el-form-item label="可用区编码"><el-input v-model="zoneForm.code" /></el-form-item></el-col>
          </el-row>
          <el-form-item label="可用区名称"><el-input v-model="zoneForm.name" /></el-form-item>
          <div class="switch-row"><el-switch v-model="zoneForm.isActive" active-text="启用" inactive-text="停用" /></div>
        </template>

        <template v-else-if="dialogKind === 'flavor'">
          <el-row :gutter="12">
            <el-col :span="12"><el-form-item label="规格编码"><el-input v-model="flavorForm.code" /></el-form-item></el-col>
            <el-col :span="12"><el-form-item label="规格名称"><el-input v-model="flavorForm.name" /></el-form-item></el-col>
          </el-row>
          <el-row :gutter="12">
            <el-col :span="12">
              <el-form-item label="产品分类">
                <el-select v-model="flavorForm.category" style="width: 100%">
                  <el-option
                    v-for="(label, value) in productCategoryMap"
                    :key="value"
                    :label="label"
                    :value="value"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="4"><el-form-item label="CPU"><el-input-number v-model="flavorForm.cpu" :min="1" style="width: 100%" /></el-form-item></el-col>
            <el-col :span="4"><el-form-item label="内存"><el-input-number v-model="flavorForm.memoryGb" :min="1" style="width: 100%" /></el-form-item></el-col>
            <el-col :span="4"><el-form-item label="系统盘"><el-input-number v-model="flavorForm.storageGb" :min="0" style="width: 100%" /></el-form-item></el-col>
          </el-row>
          <div class="switch-row"><el-switch v-model="flavorForm.isActive" active-text="启用" inactive-text="停用" /></div>
        </template>

        <template v-else>
          <el-row :gutter="12">
            <el-col :span="12"><el-form-item label="镜像编码"><el-input v-model="imageForm.code" /></el-form-item></el-col>
            <el-col :span="12"><el-form-item label="镜像名称"><el-input v-model="imageForm.name" /></el-form-item></el-col>
          </el-row>
          <el-row :gutter="12">
            <el-col :span="8"><el-form-item label="系统类型"><el-input v-model="imageForm.osType" /></el-form-item></el-col>
            <el-col :span="8"><el-form-item label="版本"><el-input v-model="imageForm.version" /></el-form-item></el-col>
            <el-col :span="8"><el-form-item label="架构"><el-input v-model="imageForm.architecture" /></el-form-item></el-col>
          </el-row>
          <el-form-item label="限定地域">
            <el-select v-model="imageForm.regionId" clearable style="width: 100%">
              <el-option v-for="item in regions" :key="item.id" :label="item.name" :value="item.id" />
            </el-select>
          </el-form-item>
          <div class="switch-row">
            <el-switch v-model="imageForm.isPublic" active-text="公开镜像" inactive-text="内部镜像" />
            <el-switch v-model="imageForm.isActive" active-text="启用" inactive-text="停用" />
          </div>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">{{ editingId ? "保存修改" : "确认创建" }}</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<style scoped>
.card-header,
.toolbar-inline {
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
.helper-copy {
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
</style>
