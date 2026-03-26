<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchAdmins,
  fetchTicketDepartments,
  fetchTicketStatistics,
  updateTicketDepartments,
  type TicketDepartment,
  type TicketStatisticsResponse
} from "@/api/admin";

type AdminMember = {
  id: number;
  username: string;
  displayName: string;
  status: string;
  roles: string[];
};

const loading = ref(false);
const saving = ref(false);
const departments = ref<TicketDepartment[]>([]);
const statistics = ref<TicketStatisticsResponse | null>(null);
const admins = ref<AdminMember[]>([]);

const summaryCards = computed(() => {
  const summary = statistics.value?.summary;
  if (!summary) return [];
  return [
    { title: "\u5de5\u5355\u603b\u6570", value: summary.total, hint: "\u5f53\u524d\u7cfb\u7edf\u5185\u5168\u90e8\u5de5\u5355" },
    {
      title: "\u5f85\u5ba2\u6237\u56de\u590d",
      value: summary.waitingCustomer,
      hint: "\u5df2\u8fdb\u5165\u5ba2\u6237\u56de\u590d\u7b49\u5f85\u9636\u6bb5\u7684\u5de5\u5355"
    },
    { title: "\u8d85\u65f6\u5de5\u5355", value: summary.breached, hint: "SLA \u5df2\u7ecf\u8d85\u65f6\u7684\u5de5\u5355" },
    { title: "\u672a\u5206\u914d", value: summary.unassigned, hint: "\u8fd8\u6ca1\u6709\u88ab\u5ba2\u670d\u63a5\u5355\u7684\u5de5\u5355" },
    { title: "\u5df2\u5173\u95ed", value: summary.closed, hint: "\u5df2\u5b8c\u6210\u6216\u5df2\u5173\u95ed\u7684\u5de5\u5355" }
  ];
});

async function loadData() {
  loading.value = true;
  try {
    const [departmentResult, stats, adminItems] = await Promise.all([
      fetchTicketDepartments(),
      fetchTicketStatistics(),
      fetchAdmins()
    ]);
    departments.value = departmentResult.items.map(item => ({
      ...item,
      adminIds: [...(item.adminIds || [])]
    }));
    statistics.value = stats;
    admins.value = adminItems;
  } finally {
    loading.value = false;
  }
}

function addDepartment() {
  departments.value.push({
    key: `department-${Date.now()}`,
    name: "",
    description: "",
    enabled: true,
    isDefault: departments.value.length === 0,
    sort: (departments.value.length + 1) * 10,
    adminIds: []
  });
}

function removeDepartment(index: number) {
  departments.value.splice(index, 1);
  if (!departments.value.some(item => item.isDefault && item.enabled) && departments.value[0]) {
    departments.value[0].enabled = true;
    departments.value[0].isDefault = true;
  }
}

function setDefault(index: number) {
  departments.value.forEach((item, currentIndex) => {
    item.isDefault = currentIndex === index;
    if (currentIndex === index) {
      item.enabled = true;
    }
  });
}

async function saveDepartments() {
  saving.value = true;
  try {
    const result = await updateTicketDepartments(departments.value);
    departments.value = result.items.map(item => ({
      ...item,
      adminIds: [...(item.adminIds || [])]
    }));
    ElMessage.success("\u5de5\u5355\u90e8\u95e8\u914d\u7f6e\u5df2\u4fdd\u5b58");
    statistics.value = await fetchTicketStatistics();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "\u5de5\u5355\u90e8\u95e8\u914d\u7f6e\u4fdd\u5b58\u5931\u8d25");
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  void loadData();
});
</script>

