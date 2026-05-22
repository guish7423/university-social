<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="logo"><span class="logo-dot">●</span> 校园社</div>
        <p class="subtitle">洛阳高校社区 · 连接每一个你</p>
      </div>
      <el-form @submit.prevent="handleLogin">
        <el-input v-model="nickname" placeholder="输入昵称登录（开发环境）" size="large" />
        <el-button type="primary" size="large" :loading="loading" native-type="submit" class="login-btn">
          {{ loading ? '登录中...' : '进入校园社' }}
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useUserStore } from "@/stores/user"
import { devLogin, getProfile } from "@/api/auth"

const router = useRouter()
const userStore = useUserStore()
const nickname = ref("")
const loading = ref(false)

async function handleLogin() {
  if (!nickname.value.trim()) return
  loading.value = true
  try {
    const res = await devLogin(nickname.value.trim())
    userStore.setAuth(res.token, res.user)
    router.push("/home")
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">

.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, oklch(0.11 0.01 30) 0%, oklch(0.15 0.02 25) 50%, oklch(0.11 0.01 30) 100%);
}

.login-card {
  width: 420px;
  padding: 56px 44px 48px;
  background: oklch(0.17 0.015 30);
  border-radius: 16px;
  border: 1px solid oklch(0.22 0.015 30);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
  .logo {
    font-size: 30px;
    font-weight: 700;
    color: oklch(0.92 0.005 30);
    .logo-dot { color: oklch(0.62 0.12 16); }
  }
  .subtitle {
    font-size: 14px;
    color: oklch(0.55 0.01 30);
    margin-top: 10px;
  }
}

.login-btn {
  width: 100%;
  margin-top: 28px;
  background: oklch(0.62 0.12 16);
  border-color: oklch(0.62 0.12 16);
  &:hover { background: oklch(0.70 0.10 16); border-color: oklch(0.70 0.10 16); }
}
</style>
