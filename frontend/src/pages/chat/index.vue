<template>
  <view class="chat-list">
    <view v-if="convs.length === 0" class="empty-state">
      <text class="empty-title">暂无对话</text>
      <text class="empty-desc">去广场发现有趣的人，开始聊天吧</text>
      <u-button type="primary" shape="circle" class="empty-btn" @click="uni.navigateTo({ url: '/pages/square/index' })">去广场逛逛</u-button>
    </view>
    <view v-for="c in convs" :key="c.other_user_id" class="conv-item" @click="goChat(c.other_user_id, c.other_nickname)">
      <image class="avatar" :src="c.other_avatar || '/static/default-avatar.png'" />
      <view class="conv-info">
        <view class="conv-top">
          <text class="nickname">{{ c.other_nickname }}</text>
          <text class="time">{{ formatTime(c.last_time) }}</text>
        </view>
        <view class="conv-bottom">
          <text class="last-msg">{{ c.last_content }}</text>
          <view v-if="c.unread_count > 0" class="badge">{{ c.unread_count > 99 ? '99+' : c.unread_count }}</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { getConversations, type ConversationData } from "@/api/chat"
import { onShow } from "@dcloudio/uni-app"


const convs = ref<ConversationData[]>([])

onShow(() => {
  loadConvs()
})

function loadConvs() {
  getConversations().then(res => {
    convs.value = res.data || []
  })
}

function goChat(userId: number, nickname: string) {
  uni.navigateTo({ url: `/pages/chat/detail?userId=${userId}&nickname=${encodeURIComponent(nickname)}` })
}

function formatTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 86400000) return d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })
  if (diff < 172800000) return "昨天"
  return d.toLocaleDateString("zh-CN", { month: "2-digit", day: "2-digit" })
}
</script>

<style scoped>
.empty-state { display: flex; flex-direction: column; align-items: center; padding: 120rpx 0; gap: 16rpx; }
.empty-title { font-size: 30rpx; font-weight: 600; color: #1E2A3A; }
.empty-desc { font-size: 26rpx; color: #8E9BAB; }
.empty-btn { margin-top: 8rpx; }

.conv-item { display: flex; padding: 24rpx 30rpx; background: #fff; border-bottom: 1rpx solid #eee; align-items: center; }
.conv-item:active { background: #f5f5f5; }
.avatar { width: 88rpx; height: 88rpx; border-radius: 50%; margin-right: 24rpx; flex-shrink: 0; }
.conv-info { flex: 1; min-width: 0; }
.conv-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8rpx; }
.nickname { font-size: 30rpx; font-weight: 600; color: #333; }
.time { font-size: 24rpx; color: #8E9BAB; flex-shrink: 0; margin-left: 16rpx; }
.conv-bottom { display: flex; justify-content: space-between; align-items: center; }
.last-msg { font-size: 26rpx; color: #8E9BAB; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }
.badge { min-width: 32rpx; height: 32rpx; background: #C67A6A; color: #fff; border-radius: 16rpx; font-size: 22rpx; text-align: center; line-height: 32rpx; padding: 0 8rpx; margin-left: 12rpx; flex-shrink: 0; }
</style>
