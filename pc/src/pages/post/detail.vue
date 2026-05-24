<template>
  <SkeletonCard v-if="loading" variant="detail" />
  <div v-else-if="!post" class="empty-state"><el-empty description="动态不存在" /></div>
  <div v-else class="detail-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>

    <div class="post-detail">
      <div class="post-header">
        <el-avatar :size="40" :src="post.author?.avatar" @click="$router.push('/user/' + post.user_id)" style="cursor:pointer" />
        <div class="post-author-info">
          <span class="name" @click="$router.push('/user/' + post.user_id)">{{ post.author?.nickname || '匿名' }}</span>
          <span class="meta">
            <span class="time">{{ formatTime(post.created_at) }}</span>
            <span v-if="post.school" class="school">· {{ post.school }}</span>
          </span>
        </div>
        <el-button v-if="post.user_id === userStore.userId && !editing" text :icon="Edit" @click="handleEdit" size="small">编辑</el-button>
        <el-button v-if="post.user_id === userStore.userId" text type="danger" :icon="Delete" @click="handleDelete" size="small">删除</el-button>
      </div>

      <p v-if="post.title" class="post-title">{{ post.title }}</p>
      <el-input v-if="editing" v-model="editForm.content" type="textarea" :rows="4" maxlength="2000" show-word-limit />
      <p v-else class="post-content">{{ post.content }}</p>

      <div v-if="post.images?.length" class="post-images">
        <img v-for="(img, i) in post.images" :key="i" :src="img" @click="lightboxIndex = i" :class="{ single: post.images.length === 1 }" />
      </div>
      <div v-if="editing" class="edit-actions">
        <el-button type="primary" @click="handleSave" :disabled="!editForm.content.trim()">保存</el-button>
        <el-button @click="handleCancel">取消</el-button>
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
        <el-input
          v-model="commentText"
          :placeholder="replyTo ? '回复 @' + replyTo.author?.nickname + '...' : '写下你的评论...'"
          size="large"
          @keyup.enter="handleComment"
        >
          <template v-if="replyTo" #prefix>
            <el-tag closable size="small" @close="replyTo = null" style="margin-right:4px">回复 @{{ replyTo.author?.nickname }}</el-tag>
          </template>
        </el-input>
        <el-button type="primary" @click="handleComment" :disabled="!commentText.trim()">发送</el-button>
      </div>

      <SkeletonCard v-if="commentsLoading" variant="post-card" :rows="3" />
      <div v-else-if="!comments.length" class="empty-state"><el-empty description="暂无评论，来抢沙发吧" :image-size="80" /></div>
      <div v-else class="comment-list">
        <div v-for="c in comments" :key="c.id" class="comment-item">
          <el-avatar :size="28" :src="c.author?.avatar" @click="$router.push('/user/' + c.user_id)" style="cursor:pointer" />
          <div class="comment-body">
            <span class="comment-author" @click="$router.push('/user/' + c.user_id)">{{ c.author?.nickname || '匿名' }}</span>
            <p>{{ c.content }}</p>
            <div class="comment-footer">
              <span class="comment-time">{{ formatTime(c.created_at) }}</span>
              <span class="comment-reply" @click="startReply(c)">回复</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ImagePreview :images="post?.images || []" v-model="lightboxIndex" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useUserStore } from "@/stores/user"
import { getPost, deletePost, listComments, createComment, toggleLike, updatePost } from "@/api/post"
import type { PostData, CommentData } from "@/api/post"
import ImagePreview from "@/components/ImagePreview.vue"

import SkeletonCard from "@/components/SkeletonCard.vue"
import dayjs from "dayjs"
import relativeTime from "dayjs/plugin/relativeTime"
import "dayjs/locale/zh-cn"
import { ArrowLeft, Edit, Delete } from "@element-plus/icons-vue"
import { ElMessageBox, ElMessage } from "element-plus"

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
const replyTo = ref<CommentData | null>(null)
const lightboxIndex = ref<number | null>(null)

const editing = ref(false)
const editForm = reactive({
	content: "",
	images: [] as string[],
	topic_id: undefined as number | undefined,
})

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
    await ElMessageBox.confirm("确定删除这条动态吗？此操作不可撤销。", "确认删除", {
      confirmButtonText: "删除", cancelButtonText: "取消", type: "warning",
    })
    await deletePost(post.value.id)
    ElMessage.success("已删除")
    router.push("/square")
  } catch { /* cancelled or error */ }
}

function handleEdit() {
	if (!post.value) return
	editForm.content = post.value.content
	editForm.images = [...(post.value.images || [])]
	editForm.topic_id = post.value.topic_id
	editing.value = true
}
async function handleSave() {
	if (!post.value || !editForm.content.trim()) return
	try {
		await updatePost(post.value.id, {
			content: editForm.content.trim(),
			images: editForm.images,
			topic_id: editForm.topic_id,
		})
		post.value.content = editForm.content.trim()
		post.value.images = editForm.images
		post.value.topic_id = editForm.topic_id
		editing.value = false
		ElMessage.success("编辑成功")
	} catch { /* handled by interceptor */ }
}
function handleCancel() {
	editing.value = false
}

