<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createTicket,
  fetchAdmins,
  fetchCustomers,
  fetchServices,
  fetchTicketDepartments,
  type Customer,
  type ServiceRecord,
  type TicketDepartment
} from "@/api/admin";

type AdminMember = {
  id: number;
  username: string;
  displayName: string;
  status: string;
  roles: string[];
};

const route = useRoute();
const router = useRouter();

const text = {
  eyebrow: "客服 / 工单",
  title: "新建工单",
  subtitle:
    "按照后台代客建单的工作台思路，先选客户，再绑定服务、部门和负责人，最后直接进入工单详情继续跟进。",
  refresh: "刷新基础数据",
  back: "返回工单中心",
  submit: "创建工单",
  customer: "客户",
  service: "关联服务",
  department: "受理部门",
  assignee: "负责人",
  priority: "优先级",
  titleLabel: "工单标题",
  content: "工单内容",
  customerPlaceholder: "选择需要代建单的客户",
  servicePlaceholder: "可不选，用于绑定这次售后来源",
  departmentPlaceholder: "选择工单归属部门",
  assigneePlaceholder: "可不选，后续也可在工单工作台再转交",
  titlePlaceholder: "例如：账单金额需手动调整",
  contentPlaceholder: "写清问题现象、预期结果、处理背景和需要客服重点跟进的说明。",
  createSuccess: "工单已创建",
  selectHint: "这里会显示该客户当前可关联的服务，也支持从服务详情页直接带入。",
  customerNo: "客户编号",
  email: "邮箱",
  group: "分组 / 等级",
  noService: "当前没有可关联的服务，仍可直接建立通用工单",
  defaultBadge: "默认",
  ruleBadge: "限定负责人",
  sourceSummary: "创建后来源会记录为 ADMIN，便于后续审计和追踪。"
} as const;

const loading = ref(false);
const saving = ref(false);
const customers = ref<Customer[]>([]);
const services = ref<ServiceRecord[]>([]);
const departments = ref<TicketDepartment[]>([]);
const admins = ref<AdminMember[]>([]);

const form = reactive({
  customerId: 0,
  serviceId: 0,
  departmentName: "",
  assignedAdminId: 0,
  priority: "NORMAL",
  title: "",
  content: ""
});

const selectedCustomer = computed(() => customers.value.find(item => item.id === form.customerId) || null);

const customerServices = computed(() => services.value.filter(item => item.customerId === form.customerId));

const enabledDepartments = computed(() =>
  departments.value.filter(item => item.enabled).sort((a, b) => a.sort - b.sort)
);

const currentDepartment = computed(
  () => enabledDepartments.value.find(item => item.name === form.departmentName) || null
);

const availableAdmins = computed(() => {
  const limitIds = currentDepartment.value?.adminIds || [];
  if (!limitIds.length) {
    return admins.value;
  }
  return admins.value.filter(item => limitIds.includes(item.id));
});

function assigneeLabel(id: number) {
  if (!id) return "";
  const matched = admins.value.find(item => item.id === id);
  return matched?.displayName || matched?.username || "";
}

function applyRoutePrefill() {
  const customerIdQuery = typeof route.query.customerId === "string" ? Number(route.query.customerId) : NaN;
  const serviceIdQuery = typeof route.query.serviceId === "string" ? Number(route.query.serviceId) : NaN;
  const assignedAdminIdQuery = typeof route.query.assignedAdminId === "string" ? Number(route.query.assignedAdminId) : NaN;
  const titleQuery = typeof route.query.title === "string" ? route.query.title : "";
  const contentQuery = typeof route.query.content === "string" ? route.query.content : "";
  const departmentQuery = typeof route.query.departmentName === "string" ? route.query.departmentName : "";
  const priorityQuery = typeof route.query.priority === "string" ? route.query.priority : "";

  if (Number.isFinite(customerIdQuery) && customerIdQuery > 0) {
    form.customerId = customerIdQuery;
  }

  if (Number.isFinite(serviceIdQuery) && serviceIdQuery > 0) {
    const linkedService = services.value.find(item => item.id === serviceIdQuery);
    if (linkedService) {
      form.serviceId = linkedService.id;
      form.customerId = linkedService.customerId;
    }
  }

  if (departmentQuery && enabledDepartments.value.some(item => item.name === departmentQuery)) {
    form.departmentName = departmentQuery;
  } else if (!form.departmentName) {
    form.departmentName =
      enabledDepartments.value.find(item => item.isDefault)?.name || enabledDepartments.value[0]?.name || "";
  }

  if (priorityQuery && ["LOW", "NORMAL", "HIGH", "URGENT"].includes(priorityQuery)) {
    form.priority = priorityQuery;
  }

  if (titleQuery) {
    form.title = titleQuery;
  }

  if (contentQuery) {
    form.content = contentQuery;
  }

  if (Number.isFinite(assignedAdminIdQuery) && assignedAdminIdQuery > 0) {
    if (availableAdmins.value.some(item => item.id === assignedAdminIdQuery)) {
      form.assignedAdminId = assignedAdminIdQuery;
    }
  }
}

