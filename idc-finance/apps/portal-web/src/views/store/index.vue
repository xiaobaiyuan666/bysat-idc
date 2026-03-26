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
const keyword = ref("");
const groupFilter = ref("");
const typeFilter = ref("");
const activeProductId = ref<number | null>(null);

const selectedCycles = reactive<Record<number, string>>({});
const selectedOptions = reactive<Record<number, Record<string, string>>>({});

const groupOptions = computed(() => Array.from(new Set(products.value.map(item => item.groupName).filter(Boolean))));
const typeOptions = computed(() => Array.from(new Set(products.value.map(item => item.productType).filter(Boolean))));

const filteredProducts = computed(() =>
  products.value.filter(product => {
    const text = keyword.value.trim().toLowerCase();
    const matchesKeyword =
      !text ||
      product.name.toLowerCase().includes(text) ||
      product.productNo.toLowerCase().includes(text) ||
      product.groupName.toLowerCase().includes(text);
    const matchesGroup = !groupFilter.value || product.groupName === groupFilter.value;
    const matchesType = !typeFilter.value || product.productType === typeFilter.value;
    return matchesKeyword && matchesGroup && matchesType;
  })
);

const activeProduct = computed(
  () => filteredProducts.value.find(item => item.id === activeProductId.value) ?? filteredProducts.value[0] ?? null
);

const visibleTotal = computed(() => filteredProducts.value.reduce((sum, item) => sum + productTotal(item), 0));
const productCount = computed(() => filteredProducts.value.length);
const groupCount = computed(() => groupOptions.value.length);
const typeCount = computed(() => typeOptions.value.length);

