package main

import (
	"log"

	"github.com/iamdebopriya/Inventory-Management-System-go/database"
	"github.com/iamdebopriya/Inventory-Management-System-go/internal/delivery/http"
	"github.com/iamdebopriya/Inventory-Management-System-go/internal/repository"
	"github.com/iamdebopriya/Inventory-Management-System-go/internal/service"
	"github.com/iamdebopriya/Inventory-Management-System-go/internal/usecase"

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
