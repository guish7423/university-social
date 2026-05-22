<template>
  <div class="chat-detail-page">
    <div class="chat-header">
      <el-button text :icon="ArrowLeft" @click="$router.push('/chat')" />
      <span class="chat-title">{{ peerName }}</span>
    </div>
    <div class="messages" ref="msgContainer">
      <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
      <div v-else-if="!messages.length" class="empty-state"><el-empty description="开始聊天吧" /></div>
      <div v-else v-for="m in messages" :key="m.id" :class="['msg', m.from_user_id === userId ? 'self' : 'peer']">
        <div class="bubble">{{ m.content }}</div>
        <div class="time">{{ formatTime(m.created_at) }}</div>
      </div>
    </div>
    <div class="input-area">
      <el-input v-model="text" placeholder="输入消息..." @keyup.enter="handleSend" />
      <el-button type="primary" :disabled="!text.trim()" @click="handleSend">发送</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from "vue"
import { useRoute } from "vue-router"
import { useUserStore } from "@/stores/user"
import { getMessages, sendMessage } from "@/api/chat"
import type { MessageData } from "@/api/chat"
import dayjs from "dayjs"
import { ArrowLeft } from "@element-plus/icons-vue"

const route = useRoute()
const userStore = useUserStore()
const userId = userStore.userId
const peerId = Number(route.params.userId)

const messages = ref<MessageData[]>([])
const text = ref("")
const loading = ref(true)
const peerName = ref("")
const msgContainer = ref<HTMLElement | null>(null)

function formatTime(t: string) { return dayjs(t).format("MM-DD HH:mm") }

async function handleSend() {
  if (!text.value.trim()) return
  try {
    await sendMessage(peerId, text.value.trim())
    text.value = ""
    messages.value = await getMessages(peerId)
    nextTick(() => msgContainer.value?.scrollTo({ top: msgContainer.value.scrollHeight, behavior: "smooth" }))
  } catch { /* handled by interceptor */ }
}

onMounted(async () => {
  try {
    const msgs = await getMessages(peerId)
    messages.value = msgs
    if (msgs.length) {
      const peer = msgs.find(m => m.from_user_id !== userId)
      if (peer) peerName.value = `用户 ${peer.from_user_id}`
    }
  } catch { /* handled by interceptor */ } finally { loading.value = false }

  nextTick(() => msgContainer.value?.scrollTo({ top: msgContainer.value.scrollHeight }))
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
