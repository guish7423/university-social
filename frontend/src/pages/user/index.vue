<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useUserStore } from "@/store/user"
import { getProfile } from "@/api/user"

const userStore = useUserStore()

const user = ref<any>(null)
const stats = ref({ postCount: 0, friendCount: 0, pointsBalance: 0 })

const isLogin = computed(() => userStore.isLogin)

onMounted(async () => {
  if (userStore.isLogin) {
    try {
      const res = await getProfile()
      user.value = res
      userStore.setUserInfo(res)
    } catch {}
  }
})

function goLogin() { uni.navigateTo({ url: "/pages/login/index" }) }
function goEdit() { uni.navigateTo({ url: "/pages/user/edit" }) }
function goPost() { uni.navigateTo({ url: "/pages/user/posts" }) }
function goFollower() { uni.navigateTo({ url: "/pages/user/followers?id=" + userStore.id }) }
function goFollowing() { uni.navigateTo({ url: "/pages/user/following?id=" + userStore.id }) }
function goPoints() { uni.navigateTo({ url: "/pages/points/index" }) }
function goVerification() { uni.navigateTo({ url: "/pages/verification/index" }) }

const menus = [
  [
    { icon: "📝", label: "我的帖子", action: goPost },
    { icon: "🛍️", label: "我的发布", action: () => uni.navigateTo({ url: "/pages/product/mine" }) },
    { icon: "🎯", label: "我的活动", action: () => uni.navigateTo({ url: "/pages/activity/mine" }) },
  ],
  [
    { icon: "🔎", label: "失物招领", action: () => uni.navigateTo({ url: "/pages/found/mine" }) },
    { icon: "⭐", label: "积分", action: goPoints },
    { icon: "✅", label: "校园认证", action: goVerification },
  ],
  [
    { icon: "📋", label: "编辑资料", action: goEdit },
    { icon: "🔒", label: "退出登录", action: () => { userStore.logout(); uni.reLaunch({ url: "/pages/login/index" }) } },
  ],
]
</script>

<template>
  <view class="page-root">
    <view class="page-glow" />
    <view v-if="!isLogin" class="login-prompt animate-fade-in-up">
      <text class="login-emoji">🌸</text>
      <text class="login-title">登录梦园社</text>
      <text class="login-desc">登录后即可体验全部功能</text>
      <view class="login-btn" @click="goLogin">立即登录</view>
    </view>

    <template v-else>
      <view class="profile-header">
        <view class="header-bg" />
        <view class="header-content">
          <view class="avatar-wrap">
            <image class="avatar" :src="userStore.avatar || ''" mode="aspectFill" />
            <view class="avatar-badge" @click="goEdit">✏️</view>
          </view>
          <text class="nickname">{{ userStore.nickname || "用户" }}</text>
          <text class="school">{{ userStore.school || "未设置学校" }}</text>
          <view class="verified-badge" v-if="userStore.isVerified">
            <text class="verified-icon">✅</text>
            <text class="verified-text">已认证</text>
          </view>
        </view>
      </view>

      <view class="stats-row glass-card animate-fade-in-up">
        <view class="stat-item" @click="goPost">
          <text class="stat-value">--</text>
          <text class="stat-label">帖子</text>
        </view>
        <view class="stat-divider" />
        <view class="stat-item" @click="goFollowing">
          <text class="stat-value">--</text>
          <text class="stat-label">关注</text>
        </view>
        <view class="stat-divider" />
        <view class="stat-item" @click="goFollower">
          <text class="stat-value">--</text>
          <text class="stat-label">粉丝</text>
        </view>
      </view>

      <view class="menu-section animate-fade-in-up stagger-2">
        <view v-for="(group, gi) in menus" :key="gi" class="menu-group glass-card">
          <view v-for="(item, ii) in group" :key="ii" class="menu-item" @click="item.action">
            <text class="menu-icon">{{ item.icon }}</text>
            <text class="menu-label">{{ item.label }}</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>
      </view>
    </template>
  </view>