async function loadData() {
  loading.value = true;
  try {
    const [customerResult, serviceResult, departmentResult, adminResult] = await Promise.all([
      fetchCustomers(),
      fetchServices({ page: 1, limit: 200 }),
      fetchTicketDepartments(),
      fetchAdmins()
    ]);
    customers.value = customerResult.items;
    services.value = serviceResult.items;
    departments.value = departmentResult.items;
    admins.value = adminResult;
    applyRoutePrefill();
  } finally {
    loading.value = false;
  }
}

async function submit() {
  if (!form.customerId || !form.title.trim() || !form.content.trim()) {
    ElMessage.warning("请先选择客户并填写标题和内容");
    return;
  }

  saving.value = true;
  try {
    const created = await createTicket({
      customerId: form.customerId,
      serviceId: form.serviceId || undefined,
      departmentName: form.departmentName || undefined,
      assignedAdminId: form.assignedAdminId || undefined,
      assignedAdminName: form.assignedAdminId ? assigneeLabel(form.assignedAdminId) : undefined,
      priority: form.priority,
      title: form.title.trim(),
      content: form.content.trim()
    });
    ElMessage.success(text.createSuccess);
    void router.push(`/tickets/${created.id}`);
  } finally {
    saving.value = false;
  }
}

watch(
  () => form.customerId,
  () => {
    if (form.serviceId && !customerServices.value.some(item => item.id === form.serviceId)) {
      form.serviceId = 0;
    }
  }
);

watch(
  () => form.departmentName,
  () => {
    if (form.assignedAdminId && !availableAdmins.value.some(item => item.id === form.assignedAdminId)) {
      form.assignedAdminId = 0;
    }
  }
);

watch(
  () => [route.query.customerId, route.query.serviceId, route.query.title, route.query.content, route.query.departmentName, route.query.priority, route.query.assignedAdminId],
  () => {
    applyRoutePrefill();
  }
);

