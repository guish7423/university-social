<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="8" animated /></div>
  <div v-else-if="!target" class="empty-state"><el-empty description="评价对象不存在" /></div>
  <div v-else class="detail-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn stagger-item">返回</el-button>

    <!-- Score Overview -->
    <div class="score-section stagger-item">
      <div class="score-left">
        <div class="score-big" :class="scoreClass(target.avgScore)">{{ target.avgScore.toFixed(1) }}</div>
        <div class="score-label">{{ target.totalRatings }} 人评价</div>
      </div>
      <div class="score-right">
        <div class="dist-row" v-for="(val, i) in distData" :key="i">
          <span class="dist-label">{{ distLabels[i] }}</span>
          <div class="dist-bar-bg">
            <div class="dist-bar-fill" :style="{ width: pct(val, target.totalRatings) }" />
          </div>
          <span class="dist-val">{{ val }}</span>
        </div>
      </div>
    </div>

    <div class="meta-row stagger-item">
      <span class="meta-cat">{{ categoryLabel(target.category) }}</span>
      <span class="meta-desc">{{ target.description }}</span>
      <span class="meta-by">发起人：{{ target.creator.nickname }}</span>
    </div>

    <!-- Rate Button -->
    <div class="rate-section stagger-item">
      <el-button v-if="!userRating" type="primary" size="large" @click="showRatingDialog = true">
        <el-icon><Star /></el-icon> 写评价
      </el-button>
      <el-button v-else text type="warning" @click="showRatingDialog = true">
        <el-icon><StarFilled /></el-icon> 已评价 {{ userRating.score }} 分 · 编辑
      </el-button>
    </div>

    <!-- Rating List -->
    <div class="filter-bar">
      <el-radio-group v-model="ratingSort" @change="handleSortChange" size="small">
        <el-radio-button value="hot">热门</el-radio-button>
        <el-radio-button value="latest">最新</el-radio-button>
        <el-radio-button value="high">高分</el-radio-button>
      </el-radio-group>
    </div>

    <div v-if="ratingsLoading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="!ratings.length" class="empty-card">
      <el-empty description="还没有评价，来说两句？" :image-size="60" />
    </div>
    <div v-else class="ratings-list">
      <div v-for="r in ratings" :key="r.id" class="rating-card stagger-item">
        <div class="rc-header">
          <el-avatar :size="32" :src="r.user.avatar">{{ r.user.nickname[0] }}</el-avatar>
          <span class="rc-name">{{ r.user.nickname }}</span>
          <span class="rc-score" :class="scoreClass(r.score)">{{ r.score }}</span>
          <span class="rc-time">{{ smartTime(r.createdAt) }}</span>
        </div>
        <div v-if="r.tags?.length" class="rc-tags">
          <span v-for="tag in r.tags" :key="tag" class="rc-tag">{{ tag }}</span>
        </div>
        <p v-if="r.content" class="rc-content">{{ r.content }}</p>
        <div class="rc-actions">
          <el-button :type="r.isHelpfulByMe ? 'warning' : 'default'" text size="small" @click="handleHelpful(r)">
            <el-icon><ThumbsUp /></el-icon> {{ r.helpfulCount || '' }} 有用
          </el-button>
        </div>
      </div>
    </div>

    <!-- Rating Dialog -->
    <el-dialog v-model="showRatingDialog" title="写评价" width="500px" top="5vh">
      <div class="rate-dialog-body">
        <div class="rate-stars">
          <span class="rate-label">评分</span>
          <el-rate v-model="newRating.score" :max="10" show-score :colors="['#e74c3c','#d4a85e','#67c23a']" score-template="{value} 分" />
        </div>
        <div v-if="tagOptions.length" class="rate-tags">
          <span class="rate-label">标签（选填，最多 3 个）</span>
          <el-checkbox-group v-model="newRating.tags" :max="3">
            <el-checkbox v-for="tag in tagOptions" :key="tag" :label="tag" border size="small">{{ tag }}</el-checkbox>
          </el-checkbox-group>
        </div>
        <div class="rate-content">
          <span class="rate-label">详细评价（选填）</span>
          <el-input v-model="newRating.content" type="textarea" :rows="4" maxlength="500" show-word-limit placeholder="说说你的看法..." />
        </div>
      </div>
      <template #footer>
        <el-button @click="showRatingDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" :disabled="!newRating.score" @click="handleSubmitRating">提交评价</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { getReviewTarget, listRatings, createRating, toggleHelpful } from "@/api/review"
import type { ReviewTargetData, RatingData } from "@/api/review"
import { ArrowLeft, Star, StarFilled } from "@element-plus/icons-vue"
import { ElMessage } from "element-plus"

const route = useRoute()
const router = useRouter()
const targetId = Number(route.params.id)

const target = ref<ReviewTargetData | null>(null)
const ratings = ref<RatingData[]>([])
const loading = ref(true)
const ratingsLoading = ref(true)
const ratingSort = ref("hot")
const showRatingDialog = ref(false)
const submitting = ref(false)
const userRating = ref<RatingData | null>(null)

const newRating = ref({ score: 0, tags: [] as string[], content: "" })

const distLabels = ["1-2", "3-4", "5-6", "7-8", "9-10"]
const distData = computed(() => target.value?.distribution || [0, 0, 0, 0, 0])

const tagOptions = computed(() => {
  if (!target.value) return []
  const map: Record<string, string[]> = {
    course: ["讲得清楚", "给分高", "作业多", "收获大", "水分大", "推荐"],
    teacher: ["认真负责", "风趣幽默", "要求严格", "给分好", "干货多"],
    life: ["实惠", "拥挤", "干净", "方便", "环境好", "服务差"],
    product: ["性价比高", "质量好", "成色新", "价格偏高"],
    default: ["好评", "一般", "差评", "推荐", "不推荐"],
  }
  return map[target.value.category] || map.default
})

