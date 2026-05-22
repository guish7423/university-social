<template>
  <el-container class="app-container">
    <el-aside width="220px" class="app-aside">
      <div class="logo">校园社 管理</div>
      <el-menu router :default-active="$route.path" background-color="#304156" text-color="#bfcbd9" active-text-color="#409eff">
        <el-menu-item index="/dashboard">
          <el-icon><DataBoard /></el-icon>
          <span>数据看板</span>
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
    </el-aside>
    <el-container>
      <el-header class="app-header">
        <span class="text-sm text-gray-500">{{ route.meta.title || "" }}</span>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router"
const route = useRoute()
</script>

<style>
body { margin: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
.app-container { height: 100vh; }
.app-aside { background-color: #304156; }
.app-aside .logo { height: 60px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 18px; font-weight: 600; letter-spacing: 2px; }
.app-header { background: #fff; border-bottom: 1px solid #e6e6e6; display: flex; align-items: center; padding: 0 20px; height: 60px !important; }
</style>
