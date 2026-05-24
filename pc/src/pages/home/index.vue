<template>
  <div class="home-page">
    <section class="hero">
      <div class="hero-bg" />
      <div class="hero-content">
        <div class="hero-text">
          <p class="hero-eyebrow">洛阳高校 · 校园社区</p>
          <h1 class="hero-title">欢迎回来，{{ userStore.nickname }}</h1>
          <p class="hero-desc">发现校园新鲜事，连接洛阳高校的每一个角落</p>
          <div class="hero-actions">
            <el-button type="primary" size="large" :icon="Edit" @click="$router.push('/post/create')">
              发布动态
            </el-button>
            <el-button size="large" @click="$router.push('/square')">
              浏览广场
            </el-button>
          </div>
        </div>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-value">{{ trendingPostsData.length }}+</span>
            <span class="stat-label">今日动态</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ userStore.nickname ? '在线' : '--' }}</span>
            <span class="stat-label">你的状态</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ circles.length }}</span>
            <span class="stat-label">活跃圈子</span>
          </div>
        </div>
      </div>
    </section>

    <section v-if="banners.length" class="banner-section">
      <el-carousel height="200px" indicator-position="outside" :interval="5000" arrow="hover">
        <el-carousel-item v-for="b in banners" :key="b.id">
          <a v-if="b.link_url" :href="b.link_url" target="_blank" class="banner-link">
            <img :src="b.image_url" :alt="b.title" class="banner-img" />
            <div class="banner-overlay">
              <h3 class="banner-title">{{ b.title }}</h3>
            </div>
          </a>
          <div v-else class="banner-link">
            <img :src="b.image_url" :alt="b.title" class="banner-img" />
            <div class="banner-overlay">
              <h3 class="banner-title">{{ b.title }}</h3>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </section>


    <div class="home-grid">
      <div class="grid-main">
        <section class="section">
          <div class="section-header">
            <h2>精选动态</h2>
            <el-button text @click="$router.push('/square')">查看全部 →</el-button>
          </div>

          <div v-if="trendingLoading" class="loading-area">
            <SkeletonCard v-for="i in 3" :key="i" variant="post-card" />
          </div>
          <div v-else-if="!trendingPostsData.length" class="empty-card">
            <el-empty description="还没有动态" :image-size="80">
              <el-button type="primary" @click="$router.push('/post/create')" :icon="Edit">成为第一个分享的人</el-button>
            </el-empty>
          </div>
          <PostCard
            v-for="post in trendingPostsData"
            :key="post.id"
            :post="post"
            @like="handleLike"
            @comment="$router.push('/post/' + post.id)"
            @circle="$router.push('/circles/' + post.circle?.id)"
          />
        </section>

        <section class="section">
          <div class="section-header">
            <h2>热门圈子</h2>
            <el-button text @click="$router.push('/circles')">全部圈子 →</el-button>
          </div>

          <div v-if="circlesLoading" class="loading-area">
            <SkeletonCard v-for="i in 2" :key="i" variant="circle-card" />
          </div>
          <div v-else-if="!circles.length" class="empty-card" @click="$router.push('/circle/create')">
            <p>还没有圈子，创建一个吧！</p>
          </div>
          <CircleCard
            v-for="circle in circles"
            :key="circle.id"
            :circle="circle"
            @join="handleJoin"
            @leave="handleLeave"
          />
        </section>
      </div>

      <aside class="grid-side">
        <section class="side-card">
          <h3 class="side-title">快捷入口</h3>
          <div class="quick-links">
            <el-button text class="quick-btn" @click="$router.push('/activities')">
              <el-icon><Calendar /></el-icon> 校园活动
            </el-button>
            <el-button text class="quick-btn" @click="$router.push('/products')">
              <el-icon><ShoppingCart /></el-icon> 二手集市
            </el-button>
            <el-button text class="quick-btn" @click="$router.push('/found')">
              <el-icon><WarningFilled /></el-icon> 失物招领
            </el-button>
            <el-button text class="quick-btn" @click="$router.push('/whispers')">
              <el-icon><ChatDotSquare /></el-icon> 树洞
            </el-button>
            <el-button text class="quick-btn" @click="$router.push('/courses')">
              <el-icon><Reading /></el-icon> 课程表
            </el-button>
          </div>
        </section>

        <section class="side-card">
          <h3 class="side-title">洛阳高校</h3>
          <div class="campus-list">
            <div class="campus-item" v-for="uni in universities" :key="uni.name">
              <div class="campus-dot" :style="{ background: uni.color }" />
              <span>{{ uni.name }}</span>
            </div>
          </div>
        </section>

        <section class="side-card">
          <h3 class="side-title">个人信息</h3>
          <div class="profile-mini">
            <el-avatar :size="56" :src="userStore.avatar" />
            <div class="profile-name">{{ userStore.nickname }}</div>
            <div class="profile-school">{{ userStore.school || '未设置学校' }}</div>
            <el-button type="primary" size="small" style="width:100%;margin-top:12px" @click="$router.push('/user/' + userStore.userId)">
              查看主页
            </el-button>
          </div>
        </section>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useUserStore } from "@/stores/user"