function categoryLabel(cat: string) {
  const map: Record<string, string> = { course: "课程", teacher: "教师", product: "商品", activity: "活动", life: "生活", other: "其他" }
  return map[cat] || cat
}

function scoreClass(s: number) {
  if (s >= 8) return "score-high"
  if (s >= 6) return "score-mid"
  return "score-low"
}

function pct(val: number, total: number) {
  return total > 0 ? `${(val / total * 100).toFixed(0)}%` : "0%"
}

function smartTime(t: string) {
  const d = new Date(t); const now = new Date(); const diff = now.getTime() - d.getTime()
  if (diff < 86400000) return d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })
  if (diff < 172800000) return "昨天 " + d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })
  return `${d.getMonth() + 1}-${d.getDate()}`
}

async function handleSortChange() {
  ratingsLoading.value = true
  try {
    const res = await listRatings(targetId, { sort: ratingSort.value })
    ratings.value = res.ratings || []
  } catch { /* */ } finally { ratingsLoading.value = false }
}

async function handleHelpful(r: RatingData) {
  try {
    const res = await toggleHelpful(r.id)
    r.helpfulCount = res.helpfulCount
    r.isHelpfulByMe = res.isHelpfulByMe
  } catch { /* */ }
}

async function handleSubmitRating() {
  if (!newRating.value.score) return
  submitting.value = true
  try {
    const rating = await createRating(targetId, {
      score: newRating.value.score,
      tags: newRating.value.tags,
      content: newRating.value.content,
    })
    userRating.value = rating
    showRatingDialog.value = false
    ElMessage.success("评价成功")
    await handleSortChange()
    // Refresh target stats
    const res = await getReviewTarget(targetId)
    target.value = res.target
  } catch { /* */ } finally { submitting.value = false }
}

onMounted(async () => {
  try {
    const res = await getReviewTarget(targetId)
    target.value = res.target
  } catch { /* */ } finally { loading.value = false }
  try {
    const res = await listRatings(targetId, { sort: "hot" })
    ratings.value = res.ratings || []
  } catch { /* */ } finally { ratingsLoading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;
@use "sass:color";

 margin: 0 auto;
.detail-page { max-width: 700px; }

.back-btn { margin-bottom: $space-4; }

// ═══ Score Overview ═══
.score-section {
  display: flex;
  gap: $space-8;
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-lg;
  padding: $space-6;
  margin-bottom: $space-4;
}

.score-left {
  text-align: center;
  min-width: 120px;
}

.score-big {
  font-size: 3rem;
  font-weight: 700;
  line-height: 1;
}

.score-label {
  font-size: $text-xs;
  color: $text-secondary;
  margin-top: $space-2;
}

.score-high { color: $color-success; }
.score-mid { color: #d4a85e; }
.score-low { color: #e74c3c; }

.score-right { flex: 1; }

.dist-row {
  display: flex;
  align-items: center;
  gap: $space-2;
  margin-bottom: 6px;
}

.dist-label {
  font-size: $text-xs;
  color: $text-secondary;
  width: 32px;
  text-align: right;
}

.dist-bar-bg {
  flex: 1;
  height: 8px;
  background: $bg-hover;
  border-radius: 4px;
  overflow: hidden;
}

.dist-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, oklch(0.72 0.18 55), oklch(0.62 0.12 16));
  border-radius: 4px;
  transition: width 0.3s ease;
}

.dist-val {
  font-size: $text-xs;
  color: $text-secondary;
  width: 24px;
  text-align: left;
}

// ═══ Meta ═══
.meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: $space-3;
  margin-bottom: $space-4;
  font-size: $text-sm;
  color: $text-secondary;
  .meta-cat {
    background: rgba($brand-primary-hex, 0.1);
    color: $brand-primary;
    padding: 2px 10px;
    border-radius: 999px;
    font-size: $text-xs;
  }
}

// ═══ Rate Button ═══
.rate-section {
  margin-bottom: $space-5;
}

// ═══ Filter Bar ═══
.filter-bar { margin-bottom: $space-4; }

// ═══ Rating Cards ═══
.ratings-list { display: flex; flex-direction: column; gap: $space-3; }

.rating-card {
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-md;
  padding: $space-4;
  transition: all 0.2s ease;
  &:hover { border-color: rgba($brand-primary-hex, 0.2); }
}

.rc-header {
  display: flex;
  align-items: center;
  gap: $space-2;
  margin-bottom: $space-2;
}

.rc-name {
  font-size: $text-sm;
  font-weight: 600;
  color: $text-primary;
}

.rc-score {
  font-size: 18px;
  font-weight: 700;
  margin-left: auto;
}

.rc-time {
  font-size: $text-xs;
  color: $text-secondary;
  margin-left: $space-2;
}

.rc-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: $space-2;
}

.rc-tag {
  font-size: $text-xs;
  background: rgba($brand-primary-hex, 0.08);
  color: $brand-primary;
  padding: 2px 10px;
  border-radius: 999px;
}

.rc-content {
  font-size: $text-sm;
  color: $text-primary;
  line-height: 1.7;
  margin: 0 0 $space-2;
}

.rc-actions {
  display: flex;
  gap: $space-2;
}

// ═══ Rating Dialog ═══
.rate-dialog-body {
  display: flex;
  flex-direction: column;
  gap: $space-5;
}

.rate-label {
  display: block;
  font-size: $text-sm;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $space-2;
}

.rate-tags {
  :deep(.el-checkbox) {
    margin-right: 8px;
    margin-bottom: 8px;
  }
}
</style>
