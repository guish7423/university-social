<template>
  <div class="circles-page">
    <div class="page-header">
      <h1>圈子</h1>
      <el-button type="primary" @click="$router.push('/circle/create')">创建圈子</el-button>
    </div>
    <LoadingWrapper :loading="loading && !circles.length" :data="circles.length" skeleton-variant="post-card">
      <template #empty>
        <el-empty description="暂无圈子" />
      </template>
      <div class="circle-grid">
        <div v-for="c in circles" class="circle-card stagger-item" :key="c.id" @click="$router.push('/circles/' + c.id)">
          <el-avatar :size="52" :src="c.avatar" />
          <div class="circle-info">
            <div class="name">{{ c.name }}</div>
            <div class="desc">{{ c.description?.slice(0, 60) }}{{ (c.description?.length || 0) > 60 ? '...' : '' }}</div>
            <div class="meta">{{ c.member_count }} 人 · {{ c.post_count }} 帖</div>
          </div>
          <el-tag v-if="c.is_joined" size="small" type="success">已加入</el-tag>
        </div>
      </div>
    </LoadingWrapper>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listCircles } from "@/api/circle"
import type { CircleData } from "@/api/circle"
import LoadingWrapper from "@/components/LoadingWrapper.vue"

const circles = ref<CircleData[]>([])
const loading = ref(true)

onMounted(async () => {
  try { circles.value = await listCircles() }
  catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

 margin: 0 auto;
.circles-page { max-width: 700px; }
</style>
