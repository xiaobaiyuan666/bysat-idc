<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { loadWallet, type PortalWalletOverview } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import { formatPortalMoney, formatPortalPaymentChannel } from "@/utils/business";

const localeStore = useLocaleStore();
const router = useRouter();
const loading = ref(true);
const error = ref("");
const detail = ref<PortalWalletOverview | null>(null);

const rechargeCount = computed(
  () => detail.value?.transactions.filter(item => item.transactionType === "RECHARGE").length ?? 0
);
const consumeCount = computed(
  () => detail.value?.transactions.filter(item => item.transactionType === "CONSUME").length ?? 0
);

const copy = computed(() => ({
  title: pickLabel(localeStore.locale, "钱包中心", "Wallet"),
  subtitle: pickLabel(
    localeStore.locale,
    "查看账户余额、授信额度以及每一笔资金流水。",
    "Review balance, credit limit, and each wallet transaction."
  ),
  balance: pickLabel(localeStore.locale, "账户余额", "Balance"),
  creditLimit: pickLabel(localeStore.locale, "授信额度", "Credit Limit"),
  creditUsed: pickLabel(localeStore.locale, "已用授信", "Credit Used"),
  availableCredit: pickLabel(localeStore.locale, "可用授信", "Available Credit"),
  recent: pickLabel(localeStore.locale, "资金流水", "Transactions"),
  transactionNo: pickLabel(localeStore.locale, "流水号", "Transaction No."),
  type: pickLabel(localeStore.locale, "类型", "Type"),
  direction: pickLabel(localeStore.locale, "方向", "Direction"),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  channel: pickLabel(localeStore.locale, "渠道", "Channel"),
  beforeAfter: pickLabel(localeStore.locale, "变更前后", "Before / After"),
  summary: pickLabel(localeStore.locale, "摘要", "Summary"),
  occurredAt: pickLabel(localeStore.locale, "时间", "Occurred At"),
  empty: pickLabel(localeStore.locale, "暂无钱包流水", "No wallet transactions"),
  loadError: pickLabel(localeStore.locale, "钱包数据加载失败", "Failed to load wallet"),
  usageHint: pickLabel(
    localeStore.locale,
    "账户余额可直接用于账单支付；支付时系统会优先尝试余额抵扣。",
    "Wallet balance can be used to pay invoices directly; the system tries balance first."
  ),
  rechargeCount: pickLabel(localeStore.locale, "充值记录", "Recharges"),
  consumeCount: pickLabel(localeStore.locale, "余额支付", "Balance Pay"),
  goInvoices: pickLabel(localeStore.locale, "去处理账单", "Handle Invoices"),
  goOrders: pickLabel(localeStore.locale, "查看订单", "Orders"),
  goAccount: pickLabel(localeStore.locale, "账户资料", "Account"),
  countUnit: pickLabel(localeStore.locale, "条", "records")
}));

const transactionTypeLabel = (value: string) =>
  pickLabel(
    localeStore.locale,
    (
      {
        RECHARGE: "充值",
        CONSUME: "余额支付",
        REFUND: "退款回退",
        ADJUSTMENT: "手工调整",
        CREDIT_LIMIT: "授信调整"
      }[value] ?? value
    ),
    (
      {
        RECHARGE: "Recharge",
        CONSUME: "Balance Pay",
        REFUND: "Refund",
        ADJUSTMENT: "Adjustment",
        CREDIT_LIMIT: "Credit Limit"
      }[value] ?? value
    )
  );

const directionLabel = (value: string) =>
  pickLabel(
    localeStore.locale,
    (
      {
        IN: "收入",
        OUT: "支出",
        FLAT: "平移"
      }[value] ?? value
    ),
    (
      {
        IN: "In",
        OUT: "Out",
        FLAT: "Flat"
      }[value] ?? value
    )
  );

