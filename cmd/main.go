// @title Game Integration API
// @version 1.0
// @description A Game Integration API for casino games with wallet management. Provides authentication, player information, bet placement (withdraw), bet settlement (deposit), and transaction cancellation endpoints.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description IMPORTANT: Enter your JWT token with "Bearer " prefix. Example: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyOSwiZXhwIjoxNzUxMjg4ODc0fQ.IwLr7sPvhXb_3HxI4d8F_UQinvJxc3ePfuM30ztMcdU
package main

import (
	_ "gameintegrationapi/docs"
	"gameintegrationapi/internal/delivery/http"
	"gameintegrationapi/internal/domain"
	"gameintegrationapi/internal/infrastructure"
	"gameintegrationapi/internal/repository"
	"gameintegrationapi/internal/usecase"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

func runSQLMigrations(db *gorm.DB, dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read migrations dir: %v", err)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			content, err := os.ReadFile(dir + "/" + file.Name())
			if err != nil {
				log.Fatalf("failed to read migration %s: %v", file.Name(), err)
			}
			if err := db.Exec(string(content)).Error; err != nil {
				log.Fatalf("failed to execute migration %s: %v", file.Name(), err)
			}
		}
	}
}

func main() {
	cfg := infrastructure.LoadConfig()
	infrastructure.Logger.Println("Loaded config")

	db, err := infrastructure.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	infrastructure.Logger.Println("Connected to DB", db.Name())

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Transaction{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	infrastructure.SeedTestUsers(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	// Initialize use cases
	walletClient := infrastructure.NewWalletClient(cfg.WalletURL, cfg.WalletToken)
	log.Printf("WalletClient initialized with URL: %s", cfg.WalletURL)

	authUseCase := usecase.NewAuthUseCase(userRepo)
	playerUseCase := usecase.NewPlayerUseCase(userRepo, walletClient)
	walletUseCase := usecase.NewWalletUseCase(userRepo, txRepo, db, walletClient)

	// Initialize handlers
	handlers := http.NewHandlers(authUseCase, playerUseCase, walletUseCase)

	// Setup router
	r := http.NewRouter(handlers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	infrastructure.Logger.Printf("Server starting on :%s", port)
	r.Run(":" + port)
}
