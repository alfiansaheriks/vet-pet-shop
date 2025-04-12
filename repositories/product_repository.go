package repositories

import (
	"vet-pet-shop/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	return r.DB.Create(product).Error
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Preload("ProductImages").Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetProductById(id uint) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Where("id = ?", id).Find(&products).Error
	return products, err
}

func (r *ProductRepository) UpdateProduct(product models.Product) error {
	return r.DB.Save(product).Error
}

func (r *ProductRepository) DeleteProduct(id uint) error {
	return r.DB.Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) GetProductByCategory(category string) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Preload("ProductImages").Where("category = ?", category).Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetProductByPrice(orderBy string) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Preload("ProductImages").Order("price " + orderBy).Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetProductBySearch(search string) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Where("name LIKE ?", "%"+search+"%").Find(&products).Error
	return products, err
}

func (r *ProductRepository) CreateProductImage(productImage *models.ProductImages) error {
	return r.DB.Create(productImage).Error
}

func (r *ProductRepository) CreateInventory(inventory *models.Inventory) error {
	return r.DB.Create(inventory).Error
}

func (r *ProductRepository) UpdateInventory(inventory models.Inventory) error {
	return r.DB.Save(inventory).Error
}

func (r *ProductRepository) DeleteInventory(id uint) error {
	return r.DB.Delete(&models.Inventory{}, id).Error
}

func (r *ProductRepository) GetInventories() ([]models.Inventory, error) {
	var inventory []models.Inventory
	err := r.DB.Find(&inventory).Error
	return inventory, err
}

func (r *ProductRepository) GetInventory(id uint) ([]models.Inventory, error) {
	var inventory []models.Inventory
	err := r.DB.Where("id = ?", id).Find(&inventory).Error
	return inventory, err
}

func (r *ProductRepository) GetInventoryByProductID(productID uint) ([]models.Inventory, error) {
	var inventory []models.Inventory
	err := r.DB.Where("product_id = ?", productID).Find(&inventory).Error
	return inventory, err
}

func (r *ProductRepository) GetInventoryByBranchID(branchID uint) ([]models.Inventory, error) {
	var inventory []models.Inventory
	err := r.DB.Where("branch_id = ?", branchID).Find(&inventory).Error
	return inventory, err
}
