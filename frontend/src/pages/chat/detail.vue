<template>
  <view class="chat-detail">
    <scroll-view class="msg-list" scroll-y :scroll-into-view="scrollTo" :scroll-with-animation="true" @scrolltolower="loadMore">
      <view v-for="m in msgs" :key="m.id" :id="'msg-'+m.id" class="msg-row" :class="{ mine: m.sender_id === myId }">
        <image class="msg-avatar" :src="m.sender?.avatar || '/static/default-avatar.png'" />
        <view class="msg-content">
          <text class="msg-text">{{ m.content }}</text>
          <text class="msg-time">{{ formatTime(m.created_at) }}</text>
        </view>
      </view>
      <u-loading mode="flower" size="48" v-if="loading" />
    </scroll-view>
    <view class="input-bar">
      <input class="input" v-model="text" :disabled="sending" confirm-type="send" @confirm="send" placeholder="输入消息..." />
      <button class="send-btn" :disabled="!text.trim() || sending" @click="send">发送</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, nextTick } from "vue"
import { getMessages, sendMessage, markRead, type MessageData } from "@/api/chat"
import { onLoad } from "@dcloudio/uni-app"


const msgs = ref<MessageData[]>([])
const text = ref("")
const sending = ref(false)
const loading = ref(false)
const scrollTo = ref("")
const myId = ref(0)
const page = ref(0)
const hasMore = ref(true)
const userId = ref(0)

onLoad((opts: any) => {
  userId.value = parseInt(opts.userId || "0")
  const nickname = decodeURIComponent(opts.nickname || "")
  uni.setNavigationBarTitle({ title: nickname || "聊天" })

  const token = uni.getStorageSync("token")
  if (token) {
    const parts = token.split(".")
    if (parts.length === 3) {
      try {
        myId.value = JSON.parse(atob(parts[1])).user_id
      } catch (e) { console.error(e) }
    }
  }
  loadMessages()
})

function loadMessages() {
  if (!hasMore.value || loading.value) return
  loading.value = true
  getMessages(userId.value, { offset: page.value * 20, limit: 20 }).then(res => {
    const data = res.data || []
    if (data.length < 20) hasMore.value = false
    msgs.value = [...data.reverse(), ...msgs.value]
    page.value++
    nextTick(() => { scrollTo.value = "msg-" + (msgs.value[msgs.value.length - 1]?.id || 0) })
    loading.value = false
  }).catch(() => { loading.value = false })
}

function loadMore() {
  loadMessages()
}

function send() {
  const content = text.value.trim()
  if (!content || sending.value) return
  sending.value = true
  sendMessage({ receiver_id: userId.value, content }).then(res => {
    if (res.data) {
      msgs.value.push(res.data)
      text.value = ""
      nextTick(() => { scrollTo.value = "msg-" + res.data.id })
    }
  }).finally(() => { sending.value = false })
}

function formatTime(t: string) {
  const d = new Date(t)
  return d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })
}
</script>

<style scoped>
.chat-detail { display: flex; flex-direction: column; height: 100vh; background: #F7F4F0; }
.msg-list { flex: 1; overflow-y: auto; padding: 20rpx 30rpx; }
.msg-row { display: flex; margin-bottom: 24rpx; align-items: flex-start; }
.msg-row.mine { flex-direction: row-reverse; }
.msg-avatar { width: 64rpx; height: 64rpx; border-radius: 50%; flex-shrink: 0; }
.msg-content { max-width: 65%; margin: 0 16rpx; display: flex; flex-direction: column; }
.msg-row.mine .msg-content { align-items: flex-end; }
.msg-text { background: #fff; padding: 16rpx 20rpx; border-radius: 12rpx; font-size: 28rpx; color: #333; line-height: 1.5; word-break: break-word; }
.msg-row.mine .msg-text { background: #C67A6A; color: #fff; }
.msg-time { font-size: 22rpx; color: #8E9BAB; margin-top: 6rpx; }
.loading { text-align: center; color: #8E9BAB; font-size: 26rpx; padding: 20rpx; }
.input-bar { display: flex; align-items: center; padding: 16rpx 20rpx; background: #fff; border-top: 1rpx solid #eee; }
.input { flex: 1; height: 72rpx; background: #F7F4F0; border-radius: 36rpx; padding: 0 24rpx; font-size: 28rpx; }
.send-btn { margin-left: 16rpx; height: 72rpx; padding: 0 32rpx; background: #C67A6A; color: #fff; border-radius: 36rpx; font-size: 28rpx; line-height: 72rpx; border: none; }
.send-btn:disabled { opacity: 0.5; }
</style>
