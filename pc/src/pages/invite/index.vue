<template>
  <div class="invite-page">
    <div class="page-header">
      <h1>邀请码</h1>
      <el-button type="primary" :loading="creating" @click="handleCreate">生成邀请码</el-button>
    </div>

    <div class="redeem-card">
      <h3>兑换邀请码</h3>
      <div class="redeem-input">
        <el-input v-model="redeemCode" placeholder="输入邀请码" size="large" clearable />
        <el-button type="primary" :loading="redeeming" :disabled="!redeemCode.trim()" @click="handleRedeem">兑换</el-button>
      </div>
    </div>

    <div class="section-header"><h2>我的邀请码</h2></div>

    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="!codes.length" class="empty-state"><el-empty description="暂无邀请码" :image-size="80" /></div>
    <div v-else class="codes-grid">
      <div v-for="c in codes" :key="c.id" class="code-card stagger-item">
        <div class="code-value">{{ c.code }}</div>
        <div class="code-meta">
          <span class="code-time">{{ formatDateTime(c.created_at) }}</span>
          <el-tag :type="c.is_used ? 'info' : 'success'" size="small" effect="dark">
            {{ c.is_used ? '已使用' : '未使用' }}
          </el-tag>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { ElMessage } from "element-plus"
import { myInviteCodes, createInviteCode, redeemInviteCode } from "@/api/invite"
import type { InviteCodeData } from "@/api/invite"
import { useTimeFormat } from "@/composables/useTimeFormat"

const { formatDateTime } = useTimeFormat()
const codes = ref<InviteCodeData[]>([])
const loading = ref(true)
const creating = ref(false)
const redeeming = ref(false)
const redeemCode = ref("")

async function handleCreate() {
  creating.value = true
  try {
    await createInviteCode()
    ElMessage.success("邀请码生成成功")
    codes.value = await myInviteCodes()
  } catch { /* handled by interceptor */ }
  finally { creating.value = false }
}

async function handleRedeem() {
  if (!redeemCode.value.trim()) return
  redeeming.value = true
  try {
    await redeemInviteCode(redeemCode.value.trim())
    ElMessage.success("兑换成功")
    redeemCode.value = ""
    codes.value = await myInviteCodes()
  } catch { /* handled by interceptor */ }
  finally { redeeming.value = false }
}

onMounted(async () => {
  try { codes.value = await myInviteCodes() }
  catch { /* handled by interceptor */ }
  finally { loading.value = false }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables" as *;

 margin: 0 auto;
.invite-page { max-width: 700px; }

.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px;
  h1 { font-size: 22px; font-weight: 700; }
}

.redeem-card {
  background: $bg-card; border: 1px solid $border-color; border-radius: $radius-md;
  padding: 24px; margin-bottom: 24px;

  h3 { font-size: 15px; font-weight: 600; margin-bottom: 14px; }

  .redeem-input {
    display: flex; gap: 10px;
  }
}

.section-header {
  margin-bottom: 14px;
  h2 { font-size: 17px; font-weight: 600; }
}

.codes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}

.code-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; padding: 20px;
  transition: transform 0.2s, box-shadow 0.2s, border-color 0.2s;
  cursor: default;

  &:hover {
    border-color: $primary-light;
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0,0,0,0.3);
  }

  .code-value {
    font-family: "DM Mono", "SF Mono", "Fira Code", monospace;
    font-size: 18px;
    font-weight: 600;
    letter-spacing: 2px;
    color: $primary;
    margin-bottom: 12px;
  }

  .code-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .code-time { font-size: 12px; color: $text-muted; }
  }
}

.loading-wrap, .empty-state { padding: 40px 0; }
</style>
