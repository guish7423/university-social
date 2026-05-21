<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useUserStore } from "@/store/user"
import { devLogin as devLoginApi } from "@/api/user"

const userStore = useUserStore()
const devNickname = ref("")
const showDevLogin = ref(false)

onMounted(() => {
  showDevLogin.value = !userStore.isLogin
})

async function devLogin() {
  const nick = devNickname.value.trim() || "测试用户"
  try {
    const res = await devLoginApi(nick)
    userStore.setToken(res.token)
    userStore.setUserInfo(res.user)
    showDevLogin.value = false
    uni.switchTab({ url: "/pages/square/index" })
  } catch {
    uni.showToast({ title: "登录失败", icon: "none" })
  }
}

function goSquare() { uni.switchTab({ url: '/pages/square/index' }) }
function goCircles() { uni.switchTab({ url: '/pages/circle/list' }) }
function goFriends() { uni.navigateTo({ url: '/pages/friend/index' }) }
function goNotifications() { uni.navigateTo({ url: '/pages/notification/index' }) }
function goMarket() { uni.navigateTo({ url: '/pages/product/list' }) }
function goActivities() { uni.navigateTo({ url: '/pages/activity/list' }) }
function goPoints() { uni.navigateTo({ url: '/pages/points/index' }) }
function goLostItems() { uni.navigateTo({ url: '/pages/found/list' }) }

</script>

  <template>
  <view class="container">
    <view class="hero">
      <view class="hero-bg" />
      <view class="hero-deco hero-deco-1" />
      <view class="hero-deco hero-deco-2" />
      <view class="hero-deco hero-deco-3" />
      <view class="hero-content animate-fade-in-up">
        <text class="logo">校园社</text>
        <text class="tagline">你的大学生活，从这里开始</text>
        <view class="hero-indicator">
          <view class="scroll-dot" />
        </view>
      </view>
    </view>

    <view v-if="showDevLogin" class="login-card animate-fade-in-up stagger-1">
      <view class="login-header">
        <text class="login-title">🎓 体验校园社</text>
        <text class="login-desc">输入昵称即可开始，无需注册</text>
      </view>
      <view class="login-form">
        <u-input
          v-model="devNickname"
          placeholder="输入你的昵称"
          :border="true"
          shape="square"
          :customStyle="{ padding: '20rpx 24rpx', fontSize: '28rpx', height: '88rpx' }"
        />
        <u-button
          type="primary"
          shape="circle"
          size="large"
          :customStyle="{ marginTop: '28rpx', height: '88rpx', fontSize: '32rpx', fontWeight: '600' }"
          @click="devLogin"
        >
          进入校园社
        </u-button>
      </view>
    </view>

    <view v-else class="welcome-card animate-fade-in-up stagger-1">
      <view class="welcome-user">
        <u-avatar
          :src="userStore.avatar"
          size="80"
          shape="circle"
        />
        <view class="welcome-info">
          <text class="welcome-name">{{ userStore.nickname || '用户' }}</text>
          <text class="welcome-text">欢迎回来！</text>
        </view>
      </view>
    </view>

    <view class="feature-grid">
      <view class="feature-item animate-fade-in-up stagger-2" @click="goSquare">
        <view class="feature-icon-wrap" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
          <u-icon name="home" size="40" color="#fff" />
        </view>
        <text class="feature-name">动态广场</text>
        <text class="feature-desc">发现校园新鲜事</text>
      </view>
      <view class="feature-item animate-fade-in-up stagger-3" @click="goCircles">
        <view class="feature-icon-wrap" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
          <u-icon name="grid" size="40" color="#fff" />
        </view>
        <text class="feature-name">校园圈子</text>
        <text class="feature-desc">找到你的同好</text>
      </view>
<view class="feature-item animate-fade-in-up stagger-5" @click="goMarket">
<view class="feature-icon-wrap" style="background: linear-gradient(135deg, #f9d423 0%, #ff4e50 100%);">
<u-icon name="shopping" size="40" color="#fff" />
</view>
<text class="feature-name">二手市场</text>
<text class="feature-desc">买卖校园闲置</text>
</view>
<view class="feature-item animate-fade-in-up stagger-5" @click="goLostItems">
<view class="feature-icon-wrap" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
<u-icon name="info-circle" size="40" color="#fff" />
</view>
<text class="feature-name">失物招领</text>
</view>

<view class="feature-item animate-fade-in-up stagger-6" @click="goActivities">
<view class="feature-icon-wrap" style="background: linear-gradient(135deg, #667eea 0%, #43e97b 100%);">
<u-icon name="account" size="40" color="#fff" />
</view>
<text class="feature-name">校园活动</text>
      <view class="feature-item animate-fade-in-up stagger-7" @click="goPoints">
        <view class="feature-icon" style="background: linear-gradient(135deg, #667eea, #764ba2)">
          <view class="feature-icon-inner">🏆</view>
        </view>
        <text class="feature-name">积分排行</text>
      </view>
