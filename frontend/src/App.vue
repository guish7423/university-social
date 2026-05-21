<script setup lang="ts">
import { onLaunch } from "@dcloudio/uni-app"
import { useUserStore } from "@/store/user"

onLaunch(() => {
  // #ifdef H5
  const userStore = useUserStore()

  const navItems = [
    { icon: "◈", label: "动态广场", path: "/pages/square/index", tab: true },
    { icon: "◎", label: "校园圈子", path: "/pages/circle/list", tab: true },
    { icon: "♢", label: "树洞", path: "/pages/whisper/index", tab: true },
    { icon: "⌕", label: "搜索", path: "/pages/search/index" },
    { icon: "▤", label: "二手市场", path: "/pages/product/list" },
    { icon: "▣", label: "校园活动", path: "/pages/activity/list" },
    { icon: "★", label: "积分排行", path: "/pages/points/index" },
    { icon: "▤", label: "课程评价", path: "/pages/course/search" },
    { icon: "⌕", label: "失物招领", path: "/pages/found/list" },
    { icon: "♢", label: "好友", path: "/pages/friend/index" },
    { icon: "◉", label: "通知", path: "/pages/notification/index" },
    { icon: "◎", label: "个人中心", path: "/pages/user/index", tab: true },
  ]

  function navigate(path: string, tab: boolean) {
    if (tab) uni.switchTab({ url: path })
    else uni.navigateTo({ url: path })
  }

  const sidebar = document.createElement("div")
  sidebar.className = "pc-sidebar"

  sidebar.innerHTML = [
    '<div class="sidebar-brand"><span class="brand-name">校园社</span></div>',
    '<div class="sidebar-nav">',
    ...navItems.map(item => {
      const active = item.path === location.hash.replace("#", "") ? " nav-active" : ""
      return `<div class="nav-item${active}" data-path="${item.path}" data-tab="${item.tab}"><span class="nav-icon">${item.icon}</span><span class="nav-label">${item.label}</span></div>`
    }),
    '</div>',
    '<div class="sidebar-footer" id="sidebar-footer" style="display:none">',
    '<div class="user-info" data-path="/pages/user/index" data-tab="true">',
    '<img class="user-avatar" id="sidebar-avatar" alt="" />',
    '<div class="user-text">',
    '<span class="user-name" id="sidebar-nickname"></span>',
    '<span class="user-school" id="sidebar-school">校园社</span>',
    '</div></div></div>',
  ].join("")

  function updateActive(path: string) {
    sidebar.querySelectorAll(".nav-item").forEach(el => {
      el.classList.toggle("nav-active", (el as HTMLElement).dataset.path === path)
    })
  }

  sidebar.addEventListener("click", (e) => {
    const item = (e.target as HTMLElement).closest(".nav-item, .user-info") as HTMLElement
    if (!item) return
    const path = item.dataset.path
    const tab = item.dataset.tab === "true"
    if (path) {
      updateActive(path)
      navigate(path, tab)
    }
  })

  window.addEventListener("hashchange", () => {
    updateActive(location.hash.replace("#", ""))
  })

  document.body.prepend(sidebar)

  setTimeout(() => {
    const footer = document.getElementById("sidebar-footer")
    const avatar = document.getElementById("sidebar-avatar") as HTMLImageElement
    const nickname = document.getElementById("sidebar-nickname")
    const school = document.getElementById("sidebar-school")
    if (footer && userStore.isLogin) {
      footer.style.display = ""
      if (avatar && userStore.avatar) avatar.src = userStore.avatar
      if (nickname) nickname.textContent = userStore.nickname
      if (school) school.textContent = userStore.school || "校园社"
    }
  }, 100)
  // #endif
})
</script>

<template><view /></template>

<style>
page { background: var(--color-canvas, #F7F4F0); font-family: 'Noto Sans SC', 'PingFang SC', -apple-system, BlinkMacSystemFont, sans-serif; overscroll-behavior: none; -webkit-font-smoothing: antialiased; }
@media (min-width: 1024px) {
  .uni-tabbar, .uni-tabbar-bottom, uni-tabbar { display: none !important; }
  .uni-page-head { display: none !important; }
  .uni-page-body { margin-left: 220px; padding-top: 0; }
  ::-webkit-scrollbar { display: block; width: 6px; }
  ::-webkit-scrollbar-thumb { background: #B8C2CE; border-radius: 3px; }
  ::-webkit-scrollbar-track { background: transparent; }
}
@media (min-width: 1440px) {
  .uni-page-body { margin-left: 220px; max-width: 1200px; }
}
.pc-sidebar { display: none; flex-direction: column; position: fixed; top: 0; left: 0; width: 220px; height: 100vh; background: #1E2A3A; z-index: 999; overflow-y: auto; }
@media (min-width: 1024px) { .pc-sidebar { display: flex; } }
.sidebar-brand { padding: 32px 24px 24px; border-bottom: 1px solid rgba(255,255,255,0.06); }
.brand-name { font-size: 22px; font-weight: 600; color: #C67A6A; letter-spacing: 2px; }
.sidebar-nav { flex: 1; padding: 12px 0; }
.nav-item { display: flex; align-items: center; gap: 12px; padding: 10px 24px; cursor: pointer; transition: all 0.15s ease; margin: 1px 8px; border-radius: 8px; border-left: 3px solid transparent; }
.nav-item:hover { background: rgba(255,255,255,0.06); }
.nav-active { background: rgba(198,122,106,0.15); border-left-color: #C67A6A; }
.nav-active .nav-label { color: #fff; font-weight: 500; }
.nav-icon { font-size: 14px; width: 20px; text-align: center; color: rgba(255,255,255,0.4); }
.nav-active .nav-icon { color: #C67A6A; }
.nav-label { font-size: 14px; color: rgba(255,255,255,0.55); white-space: nowrap; font-weight: 400; transition: color 0.15s; }
.sidebar-footer { border-top: 1px solid rgba(255,255,255,0.06); padding: 16px; }
.user-info { display: flex; align-items: center; gap: 10px; cursor: pointer; padding: 8px 12px; border-radius: 8px; transition: background 0.15s; }
.user-info:hover { background: rgba(255,255,255,0.06); }
.user-avatar { width: 32px; height: 32px; border-radius: 50%; background: rgba(255,255,255,0.08); border: 1.5px solid rgba(255,255,255,0.12); }
.user-name { font-size: 13px; font-weight: 500; color: rgba(255,255,255,0.85); }
.user-school { font-size: 11px; color: rgba(255,255,255,0.4); }
</style>
