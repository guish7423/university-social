<template>
  <div class="whisper-page">
    <div class="page-header">
      <h1>树洞</h1>
      <el-button type="primary" @click="showCreate = true">倾诉</el-button>
    </div>

    <div v-if="loading && !whispers.length" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
    <div v-else-if="!whispers.length" class="empty-state"><el-empty description="暂无树洞消息" /></div>
    <div v-else class="whisper-list">
      <div v-for="w in whispers" class="whisper-card stagger-item" :key="w.id" @click="$router.push('/whispers/' + w.id)">
        <p class="content">{{ w.content }}</p>
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

loadMore()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.whisper-page { max-width: 640px; }
.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;
  h1 { font-size: 22px; font-weight: 700; }
}

.whisper-list { display: flex; flex-direction: column; gap: 10px; }

.whisper-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px; cursor: pointer; transition: all 0.2s;
  &:hover { border-color: rgba($brand-primary-hex, 0.3); }
  .content {
    font-size: 14px; line-height: 1.7; margin-bottom: 10px;
    font-style: italic; color: $text-secondary;
  }
  .footer {
    display: flex; gap: 14px; font-size: 12px; color: $text-muted;
    .time { margin-right: auto; }
    .el-icon { vertical-align: middle; margin-right: 3px; }
  }
}
.load-more { text-align: center; padding: 20px 0; }
.no-more { text-align: center; padding: 20px 0; color: $text-muted; font-size: 13px; }
</style>
