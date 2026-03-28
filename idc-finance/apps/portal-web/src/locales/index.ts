import en from "element-plus/es/locale/lang/en";
import zhCn from "element-plus/es/locale/lang/zh-cn";

export type LocaleCode = "zh-CN" | "en-US";

export const localeOptions = [
  { value: "zh-CN" as const, label: "中文", shortLabel: "中" },
  { value: "en-US" as const, label: "English", shortLabel: "EN" }
];

export function pickLabel(locale: LocaleCode, zh: string, en: string) {
  return locale === "en-US" ? en : zh;
}

export function resolveMetaTitle(
  meta: { title?: string; titleEn?: string } | undefined,
  locale: LocaleCode,
  fallbackZh = "",
  fallbackEn = ""
) {
  if (!meta) {
    return pickLabel(locale, fallbackZh, fallbackEn);
  }

  return locale === "en-US"
    ? String(meta.titleEn || meta.title || fallbackEn || fallbackZh)
    : String(meta.title || meta.titleEn || fallbackZh || fallbackEn);
}

export function isLegacyMojibake(value: string | undefined) {
  if (!value) return false;
  const patterns = [/锟/, /�/, /鏃犵┓/, /鏈嶅姟/, /璐︽埛/, /鐧诲綍/, /璐㈠姟/, /宸ュ崟/];
  return patterns.some(pattern => pattern.test(value)) || value.includes("??");
}

export const elementLocales = {
  "zh-CN": zhCn,
  "en-US": en
} as const;
