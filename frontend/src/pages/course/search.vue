<template>
  <view class="page">
    <view class="search-box">
      <u-search
        v-model="keyword"
        placeholder="搜索课程名/老师/学院"
        :show-action="false"
        @search="doSearch"
        @custom="doSearch"
      ></u-search>
    </view>

    <view class="quick-tips" v-if="!searched">
      <view class="tip-icon">📚</view>
      <text class="tip-text">输入课程名称、老师姓名或学院名称搜索课程评价</text>
    </view>

    <view v-if="searched">
      <view v-if="loading" class="loading"><u-loading mode="flower"></u-loading></view>
      <view v-else-if="courses.length === 0" class="empty">
        <text class="empty-text">未找到相关课程</text>
      </view>
      <view v-else class="course-list">
        <view
          class="course-card"
          v-for="item in courses"
          :key="item.id"
          @click="goDetail(item.id)"
        >
          <view class="card-top">
            <text class="course-name">{{ item.name }}</text>
            <text class="course-teacher">{{ item.teacher || '未知老师' }}</text>
          </view>
          <view class="card-meta" v-if="item.school || item.department">
            <text class="meta-text">{{ item.school }} {{ item.department }}</text>
          </view>
          <u-icon name="arrow-right" color="#bbb" size="14"></u-icon>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { searchCourses } from "@/api/course"

export default {
  data() {
    return {
      keyword: "",
      courses: [],
      loading: false,
      searched: false,
    }
  },
  methods: {
    async doSearch() {
      if (!this.keyword.trim()) return
      this.loading = true
      this.searched = true
      try {
        const res = await searchCourses(this.keyword)
        this.courses = res || []
      } finally {
        this.loading = false
      }
    },
    goDetail(id) {
      uni.navigateTo({ url: `/pages/course/detail?id=${id}` })
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
.search-box {
  margin-bottom: var(--space-lg, 32rpx);
}
.quick-tips {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 40rpx 0;
  .tip-icon {
    font-size: 64rpx;
    margin-bottom: 24rpx;
  }
  .tip-text {
    font-size: 28rpx;
    color: var(--ink-muted, #999);
    text-align: center;
    line-height: 1.6;
  }
}
.loading {
  padding: 80rpx 0;
  display: flex;
  justify-content: center;
}
.empty {
  padding: 120rpx 0;
  text-align: center;
  .empty-text {
    font-size: 28rpx;
    color: var(--ink-muted, #999);
  }
}
.course-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-md, 24rpx);
}
.course-card {
  background: var(--color-surface, #fff);
  border-radius: var(--rounded-lg, 16rpx);
  padding: var(--space-lg, 32rpx);
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: var(--shadow-sm, 0 1px 3px rgba(0,0,0,0.06));
  .card-top {
    display: flex;
    flex-direction: column;
    gap: 8rpx;
    flex: 1;
  }
  .course-name {
    font-size: 32rpx;
    font-weight: 600;
    color: var(--ink, #333);
  }
  .course-teacher {
    font-size: 26rpx;
    color: var(--ink-muted, #999);
  }
  .card-meta {
    margin-top: 4rpx;
    .meta-text {
      font-size: 22rpx;
      color: var(--ink-subtle, #bbb);
    }
  }
}
</style>
