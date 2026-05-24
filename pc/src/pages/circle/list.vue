<template>
  <div class="circles-page">
    <div class="page-header">
      <h1>圈子</h1>
      <el-button type="primary" @click="$router.push('/circle/create')">创建圈子</el-button>
    </div>

    <!-- Tabs + Search -->
    <div class="toolbar">
      <div class="tabs">
        <button v-for="tab in tabs" :key="tab.key" :class="['tab-btn', { active: activeTab === tab.key }]" @click="activeTab = tab.key">
          {{ tab.label }}
        </button>
      </div>
      <el-input v-model="searchQuery" placeholder="搜索圈子..." size="small" clearable class="search-input" />
    </div>

    <LoadingWrapper :loading="loading && !circles.length" :data="filtered.length" skeleton-variant="post-card">
      <template #empty>
        <el-empty description="暂无圈子" />
      </template>
      <div class="circle-grid">
        <div v-for="c in filtered" class="circle-card stagger-item" :key="c.id" @click="$router.push('/circles/' + c.id)">
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
import { ref, computed, onMounted } from "vue"
import { listCircles } from "@/api/circle"
import type { CircleData } from "@/api/circle"
import LoadingWrapper from "@/components/LoadingWrapper.vue"

const circles = ref<CircleData[]>([])
const loading = ref(true)
const activeTab = ref('all')
const searchQuery = ref('')

const tabs = [
  { key: 'all', label: '所有圈子' },
  { key: 'joined', label: '已加入' },
  { key: 'popular', label: '热门' },
]

const filtered = computed(() => {
  let list = circles.value
  if (activeTab.value === 'joined') {
    list = list.filter(c => c.is_joined)
  } else if (activeTab.value === 'popular') {
    list = [...list].sort((a, b) => (b.member_count || 0) - (a.member_count || 0))
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(c => c.name?.toLowerCase().includes(q) || c.description?.toLowerCase().includes(q))
  }
  return list
})

onMounted(async () => {
  try { circles.value = await listCircles() }
  catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.circles-page { max-width: 700px; }
.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;
  h1 { font-size: 22px; font-weight: 700; }
}
.toolbar {
  display: flex; align-items: center; gap: 12px; margin-bottom: 16px;
  .tabs { display: flex; gap: 4px; }
  .tab-btn {
    padding: 6px 16px; border-radius: 20px; border: none; font-size: 13px; cursor: pointer;
    background: transparent; color: $text-secondary; transition: all 0.2s;
    &:hover { background: color-mix(in oklch, $primary 8%, transparent); color: $primary; }
    &.active { background: $primary; color: #fff; }
  }
  .search-input { width: 200px; }
}
.circle-grid {
  display: grid; grid-template-columns: 1fr; gap: 10px;
}
.circle-card {
  display: flex; align-items: center; gap: 14px;
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px; cursor: pointer; transition: all 0.2s;
  &:hover { border-color: $primary-light; transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.12); }
  .circle-info { flex: 1; min-width: 0; }
  .name { font-size: 15px; font-weight: 600; margin-bottom: 4px; }
  .desc { font-size: 13px; color: $text-secondary; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; margin-bottom: 4px; }
  .meta { font-size: 12px; color: $text-muted; }
}
</style>
