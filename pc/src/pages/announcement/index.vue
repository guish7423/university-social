<template>
  <div class="announcement-page">
    <PageHeader title="校园公告" />
    <div v-if="loading" class="loading-area">
      <SkeletonCard v-for="i in 4" :key="i" variant="post-card" />
    </div>
    <div v-else-if="!list.length" class="empty-card">
      <el-empty description="暂无公告" :image-size="80" />
    </div>
    <div v-else class="announce-list">
      <div
        v-for="item in list"
        :key="item.id"
        class="announce-card"
      >
        <div class="ac-header">
          <h3 class="ac-title">{{ item.title }}</h3>
          <span class="ac-time">{{ formatTime(item.created_at) }}</span>
        </div>
        <p class="ac-content">{{ item.content }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listAnnouncements } from "@/api/social"
import PageHeader from "@/components/PageHeader.vue"
import SkeletonCard from "@/components/SkeletonCard.vue"

const list = ref<{ id: number; title: string; content: string; created_at: string }[]>([])
const loading = ref(true)

function formatTime(ts: string) {
  const d = new Date(ts)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 86400000) return "今天"
  if (diff < 172800000) return "昨天"
  return `${d.getMonth() + 1}-${d.getDate()}`
}

onMounted(async () => {
  try {
    const res = await listAnnouncements()
    list.value = (res as any) || []
  } catch { /* */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.announcement-page { max-width: 1100px; }

.announce-list {
  display: flex;
  flex-direction: column;
  gap: $space-4;
}

.announce-card {
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-md;
  padding: $space-5;
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  }
}

.ac-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: $space-3;
}

.ac-title {
  font-size: $text-lg;
  font-weight: 600;
  color: $text-primary;
  margin: 0;
}

.ac-time {
  font-size: $text-xs;
  color: $text-secondary;
  white-space: nowrap;
  margin-left: $space-3;
}

.ac-content {
  font-size: $text-sm;
  color: $text-secondary;
  line-height: 1.7;
  margin: 0;
  white-space: pre-line;
}
</style>
