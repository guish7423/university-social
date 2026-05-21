<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useUserStore } from "@/store/user"
import { getProfile, updateProfile } from "@/api/user"
import { getFollowCounts } from "@/api/follow"
import { userPosts, type PostData } from "@/api/post"

const userStore = useUserStore()
const showEdit = ref(false)
const editSchool = ref(userStore.school)
const followCounts = ref({ following_count: 0, followers_count: 0 })
const myPosts = ref<PostData[]>([])

async function loadProfile() {
  try {
    const profile = await getProfile()
    userStore.setUserInfo(profile)
    const counts = await getFollowCounts(profile.id)
    followCounts.value = counts
  } catch {}
}
loadProfile()
async function loadMyPosts() {
	try {
		myPosts.value = await userPosts(userStore.id || 0, 0, 3)
	} catch {}
}
loadProfile()
loadMyPosts()
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
function goFollowing() { uni.navigateTo({ url: '/pages/user/following?id=' + userStore.id }) }
function goFollowers() { uni.navigateTo({ url: '/pages/user/followers?id=' + userStore.id }) }
function goMyPosts() { uni.navigateTo({ url: '/pages/user/posts' }) }
function goMyProducts() { uni.navigateTo({ url: '/pages/product/mine' }) }
function goMyActivities() { uni.navigateTo({ url: '/pages/activity/mine' }) }
function goPoints() { uni.navigateTo({ url: '/pages/points/index' }) }
function goMyLostItems() { uni.navigateTo({ url: '/pages/found/mine' }) }

</script>

<template>
  <view class="container">
    <view class="profile-header">
      <view class="profile-bg" />
      <view class="profile-deco profile-deco-1" />
      <view class="profile-deco profile-deco-2" />
      <view class="profile-info animate-fade-in-up">
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
      <view class="follow-stats animate-fade-in-up stagger-1">
        <view class="stat-item" @click="goFollowers">
          <text class="stat-num">{{ followCounts.followers_count }}</text>
          <text class="stat-label">粉丝</text>
        </view>
        <view class="stat-divider" />
        <view class="stat-item" @click="goFollowing">
          <text class="stat-num">{{ followCounts.following_count }}</text>
          <text class="stat-label">关注</text>
        </view>
      </view>
    </view>

    <view class="section animate-fade-in-up stagger-1">
      <u-cell-group>
        <u-cell-item
          title="学校信息"
          :value="userStore.school || '未设置'"
          :arrow="true"
          @click="openEdit"
        />
        <u-cell-item
          title="校园认证"
          :value="userStore.isVerified ? '已认证' : '未认证'"
          :arrow="true"
          @click="uni.navigateTo({ url: '/pages/verification/index' })"
        />
        <u-cell-item
          title="编辑资料"
          :arrow="true"
          @click="uni.navigateTo({ url: '/pages/user/edit' })"
        />
      </u-cell-group>
    </view>

    <view class="section animate-fade-in-up stagger-2">
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
			<u-cell-item
				title="我的发布"
				:arrow="true"
				@click="goMyProducts"
			/>
      <u-cell-item
          title="我的活动"
          :arrow="true"
          @click="goMyActivities"
      />
      <u-cell-item
          title="失物招领"
          :arrow="true"
          @click="goMyLostItems"
      />

      <u-cell-item
          title="我的积分"
          :arrow="true"
          @click="goPoints"
      />
		</u-cell-group>
	</view>

    <view class="section animate-fade-in-up stagger-3">
      <view class="section-title">
        <text class="section-title-text">我的帖子</text>
        <text class="section-title-more" @click="goMyPosts">全部</text>
      </view>
      <view v-if="myPosts.length" class="post-preview-list">
        <view v-for="post in myPosts" :key="post.id" class="post-preview-item" @click="uni.navigateTo({ url: '/pages/post/detail?id=' + post.id })">
          <text class="post-preview-content">{{ post.content.slice(0, 60) }}{{ post.content.length > 60 ? '...' : '' }}</text>
          <view class="post-preview-meta">
            <text class="post-preview-stat">{{ post.like_count }} 赞</text>
            <text class="post-preview-stat">{{ post.comment_count }} 评论</text>
          </view>
        </view>
      </view>
      <view v-else class="post-preview-empty">
        <text class="post-preview-empty-text">还没有发过动态</text>
      </view>
    </view>
    <view class="section animate-fade-in-up stagger-3">
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
      <view class="edit-popup animate-fade-in-up">
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
  height: 320rpx;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  border-radius: 0 0 56rpx 56rpx;
  &::after {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(ellipse at 80% 20%, rgba(102,126,234,0.12) 0%, transparent 60%);
  }
}

.profile-deco {
  position: absolute;
  z-index: 1;
  border-radius: 50%;
  opacity: 0.08;
  pointer-events: none;
  &.profile-deco-1 {
    width: 300rpx; height: 300rpx;
    background: radial-gradient(circle, #667eea, transparent 70%);
    top: -60rpx; right: -60rpx;
    animation: pulse 6s ease-in-out infinite;
  }
  &.profile-deco-2 {
    width: 180rpx; height: 180rpx;
    background: radial-gradient(circle, #764ba2, transparent 70%);
    bottom: 0; left: -40rpx;
    animation: pulse 8s ease-in-out infinite reverse;
  }
}

.profile-info {
  position: relative;
  z-index: 2;
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
  text-shadow: 0 2rpx 20rpx rgba(0,0,0,0.15);
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

.follow-stats {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 40rpx;
  padding: 24rpx 0 0;
  position: relative;
  z-index: 2;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8rpx 24rpx;
}

.stat-num {
  font-size: 34rpx; font-weight: 700; color: #fff;
}

.stat-label {
  font-size: 24rpx; color: rgba(255,255,255,0.7); margin-top: 4rpx;
}

.stat-divider {
  width: 2rpx; height: 40rpx; background: rgba(255,255,255,0.2);
}

.section {
  margin: 20rpx 30rpx 0;
  border-radius: 22rpx;
  overflow: hidden;
  transition: transform 0.3s ease;
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

.section-title {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 24rpx 32rpx 16rpx;
	background: #fff;
	border-bottom: 1rpx solid var(--hairline, #e4e7ed);
}

.section-title-text {
	font-size: 30rpx;
	font-weight: 600;
	color: #303133;
}
.section-title-more {
	font-size: 26rpx;
	color: var(--brand-primary, #667eea);
}

.post-preview-list {
	background: #fff;
	padding: 0 32rpx;
}

.post-preview-item {
	padding: 24rpx 0;
	border-bottom: 1rpx solid var(--hairline, #e4e7ed);
	&:last-child { border-bottom: none; }
	&:active { opacity: 0.7; }
}

.post-preview-content {
	font-size: 28rpx;
	color: #303133;
	line-height: 1.6;
	display: block;
}

.post-preview-meta {
	display: flex;
	gap: 24rpx;
	margin-top: 12rpx;
}

.post-preview-stat {
	font-size: 24rpx;
	color: #c0c4cc;
}

.post-preview-empty {
	padding: 40rpx;
	background: #fff;
	display: flex;
	align-items: center;
	justify-content: center;
}

.post-preview-empty-text {
	font-size: 26rpx;
	color: #c0c4cc;
}
</style>
