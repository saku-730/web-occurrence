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

	// DI (依存性注入)
	couchClient := infrastructure.NewCouchDBClient()
	userRepo := repository.NewUserRepository(db)
	
	// User Service
	userService := service.NewUserService(userRepo, couchClient)
	
	// CouchDB Service (引数が増えたのだ！)
	couchDBService := service.NewCouchDBService(userRepo, couchClient, couchSecret, couchURL)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	couchDBHandler := handler.NewCouchDBHandler(couchDBService)

	// Router
	r := router.SetupRouter(userHandler, couchDBHandler)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("サーバーを起動中 (http://localhost:%s)...", port)
	r.Run(":" + port)
}
