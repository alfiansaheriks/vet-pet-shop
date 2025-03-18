package middlewares

import (
	"net/http"
	"strings"
	"vet-pet-shop/models"
	"vet-pet-shop/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "Token not found",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token salah"})
			c.Abort()
			return
		}

		var blacklistedToken models.BlacklistedToken
		if err := db.Where("token = ?", tokenString).First(&blacklistedToken).Error; err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token sudah tidak valid"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return utils.GetJWTKey(), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid",
				"token": tokenString,
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*utils.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "Token tidak valid",
			})
			c.Abort()
			return
		}

		c.Set("id", claims.ID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
