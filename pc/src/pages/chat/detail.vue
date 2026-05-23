<template>
  <div class="chat-detail-page">
    <div class="chat-header">
      <el-button text :icon="ArrowLeft" @click="$router.push('/chat')" />
      <span class="chat-title">{{ peerName }}</span>
      <el-tag v-if="online" size="small" type="success" effect="dark">在线</el-tag>
      <el-tag v-else size="small" type="info" effect="plain">离线</el-tag>
    </div>
    <div class="messages" ref="msgContainer">
      <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
      <div v-else-if="!messages.length" class="empty-state"><el-empty description="开始聊天吧" /></div>
      <div v-else v-for="m in messages" :key="m.id" :class="['msg', m.from_user_id === userId ? 'self' : 'peer']">
        <div class="bubble">{{ m.content }}</div>
        <div class="time">{{ formatTime(m.created_at) }}</div>
      </div>
    </div>
    <div class="typing-hint" v-if="typing">对方正在输入...</div>
    <div class="input-area">
      <el-input v-model="text" placeholder="输入消息..." @keyup.enter="handleSend" @input="onTyping" />
      <el-button type="primary" :disabled="!text.trim()" @click="handleSend">发送</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch } from "vue"
import { useRoute } from "vue-router"
import { useUserStore } from "@/stores/user"
import { getMessages } from "@/api/chat"
import type { MessageData } from "@/api/chat"
import { useWebSocket } from "@/composables/useWebSocket"
import dayjs from "dayjs"
import { ArrowLeft } from "@element-plus/icons-vue"

const route = useRoute()
const userStore = useUserStore()
const userId = userStore.userId
const peerId = Number(route.params.userId)

const messages = ref<MessageData[]>([])
const text = ref("")
const loading = ref(true)
const online = ref(false)
const typing = ref(false)
let typingTimer: ReturnType<typeof setTimeout> | null = null
const peerName = ref("")
const msgContainer = ref<HTMLElement | null>(null)

const { connect, sendMessage: wsSend, lastMessage, connected } = useWebSocket(userId)

function formatTime(t: string) { return dayjs(t).format("MM-DD HH:mm") }

watch(lastMessage, (msg) => {
  if (!msg) return
  if (msg.type === "connected") {
    online.value = true
    return
  }
  if (msg.type === "message" && msg.data) {
    const exists = messages.value.some(m => m.id === msg.data!.id)
    if (!exists) {
      if (msg.data.sender_id !== userId) {
        messages.value.push(msg.data as unknown as MessageData)
        nextTick(() => scrollToBottom())
      }
    }
  }
  if (msg.type === "typing" && msg.sender_id === peerId) {
    typing.value = true
    if (typingTimer) clearTimeout(typingTimer)
    typingTimer = setTimeout(() => { typing.value = false }, 3000)
  }
})

watch(connected, (v) => { online.value = v })

function scrollToBottom() {
  msgContainer.value?.scrollTo({ top: msgContainer.value.scrollHeight, behavior: "smooth" })
}

async function handleSend() {
  if (!text.value.trim()) return
  wsSend(peerId, text.value.trim())
  text.value = ""
  const msgs = await getMessages(peerId)
  messages.value = msgs
  nextTick(() => scrollToBottom())
}

function onTyping() {
  wsSend(peerId, "", "typing")
}

onMounted(async () => {
  connect()
  try {
    const msgs = await getMessages(peerId)
    messages.value = msgs
    if (msgs.length) {
      const peer = msgs.find(m => m.from_user_id !== userId)
      if (peer) peerName.value = `用户 ${peer.from_user_id}`
    }
  } catch { /* handled */ } finally { loading.value = false }
  nextTick(() => scrollToBottom())
})

onUnmounted(() => {
  if (typingTimer) clearTimeout(typingTimer)
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.chat-detail-page {
  display: flex; flex-direction: column; height: calc(100vh - 120px);
  max-width: 640px;
}

.chat-header {
  display: flex; align-items: center; gap: 8px;
  padding-bottom: 12px; border-bottom: 1px solid $border-color; margin-bottom: 12px;
  .chat-title { font-size: 16px; font-weight: 600; }
}

.typing-hint {
  font-size: 12px; color: $text-muted; padding: 4px 0; font-style: italic;
}

.messages {
  flex: 1; overflow-y: auto; padding: 12px 0;
  display: flex; flex-direction: column; gap: 12px;
}

.msg {
  max-width: 70%;
  .bubble {
    padding: 10px 14px; border-radius: $radius-md; font-size: 14px; line-height: 1.5;
  }
  .time { font-size: 11px; color: $text-muted; margin-top: 4px; }

  &.self {
    align-self: flex-end;
    .bubble { background: $primary; color: white; border-bottom-right-radius: 4px; }
    .time { text-align: right; }
  }
  &.peer {
    align-self: flex-start;
    .bubble { background: $bg-card; border: 1px solid $border-color; border-bottom-left-radius: 4px; }
  }
}

.input-area {
  display: flex; gap: 10px; padding-top: 12px; border-top: 1px solid $border-color;
}
</style>
