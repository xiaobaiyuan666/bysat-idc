import { defineStore } from "pinia";

export const useSettingsStore = defineStore("settings", {
  state: () => ({
    title: "白猿科技云业务运营后台",
    titleEn: "BYSAT Cloud Operations Admin",
    subtitle: "江苏白猿网络科技有限公司 · 猿创软件业务组",
    subtitleEn: "Jiangsu Baiyuan Network Technology Co., Ltd. · Yuanchuang Software Business Group"
  })
});
