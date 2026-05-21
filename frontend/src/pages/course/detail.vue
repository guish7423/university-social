<template>
  <view class="page">
    <view v-if="loading" class="loading"><u-loading mode="flower"></u-loading></view>
    <template v-else-if="course">
      <view class="course-header">
        <text class="course-name">{{ course.name }}</text>
        <text class="course-teacher">{{ course.teacher || '未知老师' }}</text>
        <view class="course-meta" v-if="course.school">
          <text class="meta-text">{{ course.school }}</text>
        </view>
      </view>

      <view class="stats-row">
        <view class="stat-item">
          <text class="stat-value">{{ course.avg_teaching?.toFixed(1) || '-' }}</text>
          <text class="stat-label">教学质量</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ course.avg_difficulty?.toFixed(1) || '-' }}</text>
          <text class="stat-label">课程难度</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ course.avg_grading?.toFixed(1) || '-' }}</text>
          <text class="stat-label">给分好坏</text>
        </view>
      </view>

      <view class="rating-section">
        <view class="section-header">
          <text class="section-title">{{ course.user_rating ? '我的评价' : '评价该课程' }}</text>
        </view>
        <view v-if="!course.user_rating" class="rate-form">
          <view class="rate-row" v-for="r in rateFields" :key="r.key">
            <text class="rate-label">{{ r.label }}</text>
            <u-rate v-model="form[r.key]" :count="5" active-color="#667eea"></u-rate>
          </view>
          <u-textarea v-model="form.comment" placeholder="写下你的评价（选填）" height="120" :count="true" maxlength="500"></u-textarea>
          <u-button type="primary" :loading="submitting" @click="submitRating" class="submit-btn">提交评价</u-button>
        </view>
        <view v-else class="my-rating">
          <view class="rate-row" v-for="r in rateFields" :key="r.key">
            <text class="rate-label">{{ r.label }}</text>
            <u-rate :modelValue="course.user_rating[r.key]" :count="5" disabled active-color="#667eea"></u-rate>
          </view>
          <text class="my-comment" v-if="course.user_rating.comment">{{ course.user_rating.comment }}</text>
        </view>
      </view>

      <view class="comments-section">
        <view class="section-header">
          <text class="section-title">全部评价 ({{ course.rating_count || 0 }})</text>
          <u-icon name="order" size="16" color="#667eea"></u-icon>
        </view>
        <view v-if="!ratings || ratings.length === 0" class="empty-ratings">
          <text class="empty-text">暂无评价，来做第一个吧</text>
        </view>
        <view v-else class="ratings-list">
          <view class="rating-item" v-for="item in ratings" :key="item.id">
            <view class="rating-user">
              <u-avatar :src="item.avatar || ''" size="32"></u-avatar>
              <text class="rating-nickname">{{ item.is_anonymous ? '匿名用户' : (item.nickname || '用户') }}</text>
            </view>
            <view class="rating-stars">
              <text class="star-label">教 {{ item.teaching_rating }}★</text>
              <text class="star-label">难 {{ item.difficulty }}★</text>
              <text class="star-label">分 {{ item.grading }}★</text>
            </view>
            <text class="rating-comment" v-if="item.comment">{{ item.comment }}</text>
            <view class="rating-bottom">
              <text class="rating-time">{{ item.created_at?.slice(0, 10) }}</text>
              <view class="helpful-btn" @click="doHelpful(item)">
                <u-icon name="thumb-up" size="14" color="#999"></u-icon>
                <text class="helpful-count">{{ item.helpful_count || 0 }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>
    </template>
  </view>
</template>

<script>
import { getCourseDetail, createRating, listRatings, markHelpful } from "@/api/course"

export default {
  data() {
    return {
      course: null,
      ratings: [],
      loading: true,
      submitting: false,
      form: { teaching_rating: 5, difficulty: 3, grading: 3, comment: "" },
      rateFields: [
        { key: "teaching_rating", label: "教学质量" },
        { key: "difficulty", label: "课程难度" },
        { key: "grading", label: "给分好坏" },
      ],
    }
  },
  onLoad(opts) {
    if (opts.id) this.loadCourse(Number(opts.id))
  },
  methods: {
    async loadCourse(id) {
      this.loading = true
      try {
        const [detail, ratings] = await Promise.all([
          getCourseDetail(id),
          listRatings(id)
        ])
        this.course = detail
        this.ratings = ratings || []
      } finally {
        this.loading = false
      }
    },
    async submitRating() {
      if (!this.course) return
      this.submitting = true
      try {
        await createRating(this.course.id, { ...this.form, is_anonymous: true })
        uni.showToast({ title: "评价成功", icon: "success" })
        this.loadCourse(this.course.id)
      } finally {
        this.submitting = false
      }
    },
    async doHelpful(item) {
      try {
        await markHelpful(item.id)
        item.helpful_count = (item.helpful_count || 0) + 1
      } catch {}
    }
  }
}
</script>

<style lang="scss">
.page {
  min-height: 100vh;
  background: var(--color-canvas, #f5f5f5);
  padding: var(--space-md, 24rpx);
}
.loading {
  padding: 80rpx 0;
  display: flex;
  justify-content: center;
}
.course-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: var(--rounded-lg, 16rpx);
  padding: var(--space-lg, 32rpx);
  color: #fff;
  margin-bottom: var(--space-md, 24rpx);
  .course-name {
    font-size: 36rpx; font-weight: 700; display: block; margin-bottom: 8rpx;
  }
  .course-teacher {
    font-size: 28rpx; opacity: 0.9;
  }
  .course-meta { margin-top: 8rpx; }
  .meta-text { font-size: 24rpx; opacity: 0.7; }
}
.stats-row {
  display: flex; gap: var(--space-md, 24rpx); margin-bottom: var(--space-md, 24rpx);
}
.stat-item {
  flex: 1; background: var(--color-surface, #fff); border-radius: var(--rounded-md, 12rpx);
  padding: var(--space-lg, 32rpx) var(--space-md, 24rpx); text-align: center;
  box-shadow: var(--shadow-sm, 0 1px 3px rgba(0,0,0,0.06));
  .stat-value {
    font-size: 40rpx; font-weight: 700; color: #667eea; display: block; margin-bottom: 8rpx;
  }
  .stat-label { font-size: 24rpx; color: var(--ink-muted, #999); }
}
.rating-section, .comments-section {
  background: var(--color-surface, #fff);
  border-radius: var(--rounded-lg, 16rpx);
  padding: var(--space-lg, 32rpx);
  margin-bottom: var(--space-md, 24rpx);
  box-shadow: var(--shadow-sm, 0 1px 3px rgba(0,0,0,0.06));
}
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: var(--space-md, 24rpx); }
.section-title { font-size: 30rpx; font-weight: 600; color: var(--ink, #333); }
.rate-row { display: flex; justify-content: space-between; align-items: center; margin-bottom: var(--space-md, 24rpx); }
.rate-label { font-size: 28rpx; color: var(--ink, #333); }
.submit-btn { margin-top: var(--space-lg, 32rpx); }
.my-rating { }
.my-comment { display: block; font-size: 26rpx; color: var(--ink-muted, #666); margin-top: var(--space-md, 24rpx); }
.ratings-list { display: flex; flex-direction: column; gap: var(--space-lg, 32rpx); }
.rating-item {
  border-bottom: 1px solid var(--hairline, #f0f0f0); padding-bottom: var(--space-md, 24rpx);
}
.rating-user {
  display: flex; align-items: center; gap: 16rpx; margin-bottom: 8rpx;
}
.rating-nickname { font-size: 26rpx; color: var(--ink-muted, #666); }
.rating-stars {
  display: flex; gap: 24rpx; margin-bottom: 8rpx;
  .star-label { font-size: 24rpx; color: #667eea; }
}
.rating-comment { display: block; font-size: 26rpx; color: var(--ink, #333); margin-bottom: 8rpx; line-height: 1.5; }
.rating-bottom {
  display: flex; justify-content: space-between; align-items: center;
}
.rating-time { font-size: 22rpx; color: var(--ink-subtle, #bbb); }
.helpful-btn {
  display: flex; align-items: center; gap: 6rpx;
  .helpful-count { font-size: 22rpx; color: #999; }
}
.empty-ratings { text-align: center; padding: 60rpx 0; }
.empty-text { font-size: 26rpx; color: var(--ink-muted, #999); }
</style>