<template>
  <PageWorkbench
    eyebrow="财务 / 工单"
    title="工单配置"
    subtitle="先把部门、默认归类和负责管理员做成后台可编辑配置，再配合工单中心和工单工作台完成接单、转派和规则约束。"
  >
    <template #actions>
      <el-button :loading="loading" @click="loadData">刷新</el-button>
      <el-button @click="addDepartment">新增部门</el-button>
      <el-button type="primary" :loading="saving" @click="saveDepartments">保存配置</el-button>
    </template>

    <template #metrics>
      <div class="summary-grid">
        <div v-for="card in summaryCards" :key="card.title" class="summary-card">
          <span>{{ card.title }}</span>
          <strong>{{ card.value }}</strong>
          <p>{{ card.hint }}</p>
        </div>
      </div>
    </template>

    <div class="settings-layout" v-loading="loading">
      <div class="page-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">工单部门</h2>
            <p class="page-subtitle">
              部门会用于后台筛单、客户前台提单归类、负责人约束和后续自动化规则。这里先把最核心的部门档案和负责人绑定做实。
            </p>
          </div>
        </div>

        <el-table :data="departments" border stripe empty-text="暂无工单部门">
          <el-table-column label="排序" width="100">
            <template #default="{ row }">
              <el-input-number v-model="row.sort" :min="1" :precision="0" controls-position="right" style="width: 100%" />
            </template>
          </el-table-column>
          <el-table-column label="部门名称" min-width="180">
            <template #default="{ row }">
              <el-input v-model="row.name" placeholder="例如：Technical Support" />
            </template>
          </el-table-column>
          <el-table-column label="说明" min-width="240">
            <template #default="{ row }">
              <el-input v-model="row.description" placeholder="说明这个部门主要处理哪些问题" />
            </template>
          </el-table-column>
          <el-table-column label="负责管理员" min-width="260">
            <template #default="{ row }">
              <el-select v-model="row.adminIds" multiple collapse-tags collapse-tags-tooltip style="width: 100%">
                <el-option
                  v-for="item in admins"
                  :key="item.id"
                  :label="item.displayName || item.username"
                  :value="item.id"
                />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="启用" width="90">
            <template #default="{ row }">
              <el-switch v-model="row.enabled" />
            </template>
          </el-table-column>
          <el-table-column label="默认" width="110">
            <template #default="{ row, $index }">
              <el-button size="small" :type="row.isDefault ? 'primary' : 'default'" @click="setDefault($index)">
                默认
              </el-button>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="90" fixed="right">
            <template #default="{ $index }">
              <el-button type="danger" link @click="removeDepartment($index)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="stats-grid">
        <div class="page-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">部门负载</h2>
              <p class="page-subtitle">按部门查看当前工单量、处理状态和超时情况。</p>
            </div>
          </div>

          <el-table :data="statistics?.departmentStats || []" border stripe empty-text="暂无统计数据">
            <el-table-column prop="name" label="部门" min-width="150" />
            <el-table-column prop="total" label="总数" width="90" />
            <el-table-column prop="open" label="待处理" width="90" />
            <el-table-column prop="processing" label="处理中" width="90" />
            <el-table-column prop="waitingCustomer" label="待客户回复" width="120" />
            <el-table-column prop="breached" label="超时" width="90" />
            <el-table-column prop="closed" label="已关闭" width="90" />
          </el-table>
        </div>

        <div class="page-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">客服负载</h2>
              <p class="page-subtitle">按负责人查看接单量和当前积压情况，便于接单与转派。</p>
            </div>
          </div>

          <el-table :data="statistics?.adminStats || []" border stripe empty-text="暂无统计数据">
            <el-table-column prop="adminName" label="负责人" min-width="150" />
            <el-table-column prop="total" label="总数" width="90" />
            <el-table-column prop="open" label="待处理" width="90" />
            <el-table-column prop="processing" label="处理中" width="90" />
            <el-table-column prop="waitingCustomer" label="待客户回复" width="120" />
            <el-table-column prop="breached" label="超时" width="90" />
            <el-table-column prop="closed" label="已关闭" width="90" />
          </el-table>
        </div>
      </div>
    </div>
  </PageWorkbench>
</template>

<style scoped>
.settings-layout {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 14px;
}

.summary-card {
  border-radius: 18px;
  padding: 18px 20px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(246, 248, 252, 0.98));
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 14px 40px rgba(15, 23, 42, 0.06);
}

.summary-card span {
  display: block;
  font-size: 13px;
  color: #64748b;
}

.summary-card strong {
  display: block;
  margin-top: 10px;
  font-size: 28px;
  color: #0f172a;
}

.summary-card p {
  margin: 10px 0 0;
  line-height: 1.6;
  color: #475569;
  font-size: 13px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

@media (max-width: 1280px) {
  .summary-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }
}
</style>
