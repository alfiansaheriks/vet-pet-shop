package controllers

import (
	"net/http"
	"strconv"
	"time"
	"vet-pet-shop/config"
	"vet-pet-shop/models"
	"vet-pet-shop/repositories"
	"vet-pet-shop/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func CreateBranchDoctor(c *gin.Context, db *gorm.DB) {
	var request models.BranchDoctorRequest
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

	branchDoctor := models.BranchDoctor{
		BranchID: request.BranchID,
		UserID:   request.UserID,
	}

	//check if user is not a doctor
	user, err := repositories.GetUserByID(uint(branchDoctor.UserID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "User not found",
		})
		return
	}

	userRole := user.Role
	if userRole != "doctor" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "User is not a doctor",
		})
		return
	}

	branchDoctorRepository := repositories.BranchDoctor{DB: db}
	err = branchDoctorRepository.CreateBranchDoctor(&branchDoctor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   branchDoctor,
	})
}

func UpdateBranchDoctor(c *gin.Context, db *gorm.DB) {
	var request models.BranchDoctorRequest
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

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid branch doctor ID",
		})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	branchDoctorRepository := repositories.BranchDoctor{DB: db}
	branchDoctor, err := branchDoctorRepository.GetBranchDoctorById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "Branch doctor not found",
		})
		return
	}

	branchDoctor.BranchID = request.BranchID
	branchDoctor.UserID = request.UserID
	branchDoctor.UpdatedAt = time.Now()

	if err := branchDoctorRepository.UpdateBranchDoctor(branchDoctor); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   branchDoctor,
	})

}

func DeleteBranchDoctor(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid branch doctor ID",
		})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	branchDoctorRepository := repositories.BranchDoctor{DB: db}
	err = branchDoctorRepository.DeleteBranchDoctor(uint(id))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to delete branch doctor",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Branch doctor deleted successfully!",
	})
}

func GetBranchDoctors(c *gin.Context, db *gorm.DB) {
	branchRepository := repositories.BranchDoctor{DB: db}
	branches, err := branchRepository.GetAllBranchDoctors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   branches,
	})
}

func GetBranchDoctorById(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid branch doctor ID",
		})
		return
	}

	branchDoctorRepository := repositories.BranchDoctor{DB: db}
	branchDoctor, err := branchDoctorRepository.GetBranchDoctorById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "Branch doctor not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   branchDoctor,
	})
}
