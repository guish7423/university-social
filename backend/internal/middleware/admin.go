package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/config"
)

func AdminRequired(cfg *config.Config) gin.HandlerFunc {
	adminSet := make(map[string]struct{})
	for _, id := range strings.Split(cfg.AdminUserIDs, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			adminSet[id] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")
		if _, ok := adminSet[strconv.FormatInt(userID, 10)]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "无管理员权限"})
			return
		}
		c.Next()
	}
}
