<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
  <div v-else-if="!whisper" class="empty-state"><el-empty description="树洞不存在" /></div>
  <div v-else class="detail-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>

    <div class="whisper-detail">
      <p class="content">{{ whisper.content }}</p>
      <div class="actions">
        <span :class="['action', { active: whisper.is_liked, 'like-animating': likeAnimating }]" @click="handleLike">
          <el-icon><Goods /></el-icon> {{ whisper.like_count }}
        </span>
      </div>
    </div>

    <div class="comment-input-wrap" ref="commentInputRef">
      <el-input v-model="commentText" type="textarea" :rows="2" placeholder="写下你的回应..." maxlength="500" show-word-limit />
      <div class="comment-actions">
        <el-checkbox v-model="commentAnon">匿名回应</el-checkbox>
        <el-button type="primary" :loading="commentSubmitting" :disabled="!commentText.trim()" @click="handleComment">发送</el-button>
      </div>
    </div>

    <div class="comments">
      <h3>回应 <span class="count">{{ comments.length }}</span></h3>
      <div v-if="!comments.length" class="empty-state"><el-empty description="暂无回应" /></div>
      <div v-else v-for="c in comments" :key="c.id" class="comment-item stagger-item">
        <div class="comment-author">
          <img v-if="c.author" :src="c.author.avatar" class="avatar" loading="lazy" />
          <div v-else class="anon-avatar" :style="{ background: anonGradient(c.id) }">{{ c.content.charAt(0) }}</div>
          <div class="author-info">
            <span class="author">{{ c.author?.nickname || '匿名' }}</span>
            <span class="time">{{ formatTime(c.created_at) }}</span>
          </div>
        </div>
        <p class="comment-content">{{ c.content }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from "vue"
import { useRoute } from "vue-router"
import { getWhisper, toggleWhisperLike, listWhisperComments, createWhisperComment } from "@/api/whisper"
import type { WhisperData } from "@/api/whisper"
import dayjs from "dayjs"
import { ArrowLeft } from "@element-plus/icons-vue"

const route = useRoute()
const whisper = ref<WhisperData | null>(null)
const comments = ref<{ id: number; content: string; created_at: string; author?: { nickname: string; avatar: string } }[]>([])
const loading = ref(true)
const likeAnimating = ref(false)
const commentText = ref("")
const commentAnon = ref(true)
const commentSubmitting = ref(false)
const commentInputRef = ref<HTMLElement | null>(null)

function formatTime(t: string) { return dayjs(t).format("MM-DD HH:mm") }

const anonColors = [
  "linear-gradient(135deg, #D4A574, #B8865E)",
  "linear-gradient(135deg, #5B8C5A, #3D6B3E)",
  "linear-gradient(135deg, #5B7FB4, #3D5F8E)",
  "linear-gradient(135deg, #8B6FAF, #6B4F8E)",
  "linear-gradient(135deg, #4A9E9E, #2C7E7E)",
  "linear-gradient(135deg, #C47A5A, #A45A3A)",
  "linear-gradient(135deg, #9A5A8A, #7A3A6A)",
  "linear-gradient(135deg, #5A8A7A, #3A6A5A)",
]

function anonGradient(id: number) { return anonColors[id % anonColors.length] }

async function handleLike() {
  if (!whisper.value) return
  try {
    const res = await toggleWhisperLike(whisper.value.id)
    whisper.value.is_liked = res.liked
    whisper.value.like_count += res.liked ? 1 : -1
    likeAnimating.value = true
    setTimeout(() => likeAnimating.value = false, 300)
  } catch { /* handled by interceptor */ }
}

async function handleComment() {
  if (!commentText.value.trim() || !whisper.value) return
  commentSubmitting.value = true
  try {
    await createWhisperComment(whisper.value.id, commentText.value.trim())
    commentText.value = ""
    comments.value = await listWhisperComments(whisper.value.id)
    await nextTick()
    commentInputRef.value?.scrollIntoView({ behavior: "smooth", block: "center" })
  } catch { /* handled */ } finally { commentSubmitting.value = false }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    whisper.value = await getWhisper(id)
    comments.value = await listWhisperComments(id)
    if (route.query.scroll === "comment") {
      await nextTick()
      document.querySelector(".comments")?.scrollIntoView({ behavior: "smooth" })
    }
  } catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

  margin: 0 auto;
.detail-page { max-width: 640px; margin: 0 auto; }
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

@keyframes like-pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.3); }
  100% { transform: scale(1); }
}
.action.like-animating { animation: like-pulse 0.3s ease; }

.comment-input-wrap {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px; margin-bottom: 20px;
}
.comment-actions {
  display: flex; align-items: center; justify-content: space-between; margin-top: 10px;
}

.comments {
  h3 { font-size: 16px; font-weight: 600; margin-bottom: 12px;
    .count { color: $text-muted; font-weight: 400; font-size: 14px; }
  }
  .comment-item {
    background: $bg-card; border-radius: $radius-md; padding: 14px; margin-bottom: 10px;
    .comment-author { display: flex; align-items: center; gap: 10px; margin-bottom: 6px;
      .avatar { width: 32px; height: 32px; border-radius: 50%; object-fit: cover; }
      .anon-avatar {
        width: 32px; height: 32px; border-radius: 50%;
        display: flex; align-items: center; justify-content: center;
        font-size: 12px; font-weight: 700; color: #fff; flex-shrink: 0;
      }
      .author-info {
        .author { font-size: 13px; font-weight: 600; }
        .time { font-size: 12px; color: $text-muted; margin-left: 8px; }
      }
    }
    .comment-content { font-size: 14px; line-height: 1.6; color: $text-secondary; }
  }
}
</style>
