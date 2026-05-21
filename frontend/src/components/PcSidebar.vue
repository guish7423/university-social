<template>
  <view class="pc-sidebar" :class="{ collapsed }">
    <view class="sidebar-brand">
      <text class="brand-logo">◈</text>
      <text class="brand-name" v-show="!collapsed">校园社</text>
    </view>
    <view class="sidebar-nav">
      <view
        v-for="item in navItems"
        :key="item.path"
        :class="['nav-item', activePath === item.path && 'nav-active']"
        @click="navigate(item)"
      >
        <text class="nav-icon">{{ item.icon }}</text>
        <text class="nav-label" v-show="!collapsed">{{ item.label }}</text>
      </view>
    </view>
    <view class="sidebar-footer" v-if="userStore.isLogin && !collapsed">
      <view class="user-info" @click="goProfile">
        <image class="user-avatar" :src="userStore.avatar" mode="aspectFill" />
        <view class="user-text">
          <text class="user-name">{{ userStore.nickname }}</text>
          <text class="user-school">{{ userStore.school || '校园社' }}</text>
        </view>
      </view>
    </view>
    <view class="sidebar-toggle" @click="collapsed = !collapsed">
      <text>{{ collapsed ? '▶' : '◀' }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const collapsed = ref(false)
const activePath = ref("")

try {
  const pages = getCurrentPages()
  if (pages.length > 0) {
    const route = (pages[pages.length - 1] as any).$page?.fullPath
    if (route) activePath.value = route
  }
} catch {}

const navItems = [
  { icon: "◈", label: "动态广场", path: "/pages/square/index", tab: true },
  { icon: "◎", label: "校园圈子", path: "/pages/circle/list", tab: true },
  { icon: "♢", label: "树洞", path: "/pages/whisper/index", tab: true },
  { icon: "⌕", label: "搜索", path: "/pages/search/index" },
  { icon: "▤", label: "二手市场", path: "/pages/product/list" },
  { icon: "▣", label: "校园活动", path: "/pages/activity/list" },
  { icon: "★", label: "积分排行", path: "/pages/points/index" },
  { icon: "◉", label: "课程评价", path: "/pages/course/search" },
  { icon: "♢", label: "失物招领", path: "/pages/found/list" },
  { icon: "◯", label: "好友", path: "/pages/friend/index" },
  { icon: "⌕", label: "通知", path: "/pages/notification/index" },
  { icon: "◈", label: "个人中心", path: "/pages/user/index", tab: true },
]

function navigate(item: typeof navItems[0]) {
  activePath.value = item.path
  if (item.tab) {
    uni.switchTab({ url: item.path })
  } else {
    uni.navigateTo({ url: item.path })
  }
}

function goProfile() {
  activePath.value = "/pages/user/index"
  uni.switchTab({ url: "/pages/user/index" })
}
</script>

<style>
@media (min-width: 1024px) {
  .pc-sidebar {
    display: flex;
    flex-direction: column;
    position: fixed;
    top: 0;
    left: 0;
    width: 220px;
    height: 100vh;
    background: linear-gradient(180deg, #1E2A3A 0%, #2A3A4E 100%);
    z-index: 999;
    transition: width 0.3s ease;
    overflow: hidden;
  }
  .pc-sidebar.collapsed { width: 64px; }
  .sidebar-brand {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 24px 20px;
    border-bottom: 1px solid rgba(255,255,255,0.08);
  }
  .brand-logo { font-size: 24px; color: #C67A6A; font-weight: 300; flex-shrink: 0; }
  .brand-name { font-size: 18px; font-weight: 700; color: #fff; white-space: nowrap; }
  .sidebar-nav { flex: 1; overflow-y: auto; padding: 8px 0; }
  .nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 20px;
    cursor: pointer;
    transition: all 0.2s ease;
    border-left: 3px solid transparent;
  }
  .nav-item:hover { background: rgba(255,255,255,0.06); }
  .nav-active { background: rgba(198,122,106,0.15); border-left-color: #C67A6A; }
  .nav-icon { font-size: 18px; width: 24px; text-align: center; flex-shrink: 0; }
  .nav-label { font-size: 14px; color: rgba(255,255,255,0.75); white-space: nowrap; font-weight: 500; }
  .nav-active .nav-label { color: #fff; font-weight: 600; }
  .sidebar-footer { border-top: 1px solid rgba(255,255,255,0.08); padding: 16px; }
  .user-info { display: flex; align-items: center; gap: 10px; cursor: pointer; }
  .user-avatar { width: 36px; height: 36px; border-radius: 50%; flex-shrink: 0; background: rgba(255,255,255,0.1); }
  .user-text { flex: 1; min-width: 0; }
  .user-name { font-size: 13px; font-weight: 600; color: #fff; display: block; }
  .user-school { font-size: 11px; color: rgba(255,255,255,0.5); display: block; }
  .sidebar-toggle {
    cursor: pointer;
    color: rgba(255,255,255,0.3);
    font-size: 12px;
    padding: 8px;
    text-align: center;
  }
}
@media (max-width: 1023px) {
  .pc-sidebar { display: none; }
}
</style>
