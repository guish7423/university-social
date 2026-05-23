<template>
  <div class="notif-page">
    <div class="page-header">
      <h1>通知</h1>
      <el-button v-if="notifications.length" text @click="handleMarkRead">全部已读</el-button>
    </div>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
    <div v-else-if="!notifications.length" class="empty-state"><el-empty description="暂无通知" /></div>
    <div v-else class="notif-list">
      <div v-for="n in notifications" :key="n.id" class="stagger-item" :class="['notif-item', { unread: !n.is_read }]" @click="handleClick(n)">
        <div class="notif-icon" :class="n.type">
          <el-icon :size="18"><component :is="iconMap[n.type] || iconMap.default" /></el-icon>
        </div>
        <el-avatar :size="36" :src="n.from_user?.avatar" class="notif-avatar" />
        <div class="notif-body">
          <div class="notif-content">{{ n.content }}</div>
          <div class="notif-time">{{ formatTime(n.created_at) }}</div>
        </div>
        <div v-if="!n.is_read" class="unread-dot" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { listNotifications, markNotificationsRead } from "@/api/social"
import type { NotificationData } from "@/api/social"
import dayjs from "dayjs"
import relativeTime from "dayjs/plugin/relativeTime"
import "dayjs/locale/zh-cn"
import { Star, ChatDotSquare, User, Bell } from "@element-plus/icons-vue"

dayjs.extend(relativeTime)
dayjs.locale("zh-cn")

const router = useRouter()
const notifications = ref<NotificationData[]>([])
const loading = ref(true)

const iconMap: Record<string, object> = {
  like: Star,
  comment: ChatDotSquare,
  reply: ChatDotSquare,
  follow: User,
  system: Bell,
  default: Bell,
}

function formatTime(t: string) {
  const d = dayjs(t)
  if (dayjs().diff(d, "day") < 1) return d.fromNow()
  if (dayjs().diff(d, "day") < 7) return d.format("dddd HH:mm")
  return d.format("MM-DD HH:mm")
}

function handleClick(n: NotificationData) {
  if (n.type === "follow" && n.from_user_id) {
    router.push("/user/" + n.from_user_id)
  } else if (n.ref_id) {
    router.push("/post/" + n.ref_id)
  }
}

async function handleMarkRead() {
  try {
    await markNotificationsRead()
    notifications.value.forEach(n => n.is_read = true)
  } catch { /* handled */ }
}

onMounted(async () => {
  try { notifications.value = await listNotifications() }
  catch { /* handled */ } finally { loading.value = false }
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
  padding: 12px; border-radius: $radius-md; transition: background 0.2s; cursor: pointer;
  &:hover { background: $bg-hover; }
  &.unread { background: rgba($brand-primary-hex, 0.04); position: relative; }

  .notif-icon {
    display: flex; align-items: center; justify-content: center;
    width: 32px; height: 32px; border-radius: 50%; flex-shrink: 0;
    &.like { background: rgba(#e74c3c, 0.1); color: #e74c3c; }
    &.comment, &.reply { background: rgba(#3498db, 0.1); color: #3498db; }
    &.follow { background: rgba(#2ecc71, 0.1); color: #2ecc71; }
    &.system { background: rgba($brand-primary-hex, 0.1); color: $brand-primary-hex; }
  }

  .notif-avatar { flex-shrink: 0; }

  .notif-body { flex: 1; min-width: 0;
    .notif-content { font-size: 14px; line-height: 1.5; color: $text-primary; }
    .notif-time { font-size: 12px; color: $text-muted; margin-top: 4px; }
  }

  .unread-dot {
    width: 8px; height: 8px; border-radius: 50%; background: $primary; flex-shrink: 0;
  }
}
</style>
