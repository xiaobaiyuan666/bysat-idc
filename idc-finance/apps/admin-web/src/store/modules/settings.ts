import { defineStore } from "pinia";

export const useSettingsStore = defineStore("settings", {
  state: () => ({
    title: "IDC 财务后台",
    titleEn: "IDC Finance Admin",
    subtitle: "公有云业务与财务运营中心",
    subtitleEn: "Public cloud operations and finance control center"
  })
});
