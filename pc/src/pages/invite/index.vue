<template>
  <div class="invite-page">
    <div class="page-header">
      <h1>邀请码</h1>
      <el-button type="primary" :loading="creating" @click="handleCreate">生成邀请码</el-button>
    </div>

    <el-card shadow="never" class="redeem-card">
      <h3>兑换邀请码</h3>
      <div class="redeem-input">
        <el-input v-model="redeemCode" placeholder="输入邀请码" size="large" />
        <el-button type="primary" :loading="redeeming" :disabled="!redeemCode.trim()" @click="handleRedeem">兑换</el-button>
      </div>
    </el-card>

    <div class="section-header"><h2>我的邀请码</h2></div>
    <div v-if="loading" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="!codes.length" class="empty-state"><el-empty description="暂无邀请码" /></div>
    <el-table v-else :data="codes" stripe>
      <el-table-column prop="code" label="邀请码">
        <template #default="{row}"><el-tag>{{ row.code }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column label="状态">
        <template #default="{row}">
          <el-tag :type="row.is_used ? 'info' : 'success'">{{ row.is_used ? '已使用' : '未使用' }}</el-tag>
        </template>
      </el-table-column>
    </el-table>
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
@use "@/styles/variables.scss" as *;
.invite-page { max-width: 700px; }
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; h1 { font-size: 22px; font-weight: 700; } }
.redeem-card { margin-bottom: 24px; h3 { font-size: 15px; margin-bottom: 12px; } }
.redeem-input { display: flex; gap: 10px; }
.section-header { margin-bottom: 14px; h2 { font-size: 17px; font-weight: 600; } }
.loading-wrap, .empty-state { padding: 40px 0; }
</style>
