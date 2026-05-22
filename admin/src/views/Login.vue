
<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-logo">🏛️</div>
      <h1>校园社区管理系统</h1>
      <p class="login-desc">请通过主站登录后访问管理后台</p>
      <a :href="mainSiteUrl" class="login-btn">前往主站登录</a>
    </div>
  </div>
</template>

<script setup lang="ts">
const mainSiteUrl = location.origin.replace(/\/admin.*$/, "")
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
.login-logo { font-size: 64px; margin-bottom: 16px; }
h1 { color: #e8e4e0; margin-bottom: 12px; font-size: 20px; }
.login-desc { color: rgba(255,255,255,0.6); margin-bottom: 24px; }
.login-btn {
  display: inline-block;
  padding: 12px 32px;
  background: #c67a6a;
  color: #fff;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  transition: opacity .2s;
}
.login-btn:hover { opacity: .85; }
</style>