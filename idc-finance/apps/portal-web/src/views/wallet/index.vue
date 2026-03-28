<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { useRoute, useRouter } from "vue-router";
import {
  createRechargeOrder,
  loadWallet,
  loadWalletRecharges,
  loadWalletTransactions,
  type PortalWalletOverview,
  type PortalWalletTransaction
} from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import { formatPortalMoney, formatPortalPaymentChannel } from "@/utils/business";

const localeStore = useLocaleStore();
const route = useRoute();
const router = useRouter();
const loading = ref(false);
const tabLoading = ref(false);
const error = ref("");
const activeTab = ref("transactions");
const detail = ref<PortalWalletOverview | null>(null);
const transactions = ref<PortalWalletTransaction[]>([]);
const recharges = ref<PortalWalletTransaction[]>([]);
const rechargeDialogVisible = ref(false);
const rechargeSubmitting = ref(false);
const rechargeForm = reactive({
  amount: 100
});
const rechargePresets = [50, 100, 200, 500, 1000];

const rechargeCount = computed(() => recharges.value.length);
const consumeCount = computed(
  () => transactions.value.filter(item => item.transactionType === "CONSUME").length
);
const inflowAmount = computed(() =>
  transactions.value.filter(item => item.direction === "IN").reduce((sum, item) => sum + item.amount, 0)
);
const outflowAmount = computed(() =>
  transactions.value.filter(item => item.direction === "OUT").reduce((sum, item) => sum + item.amount, 0)
);
const latestTransactions = computed(() =>
  transactions.value
    .slice()
    .sort((a, b) => String(b.occurredAt).localeCompare(String(a.occurredAt)))
    .slice(0, 5)
);
const fundingActions = computed(() => [
  {
    title: copy.value.goInvoices,
    description: pickLabel(localeStore.locale, "优先处理未支付账单和资金压力。", "Handle unpaid invoices and funding pressure."),
    action: () => router.push("/invoices")
  },
  {
    title: copy.value.goRecharges,
    description: pickLabel(localeStore.locale, "查看最近充值和线下入账记录。", "Review recent recharges and offline funding records."),
    action: () => router.push("/wallet/recharges")
  },
  {
    title: copy.value.goOrders,
    description: pickLabel(localeStore.locale, "回看与资金相关的订单链路。", "Review order flows related to funding."),
    action: () => router.push("/orders")
  },
  {
    title: copy.value.goAccount,
    description: pickLabel(localeStore.locale, "维护账户信息与实名资料。", "Maintain account profile and identity information."),
    action: () => router.push("/account/profile")
  }
]);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "钱包中心", "Wallet"),
  title: pickLabel(localeStore.locale, "钱包中心", "Wallet"),
  subtitle: pickLabel(
    localeStore.locale,
    "集中查看余额、授信、资金流水和充值记录，并与账单中心形成联动闭环。",
    "Review balance, credit, wallet transactions, and recharge history in one place."
  ),
  balance: pickLabel(localeStore.locale, "账户余额", "Balance"),
  creditLimit: pickLabel(localeStore.locale, "授信额度", "Credit Limit"),
  creditUsed: pickLabel(localeStore.locale, "已用授信", "Credit Used"),
  availableCredit: pickLabel(localeStore.locale, "可用授信", "Available Credit"),
  rechargeCount: pickLabel(localeStore.locale, "充值记录", "Recharges"),
  consumeCount: pickLabel(localeStore.locale, "余额支出", "Balance Spend"),
  inflow: pickLabel(localeStore.locale, "累计流入", "Inflow"),
  outflow: pickLabel(localeStore.locale, "累计流出", "Outflow"),
  fundingHealth: pickLabel(localeStore.locale, "资金健康", "Funding Health"),
  fundingHealthDesc: pickLabel(localeStore.locale, "从余额、授信和资金流向判断当前资金状态。", "Understand current funding state through balance, credit, and cash flow."),
  latest: pickLabel(localeStore.locale, "最近流水", "Recent Transactions"),
  latestDesc: pickLabel(localeStore.locale, "优先关注最近的入账和支出动作。", "Review the latest inflow and outflow events first."),
  fundingDesk: pickLabel(localeStore.locale, "资金工作台", "Funding Desk"),
  fundingDeskDesc: pickLabel(localeStore.locale, "把账单处理、订单查看和账户维护收拢到同一块。", "Keep invoice handling, orders, and account actions in one place."),
  transactions: pickLabel(localeStore.locale, "资金流水", "Transactions"),
  recharges: pickLabel(localeStore.locale, "充值记录", "Recharge History"),
  transactionNo: pickLabel(localeStore.locale, "流水号", "Transaction No."),
  type: pickLabel(localeStore.locale, "类型", "Type"),
  direction: pickLabel(localeStore.locale, "方向", "Direction"),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  channel: pickLabel(localeStore.locale, "渠道", "Channel"),
  summary: pickLabel(localeStore.locale, "摘要", "Summary"),
  occurredAt: pickLabel(localeStore.locale, "发生时间", "Occurred At"),
  empty: pickLabel(localeStore.locale, "暂无资金流水", "No wallet transactions"),
  loadError: pickLabel(localeStore.locale, "钱包数据加载失败", "Failed to load wallet"),
  usageHint: pickLabel(
    localeStore.locale,
    "余额可直接用于支付账单；当余额不足时，系统会回落到在线支付或其他支付渠道。",
    "Balance can pay invoices directly. When balance is insufficient, the system falls back to online payment."
  ),
  goInvoices: pickLabel(localeStore.locale, "去处理账单", "Handle Invoices"),
  goOrders: pickLabel(localeStore.locale, "查看订单", "Orders"),
  goAccount: pickLabel(localeStore.locale, "账户资料", "Account"),
  goRecharges: pickLabel(localeStore.locale, "充值记录", "Recharge History")
}));

