<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createCustomerLevel,
  deleteCustomerLevel,
  fetchCustomerLevels,
  updateCustomerLevel,
  type CustomerLevelRecord
} from "@/api/admin";

const loading = ref(false);
const submitting = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const items = ref<CustomerLevelRecord[]>([]);

const form = reactive({
  name: "",
  priority: 50,
  description: ""
});

const highestPriority = computed(() =>
  items.value.length ? Math.max(...items.value.map(item => item.priority)) : 0
);

function resetForm() {
  editingId.value = null;
  form.name = "";
  form.priority = 50;
  form.description = "";
}

async function loadLevels() {
  loading.value = true;
  try {
    items.value = await fetchCustomerLevels();
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(item: CustomerLevelRecord) {
  editingId.value = item.id;
  form.name = item.name;
  form.priority = item.priority;
  form.description = item.description;
  dialogVisible.value = true;
}

async function submitForm() {
  submitting.value = true;
  try {
    if (editingId.value) {
      await updateCustomerLevel(editingId.value, form);
      ElMessage.success("客户等级已更新");
    } else {
      await createCustomerLevel(form);
      ElMessage.success("客户等级已创建");
    }
    dialogVisible.value = false;
    resetForm();
    await loadLevels();
  } finally {
    submitting.value = false;
  }
}

async function removeLevel(item: CustomerLevelRecord) {
  await ElMessageBox.confirm(`确认删除客户等级“${item.name}”吗？`, "删除客户等级", {
    type: "warning"
  });
  await deleteCustomerLevel(item.id);
  ElMessage.success("客户等级已删除");
  await loadLevels();
}

onMounted(() => {
  void loadLevels();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="客户"
      title="客户等级"
      subtitle="维护优先级、服务等级和默认策略，让客户等级真正可配置。"
    >
      <template #actions>
        <el-button @click="loadLevels">刷新</el-button>
        <el-button type="primary" @click="openCreate">新建等级</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>等级总数</span><strong>{{ items.length }}</strong></div>
          <div class="summary-pill"><span>最高优先级</span><strong>{{ highestPriority }}</strong></div>
          <div class="summary-pill">
            <span>已绑定客户</span>
            <strong>{{ items.reduce((sum, item) => sum + item.customerCount, 0) }}</strong>
          </div>
        </div>
      </template>

      <el-table :data="items" border stripe empty-text="暂无客户等级">
        <el-table-column prop="name" label="等级名称" min-width="180" />
        <el-table-column prop="priority" label="优先级" min-width="100" />
        <el-table-column prop="description" label="说明" min-width="280">
          <template #default="{ row }">
            {{ row.description || "暂无说明" }}
          </template>
        </el-table-column>
        <el-table-column prop="customerCount" label="客户数" min-width="100" />
        <el-table-column label="操作" min-width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button
              link
              type="danger"
              :disabled="row.customerCount > 0"
              @click="removeLevel(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-dialog
        v-model="dialogVisible"
        :title="editingId ? '编辑客户等级' : '新建客户等级'"
        width="520px"
      >
        <el-form label-position="top">
          <el-form-item label="等级名称">
            <el-input v-model="form.name" maxlength="64" />
          </el-form-item>
          <el-form-item label="优先级">
            <el-input-number v-model="form.priority" :min="0" :max="999" style="width: 100%" />
          </el-form-item>
          <el-form-item label="说明">
            <el-input v-model="form.description" type="textarea" :rows="3" maxlength="255" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="submitForm">保存</el-button>
        </template>
      </el-dialog>
    </PageWorkbench>
  </div>
</template>
