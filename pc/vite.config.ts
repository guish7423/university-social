import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import { VitePWA } from 'vite-plugin-pwa'
import path from "path"

export default defineConfig({
  plugins: [vue(),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['icons/*.svg'],
      manifest: {
        name: '大学社交',
        short_name: '校园社',
        description: '校园社交平台',
        theme_color: '#C67A6A',
        background_color: '#1a1a2e',
        display: 'standalone',
        icons: [
          { src: 'icons/icon-192x192.svg', sizes: '192x192', type: 'image/svg+xml' },
          { src: 'icons/icon-512x512.svg', sizes: '512x512', type: 'image/svg+xml' },
        ],
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg,woff2}'],
        runtimeCaching: [
          {
            urlPattern: /^https?:\/\/.*\/api\/v1\/posts/,
            handler: 'NetworkFirst',
          },
          {
            urlPattern: /^https?:\/\/.*\/uploads\//,
            handler: 'CacheFirst',
          },
        ],
      },
    }),
  ],
  base: process.env.VITE_BASE ? process.env.VITE_BASE : '/',
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
  server: {
    port: 5174,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        ws: true,
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
