<template>
  <div class="circles-page">
    <div class="page-header">
      <h1>圈子</h1>
      <el-button type="primary" @click="$router.push('/circle/create')">创建圈子</el-button>
    </div>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="!circles.length" class="empty-state"><el-empty description="暂无圈子" /></div>
    <div v-else class="circle-grid">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listCircles } from "@/api/circle"
import type { CircleData } from "@/api/circle"

const circles = ref<CircleData[]>([])
const loading = ref(true)

onMounted(async () => {
  try { circles.value = await listCircles() }
  catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.circles-page { max-width: 700px; }
.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px;
  h1 { font-size: 22px; font-weight: 700; }
}

.circle-grid { display: flex; flex-direction: column; gap: 12px; }

.circle-card {
  display: flex; align-items: center; gap: 16px;
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 18px; cursor: pointer; transition: all 0.2s;
  &:hover { border-color: $primary-light; }
  .circle-info {
    flex: 1; min-width: 0;
    .name { font-size: 16px; font-weight: 600; }
    .desc { font-size: 13px; color: $text-secondary; margin: 4px 0; }
    .meta { font-size: 12px; color: $text-muted; }
  }
}
</style>
