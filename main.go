package main

import (
	"fmt"
	"vet-pet-shop/config"
	"vet-pet-shop/models"
	"vet-pet-shop/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()

	r := gin.Default()

	db.AutoMigrate(&models.Branch{})

	routes.AuthRoutes(r, db)
	routes.BranchRoutes(r, db)

	fmt.Println("Server started at http://localhost:8080")
	r.Run(":8080")
}
