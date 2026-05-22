<template>
  <div class="found-page">
    <div class="page-header">
      <h1>失物招领</h1>
      <el-button type="primary" @click="$router.push('/found/create')">发布信息</el-button>
    </div>
    <div v-if="loading && !items.length" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="!items.length" class="empty-state"><el-empty description="暂无失物信息" /></div>
    <template v-else>
      <div class="found-grid">
        <div v-for="item in items" class="found-card stagger-item" :key="item.id" @click="$router.push('/found/' + item.id)">
          <div class="card-header">
            <el-tag :type="item.category === '寻物' ? 'warning' : 'success'" size="small">{{ item.category }}</el-tag>
            <span :class="['status', item.status === 0 ? 'active' : 'resolved']">{{ item.status === 0 ? '进行中' : '已解决' }}</span>
          </div>
          <h3>{{ item.title }}</h3>
          <p class="desc">{{ item.description?.slice(0, 80) }}{{ item.description?.length > 80 ? '...' : '' }}</p>
          <div class="meta">
            <span v-if="item.location"><el-icon><Location /></el-icon> {{ item.location }}</span>
            <span class="time">{{ formatTime(item.created_at) }}</span>
          </div>
        </div>
      </div>
      <div v-if="hasMore" class="load-more">
        <el-button :loading="loading" @click="loadMore" text>加载更多</el-button>
      </div>
      <div v-if="!hasMore && items.length > 0" class="no-more">没有更多了</div>
</template>

  </div>
</template>


<script setup lang="ts">
import type { LostItemData } from "@/api/found"
import { usePagination } from "@/composables/usePagination"
import { useTimeFormat } from "@/composables/useTimeFormat"
import { listLostItems } from "@/api/found"

const { formatTime } = useTimeFormat()
const { items, loading, hasMore, loadMore } = usePagination<LostItemData>({
  fetchFn: (offset, limit) => listLostItems({ offset, limit }),
})

loadMore()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.found-page { max-width: 780px; }
.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px;
  h1 { font-size: 22px; font-weight: 700; }
}
.found-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 14px; }
.found-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px; cursor: pointer; transition: all 0.2s;
  &:hover { border-color: $primary-light; }
  .card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
  .status.active { color: $accent-green; font-size: 12px; }
  .status.resolved { color: $text-muted; font-size: 12px; }
  h3 { font-size: 15px; font-weight: 600; margin-bottom: 6px; }
  .desc { font-size: 13px; color: $text-secondary; margin-bottom: 8px; }
  .meta { display: flex; justify-content: space-between; font-size: 12px; color: $text-muted; .el-icon { vertical-align: middle; margin-right: 3px; } }
}
.load-more { text-align: center; padding: 20px 0; }
.no-more { text-align: center; padding: 20px 0; color: $text-muted; font-size: 13px; }
</style>