const directionTag = (value: string) =>
  (
    {
      IN: "success",
      OUT: "danger",
      FLAT: "info"
    }[value] ?? "info"
  ) as "success" | "danger" | "info";

async function fetchWallet() {
  loading.value = true;
  error.value = "";
  try {
    detail.value = await loadWallet();
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

onMounted(fetchWallet);
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />

    <section class="portal-card">
      <div class="portal-card-head">
        <div>
          <h1 class="portal-title">{{ copy.title }}</h1>
          <p class="portal-subtitle">{{ copy.subtitle }}</p>
        </div>
      </div>

      <el-alert :title="copy.usageHint" type="info" :closable="false" show-icon style="margin-top: 20px" />

      <div class="portal-toolbar" style="margin-top: 18px">
        <el-button type="primary" plain @click="router.push('/invoices')">{{ copy.goInvoices }}</el-button>
        <el-button plain @click="router.push('/orders')">{{ copy.goOrders }}</el-button>
        <el-button plain @click="router.push('/account')">{{ copy.goAccount }}</el-button>
      </div>

      <div class="portal-grid portal-grid--two" style="margin-top: 20px">
        <div class="portal-stat">
          <h3>{{ copy.balance }}</h3>
          <strong>{{ formatPortalMoney(localeStore.locale, detail?.wallet.balance ?? 0) }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.creditLimit }}</h3>
          <strong>{{ formatPortalMoney(localeStore.locale, detail?.wallet.creditLimit ?? 0) }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.creditUsed }}</h3>
          <strong>{{ formatPortalMoney(localeStore.locale, detail?.wallet.creditUsed ?? 0) }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.availableCredit }}</h3>
          <strong>{{ formatPortalMoney(localeStore.locale, detail?.wallet.availableCredit ?? 0) }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.rechargeCount }}</h3>
          <strong>{{ rechargeCount }} {{ copy.countUnit }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.consumeCount }}</h3>
          <strong>{{ consumeCount }} {{ copy.countUnit }}</strong>
        </div>
      </div>
    </section>

    <section class="portal-card portal-table-card">
      <div class="portal-card-head">
        <div>
          <h2 class="portal-panel__title">{{ copy.recent }}</h2>
          <div class="portal-panel__meta">
            {{ detail?.transactions.length ?? 0 }} {{ copy.countUnit }}
          </div>
        </div>
      </div>

      <el-table v-if="detail?.transactions.length" :data="detail.transactions" border>
        <el-table-column prop="transactionNo" :label="copy.transactionNo" min-width="150" />
        <el-table-column :label="copy.type" min-width="110">
          <template #default="{ row }">{{ transactionTypeLabel(row.transactionType) }}</template>
        </el-table-column>
        <el-table-column :label="copy.direction" min-width="90">
          <template #default="{ row }">
            <el-tag :type="directionTag(row.direction)" effect="light">{{ directionLabel(row.direction) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="copy.amount" min-width="120">
          <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
        </el-table-column>
        <el-table-column :label="copy.channel" min-width="120">
          <template #default="{ row }">{{ formatPortalPaymentChannel(localeStore.locale, row.channel) }}</template>
        </el-table-column>
        <el-table-column :label="copy.beforeAfter" min-width="220">
          <template #default="{ row }">
            余额 {{ formatPortalMoney(localeStore.locale, row.balanceBefore) }} -> {{ formatPortalMoney(localeStore.locale, row.balanceAfter) }}
            <br />
            授信 {{ formatPortalMoney(localeStore.locale, row.creditBefore) }} -> {{ formatPortalMoney(localeStore.locale, row.creditAfter) }}
          </template>
        </el-table-column>
        <el-table-column prop="summary" :label="copy.summary" min-width="220" show-overflow-tooltip />
        <el-table-column prop="occurredAt" :label="copy.occurredAt" min-width="170" />
      </el-table>
      <div v-else class="portal-empty-state" style="margin: 20px">
        {{ copy.empty }}
      </div>
    </section>
  </div>
</template>
