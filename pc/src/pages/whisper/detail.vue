<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
  <div v-else-if="!whisper" class="empty-state"><el-empty description="树洞不存在" /></div>
  <div v-else class="detail-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>
    <div class="whisper-detail">
      <p class="content">{{ whisper.content }}</p>
      <div class="actions">
        <span :class="['action', { active: whisper.is_liked }]" @click="handleLike">
          <el-icon><Goods /></el-icon> {{ whisper.like_count }}
        </span>
      </div>
    </div>
    <div class="comments">
      <h3>回应</h3>
      <div v-if="!comments.length" class="empty-state"><el-empty description="暂无回应" /></div>
      <div v-else v-for="c in comments" :key="c.id" class="comment-item">
        <span class="author">{{ c.author?.nickname || '匿名' }}</span>
        <p>{{ c.content }}</p>
        <span class="time">{{ formatTime(c.created_at) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"
import { getWhisper, toggleWhisperLike, listWhisperComments } from "@/api/whisper"
import type { WhisperData } from "@/api/whisper"
import dayjs from "dayjs"
import { ArrowLeft } from "@element-plus/icons-vue"

const route = useRoute()
const whisper = ref<WhisperData | null>(null)
const comments = ref<{ id: number; content: string; created_at: string; author?: { nickname: string } }[]>([])
const loading = ref(true)

function formatTime(t: string) { return dayjs(t).format("MM-DD HH:mm") }

async function handleLike() {
  if (!whisper.value) return
  try {
    const res = await toggleWhisperLike(whisper.value.id)
    whisper.value.is_liked = res.liked
    whisper.value.like_count += res.liked ? 1 : -1
  } catch { /* handled by interceptor */ }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    whisper.value = await getWhisper(id)
    comments.value = await listWhisperComments(id)
  } catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.detail-page { max-width: 640px; }
.back-btn { margin-bottom: 16px; }

.whisper-detail {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 24px; margin-bottom: 20px;
  .content { font-size: 15px; line-height: 1.8; font-style: italic; color: $text-secondary; margin-bottom: 16px; }
  .actions { display: flex; gap: 16px;
    .action { display: flex; align-items: center; gap: 4px; cursor: pointer; color: $text-secondary;
      &.active { color: $primary; }
      &:hover { color: $primary; }
    }
  }
}

.comments {
  h3 { font-size: 16px; font-weight: 600; margin-bottom: 12px; }
  .comment-item { padding: 10px 0; border-bottom: 1px solid $border-color;
    .author { font-size: 13px; font-weight: 600; }
    p { font-size: 14px; margin: 4px 0; }
    .time { font-size: 12px; color: $text-muted; }
  }
}
</style>
