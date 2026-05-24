<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="8" animated /></div>
  <div v-else-if="!course" class="empty-state"><el-empty description="课程不存在" /></div>
  <div v-else class="detail-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn stagger-item">返回</el-button>

    <!-- Hero Banner -->
    <div class="hero-banner stagger-item">
      <div class="hero-content">
        <div class="hero-badges">
          <el-tag size="small" class="hero-tag">{{ course.department }}</el-tag>
          <el-tag size="small" type="info" class="hero-tag">{{ course.credits }} 学分</el-tag>
        </div>
        <h1 class="hero-title">{{ course.name }}</h1>
        <p class="hero-subtitle">{{ course.teacher }}</p>
        <div class="hero-score" v-if="course.rating">
          <span class="score-big">{{ course.rating.toFixed(1) }}</span>
          <span class="score-label"><el-rate v-model="course.rating" disabled :max="5" /></span>
        </div>
      </div>
      <div class="hero-glow" />
    </div>

    <!-- Tabs -->
    <el-tabs v-model="activeTab" class="detail-tabs stagger-item">
      <!-- 介绍 Tab -->
      <el-tab-pane label="介绍" name="intro">
        <div class="info-grid">
          <div class="info-card">
            <h3><el-icon><InfoFilled /></el-icon> 基本信息</h3>
            <div class="info-rows">
              <div class="info-row"><span class="label">教师</span><span>{{ course.teacher || '未知' }}</span></div>
              <div class="info-row"><span class="label">学院</span><span>{{ course.department || '未知' }}</span></div>
              <div class="info-row"><span class="label">学分</span><span>{{ course.credits }}</span></div>
              <div class="info-row"><span class="label">人数</span><span>{{ course.enrolled }}/{{ course.capacity }}</span></div>
            </div>
          </div>
          <div class="info-card" v-if="course.week || course.time || course.location">
            <h3><el-icon><Calendar /></el-icon> 课程安排</h3>
            <div class="info-rows">
              <div class="info-row" v-if="course.week"><span class="label">周次</span><span>{{ course.week }}</span></div>
              <div class="info-row" v-if="course.time"><span class="label">时间</span><span>{{ course.time }}</span></div>
              <div class="info-row" v-if="course.location"><span class="label">地点</span><span>{{ course.location }}</span></div>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- 评价 Tab -->
      <el-tab-pane label="评价" name="ratings">
        <ReviewSection
          :ratings="ratings"
          :loading="ratingLoading"
          :my-rating="myRating ?? null"
          :rating-dialog-visible="ratingDialogVisible"
          @submit="handleSubmitRating"
          @open-dialog="ratingDialogVisible = true"
          @close-dialog="ratingDialogVisible = false"
        />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"
import { getCourseDetail, createRating, listRatings } from "@/api/course"
import type { CourseData } from "@/api/course"
import { ArrowLeft, InfoFilled, Calendar } from "@element-plus/icons-vue"
import ReviewSection from "@/components/ReviewSection.vue"

const route = useRoute()
const course = ref<CourseData | null>(null)
const loading = ref(true)
const activeTab = ref("intro")

const ratings = ref<any[]>([])
const myRating = ref<any>(null)
const ratingLoading = ref(false)
const ratingDialogVisible = ref(false)

async function handleSubmitRating(data: { score: number; comment: string }) {
  if (!course.value) return
  const res = await createRating(course.value.id, { rating: data.score, comment: data.comment })
  ratingDialogVisible.value = false
  myRating.value = { ...data, id: res.id }
  await loadRatings()
}

async function loadRatings() {
  if (!course.value) return
  ratingLoading.value = true
  try {
    ratings.value = await listRatings(course.value.id)
  } catch { /* handled */ }
  finally { ratingLoading.value = false }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    course.value = await getCourseDetail(id)
    await loadRatings()
  } catch { /* handled */ }
  finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;
@use "sass:color";

 margin: 0 auto;
.detail-page { max-width: 780px; }
.back-btn { margin-bottom: 16px; }

// ── Hero Banner ──
.hero-banner {
  position: relative;
  overflow: hidden;
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-md;
  padding: 32px;
  margin-bottom: 20px;
}

.hero-glow {
  position: absolute;
  top: -50%;
  right: -20%;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, color.change($brand-primary, $alpha: 0.08), transparent 70%);
  pointer-events: none;
}

.hero-content { position: relative; z-index: 1; }
.hero-badges { display: flex; gap: 8px; margin-bottom: 12px; }
.hero-tag { border: none; }

.hero-title {
  font-size: $text-3xl;
  font-weight: 700;
  margin-bottom: 6px;
  line-height: 1.2;
}

.hero-subtitle {
  font-size: $text-base;
  color: $text-secondary;
  margin-bottom: 16px;
}

.hero-score {
  display: flex; align-items: center; gap: 10px;
  .score-big {
    font-size: 2.5rem; font-weight: 700; color: $brand-primary; line-height: 1;
  }
  .score-label { :deep(.el-rate) { line-height: 1; } }
}

// ── Tabs ──
.detail-tabs {
  :deep(.el-tabs__item) { font-size: 15px; }
}

// ── Info Grid ──
.info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }

.info-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 20px;
  h3 {
    font-size: 14px; font-weight: 600; margin-bottom: 14px;
    display: flex; align-items: center; gap: 6px;
    .el-icon { font-size: 16px; }
  }
}

.info-rows { display: flex; flex-direction: column; gap: 10px; }
.info-row {
  display: flex; justify-content: space-between; font-size: 13px;
  .label { color: $text-muted; min-width: 60px; }
}
</style>
