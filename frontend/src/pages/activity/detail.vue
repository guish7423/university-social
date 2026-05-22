<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getActivity, joinActivity, leaveActivity, deleteActivity, activityParticipants, listActivityComments, createActivityComment, type ActivityData, type ActivityComment } from "@/api/activity"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const activity = ref<ActivityData | null>(null)
const comments = ref<ActivityComment[]>([])
const participants = ref<any[]>([])
const loading = ref(false)
const commentText = ref("")
const sending = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const pages = getCurrentPages()
    const page = pages[pages.length - 1] as any
    const id = parseInt(page?.$page?.query?.id || '0')
    if (!id) { uni.showToast({ title: '参数错误', icon: 'none' }); return }
    const [a, c, p] = await Promise.all([
      getActivity(id),
      listActivityComments(id),
      activityParticipants(id)
    ])
    activity.value = a
    comments.value = c
    participants.value = p
  } catch { uni.showToast({ title: '活动不存在', icon: 'none' }) }
  loading.value = false
})

async function handleJoin() {
  if (!activity.value) return
  if (!userStore.isLogin) { uni.navigateTo({ url: "/pages/login/index" }); return }
  try {
    await joinActivity(activity.value.id)
    activity.value.is_participant = true
    activity.value.participant_count++
    uni.showToast({ title: "报名成功", icon: "success" })
  } catch { uni.showToast({ title: "报名失败", icon: "none" }) }
}

async function handleLeave() {
  if (!activity.value) return
  try {
    await leaveActivity(activity.value.id)
    activity.value.is_participant = false
    activity.value.participant_count--
    uni.showToast({ title: "已取消报名", icon: "success" })
  } catch { uni.showToast({ title: "操作失败", icon: "none" }) }
}

async function handleDelete() {
  if (!activity.value) return
  uni.showModal({
    title: "提示",
    content: "确定删除此活动？",
    success: async (r) => {
      if (r.confirm) {
        await deleteActivity(activity.value!.id)
        uni.showToast({ title: "已删除", icon: "success" })
        setTimeout(() => uni.navigateBack(), 1000)
      }
    },
  })
}

async function sendComment() {
  if (!activity.value || !commentText.value.trim()) return
  sending.value = true
  try {
    await createActivityComment(activity.value.id, commentText.value.trim())
    commentText.value = ""
    comments.value = await listActivityComments(activity.value.id)
  } catch {}
  sending.value = false
}

function formatTime(t: string) { return t ? t.slice(0, 16).replace("T", " ") : "" }
</script>

