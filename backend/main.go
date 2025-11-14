package main

import (
	"log"
	"github.com/saku-730/web-occurrence/backend/internal/handler"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
	"github.com/saku-730/web-occurrence/backend/internal/router"
	"github.com/saku-730/web-occurrence/backend/internal/service"
	"os"

	"github.com/joho/godotenv"
)


func main() {
	// 1. 環境変数を .env ファイルから読み込む
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Println("警告: .env ファイルが読み込めませんでした。")
		}
	}

	// 2. データベース接続 (infrastructure)
	db, sqlDB := infrastructure.InitDB()
	defer sqlDB.Close()

	// 3. 依存性注入 (DI)


	// --- Infrastructure層 ---
	couchClient := infrastructure.NewCouchDBClient()
	
	// --- Repository層 ---
	userRepo := repository.NewUserRepository(db)
	
	// --- Service層 ---
	// UserService に userRepo と couchClient の両方を渡す
	userService := service.NewUserService(userRepo, couchClient) 
	couchDBService := service.NewCouchDBService(userRepo, couchClient)

	// --- Handler層 ---
	userHandler := handler.NewUserHandler(userService)
	couchDBHandler := handler.NewCouchDBHandler(couchDBService)



	// 4. ルーターの設定 (Handlerを渡す)
	r := router.SetupRouter(userHandler, couchDBHandler)

	// 5. サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("サーバーを起動中 (http://localhost:%s)...", port)
	r.Run(":" + port)
}