</template>

<style lang="scss" scoped>
.page-root { min-height: 100vh; }
.page-glow {
  position: fixed; top: -200rpx; right: -200rpx;
  width: 600rpx; height: 600rpx;
  background: radial-gradient(circle, rgba(192,132,252,0.1) 0%, transparent 60%);
  pointer-events: none; z-index: 0;
}

.glass-card {
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  box-shadow: var(--glass-shadow);
}

.login-prompt {
  display: flex; flex-direction: column; align-items: center;
  padding: 300rpx 40rpx; gap: 16rpx; position: relative; z-index: 1;
}
.login-emoji { font-size: 100rpx; }
.login-title { font-size: 40rpx; font-weight: 700; color: var(--ink); font-family: 'ZCOOL KuaiLe', 'PingFang SC', sans-serif; }
.login-desc { font-size: 28rpx; color: var(--ink-subtle); }
.login-btn {
  margin-top: 24rpx; padding: 22rpx 80rpx; border-radius: 50rpx;
  background: var(--brand-gradient); color: #fff; font-size: 30rpx; font-weight: 600;
  transition: all 0.3s ease;
  &:active { transform: scale(0.95); }
}

.profile-header {
  position: relative; padding: 120rpx 40rpx 60rpx; z-index: 1;
  text-align: center;
}
.header-bg {
  position: absolute; top: 0; left: 0; right: 0; height: 280rpx;
  background: linear-gradient(180deg, rgba(192,132,252,0.12) 0%, transparent 100%);
}
.header-content { position: relative; z-index: 1; display: flex; flex-direction: column; align-items: center; gap: 12rpx; }
.avatar-wrap {
  position: relative; width: 140rpx; height: 140rpx;
  border-radius: 50%; padding: 4rpx; background: var(--brand-gradient);
  margin-bottom: 8rpx;
}
.avatar { width: 100%; height: 100%; border-radius: 50%; }
.avatar-badge {
  position: absolute; bottom: 0; right: 0;
  width: 44rpx; height: 44rpx; border-radius: 50%;
  background: #fff; display: flex; align-items: center; justify-content: center;
  font-size: 24rpx; box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
.nickname { font-size: 36rpx; font-weight: 700; color: var(--ink); }
.school { font-size: 26rpx; color: var(--ink-muted); }
.verified-badge { display: flex; align-items: center; gap: 6rpx; margin-top: 4rpx; }
.verified-icon { font-size: 24rpx; }
.verified-text { font-size: 24rpx; color: var(--color-success); font-weight: 500; }

.stats-row {
  display: flex; margin: 0 32rpx; padding: 24rpx 0;
  border-radius: 20rpx; position: relative; z-index: 1;
}
.stat-item { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 6rpx; cursor: pointer; }
.stat-value { font-size: 34rpx; font-weight: 700; color: var(--ink); }
.stat-label { font-size: 24rpx; color: var(--ink-subtle); }
.stat-divider { width: 1rpx; height: 48rpx; background: var(--hairline); align-self: center; }

.menu-section { padding: 32rpx; display: flex; flex-direction: column; gap: 20rpx; position: relative; z-index: 1; }
.menu-group { border-radius: 20rpx; overflow: hidden; }
.menu-item {
  display: flex; align-items: center; gap: 20rpx;
  padding: 28rpx 24rpx; cursor: pointer;
  transition: background 0.2s;
  &:not(:last-child) { border-bottom: 1rpx solid var(--hairline); }
  &:active { background: var(--color-surface-1); }
}
.menu-icon { font-size: 32rpx; width: 44rpx; text-align: center; }
.menu-label { flex: 1; font-size: 28rpx; color: var(--ink); }
.menu-arrow { font-size: 36rpx; color: var(--ink-tertiary); }
</style>
