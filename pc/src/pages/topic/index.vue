<template>
  <div class="topics-page">
    <PageHeader title="话题热榜" />
    <div v-if="loading" class="loading-area">
      <el-skeleton :rows="6" animated />
    </div>
    <div v-else-if="!topics.length" class="empty-card">
      <el-empty description="暂无话题" :image-size="80" />
    </div>
    <div v-else class="topic-grid">
      <div
        v-for="t in topics"
        :key="t.id"
        class="topic-card stagger-item"
        @click="$router.push('/square?topic=' + t.id)"
      >
        <div class="topic-icon">{{ t.icon }}</div>
        <div class="topic-info">
          <div class="topic-name">{{ t.name }}</div>
          <div class="topic-count">{{ t.post_count }} 条讨论</div>
        </div>
        <el-icon class="topic-arrow"><ArrowRight /></el-icon>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { listTopics } from "@/api/post"
import type { TopicData } from "@/api/post"
import PageHeader from "@/components/PageHeader.vue"
import { ArrowRight } from "@element-plus/icons-vue"

const topics = ref<TopicData[]>([])
const loading = ref(true)

onMounted(async () => {
  try { topics.value = await listTopics() }
  catch { /* */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.topics-page { max-width: 1100px; }

.topic-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: $space-4;
}

.topic-card {
  display: flex;
  align-items: center;
  gap: $space-4;
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-md;
  padding: $space-5;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
    .topic-arrow { transform: translateX(4px); opacity: 1; }
  }
}

.topic-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  background: rgba($brand-primary-hex, 0.1);
  border-radius: $radius-md;
  flex-shrink: 0;
}

.topic-info {
  flex: 1;
  min-width: 0;
}

.topic-name {
  font-size: $text-base;
  font-weight: 600;
  color: $text-primary;
}

.topic-count {
  font-size: $text-xs;
  color: $text-secondary;
  margin-top: 2px;
}

.topic-arrow {
  color: $text-secondary;
  font-size: 16px;
  transition: all 0.2s ease;
  opacity: 0.5;
}
</style>
