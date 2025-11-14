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
		err := godotenv.Load() // .env ファイルを読み込む
		if err != nil {
			log.Println("警告: .env ファイルが読み込めませんでした。")
		}
	}

	// 2. データベース接続 (infrastructure)
	// gorm.DB (*db) と sql.DB (*sqlDB) の両方を受け取るように変更
	db, sqlDB := infrastructure.InitDB()
	// main関数終了時にDB接続を閉じるのは *sqlDB の方を使うのだ
	defer sqlDB.Close() 

	// 3. 依存性注入 (DI)
	// Repository層 (DB接続は *gorm.DB の方を渡す)
	userRepo := repository.NewUserRepository(db)

	// Service層 (Repositoryを渡す)
	userService := service.NewUserService(userRepo)

	// Handler層 (Serviceを渡す)
	userHandler := handler.NewUserHandler(userService)

	// 4. ルーターの設定 (Handlerを渡す)
	r := router.SetupRouter(
		userHandler,
	)

	// 5. サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}
	log.Printf("サーバーを起動中 (http://localhost:%s)...", port)
	r.Run(":" + port)
}
