<template>
  <div class="points-page">
    <div class="points-banner">
      <div class="balance-wrap">
        <div class="label">当前积分</div>
        <div class="balance">{{ balance }}</div>
      </div>
      <div class="actions">
        <el-button type="primary" :disabled="claimedToday" @click="handleDailyClaim">
          {{ claimedToday ? '已签到' : '每日签到' }}
        </el-button>
      </div>
    </div>

    <div class="section">
      <h2>排行榜</h2>
      <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
      <div v-else-if="!leaderboard.length" class="empty-state"><el-empty description="暂无数据" /></div>
      <div v-else class="leaderboard">
        <div v-for="(u, i) in leaderboard" class="lb-row stagger-item" :key="u.id" @click="$router.push('/user/' + u.id)">
          <span class="rank" :class="{ top3: i < 3 }">{{ i + 1 }}</span>
          <el-avatar :size="32" :src="u.avatar" />
          <span class="name">{{ u.nickname }}</span>
          <span class="points">{{ u.points }} 分</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { getLeaderboard, claimDailyLogin } from "@/api/points"

const balance = ref(0)
const leaderboard = ref<{ id: number; nickname: string; avatar: string; points: number }[]>([])
const claimedToday = ref(false)
const loading = ref(true)

async function handleDailyClaim() {
  try {
    const res = await claimDailyLogin()
    balance.value += res.amount
    claimedToday.value = true
  } catch { /* handled by interceptor */ }
}

onMounted(async () => {
  try {
    const [lb] = await Promise.all([getLeaderboard()])
    leaderboard.value = lb
  } catch { /* handled by interceptor */ } finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.points-page { max-width: 640px; }

.points-banner {
  display: flex; align-items: center; justify-content: space-between;
  background: linear-gradient(135deg, $bg-card, rgba($brand-primary-hex, 0.15));
  border: 1px solid $border-color; border-radius: $radius-md;
  padding: 28px; margin-bottom: 24px;
  .balance-wrap {
    .label { font-size: 13px; color: $text-secondary; }
    .balance { font-size: 36px; font-weight: 700; color: $primary; margin-top: 4px; }
  }
}

.section {
  h2 { font-size: 17px; font-weight: 600; margin-bottom: 14px; }
}

.leaderboard { display: flex; flex-direction: column; gap: 6px; }
.lb-row {
  display: flex; align-items: center; gap: 12px;
  padding: 10px; border-radius: $radius-sm; cursor: pointer;
  &:hover { background: $bg-card; }
  .rank {
    width: 24px; font-size: 14px; font-weight: 700; color: $text-muted; text-align: center;
    &.top3 { color: $accent-gold; }
  }
  .name { flex: 1; font-size: 14px; font-weight: 600; }
  .points { font-size: 13px; color: $primary; font-weight: 600; }
}
</style>
