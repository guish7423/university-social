<template>
	<view class="page">
		<view class="header" v-if="item">
			<view class="header-top">
				<view class="tag" :class="item.category === 'lost' ? 'tag-lost' : 'tag-found'">
					{{ item.category === 'lost' ? '寻物' : '招领' }}
				</view>
				<view class="status-tag" v-if="item.status !== 0">
					{{ item.status === 1 ? '已解决' : '已关闭' }}
				</view>
			</view>
			<text class="title">{{ item.title }}</text>
			<text class="desc" v-if="item.description">{{ item.description }}</text>

			<view class="info-grid">
				<view class="info-item" v-if="item.item_type">
					<text class="info-label">物品类型</text>
					<text class="info-value">{{ item.item_type }}</text>
				</view>
				<view class="info-item" v-if="item.location">
					<text class="info-label">地点</text>
					<text class="info-value">{{ item.location }}</text>
				</view>
				<view class="info-item" v-if="item.contact">
					<text class="info-label">联系方式</text>
					<text class="info-value">{{ item.contact }}</text>
				</view>
			</view>

			<view class="footer">
				<view class="author">
					<u-avatar :src="item.avatar" size="40"></u-avatar>
					<text class="name">{{ item.nickname || '匿名' }}</text>
				</view>
				<text class="time">{{ item.created_at }}</text>
			</view>
		</view>

		<view class="actions" v-if="item && isOwner">
			<u-button type="warning" shape="circle" @click="markSolved" v-if="item.status === 0">{{ item.category === 'lost' ? '已找到' : '已完成' }}</u-button>
			<u-button type="error" shape="circle" plain @click="doDelete">删除</u-button>
		</view>
	</view>
</template>

<script>
import { getLostItem, deleteLostItem, updateLostItemStatus } from '@/api/found'
import { mapState } from 'pinia'
import { useUserStore } from '@/store/user'

export default {
	computed: { ...mapState(useUserStore, ['userInfo']), isOwner() { return this.item && this.userInfo && this.item.user_id === this.userInfo.id } },
	data() {
		return { item: null }
	},
	onLoad(opts) {
		if (opts.id) this.loadItem(parseInt(opts.id))
	},
	methods: {
		async loadItem(id) {
			try {
				const { data } = await getLostItem(id)
				this.item = data
			} catch (e) { uni.showToast({ title: '加载失败', icon: 'none' }) }
		},
		async markSolved() {
			try {
				await updateLostItemStatus(this.item.id, 1)
				uni.showToast({ title: '已更新' })
				this.loadItem(this.item.id)
			} catch (e) { uni.showToast({ title: '操作失败', icon: 'none' }) }
		},
		async doDelete() {
			uni.showModal({ title: '确认', content: '确定删除？', success: async (res) => {
				if (res.confirm) {
					try {
						await deleteLostItem(this.item.id)
						uni.showToast({ title: '已删除' })
						setTimeout(() => uni.navigateBack(), 1000)
					} catch (e) { uni.showToast({ title: '删除失败', icon: 'none' }) }
				}
			}})
		},
	}
}
</script>

<style lang="scss" scoped>
.page { padding: 24rpx; }
.header { background: #fff; border-radius: 22rpx; padding: 32rpx; margin-bottom: 24rpx; }
.header-top { flex-direction: row; align-items: center; margin-bottom: 16rpx; }
.tag { font-size: 22rpx; padding: 6rpx 20rpx; border-radius: 8rpx; font-weight: 500; margin-right: 12rpx; }
.tag-lost { background: #fff1f0; color: #ff4d4f; }
.tag-found { background: #f0fff4; color: #52c41a; }
.status-tag { font-size: 22rpx; color: #999; background: #f5f5f5; padding: 4rpx 14rpx; border-radius: 8rpx; }
.title { font-size: 34rpx; font-weight: 700; color: #1a1a2e; display: block; margin-bottom: 16rpx; }
.desc { font-size: 28rpx; color: #555; line-height: 1.6; display: block; margin-bottom: 24rpx; }
.info-grid { background: #f9f9fb; border-radius: 16rpx; padding: 20rpx; margin-bottom: 24rpx; }
.info-item { flex-direction: row; justify-content: space-between; padding: 12rpx 0; }
.info-label { font-size: 26rpx; color: #888; }
.info-value { font-size: 26rpx; color: #333; font-weight: 500; }
.footer { flex-direction: row; justify-content: space-between; align-items: center; }
.author { flex-direction: row; align-items: center; gap: 12rpx; }
.name { font-size: 26rpx; color: #5b6ef5; }
.time { font-size: 24rpx; color: #bbb; }
.actions { flex-direction: row; gap: 20rpx; }
</style>
