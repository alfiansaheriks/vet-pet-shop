package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DoctorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role dari context (sudah di-set di AuthMiddleware)
		userRole, exists := c.Get("role")
		if !exists || userRole != "doctor" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "Unauthorized: You are not a doctor!",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
