<script setup lang="ts">
import { ref } from "vue"
import { createPost, listTopics, uploadImages, type TopicData } from "@/api/post"

const content = ref("")
const topics = ref<TopicData[]>([])
const selectedTopic = ref<number>()
const submitting = ref(false)
const imageFiles = ref<string[]>([])

async function loadTopics() {
  try { topics.value = await listTopics() } catch {}
}
loadTopics()

async function chooseImage() {
  if (imageFiles.value.length >= 9) {
    uni.showToast({ title: "最多9张图片", icon: "none" })
    return
  }
  try {
    const res = await uni.chooseImage({ count: 9 - imageFiles.value.length })
    imageFiles.value = [...imageFiles.value, ...res.tempFilePaths]
  } catch {}
}

function removeImage(index: number) {
  imageFiles.value.splice(index, 1)
}

async function handleSubmit() {
  if (!content.value.trim()) {
    uni.showToast({ title: "写点内容吧", icon: "none" })
    return
  }
  submitting.value = true
  let images: string[] = []
  if (imageFiles.value.length > 0) {
    uni.showLoading({ title: "上传图片中..." })
    try {
      images = await uploadImages(imageFiles.value)
    } catch {
      uni.hideLoading()
      submitting.value = false
      uni.showToast({ title: "图片上传失败", icon: "none" })
      return
    }
    uni.hideLoading()
  }
  try {
    await createPost({
      content: content.value,
      images,
      topic_id: selectedTopic.value,
    })
    uni.$emit('post-created')
    uni.switchTab({ url: '/pages/square/index' })
  } catch {
    uni.showToast({ title: "发布失败", icon: "none" })
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <view class="container">
    <view class="form-card">
      <u-textarea
        v-model="content"
        placeholder="分享你的校园生活..."
        maxlength="2000"
        height="300"
        count
        :autoHeight="false"
        :customStyle="{ fontSize: '30rpx', padding: '24rpx', borderRadius: '16rpx', background: '#f8f9fa' }"
      />

      <view class="image-section">
        <text class="section-label">添加图片</text>
        <view class="image-grid">
          <view
            v-for="(img, idx) in imageFiles"
            :key="idx"
            class="image-item"
          >
            <image :src="img" mode="aspectFill" class="preview-img" />
            <view class="remove-btn" @click="removeImage(idx)">&times;</view>
          </view>
          <view
            v-if="imageFiles.length < 9"
            class="image-picker"
            @click="chooseImage"
          >
            <text class="picker-icon">+</text>
            <text class="picker-text">添加图片</text>
          </view>
        </view>
      </view>

      <view v-if="topics?.length" class="topic-section">
        <text class="section-label">选择话题</text>
        <scroll-view class="topic-list" scroll-x>
          <view
            v-for="t in topics"
            :key="t.id"
            :class="['topic-tag', selectedTopic === t.id && 'active']"
            @click="selectedTopic = selectedTopic === t.id ? undefined : t.id"
          >
            {{ t.name }}
          </view>
        </scroll-view>
      </view>

      <view class="submit-area">
        <u-button
          type="primary"
          shape="circle"
          size="large"
          :loading="submitting"
          :disabled="submitting || !content.trim()"
          :customStyle="{ height: '88rpx', fontSize: '32rpx', fontWeight: '600' }"
          @click="handleSubmit"
        >
          发布动态
        </u-button>
      </view>
    </view>
  </view>
</template>

<style scoped>
.container {
  min-height: 100vh;
  background: var(--color-canvas, #f8f6f3);
  padding: 20rpx;
}

.form-card {
  background: var(--color-surface, #fefcfb);
  border-radius: 24rpx;
  padding: 30rpx;
  box-shadow: 0 2rpx 16rpx rgba(#2d2b28, 0.04);
}

.image-section {
  margin-top: 24rpx;
}

.image-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
  margin-top: 12rpx;
}

.image-item {
  position: relative;
  width: 180rpx;
  height: 180rpx;
  border-radius: 12rpx;
  overflow: hidden;
}

.preview-img {
  width: 100%;
  height: 100%;
}

.remove-btn {
  position: absolute;
  top: 4rpx;
  right: 4rpx;
  width: 36rpx;
  height: 36rpx;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  font-size: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-picker {
  width: 180rpx;
  height: 180rpx;
  border: 2rpx dashed var(--ink-tertiary, #c4c0ba);
  border-radius: 12rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--color-surface-1, #fcfaf7);
}

.picker-icon {
  font-size: 48rpx;
  color: #999;
  line-height: 1;
}

.picker-text {
  font-size: 22rpx;
  color: #999;
  margin-top: 8rpx;
}

.section-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #606266;
  display: block;
}

.topic-section {
  margin-top: 24rpx;
}

.topic-list {
  display: flex;
  gap: 16rpx;
  white-space: nowrap;
  margin-top: 12rpx;
  margin-bottom: 32rpx;
}

.topic-tag {
  display: inline-block;
  padding: 14rpx 32rpx;
  background: #f0f0f0;
  border-radius: 40rpx;
  font-size: 26rpx;
  color: #606266;
  transition: all 0.2s;
}

.topic-tag.active {
  background: var(--brand-primary, #FF6B6B);
  color: #fff;
}

.submit-area {
  margin-top: 8rpx;
}
</style>
