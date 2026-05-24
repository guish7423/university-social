<template>
  <div class="activities-page">
    <div class="page-header">
      <h1>活动</h1>
      <el-button type="primary" @click="$router.push('/activity/create')">发起活动</el-button>
    </div>
    <div v-if="loading && !activities.length" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="!activities.length" class="empty-state"><el-empty description="暂无活动" /></div>
    <template v-else>
      <div class="activity-grid">
        <div v-for="a in activities" class="activity-card stagger-item" :key="a.id" @click="$router.push('/activities/' + a.id)">
          <div class="card-header">
            <el-tag size="small">{{ a.activity_type }}</el-tag>
            <span :class="['status', a.status === 0 ? 'open' : 'closed']">{{ a.status === 0 ? '报名中' : '已结束' }}</span>
          </div>
          <h3>{{ a.title }}</h3>
          <p class="desc">{{ a.description?.slice(0, 80) }}{{ (a.description?.length || 0) > 80 ? '...' : '' }}</p>
          <div class="card-meta">
            <span><el-icon><Timer /></el-icon> {{ formatDate(a.start_time) }}</span>
            <span><el-icon><User /></el-icon> {{ a.participant_count }}/{{ a.max_participants }}</span>
            <span v-if="a.location"><el-icon><Location /></el-icon> {{ a.location }}</span>
          </div>
        </div>
      </div>
      <div v-if="hasMore" class="load-more">
        <el-button :loading="loading" @click="loadMore" text>加载更多</el-button>
      </div>
      <div v-if="!hasMore && activities.length > 0" class="no-more">没有更多了</div>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { ActivityData } from "@/api/activity"
import { usePagination } from "@/composables/usePagination"
import { useTimeFormat } from "@/composables/useTimeFormat"
import { listActivities } from "@/api/activity"

const { formatDate } = useTimeFormat()
const { items: activities, loading, hasMore, loadMore } = usePagination<ActivityData>({
  fetchFn: (offset, limit) => listActivities({ offset, limit }),
})

loadMore()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.activities-page { max-width: 780px; }
.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px;
  h1 { font-size: 22px; font-weight: 700; }
}

.activity-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 14px; }

.activity-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 18px; cursor: pointer; transition: all 0.2s;
  &:hover { border-color: $primary-light; transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }

  .card-header {
    display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px;
  }
  .status.open { color: $accent-green; font-size: 12px; }
  .status.closed { color: $text-muted; font-size: 12px; }

  h3 { font-size: 16px; font-weight: 600; margin-bottom: 6px; }
  .desc { font-size: 13px; color: $text-secondary; margin-bottom: 10px; }

  .card-meta {
    display: flex; flex-wrap: wrap; gap: 12px; font-size: 12px; color: $text-muted;
    .el-icon { vertical-align: middle; margin-right: 3px; }
  }
}
.load-more { text-align: center; padding: 20px 0; }
.no-more { text-align: center; padding: 20px 0; color: $text-muted; font-size: 13px; }
</style>