onMounted(() => {
  void loadData();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench :eyebrow="text.eyebrow" :title="text.title" :subtitle="text.subtitle">
      <template #actions>
        <el-button @click="loadData">{{ text.refresh }}</el-button>
        <el-button plain @click="router.push('/tickets/list')">{{ text.back }}</el-button>
        <el-button type="primary" :loading="saving" @click="submit">{{ text.submit }}</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>{{ text.customer }}</span>
            <strong>{{ selectedCustomer?.name || "-" }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.department }}</span>
            <strong>{{ form.departmentName || "-" }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.assignee }}</span>
            <strong>{{ assigneeLabel(form.assignedAdminId) || "-" }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.service }}</span>
            <strong>{{ customerServices.length }}</strong>
          </div>
        </div>
      </template>

      <div class="ticket-create-layout">
        <div class="page-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">{{ text.customer }}</h2>
              <p class="page-subtitle">{{ text.selectHint }}</p>
            </div>
          </div>

          <el-form label-position="top" class="form-grid">
            <el-form-item :label="text.customer">
              <el-select v-model="form.customerId" filterable :placeholder="text.customerPlaceholder">
                <el-option
                  v-for="item in customers"
                  :key="item.id"
                  :label="`${item.name} / ${item.customerNo}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item :label="text.service">
              <el-select
                v-model="form.serviceId"
                clearable
                filterable
                :placeholder="text.servicePlaceholder"
                :disabled="!form.customerId"
              >
                <el-option
                  v-for="item in customerServices"
                  :key="item.id"
                  :label="`${item.serviceNo} / ${item.productName}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item :label="text.department">
              <el-select v-model="form.departmentName" :placeholder="text.departmentPlaceholder">
                <el-option v-for="item in enabledDepartments" :key="item.key" :label="item.name" :value="item.name">
                  <div class="department-option">
                    <span>{{ item.name }}</span>
                    <div class="department-option__badges">
                      <el-tag v-if="item.isDefault" size="small" effect="plain">{{ text.defaultBadge }}</el-tag>
                      <el-tag v-if="item.adminIds.length" size="small" type="warning" effect="plain">
                        {{ text.ruleBadge }}
                      </el-tag>
                    </div>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>

            <el-form-item :label="text.assignee">
              <el-select v-model="form.assignedAdminId" clearable filterable :placeholder="text.assigneePlaceholder">
                <el-option :label="'-'" :value="0" />
                <el-option
                  v-for="item in availableAdmins"
                  :key="item.id"
                  :label="item.displayName || item.username"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item :label="text.priority">
              <el-select v-model="form.priority">
                <el-option label="低" value="LOW" />
                <el-option label="普通" value="NORMAL" />
                <el-option label="高" value="HIGH" />
                <el-option label="紧急" value="URGENT" />
              </el-select>
            </el-form-item>

            <el-form-item class="form-span-2" :label="text.titleLabel">
              <el-input v-model="form.title" maxlength="120" show-word-limit :placeholder="text.titlePlaceholder" />
            </el-form-item>

            <el-form-item class="form-span-2" :label="text.content">
              <el-input
                v-model="form.content"
                type="textarea"
                :rows="8"
                maxlength="5000"
                show-word-limit
                :placeholder="text.contentPlaceholder"
              />
            </el-form-item>
          </el-form>
        </div>

        <div class="page-card side-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">{{ text.customer }}</h2>
              <p class="page-subtitle">{{ text.sourceSummary }}</p>
            </div>
          </div>

          <div class="info-list">
            <div>
              <span>{{ text.customerNo }}</span>
              <strong>{{ selectedCustomer?.customerNo || "-" }}</strong>
            </div>
            <div>
              <span>{{ text.email }}</span>
              <strong>{{ selectedCustomer?.email || "-" }}</strong>
            </div>
            <div>
              <span>{{ text.group }}</span>
              <strong>
                {{ selectedCustomer ? `${selectedCustomer.groupName || "-"} / ${selectedCustomer.levelName || "-"}` : "-" }}
              </strong>
            </div>
          </div>

          <div class="service-board">
            <div class="service-board__head">
              <strong>{{ text.service }}</strong>
            </div>
            <div v-if="customerServices.length" class="service-list">
              <div v-for="item in customerServices" :key="item.id" class="service-item">
                <strong>{{ item.serviceNo }}</strong>
                <span>{{ item.productName }}</span>
                <small>{{ item.status }}</small>
              </div>
            </div>
            <el-empty v-else :description="text.noService" :image-size="72" />
          </div>
        </div>
      </div>
    </PageWorkbench>
  </div>
</template>

<style scoped>
.ticket-create-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.8fr) minmax(320px, 0.9fr);
  gap: 18px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.form-span-2 {
  grid-column: 1 / -1;
}

.department-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.department-option__badges {
  display: inline-flex;
  gap: 6px;
}

.side-card {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.info-list {
  display: grid;
  gap: 14px;
}

.info-list div {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.18);
}

.info-list span {
  font-size: 13px;
  color: #64748b;
}

.info-list strong {
  color: #0f172a;
  font-size: 15px;
}

.service-board {
  padding: 16px;
  border-radius: 16px;
  background: linear-gradient(180deg, rgba(248, 250, 252, 0.98), rgba(241, 245, 249, 0.94));
  border: 1px solid rgba(148, 163, 184, 0.18);
}

.service-board__head {
  margin-bottom: 12px;
  color: #0f172a;
}

.service-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.service-item {
  display: grid;
  gap: 4px;
  padding: 12px 14px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid rgba(148, 163, 184, 0.16);
}

.service-item strong {
  color: #0f172a;
}

.service-item span,
.service-item small {
  color: #475569;
}

@media (max-width: 1280px) {
  .ticket-create-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .form-span-2 {
    grid-column: auto;
  }
}
</style>
