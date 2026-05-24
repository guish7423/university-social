<template>
  <div class="friends-page">
    <div class="page-header">
      <h1>好友</h1>
      <el-button @click="$router.push('/search')">添加好友</el-button>
    </div>
    <el-tabs v-model="tab">
      <el-tab-pane label="好友列表" name="friends">
        <div v-if="friendsLoading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
        <div v-else-if="!friends.length" class="empty-state"><el-empty description="暂无好友" /></div>
        <div v-else class="user-list">
          <div v-for="f in friends" class="user-row stagger-item" :key="f.id" @click="$router.push('/user/' + f.id)">
            <el-avatar :size="40" :src="f.avatar" />
            <div class="user-info">
              <div class="name">{{ f.nickname }}</div>
              <div class="school" v-if="f.school">{{ f.school }}</div>
            </div>
            <el-button text @click.stop="handleRemove(f.id)">移除</el-button>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="好友请求" name="requests">
        <div class="requests-header">
          <span>好友请求</span>
          <el-tag v-if="requests.length" type="danger" size="small">{{ requests.length }}</el-tag>
        </div>
        <div v-if="requestsLoading" class="loading-wrap"><el-skeleton :rows="3" animated /></div>
        <div v-else-if="!requests.length" class="empty-state"><el-empty description="暂无好友请求" /></div>
        <div v-else class="request-grid">
          <div v-for="r in requests" :key="r.id" class="request-card stagger-item">
            <el-avatar :size="48" :src="r.avatar" />
            <div class="req-name">{{ r.nickname }}</div>
            <el-tag v-if="actionState[r.id]" :type="actionState[r.id].status === 'accepted' ? 'success' : 'info'" size="small" effect="dark">
              {{ actionState[r.id].status === 'accepted' ? '已通过' : '已拒绝' }}
            </el-tag>
            <div v-else class="req-actions">
              <el-button type="primary" size="small" :loading="actionState[r.id]?.loading" @click="handleAccept(r.id)">接受</el-button>
              <el-button size="small" :loading="actionState[r.id]?.loading" @click="handleReject(r.id)">拒绝</el-button>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue"
import { listFriends, listFriendRequests, acceptFriendRequest, rejectFriendRequest, removeFriend } from "@/api/social"
import type { UserInfo } from "@/api/auth"

const tab = ref("friends")
const friends = ref<UserInfo[]>([])
const requests = ref<UserInfo[]>([])
const friendsLoading = ref(true)
const requestsLoading = ref(true)

interface ActionState { loading: boolean; status: "accepted" | "rejected" }
const actionState = reactive<Record<number, ActionState>>({})

async function handleAccept(id: number) {
  actionState[id] = { loading: true, status: "accepted" }
  try {
    await acceptFriendRequest(id)
    actionState[id] = { loading: false, status: "accepted" }
  } catch {
    delete actionState[id]
  }
}

async function handleReject(id: number) {
  actionState[id] = { loading: true, status: "rejected" }
  try {
    await rejectFriendRequest(id)
    actionState[id] = { loading: false, status: "rejected" }
  } catch {
    delete actionState[id]
  }
}

async function handleRemove(id: number) {
  try { await removeFriend(id); friends.value = friends.value.filter(f => f.id !== id) }
  catch { /* handled by interceptor */ }
}

onMounted(async () => {
  try { friends.value = await listFriends() } catch { /* handled by interceptor */ } finally { friendsLoading.value = false }
  try { requests.value = await listFriendRequests() } catch { /* handled by interceptor */ } finally { requestsLoading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

 margin: 0 auto;
.friends-page { max-width: 640px; }

.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;
  h1 { font-size: 22px; font-weight: 700; }
}

.requests-header {
  display: flex; align-items: center; gap: 8px; margin-bottom: 14px;
  font-size: 15px; font-weight: 600;
}

.user-list { display: flex; flex-direction: column; gap: 8px; }

.user-row {
  display: flex; align-items: center; gap: 12px;
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 12px; cursor: pointer; transition: all $duration-fast;
  &:hover { border-color: $primary-light; }
  .user-info { flex: 1; min-width: 0;
    .name { font-size: 14px; font-weight: 600; }
    .school { font-size: 12px; color: $text-muted; margin-top: 2px; }
  }
}

.request-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
}

.request-card {
  display: flex; flex-direction: column; align-items: center; gap: 10px;
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-lg; padding: 20px 16px;
  transition: all $duration-fast;
  &:hover { border-color: $primary-light; transform: translateY(-2px); box-shadow: $shadow-sm; }

  .req-name { font-size: 15px; font-weight: 600; text-align: center; }
  .req-actions { display: flex; gap: 8px; margin-top: 4px; }
}
</style>
