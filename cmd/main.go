package main

import (
	"inventory-service-go/database"
	"inventory-service-go/internal/delivery/http"
	"inventory-service-go/internal/repository"
	"inventory-service-go/internal/service"
	"inventory-service-go/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	db := database.ConnectDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	r := gin.Default()

	repo := repository.NewRepository()
	uc := usecase.NewUsecase(repo)
	emailService := service.NewEmailService()
	handler := http.NewHandler(uc, emailService)

	http.SetupRoutes(r, handler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
