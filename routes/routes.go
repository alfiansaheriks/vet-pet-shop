package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	AuthRoutes(r, db)
	BranchRoutes(r, db)
	ProductRoutes(r, db)
	AppointmentRoutes(r, db)
	PetRoutes(r, db)
}
