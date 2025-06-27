package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gameintegrationapi/internal/infrastructure"
)

func main() {
	cfg := infrastructure.LoadConfig()
	infrastructure.Logger.Println("Loaded config")

	db, err := infrastructure.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	infrastructure.Logger.Println("Connected to DB", db.Name())

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	infrastructure.Logger.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
