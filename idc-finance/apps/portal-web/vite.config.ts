import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import path from "node:path";

export default defineConfig(({ command }) => ({
  base: command === "build" ? "/portal/" : "/",
  plugins: [vue()],
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes("node_modules")) {
            return undefined;
          }

          if (id.includes("element-plus")) {
            return "vendor-element-plus";
          }

          if (id.includes("vue-router") || id.includes("pinia")) {
            return "vendor-router-store";
          }

          return "vendor-misc";
        }
      }
    }
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src")
    },
    extensions: [".ts", ".tsx", ".vue", ".js", ".jsx", ".json"]
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:18080",
        changeOrigin: true
      }
    }
  }
}));
