package models

import (
	"mime/multipart"
	"time"
)

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" binding:"required"`
	Category    string    `json:"category" binding:"required" gorm:"type:product_category;not null"`
	Description string    `json:"description" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	Unit        string    `json:"unit" binding:"required"` // e.g. kg, pcs, etc.
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time
	// Inventory     []Inventory     `json:"inventory,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	ProductImages []ProductImages `json:"product_images,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

type Inventory struct {
	ID             uint      `gorm:"primaryKey"`
	ProductID      uint      `json:"product_id"`
	Product        *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	BranchID       uint      `json:"branch_id"`
	Branch         *Branch   `json:"branch,omitempty" gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;"`
	Stock_Quantity int       `json:"stock" binding:"required"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time
}

type ProductImages struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID uint      `json:"product_id"`
	Product   *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	ImageURL  string    `json:"image_url" binding:"required"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type ProductDiscount struct {
	ID                 uint      `gorm:"primaryKey"`
	ProductID          uint      `json:"product_id"`
	Product            *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	DiscountPercentage float64   `json:"discount_percentage" binding:"required"`
	StartDate          time.Time `json:"start_date" binding:"required"`
	EndDate            time.Time `json:"end_date" binding:"required"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ProductRequest struct {
	Name        string                  `form:"name" binding:"required"`
	Category    string                  `form:"category" binding:"required"`
	Description string                  `form:"description" binding:"required"`
	Price       float64                 `form:"price" binding:"required"`
	Unit        string                  `form:"unit" binding:"required"`
	Images      []*multipart.FileHeader `form:"images"`
}

type InventoryRequest struct {
	ProductID      uint `form:"product_id" binding:"required"`
	BranchID       uint `form:"branch_id" binding:"required"`
	Stock_Quantity int  `form:"stock" binding:"required"`
}

type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Unit        string    `json:"unit"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	ProductImages []struct {
		ID       uint   `json:"id"`
		ImageURL string `json:"image_url"`
	} `json:"product_images,omitempty"`
}
