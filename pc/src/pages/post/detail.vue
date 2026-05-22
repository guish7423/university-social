<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="8" animated /></div>
  <div v-else-if="!post" class="empty-state"><el-empty description="动态不存在" /></div>
  <div v-else class="detail-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>

    <div class="post-detail">
      <div class="post-header">
        <el-avatar :size="40" :src="post.author?.avatar" />
        <div class="post-author-info">
          <span class="name">{{ post.author?.nickname || '匿名' }}</span>
          <span class="time">{{ formatTime(post.created_at) }}</span>
        </div>
        <el-button v-if="post.user_id === userStore.userId" text type="danger" :icon="Delete" @click="handleDelete" size="small">
          删除
        </el-button>
      </div>

      <p class="post-content">{{ post.content }}</p>

      <div v-if="post.images?.length" class="post-images">
        <img v-for="(img, i) in post.images" :key="i" :src="img" />
      </div>

      <div class="post-actions">
        <span :class="['action', { active: post.is_liked }]" @click="handleLike">
          <el-icon><Goods /></el-icon> {{ post.like_count }}
        </span>
        <span class="action"><el-icon><ChatDotSquare /></el-icon> {{ post.comment_count }}</span>
      </div>
    </div>

    <div class="comments-section">
      <h3>评论 ({{ post.comment_count }})</h3>
      <div class="comment-input">
        <el-input v-model="commentText" placeholder="写下你的评论..." size="large" @keyup.enter="handleComment" />
        <el-button type="primary" @click="handleComment" :disabled="!commentText.trim()">发送</el-button>
      </div>

      <div v-if="commentsLoading" class="loading-wrap"><el-skeleton :rows="3" animated /></div>
      <div v-else-if="!comments.length" class="empty-state"><el-empty description="暂无评论" /></div>
      <div v-else class="comment-list">
        <div v-for="c in comments" :key="c.id" class="comment-item">
          <el-avatar :size="28" :src="c.author?.avatar" />
          <div class="comment-body">
            <span class="comment-author">{{ c.author?.nickname || '匿名' }}</span>
            <p>{{ c.content }}</p>
            <span class="comment-time">{{ formatTime(c.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useUserStore } from "@/stores/user"
import { getPost, deletePost, listComments, createComment, toggleLike } from "@/api/post"
import type { PostData, CommentData } from "@/api/post"
import dayjs from "dayjs"
import relativeTime from "dayjs/plugin/relativeTime"
import "dayjs/locale/zh-cn"
import { ArrowLeft } from "@element-plus/icons-vue"
import { ElMessageBox } from "element-plus"
import { Delete } from "@element-plus/icons-vue"

dayjs.extend(relativeTime)
dayjs.locale("zh-cn")

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const post = ref<PostData | null>(null)
const comments = ref<CommentData[]>([])
const loading = ref(true)
const commentsLoading = ref(true)
const commentText = ref("")

function formatTime(t: string) {
  const d = dayjs(t)
  if (dayjs().diff(d, "day") < 1) return d.fromNow()
  return d.format("MM-DD HH:mm")
}

async function handleLike() {
  if (!post.value) return
  try {
    const res = await toggleLike(post.value.id)
    post.value.is_liked = res.liked
    post.value.like_count += res.liked ? 1 : -1
  } catch { /* handled by interceptor */ }
}

async function handleDelete() {
  if (!post.value) return
  try {
    await deletePost(post.value.id)
      await ElMessageBox.confirm("确定删除这条动态吗？", "确认", {
        confirmButtonText: "删除",
        cancelButtonText: "取消",
        type: "warning",
      })
    router.push("/square")
  } catch { /* handled by interceptor */ }
}

async function handleComment() {
  if (!commentText.value.trim() || !post.value) return
  try {
    await createComment(post.value.id, commentText.value.trim())
    commentText.value = ""
    const newComments = await listComments(post.value.id)
    comments.value = newComments
    if (post.value) post.value.comment_count = newComments.length
  } catch { /* handled by interceptor */ }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    post.value = await getPost(id)
  } catch { /* handled by interceptor */ } finally { loading.value = false }

  try {
    comments.value = await listComments(id)
  } catch { /* handled by interceptor */ } finally { commentsLoading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.detail-page { max-width: 680px; }

.back-btn { margin-bottom: 16px; }

.post-detail {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 20px; margin-bottom: 20px;

  .post-header {
    display: flex; align-items: center; gap: 10px; margin-bottom: 14px;
    .post-author-info {
      flex: 1;
      .name { font-size: 15px; font-weight: 600; }
      .time { font-size: 12px; color: $text-muted; margin-left: 8px; }
    }
  }

  .post-content {
    font-size: 15px; line-height: 1.8; white-space: pre-wrap; margin-bottom: 14px;
  }

  .post-images {
    display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 14px;
    img { width: 120px; height: 120px; object-fit: cover; border-radius: $radius-sm; }
  }

  .post-actions {
    display: flex; gap: 24px; padding-top: 14px; border-top: 1px solid $border-color;
    .action {
      display: flex; align-items: center; gap: 4px; font-size: 14px;
      color: $text-secondary; cursor: pointer;
      &.active { color: $primary; }
      &:hover { color: $primary; }
    }
  }
}

.comments-section {
  h3 { font-size: 16px; font-weight: 600; margin-bottom: 14px; }
}

.comment-input {
  display: flex; gap: 10px; margin-bottom: 16px;
}

.comment-list { margin-top: 16px; }

.comment-item {
  display: flex; gap: 10px; padding: 12px 0; border-bottom: 1px solid $border-color;
  .comment-body {
    flex: 1;
    .comment-author { font-size: 13px; font-weight: 600; }
    p { font-size: 14px; line-height: 1.6; margin: 4px 0; }
    .comment-time { font-size: 12px; color: $text-muted; }
  }
}

.loading-wrap, .empty-state { padding: 40px 0; }
</style>
