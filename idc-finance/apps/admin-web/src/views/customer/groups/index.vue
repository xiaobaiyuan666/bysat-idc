<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createCustomerGroup,
  deleteCustomerGroup,
  fetchCustomerGroups,
  updateCustomerGroup,
  type CustomerGroupRecord
} from "@/api/admin";

const loading = ref(false);
const submitting = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const items = ref<CustomerGroupRecord[]>([]);

const form = reactive({
  name: "",
  description: ""
});

const usedCount = computed(() => items.value.filter(item => item.customerCount > 0).length);

function resetForm() {
  editingId.value = null;
  form.name = "";
  form.description = "";
}

async function loadGroups() {
  loading.value = true;
  try {
    items.value = await fetchCustomerGroups();
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(item: CustomerGroupRecord) {
  editingId.value = item.id;
  form.name = item.name;
  form.description = item.description;
  dialogVisible.value = true;
}

async function submitForm() {
  submitting.value = true;
  try {
    if (editingId.value) {
      await updateCustomerGroup(editingId.value, form);
      ElMessage.success("客户分组已更新");
    } else {
      await createCustomerGroup(form);
      ElMessage.success("客户分组已创建");
    }
    dialogVisible.value = false;
    resetForm();
    await loadGroups();
  } finally {
    submitting.value = false;
  }
}

async function removeGroup(item: CustomerGroupRecord) {
  await ElMessageBox.confirm(`确认删除客户分组“${item.name}”吗？`, "删除客户分组", {
    type: "warning"
  });
  await deleteCustomerGroup(item.id);
  ElMessage.success("客户分组已删除");
  await loadGroups();
}

onMounted(() => {
  void loadGroups();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="客户"
      title="客户分组与折扣"
      subtitle="维护客户分组、营销分层和专属策略入口，让这页真正可编辑、可运营。"
    >
      <template #actions>
        <el-button @click="loadGroups">刷新</el-button>
        <el-button type="primary" @click="openCreate">新建分组</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>分组总数</span><strong>{{ items.length }}</strong></div>
          <div class="summary-pill"><span>已绑定客户</span><strong>{{ usedCount }}</strong></div>
          <div class="summary-pill"><span>空闲分组</span><strong>{{ items.length - usedCount }}</strong></div>
        </div>
      </template>

      <el-table :data="items" border stripe empty-text="暂无客户分组">
        <el-table-column prop="name" label="分组名称" min-width="180" />
        <el-table-column prop="description" label="说明" min-width="320">
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
              @click="removeGroup(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-dialog
        v-model="dialogVisible"
        :title="editingId ? '编辑客户分组' : '新建客户分组'"
        width="520px"
      >
        <el-form label-position="top">
          <el-form-item label="分组名称">
            <el-input v-model="form.name" maxlength="64" />
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
