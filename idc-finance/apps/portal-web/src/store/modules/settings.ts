import { defineStore } from "pinia";

export const useSettingsStore = defineStore("settings", {
  state: () => ({
    title: "IDC 客户中心",
    titleEn: "IDC Client Area",
    subtitle: "云业务自助门户",
    subtitleEn: "Self-service cloud portal"
  })
});
