package router

import (
	"github.com/saku-730/web-occurrence/backend/internal/handler"
	"github.com/saku-730/web-occurrence/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)


func SetupRouter(
	userHandler *handler.UserHandler,
	couchDBHandler *handler.CouchDBHandler, // (今回追加)
) *gin.Engine {
	r := gin.Default()

	// --- Public API グループ (認証不要) ---
	apiPublic := r.Group("/api")
	{
		// ユーザー登録エンドポイント
		apiPublic.POST("/register", userHandler.Register)

		// ログインエンドポイント
		apiPublic.POST("/login", userHandler.Login)
	}

	// --- Protected API グループ (認証ミドルウェアを使用)  ---
	apiProtected := r.Group("/api")
	apiProtected.Use(middleware.AuthMiddleware()) // ★このグループはJWT認証が必須
	{
		// CouchDBセッション発行エンドポイント
		// GET /api/couchdb-session
		apiProtected.GET("/couchdb-session", couchDBHandler.GetCouchDBSession)

		// (ここに /api/upload-image などの認証必須APIを追加していく)
	}

	return r
}
