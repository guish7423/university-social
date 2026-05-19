<script setup lang="ts">
import { ref } from "vue"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const title = ref("校园社")
</script>

<template>
  <view class="container">
    <view class="header">
      <text class="logo">{{ title }}</text>
      <text v-if="userStore.isLogin" class="greeting">欢迎, {{ userStore.nickname }}</text>
    </view>
    <view class="content">
      <text class="placeholder">动态广场 — 即将上线</text>
    </view>
  </view>
</template>

<style scoped>
.container {
  min-height: 100vh;
}
.header {
  padding: 40rpx 30rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.logo {
  font-size: 48rpx;
  font-weight: 700;
  color: #fff;
}
.greeting {
  display: block;
  margin-top: 10rpx;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.9);
}
.content {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 200rpx 0;
}
.placeholder {
  font-size: 32rpx;
  color: #999;
}
</style>
