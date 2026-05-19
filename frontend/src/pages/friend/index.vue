<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listFriends, sendFriendRequest, removeFriend, type UserInfo } from "@/api/social"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const friends = ref<UserInfo[]>([])
const searchMode = ref(false)
const searchQuery = ref("")

onMounted(() => loadFriends())

async function loadFriends() {
  try { friends.value = await listFriends() } catch {}
}

function goSearch() {
  searchMode.value = !searchMode.value
  if (!searchMode.value) loadFriends()
}

function goChat(friend: UserInfo) {
  uni.navigateTo({ url: `/pages/friend/chat?id=${friend.id}&name=${friend.nickname}` })
}
</script>

<template>
  <view class="container">
    <view class="header">
      <text class="title">我的好友</text>
      <view class="header-actions">
        <button class="action-btn" @click="goSearch">{{ searchMode ? "完成" : "添加" }}</button>
      </view>
    </view>

    <view v-if="searchMode" class="search-section">
      <input v-model="searchQuery" class="search-input" placeholder="搜索用户昵称..." />
      <button class="search-btn" @click="goSearch">搜索</button>
    </view>

    <view v-if="friends.length === 0" class="empty">还没有好友</view>
    <view v-for="f in friends" :key="f.id" class="friend-item">
      <image v-if="f.avatar" class="avatar" :src="f.avatar" mode="aspectFill" />
      <view class="avatar-placeholder">{{ f.nickname?.[0] || "?" }}</view>
      <view class="friend-info">
        <text class="friend-name">{{ f.nickname }}</text>
        <text v-if="f.school" class="friend-school">{{ f.school }}</text>
      </view>
      <button class="chat-btn" @click="goChat(f)">私信</button>
    </view>
  </view>
</template>

<style scoped>
.container { min-height: 100vh; background: #f5f5f5; }
.header { display: flex; justify-content: space-between; align-items: center; padding: 30rpx; background: #fff; }
.title { font-size: 36rpx; font-weight: 700; }
.action-btn { font-size: 26rpx; color: #667eea; background: none; border: none; }
.search-section { display: flex; gap: 16rpx; padding: 20rpx 30rpx; background: #fff; }
.search-input { flex: 1; height: 64rpx; background: #f0f0f0; border-radius: 32rpx; padding: 0 24rpx; font-size: 26rpx; }
.search-btn { font-size: 26rpx; color: #667eea; }
.empty { text-align: center; padding: 200rpx 0; color: #999; }
.friend-item { display: flex; align-items: center; gap: 20rpx; padding: 24rpx 30rpx; background: #fff; border-bottom: 1rpx solid #eee; }
.avatar, .avatar-placeholder { width: 80rpx; height: 80rpx; border-radius: 50%; }
.avatar-placeholder { display: flex; justify-content: center; align-items: center; background: #667eea; color: #fff; font-size: 32rpx; }
.friend-info { flex: 1; }
.friend-name { font-size: 28rpx; font-weight: 600; display: block; }
.friend-school { font-size: 22rpx; color: #999; }
.chat-btn { font-size: 24rpx; padding: 8rpx 24rpx; background: #667eea; color: #fff; border-radius: 24rpx; }
</style>
