package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role dari context (sudah di-set di AuthMiddleware)
		userRole, exists := c.Get("role")
		fmt.Println(userRole)
		if !exists || userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error":  "Forbidden: You are not an admin!",
			})
			c.Abort()
			return
		}
		c.Next()

	}
}
