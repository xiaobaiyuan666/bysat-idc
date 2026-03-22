import { defineStore } from "pinia";

type DeviceType = "desktop" | "mobile";

export const useAppStore = defineStore("pure-app", {
  state: () => ({
    sidebar: {
      opened: true,
    },
    device: "desktop" as DeviceType,
    layout: "vertical",
    fixedHeader: true,
    hideTabs: false,
  }),
  actions: {
    toggleSidebar(opened?: boolean) {
      this.sidebar.opened = typeof opened === "boolean" ? opened : !this.sidebar.opened;
    },
    setDevice(device: DeviceType) {
      this.device = device;
      if (device === "mobile") {
        this.sidebar.opened = false;
      }
    },
    setHideTabs(value: boolean) {
      this.hideTabs = value;
    },
  },
});

export function useAppStoreHook() {
  return useAppStore();
}
