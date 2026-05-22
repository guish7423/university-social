<template>
  <div class="square-page">
    <PageHeader title="广场" description="发现洛阳高校的新鲜事">
      <template #actions>
        <el-button type="primary" :icon="Edit" @click="$router.push('/post/create')">发布动态</el-button>
      </template>
    </PageHeader>

    <div class="tabs">
      <el-radio-group v-model="currentTab" @change="handleTabChange">
        <el-radio-button value="all">全部</el-radio-button>
        <el-radio-button value="following">关注</el-radio-button>
        <el-radio-button value="hot">精华</el-radio-button>
      </el-radio-group>
    </div>

    <LoadingWrapper :loading="loading && !posts.length" :empty="{}">
      <template #loading>
        <div class="loading-area"><el-skeleton :rows="5" animated /></div>
      </template>

      <template #empty>
        <div class="empty-state">
          <el-empty description="暂无动态">
            <el-button type="primary" @click="$router.push('/post/create')">发布第一条动态</el-button>
          </el-empty>
        </div>
      </template>

      <template #default>
        <div class="post-list">
          <PostCard
            v-for="post in posts"
            :key="post.id"
            :post="post"
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
import { ref } from "vue"
import { toggleLike, listPosts, followingPosts } from "@/api/post"
import type { PostData } from "@/api/post"
import { usePagination } from "@/composables/usePagination"
import PostCard from "@/components/PostCard.vue"
import PageHeader from "@/components/PageHeader.vue"
import LoadingWrapper from "@/components/LoadingWrapper.vue"
import { Edit } from "@element-plus/icons-vue"

const currentTab = ref("all")
const fetchFn = (offset: number, limit: number) =>
  currentTab.value === "following" ? followingPosts(offset, limit) : listPosts({ offset, limit })



const { items: posts, loading, hasMore, loadMore, reset } = usePagination<PostData>({ fetchFn })

async function handleLike(post: PostData) {
  try {
    const res = await toggleLike(post.id)
    post.is_liked = res.liked
    post.like_count += res.liked ? 1 : -1
  } catch { /* handled by interceptor */ }
}

function handleTabChange() {
  reset()
}

loadMore()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.square-page { max-width: 680px; margin: 0 auto; }

.tabs { margin-bottom: $space-5; }

.loading-area { padding: $space-6 0; }

.empty-state { padding: $space-12 0; }

.post-list {
  display: flex; flex-direction: column; gap: $space-3;
}

.load-more { text-align: center; padding: $space-6 0; }
.no-more { text-align: center; padding: $space-6 0; color: $text-muted; font-size: $text-sm; }
</style>
