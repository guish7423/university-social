<template>
  <div class="chat-page">
    <div class="page-header"><h1>聊天</h1></div>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
    <div v-else-if="!conversations.length" class="empty-state"><el-empty description="暂无对话" /></div>
    <div v-else class="conv-list">
      <div v-for="c in conversations" :key="c.user_id" class="conv-row" @click="$router.push('/chat/' + c.user_id)">
        <el-avatar :size="44" :src="c.avatar" />
        <div class="conv-info">
          <div class="conv-name">{{ c.nickname }}</div>
          <div class="conv-msg">{{ c.last_message }}</div>
        </div>
        <div class="conv-meta">
          <div class="conv-time">{{ formatTime(c.last_time) }}</div>
          <el-badge v-if="c.unread" :value="c.unread" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listConversations, unreadCount } from "@/api/chat"
import type { ConversationData } from "@/api/chat"
import dayjs from "dayjs"

const conversations = ref<ConversationData[]>([])
const loading = ref(true)

function formatTime(t: string) {
  const d = dayjs(t)
  if (dayjs().diff(d, "day") < 1) return d.format("HH:mm")
  return d.format("MM-DD")
}

onMounted(async () => {
  try { conversations.value = await listConversations() }
  catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.chat-page { max-width: 640px; }
.page-header { margin-bottom: 16px; h1 { font-size: 22px; font-weight: 700; } }

.conv-list { display: flex; flex-direction: column; gap: 4px; }
.conv-row {
  display: flex; align-items: center; gap: 12px;
  padding: 12px; border-radius: $radius-md; cursor: pointer;
  transition: background 0.2s;
  &:hover { background: $bg-card; }
  .conv-info { flex: 1; min-width: 0;
    .conv-name { font-size: 14px; font-weight: 600; }
    .conv-msg { font-size: 13px; color: $text-muted; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  }
  .conv-meta { text-align: right;
    .conv-time { font-size: 11px; color: $text-muted; margin-bottom: 4px; }
  }
}
</style>
