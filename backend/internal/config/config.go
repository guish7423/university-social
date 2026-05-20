package config

import "os"

type Config struct {
	ServerAddr   string
	DatabaseURL  string
	RedisAddr    string
	JWTSecret    string
	WechatAppID  string
	WechatSecret string
	AdminUserIDs string
	UploadDir    string
}

func Load() *Config {
	return &Config{
		ServerAddr:   getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/university_social?sslmode=disable"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:    getEnv("JWT_SECRET", "change-me-in-production"),
		WechatAppID:  getEnv("WECHAT_APP_ID", ""),
		WechatSecret: getEnv("WECHAT_SECRET", ""),
		AdminUserIDs: getEnv("ADMIN_USER_IDS", ""),
		UploadDir:    getEnv("UPLOAD_DIR", "./uploads"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
