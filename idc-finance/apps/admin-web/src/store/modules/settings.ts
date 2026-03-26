import { defineStore } from "pinia";

export const useSettingsStore = defineStore("settings", {
  state: () => ({
    title: "无穷云IDC业务管理系统",
    titleEn: "Infinity Cloud IDC Management System",
    subtitle: "江苏白猿网络科技有限公司-猿创软件开发",
    subtitleEn: "Jiangsu Baiyuan Network Technology Co., Ltd. - Yuanchuang Software Development"
  })
});
