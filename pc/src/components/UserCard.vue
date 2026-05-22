<template>
  <div class="user-card" @click="$router.push('/user/' + user.id)">
    <el-avatar :size="40" :src="user.avatar" />
    <div class="user-info">
      <div class="user-name">{{ user.nickname }}</div>
      <div class="user-school">{{ user.school || '未设置学校' }}</div>
    </div>
    <div v-if="!noAction" class="user-action">
      <el-button
        v-if="isFollowing !== undefined"
        :type="isFollowing ? 'default' : 'primary'"
        size="small"
        :plain="isFollowing"
        :loading="followLoading"
        @click.stop="emit('toggleFollow', user)"
      >
        {{ isFollowing ? '已关注' : '关注' }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface UserCardData {
  id: number
  nickname: string
  avatar?: string
  school?: string
}

defineProps<{
  user: UserCardData
  isFollowing?: boolean
  followLoading?: boolean
  noAction?: boolean
}>()

const emit = defineEmits<{
  toggleFollow: [user: UserCardData]
}>()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.user-card {
  display: flex;
  align-items: center;
  gap: $space-3;
  padding: $space-3 $space-4;
  border-radius: $radius-sm;
  cursor: pointer;
  transition: background $duration-fast;

  &:hover { background: $bg-hover; }

  .user-info {
    flex: 1;
    min-width: 0;
    .user-name { font-size: $text-sm; font-weight: 600; }
    .user-school { font-size: $text-xs; color: $text-muted; }
  }

  .user-action { flex-shrink: 0; }
}
</style>
