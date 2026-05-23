<template>
  <div class="login-page">
    <div class="login-card">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="#c67a6a" stroke-width="1.5">
        <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
      </svg>
      <h1>校园社区管理系统</h1>
      <el-form :model="form" class="login-form" @submit.prevent="handleLogin">
        <el-input v-model="form.username" placeholder="管理员账号" :prefix-icon="User" class="login-input" />
        <el-input v-model="form.password" type="password" placeholder="密码" :prefix-icon="Lock" class="login-input" show-password />
        <el-button type="primary" class="login-btn" :loading="logining" native-type="submit">登录</el-button>
      </el-form>
      <el-divider><span class="divider-text">或者</span></el-divider>
      <button class="dev-login-btn" @click="devLogin">开发模式快捷登录</button>
      <p class="login-hint">生产环境请使用管理员账号登录</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue"
import { User, Lock } from "@element-plus/icons-vue"
import { ElMessage } from "element-plus"

const form = reactive({ username: "", password: "" })
const logining = ref(false)

async function handleLogin() {
  if (!form.username || !form.password) { ElMessage.warning("请输入账号和密码"); return }
  logining.value = true
  try {
    const res = await fetch("/api/v1/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ nickname: form.username, password: form.password })
    })
    if (!res.ok) { ElMessage.error("账号或密码错误"); return }
    const data = await res.json()
    if (data.token) {
      localStorage.setItem("token", data.token)
      localStorage.setItem("user", JSON.stringify(data.user))
      window.location.hash = "#/"
    }
  } catch {
    ElMessage.error("登录失败，请确保后端已启动")
  } finally { logining.value = false }
}

async function devLogin() {
  try {
    const res = await fetch("/api/v1/auth/dev-login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ nickname: "管理员", avatar: "" })
    })
    const data = await res.json()
    if (data.token) {
      localStorage.setItem("token", data.token)
      localStorage.setItem("user", JSON.stringify(data.user))
      window.location.hash = "#/"
    }
  } catch {
    ElMessage.error("登录失败，请确保后端已启动")
  }
}
</script>

<style scoped>
.login-page {
  display: flex; align-items: center; justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1e2a3a 0%, #2a3a4a 100%);
}
.login-card {
  text-align: center;
  padding: 48px 40px;
  background: rgba(255,255,255,0.08);
  border-radius: 16px;
  backdrop-filter: blur(8px);
  width: 380px;
}
.login-card svg { display: block; margin: 0 auto 16px; }
h1 { color: #e8e4e0; margin-bottom: 24px; font-size: 20px; }
.login-form { display: flex; flex-direction: column; gap: 16px; }
.login-input :deep(.el-input__wrapper) { background: rgba(255,255,255,0.1); box-shadow: none; }
.login-input :deep(.el-input__inner) { color: #e8e4e0; }
.login-input :deep(.el-input__prefix) { color: rgba(255,255,255,0.5); }
.login-btn { width: 100%; }
.dev-login-btn {
  display: inline-block; padding: 10px 24px;
  background: transparent; color: rgba(255,255,255,0.6);
  border: 1px solid rgba(255,255,255,0.2); border-radius: 8px;
  cursor: pointer; font-size: 13px; transition: all .2s;
}
.dev-login-btn:hover { color: #c67a6a; border-color: #c67a6a; }
.login-hint { margin-top: 12px; font-size: 12px; color: rgba(255,255,255,0.35); }
.divider-text { color: rgba(255,255,255,0.35); font-size: 12px; }
</style>
