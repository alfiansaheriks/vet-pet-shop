package repositories

import (
	"vet-pet-shop/models"

	"gorm.io/gorm"
)

type BranchDoctor struct {
	DB *gorm.DB
}

func (r *BranchDoctor) CreateBranchDoctor(branchDoctor *models.BranchDoctor) error {
	return r.DB.Create(branchDoctor).Error
}

func (r *BranchDoctor) GetAllBranchDoctors() ([]models.BranchDoctor, error) {
	var branchDoctors []models.BranchDoctor
	err := r.DB.Find(&branchDoctors).Error
	return branchDoctors, err
}

func (r *BranchDoctor) GetDoctorByBranchId(id uint) ([]models.BranchDoctor, error) {
	var branchDoctors []models.BranchDoctor
	err := r.DB.Where("branch_id = ?", id).Find(&branchDoctors).Error
	return branchDoctors, err
}

func (r *BranchDoctor) GetBranchDoctorById(id uint) (*models.BranchDoctor, error) {
	var branchDoctor models.BranchDoctor
	err := r.DB.First(&branchDoctor, id).Error
	return &branchDoctor, err
}

func (r *BranchDoctor) UpdateBranchDoctor(branchDoctor *models.BranchDoctor) error {
	return r.DB.Save(branchDoctor).Error
}

func (r *BranchDoctor) DeleteBranchDoctor(id uint) error {
	return r.DB.Delete(&models.BranchDoctor{}, id).Error
}
