package middleware

import (
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware はJWTトークンを検証し、user_id をコンテキストにセットするのだ
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Authorization ヘッダーを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization ヘッダーがありません"})
			c.Abort() // 処理を中止
			return
		}

		// 2. "Bearer " スキーマをチェック
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer スキーマが正しくありません"})
			c.Abort()
			return
		}

		// 3. トークンを検証 (infrastructureにヘルパーを呼ぶ)
		userID, err := infrastructure.ValidateAndExtractUserID(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}

		// 4. 検証成功！ user_id をGinのコンテキストに保存
		// これで、この後のHandlerが user_id を取り出せるようになる
		c.Set("user_id", userID)

		// 5. 次の処理（Handler）へ進む
		c.Next()
	}
}
