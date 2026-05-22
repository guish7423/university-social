<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
  <div v-else-if="!profile" class="empty-state"><el-empty description="用户不存在" /></div>
  <div v-else class="user-page">
    <div class="user-banner">
      <el-avatar :size="64" :src="profile.avatar" />
      <div class="user-info">
        <h1>{{ profile.nickname }}</h1>
        <div class="meta">
          <span>{{ profile.school || '未设置学校' }}</span>
          <el-tag v-if="profile.is_verified" size="small" type="success">已认证</el-tag>
        </div>
      </div>
      <div v-if="profile.id !== userStore.userId" class="actions">
        <el-button @click="$router.push('/chat/' + profile.id)">发消息</el-button>
        <el-button :type="following ? 'default' : 'primary'" @click="handleFollow">
          {{ following ? '已关注' : '关注' }}
        </el-button>
      </div>
    </div>

    <div class="nav-tabs">
      <el-button :type="tab === 'posts' ? 'primary' : 'default'" @click="tab = 'posts'">他的帖子</el-button>
      <el-button :type="tab === 'followers' ? 'primary' : 'default'" @click="$router.push(`/user/${profile.id}/followers`)">粉丝</el-button>
      <el-button :type="tab === 'following' ? 'primary' : 'default'" @click="$router.push(`/user/${profile.id}/following`)">关注</el-button>
    </div>

    <div v-if="tab === 'posts'" class="section">
      <div v-if="userPostsLoading" class="loading-wrap"><el-skeleton :rows="3" animated /></div>
      <div v-else-if="!userPostsData.length" class="empty-state"><el-empty description="暂无帖子" /></div>
      <div v-else class="post-list">
        <div v-for="p in userPostsData" :key="p.id" class="post-card" @click="$router.push('/post/' + p.id)">
          <p>{{ p.content.slice(0, 120) }}</p>
          <div class="meta">
            <span>{{ p.like_count }} 赞</span>
            <span>{{ p.comment_count }} 评论</span>
            <span>{{ formatTime(p.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"
import { useUserStore } from "@/stores/user"
import { getUserInfo } from "@/api/auth"
import { userPosts } from "@/api/post"
import { followUser, unfollowUser, isFollowing } from "@/api/follow"
import { useTimeFormat } from "@/composables/useTimeFormat"
import type { UserInfo } from "@/api/auth"
import type { PostData } from "@/api/post"
import { ElMessage } from "element-plus"

const route = useRoute()
const userStore = useUserStore()
const profile = ref<UserInfo | null>(null)
const userPostsData = ref<PostData[]>([])
const loading = ref(true)
const userPostsLoading = ref(false)
const following = ref(false)
const tab = ref("posts")

const { formatTime } = useTimeFormat()

async function handleFollow() {
  if (!profile.value) return
  try {
    if (following.value) {
      await unfollowUser(profile.value.id)
      following.value = false
      ElMessage.success("已取消关注")
    } else {
      await followUser(profile.value.id)
      following.value = true
      ElMessage.success("关注成功")
    }
  } catch { /* handled by interceptor */ }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try { profile.value = await getUserInfo(id) }
  catch { /* handled by interceptor */ } finally { loading.value = false }

  if (id !== userStore.userId) {
    try {
      const res = await isFollowing(id)
      following.value = res.following
    } catch { /* handled by interceptor */ }
  }

  try { userPostsData.value = await userPosts(id) }
  catch { /* handled by interceptor */ } finally { userPostsLoading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.user-page { max-width: 680px; }

.user-banner {
  display: flex; align-items: center; gap: 16px;
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 24px; margin-bottom: 20px;
  .user-info { flex: 1;
    h1 { font-size: 20px; font-weight: 700; }
    .meta { display: flex; align-items: center; gap: 8px; font-size: 13px; color: $text-secondary; margin-top: 4px; }
  }
}

.nav-tabs { display: flex; gap: 8px; margin-bottom: 16px; }

.post-list { display: flex; flex-direction: column; gap: 8px; }
.post-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 14px; cursor: pointer;
  p { font-size: 14px; line-height: 1.6; margin-bottom: 8px; }
  .meta { display: flex; gap: 12px; font-size: 12px; color: $text-muted; }
  &:hover { border-color: $primary-light; }
}
</style>
