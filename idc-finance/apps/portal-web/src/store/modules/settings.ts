import { defineStore } from "pinia";

export const useSettingsStore = defineStore("settings", {
  state: () => ({
    title: "白猿科技客户中心",
    titleEn: "BYSAT Client Portal",
    subtitle: "江苏白猿网络科技有限公司 · 猿创软件业务组",
    subtitleEn: "Jiangsu Baiyuan Network Technology Co., Ltd. · Yuanchuang Software Business Group"
  })
});
