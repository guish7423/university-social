<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getPointsBalance, claimDailyLogin, getPointsLeaderboard } from "@/api/points"

const loading = ref(true)
const balance = ref<any>({})
const leaderboard = ref<any[]>([])

const rules = [
  { action: "每日登录", points: "+1" },
  { action: "发布动态", points: "+3" },
  { action: "评论", points: "+1" },
  { action: "邀请注册", points: "+5" },
  { action: "帖子加精", points: "+10" },
]

async function load() {
  loading.value = true
  try {
    balance.value = await getPointsBalance()
    leaderboard.value = (await getPointsLeaderboard()).leaderboard ?? []
  } catch {}
  loading.value = false
}
onMounted(load)
</script>

<template>
  <view class="page">
    <u-loading mode="flower" size="60" v-if="loading" />
    <template v-else>
    <view class="balance-card">
      <text class="balance-label">我的积分</text>
      <text class="balance-value">{{ balance.points ?? 0 }}</text>
      <text class="balance-level">{{ balance.level ?? "Lv.0" }}</text>
      <view v-if="balance.next_level" class="next-level">
        距 {{ balance.next_level }} 还需 {{ balance.next_threshold - (balance.points ?? 0) }} 分
      </view>
    </view>

    <view class="section">
      <text class="section-title">积分规则</text>
      <view class="card">
        <view class="rule-item" v-for="(r, i) in rules" :key="i">
          <text class="rule-action">{{ r.action }}</text>
          <text class="rule-points">{{ r.points > 0 ? '+' : '' }}{{ r.points }}</text>
        </view>
      </view>
    </view>

    <view class="section">
      <view class="section-header">
        <text class="section-title">积分排行</text>
        <text class="section-more">更多</text>
      </view>
      <view class="card">
        <view class="lb-item" v-for="(e, i) in leaderboard" :key="e.user_id">
          <view class="lb-rank" :class="{ 'lb-gold': i===0, 'lb-silver': i===1, 'lb-bronze': i===2 }">
            {{ e.rank }}
          </view>
          <image class="lb-avatar" :src="e.avatar || '/static/default-avatar.png'" mode="aspectFill" />
          <view class="lb-info">
            <text class="lb-name">{{ e.nickname }}</text>
            <text class="lb-level">{{ e.level }}</text>
          </view>
          <text class="lb-points">{{ e.points }}</text>
        </view>
        <view v-if="leaderboard.length === 0" class="empty">暂无排行数据</view>
      </view>
    </view>
    </template>
  </view>
</template>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); padding-bottom: 40rpx; }
.balance-card {
  margin: 30rpx; padding: 40rpx; border-radius: 24rpx;
  display: flex; flex-direction: column; align-items: center;
  background: linear-gradient(135deg, #C67A6A 0%, #1E2A3A 100%);
  color: #fff;
}
.balance-label { font-size: 28rpx; opacity: 0.9; }
.balance-value { font-size: 72rpx; font-weight: 700; margin: 10rpx 0; }
.balance-level { font-size: 32rpx; font-weight: 500; }
.next-level { font-size: 24rpx; opacity: 0.8; margin-top: 16rpx; }

.section { margin: 30rpx; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20rpx; }
.section-title { font-size: 30rpx; font-weight: 600; color: var(--ink, #1E2A3A); }
.section-more { font-size: 24rpx; color: var(--ink-tertiary, #B8C2CE); }

.card { background: var(--color-surface, #fff); border-radius: 20rpx; padding: 20rpx; }
.rule-item { display: flex; justify-content: space-between; padding: 16rpx 0; border-bottom: 1rpx solid var(--hairline-light, #EAE6E0); }
.rule-item:last-child { border-bottom: none; }
.rule-action { font-size: 28rpx; color: var(--ink, #1E2A3A); }
.rule-points { font-size: 28rpx; color: var(--brand-primary, #C67A6A); font-weight: 600; }

.lb-item { display: flex; align-items: center; padding: 16rpx 0; border-bottom: 1rpx solid var(--hairline-light, #EAE6E0); }
.lb-item:last-child { border-bottom: none; }
.lb-rank { width: 48rpx; height: 48rpx; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 26rpx; font-weight: 700; color: var(--ink-muted, #5C6B7E); background: var(--color-surface-1, #F0EDE8); }
.lb-gold { background: linear-gradient(135deg, #E8C96A, #C67A6A); color: #fff; }
.lb-silver { background: linear-gradient(135deg, #E0DBD4, #B8C2CE); color: #fff; }
.lb-bronze { background: linear-gradient(135deg, #E8C6A0, #C67A6A); color: #fff; }
.lb-avatar { width: 56rpx; height: 56rpx; border-radius: 50%; margin: 0 16rpx; background: var(--color-surface-1, #F0EDE8); }
.lb-info { flex: 1; display: flex; flex-direction: column; }
.lb-name { font-size: 28rpx; color: var(--ink, #1E2A3A); font-weight: 500; }
.lb-level { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); }
.lb-points { font-size: 28rpx; color: var(--brand-primary, #C67A6A); font-weight: 600; }
.empty { text-align: center; padding: 40rpx; color: var(--ink-tertiary, #B8C2CE); font-size: 26rpx; }
</style>
