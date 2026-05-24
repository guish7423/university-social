<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
  <div v-else-if="!circle" class="empty-state"><el-empty description="圈子不存在" /></div>
  <div v-else class="detail-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>

    <div class="circle-banner">
      <el-avatar :size="60" :src="circle.avatar" />
      <div class="circle-info">
        <h1>{{ circle.name }}</h1>
        <p>{{ circle.description }}</p>
        <span class="meta">{{ circle.member_count }} 人 · {{ circle.post_count }} 帖</span>
        </div>
      <el-button v-if="!circle.is_joined" type="primary" @click="handleJoin">加入圈子</el-button>
      <el-button v-else @click="handleLeave">退出圈子</el-button>
      <el-button v-if="circle.is_joined" type="primary" @click="showCreate = true">发帖</el-button>
    </div>

    <div v-if="postsLoading" class="loading-wrap"><el-skeleton :rows="3" animated /></div>
    <div v-else-if="!posts.length" class="empty-state"><el-empty description="暂无帖子" /></div>
    <div v-else class="post-list">
      <div v-for="p in posts" :key="p.id" class="post-card">
        <div class="post-header">
          <el-avatar :size="28" :src="p.author?.avatar" />
          <span class="name">{{ p.author?.nickname || '匿名' }}</span>
          <span class="time">{{ formatTime(p.created_at) }}</span>
        </div>
        <p class="post-content">{{ p.content }}</p>
        <div class="post-meta">
          <span><el-icon><Goods /></el-icon> {{ p.like_count }}</span>
          <span><el-icon><ChatDotSquare /></el-icon> {{ p.comment_count }}</span>
        </div>
      </div>
    </div>
    
    <el-dialog v-model="showCreate" title="发帖" width="520">
      <el-input v-model="postContent" type="textarea" :rows="5" placeholder="写点什么..." maxlength="2000" show-word-limit />
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" :loading="submitting" :disabled="!postContent.trim()" @click="handleCreatePost">发布</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { getCircle, joinCircle, leaveCircle, listCirclePosts, createCirclePost } from "@/api/circle"
import type { CircleData, CirclePostData } from "@/api/circle"
import { useTimeFormat } from "@/composables/useTimeFormat"
import { ElMessage } from "element-plus"
import { ArrowLeft } from "@element-plus/icons-vue"

const route = useRoute()
const router = useRouter()

const circle = ref<CircleData | null>(null)
const posts = ref<CirclePostData[]>([])
const loading = ref(true)
const postsLoading = ref(true)
const showCreate = ref(false)
const postContent = ref("")
const submitting = ref(false)

const { formatTime } = useTimeFormat()

async function handleJoin() {
  if (!circle.value) return
  try {
    await joinCircle(circle.value.id)
    circle.value.is_joined = true
    circle.value.member_count++
    ElMessage.success("已加入圈子")
  } catch { /* handled by interceptor */ }
}

async function handleLeave() {
  if (!circle.value) return
  try {
    await leaveCircle(circle.value.id)
    circle.value.is_joined = false
    circle.value.member_count--
    ElMessage.success("已退出圈子")
  } catch { /* handled by interceptor */ }
}

async function handleCreatePost() {
  if (!postContent.value.trim() || !circle.value) return
  submitting.value = true
  try {
    await createCirclePost(circle.value.id, { content: postContent.value.trim() })
    ElMessage.success("发布成功")
    showCreate.value = false
    postContent.value = ""
    posts.value = await listCirclePosts(circle.value.id)
  } catch { /* handled by interceptor */ }
  finally { submitting.value = false }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try { circle.value = await getCircle(id) }
  catch { /* handled by interceptor */ } finally { loading.value = false }

  try { posts.value = await listCirclePosts(id) }
  catch { /* handled by interceptor */ } finally { postsLoading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.detail-page { max-width: 700px; }
.back-btn { margin-bottom: 16px; }

.circle-banner {
  display: flex; align-items: center; gap: 16px;
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 24px; margin-bottom: 20px;
  .circle-info {
    flex: 1;
    h1 { font-size: 20px; font-weight: 700; }
    p { font-size: 13px; color: $text-secondary; margin: 4px 0; }
    .meta { font-size: 12px; color: $text-muted; }
  }
}

.post-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px; margin-bottom: 12px;
  .post-header {
    display: flex; align-items: center; gap: 8px; margin-bottom: 10px;
    .name { font-size: 13px; font-weight: 600; }
    .time { font-size: 11px; color: $text-muted; margin-left: auto; }
  }
  .post-content { font-size: 14px; line-height: 1.6; margin-bottom: 10px; }
  .post-meta {
    display: flex; gap: 16px; font-size: 12px; color: $text-muted;
    .el-icon { vertical-align: middle; margin-right: 3px; }
  }
}
</style>
