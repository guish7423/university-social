<script setup lang="ts">
import { ref } from "vue"
import { createActivity } from "@/api/activity"

const title = ref("")
const description = ref("")
const activityType = ref("竞赛组队")
const location = ref("")
const maxParticipants = ref("")
const startTime = ref("")
const endTime = ref("")
const images = ref<string[]>([])
const submitting = ref(false)

const types = ["竞赛组队", "约球运动", "学习局", "桌游局", "聚餐", "其他"]
const showTypePicker = ref(false)

function openTypePicker() {
  uni.showActionSheet({
    itemList: types,
    success: (res) => {
      activityType.value = types[res.tapIndex]
    }
  })
}

function pickDate(cb: (val: string) => void) {
  uni.showDatePicker?.({
    success: (res: any) => cb(res.value)
  })
}

function pickStart() {
  pickDate((val) => { startTime.value = val })
}

function pickEnd() {
  pickDate((val) => { endTime.value = val })
}

async function handleSubmit() {
  if (!title.value.trim()) { uni.showToast({ title: "请输入活动标题", icon: "none" }); return }
  if (!startTime.value) { uni.showToast({ title: "请选择开始时间", icon: "none" }); return }

  submitting.value = true
  try {
    const data: any = {
      title: title.value.trim(),
      description: description.value.trim(),
      activity_type: activityType.value,
      location: location.value.trim(),
      start_time: startTime.value,
      images: images.value,
    }
    if (maxParticipants.value) data.max_participants = parseInt(maxParticipants.value)
    if (endTime.value) data.end_time = endTime.value

    await createActivity(data)
    uni.showToast({ title: "创建成功", icon: "success" })
    setTimeout(() => uni.navigateBack(), 1000)
  } catch {
    uni.showToast({ title: "创建失败", icon: "none" })
  }
  submitting.value = false
}

function chooseImage() {
  uni.chooseImage({
    count: 9,
    success(res) {
      images.value = images.value.concat(res.tempFilePaths).slice(0, 9)
    },
  })
}

function removeImage(idx: number) {
  images.value.splice(idx, 1)
}
</script>

<template>
  <view class="container">
    <view class="form-group">
      <text class="form-label">活动标题 <text class="required">*</text></text>
      <input v-model="title" class="form-input" placeholder="请输入活动标题" maxlength="50" />
    </view>

    <view class="form-group">
      <text class="form-label">活动类型 <text class="required">*</text></text>
      <view class="form-select" @click="openTypePicker">
        <text>{{ activityType }}</text>
        <text class="arrow">▾</text>
      </view>
    </view>

    <view class="form-group">
      <text class="form-label">地点</text>
      <input v-model="location" class="form-input" placeholder="活动地点（可选）" />
    </view>

    <view class="form-row">
      <view class="form-group flex-1">
        <text class="form-label">开始时间 <text class="required">*</text></text>
        <view class="form-select" @click="pickStart">
          <text>{{ startTime ? startTime.slice(0, 16).replace('T', ' ') : '请选择' }}</text>
          <text class="arrow">▾</text>
        </view>
      </view>
      <view class="form-group flex-1">
        <text class="form-label">结束时间</text>
        <view class="form-select" @click="pickEnd">
          <text>{{ endTime ? endTime.slice(0, 16).replace('T', ' ') : '可选' }}</text>
          <text class="arrow">▾</text>
        </view>
      </view>
    </view>

    <view class="form-group">
      <text class="form-label">人数上限</text>
      <input v-model="maxParticipants" class="form-input" placeholder="不限制留空" type="digit" />
    </view>

    <view class="form-group">
      <text class="form-label">描述</text>
      <textarea v-model="description" class="form-textarea" placeholder="活动详情、注意事项等" maxlength="500" />
    </view>

    <view class="form-group">
      <text class="form-label">活动图片</text>
      <view class="image-grid">
        <view v-for="(img, idx) in images" :key="idx" class="image-item">
          <image :src="img" class="preview-img" mode="aspectFill" />
          <text class="remove-img" @click="removeImage(idx)">×</text>
        </view>
        <view v-if="images.length < 9" class="image-add" @click="chooseImage">
          <text class="add-icon">+</text>
        </view>
      </view>
    </view>

    <view class="submit-bar">
      <view class="submit-btn" :class="{ disabled: submitting }" @click="handleSubmit">
        <text>{{ submitting ? '发布中...' : '发布活动' }}</text>
      </view>
    </view>
  </view>
</template>

<style lang="scss">
page { background: #f5f7fa; }
.container { padding: 20rpx; padding-bottom: 160rpx; }

.form-group { margin-bottom: 24rpx; }
.form-label {
  font-size: 26rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 12rpx;
  display: block;
}
.required { color: #e74c3c; }
.form-input {
  height: 72rpx;
  background: #fff;
  border-radius: 12rpx;
  padding: 0 20rpx;
  font-size: 28rpx;
  color: #333;
  border: 1rpx solid #eee;
}
.form-textarea {
  width: 100%;
  height: 160rpx;
  background: #fff;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  font-size: 28rpx;
  color: #333;
  border: 1rpx solid #eee;
  box-sizing: border-box;
}
.form-select {
  height: 72rpx;
  background: #fff;
  border-radius: 12rpx;
  padding: 0 20rpx;
  font-size: 28rpx;
  color: #333;
  border: 1rpx solid #eee;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.arrow { color: #999; font-size: 24rpx; }
.form-row { display: flex; gap: 20rpx; }
.flex-1 { flex: 1; }

.image-grid { display: flex; flex-wrap: wrap; gap: 16rpx; }
.image-item { position: relative; }
.preview-img {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
}
.remove-img {
  position: absolute;
  top: -8rpx;
  right: -8rpx;
  width: 36rpx;
  height: 36rpx;
  border-radius: 50%;
  background: rgba(0,0,0,0.5);
  color: #fff;
  text-align: center;
  line-height: 36rpx;
  font-size: 24rpx;
}
.image-add {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  border: 2rpx dashed #ddd;
  display: flex;
  align-items: center;
  justify-content: center;
}
.add-icon { font-size: 48rpx; color: #ccc; }

.submit-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: #fff;
  box-shadow: 0 -2rpx 12rpx rgba(0,0,0,0.06);
}
.submit-btn {
  height: 88rpx;
  line-height: 88rpx;
  text-align: center;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  border-radius: 44rpx;
  font-size: 32rpx;
  font-weight: 600;
  &.disabled { opacity: 0.5; }
}
</style>
