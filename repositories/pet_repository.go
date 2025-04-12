package repositories

import (
	"vet-pet-shop/models"

	"gorm.io/gorm"
)

type PetRepository struct {
	DB *gorm.DB
}

func (r *PetRepository) CreatePet(pet *models.Pet) error {
	return r.DB.Create(pet).Error
}

func (r *PetRepository) GetAllPets() ([]models.Pet, error) {
	var pets []models.Pet
	err := r.DB.Preload("Customer").Find(&pets).Error
	return pets, err
}

func (r *PetRepository) GetPetById(id uint) (*models.Pet, error) {
	var pet models.Pet
	err := r.DB.Preload("Customer").First(&pet, id).Error
	return &pet, err
}

func (r *PetRepository) UpdatePet(pet *models.Pet) error {
	return r.DB.Debug().Save(pet).Error
}

func (r *PetRepository) DeletePet(pet *models.Pet) error {
	return r.DB.Delete(pet).Error
}

func (r *PetRepository) GetPetsByCustomerId(customerId uint) ([]models.Pet, error) {
	var pets []models.Pet
	err := r.DB.Preload("Customer").Where("customer_id = ?", customerId).Find(&pets).Error
	return pets, err
}
