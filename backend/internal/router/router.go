package router

import (
	"myapp/handler" // (sのプロジェクト名)

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