<text class="feature-desc">组队与活动</text>
</view>
			<view class="feature-item animate-fade-in-up stagger-7" @click="goNotifications">
				<view class="feature-icon-wrap" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
					<u-icon name="bell" size="40" color="#fff" />
				</view>
				<text class="feature-name">通知</text>
				<text class="feature-desc">查看消息提醒</text>
			</view>
		</view>
	</view>
</template>

.hero {
  position: relative;
  padding: 120rpx 40rpx 100rpx;
  overflow: hidden;
}

.hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  border-radius: 0 0 56rpx 56rpx;
  &::after {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(ellipse at 20% 50%, rgba(102,126,234,0.15) 0%, transparent 60%),
                radial-gradient(ellipse at 80% 20%, rgba(118,75,162,0.1) 0%, transparent 50%);
  }
}

.hero-deco {
  position: absolute;
  z-index: 1;
  border-radius: 50%;
  opacity: 0.08;
  pointer-events: none;
  &.hero-deco-1 {
    width: 400rpx; height: 400rpx;
    background: radial-gradient(circle, #667eea, transparent 70%);
    top: -100rpx; right: -80rpx;
    animation: pulse 6s ease-in-out infinite;
  }
  &.hero-deco-2 {
    width: 200rpx; height: 200rpx;
    background: radial-gradient(circle, #764ba2, transparent 70%);
    bottom: -40rpx; left: -60rpx;
    animation: pulse 8s ease-in-out infinite reverse;
  }
  &.hero-deco-3 {
    width: 120rpx; height: 120rpx;
    background: radial-gradient(circle, #4facfe, transparent 70%);
    top: 40%; right: 10%;
    animation: pulse 5s ease-in-out infinite;
    animation-delay: 1s;
  }
}

.hero-content {
  position: relative;
  z-index: 2;
  text-align: center;
}

.logo {
  font-size: 72rpx;
  font-weight: 800;
  color: #fff;
  display: block;
  letter-spacing: 6rpx;
  text-shadow: 0 2rpx 20rpx rgba(102,126,234,0.3);
}

.tagline {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.65);
  margin-top: 20rpx;
  display: block;
  letter-spacing: 2rpx;
}

.hero-indicator {
  margin-top: 40rpx;
  display: flex;
  justify-content: center;
}

.scroll-dot {
  width: 12rpx; height: 12rpx;
  border-radius: 50%;
  background: rgba(255,255,255,0.4);
  animation: pulse 2s ease-in-out infinite;
}

.login-card {
  margin: -60rpx 30rpx 0;
  background: #fff;
  border-radius: 28rpx;
  padding: 44rpx;
  position: relative;
  z-index: 3;
  box-shadow: 0 8rpx 40rpx rgba(0, 0, 0, 0.06);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  &:active { transform: scale(0.99); }
}

.login-header {
  margin-bottom: 36rpx;
  text-align: center;
}

.login-title {
  font-size: 38rpx;
  font-weight: 700;
  color: #1a1a2e;
  display: block;
  margin-bottom: 12rpx;
}

.login-desc {
  font-size: 26rpx;
  color: #909399;
}

.welcome-card {
  margin: -60rpx 30rpx 0;
  background: #fff;
  border-radius: 28rpx;
  padding: 40rpx;
  position: relative;
  z-index: 3;
  box-shadow: 0 8rpx 40rpx rgba(0, 0, 0, 0.06);
  transition: transform 0.3s ease;
  &:active { transform: scale(0.99); }
}

.welcome-user {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.welcome-info {
  flex: 1;
}

.welcome-name {
  font-size: 34rpx;
  font-weight: 700;
  display: block;
}

.welcome-text {
  font-size: 26rpx;
  color: #909399;
  margin-top: 6rpx;
  display: block;
}

.feature-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24rpx;
  padding: 36rpx 30rpx;
}

.feature-item {
  background: #fff;
  border-radius: 22rpx;
  padding: 40rpx 24rpx;
  text-align: center;
  box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.04);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1rpx solid transparent;
  &:active {
    transform: scale(0.95);
    box-shadow: 0 1rpx 8rpx rgba(0, 0, 0, 0.04);
  }
}

.feature-icon-wrap {
  width: 88rpx;
  height: 88rpx;
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20rpx;
  transition: transform 0.3s ease;
  .feature-item:active & {
    transform: scale(0.9);
  }
}

.feature-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #303133;
  display: block;
}

.feature-desc {
  font-size: 22rpx;
  color: #909399;
  margin-top: 8rpx;
  display: block;
}
