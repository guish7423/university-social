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

    <u-empty v-if="notes && notes.length === 0" mode="message" text="暂无通知" />
    <view v-for="(n, i) in (notes || [])" :key="n.id"
          :class="['note-item', !n.is_read && 'unread', 'stagger-' + ((i % 8) + 1)]">
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
.container { min-height: 100vh; background: $color-canvas; }
.header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 30rpx; background: $color-surface;
  border-bottom: 1rpx solid $color-hairline;
}
.title { font-size: 36rpx; font-weight: 700; color: $ink; }
.unread-badge { font-size: 24rpx; color: $ink-muted; }
.empty { text-align: center; padding: 200rpx 0; color: $ink-muted; font-size: 26rpx; }
.note-item {
  display: flex; gap: 20rpx; padding: 24rpx 30rpx;
  background: $color-surface; border-bottom: 1rpx solid $color-hairline;
  opacity: 0; transform: translateY(10rpx);
  animation: fadeInUp 0.4s ease-out forwards;
}
.note-item.unread { background: rgba(198,122,106,0.06); }
.avatar { width: 64rpx; height: 64rpx; border-radius: 50%; background: $color-hairline; }
.note-body { flex: 1; }
.note-name { font-size: 26rpx; font-weight: 600; color: $ink; display: block; }
.note-content { font-size: 26rpx; color: $ink-muted; display: block; margin: 4rpx 0; }
.note-time { font-size: 20rpx; color: $ink-tertiary; }

.stagger-1 { animation-delay: 0ms; }
.stagger-2 { animation-delay: 60ms; }
.stagger-3 { animation-delay: 120ms; }
.stagger-4 { animation-delay: 180ms; }
.stagger-5 { animation-delay: 240ms; }
.stagger-6 { animation-delay: 300ms; }
.stagger-7 { animation-delay: 360ms; }
.stagger-8 { animation-delay: 420ms; }
</style>
