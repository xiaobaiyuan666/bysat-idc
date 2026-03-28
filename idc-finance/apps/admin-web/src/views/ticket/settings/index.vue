<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
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

const router = useRouter();
const loading = ref(false);
const saving = ref(false);
const departments = ref<TicketDepartment[]>([]);
const statistics = ref<TicketStatisticsResponse | null>(null);
const admins = ref<AdminMember[]>([]);

const summaryCards = computed(() => {
  const summary = statistics.value?.summary;
  if (!summary) return [];
  return [
    { title: "工单总数", value: summary.total, hint: "当前系统内全部工单", query: {} as Record<string, string> },
    {
      title: "待客户回复",
      value: summary.waitingCustomer,
      hint: "已进入客户回复等待阶段的工单",
      query: { status: "WAITING_CUSTOMER" } as Record<string, string>
    },
    {
      title: "超时工单",
      value: summary.breached,
      hint: "SLA 已经超时的工单，需要优先跟进",
      query: { priority: "URGENT" } as Record<string, string>
    },
    {
      title: "未分配",
      value: summary.unassigned,
      hint: "还没有被客服接单的工单",
      query: { status: "OPEN" } as Record<string, string>
    },
    {
      title: "已关闭",
      value: summary.closed,
      hint: "已完成或已关闭的工单",
      query: { status: "CLOSED" } as Record<string, string>
    }
  ];
});

