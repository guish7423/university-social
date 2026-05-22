<template>
  <div class="login-page">
    <div class="login-card">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="#c67a6a" stroke-width="1.5">
        <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
      </svg>
      <h1>校园社区管理系统</h1>
      <p class="login-desc">请登录后访问管理后台</p>
      <button class="login-btn" @click="devLogin">开发模式登录</button>
      <p class="login-hint">开发环境直接登录，生产环境请通过主站认证</p>
    </div>
  </div>
</template>

<script setup lang="ts">
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
  } catch (e) {
    alert("登录失败，请确保后端已启动")
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1e2a3a 0%, #2a3a4a 100%);
}
.login-card {
  text-align: center;
  padding: 48px;
  background: rgba(255,255,255,0.08);
  border-radius: 16px;
  backdrop-filter: blur(8px);
}
.login-card svg { display: block; margin: 0 auto 16px; }
h1 { color: #e8e4e0; margin-bottom: 12px; font-size: 20px; }
.login-desc { color: rgba(255,255,255,0.6); margin-bottom: 24px; }
.login-btn {
  display: inline-block;
  padding: 12px 32px;
  background: #c67a6a;
  color: #fff;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
  font-weight: 600;
  transition: opacity .2s;
}
.login-btn:hover { opacity: .85; }
.login-hint {
  margin-top: 12px;
  font-size: 12px;
  color: rgba(255,255,255,0.35);
}
</style>
