<template>
  <div class="post-card" @click="handleClick">
    <div class="post-header">
      <el-avatar :size="36" :src="post.author?.avatar" />
      <div class="post-author-info">
        <span class="name">{{ post.author?.nickname || '匿名' }}</span>
        <span class="meta">
          <span class="time">{{ formatTime(post.created_at) }}</span>
          <span v-if="post.school" class="school">· {{ post.school }}</span>
        </span>
      </div>
    </div>

    <p v-if="post.title" class="post-title">{{ post.title }}</p>
    <p class="post-content">{{ post.content }}</p>

    <div v-if="post.images?.length" class="post-images">
      <AppImage
        v-for="(img, i) in post.images" :key="i"
        :src="img" img-class="post-image"
        @click.stop="lightboxIndex = i"
      />
    </div>

    <div class="post-footer">
      <span
        :class="['action', { liked: post.is_liked }]"
        @click.stop="emit('like', post)"
      >
        <el-icon><Goods /></el-icon>
        <span>{{ post.like_count || 0 }}</span>
      </span>
      <span class="action" @click.stop="emit('comment', post)">
        <el-icon><ChatDotSquare /></el-icon>
        <span>{{ post.comment_count || 0 }}</span>
      </span>
      <span
        :class="['action', { favorited: post.is_favorited }]"
        @click.stop="emit('favorite', post)"
      >
        <el-icon><Star /></el-icon>
        <span>{{ post.is_favorited ? '已收藏' : '收藏' }}</span>
      </span>
      <span v-if="post.circle" class="circle-tag" @click.stop="emit('circle', post)">
        <el-icon><Connection /></el-icon>
        {{ post.circle.name }}
      </span>
    </div>
  </div>

    <ImagePreview :images="post.images || []" v-model="lightboxIndex" />
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useTimeFormat } from "@/composables/useTimeFormat"
import ImagePreview from "@/components/ImagePreview.vue"
import AppImage from '@/components/AppImage.vue'

export interface PostCardData {
  id: number
  content: string
  title?: string
  images?: string[]
  like_count: number
  comment_count: number
  is_liked?: boolean
  is_favorited?: boolean
  created_at: string
  school?: string
  author?: { nickname?: string; avatar?: string }
  circle?: { id: number; name: string }
}

const props = defineProps<{
  post: PostCardData
}>()

const emit = defineEmits<{
  like: [post: any]
  comment: [post: any]
  favorite: [post: any]
  circle: [post: any]
}>()

const router = useRouter()
const { formatTime } = useTimeFormat()
const lightboxIndex = ref<number | null>(null)

function handleClick() {
  router.push(`/post/${props.post.id}`)
}
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.post-card {
  background: $bg-card;
  border: 1px solid $border-default;
  border-radius: $radius-md;
  padding: $space-5;
  margin-bottom: $space-4;
  cursor: pointer;
  transition: all $duration-fast $ease-out;

  &:hover {
    border-color: rgba($brand-primary-hex, 0.5);
    box-shadow: $shadow-lg;
    transform: translateY(-2px);
  }
  &:active {
    transform: translateY(0);
  }
}

.post-header {
  display: flex;
  align-items: center;
  gap: $space-3;
  margin-bottom: $space-3;
}

.post-author-info {
  flex: 1;
  min-width: 0;
  .name {
    font-size: $text-sm;
    font-weight: 600;
  }
  .meta {
    display: block;
    font-size: $text-xs;
    color: $text-muted;
    margin-top: 1px;
  }
}

.post-title {
  font-size: $text-lg;
  font-weight: 600;
  margin-bottom: $space-2;
  line-height: 1.4;
}

.post-content {
  font-size: $text-sm;
  line-height: 1.7;
  color: $text-primary;
  margin-bottom: $space-3;
  white-space: pre-wrap;
  word-break: break-word;
}

.post-images {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: $space-1;
  margin-bottom: $space-3;

  img {
    width: 100%;
    aspect-ratio: 1;
    object-fit: cover;
    border-radius: $radius-sm;
    cursor: zoom-in;
    transition: opacity $duration-fast;

    &.single {
      max-width: 240px;
      aspect-ratio: auto;
      max-height: 240px;
    }
    &:hover { opacity: 0.9; }
  }
}

.post-footer {
  display: flex;
  align-items: center;
  gap: $space-5;
  padding-top: $space-3;
  border-top: 1px solid $border-default;
}

.action {
  display: flex;
  align-items: center;
  gap: $space-1;
  font-size: $text-sm;
  color: $text-secondary;
  cursor: pointer;
  transition: color $duration-fast;

  &:hover { color: $brand-primary; }
  &.liked { color: $brand-primary; }
  &.favorited { color: $campus-gold; }
}

.circle-tag {
  display: flex;
  align-items: center;
  gap: $space-1;
  font-size: $text-xs;
  color: $campus-blue;
  cursor: pointer;
  margin-left: auto;

  &:hover { color: $campus-cyan; }
}
</style>
