<script setup lang="ts">
import { ref, onMounted } from "vue"
import { searchUsers, sendFriendRequest, type UserInfo } from "@/api/social"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const query = ref("")
const results = ref<UserInfo[]>([])
const loading = ref(false)

async function handleSearch() {
  if (!query.value.trim()) return
  loading.value = true
  try {
    results.value = await searchUsers(query.value)
  } catch {}
  loading.value = false
}

async function addFriend(id: number) {
  await sendFriendRequest(id)
  uni.showToast({ title: "请求已发送", icon: "success" })
}
</script>

<template>
  <view class="container">
    <view class="search-bar">
      <input v-model="query" class="input" placeholder="搜索用户昵称" @confirm="handleSearch" />
      <button class="search-btn" @click="handleSearch">搜索</button>
    </view>

    <view v-if="loading" class="state">搜索中...</view>
    <view v-else-if="results.length === 0 && query" class="state">未找到用户</view>
    <view v-else-if="results.length === 0" class="state">输入昵称搜索</view>

    <view v-for="u in results" :key="u.id" class="user-item">
      <image v-if="u.avatar" class="avatar" :src="u.avatar" mode="aspectFill" />
      <view class="avatar-placeholder">{{ u.nickname?.[0] || "?" }}</view>
      <view class="user-info">
        <text class="name">{{ u.nickname }}</text>
        <text v-if="u.school" class="school">{{ u.school }}</text>
      </view>
      <button v-if="u.id !== userStore.id" class="add-btn" @click="addFriend(u.id)">加好友</button>
    </view>
  </view>
</template>

<style scoped>
.container { min-height: 100vh; background: #f5f5f5; }
.search-bar { display: flex; gap: 16rpx; padding: 30rpx; background: #fff; }
.input { flex: 1; height: 72rpx; background: #f0f0f0; border-radius: 36rpx; padding: 0 24rpx; font-size: 28rpx; }
.search-btn { font-size: 28rpx; color: #667eea; background: none; border: none; }
.state { text-align: center; padding: 100rpx 0; color: #999; font-size: 28rpx; }
.user-item { display: flex; align-items: center; gap: 20rpx; padding: 24rpx 30rpx; background: #fff; border-bottom: 1rpx solid #eee; }
.avatar, .avatar-placeholder { width: 72rpx; height: 72rpx; border-radius: 50%; }
.avatar-placeholder { display: flex; justify-content: center; align-items: center; background: #667eea; color: #fff; font-size: 28rpx; }
.user-info { flex: 1; }
.name { font-size: 28rpx; font-weight: 600; display: block; }
.school { font-size: 22rpx; color: #999; }
.add-btn { font-size: 24rpx; padding: 8rpx 24rpx; background: #667eea; color: #fff; border-radius: 24rpx; }
</style>
