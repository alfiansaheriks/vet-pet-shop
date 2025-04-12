package routes

import (
	"vet-pet-shop/controllers"
	"vet-pet-shop/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PetRoutes(r *gin.Engine, db *gorm.DB) {
	petGroup := r.Group("/pets")

	petGroup.Use(middlewares.AuthMiddleware(db))
	{
		petGroup.GET("/", func(c *gin.Context) { controllers.GetAllPets(c, db) })
		petGroup.GET("/:id", func(c *gin.Context) { controllers.GetPetById(c, db) })
		petGroup.GET("/customer/:id", func(c *gin.Context) { controllers.GetPetsByCustomerId(c, db) })

		petGroup.POST("/", middlewares.CustomerMiddleware(), func(c *gin.Context) { controllers.CreatePet(c, db) })
		petGroup.PUT("/:id", middlewares.CustomerMiddleware(), func(c *gin.Context) { controllers.UpdatePet(c, db) })
		petGroup.DELETE("/:id", middlewares.CustomerMiddleware(), func(c *gin.Context) { controllers.DeletePet(c, db) })
	}
}
