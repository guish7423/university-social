<template>
  <div class="square-page">
    <PageHeader title="广场" description="发现洛阳高校的新鲜事">
      <template #actions>
        <el-button type="primary" :icon="Edit" @click="$router.push('/post/create')">发布动态</el-button>
      </template>
    </PageHeader>

    <div class="tabs">
      <el-radio-group v-model="currentTab" @change="handleTabChange">
        <el-radio-button value="latest">最新</el-radio-button>
        <el-radio-button value="hot">🔥 热门</el-radio-button>
        <el-radio-button value="following">关注</el-radio-button>
        <el-radio-button value="featured">✨ 精选</el-radio-button>
      </el-radio-group>
    </div>

    <div v-if="currentTab === 'following' && !isLogin" class="empty-state">
      <el-empty description="登录后可查看关注动态">
        <el-button type="primary" @click="$router.push('/login')">去登录</el-button>
      </el-empty>
    </div>

    <template v-else>
      <LoadingWrapper :loading="loading && !posts.length" :data="posts.length" :empty="{}" skeleton-variant="post-card" skeleton-images>
        <template #empty>
          <div class="empty-state">
            <el-empty :description="emptyText">
              <el-button type="primary" @click="$router.push('/post/create')">发布第一条动态</el-button>
            </el-empty>
          </div>
        </template>

        <template #default>
          <div class="post-list">
            <div v-for="post in posts" :key="post.id" class="post-card-wrapper stagger-item">
              <PostCard
                :post="post"
                @like="handleLike"
                @comment="(p: any) => $router.push('/post/' + p.id)"
              />
              <div v-if="currentTab === 'hot'" class="trending-badge">
                🔥<span class="trending-label">热门</span>
              </div>
            </div>
          </div>
          <div v-if="hasMore" class="load-more">
            <el-button :loading="loading" @click="loadMore" class="load-more-btn">加载更多</el-button>
          </div>
          <div v-if="!hasMore && posts.length > 0" class="no-more">没有更多了</div>
        </template>
      </LoadingWrapper>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue"
import { toggleLike, listPosts, followingPosts, trendingPosts } from "@/api/post"
import { listFeaturedPosts } from "@/api/social"
import type { PostData } from "@/api/post"
import { usePagination } from "@/composables/usePagination"
import { useUserStore } from "@/stores/user"
import PostCard from "@/components/PostCard.vue"
import PageHeader from "@/components/PageHeader.vue"
import LoadingWrapper from "@/components/LoadingWrapper.vue"
import { Edit } from "@element-plus/icons-vue"

const userStore = useUserStore()
const isLogin = computed(() => userStore.isLogin)
const currentTab = ref("latest")

const emptyText = computed(() => {
  const map: Record<string, string> = {
    latest: "还没有动态",
    hot: "暂无热门内容",
    following: "关注的人还没有发布动态",
    featured: "暂无精选内容",
  }
  return map[currentTab.value] || "暂无动态"
})

function getFetchFn(tab: string) {
  switch (tab) {
    case "latest":
      return (offset: number, limit: number) => listPosts({ offset, limit })
    case "hot":
      return (offset: number, limit: number) => trendingPosts(7, limit)
    case "following":
      return (offset: number, limit: number) => followingPosts(offset, limit)
    case "featured":
      return (offset: number, limit: number) => listFeaturedPosts(offset, limit)
    default:
      return (offset: number, limit: number) => listPosts({ offset, limit })
  }
}

const fetchFn = (offset: number, limit: number) => getFetchFn(currentTab.value)(offset, limit)

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

  margin: 0 auto;
.square-page { max-width: 680px; margin: 0 auto; }

.tabs { margin-bottom: $space-5; }

.empty-state { padding: $space-12 0; }

.post-list { display: flex; flex-direction: column; gap: $space-3; }

.post-card-wrapper { position: relative; }

.trending-badge {
  position: absolute; top: 8px; right: 8px; z-index: 2;
  display: flex; align-items: center; gap: 2px;
  background: linear-gradient(135deg, #e74c3c, #f39c12);
  color: #fff; padding: 2px 8px; border-radius: 999px;
  font-size: 11px; font-weight: 600;
  .trending-label { opacity: 0.95; }
}

.load-more { text-align: center; padding: $space-6 0; }
.no-more { text-align: center; padding: $space-6 0; color: $text-muted; font-size: $text-sm; }

.load-more-btn {
  transition: all 0.3s ease;
  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba($brand-primary-hex, 0.15);
  }
}
.tabs {
  :deep(.el-radio-group) {
    transition: all 0.3s ease;
  }
}
</style>
