<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
  <div v-else-if="!profile" class="empty-state"><el-empty description="用户不存在" /></div>
  <div v-else class="user-page">
    <div class="profile-banner">
      <div class="cover-area" />
      <div class="profile-content">
        <el-avatar :size="72" :src="profile.avatar" class="profile-avatar">
          {{ profile.nickname[0] }}
        </el-avatar>
        <div class="profile-info">
          <h1 class="nickname">{{ profile.nickname }}</h1>
          <p class="school">{{ profile.school || '未设置学校' }}</p>
          <p class="bio" v-if="false" />
          <p class="join-date">加入于 {{ formatTime(profile.created_at) }}</p>
        </div>
        <div v-if="profile.id !== userStore.userId" class="profile-actions">
          <el-button round @click="$router.push('/chat/' + profile.id)">发消息</el-button>
          <el-button :type="following ? 'default' : 'primary'" round @click="handleFollow">
            {{ following ? '已关注' : '关注' }}
          </el-button>
        </div>
      </div>
    </div>

    <div class="stats-row">
      <div class="stat-item">
        <span class="stat-value">{{ userPostsData.length }}</span>
        <span class="stat-label">帖子</span>
      </div>
      <div class="stat-item">
        <span class="stat-value">{{ counts?.followers ?? '-' }}</span>
        <span class="stat-label">粉丝</span>
      </div>
      <div class="stat-item">
        <span class="stat-value">{{ counts?.following ?? '-' }}</span>
        <span class="stat-label">关注</span>
      </div>
    </div>

    <div class="nav-tabs">
      <el-button :type="tab === 'posts' ? 'primary' : 'default'" @click="tab = 'posts'" round size="small">帖子</el-button>
      <el-button :type="tab === 'followers' ? 'primary' : 'default'" @click="tab = 'followers'" round size="small">粉丝</el-button>
      <el-button :type="tab === 'following' ? 'primary' : 'default'" @click="tab = 'following'" round size="small">关注</el-button>
    </div>

    <div v-if="tab === 'posts'" class="section">
      <div v-if="userPostsLoading" class="loading-wrap"><el-skeleton :rows="3" animated /></div>
      <div v-else-if="!userPostsData.length" class="empty-state"><el-empty description="暂无帖子" /></div>
      <div v-else class="post-list">
        <PostCard v-for="p in userPostsData" :key="p.id" :post="p" @click="$router.push('/post/' + p.id)" />
      </div>
    </div>

    <div v-else-if="tab === 'followers'" class="section">
      <div v-if="followersLoading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
      <div v-else-if="!followersData.length" class="empty-state"><el-empty description="暂无粉丝" /></div>
      <div v-else class="user-list">
        <div v-for="f in followersData" :key="f.id" class="user-row" @click="$router.push('/user/' + f.id)">
          <el-avatar :size="36" :src="f.avatar">{{ f.nickname[0] }}</el-avatar>
          <span class="user-name">{{ f.nickname }}</span>
        </div>
      </div>
    </div>

    <div v-else class="section">
      <div v-if="followingLoading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
      <div v-else-if="!followingData.length" class="empty-state"><el-empty description="暂未关注任何人" /></div>
      <div v-else class="user-list">
        <div v-for="f in followingData" :key="f.id" class="user-row" @click="$router.push('/user/' + f.id)">
          <el-avatar :size="36" :src="f.avatar">{{ f.nickname[0] }}</el-avatar>
          <span class="user-name">{{ f.nickname }}</span>
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
import { followUser, unfollowUser, isFollowing, followCounts, followersList, followingList } from "@/api/follow"
import { useTimeFormat } from "@/composables/useTimeFormat"
import PostCard from "@/components/PostCard.vue"
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
const counts = ref<{ following: number; followers: number } | null>(null)
const followersData = ref<{ id: number; nickname: string; avatar: string }[]>([])
const followingData = ref<{ id: number; nickname: string; avatar: string }[]>([])
const followersLoading = ref(false)
const followingLoading = ref(false)

const { formatTime } = useTimeFormat()

async function handleFollow() {
  if (!profile.value) return
  try {
    if (following.value) {
      await unfollowUser(profile.value.id)
      following.value = false
      if (counts.value) counts.value.followers = Math.max(0, counts.value.followers - 1)
      ElMessage.success("已取消关注")
    } else {
      await followUser(profile.value.id)
      following.value = true
      if (counts.value) counts.value.followers++
      ElMessage.success("关注成功")
    }
  } catch { /* handled by interceptor */ }
}

async function loadFollowers() {
  if (!profile.value) return
  followersLoading.value = true
  try {
    followersData.value = await followersList(profile.value.id)
  } catch { /* handled */ } finally { followersLoading.value = false }
}

async function loadFollowing() {
  if (!profile.value) return
  followingLoading.value = true
  try {
    followingData.value = await followingList(profile.value.id)
  } catch { /* handled */ } finally { followingLoading.value = false }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    profile.value = await getUserInfo(id)
    counts.value = await followCounts(id)
  } catch { /* handled */ } finally { loading.value = false }

  if (id !== userStore.userId) {
    try {
      const res = await isFollowing(id)
      following.value = res.following
    } catch { /* handled */ }
  }

  try { userPostsData.value = await userPosts(id) }
  catch { /* handled */ } finally { userPostsLoading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.user-page { max-width: 720px; margin: 0 auto; }

.profile-banner {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-lg; overflow: hidden; margin-bottom: 20px;
  .cover-area {
    height: 120px;
    background: linear-gradient(135deg, $primary-light 0%, $primary 50%, #7a4a3e 100%);
  }
  .profile-content {
    display: flex; align-items: flex-start; gap: 16px;
    padding: 0 24px 20px; margin-top: -36px;
    .profile-avatar {
      border: 3px solid $bg-card; flex-shrink: 0;
    }
    .profile-info {
      flex: 1; min-width: 0; padding-top: 40px;
      .nickname { font-size: 22px; font-weight: 700; margin: 0 0 2px; }
      .school { font-size: 13px; color: $text-secondary; margin: 0 0 2px; }
      .join-date { font-size: 12px; color: $text-muted; margin: 0; }
    }
    .profile-actions {
      display: flex; gap: 8px; padding-top: 40px; flex-shrink: 0;
    }
  }
}

.stats-row {
  display: flex; gap: 1px; background: $border-color;
  border-radius: $radius-md; overflow: hidden; margin-bottom: 16px;
  .stat-item {
    flex: 1; background: $bg-card; padding: 14px; text-align: center;
    .stat-value { display: block; font-size: 20px; font-weight: 700; color: $primary; }
    .stat-label { display: block; font-size: 12px; color: $text-secondary; margin-top: 2px; }
  }
}

.nav-tabs { display: flex; gap: 8px; margin-bottom: 16px; }

.post-list { display: flex; flex-direction: column; gap: 8px; }

.user-list {
  display: flex; flex-direction: column; gap: 1px;
  background: $border-color; border-radius: $radius-md; overflow: hidden;
  .user-row {
    display: flex; align-items: center; gap: 10px;
    background: $bg-card; padding: 10px 14px; cursor: pointer;
    .user-name { font-size: 14px; font-weight: 500; }
    &:hover { background: $bg-hover; }
  }
}

.loading-wrap { padding: 40px 0; }
.empty-state { padding: 60px 0; }
</style>
