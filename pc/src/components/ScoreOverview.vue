<template>
  <div class="score-overview" v-if="totalRatings > 0">
    <div class="score-left">
      <div class="score-main">{{ avgScore.toFixed(1) }}</div>
      <el-rate :model-value="avgScore" disabled allow-half :max="maxStars" class="score-stars" />
      <div class="score-count">{{ totalRatings }} 条评价</div>
    </div>
    <div class="score-right">
      <div v-for="(count, i) in displayDistribution" :key="i" class="dist-row">
        <span class="dist-label">{{ maxStars - i }} 星</span>
        <el-progress
          :percentage="totalRatings > 0 ? (count / totalRatings) * 100 : 0"
          :stroke-width="8"
          class="dist-bar"
        />
        <span class="dist-count">{{ count }}</span>
      </div>
    </div>
  </div>
  <div v-else class="score-empty">
    <span>暂无评价</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"

const props = withDefaults(defineProps<{
  avgScore: number
  totalRatings: number
  distribution?: number[]
  maxStars?: number
}>(), {
  distribution: () => [],
  maxStars: 5,
})

const displayDistribution = computed(() => {
  if (props.distribution.length === props.maxStars) return props.distribution
  // Fallback: compute from avg if no distribution data
  const arr = new Array(props.maxStars).fill(0)
  if (props.totalRatings > 0) {
    arr[Math.floor(props.avgScore) - 1] = props.totalRatings
  }
  return arr
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.score-overview {
  display: flex; gap: 32px; align-items: flex-start;
  padding: 20px; background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; margin-bottom: 16px;
}

.score-left {
  text-align: center; min-width: 120px;
  .score-main {
    font-size: 3rem; font-weight: 700; color: $brand-primary; line-height: 1;
  }
  .score-stars { margin: 8px 0; justify-content: center; }
  .score-count { font-size: 12px; color: $text-muted; }
}

.score-right { flex: 1; display: flex; flex-direction: column; gap: 6px; }

.dist-row {
  display: flex; align-items: center; gap: 8px;
  .dist-label { font-size: 12px; color: $text-secondary; min-width: 30px; }
  .dist-bar { flex: 1; :deep(.el-progress-bar__outer) { background: $bg-hover; } }
  .dist-count { font-size: 12px; color: $text-muted; min-width: 20px; text-align: right; }
}

.score-empty {
  padding: 20px; background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; margin-bottom: 16px;
  text-align: center; font-size: 13px; color: $text-muted;
}
</style>
