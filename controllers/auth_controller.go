package controllers

import (
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

func Login(c *gin.Context) {
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

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully!",
		"token":   token,
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"role":      user.Role,
			"token_exp": time.Now().Add(24 * time.Hour),
		},
	})
}

func Logout(c *gin.Context, db *gorm.DB) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Missing token",
		})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Invalid token format",
		})
	}

	token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTKey(), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Invalid token",
		})
		return
	}

	claims, ok := token.Claims.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Invalid token",
		})
		return
	}

	blacklistedToken := models.BlacklistedToken{
		Token:     tokenString,
		ExpiredAt: claims.ExpiresAt.Time,
	}

	if err := db.Create(&blacklistedToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to blacklist token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User logged out successfully!",
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
