import { defineStore } from "pinia";
import { login } from "@/api/auth";

export const useUserStore = defineStore("user", {
  state: () => ({
    token: localStorage.getItem("admin-token") ?? "",
    displayName: localStorage.getItem("admin-display-name") ?? ""
  }),
  getters: {
    isLoggedIn: state => Boolean(state.token)
  },
  actions: {
    async login(payload: { username: string; password: string }) {
      const result = await login(payload);
      this.token = result.token;
      this.displayName = result.displayName;
      localStorage.setItem("admin-token", result.token);
      localStorage.setItem("admin-display-name", result.displayName);
    },
    logout() {
      this.token = "";
      this.displayName = "";
      localStorage.removeItem("admin-token");
      localStorage.removeItem("admin-display-name");
    }
  }
});
