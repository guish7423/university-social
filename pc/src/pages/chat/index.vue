<template>
  <div class="chat-page">
    <div class="page-header">
      <h1>消息</h1>
    </div>
    <div class="search-bar">
      <el-input v-model="searchQuery" placeholder="搜索聊天" :prefix-icon="SearchIcon" clearable />
    </div>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
    <div v-else-if="!filteredConvs.length" class="empty-state">
      <el-empty :description="searchQuery ? '未找到匹配的聊天' : '暂无聊天'" />
    </div>
    <div v-else class="conv-list">
      <div
        v-for="c in filteredConvs"
        :key="c.user_id"
        class="conv-row stagger-item"
        @click="$router.push('/chat/' + c.user_id)"
      >
        <div class="avatar-wrap">
          <el-avatar :size="44" :src="c.avatar" />
          <span class="online-dot" :class="{ online: connected }" />
        </div>
        <div class="conv-info">
          <div class="conv-top">
            <span class="conv-name">{{ c.nickname }}</span>
            <span class="conv-time">{{ formatTime(c.last_time) }}</span>
          </div>
          <div class="conv-bottom">
            <span class="conv-msg">{{ c.last_message }}</span>
            <el-badge v-if="c.unread" :value="c.unread" :max="99" class="unread-badge" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useUserStore } from "@/stores/user"
import { listConversations } from "@/api/chat"
import type { ConversationData } from "@/api/chat"
import { useWebSocket } from "@/composables/useWebSocket"
import dayjs from "dayjs"
import { Search as SearchIcon } from "@element-plus/icons-vue"

const userStore = useUserStore()
const { connect, connected } = useWebSocket(userStore.userId)

const searchQuery = ref("")
const conversations = ref<ConversationData[]>([])
const loading = ref(true)

const filteredConvs = computed(() => {
  if (!searchQuery.value) return conversations.value
  const q = searchQuery.value.toLowerCase()
  return conversations.value.filter(
    (c) => c.nickname.toLowerCase().includes(q) || c.last_message.toLowerCase().includes(q)
  )
})

function formatTime(t: string) {
  const d = dayjs(t)
  const now = dayjs()
  if (d.isSame(now, "day")) return d.format("HH:mm")
  if (d.isSame(now.subtract(1, "day"), "day")) return "昨天"
  return d.format("MM-DD")
}

onMounted(async () => {
  connect()
  try {
    conversations.value = await listConversations()
  } catch { /* handled by interceptor */ } finally {
    loading.value = false
  }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.chat-page { max-width: 640px; }

.page-header {
  display: flex; align-items: center; margin-bottom: 12px;
  h1 { font-size: 22px; font-weight: 700; }
}

.search-bar { margin-bottom: 12px; }

.conv-list { display: flex; flex-direction: column; gap: 2px; }

.conv-row {
  display: flex; align-items: center; gap: 12px;
  padding: 12px; border-radius: $radius-md; cursor: pointer;
  transition: all $duration-fast;
  &:hover { background: $bg-card; }

  .avatar-wrap {
    position: relative; flex-shrink: 0;
    .online-dot {
      position: absolute; bottom: 0; right: 0;
      width: 10px; height: 10px; border-radius: 50%;
      background: $text-muted; border: 2px solid $bg-app;
      transition: background $duration-fast;
      &.online { background: $color-success; }
    }
  }

  .conv-info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 4px;
    .conv-top { display: flex; justify-content: space-between; align-items: center;
      .conv-name { font-size: 15px; font-weight: 600; }
      .conv-time { font-size: 11px; color: $text-muted; white-space: nowrap; }
    }
    .conv-bottom { display: flex; justify-content: space-between; align-items: center;
      .conv-msg { font-size: 13px; color: $text-muted; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
    }
  }
}

.unread-badge :deep(.el-badge__content) { font-size: 11px; }
</style>
