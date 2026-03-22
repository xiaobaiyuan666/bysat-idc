<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import type { NoVncModule, NoVncRfb } from "@/types/novnc";

const route = useRoute();
const containerRef = ref<HTMLDivElement | null>(null);
const status = ref("正在连接控制台...");
const error = ref("");
let rfb: NoVncRfb | undefined;

function readQueryParam(key: string) {
  const value = route.query[key];
  return typeof value === "string" ? value : "";
}

const wsUrl = computed(() => readQueryParam("ws"));
const password = computed(() => readQueryParam("password"));
const ip = computed(() => readQueryParam("ip"));
const title = computed(() => readQueryParam("title"));

function resolveNoVncModuleUrl() {
  // Keep noVNC as a local static bundle. The npm package currently pulls in
  // code that breaks this Vite build pipeline during production bundling.
  const basePath = import.meta.env.BASE_URL || "/";
  return `${basePath.replace(/\/?$/, "/")}vendor/novnc/rfb.bundle.mjs`;
}

onMounted(async () => {
  if (!containerRef.value || !wsUrl.value) {
    status.value = "缺少 VNC 连接参数";
    return;
  }

  try {
    const moduleUrl = resolveNoVncModuleUrl();
    const { default: RFB } = (await import(/* @vite-ignore */ moduleUrl)) as NoVncModule;

    if (!containerRef.value) {
      return;
    }

    rfb = new RFB(
      containerRef.value,
      wsUrl.value,
      password.value ? { credentials: { password: password.value } } : undefined
    );
    rfb.scaleViewport = true;
    rfb.resizeSession = true;
    rfb.clipViewport = false;
    rfb.background = "#08111f";

    rfb.addEventListener("connect", () => {
      status.value = "控制台已连接";
      error.value = "";
    });

    rfb.addEventListener("disconnect", (event: { detail?: { clean?: boolean } }) => {
      status.value = event?.detail?.clean ? "控制台已断开" : "控制台连接中断";
    });

    rfb.addEventListener("credentialsrequired", () => {
      if (password.value && rfb) {
        rfb.sendCredentials({ password: password.value });
        return;
      }
      status.value = "控制台需要 VNC 密码";
    });
  } catch (connectError) {
    error.value = connectError instanceof Error ? connectError.message : "VNC 连接初始化失败";
    status.value = "控制台初始化失败";
  }
});

onBeforeUnmount(() => {
  if (!rfb) {
    return;
  }
  try {
    rfb.disconnect();
  } catch {
    // Ignore disconnect errors during navigation.
  }
});
</script>

<template>
  <main class="vnc-page">
    <section class="vnc-hero">
      <div>
        <p class="vnc-eyebrow">Remote Console</p>
        <h1>实例控制台</h1>
        <p class="vnc-subtitle">通过新后台直接承接魔方云 VNC 会话，不再依赖旧版 Next 后台。</p>
      </div>
      <div class="vnc-meta">
        <span>实例：{{ title || "-" }}</span>
        <span>目标 IP：{{ ip || "-" }}</span>
      </div>
    </section>

    <section class="vnc-status-card">
      <div>
        <strong>{{ status }}</strong>
        <p>连接地址：{{ wsUrl || "-" }}</p>
      </div>
      <span class="vnc-status-badge">{{ wsUrl ? "WebSocket VNC" : "Missing Params" }}</span>
    </section>

    <section v-if="error" class="vnc-error">
      {{ error }}
    </section>

    <section v-if="!wsUrl" class="vnc-empty">
      缺少 VNC 连接参数，请返回服务详情重新点击“获取控制台”。
    </section>
    <section v-else ref="containerRef" class="vnc-canvas" />
  </main>
</template>

<style scoped>
.vnc-page {
  min-height: 100vh;
  padding: 28px;
  background:
    radial-gradient(circle at top left, rgba(46, 109, 255, 0.16), transparent 24%),
    linear-gradient(180deg, #eef4fb 0%, #f7faff 100%);
}

.vnc-hero,
.vnc-status-card,
.vnc-empty,
.vnc-error {
  max-width: 1440px;
  margin: 0 auto 18px;
}

.vnc-hero {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  padding: 24px 28px;
  border: 1px solid #d8e3f2;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.88);
  box-shadow: 0 20px 48px rgba(15, 23, 42, 0.08);
  backdrop-filter: blur(14px);
}

.vnc-eyebrow {
  margin: 0 0 10px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #3556a8;
}

.vnc-hero h1 {
  margin: 0;
  font-size: 34px;
}

.vnc-subtitle {
  margin: 10px 0 0;
  max-width: 720px;
  color: #5d6e88;
  line-height: 1.7;
}

.vnc-meta {
  display: grid;
  gap: 10px;
  align-content: start;
  min-width: 220px;
  color: #27406d;
  font-size: 13px;
}

.vnc-status-card {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  padding: 18px 22px;
  border: 1px solid #d8e3f2;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.82);
  box-shadow: 0 16px 36px rgba(15, 23, 42, 0.06);
}

.vnc-status-card strong {
  display: block;
  font-size: 16px;
}

.vnc-status-card p {
  margin: 6px 0 0;
  color: #60718b;
  font-size: 13px;
  word-break: break-all;
}

.vnc-status-badge {
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(35, 85, 214, 0.1);
  color: #1e4fd0;
  font-size: 12px;
  font-weight: 700;
}

.vnc-error,
.vnc-empty {
  padding: 16px 18px;
  border-radius: 18px;
  border: 1px solid #eac3bc;
  background: rgba(220, 38, 38, 0.08);
  color: #a03425;
}

.vnc-canvas {
  max-width: 1440px;
  min-height: 760px;
  margin: 0 auto;
  overflow: hidden;
  border-radius: 28px;
  border: 1px solid rgba(13, 24, 43, 0.14);
  background: #08111f;
  box-shadow: 0 28px 60px rgba(8, 17, 31, 0.28);
}

@media (max-width: 960px) {
  .vnc-page {
    padding: 18px;
  }

  .vnc-hero,
  .vnc-status-card {
    flex-direction: column;
    align-items: flex-start;
  }

  .vnc-canvas {
    min-height: 560px;
  }
}
</style>
