<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { loadAccount, type PortalAccount } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalCustomerStatus,
  formatPortalIdentityStatus,
  formatPortalMoney
} from "@/utils/business";

const localeStore = useLocaleStore();
const loading = ref(true);
const error = ref("");
const account = ref<PortalAccount | null>(null);

const copy = computed(() => ({
  title: pickLabel(localeStore.locale, "账户中心", "Account Center"),
  subtitle: pickLabel(
    localeStore.locale,
    "查看客户资料、默认联系人、实名认证状态和钱包概览。",
    "Review profile, default contact, identity verification, and wallet summary."
  ),
  customerName: pickLabel(localeStore.locale, "客户名称", "Customer"),
  customerStatus: pickLabel(localeStore.locale, "账户状态", "Account Status"),
  balance: pickLabel(localeStore.locale, "账户余额", "Balance"),
  basicInfo: pickLabel(localeStore.locale, "基础资料", "Profile"),
  basicInfoDesc: pickLabel(
    localeStore.locale,
    "客户编号、联系方式、客户分组和销售归属。",
    "Customer number, contact info, grouping, and sales ownership."
  ),
  identityInfo: pickLabel(localeStore.locale, "实名信息", "Identity"),
  identityInfoDesc: pickLabel(
    localeStore.locale,
    "当前实名认证审核状态和证件主体信息。",
    "Current verification status and identity holder information."
  ),
  primaryContact: pickLabel(localeStore.locale, "默认联系人", "Primary Contact"),
  primaryContactDesc: pickLabel(
    localeStore.locale,
    "优先用于账单通知、工单和业务联系。",
    "Used first for invoice, ticket, and service communication."
  ),
  customerNo: pickLabel(localeStore.locale, "客户编号", "Customer No."),
  email: pickLabel(localeStore.locale, "邮箱", "Email"),
  mobile: pickLabel(localeStore.locale, "手机", "Mobile"),
  groupName: pickLabel(localeStore.locale, "客户分组", "Customer Group"),
  levelName: pickLabel(localeStore.locale, "客户等级", "Customer Level"),
  salesOwner: pickLabel(localeStore.locale, "销售归属", "Sales Owner"),
  subjectName: pickLabel(localeStore.locale, "主体名称", "Subject"),
  certNo: pickLabel(localeStore.locale, "证件号码", "Document No."),
  verifyStatus: pickLabel(localeStore.locale, "审核状态", "Verification Status"),
  reviewRemark: pickLabel(localeStore.locale, "审核备注", "Review Remark"),
  name: pickLabel(localeStore.locale, "姓名", "Name"),
  role: pickLabel(localeStore.locale, "角色", "Role"),
  loadError: pickLabel(localeStore.locale, "账户信息加载失败", "Failed to load account")
}));

async function fetchAccount() {
  loading.value = true;
  error.value = "";

  try {
    account.value = await loadAccount();
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

onMounted(fetchAccount);
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />

    <section v-if="account" class="portal-card">
      <div class="portal-card-head">
        <div>
          <h1 class="portal-title">{{ copy.title }}</h1>
          <p class="portal-subtitle">{{ copy.subtitle }}</p>
        </div>
      </div>

      <div class="portal-grid portal-grid--three" style="margin-top: 20px">
        <div class="portal-stat">
          <h3>{{ copy.customerName }}</h3>
          <strong>{{ account.customer.name }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.customerStatus }}</h3>
          <strong>{{ formatPortalCustomerStatus(localeStore.locale, account.customer.status) }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.balance }}</h3>
          <strong>{{ formatPortalMoney(localeStore.locale, account.wallet.balance) }}</strong>
        </div>
      </div>
    </section>

    <section v-if="account" class="portal-grid portal-grid--two">
      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.basicInfo }}</h2>
            <div class="portal-panel__meta">{{ copy.basicInfoDesc }}</div>
          </div>
        </div>

        <el-descriptions :column="1" border style="margin-top: 18px">
          <el-descriptions-item :label="copy.customerNo">
            {{ account.customer.customerNo }}
          </el-descriptions-item>
          <el-descriptions-item :label="copy.email">
            {{ account.customer.email }}
          </el-descriptions-item>
          <el-descriptions-item :label="copy.mobile">
            {{ account.customer.mobile }}
          </el-descriptions-item>
          <el-descriptions-item :label="copy.groupName">
            {{ account.customer.groupName || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="copy.levelName">
            {{ account.customer.levelName || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="copy.salesOwner">
            {{ account.customer.salesOwner || "-" }}
          </el-descriptions-item>
        </el-descriptions>
      </article>

      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.identityInfo }}</h2>
            <div class="portal-panel__meta">{{ copy.identityInfoDesc }}</div>
          </div>
        </div>

        <el-descriptions :column="1" border style="margin-top: 18px">
          <el-descriptions-item :label="copy.subjectName">
            {{ account.customer.identity?.subjectName ?? "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="copy.certNo">
            {{ account.customer.identity?.certNo ?? "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="copy.verifyStatus">
            <el-tag type="info" effect="plain">
              {{ formatPortalIdentityStatus(localeStore.locale, account.customer.identity?.verifyStatus ?? "") }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="copy.reviewRemark">
            {{ account.customer.identity?.reviewRemark ?? "-" }}
          </el-descriptions-item>
        </el-descriptions>
      </article>
    </section>

    <section v-if="account" class="portal-card" style="margin-top: 20px">
      <div class="portal-card-head">
        <div>
          <h2 class="portal-panel__title">{{ copy.primaryContact }}</h2>
          <div class="portal-panel__meta">{{ copy.primaryContactDesc }}</div>
        </div>
      </div>

      <el-descriptions :column="2" border style="margin-top: 18px">
        <el-descriptions-item :label="copy.name">
          {{ account.primaryContact?.name ?? "-" }}
        </el-descriptions-item>
        <el-descriptions-item :label="copy.role">
          {{ account.primaryContact?.roleName ?? "-" }}
        </el-descriptions-item>
        <el-descriptions-item :label="copy.email">
          {{ account.primaryContact?.email ?? "-" }}
        </el-descriptions-item>
        <el-descriptions-item :label="copy.mobile">
          {{ account.primaryContact?.mobile ?? "-" }}
        </el-descriptions-item>
      </el-descriptions>
    </section>
  </div>
</template>
