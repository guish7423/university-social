<template>
	<view class="page">
		<view class="tabs">
			<view class="tab" :class="{ active: tab === 'all' }" @click="switchTab('all')">全部</view>
			<view class="tab" :class="{ active: tab === 'lost' }" @click="switchTab('lost')">寻物</view>
			<view class="tab" :class="{ active: tab === 'found' }" @click="switchTab('found')">招领</view>
		</view>

		<view class="list" v-if="items.length">
			<view class="card" v-for="(item, idx) in items" :key="item.id" :class="'stagger-' + ((idx % 8) + 1)" @click="goDetail(item.id)">
				<view class="card-top">
					<view class="tag" :class="item.category === 'lost' ? 'tag-lost' : 'tag-found'">
						{{ item.category === 'lost' ? '寻物' : '招领' }}
					</view>
					<view class="status-tag" v-if="item.status !== 0">
						{{ item.status === 1 ? '已解决' : '已关闭' }}
					</view>
				</view>
				<text class="card-title">{{ item.title }}</text>
				<text class="card-desc" v-if="item.description">{{ item.description }}</text>
				<view class="card-footer">
					<text class="time">{{ item.created_at }}</text>
				</view>
			</view>
		</view>
		<u-empty text="暂无发布" mode="list" v-else-if="!loading"></u-empty>

		<u-loadmore :status="loadStatus" v-if="items.length" />
	</view>
</template>

<script>
import { myLostItems, listLostItems } from '@/api/found'

export default {
	data() {
		return {
			tab: 'all', items: [], loading: false,
			loadStatus: 'loadmore', offset: 0, limit: 20, hasMore: true,
		}
	},
	onLoad() { this.loadItems() },
	onReachBottom() { if (this.hasMore) this.loadItems() },
	methods: {
		getParams() {
			if (this.tab === 'all') return {}
			return { category: this.tab }
		},
		async loadItems() {
			if (this.loading || !this.hasMore) return
			this.loading = true
			this.loadStatus = 'loading'
			try {
				const category = this.getParams().category
				const { data } = category
					? await listLostItems(category, -1, this.offset, this.limit)
					: await myLostItems(this.offset, this.limit)
				if (this.offset === 0) this.items = data
				else this.items = [...this.items, ...data]
				this.hasMore = data.length >= this.limit
				this.offset += data.length
				this.loadStatus = this.hasMore ? 'loadmore' : 'nomore'
			} catch (e) { this.loadStatus = 'loadmore' }
			finally { this.loading = false }
		},
		switchTab(t) {
			this.tab = t; this.offset = 0; this.hasMore = true; this.items = []; this.loadItems()
		},
		goDetail(id) { uni.navigateTo({ url: '/pages/found/detail?id=' + id }) },
	}
}
</script>

<style lang="scss" scoped>
.page { min-height: 100vh; padding-bottom: 120rpx; padding: 0 24rpx; }
.tabs { flex-direction: row; padding: 20rpx 0; gap: 20rpx; }
.tab { font-size: 28rpx; color: #888; padding: 8rpx 24rpx; border-radius: 30rpx; background: #f5f5f5; }
.tab.active { background: linear-gradient(135deg, #667eea, #764ba2); color: #fff; }
.list { padding-top: 10rpx; }
.card {
	background: #fff; border-radius: 22rpx; padding: 28rpx;
	margin-bottom: 20rpx; box-shadow: 0 2rpx 12rpx rgba(0,0,0,.06);
}
.card-top { flex-direction: row; align-items: center; margin-bottom: 12rpx; }
.tag { font-size: 22rpx; padding: 4rpx 16rpx; border-radius: 8rpx; font-weight: 500; margin-right: 12rpx; }
.tag-lost { background: #fff1f0; color: #ff4d4f; }
.tag-found { background: #f0fff4; color: #52c41a; }
.status-tag { font-size: 22rpx; color: #999; background: #f5f5f5; padding: 4rpx 14rpx; border-radius: 8rpx; }
.card-title { font-size: 30rpx; font-weight: 600; color: #1a1a2e; display: block; margin-bottom: 8rpx; }
.card-desc { font-size: 26rpx; color: #666; display: block; margin-bottom: 12rpx; }
.card-footer { flex-direction: row; justify-content: flex-end; }
.time { font-size: 22rpx; color: #bbb; }
</style>