import { trendingPosts, toggleLike } from "@/api/post"
import { listBanners } from "@/api/social"
import { listCircles, joinCircle, leaveCircle } from "@/api/circle"
import type { PostCardData } from "@/components/PostCard.vue"
import type { CircleCardData } from "@/components/CircleCard.vue"
import PostCard from "@/components/PostCard.vue"
import CircleCard from "@/components/CircleCard.vue"
import { ElMessage } from "element-plus"
import { Edit } from "@element-plus/icons-vue"
import SkeletonCard from "@/components/SkeletonCard.vue"

const banners = ref<{ id: number; title: string; image_url: string; link_url: string }[]>([])

const userStore = useUserStore()
const trendingLoading = ref(true)
const trendingPostsData = ref<PostCardData[]>([])
const circles = ref<CircleCardData[]>([])
const circlesLoading = ref(true)

const universities = [
  { name: "河南科技大学", color: "#C67A6A" },
  { name: "洛阳师范学院", color: "#5E8FD4" },
  { name: "洛阳理工学院", color: "#5EB87E" },
  { name: "河南推拿职业学院", color: "#D4A85E" },
]

async function handleLike(post: PostCardData) {
  try {
    const res = await toggleLike(post.id)
    post.is_liked = res.liked
    post.like_count += res.liked ? 1 : -1
  } catch { /* interceptor handles */ }
}

async function handleJoin(circle: CircleCardData) {
  try {
    await joinCircle(circle.id)
    circle.is_joined = true
    circle.member_count++
    ElMessage.success("已加入圈子")
  } catch { /* */ }
}

async function handleLeave(circle: CircleCardData) {
  try {
    await leaveCircle(circle.id)
    circle.is_joined = false
    circle.member_count--
    ElMessage.success("已退出圈子")
  } catch { /* */ }
}


