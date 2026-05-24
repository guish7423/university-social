<template>
  <div class="layout" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
    <div v-if="isMobile && sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false" />
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed, 'mobile-open': isMobile && sidebarOpen }">
      <div class="logo" @click="$router.push('/home')">
        <div class="logo-icon">
          <svg width="28" height="28" viewBox="0 0 28 28" fill="none">
            <rect width="28" height="28" rx="6" fill="#C67A6A"/>
            <path d="M8 14h12M14 8v12" stroke="#fff" stroke-width="2.5" stroke-linecap="round"/>
          </svg>
        </div>
        <span v-show="!sidebarCollapsed" class="logo-text">校园社</span>
        <span v-show="!sidebarCollapsed" class="logo-badge">洛阳高校</span>
      </div>

      <div class="sidebar-collapse-btn" @click="toggleSidebar">
        <el-icon :class="{ rotated: sidebarCollapsed }"><Fold /></el-icon>
      </div>

      <nav class="nav-groups">
        <div v-for="group in navGroups" :key="group.label" class="nav-group">
          <div class="nav-group-label">{{ group.label }}</div>
          <el-menu :default-active="activeRoute" router class="nav-menu">
            <template v-for="item in group.items" :key="item.path">
              <el-menu-item v-if="!item.children" :index="item.path">
                <el-icon><component :is="item.icon" /></el-icon>
                  <span v-show="!sidebarCollapsed">{{ item.label }}</span>
              </el-menu-item>
              <el-sub-menu v-else :index="item.path">
                <template #title>
                  <el-icon><component :is="item.icon" /></el-icon>
                  <span v-show="!sidebarCollapsed">{{ item.label }}</span>
                </template>
                <el-menu-item v-for="child in item.children" :key="child.path" :index="child.path">
                  {{ child.label }}
                </el-menu-item>
              </el-sub-menu>
            </template>
          </el-menu>
        </div>
      </nav>

      <div class="sidebar-footer">
        <div class="user-card" @click="$router.push('/user/' + userStore.userId)">
          <el-avatar :size="36" :src="userStore.avatar" />
          <div class="user-info">
            <div class="nickname">{{ userStore.nickname }}</div>
            <div class="school">{{ userStore.school || '未设置学校' }}</div>
          </div>
          <el-icon class="user-arrow"><ArrowRight /></el-icon>
        </div>
      </div>
    </aside>

    <main class="main-area">
      <header class="top-bar">
        <div class="top-bar-left">
          <el-button text class="hamburger-btn" @click="sidebarOpen = true" v-if="isMobile">
            <el-icon size="20"><Expand /></el-icon>
          </el-button>
          <el-button text class="collapse-btn" @click="toggleSidebar" v-else>
            <el-icon><Fold /></el-icon>
          </el-button>
          <el-button text class="search-btn" @click="$router.push('/search')">
            <el-icon><Search /></el-icon>
            <span>搜索校园社...</span>
          </el-button>
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notif-badge">
            <el-button text @click="$router.push('/notifications')">
              <el-icon><Bell /></el-icon>
            </el-button>
          </el-badge>
        </div>
        <div class="top-bar-right">
          <el-tooltip content="发布动态" placement="bottom">
            <el-button type="primary" :icon="Edit" circle @click="$router.push('/post/create')" />
          </el-tooltip>
          <el-dropdown trigger="click">
            <el-button text class="user-dropdown-btn">
              <span>{{ userStore.nickname }}</span>
              <el-icon><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="$router.push('/user/' + userStore.userId)">个人主页</el-dropdown-item>
                <el-dropdown-item @click="$router.push('/profile/edit')">编辑资料</el-dropdown-item>
                <el-dropdown-item divided @click="toggleTheme">
                  <el-icon style="margin-right: 6px; vertical-align: -2px;">
                    <component :is="currentTheme === 'dark' ? Sunny : Moon" />
                  </el-icon>
                  {{ currentTheme === 'dark' ? '浅色模式' : '深色模式' }}
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <BreadcrumbBar :breadcrumbs="breadcrumbs" />

      <div ref="contentRef" class="content-area">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>

      <transition name="fade">
        <el-button
          v-show="showBackToTop"
          class="back-to-top"
          circle
          @click="scrollToTop"
        >
          <el-icon><ArrowUp /></el-icon>
        </el-button>
      </transition>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useUserStore } from "@/stores/user"
