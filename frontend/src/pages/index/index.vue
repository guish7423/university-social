<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useUserStore } from "@/store/user"
import { devLogin as devLoginApi } from "@/api/user"

const userStore = useUserStore()
const devNickname = ref("")
const showDevLogin = ref(false)

onMounted(() => {
  showDevLogin.value = !userStore.isLogin
})

async function devLogin() {
  const nick = devNickname.value.trim() || "测试用户"
  try {
    const res = await devLoginApi(nick)
    userStore.setToken(res.token)
    userStore.setUserInfo(res.user)
    showDevLogin.value = false
    uni.switchTab({ url: "/pages/square/index" })
  } catch {
    uni.showToast({ title: "登录失败", icon: "none" })
  }
}

function goSquare() { uni.switchTab({ url: '/pages/square/index' }) }
function goCircles() { uni.switchTab({ url: '/pages/circle/list' }) }
function goFriends() { uni.navigateTo({ url: '/pages/friend/index' }) }
function goNotifications() { uni.navigateTo({ url: '/pages/notification/index' }) }
</script>

<template>
  <view class="container">
    <view class="hero">
      <view class="hero-bg" />
      <view class="hero-content">
        <text class="logo">校园社</text>
        <text class="tagline">你的大学生活，从这里开始</text>
      </view>
    </view>

    <view v-if="showDevLogin" class="login-card">
      <view class="login-header">
        <text class="login-title">🎓 体验校园社</text>
        <text class="login-desc">输入昵称即可开始，无需注册</text>
      </view>
      <view class="login-form">
        <u-input
          v-model="devNickname"
          placeholder="输入你的昵称"
          :border="true"
          shape="square"
          :customStyle="{ padding: '20rpx 24rpx', fontSize: '28rpx', height: '88rpx' }"
        />
        <u-button
          type="primary"
          shape="circle"
          size="large"
          :customStyle="{ marginTop: '28rpx', height: '88rpx', fontSize: '32rpx', fontWeight: '600' }"
          @click="devLogin"
        >
          进入校园社
        </u-button>
      </view>
    </view>

    <view v-else class="welcome-card">
      <view class="welcome-user">
        <u-avatar
          :src="userStore.avatar"
          size="80"
          shape="circle"
        />
        <view class="welcome-info">
          <text class="welcome-name">{{ userStore.nickname || '用户' }}</text>
          <text class="welcome-text">欢迎回来！</text>
        </view>
      </view>
    </view>

    <view class="feature-grid">
      <view class="feature-item" @click="goSquare">
        <view class="feature-icon-wrap" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
          <u-icon name="home" size="40" color="#fff" />
        </view>
        <text class="feature-name">动态广场</text>
        <text class="feature-desc">发现校园新鲜事</text>
      </view>
      <view class="feature-item" @click="goCircles">
        <view class="feature-icon-wrap" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
          <u-icon name="grid" size="40" color="#fff" />
        </view>
        <text class="feature-name">校园圈子</text>
        <text class="feature-desc">找到你的同好</text>
      </view>
      <view class="feature-item" @click="goFriends">
        <view class="feature-icon-wrap" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
          <u-icon name="account" size="40" color="#fff" />
        </view>
        <text class="feature-name">好友</text>
        <text class="feature-desc">结识新朋友</text>
      </view>
      <view class="feature-item" @click="goNotifications">
        <view class="feature-icon-wrap" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
          <u-icon name="bell" size="40" color="#fff" />
        </view>
        <text class="feature-name">通知</text>
        <text class="feature-desc">查看消息提醒</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.container {
  background: var(--u-bg-color, #f3f4f6);
  min-height: 100vh;
}

.hero {
  position: relative;
  padding: 100rpx 40rpx 80rpx;
  overflow: hidden;
}

.hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  border-radius: 0 0 48rpx 48rpx;
}

.hero-content {
  position: relative;
  z-index: 1;
  text-align: center;
}

.logo {
  font-size: 64rpx;
  font-weight: 800;
  color: #fff;
  display: block;
  letter-spacing: 4rpx;
}

.tagline {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.7);
  margin-top: 16rpx;
  display: block;
}

.login-card {
  margin: -40rpx 30rpx 0;
  background: #fff;
  border-radius: 24rpx;
  padding: 40rpx;
  position: relative;
  z-index: 2;
  box-shadow: 0 8rpx 40rpx rgba(0, 0, 0, 0.06);
}

.login-header {
  margin-bottom: 32rpx;
  text-align: center;
}

.login-title {
  font-size: 36rpx;
  font-weight: 700;
  color: #1a1a2e;
  display: block;
  margin-bottom: 10rpx;
}

.login-desc {
  font-size: 26rpx;
  color: #909399;
}

.welcome-card {
  margin: -40rpx 30rpx 0;
  background: #fff;
  border-radius: 24rpx;
  padding: 36rpx;
  position: relative;
  z-index: 2;
  box-shadow: 0 8rpx 40rpx rgba(0, 0, 0, 0.06);
}

.welcome-user {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.welcome-info {
  flex: 1;
}

.welcome-name {
  font-size: 34rpx;
  font-weight: 700;
  display: block;
}

.welcome-text {
  font-size: 26rpx;
  color: #909399;
  margin-top: 6rpx;
  display: block;
}

.feature-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20rpx;
  padding: 30rpx;
}

.feature-item {
  background: #fff;
  border-radius: 20rpx;
  padding: 36rpx 24rpx;
  text-align: center;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.04);
}

.feature-icon-wrap {
  width: 88rpx;
  height: 88rpx;
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16rpx;
}

.feature-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #303133;
  display: block;
}

.feature-desc {
  font-size: 22rpx;
  color: #909399;
  margin-top: 6rpx;
  display: block;
}
</style>
