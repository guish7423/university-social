<template>
  <div class="review-section">
    <ScoreOverview
      :avg-score="avgScore"
      :total-ratings="ratings.length"
      :max-stars="5"
    />

    <div class="rating-actions" v-if="!myRating">
      <el-button type="primary" size="large" @click="$emit('openDialog')">
        写评价
      </el-button>
    </div>
    <div class="rating-actions" v-else>
      <el-tag type="success" size="large" effect="plain">
        已评价 {{ myRating.score }} 星
      </el-tag>
    </div>

    <div v-if="loading" class="loading-wrap">
      <el-skeleton :rows="3" animated />
    </div>
    <div v-else-if="ratings.length === 0" class="empty-state">
      <el-empty :description="myRating ? '暂无其他评价' : '还没有评价，来说两句？'" />
    </div>
    <div v-else class="review-list">
      <div v-for="r in sortedRatings" :key="r.id" class="review-card stagger-item">
        <div class="review-header">
          <div class="review-author">
            <span class="review-name">{{ r.author?.nickname || '匿名用户' }}</span>
            <el-rate :model-value="r.rating" disabled :max="5" class="review-stars" />
          </div>
          <span class="review-time">{{ formatTime(r.created_at) }}</span>
        </div>
        <p class="review-text" v-if="r.comment">{{ r.comment }}</p>
      </div>
    </div>

    <RatingDialog
      :visible="ratingDialogVisible"
      @close="$emit('closeDialog')"
      @submit="(data) => $emit('submit', data)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import ScoreOverview from "./ScoreOverview.vue"
import RatingDialog from "./RatingDialog.vue"

const props = defineProps<{
  ratings: any[]
  loading: boolean
  myRating?: any | null
  ratingDialogVisible: boolean
}>()

defineEmits<{
  (e: "submit", data: { score: number; comment: string }): void
  (e: "openDialog"): void
  (e: "closeDialog"): void
}>()

const avgScore = computed(() => {
  if (props.ratings.length === 0) return 0
  return props.ratings.reduce((sum, r) => sum + r.rating, 0) / props.ratings.length
})

const sortedRatings = computed(() => {
  return [...props.ratings].sort((a, b) =>
    new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )
})

function formatTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 60000) return "刚刚"
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  if (d.getDate() === now.getDate() && d.getMonth() === now.getMonth()) return "今天"
  return `${d.getMonth() + 1}-${d.getDate()}`
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.rating-actions {
  margin-bottom: 16px;
}

.review-list { display: flex; flex-direction: column; gap: 10px; }

.review-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 16px;
}

.review-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;
}

.review-author {
  display: flex; align-items: center; gap: 10px;
  .review-name { font-size: 13px; font-weight: 600; }
  .review-stars { :deep(.el-rate__icon) { font-size: 14px; } }
}

.review-time { font-size: 12px; color: $text-muted; }

.review-text {
  font-size: 14px; line-height: 1.7; color: $text-primary; margin: 0;
}
</style>
