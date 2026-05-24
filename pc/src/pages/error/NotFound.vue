<template>
  <div class="not-found-page">
    <div class="glow-orb glow-orb--1"></div>
    <div class="glow-orb glow-orb--2"></div>

    <div class="card">
      <div class="code-row">
        <span class="digit" v-for="(d, i) in digits" :key="i"
          :style="{ animationDelay: `${i * 0.12}s` }"
        >{{ d }}</span>
      </div>

      <div class="divider"></div>

      <p class="message">哎呀，页面迷路了</p>
      <p class="hint">可能已被删除、移动，或地址输入有误</p>

      <div class="actions">
        <el-button class="home-btn" @click="$router.push('/home')">
          <el-icon class="btn-icon"><HomeFilled /></el-icon>
          返回首页
        </el-button>
        <el-button class="back-btn" @click="$router.back()">
          <el-icon class="btn-icon"><ArrowLeft /></el-icon>
          后退一步
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { HomeFilled, ArrowLeft } from "@element-plus/icons-vue"

const digits = ["4", "0", "4"]
</script>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.not-found-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 70vh;
  position: relative;
  overflow: hidden;
  padding: $space-8;
}

// ── Glow orbs ──
.glow-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
  opacity: 0.35;

  &--1 {
    width: 320px; height: 320px;
    background: radial-gradient(circle, $brand-primary 0%, transparent 70%);
    top: 5%; left: -5%;
    animation: float 8s ease-in-out infinite;
  }
  &--2 {
    width: 260px; height: 260px;
    background: radial-gradient(circle, $campus-blue 0%, transparent 70%);
    bottom: 5%; right: -5%;
    animation: float 10s ease-in-out infinite reverse;
  }
}

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, -20px); }
}

// ── Card ──
.card {
  position: relative;
  z-index: 1;
  text-align: center;
  background: var(--bg-card);
  border: 1px solid var(--border-light);
  border-radius: $radius-xl;
  padding: $space-16 $space-12;
  margin: 0 auto;
  max-width: 440px;
  width: 100%;
  box-shadow: var(--shadow-lg);
  backdrop-filter: blur(12px);
}

// ── 404 digits ──
.code-row {
  display: flex;
  justify-content: center;
  gap: $space-4;
  margin-bottom: $space-6;
}

.digit {
  font-size: 5rem;
  font-weight: 800;
  line-height: 1;
  background: linear-gradient(135deg, $brand-primary 0%, $campus-gold 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: popIn 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

@keyframes popIn {
  0% {
    opacity: 0;
    transform: scale(0.3) translateY(20px);
  }
  100% {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

// ── Divider ──
.divider {
  width: 60px;
  height: 3px;
  background: linear-gradient(90deg, $brand-primary, $campus-gold);
  border-radius: 2px;
  margin: 0 auto $space-6;
}

// ── Text ──
.message {
  font-size: $text-xl;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 $space-2;
}

.hint {
  font-size: $text-sm;
  color: var(--text-muted);
  margin: 0 0 $space-8;
}

// ── Actions ──
.actions {
  display: flex;
  gap: $space-3;
  justify-content: center;
  flex-wrap: wrap;
}

.home-btn {
  background: linear-gradient(135deg, $brand-primary, $brand-primary-light) !important;
  border: none !important;
  color: #fff !important;
  padding: $space-3 $space-6;
  font-size: $text-sm;
  font-weight: 500;
  border-radius: $radius-full;
  transition: all 0.25s ease;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 16px rgba(198, 122, 106, 0.4);
  }

  .btn-icon { margin-right: $space-1; }
}

.back-btn {
  background: var(--bg-raised) !important;
  border: 1px solid var(--border-light) !important;
  color: var(--text-secondary) !important;
  padding: $space-3 $space-6;
  font-size: $text-sm;
  border-radius: $radius-full;
  transition: all 0.25s ease;

  &:hover {
    background: var(--bg-hover) !important;
    color: var(--text-primary) !important;
    transform: translateY(-1px);
  }

  .btn-icon { margin-right: $space-1; }
}

// ── Responsive ──
@media (max-width: 480px) {
  .card { padding: $space-10 $space-6; }
  .digit { font-size: 3.5rem; }
  .message { font-size: $text-lg; }
}
</style>
