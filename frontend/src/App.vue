<script setup lang="ts">
import { onLaunch } from "@dcloudio/uni-app"
import { useUserStore } from "@/store/user"

onLaunch(() => {
  // #ifdef H5
  const userStore = useUserStore()

  const navItems = [
    { icon: "🏠", label: "动态广场", path: "/pages/square/index", tab: true },
    { icon: "🔵", label: "校园圈子", path: "/pages/circle/list", tab: true },
    { icon: "🌳", label: "树洞", path: "/pages/whisper/index", tab: true },
    { icon: "🔍", label: "搜索", path: "/pages/search/index" },
    { icon: "🛒", label: "二手市场", path: "/pages/product/list" },
    { icon: "🎯", label: "校园活动", path: "/pages/activity/list" },
    { icon: "⭐", label: "积分排行", path: "/pages/points/index" },
    { icon: "📚", label: "课程评价", path: "/pages/course/search" },
    { icon: "🔎", label: "失物招领", path: "/pages/found/list" },
    { icon: "👥", label: "好友", path: "/pages/friend/index" },
    { icon: "🔔", label: "通知", path: "/pages/notification/index" },
    { icon: "👤", label: "个人中心", path: "/pages/user/index", tab: true },
  ]

  function navigate(path: string, tab: boolean) {
    if (tab) uni.switchTab({ url: path })
    else uni.navigateTo({ url: path })
  }

  const sidebar = document.createElement("div")
  sidebar.className = "pc-sidebar"
  sidebar.innerHTML = `
    <div class="sidebar-brand">
      <span class="brand-logo">🎓</span>
      <span class="brand-name">校园社</span>
    </div>
    <div class="sidebar-nav">
      ${navItems.map(item => `
        <div class="nav-item" data-path="${item.path}" data-tab="${item.tab}">
          <span class="nav-icon">${item.icon}</span>
          <span class="nav-label">${item.label}</span>
        </div>
      `).join("")}
    </div>
    <div class="sidebar-footer" id="sidebar-footer" style="display:none">
      <div class="user-info" data-path="/pages/user/index" data-tab="true">
        <img class="user-avatar" id="sidebar-avatar" src="" />
        <div class="user-text">
          <span class="user-name" id="sidebar-nickname"></span>
          <span class="user-school" id="sidebar-school">校园社</span>
        </div>
      </div>
    </div>
  `

  sidebar.addEventListener("click", (e) => {
    const item = (e.target as HTMLElement).closest(".nav-item, .user-info") as HTMLElement
    if (!item) return
    const path = item.dataset.path
    const tab = item.dataset.tab === "true"
    if (path) navigate(path, tab)
  })

  document.body.prepend(sidebar)

  // Watch user store for updates
  const unwatch = userStore.$subscribe(() => {
    const footer = document.getElementById("sidebar-footer")
    const avatar = document.getElementById("sidebar-avatar") as HTMLImageElement
    const nickname = document.getElementById("sidebar-nickname")
    const school = document.getElementById("sidebar-school")
    if (footer && userStore.isLogin) {
      footer.style.display = "flex"
      if (avatar && userStore.avatar) avatar.src = userStore.avatar
      if (nickname) nickname.textContent = userStore.nickname
      if (school) school.textContent = userStore.school || "校园社"
    } else if (footer) {
      footer.style.display = "none"
    }
  })
  // #endif
})
</script>

<template><view /></template>

<style>
page { background: #f0f2f5; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", sans-serif; overscroll-behavior: none; -webkit-font-smoothing: antialiased; }

@media (min-width: 1024px) {
  .uni-tabbar, .uni-tabbar-bottom, uni-tabbar { display: none !important; }
  .uni-page-head { display: none !important; }
  .uni-page-body { margin-left: 220px; max-width: calc(100vw - 220px); padding-top: 0; }
  ::-webkit-scrollbar { display: block; width: 6px; }
  ::-webkit-scrollbar-thumb { background: #c0c4cc; border-radius: 3px; }
  ::-webkit-scrollbar-track { background: transparent; }
}
@media (min-width: 1440px) { .uni-page-body { margin-left: 220px; max-width: 1200px; } }

.pc-sidebar { display: none; flex-direction: column; position: fixed; top: 0; left: 0; width: 220px; height: 100vh; background: linear-gradient(180deg,#1a1a2e 0%,#16213e 50%,#0f3460 100%); z-index: 999; overflow-y: auto; }
@media (min-width: 1024px) { .pc-sidebar { display: flex; } }
.sidebar-brand { display: flex; align-items: center; gap: 10px; padding: 24px 20px; border-bottom: 1px solid rgba(255,255,255,0.08); }
.brand-logo { font-size: 24px; }
.brand-name { font-size: 18px; font-weight: 700; color: #fff; }
.sidebar-nav { flex: 1; padding: 8px 0; }
.nav-item { display: flex; align-items: center; gap: 10px; padding: 12px 20px; cursor: pointer; border-left: 3px solid transparent; transition: background 0.2s; }
.nav-item:hover { background: rgba(255,255,255,0.06); }
.nav-icon { font-size: 18px; width: 24px; text-align: center; }
.nav-label { font-size: 14px; color: rgba(255,255,255,0.75); white-space: nowrap; }
.sidebar-footer { border-top: 1px solid rgba(255,255,255,0.08); padding: 16px; display: flex; }
.user-info { display: flex; align-items: center; gap: 10px; cursor: pointer; }
.user-avatar { width: 36px; height: 36px; border-radius: 50%; background: rgba(255,255,255,0.1); }
.user-name { font-size: 13px; font-weight: 600; color: #fff; }
.user-school { font-size: 11px; color: rgba(255,255,255,0.5); }
</style>
