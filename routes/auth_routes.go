package routes

import (
	"vet-pet-shop/controllers"
	"vet-pet-shop/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", controllers.Register)
	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c, db)
	})
	r.POST("/refresh-token", func(c *gin.Context) {
		controllers.RefreshTokenHandler(c, db)
	})
	r.POST("/logout", func(c *gin.Context) {
		controllers.Logout(c, db)
	})

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware(db))

	protected.GET("/users", controllers.GetUsers)
	protected.GET("/users/:id", controllers.GetUserByID)
	protected.PUT("/users/:id", controllers.UpdateUser)
	protected.DELETE("/users/:id", controllers.DeleteUser)

}