import { unreadCount as fetchUnreadCount } from "@/api/chat"
import BreadcrumbBar from "@/components/BreadcrumbBar.vue"
import { useTheme } from "@/composables/useTheme"
import {
  HomeFilled, MessageBox, Connection, Calendar, ShoppingCart,
  WarningFilled, ChatDotSquare, UserFilled, ChatLineSquare,
  Reading, Coin, School, Link, Bell, Edit, Search,
  ArrowRight, ArrowDown, ArrowUp, Notebook, InfoFilled, Key, Expand, Fold, Moon, Sunny
} from "@element-plus/icons-vue"

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const { theme: currentTheme, toggle: toggleTheme } = useTheme()
const unreadCount = ref(0)
const sidebarCollapsed = ref(localStorage.getItem('sidebar_collapsed') === 'true')
const sidebarOpen = ref(false)
const isMobile = ref(window.innerWidth < 768)
const showBackToTop = ref(false)
const contentRef = ref<HTMLElement | null>(null)
let pollTimer: ReturnType<typeof setInterval>

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('sidebar_collapsed', String(sidebarCollapsed.value))
}

function handleResize() {
  isMobile.value = window.innerWidth < 768
  if (!isMobile.value) sidebarOpen.value = false
}

function onContentScroll() {
  showBackToTop.value = contentRef.value ? contentRef.value.scrollTop > 300 : false
}

function scrollToTop() {
  contentRef.value?.scrollTo({ top: 0, behavior: 'smooth' })
}

const navGroups = computed(() => [
  {
    label: "发现",
    items: [
      { path: "/home", icon: HomeFilled, label: "首页" },
      { path: "/square", icon: MessageBox, label: "广场" },
      { path: "/search", icon: Search, label: "搜索" },
    ],
  },
  {
    label: "社交",
    items: [
      { path: "/circles", icon: Connection, label: "圈子" },
      { path: "/friends", icon: UserFilled, label: "好友" },
      { path: "/chat", icon: ChatLineSquare, label: "聊天" },
      { path: "/notifications", icon: Bell, label: "通知" },
    ],
  },
  {
    label: "校园",
    items: [
      { path: "/campus", icon: School, label: "校园服务", children: [
        { path: "/campus/calendar", label: "校历" },
        { path: "/campus/directory", label: "通讯录" },
        { path: "/campus/rooms", label: "空教室" },
      ]},
      { path: "/courses", icon: Notebook, label: "课程" },
      { path: "/verification", icon: InfoFilled, label: "认证" },
      { path: "/invite", icon: Key, label: "邀请码" },
    ],
  },
  {
    label: "生活",
    items: [
      { path: "/activities", icon: Calendar, label: "活动" },
      { path: "/products", icon: ShoppingCart, label: "二手" },
      { path: "/found", icon: WarningFilled, label: "失物招领" },
      { path: "/whispers", icon: ChatDotSquare, label: "树洞" },
      { path: "/points", icon: Coin, label: "积分" },
    ],
  },
])

const breadcrumbs = computed(() => {
  const path = route.path
  const segments = path.split("/").filter(Boolean)
  const crumbs: { label: string; path?: string }[] = []
  let cum = ""
  for (const seg of segments) {
    cum += "/" + seg
    const label = seg === "campus" ? "校园服务"
      : seg === "notifications" ? "通知"
      : seg === "profile" ? "个人中心"
      : seg.charAt(0).toUpperCase() + seg.slice(1)
    crumbs.push({ label, path: cum })
  }
  if (crumbs.length > 1) crumbs[crumbs.length - 1].path = undefined
  return crumbs
})