function startReply(comment: CommentData) {
  replyTo.value = comment
  commentText.value = ""
}

async function handleComment() {
  if (!commentText.value.trim() || !post.value) return
  const text = replyTo.value
    ? "回复 @" + replyTo.value.author?.nickname + "：" + commentText.value.trim()
    : commentText.value.trim()
  try {
    await createComment(post.value.id, text)
    commentText.value = ""
    replyTo.value = null
    const newComments = await listComments(post.value.id)
    comments.value = newComments
    if (post.value) post.value.comment_count = newComments.length
    ElMessage.success("评论成功")
  } catch { /* handled by interceptor */ }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try { post.value = await getPost(id) }
  catch { /* handled by interceptor */ }
  finally { loading.value = false }

  try { comments.value = await listComments(id) }
  catch { /* handled by interceptor */ }
  finally { commentsLoading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.detail-page { max-width: 680px; margin: 0 auto; }

.back-btn { margin-bottom: $space-4; }

.post-detail {
  background: $bg-card; border: 1px solid $border-default;
  border-radius: $radius-md; padding: $space-6; margin-bottom: $space-5;

  .post-header {
    display: flex; align-items: center; gap: $space-3; margin-bottom: $space-4;
    .post-author-info { flex: 1;
      .name { font-size: $text-base; font-weight: 600; cursor: pointer;
        &:hover { color: $brand-primary; }
      }
      .meta { display: block; font-size: $text-xs; color: $text-muted; margin-top: 2px; }
    }
  }

  .post-title { font-size: $text-xl; font-weight: 700; margin-bottom: $space-3; line-height: 1.4; }
  .post-content { font-size: $text-base; line-height: 1.8; white-space: pre-wrap; margin-bottom: $space-4; color: $text-primary; }

  .post-images {
    display: flex; flex-wrap: wrap; gap: $space-2; margin-bottom: $space-4;
    img { width: 120px; height: 120px; object-fit: cover; border-radius: $radius-sm; cursor: zoom-in;
      transition: opacity $duration-fast;
      &:hover { opacity: 0.85; }
      &.single { max-width: 300px; width: auto; aspect-ratio: auto; max-height: 300px; }
    }
  }

  .post-actions {
    display: flex; gap: $space-6; padding-top: $space-4; border-top: 1px solid $border-default;
    .action {
      display: flex; align-items: center; gap: $space-1; font-size: $text-sm;
      color: $text-secondary; cursor: pointer; transition: color $duration-fast;
      &.active, &:hover { color: $brand-primary; }
    }
  }
}

.edit-actions { display: flex; gap: $space-3; margin-bottom: $space-4; }

.comments-section {
  h3 { font-size: $text-base; font-weight: 600; margin-bottom: $space-4; }
}

.comment-input {
  display: flex; gap: $space-3; margin-bottom: $space-4;
  :deep(.el-input__wrapper) { background: $bg-card; }
}

.comment-list { margin-top: $space-4; }

.comment-item {
  display: flex; gap: $space-3; padding: $space-4 0; border-bottom: 1px solid $border-default;
  &:last-child { border-bottom: none; }

  .comment-body { flex: 1;
    .comment-author { font-size: $text-sm; font-weight: 600; cursor: pointer;
      &:hover { color: $brand-primary; }
    }
    p { font-size: $text-sm; line-height: 1.6; margin: $space-1 0; color: $text-primary; }
    .comment-footer { display: flex; align-items: center; gap: $space-3; }
    .comment-time { font-size: $text-xs; color: $text-muted; }
    .comment-reply { font-size: $text-xs; color: $brand-primary; cursor: pointer;
      &:hover { text-decoration: underline; }
    }
  }
}

.empty-state { padding: $space-12 0; }
</style>

<style lang="scss">
.lightbox-overlay {
  position: fixed; inset: 0; z-index: 9999;
  background: rgba(0,0,0,0.85); display: flex; align-items: center; justify-content: center;
  .lightbox-img { max-width: 90vw; max-height: 90vh; object-fit: contain; border-radius: 4px; }
}
.lightbox-nav {
  position: absolute; width: 100%; display: flex; justify-content: space-between; padding: 0 20px; pointer-events: none; top: 50%; transform: translateY(-50%);
  .nav-btn { pointer-events: auto; font-size: 48px; color: #fff; cursor: pointer; opacity: 0.6; transition: opacity 150ms; user-select: none;
    &:hover { opacity: 1; }
  }
}
</style>
