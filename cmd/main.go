// @title Game Integration API
// @version 1.0
// @description A Game Integration API for casino games with wallet management
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	_ "gameintegrationapi/docs"
	"gameintegrationapi/internal/delivery/http"
	"gameintegrationapi/internal/domain"
	"gameintegrationapi/internal/infrastructure"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
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
		&domain.Bet{},
		&domain.Transaction{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	runSQLMigrations(db, "./migrations")

	infrastructure.SeedTestUsers(db)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	r.GET("/swagger/*any", http.SwaggerHandler())

	r.POST("/auth/login", http.Login)

	// Protected routes
	r.GET("/profile", http.JWTAuthMiddleware(db), http.Profile)
	r.POST("/bet/withdraw", http.JWTAuthMiddleware(db), http.Withdraw)
	r.POST("/bet/deposit", http.JWTAuthMiddleware(db), http.Deposit)
	r.POST("/bet/cancel", http.JWTAuthMiddleware(db), http.Cancel)

	r.GET("/healthz", http.Healthz)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	infrastructure.Logger.Printf("Server starting on :%s", port)
	r.Run(":" + port)
}
