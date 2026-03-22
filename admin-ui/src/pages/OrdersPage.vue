<script setup lang="ts">
import {
  Plus,
  Refresh,
  Search,
  DocumentCopy,
  Download,
} from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";

import ListWorkbenchShell from "@/components/common/ListWorkbenchShell.vue";
import { http } from "@/api/http";
import { downloadCsv } from "@/utils/download";
import { formatCurrency, formatDate, formatDateTime } from "@/utils/format";
import {
  cycleMap,
  getLabel,
  getStatusTagType,
  orderSourceMap,
  orderStatusMap,
} from "@/utils/maps";

const router = useRouter();

const loading = ref(false);
const dialogVisible = ref(false);
const advancedOpen = ref(false);
const statusFilter = ref("ALL");
const keyword = ref("");
const orders = ref<any[]>([]);
const customers = ref<any[]>([]);
const products = ref<any[]>([]);
const selectedRows = ref<any[]>([]);

const filters = reactive({
  source: "ALL",
  customerId: "ALL",
  cycle: "ALL",
  amountMin: undefined as number | undefined,
  amountMax: undefined as number | undefined,
  dueRange: [] as [Date, Date] | [],
});

const form = reactive({
  customerId: "",
  productId: "",
  serviceName: "",
  cycle: "MONTHLY",
  quantity: 1,
  notes: "",
});

const statusCards = computed(() => {
  const total = orders.value.length;
  const pending = orders.value.filter((item) => item.status === "PENDING").length;
  const active = orders.value.filter((item) => item.status === "ACTIVE").length;
  const provisioning = orders.value.filter((item) => item.status === "PROVISIONING").length;

  return [
    { key: "ALL", label: "全部订单", value: total, hint: "系统内所有订单" },
    { key: "PENDING", label: "待支付", value: pending, hint: "等待付款或审核" },
    { key: "PROVISIONING", label: "开通中", value: provisioning, hint: "正在交付资源" },
    { key: "ACTIVE", label: "已生效", value: active, hint: "已交付完成" },
  ];
});

const summaryTotals = computed(() => {
  const totalAmount = orders.value.reduce((sum, item) => sum + item.totalAmount, 0);
  const totalPaid = orders.value.reduce((sum, item) => sum + item.paidAmount, 0);

  return [
    {
      label: "订单总额",
      value: formatCurrency(totalAmount),
      hint: "全部订单金额",
    },
    {
      label: "已收金额",
      value: formatCurrency(totalPaid),
      hint: "订单累计入账",
    },
  ];
});

const filteredOrders = computed(() => {
  const normalizedKeyword = keyword.value.trim().toLowerCase();
  const [dueFrom, dueTo] = filters.dueRange;

  return orders.value.filter((item) => {
    const matchStatus = statusFilter.value === "ALL" || item.status === statusFilter.value;
    const matchSource = filters.source === "ALL" || item.source === filters.source;
    const matchCustomer = filters.customerId === "ALL" || item.customerId === filters.customerId;
    const matchCycle =
      filters.cycle === "ALL" || item.items.some((row: any) => row.cycle === filters.cycle);
    const matchKeyword =
      normalizedKeyword.length === 0 ||
      [
        item.orderNo,
        item.customer?.name,
        item.notes,
        ...item.items.map((row: any) => row.title),
      ]
        .some((field) => String(field ?? "").toLowerCase().includes(normalizedKeyword));

    const totalYuan = item.totalAmount / 100;
    const matchAmountMin =
      typeof filters.amountMin !== "number" || totalYuan >= filters.amountMin;
    const matchAmountMax =
      typeof filters.amountMax !== "number" || totalYuan <= filters.amountMax;

    const dueDate = item.dueDate ? new Date(item.dueDate) : null;
    const matchDueFrom = !dueFrom || !dueDate || dueDate >= dueFrom;
    const matchDueTo = !dueTo || !dueDate || dueDate <= dueTo;

    return (
      matchStatus &&
      matchSource &&
      matchCustomer &&
      matchCycle &&
      matchKeyword &&
      matchAmountMin &&
      matchAmountMax &&
      matchDueFrom &&
      matchDueTo
    );
  });
});

