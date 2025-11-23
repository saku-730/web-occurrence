package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// [DEBUG] ヘッダーの確認ログ
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("[AuthMiddleware] Error: Authorization header is empty")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Bearer トークンの形式チェック
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("[AuthMiddleware] Error: Invalid header format. Got: %s\n", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// トークンの検証
		userID, err := infrastructure.ValidateToken(tokenString)
		if err != nil {
			// [DEBUG] 検証失敗の理由を出力（これが知りたかった！）
			fmt.Printf("[AuthMiddleware] Error: Token validation failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 成功！コンテキストにセット
		c.Set("user_id", userID)
		c.Next()
	}
}
