package routes

import (
	"vet-pet-shop/controllers"
	"vet-pet-shop/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AppointmentRoutes(r *gin.Engine, db *gorm.DB) {
	appointmentGroup := r.Group("/appointments")

	appointmentGroup.Use(middlewares.AuthMiddleware(db))
	{
		appointmentGroup.POST("/", func(c *gin.Context) { controllers.CreateAppointment(c, db) })
		appointmentGroup.GET("/", func(c *gin.Context) { controllers.GetAllAppointments(c, db) })
		appointmentGroup.GET("/:id", func(c *gin.Context) { controllers.GetAppointmentByID(c, db) })
		appointmentGroup.PUT("/:id", func(c *gin.Context) { controllers.UpdateAppointment(c, db) })
		appointmentGroup.DELETE("/:id", func(c *gin.Context) { controllers.DeleteAppointment(c, db) })
		appointmentGroup.GET("/customer/:id", func(c *gin.Context) { controllers.GetAppointmentByCustomerId(c, db) })
		appointmentGroup.GET("/doctor/:id", func(c *gin.Context) { controllers.GetAppointmentByDoctorId(c, db) })
		appointmentGroup.GET("/branch/:id", func(c *gin.Context) { controllers.GetAppointmentByBranchId(c, db) })
		appointmentGroup.GET("/status/:status", func(c *gin.Context) { controllers.GetAppointmentByStatus(c, db) })
		appointmentGroup.GET("/type/:visit_type", func(c *gin.Context) { controllers.GetAppointmentByVisitType(c, db) })
		appointmentGroup.GET("/date", func(c *gin.Context) { controllers.GetAppointmentByDate(c, db) })
	}
}