const pendingDeliveryCount = computed(
  () => filteredOrders.value.filter((item) => item.status === "PROVISIONING").length,
);

function syncCycleWithProduct(productId: string) {
  const product = products.value.find((item) => item.id === productId);

  if (product?.billingCycle) {
    form.cycle = product.billingCycle;
  }
}

function resetAdvancedFilters() {
  Object.assign(filters, {
    source: "ALL",
    customerId: "ALL",
    cycle: "ALL",
    amountMin: undefined,
    amountMax: undefined,
    dueRange: [],
  });
}

function handleSelectionChange(rows: any[]) {
  selectedRows.value = rows;
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/orders");
    orders.value = data.data.orders;
    customers.value = data.data.customers;
    products.value = data.data.products;

    if (!form.customerId && customers.value[0]) {
      form.customerId = customers.value[0].id;
    }
    if (!form.productId && products.value[0]) {
      form.productId = products.value[0].id;
      syncCycleWithProduct(products.value[0].id);
    }
  } finally {
    loading.value = false;
  }
}

watch(
  () => form.productId,
  (value) => {
    syncCycleWithProduct(value);
  },
);

async function submitForm() {
  await http.post("/orders", form);
  ElMessage.success("订单创建成功");
  dialogVisible.value = false;
  form.serviceName = "";
  form.quantity = 1;
  form.notes = "";
  await loadData();
}

async function copySelectedOrderNos() {
  const orderNos = selectedRows.value.map((item) => item.orderNo).join("\n");
  await navigator.clipboard.writeText(orderNos);
  ElMessage.success("已复制所选订单号");
}

function exportOrders() {
  const rows = (selectedRows.value.length > 0 ? selectedRows.value : filteredOrders.value).map(
    (item) => [
      item.orderNo,
      item.customer?.name,
      getLabel(orderSourceMap, item.source),
      item.items.map((row: any) => row.title).join(" / "),
      formatCurrency(item.totalAmount),
      formatCurrency(item.paidAmount),
      formatDate(item.dueDate),
      getLabel(orderStatusMap, item.status),
      formatDateTime(item.createdAt),
    ],
  );

  downloadCsv(
    `orders-${new Date().toISOString().slice(0, 10)}.csv`,
    ["订单号", "客户", "来源", "订单项目", "订单金额", "已收金额", "到期日", "状态", "创建时间"],
    rows,
  );
  ElMessage.success("订单数据已导出");
}

onMounted(loadData);
</script>

