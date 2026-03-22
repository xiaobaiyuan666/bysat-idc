<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import {
  checkoutProduct,
  loadStoreProducts,
  type CheckoutSelection,
  type PortalProduct,
  type PortalProductConfigOption
} from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalBillingCycle,
  formatPortalMoney,
  formatPortalProductType
} from "@/utils/business";

const router = useRouter();
const localeStore = useLocaleStore();
const loading = ref(true);
const error = ref("");
const submittingId = ref<number | null>(null);
const products = ref<PortalProduct[]>([]);
const selectedCycles = reactive<Record<number, string>>({});
const selectedOptions = reactive<Record<number, Record<string, string>>>({});

const copy = computed(() => ({
  title: pickLabel(localeStore.locale, "产品商城", "Marketplace"),
  subtitle: pickLabel(
    localeStore.locale,
    "按配置项自定义下单，价格会实时汇总并进入账单流程。",
    "Customize products with configuration options and check out with live pricing."
  ),
  availableProducts: pickLabel(localeStore.locale, "可售商品", "Available Products"),
  currentTotal: pickLabel(localeStore.locale, "当前总价", "Current Total"),
  purchaseCycle: pickLabel(localeStore.locale, "购买周期", "Billing Cycle"),
  configOptions: pickLabel(localeStore.locale, "配置项", "Configuration"),
  required: pickLabel(localeStore.locale, "必填", "Required"),
  optionHint: pickLabel(
    localeStore.locale,
    "按所选规格自动调整订单金额。",
    "Price updates automatically based on selected configuration."
  ),
  basePrice: pickLabel(localeStore.locale, "基础售价", "Base Price"),
  setupFee: pickLabel(localeStore.locale, "开通费用", "Setup Fee"),
  optionDelta: pickLabel(localeStore.locale, "配置增量", "Configuration Delta"),
  payable: pickLabel(localeStore.locale, "本次应付", "Amount Due"),
  orderNow: pickLabel(localeStore.locale, "立即下单", "Order Now"),
  loadError: pickLabel(localeStore.locale, "商品列表加载失败", "Failed to load products"),
  cycleRequired: pickLabel(localeStore.locale, "请选择购买周期", "Please choose a billing cycle"),
  inputPlaceholder: pickLabel(localeStore.locale, "请输入配置值", "Please enter a value")
}));

function ensureSelections(product: PortalProduct) {
  if (!selectedCycles[product.id] && product.pricing.length > 0) {
    selectedCycles[product.id] = product.pricing[0].cycleCode;
  }
  if (!selectedOptions[product.id]) {
    selectedOptions[product.id] = {};
  }
  for (const option of product.configOptions) {
    if (!selectedOptions[product.id][option.code]) {
      selectedOptions[product.id][option.code] = option.defaultValue || option.choices[0]?.value || "";
    }
  }
}

