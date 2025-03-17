package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vet-pet-shop/models"
	"vet-pet-shop/repositories"
	"vet-pet-shop/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	// var user models.User
	var request models.UserRegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessage := utils.FormatValidationErrors(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  errorMessage,
			})
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to hash password",
		})
		return
	}

	user := models.User{
		Name:            request.Name,
		Email:           request.Email,
		Password:        string(hashedPassword),
		Role:            request.Role,
		Phone_Number:    request.Phone_Number,
		Wa_Phone_Number: request.Wa_Phone_Number,
		CreatedAt:       time.Now(),
	}

	if err := repositories.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully!",
		"data":    user,
	})
}

func Login(c *gin.Context, db *gorm.DB) {
	var input models.UserLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessage := utils.FormatValidationErrors(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  errorMessage,
			})
			return
		}
	}

	user, err := repositories.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Email or password is incorrect!",
		})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Email or password is incorrect!",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to generate token",
		})
		return
	}

	refreshTokenRepo := repositories.RefreshTokenRepository{DB: db}
	//delete all refresh token by user id
	_ = refreshTokenRepo.DeleteTokensByUserID(user.ID)

	refreshToken, _ := utils.GenerateRefreshToken(user.ID)
	_ = refreshTokenRepo.SaveRefreshToken(user.ID, refreshToken, time.Now().Add(7*24*time.Hour))

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":       "User logged in successfully!",
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

func Logout(c *gin.Context, db *gorm.DB) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid request",
		})
	}

	refreshTokenRepo := repositories.RefreshTokenRepository{DB: db}
	refreshToken, err := refreshTokenRepo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":        "error",
			"error":         "Failed to get refresh token",
			"refresh_token": req.RefreshToken,
		})
		return
	}

	log.Printf("deleting refresh token: %v", refreshToken)

	err = refreshTokenRepo.DeleteRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to delete refresh token",
		})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString != authHeader {
			token, _ := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return utils.GetJWTKey(), nil
			})

			if claims, ok := token.Claims.(*utils.Claims); ok {
				blacklistedToken := models.BlacklistedToken{
					Token:     tokenString,
					ExpiredAt: claims.ExpiresAt.Time,
				}
				db.Create(&blacklistedToken)
			}
		}
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User logged out successfully!",
	})
}

func RefreshTokenHandler(c *gin.Context, db *gorm.DB) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid request",
		})
		return
	}

	refreshTokenRepo := repositories.RefreshTokenRepository{DB: db}
	refreshToken, err := refreshTokenRepo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to get refresh token",
		})
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Refresh token expired",
		})
		return
	}

	user, err := repositories.GetUserByID(refreshToken.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to get user",
		})
		return
	}

	// Generate access token baru
	newAccessToken, _ := utils.GenerateJWT(user.ID, user.Email, user.Role)

	// Hapus refresh token lama dan re create
	_ = refreshTokenRepo.DeleteRefreshToken(req.RefreshToken)
	newRefreshToken, _ := utils.GenerateRefreshToken(user.ID)
	_ = refreshTokenRepo.SaveRefreshToken(user.ID, newRefreshToken, time.Now().Add(7*24*time.Hour))

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})

}

func GetUsers(c *gin.Context) {
	users, err := repositories.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to fetch users",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   users,
	})
}

func GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid user ID",
		})
		return
	}

	user, err := repositories.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid user ID",
		})
		return
	}

	var input models.UserEditRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessage := utils.FormatValidationErrors(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  errorMessage,
			})
			return
		}
	}

	user, err := repositories.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "User not found",
		})
		return
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Role = input.Role
	user.UpdatedAt = time.Now()

	if err := repositories.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully!",
		"data":    user,
	})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid user ID",
		})
		return
	}

	if err := repositories.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully!",
	})
}
