<script setup lang="ts">
import { ref } from "vue"
import { useUserStore } from "@/store/user"
import { getProfile, updateProfile } from "@/api/user"

const userStore = useUserStore()
const showEdit = ref(false)
const editSchool = ref(userStore.school)

async function loadProfile() {
  try {
    const profile = await getProfile()
    userStore.setUserInfo(profile)
  } catch {}
}
loadProfile()

function openEdit() {
  editSchool.value = userStore.school
  showEdit.value = true
}

async function saveSchool() {
  await updateProfile({ school: editSchool.value })
  userStore.school = editSchool.value
  uni.setStorageSync("school", editSchool.value)
  showEdit.value = false
  uni.showToast({ title: "保存成功", icon: "success" })
}

function handleLogout() {
  uni.showModal({
    title: "提示",
    content: "确定退出登录？",
    success: (res) => {
      if (res.confirm) {
        userStore.logout()
        uni.reLaunch({ url: "/pages/login/index" })
      }
    },
  })
}
</script>

<template>
  <view class="container">
    <view class="profile-card">
      <image v-if="userStore.avatar" class="avatar" :src="userStore.avatar" mode="aspectFill" />
      <view v-else class="avatar-placeholder">{{ userStore.nickname?.[0] || "我" }}</view>
      <text class="nickname">{{ userStore.nickname || "未登录" }}</text>
      <view v-if="userStore.isVerified" class="verified-badge">已认证</view>
      <text v-if="userStore.school" class="school">{{ userStore.school }}</text>
    </view>

    <view class="menu-list">
      <view class="menu-item" @click="openEdit">
        <text>学校信息</text>
        <text class="menu-right">
          <text class="menu-value">{{ userStore.school || "未设置" }}</text>
          <text class="arrow">›</text>
        </text>
      </view>
      <view class="menu-item" @click="uni.switchTab({ url: '/pages/square/index' })">
        <text>我的动态</text>
        <text class="arrow">›</text>
      </view>
      <view class="menu-item" @click="uni.navigateTo({ url: '/pages/friend/index' })">
        <text>我的好友</text>
        <text class="arrow">›</text>
      </view>
      <view class="menu-item" @click="uni.navigateTo({ url: '/pages/circle/list' })">
        <text>我的圈子</text>
        <text class="arrow">›</text>
      </view>
      <view class="menu-item" @click="handleLogout">
        <text class="logout-text">退出登录</text>
      </view>
    </view>

    <uni-popup v-model="showEdit" type="bottom">
      <view class="edit-popup">
        <text class="edit-title">设置学校</text>
        <input v-model="editSchool" class="edit-input" placeholder="请输入学校名称" />
        <view class="edit-actions">
          <button @click="showEdit = false" class="cancel-btn">取消</button>
          <button @click="saveSchool" type="primary" class="save-btn">保存</button>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
}
.profile-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60rpx 0;
  background: #fff;
}
.avatar, .avatar-placeholder {
  width: 140rpx;
  height: 140rpx;
  border-radius: 50%;
}
.avatar-placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  font-size: 48rpx;
}
.nickname {
  margin-top: 20rpx;
  font-size: 36rpx;
  font-weight: 600;
}
.verified-badge {
  margin-top: 10rpx;
  padding: 4rpx 16rpx;
  background: #e8f5e9;
  color: #2e7d32;
  border-radius: 20rpx;
  font-size: 22rpx;
}
.school {
  margin-top: 8rpx;
  font-size: 26rpx;
  color: #666;
}
.menu-list {
  margin-top: 20rpx;
  background: #fff;
}
.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  border-bottom: 1rpx solid #eee;
  font-size: 28rpx;
}
.menu-right {
  display: flex;
  align-items: center;
  gap: 10rpx;
}
.menu-value {
  color: #999;
  font-size: 26rpx;
}
.arrow {
  color: #ccc;
  font-size: 36rpx;
}
.logout-text {
  color: #e74c3c;
}
.edit-popup {
  padding: 40rpx;
  background: #fff;
  border-radius: 20rpx 20rpx 0 0;
}
.edit-title {
  font-size: 32rpx;
  font-weight: 600;
  margin-bottom: 30rpx;
}
.edit-input {
  border: 1rpx solid #ddd;
  border-radius: 12rpx;
  padding: 20rpx;
  font-size: 28rpx;
  margin-bottom: 30rpx;
}
.edit-actions {
  display: flex;
  gap: 20rpx;
}
.cancel-btn, .save-btn {
  flex: 1;
}
</style>
