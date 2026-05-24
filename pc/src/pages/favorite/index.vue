<template>
  <div class="favorite-page">
    <PageHeader title="我的收藏" description="收藏的精彩动态" />

    <LoadingWrapper :loading="loading">
      <template #empty>
        <div class="empty-state">
          <el-icon :size="48"><Star /></el-icon>
          <p>还没有收藏任何动态</p>
          <el-button type="primary" @click="$router.push('/square')">去发现</el-button>
        </div>
      </template>

      <template #default>
        <div v-if="posts.length" class="post-list">
          <PostCard
            v-for="post in posts"
            :key="post.id"
            :post="post"
            @like="handleLike"
            @comment="handleComment"
            @favorite="handleFavorite"
            @circle="handleCircle"
          />
        </div>

        <div v-if="total > posts.length" class="load-more">
          <el-button :loading="loadingMore" @click="loadMore">加载更多</el-button>
        </div>
      </template>
    </LoadingWrapper>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"
import { useRouter } from "vue-router"
import { listFavorites, removeFavorite } from "@/api/favorite"
import { toggleLike } from "@/api/post"
import type { PostCardData } from "@/components/PostCard.vue"
import PostCard from "@/components/PostCard.vue"
import PageHeader from "@/components/PageHeader.vue"
import LoadingWrapper from "@/components/LoadingWrapper.vue"

const router = useRouter()
const posts = ref<PostCardData[]>([])
const loading = ref(true)
const loadingMore = ref(false)
const page = ref(1)
const total = ref(0)

onMounted(async () => {
  await loadPosts()
  loading.value = false
})

async function loadPosts() {
  try {
    const res = await listFavorites(page.value, 20)
    posts.value = (res.posts || []).map((p: any) => ({
      id: p.id,
      content: p.content,
      images: p.images,
      like_count: p.like_count,
      comment_count: p.comment_count,
      is_liked: p.is_liked,
      is_favorited: true,
      created_at: p.created_at,
      author: p.author,
    }))
    total.value = res.total || 0
  } catch {
    /* */
  }
}

async function loadMore() {
  loadingMore.value = true
  page.value++
  const prev = posts.value.length
  try {
    const res = await listFavorites(page.value, 20)
    const newPosts = (res.posts || []).map((p: any) => ({
      id: p.id,
      content: p.content,
      images: p.images,
      like_count: p.like_count,
      comment_count: p.comment_count,
      is_liked: p.is_liked,
      is_favorited: true,
      created_at: p.created_at,
      author: p.author,
    }))
    posts.value.push(...newPosts)
    total.value = res.total || 0
  } catch {
    page.value--
  } finally {
    loadingMore.value = false
  }
}

async function handleLike(post: PostCardData) {
  try {
    await toggleLike(post.id)
    post.is_liked = !post.is_liked
    post.like_count += post.is_liked ? 1 : -1
  } catch { /* */ }
}

function handleComment(post: PostCardData) {
  router.push(`/post/${post.id}`)
}

async function handleFavorite(post: PostCardData) {
  try {
    await removeFavorite(post.id)
    posts.value = posts.value.filter(p => p.id !== post.id)
    total.value--
  } catch { /* */ }
}

function handleCircle(post: PostCardData) {
  if (post.circle) router.push(`/circles/${post.circle.id}`)
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.favorite-page {
  max-width: 640px;
  margin: 0 auto;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $space-4;
  padding: $space-16 $space-4;
  color: $text-muted;
  p { font-size: $text-sm; }
}

.post-list {
  margin-top: $space-4;
}

.load-more {
  text-align: center;
  padding: $space-8;
}
</style>
