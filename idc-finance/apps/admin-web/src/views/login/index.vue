<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { localeOptions, pickLabel } from "@/locales";
import { useLocaleStore, usePermissionStore, useUserStore } from "@/store";

const router = useRouter();
const localeStore = useLocaleStore();
const userStore = useUserStore();
const permissionStore = usePermissionStore();
const loading = ref(false);

const form = reactive({
  username: "admin",
  password: "Admin123!"
});

const copy = computed(() => ({
  heroTitle: pickLabel(localeStore.locale, "无穷云IDC业务管理系统后台入口", "Infinity Cloud IDC admin entry"),
  heroSubtitle: pickLabel(
    localeStore.locale,
    "江苏白猿网络科技有限公司-猿创软件开发维护，面向财务、工单、服务、客户、资源与自动化运营场景。",
    "Maintained by Jiangsu Baiyuan Network Technology Co., Ltd. - Yuanchuang Software Development for finance, tickets, services, customers, resources, and automation."
  ),
  featureA: pickLabel(localeStore.locale, "运营工作台", "Operations Workbench"),
  featureADesc: pickLabel(localeStore.locale, "集中处理实名认证、逾期账单、到期服务和工单待办。", "Handle identity review, overdue invoices, expiring services, and pending tickets in one place."),
  featureB: pickLabel(localeStore.locale, "财务主链", "Billing Chain"),
  featureBDesc: pickLabel(localeStore.locale, "订单、账单、收款、退款和服务状态回写已经打通。", "Orders, invoices, payments, refunds, and service status are linked together."),
  featureC: pickLabel(localeStore.locale, "资源与自动化中心", "Resources & Automation"),
  featureCDesc: pickLabel(localeStore.locale, "商品配置、资源快照、实例动作和服务工作台已经统一进后台视图。", "Product setup, resource snapshots, instance actions, and service workbenches are unified."),
  loginTitle: pickLabel(localeStore.locale, "登录后台", "Admin Login"),
  loginSubtitle: pickLabel(localeStore.locale, "使用管理员账号进入客户、财务和服务运营工作台。", "Use the administrator account to enter the operations console."),
  defaultAccount: pickLabel(localeStore.locale, "默认账号", "Default account"),
  declarationTitle: pickLabel(localeStore.locale, "版权声明", "Copyright"),
  declarationText: pickLabel(
    localeStore.locale,
    "无穷云IDC业务管理系统由江苏白猿网络科技有限公司-猿创软件开发 100% AI 开发，全部著作权归江苏白猿网络科技有限公司所有。",
    "Infinity Cloud IDC Management System is 100% AI-developed by Jiangsu Baiyuan Network Technology Co., Ltd. - Yuanchuang Software Development, with full copyright owned by Jiangsu Baiyuan Network Technology Co., Ltd."
  ),
  username: pickLabel(localeStore.locale, "用户名", "Username"),
  password: pickLabel(localeStore.locale, "密码", "Password"),
  submit: pickLabel(localeStore.locale, "进入后台", "Enter Admin"),
  success: pickLabel(localeStore.locale, "登录成功", "Login successful"),
  failed: pickLabel(localeStore.locale, "登录失败，请确认 API 服务已经启动。", "Login failed. Confirm the API service is running.")
}));

async function handleSubmit() {
  loading.value = true;
  try {
    await userStore.login(form);
    await permissionStore.load();
    ElMessage.success(copy.value.success);
    router.push("/workbench");
  } catch {
    ElMessage.error(copy.value.failed);
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="login-shell">
    <div class="page-card login-card">
      <section class="login-hero">
        <div class="login-hero__toolbar">
          <div class="login-hero__eyebrow">Infinity Cloud IDC</div>
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
        <h1 class="login-hero__title">{{ copy.heroTitle }}</h1>
        <p class="login-hero__subtitle">{{ copy.heroSubtitle }}</p>

        <div class="login-feature-list">
          <div class="login-feature">
            <div class="login-feature__title">{{ copy.featureA }}</div>
            <div class="login-feature__desc">{{ copy.featureADesc }}</div>
          </div>
          <div class="login-feature">
            <div class="login-feature__title">{{ copy.featureB }}</div>
            <div class="login-feature__desc">{{ copy.featureBDesc }}</div>
          </div>
          <div class="login-feature">
            <div class="login-feature__title">{{ copy.featureC }}</div>
            <div class="login-feature__desc">{{ copy.featureCDesc }}</div>
          </div>
        </div>
      </section>

      <section class="login-form">
        <div class="login-form__head">
          <div class="login-form__title">{{ copy.loginTitle }}</div>
          <div class="login-form__subtitle">{{ copy.loginSubtitle }}</div>
        </div>

        <div class="login-tips">
          <strong>{{ copy.defaultAccount }}</strong><br />
          {{ copy.username }}：`admin`<br />
          {{ copy.password }}：`Admin123!`
        </div>

        <div class="login-tips">
          <strong>{{ copy.declarationTitle }}</strong><br />
          {{ copy.declarationText }}
        </div>

        <el-form label-position="top">
          <el-form-item :label="copy.username">
            <el-input v-model="form.username" size="large" />
          </el-form-item>
          <el-form-item :label="copy.password">
            <el-input v-model="form.password" type="password" size="large" show-password />
          </el-form-item>
          <el-button type="primary" size="large" :loading="loading" style="width: 100%" @click="handleSubmit">
            {{ copy.submit }}
          </el-button>
        </el-form>
      </section>
    </div>
  </div>
</template>
