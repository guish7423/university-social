<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useUserStore } from "@/store/user"
import { getProfile } from "@/api/user"
import { listPosts, type PostData } from "@/api/post"

const userStore = useUserStore()
const posts = ref<PostData[]>([])
const loading = ref(true)

const menuGroups = [
  { title: "内容管理", items: [
    { icon: "◈", label: "我的帖子", path: "/pages/user/posts" },
    { icon: "▤", label: "我的发布", path: "/pages/product/mine" },
    { icon: "▣", label: "我的活动", path: "/pages/activity/mine" },
  ]},
  { title: "社交", items: [
    { icon: "♢", label: "好友", path: "/pages/friend/index" },
    { icon: "◉", label: "通知", path: "/pages/notification/index" },
    { icon: "✉", label: "消息", path: "/pages/chat/index" },
    { icon: "♢", label: "关注", path: "/pages/user/following" },
    { icon: "◯", label: "粉丝", path: "/pages/user/followers" },
  ]},
  { title: "其他", items: [
    { icon: "⌕", label: "校园认证", path: "/pages/verification/index" },
    { icon: "★", label: "积分排行", path: "/pages/points/index" },
  ]},
]

onMounted(async () => {
  if (userStore.isLogin) {
    try {
      const info = await getProfile()
      userStore.setUserInfo(info)
    } catch {}
  }
  loading.value = false
})

function go(path: string) {
  uni.navigateTo({ url: path })
}

function goLogin() {
  uni.navigateTo({ url: "/pages/login/index" })
}

function goEdit() {
  uni.navigateTo({ url: "/pages/user/edit" })
}
</script>

<template>
  <view class="page">
    <view class="profile-header">
      <view v-if="userStore.isLogin" class="profile-info" @click="goEdit">
        <u-avatar :src="userStore.avatar" :text="userStore.nickname?.[0] || '?'" size="120" shape="circle" />
        <view class="profile-text">
          <text class="profile-name">{{ userStore.nickname || '未设置昵称' }}</text>
          <text class="profile-school">{{ userStore.school || '校园社' }}</text>
        </view>
      </view>
      <view v-else class="profile-info" @click="goLogin">
        <u-avatar text="?" size="120" shape="circle" />
        <view class="profile-text">
          <text class="profile-name">未登录</text>
          <text class="profile-school">点击登录</text>
        </view>
      </view>
    </view>

    <view v-for="(group, gi) in menuGroups" :key="gi" class="menu-section">
      <text class="menu-section-title">{{ group.title }}</text>
      <view class="menu-list">
        <view v-for="(item, ii) in group.items" :key="ii"
          class="menu-item" @click="go(item.path)">
          <view class="menu-item-left">
            <view class="menu-icon-wrap"><text class="menu-icon">{{ item.icon }}</text></view>
            <text class="menu-label">{{ item.label }}</text>
          </view>
          <text class="menu-arrow">›</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding-bottom: 300rpx; }

.profile-header {
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 100%);
  padding: 80rpx 32rpx 48rpx;
}
.profile-info { display: flex; align-items: center; gap: 24rpx; }
.profile-text { display: flex; flex-direction: column; gap: 8rpx; }
.profile-name { font-size: 32rpx; font-weight: 600; color: #fff; }
.profile-school { font-size: 24rpx; color: rgba(255,255,255,0.5); font-weight: 300; }

.menu-section { margin: 24rpx 16rpx 0; }
.menu-section-title { font-size: 24rpx; color: var(--ink-subtle, #8E9BAB); padding: 0 16rpx 16rpx; display: block; font-weight: 500; }
.menu-list {
  background: var(--color-surface, #fff);
  border-radius: 14rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
}
.menu-item {
  display: flex; align-items: center; justify-content: space-between;
  padding: 28rpx 24rpx;
  border-bottom: 1rpx solid var(--hairline-light, #EAE6E0);
  cursor: pointer;
  transition: background 0.15s;
  &:last-child { border-bottom: none; }
  &:active { background: var(--color-surface-1, #F0EDE8); }
}
.menu-item-left { display: flex; align-items: center; gap: 16rpx; }
.menu-icon-wrap {
  width: 52rpx; height: 52rpx;
  border-radius: 12rpx;
  background: rgba(198,122,106,0.08);
  display: flex; align-items: center; justify-content: center;
}
.menu-icon { font-size: 22rpx; color: #C67A6A; }
.menu-label { font-size: 28rpx; color: var(--ink, #1E2A3A); font-weight: 400; }
.menu-arrow { font-size: 32rpx; color: var(--ink-tertiary, #B8C2CE); font-weight: 300; }
</style>
