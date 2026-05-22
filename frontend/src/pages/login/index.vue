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
    <view class="hero-area">
      <view class="hero-glow" />
      <view class="brand">
        <view class="logo-wrap">
          <text class="logo-text">C</text>
        </view>
        <text class="brand-title">校园社</text>
        <text class="brand-subtitle">连接校园里的每一个人</text>
      </view>
    </view>

    <view class="form-area">
      <view class="form-card">
        <text class="form-label">你的昵称</text>
        <view class="input-wrap">
          <u-icon name="account" size="32" color="#B8C2CE" />
          <input v-model="nickname" class="input" placeholder="输入一个昵称开始" placeholder-class="ph" maxlength="20" />
        </view>

        <view class="login-btn" :class="{ loading: loading }" @click="handleLogin">
          <u-loading v-if="loading" mode="flower" size="32" color="#fff" />
          <text v-else>进入校园社</text>
        </view>

        <text class="form-hint">登录即表示同意 用户协议 和 隐私政策</text>
      </view>

      <view class="features-preview">
        <view class="fp-item" v-for="f in [{icon:'moment',label:'动态广场'},{icon:'community',label:'校园圈子'},{icon:'favor',label:'树洞'},{icon:'gift',label:'二手市场'}]" :key="f.label">
          <u-icon :name="f.icon" size="32" color="#C67A6A" />
          <text class="fp-label">{{ f.label }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style lang="scss">
page { background: #1E2A3A; }
.ph { color: #B8C2CE; font-size: 28rpx; }
</style>

<style lang="scss" scoped>
.page { min-height: 100vh; display: flex; flex-direction: column; }

.hero-area {
  position: relative; flex: 1; display: flex; align-items: center; justify-content: center;
  padding: 80rpx 40rpx; overflow: hidden;
}
.hero-glow {
  position: absolute;
  width: 500rpx; height: 500rpx;
  background: radial-gradient(circle, rgba(198,122,106,0.15) 0%, transparent 60%);
  top: 50%; left: 50%; transform: translate(-50%, -50%);
}
.brand { display: flex; flex-direction: column; align-items: center; gap: 16rpx; z-index: 1; }
.logo-wrap {
  width: 120rpx; height: 120rpx; border-radius: 32rpx;
  background: linear-gradient(135deg, #C67A6A, #A85A4A);
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 8rpx 32rpx rgba(198,122,106,0.3);
  margin-bottom: 8rpx;
}
.logo-text { font-size: 56rpx; font-weight: 700; color: #fff; }
.brand-title { font-size: 40rpx; font-weight: 700; color: #fff; letter-spacing: 4rpx; }
.brand-subtitle { font-size: 26rpx; color: rgba(255,255,255,0.4); font-weight: 300; }

.form-area {
  background: #F7F4F0;
  border-radius: 40rpx 40rpx 0 0;
  padding: 48rpx 32rpx 60rpx;
  display: flex; flex-direction: column; gap: 40rpx;
}
.form-card { background: #fff; border-radius: 24rpx; padding: 40rpx 32rpx; box-shadow: 0 4rpx 24rpx rgba(30,42,58,0.06); }
.form-label { font-size: 28rpx; font-weight: 600; color: #1E2A3A; display: block; margin-bottom: 16rpx; }
.input-wrap {
  display: flex; align-items: center; gap: 12rpx;
  background: #F7F4F0; border-radius: 14rpx; padding: 20rpx 24rpx;
}
.input { flex: 1; font-size: 28rpx; color: #1E2A3A; background: transparent; border: none; outline: none; height: 44rpx; }

.login-btn {
  margin-top: 32rpx; height: 88rpx; border-radius: 44rpx;
  background: linear-gradient(135deg, #C67A6A, #A85A4A);
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; transition: all 0.2s ease;
  &:active { transform: scale(0.97); opacity: 0.9; }
  text { font-size: 30rpx; font-weight: 600; color: #fff; }
}
.login-btn.loading { opacity: 0.7; }

.form-hint { font-size: 22rpx; color: #B8C2CE; text-align: center; margin-top: 24rpx; display: block; }

.features-preview { display: flex; justify-content: space-around; }
.fp-item { display: flex; flex-direction: column; align-items: center; gap: 8rpx; }
.fp-label { font-size: 22rpx; color: #8E9BAB; }
</style>
