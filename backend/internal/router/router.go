package router

import (
	"github.com/saku-730/web-occurrence/backend/handler"

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
		
		// (ここに /api/login などを追加していく)
	}

	return r
}
