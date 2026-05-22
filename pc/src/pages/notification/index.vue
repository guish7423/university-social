<template>
  <div class="notif-page">
    <div class="page-header">
      <h1>通知</h1>
      <el-button v-if="notifications.length" text @click="handleMarkRead">全部已读</el-button>
    </div>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
    <div v-else-if="!notifications.length" class="empty-state"><el-empty description="暂无通知" /></div>
    <div v-else class="notif-list">
      <div v-for="n in notifications" class="stagger-item" :key="n.id" :class="['notif-item', { unread: !n.is_read }]">
        <el-avatar :size="36" :src="n.from_user?.avatar" />
        <div class="notif-body">
          <div class="notif-content">{{ n.content }}</div>
          <div class="notif-time">{{ formatTime(n.created_at) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listNotifications, markNotificationsRead } from "@/api/social"
import type { NotificationData } from "@/api/social"
import dayjs from "dayjs"
import relativeTime from "dayjs/plugin/relativeTime"
import "dayjs/locale/zh-cn"

dayjs.extend(relativeTime)
dayjs.locale("zh-cn")

const notifications = ref<NotificationData[]>([])
const loading = ref(true)

function formatTime(t: string) {
  const d = dayjs(t)
  if (dayjs().diff(d, "day") < 1) return d.fromNow()
  return d.format("MM-DD HH:mm")
}

async function handleMarkRead() {
  try {
    await markNotificationsRead()
    notifications.value.forEach(n => n.is_read = true)
  } catch { /* handled by interceptor */ }
}

onMounted(async () => {
  try { notifications.value = await listNotifications() }
  catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.notif-page { max-width: 640px; }
.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;
  h1 { font-size: 22px; font-weight: 700; }
}

.notif-list { display: flex; flex-direction: column; gap: 6px; }
.notif-item {
  display: flex; align-items: center; gap: 12px;
  padding: 12px; border-radius: $radius-md; transition: background 0.2s;
  &:hover { background: $bg-card; }
  &.unread { border-left: 3px solid $primary; }
  .notif-body { flex: 1;
    .notif-content { font-size: 14px; line-height: 1.5; }
    .notif-time { font-size: 12px; color: $text-muted; margin-top: 4px; }
  }
}
</style>
