<template>
  <div class="list-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>
    <PageHeader title="TA的帖子" />

    <LoadingWrapper :loading="loading && !posts.length" :data="posts.length" skeleton-variant="post-card" skeleton-images>
      <template #empty>
        <el-empty description="暂无帖子" />
      </template>

      <template #default>
        <div class="post-list">
          <PostCard
            v-for="p in posts"
            :key="p.id"
            :post="p"
            class="stagger-item"
            @like="handleLike"
            @comment="(p: any) => $router.push('/post/' + p.id)"
          />
        </div>
        <div v-if="hasMore" class="load-more">
          <el-button :loading="loading" @click="loadMore" text>加载更多</el-button>
        </div>
        <div v-if="!hasMore && posts.length > 0" class="no-more">没有更多了</div>
      </template>
    </LoadingWrapper>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router"
import type { PostData } from "@/api/post"
import { userPosts, toggleLike } from "@/api/post"
import { usePagination } from "@/composables/usePagination"
import PageHeader from "@/components/PageHeader.vue"
import PostCard from "@/components/PostCard.vue"
import LoadingWrapper from "@/components/LoadingWrapper.vue"
import { ArrowLeft } from "@element-plus/icons-vue"

const route = useRoute()
const userId = Number(route.params.id)

const { items: posts, loading, hasMore, loadMore } = usePagination<PostData>({
  fetchFn: (offset, limit) => userPosts(userId, offset, limit),
})

async function handleLike(post: PostData) {
  try {
    const res = await toggleLike(post.id)
    post.is_liked = res.liked
    post.like_count += res.liked ? 1 : -1
  } catch { /* handled by interceptor */ }
}

loadMore()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

 margin: 0 auto;
.list-page { max-width: 640px; }
.back-btn { margin-bottom: 4px; }

.post-list {
  display: flex;
  flex-direction: column;
  gap: $space-3;
}

.load-more { text-align: center; padding: $space-6 0; }
.no-more { text-align: center; padding: $space-6 0; color: $text-muted; font-size: $text-sm; }
</style>
