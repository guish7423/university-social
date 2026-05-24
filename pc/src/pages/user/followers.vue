<template>
  <div class="list-page">
    <el-button text :icon="ArrowLeft" @click="$router.back()" class="back-btn">返回</el-button>
    <PageHeader title="粉丝" />

    <LoadingWrapper :loading="loading && !list.length" :data="list.length" skeleton-variant="post-card">
      <template #empty>
        <el-empty description="暂无粉丝" />
      </template>

      <template #default>
        <div class="user-grid">
          <div
            v-for="u in list"
            :key="u.id"
            class="user-card stagger-item"
            @click="$router.push('/user/' + u.id)"
          >
            <el-avatar :size="48" :src="u.avatar" />
            <div class="user-info">
              <div class="name">{{ u.nickname }}</div>
              <div v-if="u.school" class="school">{{ u.school }}</div>
            </div>
          </div>
        </div>
      </template>
    </LoadingWrapper>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"
import { getFollowers } from "@/api/social"
import type { UserInfo } from "@/api/auth"
import PageHeader from "@/components/PageHeader.vue"
import LoadingWrapper from "@/components/LoadingWrapper.vue"
import { ArrowLeft } from "@element-plus/icons-vue"

const route = useRoute()
const list = ref<UserInfo[]>([])
const loading = ref(true)
const userId = Number(route.params.id)

onMounted(async () => {
  try {
    const res = await getFollowers(userId)
    list.value = res
  } catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

 margin: 0 auto;
.list-page { max-width: 640px; }
.back-btn { margin-bottom: 4px; }

.user-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-md;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s, border-color 0.2s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
    border-color: $primary-light;
  }

  .user-info {
    flex: 1;
    min-width: 0;

    .name {
      font-size: 15px;
      font-weight: 600;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .school {
      font-size: 12px;
      color: $text-muted;
      margin-top: 2px;
    }
  }
}
</style>
