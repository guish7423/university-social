<template>
  <div class="app-shell">
    <div class="sidebar-section">
      <dv-border-box-12 class="sidebar-box">
        <div class="logo">校园社 管理</div>
        <el-menu router :default-active="$route.path" background-color="transparent" text-color="var(--admin-text-secondary, '#bfcbd9')" active-text-color="#409eff" class="sidebar-menu">
          <el-menu-item index="/dashboard">
            <el-icon><DataBoard /></el-icon>
            <span>数据看板</span>
          </el-menu-item>
          <el-menu-item index="/datascreen">
            <el-icon><Monitor /></el-icon>
            <span>数据大屏</span>
          </el-menu-item>
          <el-menu-item index="/users">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/posts">
            <el-icon><Document /></el-icon>
            <span>内容管理</span>
          </el-menu-item>
          <el-menu-item index="/circles">
            <el-icon><Collection /></el-icon>
            <span>圈子管理</span>
          </el-menu-item>
          <el-menu-item index="/banners">
            <el-icon><Picture /></el-icon>
            <span>Banner管理</span>
          </el-menu-item>
          <el-menu-item index="/announcements">
            <el-icon><Message /></el-icon>
            <span>公告管理</span>
          </el-menu-item>
          <el-menu-item index="/sensitive-words">
            <el-icon><Warning /></el-icon>
            <span>敏感词管理</span>
          </el-menu-item>
          <el-menu-item index="/reports">
            <el-icon><ChatDotSquare /></el-icon>
            <span>举报管理</span>
          </el-menu-item>
        </el-menu>
      </dv-border-box-12>
    </div>
    <div class="main-section">
      <dv-border-box-7 class="header-box">
        <div class="header-content">
          <span class="text-sm text-gray-500">{{ route.meta.title || "" }}</span>
          <div class="header-actions">
            <el-button class="theme-toggle-btn" :icon="currentTheme === 'dark' ? Sunny : Moon" @click="toggleTheme">
              {{ currentTheme === 'dark' ? '浅色模式' : '深色模式' }}
            </el-button>
          </div>
        </div>
      </dv-border-box-7>
      <div class="content-area">
        <transition name="fade" mode="out-in">
          <router-view />
        </transition>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router"
import { ref, onMounted } from "vue"
import { Sunny, Moon } from "@element-plus/icons-vue"
const route = useRoute()

const currentTheme = ref("dark")

function toggleTheme() {
  const next = currentTheme.value === "dark" ? "light" : "dark"
  currentTheme.value = next
  document.documentElement.setAttribute("data-theme", next)
  localStorage.setItem("admin-theme", next)
}

onMounted(() => {
  const saved = localStorage.getItem("admin-theme")
  if (saved) {
    currentTheme.value = saved
    document.documentElement.setAttribute("data-theme", saved)
  }
})
</script>

@import "./styles/animations.css";

<style>
body { margin: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: var(--admin-bg, #f0f2f5); }
.app-shell { display: flex; height: 100vh; overflow: hidden; background: var(--admin-bg, #f0f2f5); }
.sidebar-section { width: 220px; flex-shrink: 0; padding: 4px 0 4px 4px; }
.sidebar-box { height: 100%; padding: 8px; }
.main-section { flex: 1; display: flex; flex-direction: column; padding: 4px 4px 4px 0; overflow: hidden; }
.header-box { margin-bottom: 0; flex-shrink: 0; }
.header-box .dv-border-box-7 { padding: 0 16px; }
.header-content { height: 56px; display: flex; align-items: center; justify-content: space-between; padding: 0 20px; background: var(--admin-bg-header, #fff); }
.content-area { flex: 1; overflow-y: auto; padding: 16px; }
.logo { height: 50px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 18px; font-weight: 600; letter-spacing: 2px; margin-bottom: 8px; }
.el-menu { border-right: none !important; }
.header-actions { display: flex; align-items: center; gap: 8px; }
.theme-toggle-btn { --el-button-bg-color: transparent !important; --el-button-border-color: var(--admin-border-light, #e4e7ed) !important; --el-button-text-color: var(--admin-text-secondary, #606266) !important; }
.theme-toggle-btn:hover { --el-button-text-color: var(--admin-text, #333) !important; --el-button-border-color: var(--admin-text-secondary, #606266) !important; }
</style>
