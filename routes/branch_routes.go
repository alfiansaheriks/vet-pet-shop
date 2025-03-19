package routes

import (
	"vet-pet-shop/controllers"
	"vet-pet-shop/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BranchRoutes(r *gin.Engine, db *gorm.DB) {
	branchGroup := r.Group("/branches")

	branchGroup.Use(middlewares.AuthMiddleware(db))
	{
		branchGroup.GET("/", func(c *gin.Context) { controllers.GetBranches(c, db) })
		branchGroup.GET("/:id", func(c *gin.Context) { controllers.GetBranchById(c, db) })

		branchGroup.POST("/", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.CreateBranch(c, db) })
		branchGroup.PUT("/:id", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.UpdateBranch(c, db) })
		branchGroup.DELETE("/:id", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.DeleteBranch(c, db) })

		branchGroup.GET("/doctors", func(c *gin.Context) { controllers.GetBranchDoctors(c, db) })
		branchGroup.GET("/doctors/:id", func(c *gin.Context) { controllers.GetBranchDoctorById(c, db) })

		branchGroup.POST("/doctors", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.CreateBranchDoctor(c, db) })
		branchGroup.PUT("/doctors/:id", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.UpdateBranchDoctor(c, db) })
		branchGroup.DELETE("/doctors/:id", middlewares.AdminMiddleware(), func(c *gin.Context) { controllers.DeleteBranchDoctor(c, db) })
	}
}
