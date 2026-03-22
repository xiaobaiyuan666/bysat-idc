<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import { loadDashboard, type PortalWallet } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import { formatPortalMoney } from "@/utils/business";

const localeStore = useLocaleStore();
const loading = ref(true);
const error = ref("");
const wallet = ref<PortalWallet | null>(null);

const copy = computed(() => ({
  title: pickLabel(localeStore.locale, "钱包中心", "Wallet"),
  subtitle: pickLabel(
    localeStore.locale,
    "查看账户余额与信用额度，后续可直接用于账单支付和余额流水。",
    "Review account balance and credit limit for invoice payments and wallet usage."
  ),
  balance: pickLabel(localeStore.locale, "账户余额", "Balance"),
  creditLimit: pickLabel(localeStore.locale, "信用额度", "Credit Limit"),
  helpTitle: pickLabel(localeStore.locale, "钱包说明", "Wallet Guide"),
  helpDesc: pickLabel(
    localeStore.locale,
    "余额可直接抵扣未支付账单，信用额度用于授信支付。",
    "Balance can be used for unpaid invoices and credit limit for approved credit payments."
  ),
  balanceUsage: pickLabel(localeStore.locale, "余额用途", "Balance Usage"),
  balanceUsageValue: pickLabel(localeStore.locale, "支付账单与自动续费", "Invoice payment and auto-renew"),
  creditUsage: pickLabel(localeStore.locale, "信用额度", "Credit Limit"),
  creditUsageValue: pickLabel(localeStore.locale, "授信支付与风险控制", "Credit payment and risk control"),
  loadError: pickLabel(localeStore.locale, "钱包数据加载失败", "Failed to load wallet")
}));

async function fetchWallet() {
  loading.value = true;
  error.value = "";
  try {
    const dashboard = await loadDashboard();
    wallet.value = dashboard.wallet;
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

      <div class="portal-grid portal-grid--two" style="margin-top: 20px">
        <div class="portal-stat">
          <h3>{{ copy.balance }}</h3>
          <strong>{{ wallet ? formatPortalMoney(localeStore.locale, wallet.balance) : formatPortalMoney(localeStore.locale, 0) }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.creditLimit }}</h3>
          <strong>{{ wallet ? formatPortalMoney(localeStore.locale, wallet.creditLimit) : formatPortalMoney(localeStore.locale, 0) }}</strong>
        </div>
      </div>
    </section>

    <section class="portal-card">
      <div class="portal-card-head">
        <div>
          <h2 class="portal-panel__title">{{ copy.helpTitle }}</h2>
          <div class="portal-panel__meta">{{ copy.helpDesc }}</div>
        </div>
      </div>

      <div class="portal-summary" style="margin-top: 18px">
        <div class="portal-summary-row">
          <span>{{ copy.balanceUsage }}</span>
          <strong>{{ copy.balanceUsageValue }}</strong>
        </div>
        <div class="portal-summary-row">
          <span>{{ copy.creditUsage }}</span>
          <strong>{{ copy.creditUsageValue }}</strong>
        </div>
      </div>
    </section>
  </div>
</template>
