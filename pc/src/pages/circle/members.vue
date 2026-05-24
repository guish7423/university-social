<template>
  <div class="members-page">
    <PageHeader title="圈子成员" @back="$router.back()" />
    <div v-if="loading" class="loading-area">
      <SkeletonCard v-for="i in 6" :key="i" variant="post-card" />
    </div>
    <div v-else-if="!members.length" class="empty-card">
      <el-empty description="暂无成员" :image-size="80" />
    </div>
    <div v-else class="member-grid">
      <div
        v-for="m in members"
        :key="m.id"
        class="member-card stagger-item"
        @click="$router.push('/user/' + m.id)"
      >
        <el-avatar :size="48" :src="m.avatar">{{ m.nickname[0] }}</el-avatar>
        <div class="member-name">{{ m.nickname }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"
import { listCircleMembers } from "@/api/circle"
import PageHeader from "@/components/PageHeader.vue"
import SkeletonCard from "@/components/SkeletonCard.vue"

const route = useRoute()
const members = ref<{ id: number; nickname: string; avatar: string }[]>([])
const loading = ref(true)

onMounted(async () => {
  const id = Number(route.params.id)
  try { members.value = await listCircleMembers(id) }
  catch { /* */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.members-page { max-width: 1100px; }

.member-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: $space-4;
}

.member-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $space-3;
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-md;
  padding: $space-5 $space-3;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  }
}

.member-name {
  font-size: $text-sm;
  font-weight: 600;
  color: $text-primary;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}
</style>
