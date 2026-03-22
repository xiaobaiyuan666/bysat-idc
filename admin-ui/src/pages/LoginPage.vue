<script setup lang="ts">
import { Monitor, Money, Promotion } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";

import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const authStore = useAuthStore();
const loading = ref(false);
const form = reactive({
  email: "admin@idc.local",
  password: "Admin123!",
});

async function onSubmit() {
  loading.value = true;

  try {
    await authStore.login(form);
    ElMessage.success("登录成功");
    router.push("/workbench");
  } catch {
    ElMessage.error("登录失败，请检查账号和密码");
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="login-shell">
    <section class="login-hero">
      <div class="hero-badge">Vue 3 / Element Plus / vue-pure-admin</div>
      <h1 class="hero-title">IDC 财务业务管理系统</h1>
      <p class="hero-copy">
        一套覆盖客户、订单、账单、收款、工单、资源与审计的运营后台，适合 IDC
        业务、财务与运维协同场景。
      </p>

      <div class="feature-grid">
        <div class="feature-card">
          <el-icon><Money /></el-icon>
          <div>
            <h3>财务闭环</h3>
            <p>订单、账单、收款、续费、逾期、余额流水统一纳入财务链路。</p>
          </div>
        </div>
        <div class="feature-card">
          <el-icon><Monitor /></el-icon>
          <div>
            <h3>资源运营</h3>
            <p>服务生命周期、IP、磁盘、快照、备份、VPC 与安全组集中管理。</p>
          </div>
        </div>
        <div class="feature-card">
          <el-icon><Promotion /></el-icon>
          <div>
            <h3>云平台对接</h3>
            <p>已接入魔方云真实环境，可继续扩展实例、资源和财务同步链路。</p>
          </div>
        </div>
      </div>
    </section>

    <div class="login-panel">
      <el-card shadow="never">
        <div class="login-heading">
          <div class="login-kicker">欢迎登录</div>
          <h2>进入运营管理后台</h2>
          <p>当前已预置演示账号，登录后可直接查看计费引擎、账单收款、工单协作和审计日志。</p>
        </div>

        <el-form label-position="top" @submit.prevent="onSubmit">
          <el-form-item label="邮箱账号">
            <el-input v-model="form.email" placeholder="admin@idc.local" />
          </el-form-item>
          <el-form-item label="登录密码">
            <el-input
              v-model="form.password"
              type="password"
              show-password
              placeholder="请输入密码"
            />
          </el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            class="login-button"
            @click="onSubmit"
          >
            登录系统
          </el-button>
        </el-form>

        <div class="login-tips">
          <div>默认管理员：`admin@idc.local`</div>
          <div>默认密码：`Admin123!`</div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
.login-shell {
  min-height: 100vh;
  display: grid;
  gap: 24px;
  padding: 24px;
  grid-template-columns: minmax(0, 1.15fr) minmax(360px, 440px);
}

.login-hero {
  padding: 56px;
  border-radius: 28px;
  color: white;
  background:
    radial-gradient(circle at top right, rgba(35, 198, 170, 0.24), transparent 22%),
    linear-gradient(135deg, #121b2c, #1d2d47 65%, #235ec8);
  box-shadow: 0 36px 100px rgba(18, 27, 44, 0.2);
}

.hero-badge {
  display: inline-flex;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.12);
  font-size: 13px;
}

.hero-title {
  margin: 20px 0 0;
  font-size: 54px;
  line-height: 1.08;
  letter-spacing: -0.04em;
}

.hero-copy {
  max-width: 720px;
  margin-top: 18px;
  color: rgba(255, 255, 255, 0.78);
  font-size: 17px;
  line-height: 1.9;
}

.feature-grid {
  display: grid;
  gap: 18px;
  margin-top: 42px;
}

.feature-card {
  display: grid;
  grid-template-columns: 36px minmax(0, 1fr);
  gap: 14px;
  padding: 18px 20px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(12px);
}

.feature-card h3 {
  margin: 0;
  font-size: 18px;
}

.feature-card p {
  margin: 8px 0 0;
  color: rgba(255, 255, 255, 0.68);
  line-height: 1.7;
}

.login-panel {
  display: flex;
  align-items: center;
}

.login-heading {
  margin-bottom: 24px;
}

.login-kicker {
  color: var(--primary);
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.16em;
}

.login-heading h2 {
  margin: 14px 0 0;
  font-size: 34px;
}

.login-heading p {
  margin: 12px 0 0;
  color: var(--text-secondary);
  line-height: 1.8;
}

.login-button {
  width: 100%;
  height: 46px;
  margin-top: 12px;
}

.login-tips {
  margin-top: 18px;
  padding: 14px 16px;
  border-radius: 16px;
  background: var(--bg-soft);
  color: var(--text-secondary);
  font-size: 13px;
  line-height: 1.8;
}

@media (max-width: 980px) {
  .login-shell {
    grid-template-columns: 1fr;
  }

  .login-hero {
    padding: 28px;
  }

  .hero-title {
    font-size: 38px;
  }
}
</style>
