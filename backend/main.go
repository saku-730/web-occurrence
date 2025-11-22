package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/saku-730/web-occurrence/backend/internal/handler"
	"github.com/saku-730/web-occurrence/backend/internal/infrastructure"
	"github.com/saku-730/web-occurrence/backend/internal/middleware"
	"github.com/saku-730/web-occurrence/backend/internal/model"
	"github.com/saku-730/web-occurrence/backend/internal/repository"
	"github.com/saku-730/web-occurrence/backend/internal/router"
	"github.com/saku-730/web-occurrence/backend/internal/service"
)

func main() {
	// 1. Initialize Database
	db, err := infrastructure.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2. Initialize CouchDB Client
	couchConfig := &model.CouchDBConfig{
		URL:       os.Getenv("COUCHDB_URL"),
		Secret:    os.Getenv("COUCHDB_SECRET"),
		AdminUser: os.Getenv("COUCHDB_USER"),
		AdminPass: os.Getenv("COUCHDB_PASSWORD"),
	}
	if couchConfig.URL == "" {
		couchConfig.URL = "http://localhost:5984" // Default
	}
	couchClient := infrastructure.NewCouchDBClient(couchConfig)

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	wsRepo := repository.NewWorkstationRepository(db)
	masterRepo := repository.NewMasterRepository(db)

	// 4. Initialize Services
	authService := service.NewUserService(userRepo, wsRepo)
	wsService := service.NewWorkstationService(wsRepo, masterRepo, couchClient)
	masterService := service.NewMasterService(masterRepo)
	couchService := service.NewCouchDBService(couchClient)

	// ★ 修正: SyncServiceにwsRepoを渡す
	syncService := service.NewSyncService(db, couchClient, wsRepo)

	// 5. Start Sync Polling (Background)
	syncService.StartPolling()

	// 6. Initialize Handlers
	userHandler := handler.NewUserHandler(authService)
	wsHandler := handler.NewWorkstationHandler(wsService)
	masterHandler := handler.NewMasterHandler(masterService)
	couchHandler := handler.NewCouchDBHandler(couchService)

	// 7. Setup Router
	r := gin.Default()

	// CORS Setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Define Routes
	router.SetupRoutes(r, userHandler, wsHandler, masterHandler, couchHandler, middleware.AuthMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
