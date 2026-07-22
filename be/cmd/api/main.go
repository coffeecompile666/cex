package main

import (
	template2 "icon_exchange/internal/asset"
	repository3 "icon_exchange/internal/asset/repository"
	template "icon_exchange/internal/market"
	model2 "icon_exchange/internal/market/model"
	repository2 "icon_exchange/internal/market/repository"
	"icon_exchange/internal/matching_engine"
	"log"
	"net/http"
	"os"

	"icon_exchange/internal/config"
	"icon_exchange/internal/mailer"
	"icon_exchange/internal/shared/database"
	"icon_exchange/internal/user"
	"icon_exchange/internal/user/model"
	"icon_exchange/internal/user/repository"
	"icon_exchange/internal/user/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration
	config.LoadConfig()

	// 2. Connect to Database
	db := database.ConnectDB()
	if db != nil {
		log.Println("Running AutoMigrate...")
		err := db.AutoMigrate(
			&model.User{},
			&model.OTP{},
			&model.Session{},
			&model2.Market{},
		)
		if err != nil {
			log.Fatalf("AutoMigrate failed: %v", err)
		}
	}

	// -> Mailer Module
	mailerService := mailer.NewService()

	// -> User Module: Repositories
	userRepo := repository.NewUserRepo(db)
	otpRepo := repository.NewOTPRepo(db)
	sessionRepo := repository.NewSessionRepo(db)

	// -> User Module: Auth Service (wires together repos + mailer)
	authService := service.NewAuthService(db, userRepo, otpRepo, sessionRepo, mailerService)

	// -> User Module: Handler
	userHandler := user.NewHandler(authService)

	// -> Market Module: Handler
	marketRepo := repository2.NewMarketRepo(db)
	//marketService := service3.NewMarketService(marketRepo)
	marketHandler := template.NewMarketHandler(marketRepo)

	// -> Asset Module: Handler
	assetRepo := repository3.NewAssetRepository(db)
	assetHandler := template2.NewAssetHandler(assetRepo)

	// -> Start Matching engine
	matchingEngine := matching_engine.NewMatchingEngine()
	matchingEngine.Start()

	// 4. Setup Router
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		// Register module routes
		userHandler.RegisterRoutes(v1)

		// Create a group for routes that require authentication
		authorized := v1.Group("/")
		authorized.Use(user.AuthMiddleware())
		{
			// Register module asset
			assetHandler.RegisterRoutes(authorized)
		}

		// Register module market (public routes)
		marketHandler.RegisterRoutes(v1)
	}

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
