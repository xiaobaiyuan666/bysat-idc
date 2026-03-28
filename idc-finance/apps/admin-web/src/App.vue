<script setup lang="ts">
import { computed, watchEffect } from "vue";
import { useRoute } from "vue-router";
import { useLocaleStore } from "@/store";
import Layout from "./layout/index.vue";

const route = useRoute();
const localeStore = useLocaleStore();
const useLayout = computed(() => route.path !== "/login");

watchEffect(() => {
  document.documentElement.lang = localeStore.locale;
});
</script>

<template>
  <el-config-provider :locale="localeStore.elementLocale">
    <Layout v-if="useLayout" />
    <RouterView v-else />
  </el-config-provider>
</template>
