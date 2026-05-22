import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import path from "path"

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
  server: {
    port: 5174,
    proxy: {
      "/api": {
        target: "http://localhost:8081",
        changeOrigin: true,
      },
    },
  },
  build: {
    chunkSizeWarningLimit: 600,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules")) {
            if (id.includes("vue") || id.includes("pinia") || id.includes("vue-router")) return "vue-vendor"
            if (id.includes("element-plus") || id.includes("@element-plus")) return "element-plus"
            if (id.includes("echarts") || id.includes("zrender")) return "echarts"
            if (id.includes("axios")) return "vendor"
            if (id.includes("markdown") || id.includes("highlight")) return "vendor"
          }
        },
      },
    },
  },
})
