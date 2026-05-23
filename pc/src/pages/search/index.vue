<template>
  <div class="search-page">
    <div class="page-header">
      <h1>搜索</h1>
    </div>
    <el-input
      v-model="query" ref="inputRef"
      placeholder="搜索用户、帖子、圈子..." size="large" clearable
      @keyup.enter="handleSearch"
    >
      <template #prefix><el-icon><Search /></el-icon></template>
    </el-input>

    <div v-if="query && !searched" class="hint">
      按 Enter 搜索
    </div>

    <div v-if="searched">
      <div v-if="searchLoading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
      <template v-else>
        <div class="tabs">
          <el-button :type="tab === 'all' ? 'primary' : 'default'" size="small" round @click="tab = 'all'">
            全部 <span v-if="totalCount">{{ totalCount }}</span>
          </el-button>
          <el-button :type="tab === 'users' ? 'primary' : 'default'" size="small" round @click="tab = 'users'">
            用户 <span v-if="users.length">{{ users.length }}</span>
          </el-button>
          <el-button :type="tab === 'posts' ? 'primary' : 'default'" size="small" round @click="tab = 'posts'">
            帖子 <span v-if="posts.length">{{ posts.length }}</span>
          </el-button>
          <el-button :type="tab === 'circles' ? 'primary' : 'default'" size="small" round @click="tab = 'circles'">
            圈子 <span v-if="circles.length">{{ circles.length }}</span>
          </el-button>
        </div>

        <div v-if="tab === 'users' || tab === 'all'" class="section">
          <h3 v-if="tab === 'all'">用户 <small v-if="users.length">({{ users.length }})</small></h3>
          <div v-if="!users.length && tab !== 'all'" class="empty-state"><el-empty description="未找到用户" /></div>
          <div v-else class="user-list">
            <div v-for="u in users" :key="u.id" class="user-row" @click="$router.push('/user/' + u.id)">
              <el-avatar :size="36" :src="u.avatar">{{ u.nickname[0] }}</el-avatar>
              <span class="name">{{ u.nickname }}</span>
              <el-tag v-if="u.friend_status === 'accepted'" size="small" type="success">好友</el-tag>
              <el-tag v-else-if="u.friend_status === 'pending'" size="small" type="warning">待处理</el-tag>
            </div>
          </div>
        </div>

        <div v-if="tab === 'posts' || tab === 'all'" class="section">
          <h3 v-if="tab === 'all'">帖子 <small v-if="posts.length">({{ posts.length }})</small></h3>
          <div v-if="!posts.length && tab !== 'all'" class="empty-state"><el-empty description="未找到帖子" /></div>
          <div v-else class="post-list">
            <PostCard v-for="p in posts" :key="p.id" :post="p" @click="$router.push('/post/' + p.id)" />
          </div>
        </div>

        <div v-if="tab === 'circles' || tab === 'all'" class="section">
          <h3 v-if="tab === 'all'">圈子 <small v-if="circles.length">({{ circles.length }})</small></h3>
          <div v-if="!circles.length && tab !== 'all'" class="empty-state"><el-empty description="未找到圈子" /></div>
          <div v-else class="circle-list">
            <div v-for="c in circles" :key="c.id" class="circle-item" @click="$router.push('/circles/' + c.id)">
              <el-avatar :size="32" :src="c.avatar">{{ c.name[0] }}</el-avatar>
              <span class="name">{{ c.name }}</span>
              <span class="meta">{{ c.member_count }} 人</span>
            </div>
          </div>
        </div>

        <div v-if="!users.length && !posts.length && !circles.length" class="empty-state">
          <el-empty description="未找到结果" />
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue"
import { searchUsers } from "@/api/social"
import { searchPosts } from "@/api/post"
import { searchCircles } from "@/api/circle"
import PostCard from "@/components/PostCard.vue"
import type { PostData } from "@/api/post"
import type { CircleData } from "@/api/circle"
import type { UserInfo } from "@/api/auth"

const query = ref("")
const searched = ref(false)
const searchLoading = ref(false)
const tab = ref("all")
const users = ref<(UserInfo & { friend_status: string })[]>([])
const posts = ref<PostData[]>([])
const circles = ref<CircleData[]>([])

const totalCount = computed(() => users.value.length + posts.value.length + circles.value.length)

async function handleSearch() {
  if (!query.value.trim()) return
  searched.value = true
  searchLoading.value = true
  users.value = []
  posts.value = []
  circles.value = []
  try {
    const [uRes, pRes, cRes] = await Promise.allSettled([
      searchUsers(query.value.trim()),
      searchPosts(query.value.trim()),
      searchCircles(query.value.trim()),
    ])
    if (uRes.status === "fulfilled") users.value = uRes.value
    if (pRes.status === "fulfilled") posts.value = pRes.value
    if (cRes.status === "fulfilled") circles.value = cRes.value
  } finally { searchLoading.value = false }
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.search-page { max-width: 640px; }

.page-header { margin-bottom: 16px; h1 { font-size: 22px; font-weight: 700; } }

.hint { text-align: center; padding: 40px 0; color: $text-muted; font-size: 14px; }

.tabs { display: flex; gap: 8px; margin: 16px 0;
  span { font-size: 12px; opacity: 0.7; margin-left: 2px; }
}

.section {
  margin-top: 20px;
  h3 { font-size: 15px; font-weight: 600; margin-bottom: 10px; }
}

.user-list { display: flex; flex-direction: column; gap: 6px; }
.user-row {
  display: flex; align-items: center; gap: 10px;
  padding: 10px; border-radius: $radius-sm; cursor: pointer;
  &:hover { background: $bg-hover; }
  .name { flex: 1; font-size: 14px; font-weight: 600; }
}

.post-list { display: flex; flex-direction: column; gap: 6px; }

.circle-item {
  display: flex; align-items: center; gap: 10px;
  padding: 10px; border-radius: $radius-sm; cursor: pointer;
  &:hover { background: $bg-hover; }
  .name { flex: 1; font-size: 14px; }
  .meta { font-size: 12px; color: $text-muted; }
}

.loading-wrap { padding: 40px 0; }
</style>
