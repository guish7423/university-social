<template>
	<view class="page">
		<u-tabs :list="tabs" :scrollable="false" :current="currentTab" @change="onTabChange" active-color="#5b6ef5"></u-tabs>

		<view class="toolbar">
			<u-search placeholder="搜索物品" v-model="keyword" @search="onSearch" @clear="onClear" :show-action="false" bg-color="#f5f5f5" shape="round"></u-search>
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
				<view class="card-meta">
					<text class="meta-item" v-if="item.item_type">📦 {{ item.item_type }}</text>
					<text class="meta-item" v-if="item.location">📍 {{ item.location }}</text>
				</view>
				<view class="card-footer">
					<text class="author">{{ item.nickname || '匿名' }}</text>
					<text class="time">{{ formatTime(item.created_at) }}</text>
				</view>
			</view>
		</view>

		<u-empty text="暂无内容" mode="list" v-else-if="!loading"></u-empty>
		<u-loadmore :status="loadStatus" v-if="items.length" />

		<u-toast ref="toast"></u-toast>
	</view>
</template>

<script>
export default {
	data() {
		return {
			tabs: [{ name: '全部' }, { name: '寻物' }, { name: '招领' }],
			currentTab: 0,
			keyword: '',
			items: [],
			loading: false,
			loadStatus: 'loadmore',
			offset: 0,
			limit: 20,
			hasMore: true,
		}
	},
	onLoad() {
		this.loadItems()
	},
	onReachBottom() {
		if (this.hasMore) this.loadItems()
	},
	methods: {
		getCategory() {
			if (this.currentTab === 0) return ''
			return this.currentTab === 1 ? 'lost' : 'found'
		},
		async loadItems() {
			if (this.loading || !this.hasMore) return
			this.loading = true
			this.loadStatus = 'loading'
			try {
				const { data } = this.keyword
					? await searchLostItems(this.keyword, this.getCategory(), this.offset, this.limit)
					: await listLostItems(this.getCategory(), -1, this.offset, this.limit)
				if (this.offset === 0) this.items = data
				else this.items = [...this.items, ...data]
				this.hasMore = data.length >= this.limit
				this.offset += data.length
				this.loadStatus = this.hasMore ? 'loadmore' : 'nomore'
			} catch (e) {
				this.loadStatus = 'loadmore'
			} finally {
				this.loading = false
			}
		},
		onTabChange(idx) {
			this.currentTab = idx
			this.offset = 0
			this.hasMore = true
			this.items = []
			this.loadItems()
		},
		onSearch() {
			this.offset = 0
			this.hasMore = true
			this.items = []
			this.loadItems()
		},
		onClear() {
			this.keyword = ''
			this.onSearch()
		},
		goDetail(id) { uni.navigateTo({ url: '/pages/found/detail?id=' + id }) },
		formatTime(t) {
			if (!t) return ''
			const d = new Date(t)
			return `${d.getMonth()+1}.${d.getDate()}`
		}
	}
}
</script>

<style lang="scss" scoped>
.page {
	min-height: 100vh;
	padding-bottom: 120rpx;
}
.toolbar { padding: 20rpx 24rpx; }
.list { padding: 0 24rpx; }
.card {
	background: #fff;
	border-radius: 22rpx;
	padding: 28rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 2rpx 12rpx rgba(0,0,0,.06);
	animation: fadeInUp .5s ease-out both;
}
.card-top {
	flex-direction: row; align-items: center; margin-bottom: 12rpx;
}
.tag {
	font-size: 22rpx; padding: 4rpx 16rpx; border-radius: 8rpx; font-weight: 500;
	margin-right: 12rpx;
}
.tag-lost { background: #fff1f0; color: #ff4d4f; }
.tag-found { background: #f0fff4; color: #52c41a; }
.status-tag {
	font-size: 22rpx; color: #999; background: #f5f5f5;
	padding: 4rpx 14rpx; border-radius: 8rpx;
}
.card-title { font-size: 30rpx; font-weight: 600; color: #1a1a2e; display: block; margin-bottom: 8rpx; }
.card-desc { font-size: 26rpx; color: #666; display: block; margin-bottom: 12rpx; line-height: 1.5; }
.card-meta { flex-direction: row; margin-bottom: 16rpx; gap: 20rpx; }
.meta-item { font-size: 24rpx; color: #888; }
.card-footer { flex-direction: row; justify-content: space-between; align-items: center; }
.author { font-size: 24rpx; color: #5b6ef5; }
.time { font-size: 22rpx; color: #bbb; }

@keyframes fadeInUp {
	from { opacity: 0; transform: translateY(20rpx); }
	to { opacity: 1; transform: translateY(0); }
}
.stagger-1 { animation-delay: 0ms; }
.stagger-2 { animation-delay: 60ms; }
.stagger-3 { animation-delay: 120ms; }
.stagger-4 { animation-delay: 180ms; }
.stagger-5 { animation-delay: 240ms; }
.stagger-6 { animation-delay: 300ms; }
.stagger-7 { animation-delay: 360ms; }
.stagger-8 { animation-delay: 420ms; }
</style>