<template>
  <u-loading v-if="loading" mode="flower" size="48" />

  <view v-else-if="activity" class="container">
    <view class="header-card">
      <view class="header-top">
        <text class="title">{{ activity.title }}</text>
        <u-tag :text="activity.activity_type" type="primary" size="mini" shape="circle" plain />
      </view>
      <view class="info-grid">
        <view class="info-item">
          <text class="info-label">🕐 开始</text>
          <text class="info-value">{{ formatTime(activity.start_time) }}</text>
        </view>
        <view v-if="activity.end_time" class="info-item">
          <text class="info-label">🔚 结束</text>
          <text class="info-value">{{ formatTime(activity.end_time) }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">📍 地点</text>
          <text class="info-value">{{ activity.location || '待定' }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">👥 人数</text>
          <text class="info-value">{{ activity.participant_count }}/{{ activity.max_participants || '不限' }}</text>
        </view>
      </view>
    </view>

    <view class="section">
      <text class="section-title">活动介绍</text>
      <text class="desc-text">{{ activity.description || '暂无描述' }}</text>
    </view>

    <view v-if="participants.length > 0" class="section">
      <text class="section-title">参与成员 ({{ participants.length }})</text>
      <view class="avatar-wall">
        <view v-for="u in participants.slice(0, 8)" :key="u.id" class="avatar-item">
          <image :src="u.avatar || '/static/default-avatar.png'" class="avatar-img" mode="aspectFill" />
          <text class="avatar-name">{{ u.nickname?.slice(0, 2) || '?' }}</text>
        </view>
        <view v-if="participants.length > 8" class="avatar-more">
          <text>+{{ participants.length - 8 }}</text>
        </view>
      </view>
    </view>

    <view class="section">
      <text class="section-title">评论 ({{ comments.length }})</text>
      <view class="comment-list">
        <view v-for="c in comments" :key="c.id" class="comment-item">
          <text class="comment-content">{{ c.content }}</text>
          <text class="comment-time">{{ formatTime(c.created_at) }}</text>
        </view>
      </view>
      <view v-if="comments.length === 0" class="no-comment">
        <text>暂无评论</text>
      </view>
    </view>

    <view class="bottom-bar">
      <view class="comment-input-area">
        <input v-model="commentText" class="comment-input" placeholder="说点什么..." />
        <text class="send-btn" :class="{ disabled: sending || !commentText.trim() }" @click="sendComment">
          {{ sending ? '发送中' : '发送' }}
        </text>
      </view>
      <view class="action-btns">
        <view v-if="activity.is_owner" class="action-btn delete" @click="handleDelete">
          <text>删除活动</text>
        </view>
        <view v-else-if="activity.is_participant" class="action-btn leave" @click="handleLeave">
          <text>取消报名</text>
        </view>
        <view v-else class="action-btn join" @click="handleJoin">
          <text>立即报名</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style lang="scss">
page { background: #F7F4F0; }
.loading { text-align: center; padding: 200rpx; color: #999; }
.container { padding-bottom: 200rpx; }

.header-card {
  background: #fff;
  margin: 20rpx;
  padding: 24rpx;
  border-radius: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);
}
.header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}
.title {
  font-size: 34rpx;
  font-weight: 700;
  color: #333;
  flex: 1;
  margin-right: 12rpx;
}
.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16rpx;
}
.info-item {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}
.info-label {
  font-size: 22rpx;
  color: #999;
}
.info-value {
  font-size: 26rpx;
  color: #444;
  font-weight: 500;
}

.section {
  background: #fff;
  margin: 0 20rpx 20rpx;
  padding: 24rpx;
  border-radius: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);
}
.section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}
.desc-text {
  font-size: 26rpx;
  color: #555;
  line-height: 1.6;
}

.avatar-wall {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}
.avatar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4rpx;
}
.avatar-img {
  width: 70rpx;
  height: 70rpx;
  border-radius: 50%;
  background: #f0f2f5;
}
.avatar-name {
  font-size: 20rpx;
  color: #888;
}
.avatar-more {
  width: 70rpx;
  height: 70rpx;
  border-radius: 50%;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
  color: #999;
}

.comment-list { display: flex; flex-direction: column; gap: 16rpx; }
.comment-item {
  padding: 12rpx 0;
  border-bottom: 1rpx solid #F7F4F0;
}
.comment-content { font-size: 26rpx; color: #333; display: block; }
.comment-time { font-size: 20rpx; color: #bbb; margin-top: 4rpx; display: block; }
.no-comment { text-align: center; color: #999; font-size: 24rpx; padding: 20rpx 0; }

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  padding: 16rpx 20rpx;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -2rpx 12rpx rgba(0,0,0,0.06);
}
.comment-input-area {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 12rpx;
}
.comment-input {
  flex: 1;
  height: 60rpx;
  background: #F7F4F0;
  border-radius: 30rpx;
  padding: 0 24rpx;
  font-size: 26rpx;
}
.send-btn {
  padding: 12rpx 24rpx;
  background: linear-gradient(135deg, #C67A6A, #1E2A3A);
  color: #fff;
  border-radius: 30rpx;
  font-size: 24rpx;
  &.disabled { opacity: 0.5; }
}
.action-btns { display: flex; gap: 12rpx; }
.action-btn {
  flex: 1;
  text-align: center;
  padding: 16rpx;
  border-radius: 16rpx;
  font-size: 28rpx;
  font-weight: 600;
}
.join {
  background: linear-gradient(135deg, #C67A6A, #1E2A3A);
  color: #fff;
}
.leave {
  background: #F7F4F0;
  color: #666;
}
.delete {
  background: #fff1f0;
  color: #ff4d4f;
}
</style>
