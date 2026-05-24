<template>
  <div class="whisper-page">
    <div class="page-header">
      <h1>树洞</h1>
      <el-button type="primary" @click="showCreate = true">倾诉</el-button>
    </div>

    <LoadingWrapper :loading="loading && !whispers.length" :data="whispers.length" skeleton-variant="post-card" :rows="3">
      <template #empty>
        <el-empty description="暂无树洞消息" />
      </template>
      <div class="whisper-list">
      <div v-for="w in whispers" class="whisper-card stagger-item" :key="w.id" @click="$router.push('/whispers/' + w.id)">
        <div class="author-row">
          <div v-if="w.is_anonymous" class="anon-avatar" :style="{ background: anonGradient(w.user_id) }">{{ w.content.charAt(0) }}</div>
          <template v-else-if="w.author">
            <img :src="w.author?.avatar" :alt="w.author?.nickname" class="real-avatar" loading="lazy" />
            <span class="nickname">{{ w.author?.nickname || '匿名' }}</span>
          </template>
          <span v-if="w.is_anonymous" class="whisper-id">树洞 #{{ String(w.user_id).slice(-4) }}</span>
        </div>
        <p class="content clamp-3">{{ w.content }}</p>
        <div class="footer">
          <span class="time">{{ formatTime(w.created_at) }}</span>
          <span><el-icon><Goods /></el-icon> {{ w.like_count }}</span>
          <span><el-icon><ChatDotSquare /></el-icon> {{ w.comment_count }}</span>
        </div>
      </div>
      <div v-if="hasMore" class="load-more">
        <el-button :loading="loading" @click="loadMore" text>加载更多</el-button>
      </div>
      <div v-if="!hasMore && whispers.length > 0" class="no-more">没有更多了</div>
      </div>
    </LoadingWrapper>

    <el-dialog v-model="showCreate" title="说点什么" width="480">
      <el-input v-model="text" type="textarea" :rows="4" placeholder="匿名倾诉你的心声..." maxlength="500" show-word-limit />
      <el-checkbox v-model="anonymous" style="margin: 12px 0">匿名发布</el-checkbox>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" :loading="submitting" :disabled="!text.trim()" @click="handleSubmit">发布</el-button>
      </template>
    </el-dialog>
  </div>
</template>

import LoadingWrapper from "@/components/LoadingWrapper.vue"
<script setup lang="ts">
import { ref } from "vue"
import { createWhisper } from "@/api/whisper"
import type { WhisperData } from "@/api/whisper"
import { usePagination } from "@/composables/usePagination"
import { useTimeFormat } from "@/composables/useTimeFormat"
import { listWhispers } from "@/api/whisper"

const { formatTime } = useTimeFormat()
const { items: whispers, loading, hasMore, loadMore, reset } = usePagination<WhisperData>({
  fetchFn: (offset, limit) => listWhispers(offset, limit),
})

const showCreate = ref(false)
const text = ref("")
const anonymous = ref(true)
const submitting = ref(false)

async function handleSubmit() {
  if (!text.value.trim()) return
  submitting.value = true
  try {
    await createWhisper({ content: text.value.trim(), is_anonymous: anonymous.value })
    showCreate.value = false
    text.value = ""
    reset()
  } finally { submitting.value = false }
}

const anonColors = [
  'linear-gradient(135deg, #D4A574, #B8865E)',
  'linear-gradient(135deg, #5B8C5A, #3D6B3E)',
  'linear-gradient(135deg, #5B7FB4, #3D5F8E)',
  'linear-gradient(135deg, #8B6FAF, #6B4F8E)',
  'linear-gradient(135deg, #4A9E9E, #2C7E7E)',
  'linear-gradient(135deg, #C47A5A, #A45A3A)',
  'linear-gradient(135deg, #9A5A8A, #7A3A6A)',
  'linear-gradient(135deg, #5A8A7A, #3A6A5A)',
]

function anonGradient(userId: number): string {
  return anonColors[userId % anonColors.length]
}

loadMore()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;
@use "sass:color";

.whisper-page { max-width: 640px; }
.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;
  h1 { font-size: 22px; font-weight: 700; }
}

.whisper-list { display: flex; flex-direction: column; gap: 10px; }

.whisper-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px; cursor: pointer;
  transition: all 0.25s ease;
  &:hover {
    border-color: rgba($brand-primary-hex, 0.3);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.2);
  }
  .content {
    font-size: 14px; line-height: 1.7; margin-bottom: 10px;
    font-style: italic; color: $text-secondary;
  }
  .clamp-3 {
    display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical;
    overflow: hidden; word-break: break-word;
  }
  .footer {
    display: flex; gap: 14px; font-size: 12px; color: $text-muted;
    .time { margin-right: auto; }
    .el-icon { vertical-align: middle; margin-right: 3px; }
  }
}

.author-row {
  display: flex; align-items: center; gap: 10px; margin-bottom: 8px;
  .nickname { font-size: 13px; font-weight: 600; }
  .whisper-id { font-size: 11px; color: $text-muted; margin-left: auto; }
}

.anon-avatar {
  width: 36px; height: 36px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 14px; font-weight: 700; color: #fff; flex-shrink: 0;
}

.real-avatar {
  width: 36px; height: 36px; border-radius: 50%; object-fit: cover; flex-shrink: 0;
}

.load-more { text-align: center; padding: 20px 0; }
.no-more { text-align: center; padding: 20px 0; color: $text-muted; font-size: 13px; }

.skeleton-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px;
  .skeleton-row {
    background: linear-gradient(90deg, var(--skeleton-from, oklch(0.20 0.015 30)) 25%, var(--skeleton-to, oklch(0.23 0.015 30)) 50%, var(--skeleton-from, oklch(0.20 0.015 30)) 75%);
    background-size: 200% 100%;
    animation: skeleton-shimmer 1.5s infinite;
    border-radius: 4px; height: 12px; margin-bottom: 8px;
  }
  .skeleton-author { width: 80px; height: 10px; }
  .skeleton-line { width: 100%; }
  .skeleton-line-short { width: 60%; }
  .skeleton-footer { width: 120px; height: 10px; margin-bottom: 0; }
}

@keyframes skeleton-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
