package controllers

import (
	"net/http"
	"os"
	"path/filepath"
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

func CreateProduct(c *gin.Context, db *gorm.DB) {
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var request models.ProductRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	product := models.Product{
		Name:        request.Name,
		Category:    request.Category,
		Description: request.Description,
		Price:       request.Price,
		Unit:        request.Unit,
	}

	productRepository := repositories.ProductRepository{DB: tx}
	if err := productRepository.CreateProduct(&product); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to create product"})
		return
	}

	//handle upload multiple images
	form, err := c.MultipartForm()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to upload image"})
		return
	}
	files := form.File["images"]

	//add product image
	var productImages []models.ProductImages
	for _, file := range files {
		folderPath := filepath.Join("uploads/products", strconv.FormatUint(uint64(product.ID), 10))
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to create directory"})
			return
		}

		//change filename to random
		filename := filepath.Base(file.Filename)
		var randFilename = utils.GenerateRandomString(10) + filepath.Ext(filename)

		filePath := filepath.Join(folderPath, randFilename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to upload image"})
			return
		}

		// save image url to db
		productImages = append(productImages, models.ProductImages{
			ProductID: product.ID,
			ImageURL:  filePath,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		})
	}

	if len(productImages) > 0 {
		for i := range productImages {
			if err := tx.Create(&productImages[i]).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to save image"})
				return
			}
		}
	}

	tx.Commit()

	product.ProductImages = productImages

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": product})
}

func GetAllProducts(c *gin.Context, db *gorm.DB) {
	ProductRepository := repositories.ProductRepository{DB: db}
	products, err := ProductRepository.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})
}

func GetProductById(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "success", "error": err.Error()})
		return
	}

	ProductRepository := repositories.ProductRepository{DB: db}
	products, err := ProductRepository.GetProductById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})
}

func UpdateProduct(c *gin.Context, db *gorm.DB) {
	var request models.ProductRequest
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
		c.JSON(http.StatusBadRequest, gin.H{"status": "success", "error": err.Error()})
		return
	}

	ProductRepository := repositories.ProductRepository{DB: db}
	products, err := ProductRepository.GetProductById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	product := models.Product{
		ID:          products[0].ID,
		Name:        request.Name,
		Category:    request.Category,
		Description: request.Description,
		Price:       request.Price,
		Unit:        request.Unit,
	}

	err = ProductRepository.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": product})
}

func DeleteProduct(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "success", "error": err.Error()})
		return
	}

	ProductRepository := repositories.ProductRepository{DB: db}
	err = ProductRepository.DeleteProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product has been deleted"})
}

func GetProductByCategory(c *gin.Context, db *gorm.DB) {
	category := c.Param("category")

	ProductRepository := repositories.ProductRepository{DB: db}
	products, err := ProductRepository.GetProductByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})
}

func GetProductByPrice(c *gin.Context, db *gorm.DB) {
	orderBy := c.Param("order_by")

	ProductRepository := repositories.ProductRepository{DB: db}
	products, err := ProductRepository.GetProductByPrice(orderBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})
}

func GetProductBySearch(c *gin.Context, db *gorm.DB) {
	search := c.Query("search")

	ProductRepository := repositories.ProductRepository{DB: db}
	products, err := ProductRepository.GetProductBySearch(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})
}

func CreateInventory(c *gin.Context, db *gorm.DB) {
	var request models.InventoryRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	inventory := models.Inventory{
		ProductID:      request.ProductID,
		BranchID:       request.BranchID,
		Stock_Quantity: request.Stock_Quantity,
	}

	ProductRepository := repositories.ProductRepository{DB: db}
	err := ProductRepository.CreateInventory(&inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to create inventory"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": inventory})
}

func UpdateInventory(c *gin.Context, db *gorm.DB) {
	var request models.InventoryRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "success", "error": err.Error()})
		return
	}

	inventoryRepository := repositories.ProductRepository{DB: db}
	inventories, err := inventoryRepository.GetInventory(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	inventory := models.Inventory{
		ID:             inventories[0].ID,
		ProductID:      request.ProductID,
		BranchID:       request.BranchID,
		Stock_Quantity: request.Stock_Quantity,
		UpdatedAt:      time.Now(),
	}

	err = inventoryRepository.UpdateInventory(inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to update inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": inventory})
}

func DeleteInventory(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "success", "error": err.Error()})
		return
	}

	ProductRepository := repositories.ProductRepository{DB: db}
	err = ProductRepository.DeleteInventory(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Inventory has been deleted"})
}

func GetInventories(c *gin.Context, db *gorm.DB) {
	ProductRepository := repositories.ProductRepository{DB: db}
	inventories, err := ProductRepository.GetInventories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to get inventories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": inventories})
}

func GetInventory(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "success", "error": err.Error()})
		return
	}

	ProductRepository := repositories.ProductRepository{DB: db}
	inventories, err := ProductRepository.GetInventory(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to get inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": inventories})
}

func GetInventoryByBranchID(c *gin.Context, db *gorm.DB) {
	branchID, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "success", "error": err.Error()})
		return
	}

	ProductRepository := repositories.ProductRepository{DB: db}
	inventories, err := ProductRepository.GetInventoryByBranchID(uint(branchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to get inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": inventories})
}

func GetInventoryByProductID(c *gin.Context, db *gorm.DB) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "success", "error": err.Error()})
		return
	}

	ProductRepository := repositories.ProductRepository{DB: db}
	inventories, err := ProductRepository.GetInventoryByProductID(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to get inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": inventories})
}
