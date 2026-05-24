<template>
  <div class="notif-page">
    <PageHeader title="通知">
      <template #extra>
        <el-button v-if="filtered.length" text size="small" @click="handleMarkRead">全部已读</el-button>
      </template>
    </PageHeader>

    <div class="tab-bar">
      <el-radio-group v-model="activeTab" @change="handleTabChange">
        <el-radio-button value="all">全部 <small>{{ totalCount }}</small></el-radio-button>
        <el-radio-button value="interact">互动 <small>{{ interactCount }}</small></el-radio-button>
        <el-radio-button value="system">系统 <small>{{ systemCount }}</small></el-radio-button>
        <el-radio-button value="course">课程 <small>{{ courseCount }}</small></el-radio-button>
      </el-radio-group>
    </div>

    <LoadingWrapper :loading="loading" :data="filtered.length" skeleton-variant="post-card" :rows="6">
      <template #empty>
        <el-empty :description="emptyText" :image-size="80" />
      </template>
      <div class="notif-list">
        <div v-for="n in filtered" :key="n.id" class="stagger-item" :class="['notif-item', { unread: !n.is_read }]" @click="handleClick(n)">
          <div class="notif-icon" :class="n.type">
            <el-icon :size="18"><component :is="iconMap[n.type] || iconMap.default" /></el-icon>
          </div>
          <el-avatar :size="36" :src="n.from_user?.avatar" class="notif-avatar" />
          <div class="notif-body">
            <div class="notif-content">{{ n.content }}</div>
            <div class="notif-time">{{ smartTime(n.created_at) }}</div>
          </div>
          <div v-if="!n.is_read" class="unread-dot" />
        </div>
      </div>
    </LoadingWrapper>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"
import LoadingWrapper from "@/components/LoadingWrapper.vue"
import { listNotifications, markNotificationsRead } from "@/api/social"
import type { NotificationData } from "@/api/social"
import { Star, ChatDotSquare, User, Bell, Reading } from "@element-plus/icons-vue"
import PageHeader from "@/components/PageHeader.vue"

const router = useRouter()
const notifications = ref<NotificationData[]>([])
const loading = ref(true)
const activeTab = ref("all")

const iconMap: Record<string, object> = {
  like: Star,
  comment: ChatDotSquare,
  reply: ChatDotSquare,
  follow: User,
  system: Bell,
  course: Reading,
  default: Bell,
}

const categoryMap: Record<string, string[]> = {
  all: [],
  interact: ["like", "comment", "reply", "follow"],
  system: ["system"],
  course: ["course"],
}

const emptyTextMap: Record<string, string> = {
  all: "暂无通知",
  interact: "暂时没有互动通知",
  system: "暂时没有系统通知",
  course: "暂时没有课程通知",
}

const filtered = computed(() => {
  const cats = categoryMap[activeTab.value]
  if (!cats || !cats.length) return notifications.value
  return notifications.value.filter(n => cats.includes(n.type))
})

const totalCount = computed(() => notifications.value.length)
const interactCount = computed(() => notifications.value.filter(n => ["like","comment","reply","follow"].includes(n.type)).length)
const systemCount = computed(() => notifications.value.filter(n => n.type === "system").length)
const courseCount = computed(() => notifications.value.filter(n => n.type === "course").length)

const emptyText = computed(() => emptyTextMap[activeTab.value] || "暂无通知")

function handleTabChange() { /* filtered recomputes automatically */ }

function smartTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const day = 86400000
  if (diff < day) return d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })
  if (diff < 2 * day) return "昨天 " + d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })
  if (diff < 7 * day) {
    const days = ["周日", "周一", "周二", "周三", "周四", "周五", "周六"]
    return days[d.getDay()] + " " + d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })
  }
  return `${d.getMonth() + 1}-${d.getDate()}`
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
@use "@/styles/variables" as *;
@use "sass:color";

.notif-page { max-width: 640px; }

.tab-bar {
  margin-bottom: $space-5;
  :deep(.el-radio-group) { gap: 0; }
  :deep(.el-radio-button__inner) {
    font-size: $text-sm;
    small { font-size: $text-xs; color: $text-secondary; margin-left: 4px; }
  }
}

.notif-list { display: flex; flex-direction: column; gap: 6px; }

.notif-item {
  display: flex; align-items: center; gap: 12px;
  padding: 12px; border-radius: $radius-md; transition: all 0.3s ease; cursor: pointer;
  &:hover { background: $bg-hover; transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
  &.unread { background: color.change($primary, $alpha: 0.04); position: relative; }

  .notif-icon {
    display: flex; align-items: center; justify-content: center;
    width: 32px; height: 32px; border-radius: 50%; flex-shrink: 0;
    &.like { background: rgba(#e74c3c, 0.1); color: #e74c3c; }
    &.comment, &.reply { background: rgba(#3498db, 0.1); color: #3498db; }
    &.follow { background: rgba(#2ecc71, 0.1); color: #2ecc71; }
    &.system { background: color.change($primary, $alpha: 0.1); color: $primary; }
    &.course { background: rgba(#d4a85e, 0.1); color: #d4a85e; }
  }

  .notif-avatar { flex-shrink: 0; }

  .notif-body { flex: 1; min-width: 0;
    .notif-content { font-size: 14px; line-height: 1.5; color: $text-primary; }
    .notif-time { font-size: 12px; color: $text-secondary; margin-top: 4px; }
  }

  .unread-dot {
    width: 8px; height: 8px; border-radius: 50%; background: $primary; flex-shrink: 0;
  }
}
</style>
