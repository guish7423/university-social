<script setup lang="ts">
import { ref, onMounted } from "vue"
import { myActivities, deleteActivity, type ActivityData } from "@/api/activity"

const created = ref<ActivityData[]>([])
const joined = ref<ActivityData[]>([])
const currentTab = ref<"created" | "joined">("created")
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const res = await myActivities()
    created.value = res.created
    joined.value = res.joined
  } catch (e) { console.error(e) }
  loading.value = false
})

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/activity/detail?id=${id}` })
}

async function handleDelete(id: number) {
  uni.showModal({
    title: "提示",
    content: "确定删除此活动？",
    success: async (r) => {
      if (r.confirm) {
        await deleteActivity(id)
        const res = await myActivities()
        created.value = res.created
        joined.value = res.joined
        uni.showToast({ title: "已删除", icon: "success" })
      }
    },
  })
}

function formatTime(t: string) { return t ? t.slice(0, 16).replace("T", " ") : "" }
</script>

<template>
  <view class="container">
    <view class="tab-bar">
      <view
        class="tab-item" :class="{ active: currentTab === 'created' }"
        @click="currentTab = 'created'"
      >
        <text>我创建的 ({{ created.length }})</text>
      </view>
      <view
        class="tab-item" :class="{ active: currentTab === 'joined' }"
        @click="currentTab = 'joined'"
      >
        <text>我参与的 ({{ joined.length }})</text>
      </view>
    </view>

    <u-loading v-if="loading" mode="flower" size="48" />

    <view v-else class="activity-list">
      <view
        v-for="(item, idx) in (currentTab === 'created' ? created : joined)" :key="item.id"
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
          <text class="meta-item">👥 {{ item.participant_count }}/{{ item.max_participants || '∞' }}</text>
        </view>
        <view v-if="currentTab === 'created'" class="card-actions" @click.stop>
          <text class="action-link danger" @click.stop="handleDelete(item.id)">删除</text>
        </view>
      </view>
    </view>

    <view v-if="!loading && (currentTab === 'created' ? created : joined).length === 0" class="empty">
      <text>{{ currentTab === 'created' ? '还没有创建过活动' : '还没有参与过活动' }}</text>
    </view>
  </view>
</template>

<style lang="scss">
page { background: #f5f7fa; }
.container { padding-bottom: 40rpx; }

.tab-bar {
  display: flex;
  background: #fff;
  position: sticky;
  top: 0;
  z-index: 10;
}
.tab-item {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  font-size: 28rpx;
  color: #666;
  border-bottom: 2rpx solid transparent;
  transition: all 0.2s;
  &.active {
    color: #C67A6A;
    font-weight: 600;
    border-bottom-color: #C67A6A;
  }
}

.loading { text-align: center; padding: 100rpx; color: #999; }

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
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  flex: 1;
  margin-right: 12rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.card-meta { display: flex; flex-direction: column; gap: 8rpx; }
.meta-item { font-size: 24rpx; color: #888; }
.card-actions {
  margin-top: 12rpx;
  padding-top: 12rpx;
  border-top: 1rpx solid #f0f2f5;
  display: flex;
  justify-content: flex-end;
}
.action-link { font-size: 24rpx; padding: 8rpx 16rpx; }
.danger { color: #ff4d4f; }

.empty { text-align: center; padding: 100rpx 0; color: #999; font-size: 28rpx; }

.stagger-1 { animation: fadeInUp 0.4s ease-out 0s both; }
.stagger-2 { animation: fadeInUp 0.4s ease-out 0.08s both; }
.stagger-3 { animation: fadeInUp 0.4s ease-out 0.16s both; }
.stagger-4 { animation: fadeInUp 0.4s ease-out 0.24s both; }

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20rpx); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
