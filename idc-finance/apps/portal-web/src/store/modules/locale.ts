import { defineStore } from "pinia";
import type { LocaleCode } from "@/locales";
import { elementLocales } from "@/locales";

const STORAGE_KEY = "portal-locale";

function getInitialLocale(): LocaleCode {
  const saved = localStorage.getItem(STORAGE_KEY);
  return saved === "en-US" ? "en-US" : "zh-CN";
}

export const useLocaleStore = defineStore("locale", {
  state: () => ({
    locale: getInitialLocale() as LocaleCode
  }),
  getters: {
    isEnglish: state => state.locale === "en-US",
    elementLocale: state => elementLocales[state.locale]
  },
  actions: {
    setLocale(locale: LocaleCode) {
      this.locale = locale;
      localStorage.setItem(STORAGE_KEY, locale);
    },
    toggleLocale() {
      this.setLocale(this.locale === "zh-CN" ? "en-US" : "zh-CN");
    }
  }
});
