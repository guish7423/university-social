<script setup lang="ts">
import { ref } from "vue"
import { useUserStore } from "@/store/user"
import { devLogin } from "@/api/user"

const userStore = useUserStore()
const nickname = ref("")
const loading = ref(false)

async function handleLogin() {
  if (!nickname.value.trim()) { uni.showToast({ title: "请输入昵称", icon: "none" }); return }
  loading.value = true
  try {
    const res = await devLogin(nickname.value.trim())
    userStore.setToken(res.token)
    userStore.setUserInfo(res.user)
    uni.showToast({ title: "登录成功", icon: "success" })
    setTimeout(() => uni.switchTab({ url: "/pages/square/index" }), 800)
  } catch { uni.showToast({ title: "登录失败", icon: "none" }) }
  loading.value = false
}
</script>

<template>
  <view class="page-root">
    <view class="glow-1" />
    <view class="glow-2" />
    <view class="glow-3" />

    <view class="login-card glass-card animate-fade-in-up">
      <view class="card-header">
        <text class="card-emoji">🌸</text>
        <text class="card-title">梦园社</text>
        <text class="card-desc">发现校园美好，连接有趣灵魂</text>
      </view>

      <view class="card-body">
        <view class="input-wrap">
          <text class="input-label">你的昵称</text>
          <input
            class="input-field"
            v-model="nickname"
            placeholder="输入昵称开始探索"
            placeholder-class="input-placeholder"
            maxlength="20"
          />
        </view>
        <view class="login-btn" :class="{ loading }" @click="handleLogin">
          <text>{{ loading ? '登录中...' : '✨ 进入梦园社' }}</text>
        </view>
        <text class="login-tip">首次登录自动注册</text>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page-root {
  min-height: 100vh; display: flex; align-items: center; justify-content: center;
  padding: 40rpx; position: relative; overflow: hidden;
}
.glow-1 {
  position: fixed; top: -300rpx; left: -200rpx;
  width: 700rpx; height: 700rpx;
  background: radial-gradient(circle, rgba(192,132,252,0.12) 0%, transparent 60%);
  pointer-events: none;
}
.glow-2 {
  position: fixed; bottom: -200rpx; right: -200rpx;
  width: 600rpx; height: 600rpx;
  background: radial-gradient(circle, rgba(249,168,212,0.1) 0%, transparent 60%);
  pointer-events: none;
}
.glow-3 {
  position: fixed; top: 40%; left: 30%;
  width: 400rpx; height: 400rpx;
  background: radial-gradient(circle, rgba(103,232,249,0.06) 0%, transparent 60%);
  pointer-events: none;
}

.glass-card {
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  box-shadow: var(--glass-shadow);
}

.login-card {
  width: 100%; max-width: 600rpx;
  border-radius: 32rpx; padding: 60rpx 40rpx;
  position: relative; z-index: 1;
}

.card-header { display: flex; flex-direction: column; align-items: center; gap: 12rpx; margin-bottom: 48rpx; }
.card-emoji { font-size: 80rpx; }
.card-title { font-size: 44rpx; font-weight: 700; background: var(--brand-gradient); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; font-family: 'ZCOOL KuaiLe', 'PingFang SC', sans-serif; }
.card-desc { font-size: 26rpx; color: var(--ink-muted); }

.input-wrap { margin-bottom: 32rpx; }
.input-label { font-size: 26rpx; color: var(--ink); font-weight: 500; margin-bottom: 12rpx; display: block; }
.input-field {
  width: 100%; padding: 24rpx 28rpx; border-radius: 16rpx;
  background: var(--color-surface-1); border: 2rpx solid transparent;
  font-size: 28rpx; color: var(--ink); transition: all 0.3s ease;
  &:focus { border-color: var(--brand-primary); background: #fff; }
}
.input-placeholder { color: var(--ink-tertiary); }

.login-btn {
  width: 100%; padding: 26rpx 0; border-radius: 50rpx;
  background: var(--brand-gradient); color: #fff;
  font-size: 30rpx; font-weight: 600; text-align: center;
  transition: all 0.3s ease;
  &:active { transform: scale(0.97); }
  &.loading { opacity: 0.7; }
}

.login-tip { display: block; text-align: center; margin-top: 20rpx; font-size: 22rpx; color: var(--ink-tertiary); }
</style>
