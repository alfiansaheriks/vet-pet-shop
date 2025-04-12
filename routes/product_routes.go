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
		productGroup.GET("/category/:category", func(c *gin.Context) { controllers.GetProductByCategory(c, db) })
		productGroup.GET("/price/:order_by", func(c *gin.Context) { controllers.GetProductByPrice(c, db) })
		productGroup.GET("/search", func(c *gin.Context) { controllers.GetProductBySearch(c, db) })

		productGroup.POST("/", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.CreateProduct(c, db) })
		productGroup.PUT("/:id", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.UpdateProduct(c, db) })
		productGroup.DELETE("/:id", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.DeleteProduct(c, db) })

		inventoryGroup := productGroup.Group("/inventory")

		inventoryGroup.GET("/", func(c *gin.Context) { controllers.GetInventories(c, db) })
		inventoryGroup.GET("/:id", func(c *gin.Context) { controllers.GetInventory(c, db) })
		inventoryGroup.GET("/branch/:branch_id", func(c *gin.Context) { controllers.GetInventoryByBranchID(c, db) })
		inventoryGroup.GET("/product/:product_id", func(c *gin.Context) { controllers.GetInventoryByProductID(c, db) })

		inventoryGroup.Use(middlewares.AdminMiddleware())
		inventoryGroup.POST("/", func(c *gin.Context) { controllers.CreateInventory(c, db) })
		inventoryGroup.PUT("/:id", func(c *gin.Context) { controllers.UpdateInventory(c, db) })
		inventoryGroup.DELETE("/:id", func(c *gin.Context) { controllers.DeleteInventory(c, db) })

	}

}
