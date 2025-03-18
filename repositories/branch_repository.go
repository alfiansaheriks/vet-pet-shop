package repositories

import (
	"vet-pet-shop/models"

	"gorm.io/gorm"
)

type BranchRepository struct {
	DB *gorm.DB
}

func (r *BranchRepository) CreateBranch(branch *models.Branch) error {
	return r.DB.Create(branch).Error
}

func (r *BranchRepository) GetAllBranches() ([]models.Branch, error) {
	var branches []models.Branch
	err := r.DB.Find(&branches).Error
	return branches, err
}

func (r *BranchRepository) GetBranchById(id uint) (*models.Branch, error) {
	var branch models.Branch
	err := r.DB.First(&branch, id).Error
	return &branch, err
}

func (r *BranchRepository) UpdateBranch(branch *models.Branch) error {
	return r.DB.Save(branch).Error
}

func (r *BranchRepository) DeleteBranch(branch *models.Branch) error {
	return r.DB.Delete(branch).Error
}
