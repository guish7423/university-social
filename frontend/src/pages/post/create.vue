<script setup lang="ts">
import { ref } from "vue"
import { createPost, listTopics, type TopicData } from "@/api/post"

const content = ref("")
const topics = ref<TopicData[]>([])
const selectedTopic = ref<number>()
const submitting = ref(false)

async function loadTopics() {
  try { topics.value = await listTopics() } catch {}
}
loadTopics()

async function handleSubmit() {
  if (!content.value.trim()) {
    uni.showToast({ title: "写点内容吧", icon: "none" })
    return
  }
  submitting.value = true
  try {
    await createPost({
      content: content.value,
      topic_id: selectedTopic.value,
    })
    uni.showToast({ title: "发布成功", icon: "success" })
    setTimeout(() => uni.navigateBack(), 1000)
  } catch {
    uni.showToast({ title: "发布失败", icon: "none" })
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <view class="container">
    <view class="form">
      <textarea
        v-model="content"
        class="textarea"
        placeholder="分享你的校园生活..."
        maxlength="2000"
      />
      <view class="char-count">{{ content.length }}/2000</view>

      <view v-if="topics.length" class="topic-section">
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

      <button
        class="submit-btn"
        type="primary"
        :loading="submitting"
        :disabled="submitting || !content.trim()"
        @click="handleSubmit"
      >
        发布
      </button>
    </view>
  </view>
</template>

<style scoped>
.container { min-height: 100vh; background: #fff; padding: 30rpx; }
.textarea { width: 100%; height: 400rpx; font-size: 30rpx; line-height: 1.6; }
.char-count { text-align: right; font-size: 24rpx; color: #999; margin-bottom: 30rpx; }
.section-label { font-size: 28rpx; color: #666; margin-bottom: 16rpx; display: block; }
.topic-list { display: flex; gap: 16rpx; white-space: nowrap; margin-bottom: 40rpx; }
.topic-tag { display: inline-block; padding: 12rpx 28rpx; background: #f0f0f0; border-radius: 30rpx; font-size: 26rpx; }
.topic-tag.active { background: #667eea; color: #fff; }
.submit-btn { width: 100%; height: 88rpx; border-radius: 44rpx; font-size: 32rpx; }
</style>
