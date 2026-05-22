<template>
  <div class="list-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>
    <h1>粉丝</h1>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="!list.length" class="empty-state"><el-empty description="暂无粉丝" /></div>
    <div v-else class="user-list">
      <div v-for="u in list" :key="u.id" class="user-row" @click="$router.push('/user/' + u.id)">
        <el-avatar :size="36" :src="u.avatar" />
        <span class="name">{{ u.nickname }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"
import { getFollowers } from "@/api/social"
import type { UserInfo } from "@/api/auth"
import { ArrowLeft } from "@element-plus/icons-vue"

const route = useRoute()
const list = ref<UserInfo[]>([])
const loading = ref(true)
const userId = Number(route.params.id)

onMounted(async () => {
  try {
    const res = await getFollowers(userId)
    list.value = res
  } catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.list-page { max-width: 560px; }
h1 { font-size: 20px; font-weight: 700; margin: 16px 0; }
.user-list { display: flex; flex-direction: column; gap: 6px; }
.user-row {
  display: flex; align-items: center; gap: 10px;
  padding: 10px; border-radius: $radius-sm; cursor: pointer;
  &:hover { background: $bg-card; }
  .name { font-size: 14px; font-weight: 600; }
}
</style>
