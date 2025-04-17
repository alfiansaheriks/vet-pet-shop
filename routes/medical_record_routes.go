package routes

import (
	"vet-pet-shop/controllers"
	"vet-pet-shop/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MedicalRecordRoutes(r *gin.Engine, db *gorm.DB) {
	medicalRecordGroup := r.Group("/medical-records")

	medicalRecordGroup.Use(middlewares.AuthMiddleware(db))
	{
		medicalRecordGroup.POST("/", func(c *gin.Context) { controllers.CreateMedicalRecord(c, db) })
		medicalRecordGroup.GET("/", func(c *gin.Context) { controllers.GetAllMedicalRecords(c, db) })
		medicalRecordGroup.GET("/:id", func(c *gin.Context) { controllers.GetMedicalRecordById(c, db) })
		medicalRecordGroup.PUT("/:id", func(c *gin.Context) { controllers.UpdateMedicalRecord(c, db) })
		medicalRecordGroup.DELETE("/:id", func(c *gin.Context) { controllers.DeleteMedicalRecord(c, db) })
		medicalRecordGroup.GET("/pet/:id", func(c *gin.Context) { controllers.GetMedicalRecordsByPetId(c, db) })
		medicalRecordGroup.GET("/customer", func(c *gin.Context) { controllers.GetMedicalRecordsByCustomerId(c, db) })
		medicalRecordGroup.GET("/doctor", middlewares.DoctorMiddleware(), func(c *gin.Context) { controllers.GetMedicalRecordByDoctorId(c, db) })
		medicalRecordGroup.GET("/appointment/:id", func(c *gin.Context) { controllers.GetMedicalRecordsByAppointmentId(c, db) })
	}
}
