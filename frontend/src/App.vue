<script setup lang="ts">
import { onLaunch } from "@dcloudio/uni-app"
import { useUserStore } from "@/store/user"

onLaunch(() => {
  // #ifdef H5
  const userStore = useUserStore()

  const navItems = [
    { icon: "\u{1F338}", label: "\u52A8\u6001\u5E7F\u573A", path: "/pages/square/index", tab: true },
    { icon: "\u{1F49C}", label: "\u6821\u56ED\u5708\u5B50", path: "/pages/circle/list", tab: true },
    { icon: "\u{1F33F}", label: "\u6811\u6D1E", path: "/pages/whisper/index", tab: true },
    { icon: "\u{1F50D}", label: "\u641C\u7D22", path: "/pages/search/index" },
    { icon: "\u{1F6CD}\uFE0F", label: "\u4E8C\u624B\u5E02\u573A", path: "/pages/product/list" },
    { icon: "\u{1F3AF}", label: "\u6821\u56ED\u6D3B\u52A8", path: "/pages/activity/list" },
    { icon: "\u2728", label: "\u79EF\u5206\u6392\u884C", path: "/pages/points/index" },
    { icon: "\u{1F4DA}", label: "\u8BFE\u7A0B\u8BC4\u4EF7", path: "/pages/course/search" },
    { icon: "\u{1F50E}", label: "\u5931\u7269\u62DB\u9886", path: "/pages/found/list" },
    { icon: "\u{1F465}", label: "\u597D\u53CB", path: "/pages/friend/index" },
    { icon: "\u{1F514}", label: "\u901A\u77E5", path: "/pages/notification/index" },
    { icon: "\u{1F464}", label: "\u4E2A\u4EBA\u4E2D\u5FC3", path: "/pages/user/index", tab: true },
  ]

  function navigate(path: string, tab: boolean) {
    if (tab) uni.switchTab({ url: path })
    else uni.navigateTo({ url: path })
  }

  const sidebar = document.createElement("div")
  sidebar.className = "pc-sidebar"
  sidebar.style.background = "linear-gradient(135deg, rgba(30,27,46,0.95) 0%, rgba(88,28,135,0.9) 50%, rgba(30,27,46,0.95) 100%)"
  sidebar.style.backdropFilter = "blur(20px)"
  sidebar.style.borderRight = "1px solid rgba(192,132,252,0.12)"

  sidebar.innerHTML = [
    '<div class="sidebar-brand"><span class="brand-logo">\u{1F338}</span><span class="brand-name">\u68A6\u56ED\u793E</span></div>',
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
    '<span class="user-school" id="sidebar-school">\u6821\u56ED\u793E</span>',
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

  const style = document.createElement("style")
  style.textContent = `
    @keyframes aurora {
      0% { background-position: 0% 50%; }
      50% { background-position: 100% 50%; }
      100% { background-position: 0% 50%; }
    }
    .pc-sidebar::before {
      content: '';
      position: absolute;
      top: -50%;
      left: -50%;
      width: 200%;
      height: 200%;
      background: radial-gradient(circle at 30% 20%, rgba(192,132,252,0.08) 0%, transparent 50%),
                  radial-gradient(circle at 70% 80%, rgba(249,168,212,0.06) 0%, transparent 50%);
      pointer-events: none;
      z-index: 0;
    }
    .sidebar-brand, .sidebar-nav, .sidebar-footer { position: relative; z-index: 1; }
    .pc-sidebar::-webkit-scrollbar { width: 3px; }
    .pc-sidebar::-webkit-scrollbar-thumb { background: rgba(192,132,252,0.3); border-radius: 2px; }
  `
  document.head.appendChild(style)

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
      if (school) school.textContent = userStore.school || "\u6821\u56ED\u793E"
    }
  }, 100)
  // #endif
})
</script>

<template><view /></template>

<style>
page { background: var(--color-canvas, #faf5ff); font-family: 'Noto Sans SC', 'PingFang SC', -apple-system, BlinkMacSystemFont, sans-serif; overscroll-behavior: none; -webkit-font-smoothing: antialiased; }
@media (min-width: 1024px) {
  .uni-tabbar, .uni-tabbar-bottom, uni-tabbar { display: none !important; }
  .uni-page-head { display: none !important; }
  .uni-page-body { margin-left: 220px; padding-top: 0; }
  ::-webkit-scrollbar { display: block; width: 6px; }
  ::-webkit-scrollbar-thumb { background: rgba(192,132,252,0.3); border-radius: 3px; }
  ::-webkit-scrollbar-track { background: transparent; }
}
@media (min-width: 1440px) {
  .uni-page-body { margin-left: 220px; max-width: 1200px; }
}
.pc-sidebar { display: none; flex-direction: column; position: fixed; top: 0; left: 0; width: 220px; height: 100vh; z-index: 999; overflow-y: auto; }
@media (min-width: 1024px) { .pc-sidebar { display: flex; } }
.sidebar-brand { display: flex; align-items: center; gap: 10px; padding: 28px 20px 20px; border-bottom: 1px solid rgba(192,132,252,0.12); }
.brand-logo { font-size: 28px; }
.brand-name { font-size: 20px; font-weight: 700; background: linear-gradient(135deg, #C084FC 0%, #F9A8D4 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; font-family: 'ZCOOL KuaiLe', 'PingFang SC', sans-serif; }
.sidebar-nav { flex: 1; padding: 8px 0; }
.nav-item { display: flex; align-items: center; gap: 12px; padding: 13px 20px; cursor: pointer; border-left: 3px solid transparent; transition: all 0.2s ease; margin: 2px 8px; border-radius: 12px; }
.nav-item:hover { background: rgba(192,132,252,0.12); transform: translateX(4px); }
.nav-active { background: rgba(192,132,252,0.2); border-left-color: #C084FC; }
.nav-active .nav-label { color: #fff; font-weight: 600; }
.nav-icon { font-size: 18px; width: 24px; text-align: center; }
.nav-label { font-size: 14px; color: rgba(255,255,255,0.6); white-space: nowrap; font-weight: 400; transition: color 0.2s; }
.sidebar-footer { border-top: 1px solid rgba(192,132,252,0.12); padding: 16px; }
.user-info { display: flex; align-items: center; gap: 10px; cursor: pointer; padding: 8px; border-radius: 12px; transition: background 0.2s; }
.user-info:hover { background: rgba(192,132,252,0.12); }
.user-avatar { width: 36px; height: 36px; border-radius: 50%; background: rgba(192,132,252,0.15); border: 2px solid rgba(192,132,252,0.3); }
.user-name { font-size: 13px; font-weight: 600; color: #fff; }
.user-school { font-size: 11px; color: rgba(255,255,255,0.5); }
</style>
