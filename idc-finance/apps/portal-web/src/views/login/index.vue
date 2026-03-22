<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { localeOptions, pickLabel } from "@/locales";
import { useLocaleStore, useUserStore } from "@/store";

const router = useRouter();
const localeStore = useLocaleStore();
const userStore = useUserStore();
const loading = ref(false);

const form = reactive({
  username: "portal",
  password: "Portal123!"
});

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "IDC 客户中心", "IDC Client Area"),
  title: pickLabel(localeStore.locale, "面向客户的云业务自助入口", "Self-service entry for cloud customers"),
  subtitle: pickLabel(
    localeStore.locale,
    "统一查看服务、订单、账单、钱包和工单，支持续费、支付、状态查询和提交售后请求。",
    "View services, orders, invoices, wallet, and tickets in one place with self-service actions."
  ),
  statA: pickLabel(localeStore.locale, "客户自助", "Self Service"),
  statB: pickLabel(localeStore.locale, "统一", "Unified"),
  statC: pickLabel(localeStore.locale, "实时", "Live"),
  statADesc: pickLabel(localeStore.locale, "7x24 门户", "24/7 portal"),
  statBDesc: pickLabel(localeStore.locale, "业务入口", "business access"),
  statCDesc: pickLabel(localeStore.locale, "状态同步", "status sync"),
  loginTitle: pickLabel(localeStore.locale, "登录客户中心", "Client Login"),
  loginSubtitle: pickLabel(localeStore.locale, "演示账号：portal / Portal123!", "Demo account: portal / Portal123!"),
  username: pickLabel(localeStore.locale, "账号", "Username"),
  password: pickLabel(localeStore.locale, "密码", "Password"),
  submit: pickLabel(localeStore.locale, "进入客户中心", "Enter Client Area"),
  success: pickLabel(localeStore.locale, "登录成功", "Login successful"),
  failed: pickLabel(localeStore.locale, "账号或密码错误", "Invalid username or password")
}));

async function handleSubmit() {
  loading.value = true;
  try {
    await userStore.login(form);
    ElMessage.success(copy.value.success);
    router.push("/console");
  } catch {
    ElMessage.error(copy.value.failed);
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="portal-login-shell">
    <div class="portal-login-hero">
      <div class="portal-login-hero__toolbar">
        <div class="portal-login-hero__badge">{{ copy.badge }}</div>
        <el-button-group>
          <el-button
            v-for="item in localeOptions"
            :key="item.value"
            size="small"
            :type="localeStore.locale === item.value ? 'primary' : 'default'"
            @click="localeStore.setLocale(item.value)"
          >
            {{ item.label }}
          </el-button>
        </el-button-group>
      </div>
      <h1>{{ copy.title }}</h1>
      <p>{{ copy.subtitle }}</p>
      <div class="portal-login-hero__stats">
        <div>
          <strong>24/7</strong>
          <span>{{ copy.statADesc }}</span>
        </div>
        <div>
          <strong>{{ copy.statB }}</strong>
          <span>{{ copy.statBDesc }}</span>
        </div>
        <div>
          <strong>{{ copy.statC }}</strong>
          <span>{{ copy.statCDesc }}</span>
        </div>
      </div>
    </div>

    <div class="portal-login-card portal-card">
      <h2 class="portal-title">{{ copy.loginTitle }}</h2>
      <p class="portal-subtitle">{{ copy.loginSubtitle }}</p>

      <el-form label-position="top" class="portal-login-form">
        <el-form-item :label="copy.username">
          <el-input v-model="form.username" size="large" />
        </el-form-item>
        <el-form-item :label="copy.password">
          <el-input v-model="form.password" type="password" show-password size="large" />
        </el-form-item>
        <el-button type="primary" :loading="loading" size="large" class="portal-login-submit" @click="handleSubmit">
          {{ copy.submit }}
        </el-button>
      </el-form>
    </div>
  </div>
</template>
