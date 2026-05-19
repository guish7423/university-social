<script setup lang="ts">
import { ref } from "vue"
import { useUserStore } from "@/store/user"
import { wxLogin } from "@/api/user"

const userStore = useUserStore()
const loading = ref(false)

async function handleLogin() {
  loading.value = true
  try {
    const { code } = await uni.login({ provider: "weixin" })
    const { userInfo } = await uni.getUserProfile({ desc: "用于完善用户资料" })
    const res = await wxLogin(code, userInfo.nickName, userInfo.avatarUrl)
    userStore.setToken(res.token)
    userStore.setUserInfo(res.user)
    uni.switchTab({ url: "/pages/index/index" })
  } catch (err: any) {
    if (err.errMsg !== "getUserProfile:fail auth deny") {
      uni.showToast({ title: "登录失败", icon: "none" })
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <view class="container">
    <view class="login-box">
      <text class="title">校园社</text>
      <text class="subtitle">大学社交平台</text>
      <button
        class="wechat-btn"
        type="primary"
:loading="loading"
:disabled="loading"
        @click="handleLogin"
      >
        微信一键登录
      </button>
    </view>
  </view>
</template>

<style scoped>
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-box {
  text-align: center;
  padding: 60rpx;
}
.title {
  font-size: 64rpx;
  font-weight: 700;
  color: #fff;
}
.subtitle {
  display: block;
  margin: 20rpx 0 80rpx;
  font-size: 32rpx;
  color: rgba(255, 255, 255, 0.8);
}
.wechat-btn {
  width: 500rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 44rpx;
  font-size: 32rpx;
}
</style>
