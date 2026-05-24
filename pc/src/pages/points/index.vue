<template>
  <div class="points-page">
    <div class="daily-card">
      <div class="card-left">
        <div class="label">当前积分</div>
        <div class="balance">{{ balance }}</div>
        <div class="streak" v-if="streak > 0">🔥 已连续签到 {{ streak }} 天</div>
      </div>
      <div class="card-right">
        <el-button
          type="primary"
          :disabled="claimedToday"
          @click="handleDailyClaim"
          :class="{ 'is-claimed': claimedToday }"
        >
          {{ claimedToday ? '今日已签到 ✓' : '签到 +5' }}
        </el-button>
      </div>
    </div>

    <div class="section history-section">
      <h2>积分记录</h2>
      <div v-if="logsLoading" class="loading-wrap"><el-skeleton :rows="3" animated /></div>
      <div v-else-if="!logs.length" class="empty-state"><el-empty description="暂无积分记录" /></div>
      <div v-else class="history-list">
        <div v-for="log in logs" :key="log.id" class="history-item stagger-item">
          <div class="log-info">
            <span class="log-reason">{{ log.reason }}</span>
            <span class="log-time">{{ formatTime(log.created_at) }}</span>
          </div>
          <span class="log-amount" :class="log.amount > 0 ? 'positive' : 'negative'">
            {{ log.amount > 0 ? '+' : '' }}{{ log.amount }}
          </span>
        </div>
      </div>
    </div>

    <div class="section">
      <h2>积分排行</h2>
      <div v-if="loading" class="loading-wrap"><el-skeleton :rows="5" animated /></div>
      <div v-else-if="!leaderboard.length" class="empty-state"><el-empty description="暂无数据" /></div>
      <div v-else class="leaderboard">
        <div
          v-for="(u, i) in leaderboard"
          :key="u.id"
          class="lb-row stagger-item"
          :class="{ 'is-self': u.id === userId }"
          @click="$router.push('/user/' + u.id)"
        >
          <span class="rank" :class="{ top3: i < 3 }">
            <template v-if="i === 0">🥇</template>
            <template v-else-if="i === 1">🥈</template>
            <template v-else-if="i === 2">🥉</template>
            <template v-else>{{ i + 1 }}</template>
          </span>
          <el-avatar :size="32" :src="u.avatar" />
          <span class="name">{{ u.nickname }}</span>
          <span class="points">{{ u.points }} 分</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useUserStore } from "@/stores/user"
import { getBalance, claimDailyLogin, getLogs, getLeaderboard } from "@/api/points"
import type { PointsLog } from "@/api/points"
import dayjs from "dayjs"

const userStore = useUserStore()
const userId = userStore.userId

const balance = ref(0)
const leaderboard = ref<{ id: number; nickname: string; avatar: string; points: number }[]>([])
const logs = ref<PointsLog[]>([])
const claimedToday = ref(false)
const loading = ref(true)
const logsLoading = ref(true)

const streak = computed(() => {
  const checkins = logs.value.filter((l) => l.reason === "每日签到")
  let count = 0
  const today = dayjs()
  for (let i = 0; i < checkins.length; i++) {
    if (dayjs(checkins[i].created_at).isSame(today.subtract(i, "day"), "day")) {
      count++
    } else break
  }
  return count
})

function formatTime(t: string) {
  return dayjs(t).format("MM-DD HH:mm")
}

async function handleDailyClaim() {
  try {
    const res = await claimDailyLogin()
    balance.value += res.amount
    claimedToday.value = true
    logs.value = await getLogs()
  } catch {
    claimedToday.value = true
  }
}

onMounted(async () => {
  try {
    const [bal, lb, lg] = await Promise.all([getBalance(), getLeaderboard(), getLogs()])
    balance.value = bal.balance
    leaderboard.value = lb
    logs.value = lg
    const today = dayjs().format("YYYY-MM-DD")
    claimedToday.value = lg.some(
      (l) => l.reason === "每日签到" && dayjs(l.created_at).format("YYYY-MM-DD") === today
    )
  } catch {
    /* handled by interceptor */
  } finally {
    loading.value = false
    logsLoading.value = false
  }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

 margin: 0 auto;
.points-page { max-width: 640px; }

.daily-card {
  display: flex; align-items: center; justify-content: space-between;
  background: linear-gradient(135deg, $bg-card, rgba($brand-primary-hex, 0.15));
  border: 1px solid $border-color; border-radius: $radius-lg;
  padding: 24px 28px; margin-bottom: 24px;

  .card-left {
    .label { font-size: 13px; color: $text-secondary; margin-bottom: 2px; }
    .balance { font-size: 36px; font-weight: 700; color: $primary; line-height: 1.2; }
    .streak { font-size: 13px; color: $accent-gold; margin-top: 6px; font-weight: 500; }
  }
  .card-right {
    .el-button.is-claimed {
      background: $color-success; border-color: $color-success;
      &:hover { background: $color-success; border-color: $color-success; }
    }
  }
}

.section {
  margin-bottom: 24px;
  h2 { font-size: 17px; font-weight: 600; margin-bottom: 14px; }
}

.history-list {
  display: flex; flex-direction: column; gap: 2px;
}

.history-item {
  display: flex; align-items: center; justify-content: space-between;
  padding: 10px 12px; border-radius: $radius-sm;
  transition: background 0.2s;
  &:hover { background: $bg-card; }

  .log-info {
    display: flex; flex-direction: column; gap: 2px;
    .log-reason { font-size: 14px; font-weight: 500; }
    .log-time { font-size: 11px; color: $text-muted; }
  }
  .log-amount {
    font-size: 15px; font-weight: 700;
    &.positive { color: $color-success; }
    &.negative { color: $color-danger; }
  }
}

.leaderboard { display: flex; flex-direction: column; gap: 6px; }

.lb-row {
  display: flex; align-items: center; gap: 12px;
  padding: 10px; border-radius: $radius-sm; cursor: pointer;
  transition: background 0.2s;
  &:hover { background: $bg-card; }
  &.is-self { background: rgba($brand-primary-hex, 0.08); }

  .rank {
    width: 28px; font-size: 18px; font-weight: 700; color: $text-muted; text-align: center;
    &.top3 { color: $accent-gold; }
  }
  .name { flex: 1; font-size: 14px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .points { font-size: 13px; color: $primary; font-weight: 600; white-space: nowrap; }
}
</style>
