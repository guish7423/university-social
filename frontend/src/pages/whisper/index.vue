<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listWhispers, toggleWhisperLike, whisperSummary, createWhisper, type WhisperData } from "@/api/whisper"
import { useUserStore } from "@/store/user"
import ReportPopup from "@/components/ReportPopup.vue"

const userStore = useUserStore()
const whispers = ref<WhisperData[]>([])
const summary = ref({ total: 0, today: 0 })
const loading = ref(false)
const tabCurrent = ref(0)
const showCreate = ref(false)
const formContent = ref("")
const formType = ref(0)
const creating = ref(false)

const showReport = ref(false)
const reportContentId = ref(0)
const reportContentType = ref("whisper")

const tabs = ["全部", "表白墙"]

onMounted(async () => {
  loading.value = true
  try {
    const [w, s] = await Promise.all([listWhispers(), whisperSummary()])
    whispers.value = w
    summary.value = s
  } catch {}
  loading.value = false
})

async function handleLike(item: WhisperData) {
  if (!userStore.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" })
    return
  }
  const res = await toggleWhisperLike(item.id)
  item.is_liked = res.liked
  item.like_count += res.liked ? 1 : -1
}

function handleMore(item: WhisperData) {
  uni.showActionSheet({
    itemList: ["举报"],
    success(res) {
      if (res.tapIndex === 0) {
        reportContentId.value = item.id
        reportContentType.value = "whisper"
        showReport.value = true
      }
    },
  })
}

function handleTab(idx: number) {
  tabCurrent.value = idx
}

async function submitWhisper() {
  if (!formContent.value.trim()) return
  creating.value = true
  try {
    await createWhisper({
      content: formContent.value.trim(),
      whisper_type: formType.value,
    })
    uni.showToast({ title: "已发送到树洞", icon: "success" })
    showCreate.value = false
    formContent.value = ""
    const w = await listWhispers()
    whispers.value = w
  } catch {}
  creating.value = false
}
</script>

<template>
  <view class="container">
    <view class="whisper-header">
      <text class="whisper-title">树洞</text>
      <text class="whisper-subtitle">匿名倾诉 · 安全释放</text>
      <view class="whisper-stats">
        <view class="stat-item">
          <text class="stat-num">{{ summary.total }}</text>
          <text class="stat-label">总帖子</text>
        </view>
        <view class="stat-item">
          <text class="stat-num">{{ summary.today }}</text>
          <text class="stat-label">今日</text>
        </view>
      </view>
    </view>

    <view class="tabs">
      <view
        v-for="(tab, idx) in tabs"
        :key="idx"
        :class="['tab-item', tabCurrent === idx && 'tab-active']"
        @click="handleTab(idx)"
      >
        <text>{{ tab }}</text>
      </view>
    </view>

    <view v-if="loading" class="loading-state">
      <u-loading mode="flower" size="60" />
      <text class="loading-text">加载中...</text>
    </view>

    <template v-else-if="whispers && whispers.length">
      <view class="whisper-list">
        <view
          v-for="item in whispers"
          :key="item.id"
          class="whisper-card"
          @click="goDetail(item.id)"
        >
          <view class="card-top">
            <view class="codename-badge">
              <text class="codename-icon">🫥</text>
              <text class="codename-text">{{ item.codename || "匿名用户" }}</text>
            </view>
            <text v-if="item.whisper_type === 1" class="confess-tag">💕 表白</text>
            <view class="card-top-right">
              <u-icon name="more-dot-fill" color="#c0c4cc" size="28" @click.stop="handleMore(item)" />
              <text class="card-time">{{ formatTime(item.created_at) }}</text>
            </view>
          </view>

          <view class="card-body">
            <text class="card-content">{{ item.content }}</text>
            <view v-if="item?.images?.length" class="card-images">
              <image
                v-for="(img, i) in item.images.slice(0, 3)"
                :key="i"
                class="card-image"
                :src="img"
                mode="aspectFill"
              />
            </view>
          </view>

          <view class="card-actions">
            <view class="action-btn" @click.stop="handleLike(item)">
              <u-icon :name="item.is_liked ? 'heart-fill' : 'heart'" :color="item.is_liked ? '#f56c6c' : '#c0c4cc'" size="32" />
              <text :class="['action-text', item.is_liked && 'liked']">{{ item.like_count || "赞" }}</text>
            </view>
            <view class="action-btn">
              <u-icon name="chat" color="#c0c4cc" size="32" />
              <text class="action-text">{{ item.comment_count || "评论" }}</text>
            </view>
          </view>
        </view>
      </view>

      <view class="load-more" @click="loadMore">
        <text>加载更多</text>
      </view>
    </template>

    <view v-else class="empty-state">
      <u-icon name="chat" size="120" color="#e0e0e0" />
      <text class="empty-text">还没有树洞消息</text>
      <text class="empty-desc">点击右下角开始匿名倾诉</text>
    </view>

    <view v-if="userStore.isLogin" class="fab" @click="showCreate = true">
      <u-icon name="edit-pen" color="#fff" size="40" />
    </view>

    <u-popup :show="showCreate" mode="bottom" @close="showCreate = false">
      <view class="popup-body">
        <text class="popup-title">匿名倾诉</text>
        <view class="type-select">
          <view :class="['type-btn', formType === 0 && 'type-active']" @click="formType = 0">
            <text>普通</text>
          </view>
          <view :class="['type-btn', formType === 1 && 'type-active']" @click="formType = 1">
            <text>💕 表白</text>
          </view>
        </view>
        <textarea
          v-model="formContent"
          placeholder="写下你想说的话..."
          maxlength="500"
          class="whisper-textarea"
        ></textarea>
        <u-button
          type="primary"
          :loading="creating"
          :disabled="!formContent.trim()"
          @click="submitWhisper"
        >
          匿名发送
        </u-button>
      </view>
    </u-popup>
  </view>