onMounted(async () => {
    try {
      banners.value = (await listBanners()) as any
    } catch { /* */ }


  try {
    trendingPostsData.value = (await trendingPosts()) as unknown as PostCardData[]
  } catch { /* */ } finally { trendingLoading.value = false }

  try {
    circles.value = ((await listCircles({ limit: 10 })) || []).slice(0, 8)
  } catch { /* */ } finally { circlesLoading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

  margin: 0 auto;
 margin: 0 auto;
.home-page { max-width: 1100px;

// ═══ Hero ═══
.hero {
  position: relative;
  border-radius: $radius-lg;
  overflow: hidden;
  margin-bottom: $space-8;
  min-height: 200px;
}

// ═══ Banner Carousel ═══
.banner-section {
  margin-bottom: $space-6;
  border-radius: $radius-lg;
  overflow: hidden;

  .banner-link {
    display: block;
    position: relative;
    width: 100%;
    height: 100%;
    text-decoration: none;
  }

  .banner-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }

  .banner-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: $space-6;
    background: linear-gradient(to top, rgba(0,0,0,0.7), transparent);
  }

  .banner-title {
    font-size: $text-xl;
    font-weight: 600;
    color: #fff;
    margin: 0;
    text-shadow: 0 1px 4px rgba(0,0,0,0.4);
  }
}

:deep(.el-carousel__indicator--outside) {
  padding: 6px 4px 0;
}

:deep(.el-carousel__button) {
  background: oklch(0.45 0.01 30);
  opacity: 0.5;
  width: 20px;
  height: 4px;
  border-radius: 2px;
}

:deep(.el-carousel__indicator.is-active .el-carousel__button) {
  background: $brand-primary;
  opacity: 1;
  width: 30px;
}

.hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, oklch(0.18 0.04 30) 0%, oklch(0.22 0.06 20) 50%, oklch(0.16 0.03 30) 100%);
  &::before {
    content: "";
    position: absolute;
    inset: 0;
    background:
      radial-gradient(ellipse at 20% 50%, rgba($brand-primary-hex, 0.15) 0%, transparent 60%),
      radial-gradient(ellipse at 80% 20%, rgba($campus-blue-hex, 0.1) 0%, transparent 50%);
  }
  &::after {
    content: "";
    position: absolute;
    inset: 0;
    background-image: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.03'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  }
}

.hero-content {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $space-10 $space-8;
  z-index: 1;
}

.hero-text {
  
  max-width: 500px;
}

.hero-eyebrow {
  font-size: $text-xs;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: $brand-primary-light;
  margin-bottom: $space-2;
}

.hero-title {
  font-size: $text-3xl;
  font-weight: 700;
  line-height: 1.2;
  margin-bottom: $space-2;
}

.hero-desc {
  font-size: $text-sm;
  color: $text-secondary;
  margin-bottom: $space-6;
  line-height: 1.6;
}

.hero-actions {
  display: flex;
  gap: $space-3;
}

.hero-stats {
  display: flex;
  gap: $space-8;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $space-1;

  .stat-value {
    font-size: $text-2xl;
    font-weight: 700;
    color: $brand-primary-light;
  }

  .stat-label {
    font-size: $text-xs;
    color: $text-muted;
  }
}

// ═══ Grid Layout ═══
.home-grid {
  display: grid;
  grid-template-columns: 1fr 280px;
  gap: $space-6;
  align-items: start;
}

// ═══ Sections ═══
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $space-4;

  h2 { font-size: $text-lg; font-weight: 600; }
}

.loading-area { padding: $space-6 0; }

.empty-card {
  background: $bg-card;
  border: 1px dashed $border-light;
  border-radius: $radius-md;
  padding: $space-10;
  text-align: center;
  cursor: pointer;
  transition: all $duration-fast;
  color: $text-muted;

  &:hover {
    border-color: rgba($brand-primary-hex, 0.4);
    color: $text-secondary;
  }

  .empty-icon { font-size: 32px; margin-bottom: $space-2; }
  p { font-size: $text-sm; }
}

// ═══ Side Panel ═══
.side-card {
  background: $bg-card;
  border: 1px solid $border-default;
  border-radius: $radius-md;
  padding: $space-5;
  margin-bottom: $space-4;
}

.side-title {
  font-size: $text-sm;
  font-weight: 600;
  margin-bottom: $space-3;
  color: $text-secondary;
}

.quick-links {
  display: flex;
  flex-direction: column;
  gap: $space-1;
}

.quick-btn {
  justify-content: flex-start;
  color: $text-primary;
  font-size: $text-sm;
  padding: $space-2 $space-3;
  border-radius: $radius-sm;

  &:hover {
    background: rgba($brand-primary-hex, 0.08);
  }
  .el-icon { margin-right: $space-2; }
}

.campus-list {
  display: flex;
  flex-direction: column;
  gap: $space-2;
}

.campus-item {
  display: flex;
  align-items: center;
  gap: $space-2;
  font-size: $text-sm;
}

.campus-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.profile-mini {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.profile-name {
  font-size: $text-sm;
  font-weight: 600;
  margin-top: $space-2;
}

.profile-school {
  font-size: $text-xs;
  color: $text-muted;
}
</style>
