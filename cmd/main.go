package main

import (
	"inventory-service-go/database"
	"inventory-service-go/internal/delivery/http"
	"inventory-service-go/internal/repository"
	"inventory-service-go/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.ConnectDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	r := gin.Default()

	repo := repository.NewRepository()
	uc := usecase.NewUsecase(repo)
	handler := http.NewHandler(uc)

	http.SetupRoutes(r, handler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
