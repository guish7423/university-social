<template>
  <div class="login-page">
    <div class="login-card stagger-item">
      <div class="card-accent" />
      <div class="login-header">
        <!-- Campus gate SVG icon -->
        <svg class="logo-icon" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M8 56V28L32 12L56 28V56" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M20 56V36H44V56" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
          <rect x="26" y="44" width="12" height="12" rx="1.5" stroke="currentColor" stroke-width="2"/>
          <path d="M32 12V18" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
          <path d="M38 16L50 28" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" opacity="0.4"/>
          <path d="M26 16L14 28" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" opacity="0.4"/>
        </svg>
        <h1 class="logo-text">校园社</h1>
        <p class="subtitle">洛阳高校社区 · 连接每一个你</p>
      </div>
      <el-form @submit.prevent="handleLogin" class="login-form">
        <el-input
          v-model="nickname" ref="inputRef"
          placeholder="输入昵称登录（开发环境）" size="large" clearable
          @keyup.enter="handleLogin"
        >
          <template #prefix><el-icon><User /></el-icon></template>
        </el-input>
        <el-button type="primary" size="large" :loading="loading" native-type="submit" class="login-btn">
          {{ loading ? '登录中...' : '进入校园社' }}
        </el-button>
        <p class="keyboard-hint">
          <kbd>Enter</kbd> 键快速登录
        </p>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useUserStore } from "@/stores/user"
import { devLogin, getProfile } from "@/api/auth"
import { User } from "@element-plus/icons-vue"

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
@use "@/styles/variables" as *;

.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, oklch(0.10 0.01 30) 0%, oklch(0.15 0.025 25) 40%, oklch(0.10 0.01 30) 70%, oklch(0.14 0.015 20) 100%);
  background-size: 400% 400%;
  animation: gradient-shift 12s ease-in-out infinite;
  position: relative;
  overflow: hidden;

  &::before {
    content: "";
    position: absolute;
    inset: 0;
    background: radial-gradient(ellipse at 30% 40%, rgba(198,122,106,0.06) 0%, transparent 40%),
                radial-gradient(ellipse at 70% 60%, rgba(94,143,212,0.04) 0%, transparent 40%);
    pointer-events: none;
  }
}

@keyframes gradient-shift {
  0%, 100% { background-position: 0% 50%; }
  33% { background-position: 100% 50%; }
  66% { background-position: 50% 100%; }
}

.login-card {
  width: 420px;
  padding: 48px 44px 40px;
  background: oklch(0.17 0.015 30);
  border-radius: 16px;
  border: 1px solid oklch(0.22 0.015 30);
  position: relative;
  overflow: hidden;

  .card-accent {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, transparent, oklch(0.62 0.12 16), oklch(0.70 0.10 16), oklch(0.62 0.12 16), transparent);
    opacity: 0.8;
  }
}

.login-header {
  text-align: center;
  margin-bottom: 36px;

  .logo-icon {
    width: 48px;
    height: 48px;
    color: oklch(0.62 0.12 16);
    margin-bottom: 16px;
  }

  .logo-text {
    font-size: 28px;
    font-weight: 700;
    color: oklch(0.92 0.005 30);
    margin-bottom: 8px;
  }

  .subtitle {
    font-size: 14px;
    color: oklch(0.55 0.01 30);
  }
}

.login-form {
  .el-input {
    margin-bottom: 20px;
  }

  .login-btn {
    width: 100%;
    background: oklch(0.62 0.12 16);
    border-color: oklch(0.62 0.12 16);
    font-size: 15px;
    &:hover { background: oklch(0.70 0.10 16); border-color: oklch(0.70 0.10 16); }
    &:active { transform: translateY(1px); }
  }

  .keyboard-hint {
    text-align: center;
    margin-top: 18px;
    font-size: 12px;
    color: oklch(0.45 0.01 30);

    kbd {
      display: inline-block;
      padding: 2px 8px;
      font-size: 11px;
      font-family: inherit;
      background: oklch(0.20 0.015 30);
      border: 1px solid oklch(0.26 0.015 30);
      border-radius: 4px;
      color: oklch(0.65 0.01 30);
      margin: 0 2px;
    }
  }
}
</style>
