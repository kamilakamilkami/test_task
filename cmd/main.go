package main

import (
	"context"
	"log"
	"net/http"

	"project/config"
	"project/routes"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "project/docs"
)

// @title Project API
// @version 1.0
// @description REST API for your backend
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	// 🔁 Выполни миграции перед запуском
	config.RunMigrations()

	// Init DB
	ctx := context.Background()
	db := config.InitDB(ctx)
	defer db.Close()

	// Setup routes
	r := routes.SetupRoutes(db)

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start server
	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server error:", err)
	}
}