const activeRoute = computed(() => {
  const path = route.path
  if (path.startsWith("/post/") || path.startsWith("/post/create")) return "/square"
  if (path.startsWith("/user/") || path.startsWith("/profile/")) return "/home"
  if (path.startsWith("/campus/")) return "/campus/calendar"
  if (path.startsWith("/notifications")) return "/notifications"
  if (path.startsWith("/chat/")) return "/chat"
  for (const group of navGroups.value) {
    for (const item of group.items) {
      if (path.startsWith(item.path)) {
        if (item.children) {
          for (const child of item.children) {
            if (path.startsWith(child.path)) return child.path
          }
        }
        return item.path
      }
    }
  }
  return path
})

function handleLogout() {
  userStore.logout()
  router.push("/login")
}

  onMounted(() => {
    window.addEventListener('resize', handleResize)
    contentRef.value?.addEventListener('scroll', onContentScroll)
  })
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    contentRef.value?.removeEventListener('scroll', onContentScroll)
  })

onMounted(async () => {
  try {
    const res = await fetchUnreadCount()
    unreadCount.value = res.count
  } catch { /* interceptor handles 401 */ }
  pollTimer = setInterval(async () => {
    try {
      const res = await fetchUnreadCount()
      unreadCount.value = res.count
    } catch { /* */ }
  }, 30000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: $bg-app;
  position: relative;
}

// ═══ Sidebar Overlay (mobile) ═══
.sidebar-overlay {
  display: none;
  @media (max-width: 767px) {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.5);
    z-index: 9;
  }
}
.sidebar.mobile-open {
  @media (max-width: 767px) {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 11;
    width: $sidebar-width;
  }
}
.sidebar:not(.mobile-open) {
  @media (max-width: 767px) {
    transform: translateX(-100%);
    position: fixed;
    z-index: 11;
  }
}

// ═══ Sidebar ═══
.sidebar {
  width: $sidebar-width;
  background: $bg-sidebar;
  display: flex;
  flex-direction: column;
  border-right: 1px solid $border-default;
  flex-shrink: 0;
  z-index: 10;
  transition: width $duration-normal $ease-out, transform $duration-normal $ease-out;
  overflow: hidden;

  &.collapsed {
    width: 64px;
  }
}

.sidebar-collapse-btn {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 36px;
  margin: 0 $space-3 $space-2;
  cursor: pointer;
  color: $text-muted;
  border-radius: $radius-sm;
  transition: background $duration-fast;
  &:hover { background: rgba($brand-primary-hex, 0.1); color: $text-primary; }
  .el-icon { transition: transform $duration-normal; }
  .rotated { transform: rotate(180deg); }
}

.logo {
  display: flex;
  align-items: center;
  gap: $space-3;
  padding: $space-6 $space-5 $space-5;
  cursor: pointer;
  user-select: none;

  transition: padding $duration-normal;
  .logo-icon { flex-shrink: 0; }
  .logo-text {
    font-size: $text-lg;
    font-weight: 700;
    letter-spacing: 1px;
  }
  .logo-badge {
  .logo-text, .logo-badge { transition: opacity $duration-fast; }
  .collapsed & {
    justify-content: center;
    padding: $space-4 $space-2;
    .logo-text, .logo-badge { opacity: 0; }
  }
    font-size: 10px;
    color: $brand-primary;
    background: rgba($brand-primary-hex, 0.12);
    padding: 1px $space-2;
    border-radius: $radius-full;
    margin-left: auto;
  }
}

.nav-groups {
  flex: 1;
  padding: 0 $space-2;

.collapsed & {
    :deep(.el-menu-item) {
      justify-content: center;
      padding: 0 !important;
      .el-icon { margin-right: 0 !important; }
    }
    :deep(.el-sub-menu__title) {
      justify-content: center;
      padding: 0 6px !important;
      .el-icon { margin-right: 0 !important; }
    }
    :deep(.el-sub-menu) { display: none; }
  }

.nav-group {
  margin-bottom: $space-2;
}

.nav-group-label {
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  color: $text-muted;
  letter-spacing: 1px;
  padding: $space-3 $space-3 $space-2;
}

  background: transparent;
  border-right: none;

  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    color: $text-secondary;
    border-radius: $radius-sm;
    margin: 1px 0;
    height: 38px;
    line-height: 38px;
    font-size: $text-sm;

    &:hover {
      background: rgba($brand-primary-hex, 0.08);
      color: $text-primary;
    }
    &.is-active {
      background: rgba($brand-primary-hex, 0.15);
      color: $brand-primary;
      font-weight: 600;
    }
    .el-icon { margin-right: $space-2; }
  }

  :deep(.el-sub-menu__title) { padding-right: 6px; }
  :deep(.el-sub-menu .el-menu) { background: transparent; }
  :deep(.el-sub-menu .el-menu-item) { padding-left: 50px !important; height: 34px; line-height: 34px; font-size: 13px; }
}

