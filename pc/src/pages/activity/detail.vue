<template>
  <div v-if="loading" class="loading-wrap"><el-skeleton :rows="6" animated /></div>
  <div v-else-if="!activity" class="empty-state"><el-empty description="活动不存在" /></div>
  <div v-else class="detail-page stagger-item">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn stagger-item">返回</el-button>
    <div class="activity-detail">
      <div class="header">
        <el-tag size="small">{{ activity.activity_type }}</el-tag>
        <h1>{{ activity.title }}</h1>
      </div>
      <div class="meta-grid">
        <div><el-icon><Timer /></el-icon> {{ formatDate(activity.start_time) }}</div>
        <div v-if="activity.location"><el-icon><Location /></el-icon> {{ activity.location }}</div>
        <div><el-icon><User /></el-icon> {{ activity.participant_count }}/{{ activity.max_participants }} 人</div>
      </div>
      <p class="description">{{ activity.description }}</p>
      <div v-if="activity.images?.length" class="images">
        <img v-for="(img, i) in activity.images" :key="i" :src="img" loading="lazy" @click="lightboxIndex = i" />
      </div>
      <div class="actions">
        <el-button v-if="!activity.is_participant && activity.status === 0" type="primary" @click="handleJoin">
          报名参加
        </el-button>
        <el-button v-else-if="activity.is_participant" @click="handleLeave">取消报名</el-button>
      </div>
    </div>
  </div>
    <ImagePreview :images="activity?.images || []" v-model="lightboxIndex" />

</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { getActivity, joinActivity, leaveActivity } from "@/api/activity"
import type { ActivityData } from "@/api/activity"
import dayjs from "dayjs"
import { ArrowLeft } from "@element-plus/icons-vue"
import ImagePreview from "@/components/ImagePreview.vue"

const lightboxIndex = ref<number | null>(null)

const route = useRoute()
const router = useRouter()
const activity = ref<ActivityData | null>(null)
const loading = ref(true)

function formatDate(t: string) { return dayjs(t).format("MM-DD HH:mm") }

async function handleJoin() {
  if (!activity.value) return
  try {
    await joinActivity(activity.value.id)
    activity.value.is_participant = true
    activity.value.participant_count++
  } catch { /* handled by interceptor */ }
}

async function handleLeave() {
  if (!activity.value) return
  try {
    await leaveActivity(activity.value.id)
    activity.value.is_participant = false
    activity.value.participant_count--
  } catch { /* handled by interceptor */ }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try { activity.value = await getActivity(id) }
  catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

 margin: 0 auto;
.detail-page { max-width: 680px; }
.back-btn { margin-bottom: 16px; }

.activity-detail {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 24px;
  .header { margin-bottom: 16px; h1 { font-size: 22px; font-weight: 700; margin-top: 8px; } }
  .meta-grid {
    display: flex; flex-wrap: wrap; gap: 16px; font-size: 13px; color: $text-secondary;
    margin-bottom: 16px; transition: all 0.3s ease;
.el-icon { vertical-align: middle; margin-right: 4px; }
  }
  .description { font-size: 14px; line-height: 1.8; margin-bottom: 16px; }
  .images { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 16px;
    img { width: 100px; height: 100px; object-fit: cover; border-radius: $radius-sm;
      transition: transform 0.3s ease;
      &:hover { transform: scale(1.08); }
    }
  }
  .actions { padding-top: 16px; border-top: 1px solid $border-color; }
}
</style>
