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

function goSquare() { uni.switchTab({ url: '/pages/square/index' }) }
function goFriends() { uni.navigateTo({ url: '/pages/friend/index' }) }
function goCircles() { uni.navigateTo({ url: '/pages/circle/list' }) }
</script>

<template>
  <view class="container">
    <view class="profile-header">
      <view class="profile-bg" />
      <view class="profile-info">
        <u-avatar
          :src="userStore.avatar"
          :text="userStore.nickname?.[0] || '我'"
          size="120"
          shape="circle"
          fontSize="48"
          bgColor="#667eea"
        />
        <text class="nickname">{{ userStore.nickname || "未登录" }}</text>
        <view class="profile-meta">
          <u-tag v-if="userStore.isVerified" text="已认证" type="success" shape="circle" size="mini" />
          <text v-if="userStore.school" class="school">{{ userStore.school }}</text>
        </view>
      </view>
    </view>

    <view class="section">
      <u-cell-group>
        <u-cell-item
          title="学校信息"
          :value="userStore.school || '未设置'"
          :arrow="true"
          @click="openEdit"
        />
        },
        u-cell-item(
          title="编辑资料"
          :arrow="true"
          @click="uni.navigateTo({ url: '/pages/user/edit' })"
        )
      </u-cell-group>
    </view>

    <view class="section">
      <u-cell-group>
        <u-cell-item
          title="我的动态"
          :arrow="true"
          @click="goSquare"
        />
        <u-cell-item
          title="我的好友"
          :arrow="true"
          @click="goFriends"
        />
        <u-cell-item
          title="我的圈子"
          :arrow="true"
          @click="goCircles"
        />
      </u-cell-group>
    </view>

    <view class="section">
      <u-cell-group>
        <u-cell-item
          title="退出登录"
          :arrow="false"
          @click="handleLogout"
        >
          <template #title>
            <text style="color: #e74c3c;">退出登录</text>
          </template>
        </u-cell-item>
      </u-cell-group>
    </view>

    <u-popup :show="showEdit" mode="bottom" round="20" @close="showEdit = false">
      <view class="edit-popup">
        <text class="edit-title">设置学校</text>
        <u-input
          v-model="editSchool"
          placeholder="请输入学校名称"
          :border="true"
          shape="square"
          :customStyle="{ padding: '20rpx 24rpx', fontSize: '28rpx', height: '88rpx', marginBottom: '30rpx' }"
        />
        <view class="edit-actions">
          <u-button
            shape="circle"
            :customStyle="{ flex: 1, height: '88rpx' }"
            @click="showEdit = false"
          >取消</u-button>
          <u-button
            type="primary"
            shape="circle"
            :customStyle="{ flex: 1, height: '88rpx' }"
            @click="saveSchool"
          >保存</u-button>
        </view>
      </view>
    </u-popup>
  </view>
</template>

<style scoped>
.container {
  min-height: 100vh;
  background: var(--u-bg-color, #f3f4f6);
}

.profile-header {
  position: relative;
  padding-bottom: 40rpx;
}

.profile-bg {
  position: absolute;
  inset: 0;
  height: 280rpx;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  border-radius: 0 0 48rpx 48rpx;
}

.profile-info {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 60rpx;
}

.nickname {
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
  margin-top: 20rpx;
}

.profile-meta {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-top: 12rpx;
}

.school {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.7);
}

.section {
  margin: 20rpx 30rpx 0;
  border-radius: 20rpx;
  overflow: hidden;
}

.edit-popup {
  padding: 40rpx;
}

.edit-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #303133;
  display: block;
  margin-bottom: 32rpx;
  text-align: center;
}

.edit-actions {
  display: flex;
  gap: 20rpx;
}
</style>
