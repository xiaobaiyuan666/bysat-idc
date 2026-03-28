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
      const response = await fetch("/api/v1/portal/auth/login/", {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        throw new Error("invalid credentials");
      }

      const payloadJson = (await response.json()) as {
        data?: {
          token?: string;
          displayName?: string;
        };
      };
      this.token = payloadJson.data?.token ?? "";
      this.displayName = payloadJson.data?.displayName ?? "演示客户";
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
