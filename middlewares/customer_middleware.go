package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role dari context (sudah di-set di AuthMiddleware)
		userRole, exists := c.Get("role")
		if !exists || userRole != "customer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "Unauthorized: You are not a customer!",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
