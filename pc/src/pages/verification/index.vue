<template>
  <div class="verif-page">
    <h1>实名认证</h1>
    <div v-if="status !== null" class="status-banner">
      <el-alert v-if="status" title="已认证" type="success" :closable="false" description="你已完成实名认证" show-icon />
      <el-alert v-else title="未认证" type="warning" :closable="false" description="完成认证后可获得更多功能权限" show-icon />
    </div>

    <div v-if="!status" class="verif-form">
      <el-form @submit.prevent="handleSendCode">
        <el-form-item label="手机号">
          <el-input v-model="phone" placeholder="输入手机号" />
        </el-form-item>
        <el-form-item label="验证码">
          <div class="code-row">
            <el-input v-model="code" placeholder="输入验证码" />
            <el-button :disabled="sending || countdown > 0" @click="handleSendCode">
              {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>
        <el-button type="primary" :loading="submitting" :disabled="!code.trim()" @click="handleVerify">
          提交认证
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { sendCode, verify, verificationStatus } from "@/api/verification"

const status = ref<boolean | null>(null)
const phone = ref("")
const code = ref("")
const sending = ref(false)
const submitting = ref(false)
const countdown = ref(0)

async function handleSendCode() {
  if (!phone.value.trim()) return
  sending.value = true
  try {
    await sendCode({ phone: phone.value.trim() })
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) clearInterval(timer)
    }, 1000)
  } catch { /* handled by interceptor */ } finally { sending.value = false }
}

async function handleVerify() {
  if (!code.value.trim()) return
  submitting.value = true
  try {
    await verify({ code: code.value.trim(), phone: phone.value.trim() })
    status.value = true
  } catch { /* handled by interceptor */ } finally { submitting.value = false }
}

onMounted(async () => {
  try {
    const res = await verificationStatus()
    status.value = res.is_verified
  } catch { /* handled by interceptor */ }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.verif-page { max-width: 500px; }
h1 { font-size: 20px; font-weight: 700; margin-bottom: 20px; }
.status-banner { margin-bottom: 20px; }
.verif-form { margin-top: 20px; }
.code-row { display: flex; gap: 10px; }
</style>
