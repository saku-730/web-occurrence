package router

import (
	"github.com/saku-730/web-occurrence/backend/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter は Gin のルーターを設定するのだ
func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// APIのグループ
	api := r.Group("/api")
	{
		// ユーザー登録エンドポイント
		// POST /api/register
		api.POST("/register", userHandler.Register)
		
		api.POST("/login", userHandler.Login)
	}

	return r
}
