<template>
  <div class="search-page">
    <div class="page-header"><h1>搜索</h1></div>
    <el-input v-model="query" placeholder="搜索用户、帖子、圈子..." size="large" clearable @keyup.enter="handleSearch">
      <template #prefix><el-icon><Search /></el-icon></template>
    </el-input>

    <div v-if="searched">
      <div v-if="searchLoading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
      <template v-else>
        <div v-if="users.length" class="section">
          <h3>用户</h3>
          <div class="user-list">
            <div v-for="u in users" :key="u.id" class="user-row" @click="$router.push('/user/' + u.id)">
              <el-avatar :size="36" :src="u.avatar" />
              <span class="name">{{ u.nickname }}</span>
              <el-tag v-if="u.friend_status === 'accepted'" size="small" type="success">好友</el-tag>
              <el-tag v-else-if="u.friend_status === 'pending'" size="small" type="warning">待处理</el-tag>
            </div>
          </div>
        </div>
        <div v-if="posts.length" class="section">
          <h3>帖子</h3>
          <div class="post-list">
            <div v-for="p in posts" :key="p.id" class="post-item" @click="$router.push('/post/' + p.id)">
              <p>{{ p.content.slice(0, 100) }}</p>
              <span class="meta">{{ p.like_count }} 赞 · {{ p.comment_count }} 评论</span>
            </div>
          </div>
        </div>
        <div v-if="circles.length" class="section">
          <h3>圈子</h3>
          <div class="circle-list">
            <div v-for="c in circles" :key="c.id" class="circle-item" @click="$router.push('/circles/' + c.id)">
              <el-avatar :size="32" :src="c.avatar" />
              <span>{{ c.name }}</span>
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
import { ref } from "vue"
import { searchUsers } from "@/api/social"
import { searchPosts } from "@/api/post"
import { searchCircles } from "@/api/circle"
import type { PostData } from "@/api/post"
import type { CircleData } from "@/api/circle"
import type { UserInfo } from "@/api/auth"

const query = ref("")
const searched = ref(false)
const searchLoading = ref(false)
const users = ref<(UserInfo & { friend_status: string })[]>([])
const posts = ref<PostData[]>([])
const circles = ref<CircleData[]>([])

async function handleSearch() {
  if (!query.value.trim()) return
  searched.value = true
  searchLoading.value = true
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

.section { margin-top: 20px;
  h3 { font-size: 15px; font-weight: 600; margin-bottom: 10px; }
}

.user-list { display: flex; flex-direction: column; gap: 6px; }
.user-row {
  display: flex; align-items: center; gap: 10px;
  padding: 10px; border-radius: $radius-sm; cursor: pointer;
  &:hover { background: $bg-card; }
  .name { flex: 1; font-size: 14px; font-weight: 600; }
}

.post-item {
  padding: 10px; border-radius: $radius-sm; cursor: pointer; margin-bottom: 4px;
  &:hover { background: $bg-card; }
  p { font-size: 14px; margin-bottom: 4px; }
  .meta { font-size: 12px; color: $text-muted; }
}

.circle-item {
  display: flex; align-items: center; gap: 10px;
  padding: 10px; border-radius: $radius-sm; cursor: pointer;
  &:hover { background: $bg-card; }
  .meta { margin-left: auto; font-size: 12px; color: $text-muted; }
}
</style>
