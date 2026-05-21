<script setup lang="ts">
import { ref } from "vue"
import { devLogin } from "@/api/user"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const nickname = ref("")
const loading = ref(false)

async function handleLogin() {
  if (!nickname.value.trim()) {
    uni.showToast({ title: "请输入昵称", icon: "none" })
    return
  }
  loading.value = true
  try {
    const res = await devLogin(nickname.value.trim())
    userStore.setToken(res.token)
    userStore.setUserInfo(res.user)
    uni.showToast({ title: "登录成功", icon: "success" })
    setTimeout(() => uni.switchTab({ url: "/pages/square/index" }), 500)
  } catch { uni.showToast({ title: "登录失败", icon: "none" }) }
  loading.value = false
}
</script>

<template>
  <view class="page">
    <view class="login-card">
      <view class="brand-section">
        <text class="brand-title">校园社</text>
        <text class="brand-subtitle">连接校园里的每一个人</text>
      </view>

      <view class="form-section">
        <text class="form-label">昵称</text>
        <view class="input-wrap">
          <input v-model="nickname" class="input" placeholder="输入你的昵称" placeholder-class="placeholder" maxlength="20" />
        </view>

        <view class="login-btn" :class="{ loading: loading }" @click="handleLogin">
          <text>{{ loading ? '登录中...' : '进入校园社' }}</text>
        </view>

        <text class="form-hint">登录即表示同意 用户协议 和 隐私政策</text>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page {
  min-height: 100vh; display: flex; align-items: center; justify-content: center;
  background: var(--color-canvas, #F7F4F0);
  padding: 40rpx;
}

.login-card {
  width: 100%;
  max-width: 560rpx;
  background: var(--color-surface, #fff);
  border-radius: 24rpx;
  box-shadow: 0 8rpx 32rpx rgba(30,42,58,0.06);
  overflow: hidden;
}

.brand-section {
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 100%);
  padding: 60rpx 40rpx;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}
.brand-title { font-size: 40rpx; font-weight: 600; color: #C67A6A; letter-spacing: 4rpx; }
.brand-subtitle { font-size: 24rpx; color: rgba(255,255,255,0.45); font-weight: 300; letter-spacing: 2rpx; }

.form-section {
  padding: 40rpx;
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}
.form-label { font-size: 26rpx; color: var(--ink-muted, #5C6B7E); font-weight: 500; }

.input-wrap {
  border: 1.5rpx solid var(--hairline, #E0DBD4);
  border-radius: 14rpx;
  padding: 24rpx 20rpx;
  transition: border-color 0.2s;
}
.input-wrap:focus-within {
  border-color: #C67A6A;
}
.input { font-size: 28rpx; color: var(--ink, #1E2A3A); width: 100%; }
.placeholder { color: var(--ink-tertiary, #B8C2CE); }

.login-btn {
  padding: 28rpx;
  background: #C67A6A;
  border-radius: 14rpx;
  text-align: center;
  font-size: 28rpx;
  font-weight: 500;
  color: #fff;
  transition: all 0.2s ease;
  cursor: pointer;
  &:active { opacity: 0.85; transform: scale(0.98); }
  &.loading { opacity: 0.7; }
}

.form-hint { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); text-align: center; }
</style>
