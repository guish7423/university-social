<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useUserStore } from "@/store/user"
import { getProfile, updateProfile } from "@/api/user"
import { uploadImage } from "@/api/post"

const userStore = useUserStore()
const nickname = ref("")
const school = ref("")
const avatar = ref("")

onMounted(() => {
  nickname.value = userStore.nickname
  school.value = userStore.school
  avatar.value = userStore.avatar
})

async function pickAvatar() {
  try {
    const res = await uni.chooseImage({ count: 1, sizeType: ["compressed"] })
    uni.showLoading({ title: "上传头像..." })
    const url = await uploadImage(res.tempFilePaths[0])
    avatar.value = url
    uni.hideLoading()
  } catch {
    uni.hideLoading()
  }
}

async function handleSave() {
  if (!nickname.value.trim()) {
    uni.showToast({ title: "昵称不能为空", icon: "none" })
    return
  }
  uni.showLoading({ title: "保存中..." })
  try {
    await updateProfile({
      nickname: nickname.value.trim(),
      school: school.value.trim(),
      avatar: avatar.value,
    })
    userStore.nickname = nickname.value.trim()
    userStore.school = school.value.trim()
    userStore.avatar = avatar.value
    uni.setStorageSync("nickname", nickname.value.trim())
    uni.setStorageSync("school", school.value.trim())
    uni.setStorageSync("avatar", avatar.value)
    uni.hideLoading()
    uni.showToast({ title: "保存成功", icon: "success" })
    uni.navigateBack()
  } catch {
    uni.hideLoading()
    uni.showToast({ title: "保存失败", icon: "none" })
  }
}
</script>

<template>
  <view class="container">
    <view class="form-card">
      <view class="avatar-section" @click="pickAvatar">
        <u-avatar
          :src="avatar"
          :text="nickname?.[0] || '?'"
          size="140"
          shape="circle"
          fontSize="56"
          bgColor="#667eea"
        />
        <view class="avatar-overlay">
          <text class="camera-icon">📷</text>
          <text class="camera-text">更换头像</text>
        </view>
      </view>

      <view class="field">
        <text class="label">昵称</text>
        <u-input
          v-model="nickname"
          placeholder="请输入昵称"
          maxlength="20"
          :border="true"
          shape="square"
          :customStyle="{ padding: '20rpx 24rpx', fontSize: '28rpx', height: '88rpx' }"
        />
      </view>

      <view class="field">
        <text class="label">学校</text>
        <u-input
          v-model="school"
          placeholder="请输入学校名称"
          maxlength="50"
          :border="true"
          shape="square"
          :customStyle="{ padding: '20rpx 24rpx', fontSize: '28rpx', height: '88rpx' }"
        />
      </view>
    </view>

    <view class="submit-area">
      <u-button
        type="primary"
        shape="circle"
        size="large"
        :customStyle="{ height: '88rpx', fontSize: '32rpx', fontWeight: '600' }"
        @click="handleSave"
      >
        保存
      </u-button>
    </view>
  </view>
</template>

<style scoped>
.container {
  min-height: 100vh;
  background: var(--u-bg-color, #f3f4f6);
  padding: 30rpx;
}

.form-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 40rpx;
  box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.04);
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30rpx 0 40rpx;
  position: relative;
}

.avatar-overlay {
  position: absolute;
  bottom: 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4rpx;
}

.camera-icon {
  font-size: 40rpx;
}

.camera-text {
  font-size: 22rpx;
  color: #667eea;
}

.field {
  margin-bottom: 30rpx;
}

.label {
  display: block;
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12rpx;
}

.submit-area {
  margin-top: 40rpx;
  padding: 0 10rpx;
}
</style>
