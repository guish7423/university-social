<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listFriends, listFriendRequests, removeFriend, type UserInfo } from "@/api/social"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const friends = ref<UserInfo[]>([])
const pendingCount = ref(0)

onMounted(() => {
  loadFriends()
  loadPendingCount()
})

async function loadFriends() {
  try { friends.value = (await listFriends()) || [] } catch {}
}
  try { friends.value = (await listFriends()) || [] } catch {}

async function loadPendingCount() {
  try {
    const reqs = await listFriendRequests()
    pendingCount.value = (reqs || []).length
  } catch {}
}

function goRequests() {
  uni.navigateTo({ url: "/pages/friend/requests" })
}

function goSearch() {
  uni.navigateTo({ url: "/pages/friend/search" })
}

function goChat(friend: UserInfo) {
  uni.navigateTo({ url: `/pages/friend/chat?id=${friend.id}&name=${friend.nickname}` })
}

async function handleRemove(id: number) {
  try {
    await removeFriend(id)
    friends.value = friends.value.filter(f => f.id !== id)
    uni.showToast({ title: "已删除", icon: "none" })
  } catch {
    uni.showToast({ title: "操作失败", icon: "error" })
  }
}

function confirmRemove(id: number, name: string) {
  uni.showModal({
    title: "删除好友",
    content: `确定删除好友 ${name}？`,
    success: (res) => {
      if (res.confirm) handleRemove(id)
    }
  })
}
</script>

<template>
  <view class="container">
    <view class="header">
      <text class="title">好友</text>
      <view class="header-actions">
        <view class="action-btn" @click="goRequests">
          <text>请求</text>
          <text v-if="pendingCount > 0" class="badge">{{ pendingCount }}</text>
        </view>
        <view class="action-btn" @click="goSearch">添加</view>
      </view>
    </view>

    <view v-if="(friends || []).length === 0" class="empty">还没有好友</view>
    <view v-for="f in friends" :key="f.id" class="friend-item" @longpress="confirmRemove(f.id, f.nickname)">
      <u-avatar :src="f.avatar" size="72" :name="f.nickname?.[0] || '?'" />
      <view class="friend-info">
        <text class="friend-name">{{ f.nickname }}</text>
        <text v-if="f.school" class="friend-school">{{ f.school }}</text>
      </view>
      <u-button size="small" type="primary" class="chat-btn" @click="goChat(f)">私信</u-button>
    </view>
  </view>
</template>

<style scoped>
.container { min-height: 100vh; background: #f5f5f5; }
.header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 30rpx; background: #fff;
}
.title { font-size: 36rpx; font-weight: 700; }
.header-actions { display: flex; gap: 24rpx; }
.action-btn {
  font-size: 26rpx; color: #C67A6A; position: relative;
  display: flex; align-items: center;
}
.badge {
  display: inline-flex; align-items: center; justify-content: center;
  min-width: 32rpx; height: 32rpx; border-radius: 16rpx;
  background: #e74c3c; color: #fff; font-size: 20rpx;
  margin-left: 6rpx; padding: 0 6rpx;
}
.empty { text-align: center; padding: 200rpx 0; color: #999; }
.friend-item {
  display: flex; align-items: center; gap: 20rpx;
  padding: 24rpx 30rpx; background: #fff;
  border-bottom: 1rpx solid #eee;
}
.friend-info { flex: 1; }
.friend-name { font-size: 28rpx; font-weight: 600; display: block; }
.friend-school { font-size: 22rpx; color: #999; }
.chat-btn { min-width: 100rpx; }
</style>
