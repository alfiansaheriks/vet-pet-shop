package repositories

import (
	"vet-pet-shop/models"

	"gorm.io/gorm"
)

type SalesRepository struct {
	DB *gorm.DB
}

func (r *SalesRepository) CreateTransaction(transaction *models.SalesTransaction) error {
	return r.DB.Create(transaction).Error
}

func (r *SalesRepository) GetTransactionByID(id uint) (*models.SalesTransaction, error) {
	var transaction models.SalesTransaction
	if err := r.DB.Preload("Items").Preload("Customer").Preload("Branch").Preload("SalesTransactionItem.Product").First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *SalesRepository) GetAllTransactions() ([]models.SalesTransaction, error) {
	var transactions []models.SalesTransaction
	if err := r.DB.Preload("Items").Preload("Customer").Preload("Branch").Preload("SalesTransactionItem.Product").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *SalesRepository) UpdateTransaction(transaction *models.SalesTransaction) error {
	return r.DB.Save(transaction).Error
}

func (r *SalesRepository) DeleteTransaction(id uint) error {
	return r.DB.Delete(&models.SalesTransaction{}, id).Error
}

func (r *SalesRepository) GetTransactionsByCustomerID(customerID uint) ([]models.SalesTransaction, error) {
	var transactions []models.SalesTransaction
	if err := r.DB.Preload("Items").Preload("Customer").Preload("Branch").Preload("SalesTransactionItem.Product").Where("customer_id = ?", customerID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *SalesRepository) GetTransactionsByBranchID(branchID uint) ([]models.SalesTransaction, error) {
	var transactions []models.SalesTransaction
	if err := r.DB.Preload("Items").Preload("Customer").Preload("Branch").Preload("SalesTransactionItem.Product").Where("branch_id = ?", branchID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