<template>
  <ListWorkbenchShell
    v-model:advancedOpen="advancedOpen"
    title="订单中心"
    subtitle="订单中心负责新购、续费和后台代下单，订单会自动联动账单和服务实例，核心关系已经按“订单 -> 账单 -> 服务”重构。"
    :result-count="filteredOrders.length"
    :selected-count="selectedRows.length"
    :show-advanced-toggle="true"
    advanced-button-text="高级筛选"
  >
    <template #actions>
      <el-button :icon="Refresh" @click="loadData">刷新数据</el-button>
      <el-button type="primary" :icon="Plus" @click="dialogVisible = true">新增订单</el-button>
    </template>

    <template #summary>
      <button
        v-for="card in statusCards"
        :key="card.key"
        class="status-tab"
        :class="{ 'is-active': statusFilter === card.key }"
        @click="statusFilter = card.key"
      >
        <div class="status-tab-label">{{ card.label }}</div>
        <div class="status-tab-value">{{ card.value }}</div>
        <div class="muted-line">{{ card.hint }}</div>
      </button>

      <div v-for="card in summaryTotals" :key="card.label" class="summary-card">
        <div class="status-tab-label">{{ card.label }}</div>
        <h3>{{ card.value }}</h3>
        <p>{{ card.hint }}</p>
      </div>
    </template>

    <template #filters>
      <el-input
        v-model="keyword"
        clearable
        placeholder="搜索订单号、客户名、订单备注、订单项目"
        style="width: 360px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <el-select v-model="filters.source" style="width: 180px">
        <el-option label="全部来源" value="ALL" />
        <el-option
          v-for="(label, value) in orderSourceMap"
          :key="value"
          :label="label"
          :value="value"
        />
      </el-select>
    </template>

    <template #filterMeta>
      <el-tag effect="plain">当前结果 {{ filteredOrders.length }} 条</el-tag>
      <el-tag effect="plain">待交付 {{ pendingDeliveryCount }} 单</el-tag>
    </template>

    <template #advanced>
      <div class="advanced-filter-grid">
        <el-select v-model="filters.customerId" filterable>
          <el-option label="全部客户" value="ALL" />
          <el-option
            v-for="item in customers"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          />
        </el-select>
        <el-select v-model="filters.cycle">
          <el-option label="全部周期" value="ALL" />
          <el-option
            v-for="(label, value) in cycleMap"
            :key="value"
            :label="label"
            :value="value"
          />
        </el-select>
        <el-input-number
          v-model="filters.amountMin"
          :precision="2"
          :min="0"
          style="width: 100%"
          placeholder="最小金额"
        />
        <el-input-number
          v-model="filters.amountMax"
          :precision="2"
          :min="0"
          style="width: 100%"
          placeholder="最大金额"
        />
        <el-date-picker
          v-model="filters.dueRange"
          type="daterange"
          range-separator="至"
          start-placeholder="到期开始"
          end-placeholder="到期结束"
          style="width: 100%"
        />
        <div class="inline-actions">
          <el-button @click="resetAdvancedFilters">清空高级筛选</el-button>
        </div>
      </div>
    </template>

    <template #selection>
      <el-button plain :icon="DocumentCopy" @click="copySelectedOrderNos">复制订单号</el-button>
      <el-button plain :icon="Download" @click="exportOrders">导出所选订单</el-button>
    </template>

    <el-table
      v-loading="loading"
      :data="filteredOrders"
      stripe
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column label="订单信息" min-width="240">
        <template #default="{ row }">
          <div class="primary-line">{{ row.orderNo }}</div>
          <div class="muted-line">
            {{ row.customer.name }} / {{ getLabel(orderSourceMap, row.source) }}
          </div>
          <div class="muted-line">{{ formatDateTime(row.createdAt) }}</div>
        </template>
      </el-table-column>
      <el-table-column label="订单项目" min-width="260">
        <template #default="{ row }">
          <div v-for="item in row.items" :key="item.id" style="margin-bottom: 4px">
            {{ item.title }} × {{ item.quantity }} / {{ getLabel(cycleMap, item.cycle) }}
          </div>
        </template>
      </el-table-column>
      <el-table-column label="订单金额" width="130">
        <template #default="{ row }">
          {{ formatCurrency(row.totalAmount) }}
        </template>
      </el-table-column>
      <el-table-column label="已收金额" width="130">
        <template #default="{ row }">
          {{ formatCurrency(row.paidAmount) }}
        </template>
      </el-table-column>
      <el-table-column label="到期日" width="120">
        <template #default="{ row }">
          {{ formatDate(row.dueDate) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.status)">
            {{ getLabel(orderStatusMap, row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="110" fixed="right">
        <template #default="{ row }">
          <el-button text type="primary" @click="router.push(`/orders/${row.id}`)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="新增订单" width="620px">
      <el-form label-position="top">
        <el-form-item label="客户">
          <el-select v-model="form.customerId" style="width: 100%">
            <el-option
              v-for="item in customers"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="产品">
          <el-select v-model="form.productId" style="width: 100%">
            <el-option
              v-for="item in products"
              :key="item.id"
              :label="`${item.name} / ${formatCurrency(item.price)}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="服务名称">
          <el-input v-model="form.serviceName" placeholder="例如：香港节点业务 A-01" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="计费周期">
              <el-select v-model="form.cycle" style="width: 100%">
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
            <el-form-item label="购买数量">
              <el-input-number v-model="form.quantity" :min="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注">
          <el-input
            v-model="form.notes"
            type="textarea"
            :rows="4"
            placeholder="可填写来源、合同号、上线要求等信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确认创建</el-button>
      </template>
    </el-dialog>
  </ListWorkbenchShell>
</template>
