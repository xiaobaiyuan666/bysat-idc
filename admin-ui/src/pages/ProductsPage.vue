<script setup lang="ts">
import { Edit, Plus } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { onMounted, reactive, ref } from "vue";

import { http } from "@/api/http";
import { formatCurrency } from "@/utils/format";
import {
  cycleMap,
  getLabel,
  getStatusTagType,
  productCategoryMap,
  productStatusMap,
  providerTypeMap,
} from "@/utils/maps";

const loading = ref(false);
const dialogVisible = ref(false);
const editingId = ref("");
const products = ref<any[]>([]);

const form = reactive({
  code: "",
  name: "",
  category: "CLOUD_SERVER",
  status: "ACTIVE",
  billingCycle: "MONTHLY",
  price: 99,
  setupFee: 0,
  stock: 50,
  autoProvision: true,
  providerType: "MOFANG_CLOUD",
  providerProductId: "",
  regionTemplate: "",
  description: "",
});

function resetForm() {
  Object.assign(form, {
    code: "",
    name: "",
    category: "CLOUD_SERVER",
    status: "ACTIVE",
    billingCycle: "MONTHLY",
    price: 99,
    setupFee: 0,
    stock: 50,
    autoProvision: true,
    providerType: "MOFANG_CLOUD",
    providerProductId: "",
    regionTemplate: "",
    description: "",
  });
}

function openCreate() {
  editingId.value = "";
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: any) {
  editingId.value = row.id;
  Object.assign(form, {
    code: row.code,
    name: row.name,
    category: row.category,
    status: row.status,
    billingCycle: row.billingCycle,
    price: row.price / 100,
    setupFee: row.setupFee / 100,
    stock: row.stock,
    autoProvision: row.autoProvision,
    providerType: row.providerType,
    providerProductId: row.providerProductId || "",
    regionTemplate: row.regionTemplate || "",
    description: row.description || "",
  });
  dialogVisible.value = true;
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/products");
    products.value = data.data;
  } finally {
    loading.value = false;
  }
}

async function submitForm() {
  if (editingId.value) {
    await http.put(`/products/${editingId.value}`, form);
    ElMessage.success("产品已更新");
  } else {
    await http.post("/products", form);
    ElMessage.success("产品已创建");
  }

  dialogVisible.value = false;
  resetForm();
  await loadData();
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">产品管理</h1>
        <div class="page-subtitle">
          产品定义是计费、售卖方案、服务实例和魔方云对接的基础层。这里改成完整可编辑，不再只支持新增。
        </div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">新增产品</el-button>
    </div>

    <el-card class="page-card">
      <el-table v-loading="loading" :data="products" stripe>
        <el-table-column label="产品信息" min-width="240">
          <template #default="{ row }">
            <div style="font-weight: 700">{{ row.name }}</div>
            <div class="muted-line">{{ row.code }} / {{ getLabel(providerTypeMap, row.providerType) }}</div>
          </template>
        </el-table-column>
        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            {{ getLabel(productCategoryMap, row.category) }}
          </template>
        </el-table-column>
        <el-table-column label="周期" width="110">
          <template #default="{ row }">
            {{ getLabel(cycleMap, row.billingCycle) }}
          </template>
        </el-table-column>
        <el-table-column label="售价" width="120">
          <template #default="{ row }">
            {{ formatCurrency(row.price) }}
          </template>
        </el-table-column>
        <el-table-column label="开通费" width="120">
          <template #default="{ row }">
            {{ formatCurrency(row.setupFee) }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="90" />
        <el-table-column label="方案 / 服务" width="110">
          <template #default="{ row }">
            {{ row._count.plans }} / {{ row._count.services }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getLabel(productStatusMap, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="110" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" :icon="Edit" @click="openEdit(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑产品' : '新增产品'" width="680px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="产品编码">
              <el-input v-model="form.code" placeholder="例如 CLOUD-GEN2-2C4G" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="产品名称">
              <el-input v-model="form.name" placeholder="请输入产品名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="产品分类">
              <el-select v-model="form.category" style="width: 100%">
                <el-option
                  v-for="(label, value) in productCategoryMap"
                  :key="value"
                  :label="label"
                  :value="value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="销售状态">
              <el-select v-model="form.status" style="width: 100%">
                <el-option
                  v-for="(label, value) in productStatusMap"
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
          <el-col :span="12">
            <el-form-item label="云平台类型">
              <el-select v-model="form.providerType" style="width: 100%">
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
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="销售价">
              <el-input-number v-model="form.price" :precision="2" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="开通费">
              <el-input-number
                v-model="form.setupFee"
                :precision="2"
                :min="0"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="库存">
              <el-input-number v-model="form.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="云平台商品 ID">
              <el-input v-model="form.providerProductId" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="默认地域模板">
              <el-input v-model="form.regionTemplate" placeholder="例如 cn-hk-1" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="产品说明">
          <el-input v-model="form.description" type="textarea" :rows="4" />
        </el-form-item>
        <el-switch
          v-model="form.autoProvision"
          active-text="支付后自动开通"
          inactive-text="人工审核开通"
        />
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">{{ editingId ? "保存修改" : "确认创建" }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.muted-line {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 12px;
}
</style>
