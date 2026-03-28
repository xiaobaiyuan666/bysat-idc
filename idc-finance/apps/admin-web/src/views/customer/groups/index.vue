<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createCustomerGroup,
  deleteCustomerGroup,
  fetchCustomerGroups,
  updateCustomerGroup,
  type CustomerGroupRecord
} from "@/api/admin";

const router = useRouter();
const loading = ref(false);
const submitting = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const items = ref<CustomerGroupRecord[]>([]);

const form = reactive({
  name: "",
  description: ""
});

const sortedItems = computed(() =>
  [...items.value].sort((left, right) => {
    if (right.customerCount !== left.customerCount) {
      return right.customerCount - left.customerCount;
    }
    return left.name.localeCompare(right.name, "zh-CN");
  })
);

const usedCount = computed(() => items.value.filter(item => item.customerCount > 0).length);
const totalCustomerCount = computed(() =>
  items.value.reduce((total, item) => total + item.customerCount, 0)
);
const largestGroup = computed(() => sortedItems.value[0]);

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

function openCustomers(item: CustomerGroupRecord) {
  void router.push({
    path: "/customer/list",
    query: {
      groupName: item.name
    }
  });
}

async function submitForm() {
  if (!form.name.trim()) {
    ElMessage.warning("请输入客户分组名称");
    return;
  }

  submitting.value = true;
  try {
    const payload = {
      name: form.name.trim(),
      description: form.description.trim()
    };
    if (editingId.value) {
      await updateCustomerGroup(editingId.value, payload);
      ElMessage.success("客户分组已更新");
    } else {
      await createCustomerGroup(payload);
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
      title="客户分组"
      subtitle="维护客户分组、分层说明和客户归属入口，让分组页面不再只是静态配置表。"
    >
      <template #actions>
        <el-button @click="loadGroups">刷新</el-button>
        <el-button type="primary" @click="openCreate">新建分组</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>分组总数</span><strong>{{ items.length }}</strong></div>
          <div class="summary-pill"><span>已绑定分组</span><strong>{{ usedCount }}</strong></div>
          <div class="summary-pill"><span>覆盖客户数</span><strong>{{ totalCustomerCount }}</strong></div>
          <div class="summary-pill">
            <span>最大分组</span>
            <strong>{{ largestGroup ? largestGroup.name : "暂无" }}</strong>
          </div>
        </div>
      </template>

      <el-table :data="sortedItems" border stripe empty-text="暂无客户分组">
        <el-table-column prop="name" label="分组名称" min-width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="openCustomers(row)">{{ row.name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="customerCount" label="绑定客户" min-width="120">
          <template #default="{ row }">
            <el-tag :type="row.customerCount > 0 ? 'success' : 'info'" effect="light">
              {{ row.customerCount }} 位
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="分组说明" min-width="360">
          <template #default="{ row }">
            {{ row.description || "暂无说明" }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="220" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openCustomers(row)">查看客户</el-button>
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
          <el-form-item label="分组名称" required>
            <el-input v-model="form.name" maxlength="64" show-word-limit />
          </el-form-item>
          <el-form-item label="分组说明">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="4"
              maxlength="255"
              show-word-limit
            />
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
