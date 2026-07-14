package main

import (
	"log"
	"net/http"
	"os"

	"icon_exchange/internal/config"
	"icon_exchange/internal/module_mailer"
	"icon_exchange/internal/module_user"
	"icon_exchange/internal/module_user/model"
	"icon_exchange/internal/module_user/repository"
	"icon_exchange/internal/module_user/service"
	"icon_exchange/internal/module_wallet"
	"icon_exchange/internal/shared/database"

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
			&module_wallet.Wallet{},
		)
		if err != nil {
			log.Fatalf("AutoMigrate failed: %v", err)
		}
	}

	// 3. Initialize Modules & Dependency Injection
	// -> Wallet Module
	walletRepo := module_wallet.NewRepository(db)
	walletService := module_wallet.NewService(walletRepo)
	walletHandler := module_wallet.NewHandler(walletService)

	// -> Mailer Module
	mailerService := module_mailer.NewService()

	// -> User Module: Repositories
	userRepo := repository.NewUserRepo(db)
	otpRepo := repository.NewOTPRepo(db)
	sessionRepo := repository.NewSessionRepo(db)

	// -> User Module: Auth Service (wires together repos + mailer)
	authService := service.NewAuthService(db, userRepo, otpRepo, sessionRepo, mailerService)

	// -> User Module: Handler
	userHandler := module_user.NewHandler(authService)

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
		walletHandler.RegisterRoutes(v1)
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
