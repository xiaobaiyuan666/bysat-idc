import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    ".next/**",
    "out/**",
    "build/**",
    "output/**",
    "vendor/**",
    "public/vendor/**",
    "next-env.d.ts",
    "idc-finance/**",
    "admin-ui/**",
    "vue-pure-admin-ref/**",
    "scui-ref/**",
    "zjmf-business-ref/**",
    ".playwright-cli/**",
    "scan_*.js",
    "tmp-*.js",
    "tmp_*.js",
    "tmp_*.html",
  ]),
]);

export default eslintConfig;
