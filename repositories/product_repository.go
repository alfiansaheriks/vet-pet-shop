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
	err := r.DB.Find(&products).Error
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
	err := r.DB.Where("category = ?", category).Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetProductHighestPrice() ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Order("price desc").Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetProductLowestPrice() ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Order("price asc").Find(&products).Error
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
