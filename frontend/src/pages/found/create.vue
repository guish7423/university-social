<template>
	<view class="page">
		<view class="form">
			<view class="form-item">
				<text class="label">类型 <text class="req">*</text></text>
				<view class="type-picker">
					<view class="type-option" :class="{ active: form.category === 'lost' }" @click="form.category = 'lost'">
						<text>😢 寻物</text>
					</view>
					<view class="type-option" :class="{ active: form.category === 'found' }" @click="form.category = 'found'">
						<text>🙌 招领</text>
					</view>
				</view>
			</view>
			<u-input placeholder="物品名称（必填）" v-model="form.title" border="none" inputClass="custom-input" />
			<u-input type="textarea" placeholder="详细描述（时间、特征等）" v-model="form.description" border="none" height="100" />
			<u-input placeholder="物品类型（如：校园卡、耳机）" v-model="form.item_type" border="none" inputClass="custom-input" />
			<u-input placeholder="地点（如：图书馆3楼）" v-model="form.location" border="none" inputClass="custom-input" />
			<u-input placeholder="联系方式" v-model="form.contact" border="none" inputClass="custom-input" />

			<view class="form-item">
				<text class="label">图片</text>
				<u-upload :action="uploadUrl" :auto-upload="true" max-count="6" @on-uploaded="onUploaded"></u-upload>
			</view>

			<u-button type="primary" shape="circle" :disabled="!form.title || !form.category" :loading="submitting" @click="submit">
				发布
			</u-button>
		</view>
		<u-toast ref="toast"></u-toast>
	</view>
</template>

<script>
import { createLostItem } from '@/api/found'

export default {
	data() {
		return {
			form: { title: '', description: '', category: 'lost', item_type: '', location: '', contact: '', images: [] },
			submitting: false,
			uploadUrl: process.env.VUE_APP_API_BASE + '/api/v1/upload',
		}
	},
	methods: {
		onUploaded(files) {
			this.form.images = files.map(f => f.response?.data?.url || '').filter(Boolean)
		},
		async submit() {
			if (!this.form.title || !this.form.category) return
			this.submitting = true
			try {
				await createLostItem(this.form)
				this.$refs.toast.show({ type: 'success', message: '发布成功' })
				setTimeout(() => uni.navigateBack(), 1200)
			} catch (e) {
				this.$refs.toast.show({ type: 'error', message: '发布失败' })
			} finally {
				this.submitting = false
			}
		},
	}
}
</script>

<style lang="scss" scoped>
.page { padding: 24rpx; }
.form { background: #fff; border-radius: 22rpx; padding: 32rpx; }
.form-item { margin-bottom: 24rpx; }
.label { font-size: 28rpx; font-weight: 600; color: #1a1a2e; display: block; margin-bottom: 12rpx; }
.req { color: #ff4d4f; }
.type-picker { flex-direction: row; gap: 20rpx; }
.type-option {
	flex: 1; padding: 24rpx; border-radius: 16rpx; text-align: center;
	background: #f5f5f5; font-size: 28rpx;
	&.active {
		background: linear-gradient(135deg, #667eea, #764ba2);
		color: #fff;
	}
}
.custom-input { font-size: 28rpx; padding: 20rpx 0; border-bottom: 2rpx solid #f0f0f0; margin-bottom: 16rpx; }
</style>