const copy = computed(() => ({
  title: pickLabel(localeStore.locale, "云产品商城", "Marketplace"),
  subtitle: pickLabel(
    localeStore.locale,
    "按配置项自定义下单，价格会实时汇总并进入账单流程。",
    "Customize products with configuration options and check out with live pricing."
  ),
  guidance: pickLabel(
    localeStore.locale,
    "商城只展示当前客户可订购的商品。选好周期和配置后，系统会先生成订单与账单，再进入支付流程。",
    "Only products currently available to this client are listed here. After choosing a cycle and configuration, the system creates an order and invoice before payment."
  ),
  availableProducts: pickLabel(localeStore.locale, "在售商品", "Available Products"),
  currentTotal: pickLabel(localeStore.locale, "当前选中总价", "Current Total"),
  groupCount: pickLabel(localeStore.locale, "商品分组", "Groups"),
  typeCount: pickLabel(localeStore.locale, "商品类型", "Types"),
  keywordPlaceholder: pickLabel(localeStore.locale, "搜索商品编号、商品名称或分组", "Search product no., name, or group"),
  groupPlaceholder: pickLabel(localeStore.locale, "商品分组", "Group"),
  typePlaceholder: pickLabel(localeStore.locale, "商品类型", "Type"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  reset: pickLabel(localeStore.locale, "重置筛选", "Reset"),
  shelfTitle: pickLabel(localeStore.locale, "商品货架", "Product Shelf"),
  shelfDesc: pickLabel(localeStore.locale, "先选择商品，再在右侧完成周期和配置项。", "Choose a product first, then finish cycle and options on the right."),
  detailTitle: pickLabel(localeStore.locale, "下单面板", "Checkout Panel"),
  detailDesc: pickLabel(localeStore.locale, "支持按周期、配置项和价格差异实时汇总。", "Live pricing updates for cycle and configuration changes."),
  purchaseCycle: pickLabel(localeStore.locale, "购买周期", "Billing Cycle"),
  configOptions: pickLabel(localeStore.locale, "配置项", "Configuration"),
  required: pickLabel(localeStore.locale, "必填", "Required"),
  optionHint: pickLabel(localeStore.locale, "按所选规格自动调整订单金额。", "Price updates automatically based on selected configuration."),
  basePrice: pickLabel(localeStore.locale, "基础售价", "Base Price"),
  setupFee: pickLabel(localeStore.locale, "开通费用", "Setup Fee"),
  optionDelta: pickLabel(localeStore.locale, "配置增量", "Configuration Delta"),
  payable: pickLabel(localeStore.locale, "本次应付", "Amount Due"),
  orderNow: pickLabel(localeStore.locale, "立即下单", "Order Now"),
  openInvoices: pickLabel(localeStore.locale, "查看账单", "Invoices"),
  openOrders: pickLabel(localeStore.locale, "查看订单", "Orders"),
  loadError: pickLabel(localeStore.locale, "商品列表加载失败", "Failed to load products"),
  cycleRequired: pickLabel(localeStore.locale, "请选择购买周期", "Please choose a billing cycle"),
  inputPlaceholder: pickLabel(localeStore.locale, "请输入配置值", "Please enter a value"),
  empty: pickLabel(localeStore.locale, "暂无匹配商品。", "No matching products."),
  noProduct: pickLabel(localeStore.locale, "当前没有可下单商品。", "No products available."),
  productNo: pickLabel(localeStore.locale, "商品编号", "Product No."),
  productType: pickLabel(localeStore.locale, "商品类型", "Product Type")
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
    if (!activeProductId.value && items.length > 0) {
      activeProductId.value = items[0].id;
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

function resetFilters() {
  keyword.value = "";
  groupFilter.value = "";
  typeFilter.value = "";
}

function focusProduct(product: PortalProduct) {
  activeProductId.value = product.id;
  ensureSelections(product);
}

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
        `下单成功，订单 ${result.order.orderNo} 与账单 ${result.invoice.invoiceNo} 已生成。`,
        `Order ${result.order.orderNo} and invoice ${result.invoice.invoiceNo} created successfully.`
      )
    );
    router.push({ path: "/invoices", query: { invoiceNo: result.invoice.invoiceNo } });
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

      <el-alert :title="copy.guidance" type="info" :closable="false" show-icon style="margin-top: 20px" />

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <div class="portal-stat">
          <h3>{{ copy.availableProducts }}</h3>
          <strong>{{ productCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.currentTotal }}</h3>
          <strong>{{ formatPortalMoney(localeStore.locale, visibleTotal) }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.groupCount }}</h3>
          <strong>{{ groupCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.typeCount }}</h3>
          <strong>{{ typeCount }}</strong>
        </div>
      </div>
    </section>

    <section class="portal-card">
      <div class="portal-toolbar">
        <el-input v-model="keyword" :placeholder="copy.keywordPlaceholder" clearable />
        <el-select v-model="groupFilter" :placeholder="copy.groupPlaceholder" clearable>
          <el-option v-for="item in groupOptions" :key="item" :label="item" :value="item" />
        </el-select>
        <el-select v-model="typeFilter" :placeholder="copy.typePlaceholder" clearable>
          <el-option
            v-for="item in typeOptions"
            :key="item"
            :label="formatPortalProductType(localeStore.locale, item)"
            :value="item"
          />
        </el-select>
        <el-button @click="resetFilters">{{ copy.reset }}</el-button>
        <div class="portal-summary" style="margin: 0; padding: 12px 16px">
          <div class="portal-summary-row" style="font-size: 12px">
            <span>{{ copy.visible }}</span>
            <strong>{{ filteredProducts.length }}</strong>
          </div>
        </div>
      </div>
    </section>

    <section class="store-grid">
      <article class="portal-card portal-table-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.shelfTitle }}</h2>
            <div class="portal-panel__meta">{{ copy.shelfDesc }}</div>
          </div>
        </div>

        <div v-if="filteredProducts.length" class="store-shelf">
          <button
            v-for="product in filteredProducts"
            :key="product.id"
            type="button"
            class="store-card"
            :class="{ 'store-card--active': activeProduct?.id === product.id }"
            @click="focusProduct(product)"
          >
            <div class="store-card__head">
              <div>
                <div class="portal-badge">{{ product.productNo }}</div>
                <h3 class="store-card__title">{{ product.name }}</h3>
              </div>
              <strong class="store-card__price">
                {{ formatPortalMoney(localeStore.locale, productTotal(product)) }}
              </strong>
            </div>
            <div class="store-card__meta">
              <span>{{ product.groupName }}</span>
              <span>{{ formatPortalProductType(localeStore.locale, product.productType) }}</span>
            </div>
            <p class="store-card__desc">{{ product.description }}</p>
          </button>
        </div>
        <div v-else class="portal-empty-state">{{ copy.empty }}</div>
      </article>

      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.detailTitle }}</h2>
            <div class="portal-panel__meta">{{ copy.detailDesc }}</div>
          </div>
        </div>

        <template v-if="activeProduct">
          <div class="portal-summary" style="margin-top: 12px">
            <div class="portal-summary-row">
              <span>{{ copy.productNo }}</span>
              <strong>{{ activeProduct.productNo }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.productType }}</span>
              <strong>{{ formatPortalProductType(localeStore.locale, activeProduct.productType) }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.groupPlaceholder }}</span>
              <strong>{{ activeProduct.groupName }}</strong>
            </div>
          </div>

          <p class="portal-subtitle" style="margin-top: 16px">{{ activeProduct.description }}</p>

          <div class="portal-section">
            <div class="portal-section-title">{{ copy.purchaseCycle }}</div>
            <el-radio-group v-model="selectedCycles[activeProduct.id]">
              <el-radio-button
                v-for="pricing in activeProduct.pricing"
                :key="pricing.cycleCode"
                :value="pricing.cycleCode"
              >
                {{ formatPortalBillingCycle(localeStore.locale, pricing.cycleCode) }} /
                {{ formatPortalMoney(localeStore.locale, pricing.price) }}
              </el-radio-button>
            </el-radio-group>
          </div>

          <div v-if="activeProduct.configOptions.length > 0" class="portal-section">
            <div class="portal-section-title">{{ copy.configOptions }}</div>
            <div style="display: grid; gap: 14px">
              <div v-for="option in activeProduct.configOptions" :key="option.code" class="portal-config-item">
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
                  <el-tag v-if="optionExtra(activeProduct, option) !== 0" type="warning" effect="plain">
                    +{{ formatPortalMoney(localeStore.locale, optionExtra(activeProduct, option)) }}
                  </el-tag>
                </div>

                <el-select
                  v-if="option.inputType === 'select' || option.inputType === 'radio'"
                  v-model="selectedOptions[activeProduct.id][option.code]"
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
                  v-model="selectedOptions[activeProduct.id][option.code]"
                  style="margin-top: 10px"
                  :placeholder="copy.inputPlaceholder"
                />
              </div>
            </div>
          </div>

          <div class="portal-summary">
            <div class="portal-summary-row">
              <span>{{ copy.basePrice }}</span>
              <strong>{{ formatPortalMoney(localeStore.locale, selectedCyclePrice(activeProduct)?.price) }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.setupFee }}</span>
              <strong>{{ formatPortalMoney(localeStore.locale, selectedCyclePrice(activeProduct)?.setupFee) }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.optionDelta }}</span>
              <strong>{{ formatPortalMoney(localeStore.locale, productOptionTotal(activeProduct)) }}</strong>
            </div>
            <div class="portal-summary-row portal-summary-row--total">
              <span>{{ copy.payable }}</span>
              <strong>{{ formatPortalMoney(localeStore.locale, productTotal(activeProduct)) }}</strong>
            </div>
          </div>

          <div class="portal-toolbar" style="margin-top: 20px">
            <el-button type="primary" :loading="submittingId === activeProduct.id" @click="handleCheckout(activeProduct)">
              {{ copy.orderNow }}
            </el-button>
            <el-button plain @click="router.push('/invoices')">{{ copy.openInvoices }}</el-button>
            <el-button plain @click="router.push('/orders')">{{ copy.openOrders }}</el-button>
          </div>
        </template>

        <div v-else class="portal-empty-state">
          {{ copy.noProduct }}
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.store-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(0, 0.9fr);
  gap: 20px;
}

.store-shelf {
  display: grid;
  gap: 14px;
}

.store-card {
  width: 100%;
  border: 1px solid rgba(148, 163, 184, 0.24);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.88);
  padding: 18px;
  text-align: left;
  cursor: pointer;
  transition: border-color 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease;
}

.store-card:hover,
.store-card--active {
  border-color: rgba(29, 78, 216, 0.45);
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
  transform: translateY(-1px);
}

.store-card__head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
}

.store-card__title {
  margin: 12px 0 0;
  font-size: 18px;
  color: #0f172a;
}

.store-card__price {
  color: #1d4ed8;
  font-size: 22px;
  white-space: nowrap;
}

.store-card__meta {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  color: #64748b;
  font-size: 13px;
  margin-top: 12px;
}

.store-card__desc {
  margin: 14px 0 0;
  color: #334155;
  line-height: 1.75;
}

@media (max-width: 1200px) {
  .store-grid {
    grid-template-columns: 1fr;
  }
}
</style>
