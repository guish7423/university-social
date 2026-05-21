<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listActivities, activityTypes, type ActivityData } from "@/api/activity"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const activities = ref<ActivityData[]>([])
const types = ref<string[]>([])
const currentType = ref("全部")
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const [typs] = await Promise.all([activityTypes()])
    types.value = typs
    await loadActivities()
  } catch {}
  loading.value = false
})

async function loadActivities() {
  const t = currentType.value === "全部" ? "" : currentType.value
  activities.value = await listActivities({ type: t })
}

function switchType(t: string) {
  currentType.value = t
  loadActivities()
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/activity/detail?id=${id}` })
}

function goCreate() {
  if (!userStore.isLogin) { uni.navigateTo({ url: "/pages/login/index" }); return }
  uni.navigateTo({ url: "/pages/activity/create" })
}

function formatTime(t: string) {
  return t ? t.slice(0, 16).replace("T", " ") : ""
}
</script>

<template>
  <view class="container">
    <view class="type-bar">
      <scroll-view scroll-x class="type-scroll">
        <view
          v-for="t in ['全部', ...types]" :key="t"
          class="type-item" :class="{ active: currentType === t }"
          @click="switchType(t)"
        >
          <text>{{ t }}</text>
        </view>
      </scroll-view>
    </view>

    <view v-if="loading" class="loading"><text>加载中...</text></view>

    <view v-else class="activity-list">
      <view
        v-for="(item, idx) in activities" :key="item.id"
        class="activity-card"
        :class="'stagger-' + ((idx % 4) + 1)"
        @click="goDetail(item.id)"
      >
        <view class="card-header">
          <text class="card-title">{{ item.title }}</text>
          <u-tag :text="item.activity_type" type="primary" size="mini" shape="circle" plain />
        </view>
        <view class="card-meta">
          <text class="meta-item">📍 {{ item.location || '待定' }}</text>
          <text class="meta-item">🕐 {{ formatTime(item.start_time) }}</text>
        </view>
        <view class="card-footer">
          <text class="participant-info">
            👥 {{ item.participant_count }}/{{ item.max_participants || '∞' }}
          </text>
          <text v-if="item.is_participant" class="joined-badge">已报名</text>
          <text v-else-if="item.is_owner" class="owner-badge">我创建的</text>
        </view>
      </view>
    </view>

    <view v-if="!loading && activities.length === 0" class="empty">
      <text>暂无活动</text>
    </view>

    <view class="fab" @click="goCreate">
      <text class="fab-icon">+</text>
    </view>
  </view>
</template>

<style lang="scss">
page { background: #f5f7fa; }
.container { padding-bottom: 120rpx; }

.type-bar {
  background: #fff;
  padding: 20rpx 0;
  position: sticky;
  top: 0;
  z-index: 10;
}
.type-scroll { white-space: nowrap; padding: 0 20rpx; }
.type-item {
  display: inline-flex;
  padding: 12rpx 28rpx;
  margin-right: 16rpx;
  border-radius: 30rpx;
  background: #f0f2f5;
  font-size: 26rpx;
  color: #666;
  transition: all 0.2s;
  &.active {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: #fff;
    font-weight: 600;
  }
}

.loading { text-align: center; padding: 60rpx; color: #999; }

.activity-list { padding: 20rpx; }
.activity-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);
  transition: all 0.2s;
  &:active { transform: scale(0.98); }
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}
.card-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  flex: 1;
  margin-right: 12rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.card-meta {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
  margin-bottom: 16rpx;
}
.meta-item {
  font-size: 24rpx;
  color: #888;
}
.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 16rpx;
  border-top: 1rpx solid #f0f2f5;
}
.participant-info {
  font-size: 24rpx;
  color: #666;
}
.joined-badge {
  font-size: 22rpx;
  color: #52c41a;
  font-weight: 500;
}
.owner-badge {
  font-size: 22rpx;
  color: #667eea;
  font-weight: 500;
}

.empty { text-align: center; padding: 100rpx 0; color: #999; font-size: 28rpx; }

.fab {
  position: fixed;
  right: 30rpx;
  bottom: 120rpx;
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 20rpx rgba(102,126,234,0.4);
  z-index: 100;
  &:active { transform: scale(0.9); }
}
.fab-icon { color: #fff; font-size: 48rpx; font-weight: 300; line-height: 1; }

.stagger-1 { animation: fadeInUp 0.4s ease-out 0s both; }
.stagger-2 { animation: fadeInUp 0.4s ease-out 0.08s both; }
.stagger-3 { animation: fadeInUp 0.4s ease-out 0.16s both; }
.stagger-4 { animation: fadeInUp 0.4s ease-out 0.24s both; }

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20rpx); }
  to { opacity: 1; transform: translateY(0); }
}

@media (min-width: 1024px) {
.activity-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24rpx;
  padding: 32rpx;
}
.activity-card { margin-bottom: 0; }
.fab { display: none; }
.container { max-width: 1400rpx; margin: 0 auto; }
}
</style>
