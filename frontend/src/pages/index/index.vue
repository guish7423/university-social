<script setup lang="ts">
import { ref, onMounted } from "vue"
import { trendingPosts, listTopics, type PostData, type Topic } from "@/api/post"
import { useUserStore } from "@/store/user"

const userStore = useUserStore()
const trending = ref<PostData[]>([])
const topics = ref<Topic[]>([])
  const loading = ref(true)

  const features = [
  { icon: "◈", title: "动态广场", desc: "发现校园新鲜事", path: "/pages/square/index", tab: true },
  { icon: "◎", title: "校园圈子", desc: "找到志同道合的伙伴", path: "/pages/circle/list", tab: true },
  { icon: "♢", title: "树洞", desc: "匿名倾诉你的心声", path: "/pages/whisper/index", tab: true },
  { icon: "▤", title: "二手市场", desc: "闲置物品交易", path: "/pages/product/list" },
  { icon: "▣", title: "校园活动", desc: "参与精彩活动", path: "/pages/activity/list" },
  { icon: "★", title: "课程评价", desc: "选课参考指南", path: "/pages/course/search" },
  { icon: "◉", title: "失物招领", desc: "找回丢失物品", path: "/pages/found/list" },
  { icon: "♢", title: "积分排行", desc: "赚积分赢奖励", path: "/pages/points/index" },
  ]

  onMounted(async () => {
  try {
    const [p, t] = await Promise.all([trendingPosts(), listTopics()])
    trending.value = p || []
    topics.value = t || []
  } catch {}
  loading.value = false
    })

    function go(path: string, tab = false) {
  if (tab) uni.switchTab({ url: path })
  else uni.navigateTo({ url: path })
        }

        function formatTime(t: string) {
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 3600000) return Math.floor(diff / 60000) + "分钟前"
  if (diff < 86400000) return Math.floor(diff / 3600000) + "小时前"
  const days = Math.floor(diff / 86400000)
  if (days < 7) return days + "天前"
  return (d.getMonth() + 1) + "/" + d.getDate()
}
    </script>

        <template>
  <view class="page">
    <view class="hero">
      <view class="hero-content">
        <text class="hero-title">校园社</text>
        <text class="hero-subtitle">连接校园里的每一个人</text>
      </view>
    </view>

    <view class="section">
      <view class="section-header">
        <text class="section-title">全部功能</text>
      </view>
      <view class="feature-grid">
        <view v-for="(f, i) in features"
          :key="i"
          :class="['feature-card', 'animate-fade-in-up', 'stagger-' + Math.min(i + 1, 8)]"
          @click="go(f.path, f.tab)">
          <view class="feature-icon-wrap"><text class="feature-icon">{{ f.icon }}</text></view>
          <text class="feature-title">{{ f.title }}</text>
          <text class="feature-desc">{{ f.desc }}</text>
        </view>
      </view>
    </view>

    <view v-if="!loading && trending.length" class="section">
      <view class="section-header">
        <text class="section-title">热门动态</text>
        <text class="section-more" @click="go('/pages/square/index', true)">查看更多</text>
      </view>
      <view class="trending-list">
        <view v-for="(post, i) in trending.slice(0, 3)" :key="post.id"
          :class="['trending-card', 'animate-fade-in-up', 'stagger-' + Math.min(i + 1, 4)]"
          @click="go('/pages/post/detail?id=' + post.id)">
          <view class="trending-header">
            <text class="trending-author">{{ post.author?.nickname || '匿名' }}</text>
            <text class="trending-time">{{ formatTime(post.created_at) }}</text>
          </view>
          <text class="trending-content">{{ post.content }}</text>
          <view class="trending-meta">
            <text class="trending-stat">{{ post.like_count || 0 }} 赞</text>
            <text class="trending-stat">{{ post.comment_count || 0 }} 评论</text>
          </view>
        </view>
      </view>
    </view>

    <view v-if="!userStore.isLogin" class="login-prompt">
      <text class="login-text">登录后解锁全部功能</text>
      <view class="login-btn" @click="go('/pages/login/index')">立即登录</view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
  .page { min-height: 100vh; background: var(--color-canvas, #F7F4F0); }

  .hero {
  background: linear-gradient(135deg, #1E2A3A 0%, #2A3A4E 100%);
  padding: 80rpx 40rpx 60rpx;
  text-align: center;
  }
  .hero-content { display: flex; flex-direction: column; align-items: center; gap: 12rpx; }
  .hero-title { font-size: 40rpx; font-weight: 600; color: #C67A6A; letter-spacing: 4rpx; }
  .hero-subtitle { font-size: 26rpx; color: rgba(255,255,255,0.55); font-weight: 300; letter-spacing: 2rpx; }

.section { padding: 40rpx 32rpx; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 28rpx; }
  .section-title { font-size: 30rpx; font-weight: 600; color: var(--ink, #1E2A3A); }
  .section-more { font-size: 24rpx; color: var(--brand-primary, #C67A6A); }

  .feature-grid { display: grid; grid-template-columns: 1fr 1fr 1fr 1fr; gap: 20rpx; }
  .feature-card {
  background: var(--color-surface, #fff);
  border-radius: 16rpx;
  padding: 28rpx 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
  transition: all 0.2s ease;
  cursor: pointer;
  &:active { transform: scale(0.97); }
  }
  .feature-icon-wrap {
  width: 64rpx; height: 64rpx;
  border-radius: 16rpx;
  background: rgba(198,122,106,0.08);
  display: flex; align-items: center; justify-content: center;
}
.feature-icon { font-size: 28rpx; color: #C67A6A; }
  .feature-title { font-size: 26rpx; font-weight: 500; color: var(--ink, #1E2A3A); }
  .feature-desc { font-size: 22rpx; color: var(--ink-subtle, #8E9BAB); }

  .trending-list { display: flex; flex-direction: column; gap: 16rpx; }
  .trending-card {
  background: var(--color-surface, #fff);
  border-radius: 14rpx;
  padding: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(30,42,58,0.04);
  cursor: pointer;
  transition: all 0.2s ease;
  &:active { transform: scale(0.99); }
}
.trending-header { display: flex; justify-content: space-between; margin-bottom: 12rpx; }
  .trending-author { font-size: 24rpx; font-weight: 500; color: var(--ink, #1E2A3A); }
  .trending-time { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); }
  .trending-content { font-size: 26rpx; color: var(--ink-muted, #5C6B7E); line-height: 1.6; display: block; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.trending-meta { display: flex; gap: 24rpx; margin-top: 16rpx; }
.trending-stat { font-size: 22rpx; color: var(--ink-tertiary, #B8C2CE); }

  .login-prompt {
  margin: 0 32rpx 40rpx;
  padding: 32rpx;
  background: rgba(198,122,106,0.06);
  border-radius: 16rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20rpx;
  border: 1rpx solid rgba(198,122,106,0.12);
}
.login-text { font-size: 26rpx; color: var(--ink-muted, #5C6B7E); }
.login-btn {
  padding: 16rpx 48rpx;
  background: #C67A6A;
  color: #fff;
  border-radius: 24rpx;
  font-size: 26rpx;
  font-weight: 500;
  transition: all 0.2s ease;
  &:active { opacity: 0.85; transform: scale(0.97); }
}

@media (min-width: 1024px) {
  .feature-grid { grid-template-columns: repeat(4, 1fr); }
  .hero { padding: 100rpx 40rpx 80rpx; }
  .section { max-width: 1200rpx; margin: 0 auto; }
}
</style>
