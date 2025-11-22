package router

import (
	"github.com/saku-730/web-occurrence/backend/internal/handler"
	"github.com/saku-730/web-occurrence/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes は main.go で作成された gin.Engine にルートを登録するのだ
func SetupRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	workstationHandler *handler.WorkstationHandler,
	masterHandler *handler.MasterHandler,
	couchDBHandler *handler.CouchDBHandler,
) {
	// --- Public API グループ (認証不要) ---
	apiPublic := r.Group("/api")
	{
		apiPublic.POST("/register", userHandler.Register)
		apiPublic.POST("/login", userHandler.Login)
	}

	// --- Protected API グループ (認証ミドルウェアを使用)  ---
	apiProtected := r.Group("/api")
	apiProtected.Use(middleware.AuthMiddleware())
	{
		apiProtected.Any("/couchdb/*path", couchDBHandler.ProxyRequest)
		apiProtected.GET("/master-data", masterHandler.GetMasterData)

		apiProtected.GET("/users/me", userHandler.GetMe)
		
		apiProtected.POST("/workstation/create", workstationHandler.Create)
		apiProtected.GET("/my-workstations", workstationHandler.List) // /api/my-workstations と合わせるか検討が必要だが、一旦Handler定義に従う
		
		// フロントエンドからのリクエストに合わせてエンドポイントを追加・調整する場合はここで行うのだ
		// 例: apiProtected.GET("/my-workstations", workstationHandler.List) 
	}
}
