<script setup lang="ts">
import { ref, onMounted } from "vue"
import { searchUsers, sendFriendRequest, type SearchUserResult } from "@/api/social"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const query = ref("")
const results = ref<SearchUserResult[]>([])
const loading = ref(false)

async function handleSearch() {
  if (!query.value.trim()) return
  loading.value = true
  try {
    results.value = await searchUsers(query.value)
  } catch {}
  loading.value = false
}

const statusMap: Record<string, string> = {
  friend: "已是好友",
  pending_sent: "已发送",
  pending_received: "待同意",
}

async function addFriend(id: number) {
  await sendFriendRequest(id)
  const u = results.value.find(r => r.id === id)
  if (u) u.friend_status = "pending_sent"
}
</script>

<template>
  <view class="page">
    <view class="search-bar">
      <input v-model="query" class="input" placeholder="搜索用户昵称" @confirm="handleSearch" />
      <view class="search-btn" @click="handleSearch">搜索</view>
    </view>

    <view v-if="loading" class="state">搜索中...</view>
    <view v-else-if="results.length === 0 && query" class="state">未找到用户</view>
    <view v-else-if="results.length === 0" class="state">输入昵称搜索</view>

    <view v-for="u in results" :key="u.id" class="user-item">
      <image v-if="u.avatar" class="avatar" :src="u.avatar" mode="aspectFill" />
      <view v-else class="avatar-placeholder">{{ u.nickname?.[0] || "?" }}</view>
      <view class="user-info">
        <text class="name">{{ u.nickname }}</text>
        <text v-if="u.school" class="school">{{ u.school }}</text>
      </view>
      <view v-if="u.id === userStore.id" class="add-btn disabled">自己</view>
      <view v-else-if="u.friend_status" class="add-btn disabled">{{ statusMap[u.friend_status] }}</view>
      <view v-else class="add-btn" @click="addFriend(u.id)">加好友</view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); }
.search-bar { display: flex; gap: 16rpx; padding: 30rpx; background: var(--color-surface, #fff); }
.input { flex: 1; height: 72rpx; background: var(--color-surface-1, #F0EDE8); border-radius: 36rpx; padding: 0 24rpx; font-size: 28rpx; color: var(--ink, #1E2A3A); }
.search-btn { font-size: 28rpx; color: var(--brand-primary, #C67A6A); background: none; padding: 0; }
.state { text-align: center; padding: 100rpx 0; color: var(--ink-tertiary, #B8C2CE); font-size: 28rpx; }
.user-item { display: flex; align-items: center; gap: 20rpx; padding: 24rpx 30rpx; background: var(--color-surface, #fff); border-bottom: 1rpx solid var(--hairline-light, #EAE6E0); }
.avatar, .avatar-placeholder { width: 72rpx; height: 72rpx; border-radius: 50%; }
.avatar-placeholder { display: flex; justify-content: center; align-items: center; background: var(--brand-primary, #C67A6A); color: #fff; font-size: 28rpx; }
.user-info { flex: 1; }
.name { font-size: 28rpx; font-weight: 600; color: var(--ink, #1E2A3A); display: block; }
.school { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); }
.add-btn { font-size: 24rpx; padding: 8rpx 24rpx; background: var(--brand-primary, #C67A6A); color: #fff; border-radius: 24rpx; }
.add-btn.disabled { background: var(--ink-tertiary, #B8C2CE); color: #fff; }
</style>