function transactionTypeLabel(value: string) {
  const mapping: Record<string, [string, string]> = {
    RECHARGE: ["充值", "Recharge"],
    CONSUME: ["余额支付", "Balance Pay"],
    REFUND: ["退款回退", "Refund"],
    ADJUSTMENT: ["人工调整", "Adjustment"],
    CREDIT_LIMIT: ["授信调整", "Credit Limit"]
  };
  const pair = mapping[value];
  return pair ? pickLabel(localeStore.locale, pair[0], pair[1]) : value;
}

function directionLabel(value: string) {
  const mapping: Record<string, [string, string]> = {
    IN: ["收入", "In"],
    OUT: ["支出", "Out"],
    FLAT: ["平移", "Flat"]
  };
  const pair = mapping[value];
  return pair ? pickLabel(localeStore.locale, pair[0], pair[1]) : value;
}

function directionTag(value: string) {
  const mapping: Record<string, "success" | "danger" | "info"> = {
    IN: "success",
    OUT: "danger",
    FLAT: "info"
  };
  return mapping[value] ?? "info";
}

function syncTabFromRoute() {
  activeTab.value = route.path === "/wallet/recharges" ? "recharges" : "transactions";
}

async function fetchWallet() {
  loading.value = true;
  error.value = "";
  try {
    detail.value = await loadWallet();
    transactions.value = (await loadWalletTransactions()).items;
    recharges.value = (await loadWalletRecharges()).items;
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

async function refreshActiveTab() {
  tabLoading.value = true;
  try {
    if (activeTab.value === "transactions") {
      transactions.value = (await loadWalletTransactions()).items;
    } else {
      recharges.value = (await loadWalletRecharges()).items;
    }
  } finally {
    tabLoading.value = false;
  }
}

function openRechargeDialog() {
  rechargeForm.amount = 100;
  rechargeDialogVisible.value = true;
}

function selectPresetAmount(amount: number) {
  rechargeForm.amount = amount;
}

async function handleRecharge() {
  if (!rechargeForm.amount || rechargeForm.amount <= 0) {
    ElMessage.warning("请输入有效的充值金额");
    return;
  }

  rechargeSubmitting.value = true;
  try {
    const order = await createRechargeOrder({
      amount: rechargeForm.amount,
      subject: "钱包在线充值"
    });
    ElMessage.success("充值订单已创建，正在打开支付页面");
    rechargeDialogVisible.value = false;
    if (order.payUrl) {
      window.open(order.payUrl, "_blank");
      window.setTimeout(() => {
        void fetchWallet();
      }, 1500);
    }
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : "充值订单创建失败");
  } finally {
    rechargeSubmitting.value = false;
  }
}

onMounted(() => {
  syncTabFromRoute();
  void fetchWallet();
});

watch(
  () => route.path,
  () => {
    syncTabFromRoute();
  }
);

watch(activeTab, tab => {
  const target = tab === "recharges" ? "/wallet/recharges" : "/wallet/transactions";
  if (route.path !== target && route.path !== "/wallet") {
    void router.replace(target);
  } else if (route.path === "/wallet") {
    void router.replace(target);
  }
  void refreshActiveTab();
});
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />

    <section class="portal-card">
      <div class="portal-card-head">
        <div>
          <div class="portal-badge">{{ copy.badge }}</div>
          <h1 class="portal-title" style="margin-top: 12px">{{ copy.title }}</h1>
          <p class="portal-subtitle">{{ copy.subtitle }}</p>
        </div>
      </div>

      <el-alert :title="copy.usageHint" type="info" :closable="false" show-icon style="margin-top: 20px" />

      <div class="portal-toolbar" style="margin-top: 20px">
        <el-button type="success" @click="openRechargeDialog">在线充值</el-button>
        <el-button type="primary" plain @click="router.push('/invoices')">{{ copy.goInvoices }}</el-button>
        <el-button plain @click="router.push('/orders')">{{ copy.goOrders }}</el-button>
        <el-button plain @click="router.push('/account/profile')">{{ copy.goAccount }}</el-button>
      </div>

      <div class="portal-grid portal-grid--three" style="margin-top: 20px">
        <article class="portal-stat"><h3>{{ copy.balance }}</h3><strong>{{ formatPortalMoney(localeStore.locale, detail?.wallet.balance ?? 0) }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.creditLimit }}</h3><strong>{{ formatPortalMoney(localeStore.locale, detail?.wallet.creditLimit ?? 0) }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.creditUsed }}</h3><strong>{{ formatPortalMoney(localeStore.locale, detail?.wallet.creditUsed ?? 0) }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.availableCredit }}</h3><strong>{{ formatPortalMoney(localeStore.locale, detail?.wallet.availableCredit ?? 0) }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.rechargeCount }}</h3><strong>{{ rechargeCount }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.consumeCount }}</h3><strong>{{ consumeCount }}</strong></article>
      </div>
    </section>

    <section class="portal-grid portal-grid--two">
      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.fundingHealth }}</h2>
            <div class="portal-panel__meta">{{ copy.fundingHealthDesc }}</div>
          </div>
        </div>
        <div class="portal-hero-grid" style="margin-top: 18px">
          <div class="portal-mini-card">
            <span class="portal-mini-card__label">{{ copy.inflow }}</span>
            <strong class="portal-mini-card__value">{{ formatPortalMoney(localeStore.locale, inflowAmount) }}</strong>
          </div>
          <div class="portal-mini-card">
            <span class="portal-mini-card__label">{{ copy.outflow }}</span>
            <strong class="portal-mini-card__value">{{ formatPortalMoney(localeStore.locale, outflowAmount) }}</strong>
          </div>
          <div class="portal-mini-card">
            <span class="portal-mini-card__label">{{ copy.rechargeCount }}</span>
            <strong class="portal-mini-card__value">{{ rechargeCount }}</strong>
          </div>
          <div class="portal-mini-card">
            <span class="portal-mini-card__label">{{ copy.consumeCount }}</span>
            <strong class="portal-mini-card__value">{{ consumeCount }}</strong>
          </div>
        </div>
      </article>

      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.fundingDesk }}</h2>
            <div class="portal-panel__meta">{{ copy.fundingDeskDesc }}</div>
          </div>
        </div>
        <div class="portal-actions-grid" style="margin-top: 18px">
          <button
            v-for="item in fundingActions"
            :key="item.title"
            type="button"
            class="portal-action-card"
            @click="item.action()"
          >
            <strong>{{ item.title }}</strong>
            <span>{{ item.description }}</span>
          </button>
        </div>
      </article>
    </section>

    <section class="portal-grid portal-grid--two">
      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.latest }}</h2>
            <div class="portal-panel__meta">{{ copy.latestDesc }}</div>
          </div>
        </div>
        <div v-if="latestTransactions.length" class="portal-list" style="margin-top: 18px">
          <div v-for="item in latestTransactions" :key="item.transactionNo" class="portal-list-item">
            <div class="portal-list-item__meta">
              <div class="portal-list-item__title">{{ item.summary || item.transactionNo }}</div>
              <div class="portal-list-item__desc">{{ item.transactionNo }}</div>
            </div>
            <div class="portal-toolbar">
              <el-tag :type="directionTag(item.direction)">{{ directionLabel(item.direction) }}</el-tag>
              <strong>{{ formatPortalMoney(localeStore.locale, item.amount) }}</strong>
            </div>
          </div>
        </div>
        <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.empty }}</div>
      </article>

      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.fundingDesk }}</h2>
            <div class="portal-panel__meta">{{ copy.fundingDeskDesc }}</div>
          </div>
        </div>
        <div class="portal-summary" style="margin-top: 18px">
          <div class="portal-summary-row">
            <span>{{ copy.inflow }}</span>
            <strong>{{ formatPortalMoney(localeStore.locale, inflowAmount) }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.outflow }}</span>
            <strong>{{ formatPortalMoney(localeStore.locale, outflowAmount) }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.rechargeCount }}</span>
            <strong>{{ rechargeCount }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.consumeCount }}</span>
            <strong>{{ consumeCount }}</strong>
          </div>
        </div>
        <div class="portal-actions-grid" style="margin-top: 18px">
          <button type="button" class="portal-action-card" @click="router.push('/invoices')">
            <strong>{{ copy.goInvoices }}</strong>
            <span>{{ copy.transactions }}</span>
          </button>
          <button type="button" class="portal-action-card" @click="router.push('/wallet/recharges')">
            <strong>{{ copy.goRecharges }}</strong>
            <span>{{ copy.recharges }}</span>
          </button>
          <button type="button" class="portal-action-card" @click="router.push('/orders')">
            <strong>{{ copy.goOrders }}</strong>
            <span>{{ copy.amount }}</span>
          </button>
          <button type="button" class="portal-action-card" @click="router.push('/account/profile')">
            <strong>{{ copy.goAccount }}</strong>
            <span>{{ copy.balance }}</span>
          </button>
        </div>
      </article>
    </section>

    <section class="portal-card" v-loading="tabLoading">
      <el-tabs v-model="activeTab">
        <el-tab-pane :label="copy.transactions" name="transactions">
          <el-table v-if="transactions.length" :data="transactions" border style="margin-top: 12px">
            <el-table-column prop="transactionNo" :label="copy.transactionNo" min-width="150" />
            <el-table-column :label="copy.type" min-width="110">
              <template #default="{ row }">{{ transactionTypeLabel(row.transactionType) }}</template>
            </el-table-column>
            <el-table-column :label="copy.direction" min-width="90">
              <template #default="{ row }">
                <el-tag :type="directionTag(row.direction)">{{ directionLabel(row.direction) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="copy.amount" min-width="120">
              <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
            </el-table-column>
            <el-table-column :label="copy.channel" min-width="140">
              <template #default="{ row }">{{ formatPortalPaymentChannel(localeStore.locale, row.channel) }}</template>
            </el-table-column>
            <el-table-column prop="summary" :label="copy.summary" min-width="220" show-overflow-tooltip />
            <el-table-column prop="occurredAt" :label="copy.occurredAt" min-width="160" />
          </el-table>
          <div v-else class="portal-empty-state" style="margin-top: 16px">{{ copy.empty }}</div>
        </el-tab-pane>

        <el-tab-pane :label="copy.recharges" name="recharges">
          <el-table v-if="recharges.length" :data="recharges" border style="margin-top: 12px">
            <el-table-column prop="transactionNo" :label="copy.transactionNo" min-width="150" />
            <el-table-column :label="copy.amount" min-width="120">
              <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
            </el-table-column>
            <el-table-column :label="copy.channel" min-width="140">
              <template #default="{ row }">{{ formatPortalPaymentChannel(localeStore.locale, row.channel) }}</template>
            </el-table-column>
            <el-table-column prop="summary" :label="copy.summary" min-width="220" show-overflow-tooltip />
            <el-table-column prop="occurredAt" :label="copy.occurredAt" min-width="160" />
          </el-table>
          <div v-else class="portal-empty-state" style="margin-top: 16px">{{ copy.empty }}</div>
        </el-tab-pane>
      </el-tabs>
    </section>

    <el-dialog v-model="rechargeDialogVisible" title="在线充值" width="480px" destroy-on-close>
      <el-form label-position="top">
        <el-form-item label="快捷金额">
          <div style="display: flex; gap: 8px; flex-wrap: wrap">
            <el-button
              v-for="preset in rechargePresets"
              :key="preset"
              :type="rechargeForm.amount === preset ? 'primary' : 'default'"
              @click="selectPresetAmount(preset)"
            >
              ¥{{ preset }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="自定义金额">
          <el-input-number
            v-model="rechargeForm.amount"
            :min="1"
            :max="100000"
            :precision="2"
            :step="10"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rechargeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="rechargeSubmitting" @click="handleRecharge">
          确认充值
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>
