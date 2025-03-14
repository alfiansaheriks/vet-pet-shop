package routes

import (
	"vet-pet-shop/controllers"
	"vet-pet-shop/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	// r.GET("/users", controllers.GetUsers)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware(db))

	protected.GET("/users", controllers.GetUsers)
	protected.GET("/users/:id", controllers.GetUserByID)
	protected.PUT("/users/:id", controllers.UpdateUser)
	protected.DELETE("/users/:id", controllers.DeleteUser)
	protected.POST("/logout", func(c *gin.Context) {
		controllers.Logout(c, db)
	})

}
