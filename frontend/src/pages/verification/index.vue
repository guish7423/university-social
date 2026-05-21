<script setup lang="ts">
import { ref, onMounted } from "vue"
import { sendVerificationCode, verifyCode, getVerificationStatus } from "@/api/verification"

const step = ref<"form" | "code" | "done">("form")
const email = ref("")
const school = ref("")
const code = ref("")
const loading = ref(false)
const verified = ref(false)
const schools = [
  "华中科技大学", "武汉大学", "华中师范大学",
  "武汉理工大学", "华中农业大学", "中南财经政法大学",
  "中国地质大学", "湖北大学", "武汉科技大学",
  "中南民族大学", "湖北工业大学", "武汉工程大学",
]

onMounted(async () => {
  try {
    const status = await getVerificationStatus()
    verified.value = status.verified
    if (status.verified) step.value = "done"
  } catch {}
})

async function handleSend() {
  if (!email.value || !school.value) return
  loading.value = true
  try {
    await sendVerificationCode(email.value, school.value)
    step.value = "code"
    uni.showToast({ title: "验证码已发送", icon: "success" })
  } catch {
    uni.showToast({ title: "发送失败", icon: "none" })
  }
  loading.value = false
}

async function handleVerify() {
  if (!code.value) return
  loading.value = true
  try {
    await verifyCode(email.value, code.value)
    verified.value = true
    step.value = "done"
    uni.showToast({ title: "认证成功", icon: "success" })
  } catch {
    uni.showToast({ title: "验证码错误", icon: "none" })
  }
  loading.value = false
}
</script>

<template>
  <view class="container">
    <view class="card">
      <view class="header-icon">
        <text class="icon-text">🎓</text>
      </view>
      <text class="title">校园认证</text>
      <text class="desc">使用学校邮箱验证身份，获得认证标识</text>

      <template v-if="step === 'form'">
        <u-input
          v-model="email"
          placeholder="学校邮箱 (如 name@hust.edu.cn)"
          border="surround"
          fontSize="28rpx"
          :customStyle="{ marginBottom: '24rpx', height: '80rpx' }"
        />
        <u-input
          v-model="school"
          placeholder="输入学校名称"
          border="surround"
          fontSize="28rpx"
          :customStyle="{ marginBottom: '24rpx', height: '80rpx' }"
        />
        <u-button
          type="primary"
          shape="circle"
          :loading="loading"
          :disabled="!email || !school"
          @click="handleSend"
        >发送验证码</u-button>
      </template>

      <template v-else-if="step === 'code'">
        <text class="info-text">验证码已发送至 {{ email }}</text>
        <u-input
          v-model="code"
          placeholder="输入6位验证码"
          border="surround"
          maxlength="6"
          fontSize="28rpx"
          :customStyle="{ marginBottom: '24rpx', height: '80rpx', textAlign: 'center' }"
        />
        <u-button
          type="primary"
          shape="circle"
          :loading="loading"
          :disabled="code.length !== 6"
          @click="handleVerify"
        >确认验证</u-button>
        <text class="back-link" @click="step = 'form'">更换邮箱</text>
      </template>

      <template v-else>
        <view class="success-icon">
          <text class="checkmark">✓</text>
        </view>
        <text class="success-text">认证通过</text>
        <text class="success-desc">您已获得 {{ school }} 校友认证标识</text>
      </template>
    </view>
  </view>
</template>

<style scoped>
.container {
  min-height: 100vh;
  background: var(--u-bg-color, #f3f4f6);
  display: flex;
  justify-content: center;
  padding: 80rpx 30rpx;
}

.card {
  background: #fff;
  border-radius: 24rpx;
  padding: 60rpx 40rpx;
  width: 100%;
  max-width: 600rpx;
  text-align: center;
}

.header-icon {
  margin-bottom: 20rpx;
}

.icon-text {
  font-size: 80rpx;
}

.title {
  font-size: 36rpx;
  font-weight: 700;
  color: #303133;
  display: block;
  margin-bottom: 12rpx;
}

.desc {
  font-size: 26rpx;
  color: #909399;
  display: block;
  margin-bottom: 48rpx;
}

.info-text {
  font-size: 24rpx;
  color: #667eea;
  display: block;
  margin-bottom: 32rpx;
}

.success-icon {
  width: 120rpx;
  height: 120rpx;
  border-radius: 60rpx;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 24rpx;
}

.checkmark {
  font-size: 60rpx;
  color: #fff;
  font-weight: bold;
}

.success-text {
  font-size: 32rpx;
  font-weight: 700;
  color: #303133;
  display: block;
  margin-bottom: 8rpx;
}

.success-desc {
  font-size: 26rpx;
  color: #909399;
  display: block;
}

.back-link {
  font-size: 24rpx;
  color: #667eea;
  display: block;
  margin-top: 24rpx;
  text-decoration: underline;
}
</style>