const enabledDepartmentCount = computed(() => departments.value.filter(item => item.enabled).length);
const boundAdminCount = computed(() => {
  const ids = new Set<number>();
  departments.value.forEach(item => item.adminIds.forEach(id => ids.add(id)));
  return ids.size;
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
    admins.value = adminItems.filter(item => item.status === "ACTIVE");
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

function openTicketList(query: Record<string, string>) {
  void router.push({
    path: "/tickets/list",
    query
  });
}

function openDepartmentTickets(name: string, status?: string) {
  const query: Record<string, string> = {
    departmentName: name
  };
  if (status) query.status = status;
  openTicketList(query);
}

function openAdminTickets(adminId: number, status?: string) {
  const query: Record<string, string> = {
    assignedAdminId: String(adminId)
  };
  if (status) query.status = status;
  openTicketList(query);
}

async function saveDepartments() {
  const normalized = departments.value.map(item => ({
    ...item,
    name: item.name.trim(),
    description: item.description.trim()
  }));

  if (normalized.some(item => !item.name)) {
    ElMessage.warning("请先补全所有工单部门名称");
    return;
  }

  if (!normalized.some(item => item.enabled && item.isDefault)) {
    ElMessage.warning("请至少保留一个启用中的默认部门");
    return;
  }

  saving.value = true;
  try {
    const result = await updateTicketDepartments(normalized);
    departments.value = result.items.map(item => ({
      ...item,
      adminIds: [...(item.adminIds || [])]
    }));
    ElMessage.success("工单部门配置已保存");
    statistics.value = await fetchTicketStatistics();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "工单部门配置保存失败");
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
    eyebrow="工单"
    title="工单配置"
    subtitle="把部门、默认归类和负责人维护做成后台可操作配置，再直接联到工单中心做接单、转派和积压排查。"
  >
    <template #actions>
      <el-button @click="openTicketList({})">工单中心</el-button>
      <el-button :loading="loading" @click="loadData">刷新</el-button>
      <el-button @click="addDepartment">新增部门</el-button>
      <el-button type="primary" :loading="saving" @click="saveDepartments">保存配置</el-button>
    </template>

    <template #metrics>
      <div class="summary-grid">
        <button
          v-for="card in summaryCards"
          :key="card.title"
          type="button"
          class="summary-card summary-card--button"
          @click="openTicketList(card.query)"
        >
          <span>{{ card.title }}</span>
          <strong>{{ card.value }}</strong>
          <p>{{ card.hint }}</p>
        </button>
        <div class="summary-card">
          <span>启用部门</span>
          <strong>{{ enabledDepartmentCount }}</strong>
          <p>当前可用于客户提单、客服筛单和自动分流的部门数量。</p>
        </div>
        <div class="summary-card">
          <span>参与客服</span>
          <strong>{{ boundAdminCount }}</strong>
          <p>当前已绑定到工单部门、可参与接单处理的管理员数量。</p>
        </div>
      </div>
    </template>

    <div class="settings-layout" v-loading="loading">
      <div class="page-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">工单部门</h2>
            <p class="page-subtitle">
              部门用于后台筛单、客户前台提单归类、负责人约束和后续自动化规则。这里先把最核心的部门档案和负责人绑定做实。
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
              <el-input v-model="row.name" placeholder="例如：技术支持" />
            </template>
          </el-table-column>
          <el-table-column label="说明" min-width="240">
            <template #default="{ row }">
              <el-input v-model="row.description" placeholder="说明该部门主要处理哪些问题" />
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
          <el-table-column label="操作" min-width="160" fixed="right">
            <template #default="{ row, $index }">
              <el-button link type="primary" @click="openDepartmentTickets(row.name)">查看工单</el-button>
              <el-button link type="danger" @click="removeDepartment($index)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="stats-grid">
        <div class="page-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">部门负载</h2>
              <p class="page-subtitle">按部门查看当前工单量、处理状态和超时情况，便于部门排班与转派。</p>
            </div>
          </div>

          <el-table :data="statistics?.departmentStats || []" border stripe empty-text="暂无统计数据">
            <el-table-column prop="name" label="部门" min-width="150">
              <template #default="{ row }">
                <el-button link type="primary" @click="openDepartmentTickets(row.name)">{{ row.name }}</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="total" label="总数" width="90" />
            <el-table-column prop="open" label="待处理" width="90" />
            <el-table-column prop="processing" label="处理中" width="90" />
            <el-table-column prop="waitingCustomer" label="待客户回复" width="120" />
            <el-table-column prop="breached" label="超时" width="90" />
            <el-table-column prop="closed" label="已关闭" width="90" />
            <el-table-column label="操作" min-width="180" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openDepartmentTickets(row.name, 'OPEN')">待处理</el-button>
                <el-button link type="primary" @click="openDepartmentTickets(row.name, 'WAITING_CUSTOMER')">
                  待客户回复
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="page-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">客服负载</h2>
              <p class="page-subtitle">按负责人查看接单量和当前积压情况，便于接单、转派和班次协调。</p>
            </div>
          </div>

          <el-table :data="statistics?.adminStats || []" border stripe empty-text="暂无统计数据">
            <el-table-column prop="adminName" label="负责人" min-width="150">
              <template #default="{ row }">
                <el-button link type="primary" @click="openAdminTickets(row.adminId)">{{ row.adminName }}</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="total" label="总数" width="90" />
            <el-table-column prop="open" label="待处理" width="90" />
            <el-table-column prop="processing" label="处理中" width="90" />
            <el-table-column prop="waitingCustomer" label="待客户回复" width="120" />
            <el-table-column prop="breached" label="超时" width="90" />
            <el-table-column prop="closed" label="已关闭" width="90" />
            <el-table-column label="操作" min-width="180" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openAdminTickets(row.adminId, 'OPEN')">待处理</el-button>
                <el-button link type="primary" @click="openAdminTickets(row.adminId, 'PROCESSING')">处理中</el-button>
              </template>
            </el-table-column>
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
  grid-template-columns: repeat(7, minmax(0, 1fr));
  gap: 14px;
}

.summary-card {
  border-radius: 18px;
  padding: 18px 20px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(246, 248, 252, 0.98));
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 14px 40px rgba(15, 23, 42, 0.06);
}

.summary-card--button {
  cursor: pointer;
  text-align: left;
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

@media (max-width: 1440px) {
  .summary-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
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
