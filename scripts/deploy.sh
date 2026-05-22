#!/bin/bash
set -euo pipefail

# ========================================
# 大学社交平台 — 一键部署脚本
# 用法: bash scripts/deploy.sh [env]
# ========================================

ENV=${1:-production}
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
echo "🚀 部署模式: $ENV"
echo "📂 项目目录: $ROOT_DIR"

cd "$ROOT_DIR"

# 1. 环境变量
if [ ! -f .env ]; then
  if [ -f .env.example ]; then
    echo "📝 从 .env.example 创建 .env"
    cp .env.example .env
  else
    echo "❌ .env 不存在，中止"
    exit 1
  fi
fi
source .env 2>/dev/null || true

# 2. 运行数据库迁移
echo "🗄️  运行数据库迁移..."
cd backend && go run cmd/migrate/main.go && cd "$ROOT_DIR"

# 3. 构建前端
if [ -d frontend ]; then
  echo "🎨 构建前端..."
  cd frontend && npm install --silent && npm run build && cd "$ROOT_DIR"
fi

# 4. 构建管理后台
if [ -d admin ]; then
  echo "⚙️  构建管理后台..."
  cd admin && npm install --silent && npm run build && cd "$ROOT_DIR"
fi

# 5. Docker Compose
if [ "$ENV" = "production" ]; then
  echo "🐳 启动 Docker Compose..."
  docker compose up -d --build

  echo "📋 服务状态:"
  docker compose ps

  echo "✅ 部署完成！"
  echo "   API: http://localhost:8081"
  echo "   前端: http://localhost:3000 (需通过 Nginx 或小程序开发者工具)"
else
  echo "✅ 本地部署准备完成。"
  echo "   运行: cd backend && go run ./cmd/server/"
fi
