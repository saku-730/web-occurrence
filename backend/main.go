package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/saku-730/web-occurrence/backend/internal/handler"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
	"github.com/saku-730/web-occurrence/backend/internal/router"
	"github.com/saku-730/web-occurrence/backend/internal/service"
)

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Println("警告: .env ファイルが読み込めませんでした。")
		}
	}
	// DB接続
	db, sqlDB := infrastructure.InitDB()
	defer sqlDB.Close()

	// 設定値の取得
	couchURL := os.Getenv("COUCHDB_URL")
	couchSecret := os.Getenv("COUCHDB_SECRET")

	// Repository
	couchClient := infrastructure.NewCouchDBClient()
	userRepo := repository.NewUserRepository(db)
	masterRepo := repository.NewMasterRepository(db)
	wsRepo := repository.NewWorkstationRepository(db) // WorkstationRepo
	
	// Service
	couchDBService := service.NewCouchDBService(userRepo, couchClient, couchSecret, couchURL)
	userService := service.NewUserService(userRepo, couchClient)
	
	// MasterService (wsRepoも必要になったので追加)
	masterService := service.NewMasterService(masterRepo, wsRepo) 
	
	// WorkstationService (wsRepo, masterRepo, couchClientが必要)
	workstationService := service.NewWorkstationService(wsRepo, masterRepo, couchClient)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	couchDBHandler := handler.NewCouchDBHandler(couchDBService)
	masterHandler := handler.NewMasterHandler(masterService)
	workstationHandler := handler.NewWorkstationHandler(workstationService) // 作成

	// Router (引数を4つ渡す)
	r := router.SetupRouter(userHandler, couchDBHandler, masterHandler, workstationHandler)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("サーバーを起動中 (http://localhost:%s)...", port)
	r.Run(":" + port)
}
