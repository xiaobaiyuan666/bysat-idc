<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { localeOptions, pickLabel } from "@/locales";
import { useLocaleStore, useSettingsStore, useUserStore } from "@/store";

const router = useRouter();
const localeStore = useLocaleStore();
const settingsStore = useSettingsStore();
const userStore = useUserStore();
const loading = ref(false);

const form = reactive({
  username: "portal",
  password: "Portal123!"
});

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, settingsStore.title, settingsStore.titleEn),
  title: pickLabel(localeStore.locale, "云业务客户门户", "Cloud Client Portal"),
  subtitle: pickLabel(
    localeStore.locale,
    "用统一工作台管理服务、财务、订单、工单与账户设置，适合企业客户的日常自助操作。",
    "Use one control plane to manage services, finance, orders, tickets, and account settings."
  ),
  statA: pickLabel(localeStore.locale, "资源联动", "Linked Resources"),
  statADesc: pickLabel(localeStore.locale, "服务与账单上下文统一收口", "Services and invoices stay in one business context"),
  statB: pickLabel(localeStore.locale, "客户自助", "Self Service"),
  statBDesc: pickLabel(localeStore.locale, "常用动作集中在门户工作台", "Common actions are grouped in a single client workbench"),
  statC: pickLabel(localeStore.locale, "统一品牌", "Unified Brand"),
  statCDesc: pickLabel(localeStore.locale, "与运营后台保持同一套信息架构", "Aligned with the operations console"),
  loginTitle: pickLabel(localeStore.locale, "登录客户门户", "Sign in"),
  loginSubtitle: pickLabel(localeStore.locale, "演示账号：portal / Portal123!", "Demo account: portal / Portal123!"),
  declarationTitle: pickLabel(localeStore.locale, "版权声明", "Copyright"),
  declarationText: pickLabel(
    localeStore.locale,
    "无穷云IDC业务管理系统由江苏白猿网络科技有限公司 - 猿创软件开发 100% AI 开发，全部著作权归白猿科技所有。",
    "Infinity Cloud IDC Management System is 100% AI-developed by Jiangsu Baiyuan Network Technology Co., Ltd. - Yuanchuang Software Development. All copyrights belong to Baiyuan Technology."
  ),
  username: pickLabel(localeStore.locale, "账号", "Username"),
  password: pickLabel(localeStore.locale, "密码", "Password"),
  submit: pickLabel(localeStore.locale, "进入客户中心", "Enter Portal"),
  success: pickLabel(localeStore.locale, "登录成功", "Login successful"),
  failed: pickLabel(localeStore.locale, "账号或密码错误", "Invalid username or password")
}));

async function handleSubmit() {
  loading.value = true;
  try {
    await userStore.login(form);
    ElMessage.success(copy.value.success);
    void router.push("/console");
  } catch {
    ElMessage.error(copy.value.failed);
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="portal-login-shell">
    <section class="portal-login-hero">
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
            {{ item.shortLabel }}
          </el-button>
        </el-button-group>
      </div>

      <div class="portal-login-hero__body">
        <div class="meta-chip">{{ settingsStore.subtitle }}</div>
        <h1>{{ copy.title }}</h1>
        <p>{{ copy.subtitle }}</p>
      </div>

      <div class="portal-login-hero__stats">
        <div>
          <strong>{{ copy.statA }}</strong>
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
    </section>

    <section class="portal-login-card">
      <div class="portal-login-card__header">
        <h2>{{ copy.loginTitle }}</h2>
        <p>{{ copy.loginSubtitle }}</p>
      </div>

      <el-form label-position="top" class="portal-login-form" @submit.prevent>
        <el-form-item :label="copy.username">
          <el-input v-model="form.username" size="large" autocomplete="username" />
        </el-form-item>
        <el-form-item :label="copy.password">
          <el-input v-model="form.password" type="password" show-password size="large" autocomplete="current-password" />
        </el-form-item>
        <el-button type="primary" :loading="loading" size="large" class="portal-login-submit" @click="handleSubmit">
          {{ copy.submit }}
        </el-button>
      </el-form>

      <div class="portal-login-card__declaration">
        <strong>{{ copy.declarationTitle }}</strong>
        <span>{{ copy.declarationText }}</span>
      </div>
    </section>
  </div>
</template>
