import { defineStore } from "pinia";

import { http } from "@/api/http";
import type { AuthUser } from "@/types";

export const useUserStore = defineStore("pure-user", {
  state: () => ({
    user: null as AuthUser | null,
    initialized: false,
  }),
  actions: {
    async fetchMe() {
      try {
        const { data } = await http.get<{ data: AuthUser }>("/auth/me");
        this.user = data.data;
      } catch {
        this.user = null;
      } finally {
        this.initialized = true;
      }
    },
    async ensureUser() {
      if (!this.initialized) {
        await this.fetchMe();
      }
      return this.user;
    },
    async login(payload: { email: string; password: string }) {
      const { data } = await http.post<{ data: AuthUser }>("/auth/login", payload);
      this.user = data.data;
      this.initialized = true;
    },
    async logout() {
      await http.post("/auth/logout");
      this.user = null;
      this.initialized = true;
    },
  },
});

export function useUserStoreHook() {
  return useUserStore();
}