// ═══ Footer ═══
.sidebar-footer {
  padding: $space-3;
  border-top: 1px solid $border-default;
  transition: opacity $duration-fast;
.collapsed & { opacity: 0; pointer-events: none; }
}

.user-card {
  display: flex;
  align-items: center;
  gap: $space-3;
  padding: $space-3;
  border-radius: $radius-sm;
  cursor: pointer;
  transition: background $duration-fast;

  &:hover { background: rgba($brand-primary-hex, 0.08); }

  .user-info {
    flex: 1;
    min-width: 0;
    .nickname {
      font-size: 13px;
      font-weight: 600;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .school {
      font-size: 11px;
      color: $text-muted;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
  .user-arrow { color: $text-muted; font-size: 14px; }
.collapsed & { display: none; }
}

// ═══ Main ═══
.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
  background: $bg-app;
}

// ═══ Top Bar ═══
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: $topbar-height;
  padding: 0 $space-6;
  border-bottom: 1px solid $border-default;
  flex-shrink: 0;
  background: $bg-app;

.top-bar-left {
    display: flex;
    align-items: center;
    gap: $space-3;
  }

.hamburger-btn, .collapse-btn {
    color: $text-muted;
    &:hover { color: $text-primary; }
  }

.search-btn {
    color: $text-muted;
    display: flex;
    align-items: center;
    gap: $space-2;
    font-size: $text-sm;
    background: $bg-card;
    border-radius: $radius-full;
    padding: $space-1 $space-4;
    border: 1px solid $border-default;
    min-width: 200px;
    justify-content: flex-start;
    @media (max-width: 767px) {
      min-width: 120px;
      span { display: none; }
    }
:hover {
      border-color: rgba($brand-primary-hex, 0.4);
      color: $text-secondary;
    }
  }

.top-bar-right {
    display: flex;
    align-items: center;
    gap: $space-2;
  }

.user-dropdown-btn {
    color: $text-secondary;
    display: flex;
    align-items: center;
    gap: $space-1;
    font-size: $text-sm;
  }
}

.notif-badge :deep(.el-badge__content) {
  background: $brand-primary;
  border: none;
  font-size: 11px;
  height: 18px;
  line-height: 18px;
}

// ═══ Content Area ═══
.content-area {
  flex: 1;
  overflow-y: auto;
  padding: $space-6;
  background: $bg-app;
  @media (max-width: 767px) {
    padding: $space-4;
  }
}
</style>

.back-to-top {
  position: fixed;
  bottom: 32px;
  right: 32px;
  z-index: 100;
  --el-button-size: 48px;
  font-size: 20px;
  background: linear-gradient(135deg, $brand-primary, #a86555) !important;
  border: none !important;
  box-shadow: 0 4px 16px rgba($brand-primary-hex, 0.3) !important;
  color: #fff !important;
  transition: all $duration-normal $ease-out;

  &:hover {
    transform: scale(1.08);
    box-shadow: 0 6px 24px rgba($brand-primary-hex, 0.45) !important;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
