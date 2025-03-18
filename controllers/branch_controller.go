package controllers

import (
	"net/http"
	"strconv"
	"vet-pet-shop/models"
	"vet-pet-shop/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateBranch(c *gin.Context, db *gorm.DB) {
	var branch models.Branch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	branchRepository := repositories.BranchRepository{DB: db}
	if err := branchRepository.CreateBranch(&branch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to create branch",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   branch,
	})
}

func GetBranches(c *gin.Context, db *gorm.DB) {
	branchRepository := repositories.BranchRepository{DB: db}
	branches, err := branchRepository.GetAllBranches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to get branches",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   branches,
	})
}

func GetBranchById(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	branchID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid branch ID",
		})
		return
	}
	branchRepository := repositories.BranchRepository{DB: db}
	branch, err := branchRepository.GetBranchById(uint(branchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "Branch not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   branch,
	})
}

func UpdateBranch(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	branchID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid branch ID",
		})
		return
	}
	branchRepository := repositories.BranchRepository{DB: db}
	branch, err := branchRepository.GetBranchById(uint(branchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "Branch not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if err := branchRepository.UpdateBranch(branch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to update branch",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   branch,
	})
}

func DeleteBranch(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	branchRepository := repositories.BranchRepository{DB: db}

	branchID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid branch ID",
		})
		return
	}
	branch, err := branchRepository.GetBranchById(uint(branchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "Branch not found",
		})
		return
	}
	if err := branchRepository.DeleteBranch(branch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to delete branch",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   "Branch deleted successfully",
	})

}
