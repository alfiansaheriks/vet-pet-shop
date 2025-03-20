package routes

import (
	"vet-pet-shop/controllers"
	"vet-pet-shop/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductRoutes(r *gin.Engine, db *gorm.DB) {
	productGroup := r.Group("/products")

	productGroup.Use(middlewares.AuthMiddleware(db))
	{
		productGroup.GET("/", func(c *gin.Context) { controllers.GetAllProducts(c, db) })
		productGroup.GET("/:id", func(c *gin.Context) { controllers.GetProductById(c, db) })

		productGroup.POST("/", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.CreateProduct(c, db) })
		productGroup.PUT("/:id", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.UpdateProduct(c, db) })
		productGroup.DELETE("/:id", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.DeleteProduct(c, db) })
	}

}
