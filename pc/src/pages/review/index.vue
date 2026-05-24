<template>
  <div class="review-hub">
    <PageHeader title="评价广场" />

    <div class="hub-toolbar">
      <el-input v-model="searchQuery" placeholder="搜索评价对象..." clearable size="large" class="search-input">
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <el-button type="primary" size="large" @click="$router.push('/reviews/create')" class="create-btn">
        <el-icon><Plus /></el-icon> 发起评价
      </el-button>
    </div>

    <div class="category-tabs">
      <el-radio-group v-model="activeCategory" @change="handleCategoryChange">
        <el-radio-button value="">全部</el-radio-button>
        <el-radio-button value="course">课程</el-radio-button>
        <el-radio-button value="teacher">教师</el-radio-button>
        <el-radio-button value="life">生活</el-radio-button>
        <el-radio-button value="other">其他</el-radio-button>
      </el-radio-group>
    </div>

    <div v-if="loading" class="loading-area">
      <SkeletonCard v-for="i in 6" :key="i" variant="post-card" />
    </div>
    <div v-else-if="!targets.length" class="empty-card">
      <el-empty description="还没有评价对象，成为第一个发起评价的人" :image-size="80">
        <el-button type="primary" @click="$router.push('/reviews/create')">发起评价</el-button>
      </el-empty>
    </div>
    <template v-else>
      <div class="target-grid">
        <div v-for="t in targets" :key="t.id" class="target-card stagger-item" @click="$router.push('/reviews/' + t.id)">
          <div class="tc-top">
            <span class="tc-cat">{{ categoryLabel(t.category) }}</span>
            <el-tag size="small" effect="plain">{{ t.totalRatings }} 评</el-tag>
          </div>
          <div class="tc-name">{{ t.name }}</div>
          <div class="tc-score-row">
            <span class="tc-score" :class="scoreClass(t.avgScore)">{{ t.avgScore.toFixed(1) }}</span>
            <span class="tc-desc">/ 10</span>
            <span class="tc-by">{{ t.creator.nickname }}</span>
          </div>
        </div>
      </div>

      <div class="feed-section">
        <h3 class="feed-title">最新评价</h3>
        <div v-if="feedLoading" class="loading-area"><el-skeleton :rows="3" animated /></div>
        <div v-else-if="!feed.length" class="empty-card"><el-empty description="暂无评价" :image-size="60" /></div>
        <div v-else class="feed-list">
          <div v-for="r in feed" :key="r.id" class="feed-card" @click="$router.push('/reviews/' + r.targetId)">
            <el-avatar :size="28" :src="r.user.avatar">{{ r.user.nickname[0] }}</el-avatar>
            <div class="feed-body">
              <span class="feed-user">{{ r.user.nickname }}</span>
              <span class="feed-target">评价了 · {{ findTargetName(r.targetId) }}</span>
            </div>
            <span class="feed-score" :class="scoreClass(r.score)">{{ r.score }}</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRouter } from "vue-router"
import { listReviewTargets, listRatings } from "@/api/review"
import type { ReviewTargetData, RatingData } from "@/api/review"
import PageHeader from "@/components/PageHeader.vue"
import SkeletonCard from "@/components/SkeletonCard.vue"
import { Search, Plus } from "@element-plus/icons-vue"
import { ElMessage } from "element-plus"

const router = useRouter()
const targets = ref<ReviewTargetData[]>([])
const feed = ref<RatingData[]>([])
const loading = ref(true)
const feedLoading = ref(true)
const searchQuery = ref("")
const activeCategory = ref("")

// Computed filtered targets (client-side filter when no backend API)
const filtered = computed(() => {
  let list = targets.value
  if (activeCategory.value) list = list.filter(t => t.category === activeCategory.value)
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(t => t.name.toLowerCase().includes(q))
  }
  return list
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

function findTargetName(targetId: number) {
  const t = targets.value.find(t => t.id === targetId)
  return t?.name || "未知对象"
}

function handleCategoryChange() { /* computed handles re-filter */ }

onMounted(async () => {
  try {
    const res = await listReviewTargets({ limit: 50 })
    targets.value = res.targets || []
  } catch { /* */ } finally { loading.value = false }

  try {
    const all = await Promise.all(
      targets.value.slice(0, 5).map(t => listRatings(t.id, { limit: 3 }).catch(() => ({ ratings: [] as RatingData[], total: 0 })))
    )
    const feedItems = all.flatMap((r: { ratings: RatingData[] }) => r.ratings)
    feed.value = feedItems.sort((a: RatingData, b: RatingData) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()).slice(0, 10)
  } catch { /* */ } finally { feedLoading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

 margin: 0 auto;
.review-hub { max-width: 1100px; }

.hub-toolbar {
  display: flex;
  gap: $space-3;
  margin-bottom: $space-5;
  .search-input { flex: 1; }
}

.category-tabs {
  margin-bottom: $space-5;
}

.target-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: $space-4;
  margin-bottom: $space-8;
}

.target-card {
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-md;
  padding: $space-5;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  }
}

.tc-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $space-3;
}

.tc-cat {
  font-size: $text-xs;
  color: $text-secondary;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.tc-name {
  font-size: $text-lg;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $space-3;
  line-height: 1.3;
}

.tc-score-row {
  display: flex;
  align-items: baseline;
  gap: $space-1;
}

.tc-score {
  font-size: 2rem;
  font-weight: 700;
  line-height: 1;
}

.tc-desc {
  font-size: $text-sm;
  color: $text-secondary;
}

.tc-by {
  margin-left: auto;
  font-size: $text-xs;
  color: $text-secondary;
}

.score-high { color: $color-success; }
.score-mid { color: #d4a85e; }
.score-low { color: #e74c3c; }

// ═══ Feed Section ═══
.feed-section {
  h3.feed-title { font-size: 18px; font-weight: 600; margin-bottom: $space-4; }
}

.feed-list { display: flex; flex-direction: column; gap: 6px; }

.feed-card {
  display: flex;
  align-items: center;
  gap: $space-3;
  padding: $space-3;
  border-radius: $radius-md;
  transition: all 0.2s ease;
  cursor: pointer;

  &:hover { background: $bg-hover; }
}

.feed-body {
  flex: 1;
  min-width: 0;
  font-size: $text-sm;
}

.feed-user { font-weight: 600; color: $text-primary; }
.feed-target { color: $text-secondary; margin-left: $space-1; }

.feed-score {
  font-size: 18px;
  font-weight: 700;
  flex-shrink: 0;
}

.create-btn {
  white-space: nowrap;
}
</style>
