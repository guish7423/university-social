#!/usr/bin/env bash
set -euo pipefail

# ==========================================
# 大学社交平台 — 一键部署脚本
# 用法: ./deploy.sh [prod|dev]
# ==========================================

MODE="${1:-prod}"
COMPOSE_FILE="docker-compose.yml"
ENV_FILE=".env"

# 检查环境变量
if [ ! -f "$ENV_FILE" ]; then
  echo "⚠️  缺少 .env 文件，从 .env.example 复制..."
  cp .env.example "$ENV_FILE"
  echo "❗ 请编辑 $ENV_FILE 设置 JWT_SECRET 等敏感配置后重新运行"
  exit 1
fi

export $(grep -v '^\s*#' "$ENV_FILE" | xargs)

if [ -z "${JWT_SECRET:-}" ] || [ "$JWT_SECRET" = "change-me-to-a-random-string-in-production" ]; then
  echo "❌ 请设置 JWT_SECRET 为随机字符串（至少32位）"
  exit 1
fi

echo "🚀 开始部署 [$MODE]..."

if [ "$MODE" = "prod" ]; then
  echo "📦 构建生产镜像..."
  docker compose -f "$COMPOSE_FILE" build --parallel
  echo "⬆️  启动服务..."
  docker compose -f "$COMPOSE_FILE" up -d
else
  docker compose -f "$COMPOSE_FILE" up -d
fi

echo ""
echo "✅ 部署完成！"
echo "   H5 前端:  http://localhost"
echo "   管理后台: http://localhost/admin/"
echo "   API:      http://localhost/api/"
