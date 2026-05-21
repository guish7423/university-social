<template>
	<view class="container">
		<view class="balance-card" :style="{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }">
			<text class="balance-label">我的积分</text>
			<text class="balance-value">{{ balance.points ?? 0 }}</text>
			<text class="balance-level">{{ balance.level ?? 'Lv.0' }}</text>
			<view v-if="balance.next_level" class="next-level">距 {{ balance.next_level }} 还需 {{ balance.next_threshold - (balance.points ?? 0) }} 分</view>
		</view>

		<view class="section">
			<view class="section-title">积分规则</view>
			<view class="rules">
				<view class="rule-item" v-for="(r, i) in rules" :key="i">
					<text class="rule-action">{{ r.action }}</text>
					<text class="rule-points">{{ r.points > 0 ? '+' : '' }}{{ r.points }}</text>
				</view>
			</view>
		</view>

		<view class="section">
			<view class="section-title" style="display:flex;justify-content:space-between;align-items:center">
				<text>积分排行</text>
				<text class="leaderboard-placeholder">更多</text>
			</view>
			<view class="leaderboard">
				<view class="lb-item" v-for="(e, i) in leaderboard" :key="e.user_id">
					<view class="lb-rank" :class="{ 'lb-gold': i===0, 'lb-silver': i===1, 'lb-bronze': i===2 }">{{ e.rank }}</view>
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
	</view>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getPointsBalance, claimDailyLogin, getPointsLeaderboard } from "@/api/points"

const balance = ref<any>({})
const leaderboard = ref<any[]>([])

const rules = [
	{ action: '每日登录', points: '+1' },
	{ action: '发布动态', points: '+3' },
	{ action: '评论', points: '+1' },
	{ action: '邀请注册', points: '+5' },
	{ action: '帖子加精', points: '+10' },
]

async function load() {
	try {
		balance.value = await getPointsBalance()
		leaderboard.value = (await getPointsLeaderboard()).leaderboard ?? []
	} catch {}
}
onMounted(load)
</script>

<style scoped lang="scss">
.container { min-height: 100vh; background: #f5f5f5; padding-bottom: 40rpx; }
.balance-card {
	margin: 30rpx; padding: 40rpx; border-radius: 24rpx; display: flex; flex-direction: column; align-items: center; color: #fff;
	.balance-label { font-size: 28rpx; opacity: 0.9; }
	.balance-value { font-size: 72rpx; font-weight: 700; margin: 10rpx 0; }
	.balance-level { font-size: 32rpx; font-weight: 500; }
	.next-level { font-size: 24rpx; opacity: 0.8; margin-top: 16rpx; }
}
.section { margin: 30rpx; }
.section-title { font-size: 30rpx; font-weight: 600; margin-bottom: 20rpx; }
.rules {
	background: #fff; border-radius: 20rpx; padding: 20rpx;
	.rule-item { display: flex; justify-content: space-between; padding: 16rpx 0; border-bottom: 1rpx solid #f0f0f0; }
	.rule-item:last-child { border-bottom: none; }
	.rule-action { font-size: 28rpx; color: #333; }
	.rule-points { font-size: 28rpx; color: #667eea; font-weight: 600; }
}
.leaderboard-placeholder { font-size: 24rpx; color: #999; }
.leaderboard {
	background: #fff; border-radius: 20rpx; padding: 20rpx;
	.lb-item { display: flex; align-items: center; padding: 16rpx 0; border-bottom: 1rpx solid #f0f0f0; }
	.lb-item:last-child { border-bottom: none; }
	.lb-rank { width: 48rpx; height: 48rpx; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 26rpx; font-weight: 700; color: #999; background: #f5f5f5; }
	.lb-gold { background: linear-gradient(135deg, #f6d365, #fda085); color: #fff; }
	.lb-silver { background: linear-gradient(135deg, #e0e0e0, #b0b0b0); color: #fff; }
	.lb-bronze { background: linear-gradient(135deg, #e6a87a, #c7794a); color: #fff; }
	.lb-avatar { width: 56rpx; height: 56rpx; border-radius: 50%; margin: 0 16rpx; }
	.lb-info { flex: 1; display: flex; flex-direction: column; }
	.lb-name { font-size: 28rpx; color: #333; font-weight: 500; }
	.lb-level { font-size: 22rpx; color: #999; }
	.lb-points { font-size: 28rpx; color: #667eea; font-weight: 600; }
}
.empty { text-align: center; padding: 40rpx; color: #999; font-size: 26rpx; }
</style>