<ReportPopup :visible="showReport" :content-type="reportContentType" :content-id="reportContentId" @close="showReport = false" />
</template>

<style scoped lang="scss">
.container {
  min-height: 100vh;
  background: linear-gradient(180deg, #f0f2ff 0%, #f8f9fa 100%);
}

.whisper-header {
  padding: 64rpx 32rpx 40rpx;
  text-align: center;
  position: relative;
  &::before {
    content: "";
    position: absolute;
    top: -120rpx;
    left: 50%;
    transform: translateX(-50%);
    width: 400rpx;
    height: 400rpx;
    background: radial-gradient(circle, rgba(102,126,234,0.15) 0%, transparent 70%);
    border-radius: 50%;
    pointer-events: none;
  }
}

.whisper-title {
  font-size: 48rpx;
  font-weight: 700;
  color: var(--ink, #1a1a2e);
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.whisper-subtitle {
  display: block;
  font-size: 26rpx;
  color: var(--ink-muted, #909399);
  margin-top: 8rpx;
}

.whisper-stats {
  display: flex;
  justify-content: center;
  gap: 48rpx;
  margin-top: 32rpx;
}

.stat-item {
  text-align: center;
}

.stat-num {
  display: block;
  font-size: 36rpx;
  font-weight: 700;
  color: var(--ink, #1a1a2e);
}

.stat-label {
  font-size: 24rpx;
  color: var(--ink-muted, #909399);
}

.tabs {
  display: flex;
  margin: 0 32rpx 24rpx;
  background: var(--color-surface, #fff);
  border-radius: var(--rounded-lg, 16rpx);
  padding: 6rpx;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 16rpx 0;
  font-size: 28rpx;
  color: var(--ink-muted, #909399);
  border-radius: 12rpx;
  transition: all 0.2s;
}

.tab-active {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  font-weight: 600;
}

.whisper-list {
  padding: 0 32rpx 32rpx;
}

.whisper-card {
  background: var(--color-surface, #fff);
  border-radius: var(--rounded-lg, 22rpx);
  padding: 28rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
  transition: transform 0.15s;
  &:active { transform: scale(0.98); }
}

.card-top {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.codename-badge {
  display: flex;
  align-items: center;
  gap: 6rpx;
  background: linear-gradient(135deg, #667eea15, #764ba215);
  padding: 6rpx 16rpx;
  border-radius: 20rpx;
}

.codename-icon { font-size: 24rpx; }

.codename-text {
  font-size: 24rpx;
  color: #667eea;
  font-weight: 600;
}

.confess-tag {
  font-size: 22rpx;
  background: #fef0f0;
  color: #f56c6c;
  padding: 4rpx 12rpx;
  border-radius: 12rpx;
}

.card-time {
  margin-left: auto;
  font-size: 22rpx;
  color: var(--ink-subtle, #c0c4cc);
}

.card-body { margin-bottom: 20rpx; }

.card-content {
  font-size: 28rpx;
  color: var(--ink, #1a1a2e);
  line-height: 1.6;
}

.card-images {
  display: flex;
  gap: 8rpx;
  margin-top: 12rpx;
}

.card-image {
  width: 200rpx;
  height: 200rpx;
  border-radius: var(--rounded-md, 12rpx);
}

.card-actions {
  display: flex;
  gap: 32rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid var(--hairline, #f0f0f0);
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.action-text {
  font-size: 24rpx;
  color: var(--ink-muted, #909399);
}

.liked { color: #f56c6c; }

.fab {
  position: fixed;
  right: 32rpx;
  bottom: 100rpx;
  width: 100rpx;
  height: 100rpx;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(102,126,234,0.4);
  transition: transform 0.15s;
  &:active { transform: scale(0.9); }
}

.popup-body {
  padding: 40rpx 32rpx 60rpx;
}

.popup-title {
  font-size: 36rpx;
  font-weight: 700;
  color: var(--ink, #1a1a2e);
  margin-bottom: 24rpx;
  display: block;
}

.type-select {
  display: flex;
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.type-btn {
  flex: 1;
  padding: 20rpx;
  text-align: center;
  border: 2rpx solid var(--hairline, #f0f0f0);
  border-radius: var(--rounded-md, 16rpx);
  font-size: 28rpx;
  color: var(--ink-muted, #909399);
  transition: all 0.2s;
}

.type-active {
  border-color: #667eea;
  color: #667eea;
  background: #667eea08;
}

.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text, .empty-text {
  margin-top: 20rpx;
  font-size: 28rpx;
  color: var(--ink-subtle, #c0c4cc);
}

.empty-desc {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: var(--ink-tertiary, #e0e0e0);
}

.load-more {
  text-align: center;
  padding: 24rpx;
  font-size: 26rpx;
  color: var(--ink-subtle, #c0c4cc);
}
</style>
