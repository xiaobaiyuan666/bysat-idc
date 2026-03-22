import { defineStore } from "pinia";

export const useUserStore = defineStore("user", {
  state: () => ({
    token: localStorage.getItem("portal-token") ?? "",
    displayName: localStorage.getItem("portal-display-name") ?? ""
  }),
  getters: {
    isLoggedIn: state => Boolean(state.token)
  },
  actions: {
    async login(payload: { username: string; password: string }) {
      if (payload.username !== "portal" || payload.password !== "Portal123!") {
        throw new Error("invalid credentials");
      }
      this.token = "phase1-portal-token";
      this.displayName = "演示客户";
      localStorage.setItem("portal-token", this.token);
      localStorage.setItem("portal-display-name", this.displayName);
    },
    logout() {
      this.token = "";
      this.displayName = "";
      localStorage.removeItem("portal-token");
      localStorage.removeItem("portal-display-name");
    }
  }
});
