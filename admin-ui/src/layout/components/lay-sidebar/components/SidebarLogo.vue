<script setup lang="ts">
defineProps<{
  collapse: boolean;
}>();
</script>

<template>
  <div class="sidebar-logo">
    <div class="sidebar-logo__mark">ID</div>
    <div v-if="!collapse" class="sidebar-logo__copy">
      <div class="sidebar-logo__title">IDC 公有云业务系统</div>
      <div class="sidebar-logo__subtitle">vue-pure-admin 控制台基座</div>
    </div>
  </div>
</template>

<style scoped>
.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 18px 14px;
}

.sidebar-logo__mark {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  color: #fff;
  font-weight: 800;
  background: linear-gradient(135deg, #2f80ff, #35c6a7);
}

.sidebar-logo__title {
  font-size: 15px;
  font-weight: 700;
  color: #fff;
}

.sidebar-logo__subtitle {
  margin-top: 4px;
  color: rgba(255, 255, 255, 0.55);
  font-size: 12px;
}
</style>
