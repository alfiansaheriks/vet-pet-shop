package main

import (
	"fmt"
	"vet-pet-shop/config"
	"vet-pet-shop/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	routes.AuthRoutes(r)

	fmt.Println("Server started at http://localhost:8080")
	r.Run(":8080")
}
