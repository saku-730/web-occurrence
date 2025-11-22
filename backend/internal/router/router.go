package router

import (
	"github.com/saku-730/web-occurrence/backend/internal/handler"
	"github.com/saku-730/web-occurrence/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	couchDBHandler *handler.CouchDBHandler,
	masterHandler *handler.MasterHandler,
	workstationHandler *handler.WorkstationHandler,
) *gin.Engine {
	r := gin.Default()

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
		
		apiProtected.POST("/workstation/create", workstationHandler.Create)
		// ▼ 追加: ワークステーション一覧
		apiProtected.GET("/workstations", workstationHandler.List)
	}

	return r
}
