package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"vet-pet-shop/models"
	"vet-pet-shop/repositories"
	"vet-pet-shop/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePet(c *gin.Context, db *gorm.DB) {
	var PetRequest models.PetRequest
	if err := c.ShouldBindBodyWithJSON(&PetRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "ID not found"})
		return
	}

	customerID, ok := id.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Invalid ID"})
		return
	}

	var birthDate *time.Time
	if PetRequest.BirthDate != "" {
		parsedDate, err := utils.ParseDatetoTimestamp(PetRequest.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}

		birthDate = &parsedDate
	}

	pet := models.Pet{
		CustomerID: customerID,
		Name:       PetRequest.Name,
		Type:       PetRequest.Type,
		Breed:      PetRequest.Breed,
		Gender:     PetRequest.Gender,
		BirthDate:  birthDate,
		Color:      PetRequest.Color,
		Weight:     PetRequest.Weight,
		CreatedAt:  time.Now(),
	}

	petRepo := repositories.PetRepository{DB: db}
	err := petRepo.CreatePet(&pet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pet"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Pet created successfully"})
}

func GetAllPets(c *gin.Context, db *gorm.DB) {
	petRepo := repositories.PetRepository{DB: db}
	pets, err := petRepo.GetAllPets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pets"})
		return
	}

	petResponse := make([]models.PetResponse, len(pets))
	for i, pet := range pets {
		var birtDate string
		if pet.BirthDate != nil {
			birthDate := utils.ParseTimeStampToDate(*pet.BirthDate)
			birtDate = birthDate
		}

		petResponse[i] = models.PetResponse{
			ID: pet.ID,
			Customer: models.CustomerResponse{
				ID:    pet.Customer.ID,
				Name:  pet.Customer.Name,
				Email: pet.Customer.Email,
			},
			CustomerID: pet.CustomerID,
			Name:       pet.Name,
			Type:       pet.Type,
			Breed:      pet.Breed,
			Gender:     pet.Gender,
			BirthDate:  birtDate,
			Color:      pet.Color,
			Weight:     pet.Weight,
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": petResponse})
}

func GetPetById(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	petRepo := repositories.PetRepository{DB: db}
	pet, err := petRepo.GetPetById(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pet"})
		}
		return
	}

	var birthDate string
	if pet.BirthDate != nil {
		parsed := utils.ParseTimeStampToDate(*pet.BirthDate)
		birthDate = parsed
	}

	petResponse := models.PetResponse{
		ID: pet.ID,
		Customer: models.CustomerResponse{
			ID:    pet.Customer.ID,
			Name:  pet.Customer.Name,
			Email: pet.Customer.Email,
		},
		CustomerID: pet.CustomerID,
		Name:       pet.Name,
		Type:       pet.Type,
		Breed:      pet.Breed,
		Gender:     pet.Gender,
		BirthDate:  birthDate,
		Color:      pet.Color,
		Weight:     pet.Weight,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": petResponse})

}

func UpdatePet(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	var petRequest models.PetRequest
	if err := c.ShouldBindJSON(&petRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	petRepo := repositories.PetRepository{DB: db}
	pet, err := petRepo.GetPetById(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pet"})
		}
		return
	}

	pet.Name = petRequest.Name
	pet.Type = petRequest.Type
	pet.Breed = petRequest.Breed
	pet.Gender = petRequest.Gender
	pet.Color = petRequest.Color
	pet.Weight = petRequest.Weight

	parsedDate, err := utils.ParseDatetoTimestamp(petRequest.BirthDate)
	if err != nil {
		fmt.Println("PARSE ERROR:", err)
	} else {
		fmt.Println("PARSE SUCCESS:", parsedDate.Format(time.RFC3339))
	}

	pet.BirthDate = &parsedDate
	pet.UpdatedAt = time.Now()

	err = petRepo.UpdatePet(pet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pet updated successfully"})
}

func DeletePet(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	petRepo := repositories.PetRepository{DB: db}
	pet, err := petRepo.GetPetById(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pet"})
		}
		return
	}

	err = petRepo.DeletePet(pet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pet deleted successfully"})
}

func GetPetsByCustomerId(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}
	petRepo := repositories.PetRepository{DB: db}
	pets, err := petRepo.GetPetsByCustomerId(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pets"})
		return
	}

	petResponse := make([]models.PetResponse, len(pets))
	for i, pet := range pets {

		var birthDate string
		if pet.BirthDate != nil {
			birthDate := utils.ParseTimeStampToDate(*pet.BirthDate)
			birthDate = birthDate
		}

		petResponse[i] = models.PetResponse{
			ID: pet.ID,
			Customer: models.CustomerResponse{
				ID:    pet.Customer.ID,
				Name:  pet.Customer.Name,
				Email: pet.Customer.Email,
			},
			CustomerID: pet.CustomerID,
			Name:       pet.Name,
			Type:       pet.Type,
			Breed:      pet.Breed,
			Gender:     pet.Gender,
			BirthDate:  birthDate,
			Color:      pet.Color,
			Weight:     pet.Weight,
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": petResponse})
}
