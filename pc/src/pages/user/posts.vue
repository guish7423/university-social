<template>
  <div class="list-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>
    <h1>TA的帖子</h1>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="!posts.length" class="empty-state"><el-empty description="暂无帖子" /></div>
    <div v-else class="post-list">
      <div v-for="p in posts" :key="p.id" class="post-card" @click="$router.push('/post/' + p.id)">
        <p>{{ p.content.slice(0, 150) }}</p>
        <div class="meta">
          <span>{{ p.like_count }} 赞</span>
          <span>{{ p.comment_count }} 评论</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"
import { userPosts } from "@/api/post"
import type { PostData } from "@/api/post"
import { ArrowLeft } from "@element-plus/icons-vue"

const route = useRoute()
const userId = Number(route.params.id)
const posts = ref<PostData[]>([])
const loading = ref(true)

onMounted(async () => {
  try { posts.value = await userPosts(userId) }
  catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.list-page { max-width: 640px; }
h1 { font-size: 20px; font-weight: 700; margin: 16px 0; }
.post-list { display: flex; flex-direction: column; gap: 8px; }
.post-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 14px; cursor: pointer;
  p { font-size: 14px; line-height: 1.6; margin-bottom: 8px; }
  .meta { display: flex; gap: 12px; font-size: 12px; color: $text-muted; }
  &:hover { border-color: $primary-light; }
}
</style>