async function fetchProducts() {
  loading.value = true;
  error.value = "";
  try {
    const items = await loadStoreProducts();
    products.value = items;
    for (const product of items) {
      ensureSelections(product);
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

function selectedCyclePrice(product: PortalProduct) {
  const cycleCode = selectedCycles[product.id];
  return product.pricing.find(item => item.cycleCode === cycleCode) ?? product.pricing[0];
}

function selectedChoice(option: PortalProductConfigOption, value: string) {
  return option.choices.find(item => item.value === value);
}

function optionExtra(product: PortalProduct, option: PortalProductConfigOption) {
  const value = selectedOptions[product.id]?.[option.code] ?? option.defaultValue;
  return selectedChoice(option, value)?.priceDelta ?? 0;
}

function productConfigSelections(product: PortalProduct): CheckoutSelection[] {
  return product.configOptions
    .map(option => ({
      code: option.code,
      value: selectedOptions[product.id]?.[option.code] ?? option.defaultValue
    }))
    .filter(item => item.code && item.value);
}

function productOptionTotal(product: PortalProduct) {
  return product.configOptions.reduce((sum, option) => sum + optionExtra(product, option), 0);
}

function productTotal(product: PortalProduct) {
  const cycle = selectedCyclePrice(product);
  if (!cycle) return 0;
  return cycle.price + cycle.setupFee + productOptionTotal(product);
}

const summaryItems = computed(() =>
  products.value.map(product => ({
    id: product.id,
    total: productTotal(product)
  }))
);

async function handleCheckout(product: PortalProduct) {
  const cycleCode = selectedCycles[product.id];
  if (!cycleCode) {
    ElMessage.warning(copy.value.cycleRequired);
    return;
  }

  for (const option of product.configOptions) {
    const value = selectedOptions[product.id]?.[option.code];
    if (option.required && !value) {
      ElMessage.warning(
        pickLabel(
          localeStore.locale,
          `请完善配置项：${option.name}`,
          `Please complete the configuration: ${option.name}`
        )
      );
      return;
    }
  }

  submittingId.value = product.id;
  try {
    const result = await checkoutProduct({
      productId: product.id,
      cycleCode,
      selectedOptions: productConfigSelections(product)
    });
    ElMessage.success(
      pickLabel(
        localeStore.locale,
        `下单成功，订单 ${result.order.orderNo} 和账单 ${result.invoice.invoiceNo} 已生成`,
        `Order ${result.order.orderNo} and invoice ${result.invoice.invoiceNo} created successfully.`
      )
    );
    router.push("/invoices");
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : pickLabel(localeStore.locale, "下单失败", "Checkout failed"));
  } finally {
    submittingId.value = null;
  }
}

onMounted(fetchProducts);
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
          <h3>{{ copy.availableProducts }}</h3>
          <strong>{{ products.length }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.currentTotal }}</h3>
          <strong>{{ formatPortalMoney(localeStore.locale, summaryItems.reduce((sum, item) => sum + item.total, 0)) }}</strong>
        </div>
      </div>
    </section>

    <div class="portal-grid portal-grid--two">
      <article v-for="product in products" :key="product.id" class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ product.name }}</h2>
            <div class="portal-panel__meta">
              {{ formatPortalProductType(localeStore.locale, product.productType) }} / {{ product.groupName }}
            </div>
            <p class="portal-subtitle" style="margin-top: 10px">{{ product.description }}</p>
          </div>
          <strong style="font-size: 24px; color: #1d4ed8">
            {{ formatPortalMoney(localeStore.locale, summaryItems.find(item => item.id === product.id)?.total) }}
          </strong>
        </div>

        <div class="portal-section">
          <div class="portal-section-title">{{ copy.purchaseCycle }}</div>
          <el-radio-group v-model="selectedCycles[product.id]">
            <el-radio-button
              v-for="pricing in product.pricing"
              :key="pricing.cycleCode"
              :value="pricing.cycleCode"
            >
              {{ formatPortalBillingCycle(localeStore.locale, pricing.cycleCode) }} /
              {{ formatPortalMoney(localeStore.locale, pricing.price) }}
            </el-radio-button>
          </el-radio-group>
        </div>

        <div v-if="product.configOptions.length > 0" class="portal-section">
          <div class="portal-section-title">{{ copy.configOptions }}</div>
          <div style="display: grid; gap: 14px">
            <div v-for="option in product.configOptions" :key="option.code" class="portal-config-item">
              <div class="portal-card-head">
                <div>
                  <strong>{{ option.name }}</strong>
                  <span v-if="option.required" style="margin-left: 8px; color: #d97706; font-size: 12px">
                    {{ copy.required }}
                  </span>
                  <div class="portal-subtitle" style="margin-top: 4px">
                    {{ option.description || copy.optionHint }}
                  </div>
                </div>
                <el-tag v-if="optionExtra(product, option) !== 0" type="warning" effect="plain">
                  +{{ formatPortalMoney(localeStore.locale, optionExtra(product, option)) }}
                </el-tag>
              </div>

              <el-select
                v-if="option.inputType === 'select' || option.inputType === 'radio'"
                v-model="selectedOptions[product.id][option.code]"
                style="margin-top: 10px; width: 100%"
              >
                <el-option
                  v-for="choice in option.choices"
                  :key="choice.value"
                  :label="`${choice.label}${choice.priceDelta ? ` / +${formatPortalMoney(localeStore.locale, choice.priceDelta)}` : ''}`"
                  :value="choice.value"
                />
              </el-select>

              <el-input
                v-else
                v-model="selectedOptions[product.id][option.code]"
                style="margin-top: 10px"
                :placeholder="copy.inputPlaceholder"
              />
            </div>
          </div>
        </div>

        <div class="portal-summary">
          <div class="portal-summary-row">
            <span>{{ copy.basePrice }}</span>
            <strong>{{ formatPortalMoney(localeStore.locale, selectedCyclePrice(product)?.price) }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.setupFee }}</span>
            <strong>{{ formatPortalMoney(localeStore.locale, selectedCyclePrice(product)?.setupFee) }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.optionDelta }}</span>
            <strong>{{ formatPortalMoney(localeStore.locale, productOptionTotal(product)) }}</strong>
          </div>
          <div class="portal-summary-row portal-summary-row--total">
            <span>{{ copy.payable }}</span>
            <strong>{{ formatPortalMoney(localeStore.locale, productTotal(product)) }}</strong>
          </div>
        </div>

        <div style="margin-top: 20px; display: flex; justify-content: flex-end">
          <el-button type="primary" :loading="submittingId === product.id" @click="handleCheckout(product)">
            {{ copy.orderNow }}
          </el-button>
        </div>
      </article>
    </div>
  </div>
</template>
