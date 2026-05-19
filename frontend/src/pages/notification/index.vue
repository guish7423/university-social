<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listNotifications, markNotificationsRead, type NotificationData } from "@/api/social"

const notes = ref<NotificationData[]>([])
const unreadCount = ref(0)

onMounted(() => load())

async function load() {
  try {
    notes.value = await listNotifications()
    unreadCount.value = notes.value.filter((n) => !n.is_read).length
    if (unreadCount.value > 0) markNotificationsRead()
  } catch {}
}

function formatTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  return `${d.getMonth() + 1}/${d.getDate()}`
}

function getTypeText(type: string) {
  const map: Record<string, string> = {
    like: "赞了你的动态",
    comment: "评论了你的动态",
    friend_request: "请求添加好友",
    friend_accepted: "通过了你的好友请求",
  }
  return map[type] || type
}
</script>

<template>
  <view class="container">
    <view class="header">
      <text class="title">消息通知</text>
      <text v-if="unreadCount" class="unread-badge">{{ unreadCount }} 条未读</text>
    </view>

    <view v-if="notes.length === 0" class="empty">暂无通知</view>
    <view v-for="n in notes" :key="n.id" :class="['note-item', !n.is_read && 'unread']">
      <image v-if="n.from_user?.avatar" class="avatar" :src="n.from_user.avatar" mode="aspectFill" />
      <view class="note-body">
        <text class="note-name">{{ n.from_user?.nickname || "系统" }}</text>
        <text class="note-content">{{ getTypeText(n.type) }}</text>
        <text class="note-time">{{ formatTime(n.created_at) }}</text>
      </view>
    </view>
  </view>
</template>

<style scoped>
.container { min-height: 100vh; background: #f5f5f5; }
.header { display: flex; justify-content: space-between; align-items: center; padding: 30rpx; background: #fff; }
.title { font-size: 36rpx; font-weight: 700; }
.unread-badge { font-size: 24rpx; color: #e74c3c; }
.empty { text-align: center; padding: 200rpx 0; color: #999; }
.note-item { display: flex; gap: 20rpx; padding: 24rpx 30rpx; background: #fff; border-bottom: 1rpx solid #eee; }
.note-item.unread { background: #f0f4ff; }
.avatar { width: 64rpx; height: 64rpx; border-radius: 50%; }
.note-body { flex: 1; }
.note-name { font-size: 26rpx; font-weight: 600; display: block; }
.note-content { font-size: 26rpx; color: #333; display: block; margin: 4rpx 0; }
.note-time { font-size: 20rpx; color: #999; }
</style>
