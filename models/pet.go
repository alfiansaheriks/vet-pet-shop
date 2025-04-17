package models

import "time"

type Pet struct {
	ID         uint       `gorm:"primaryKey"`
	CustomerID uint       `json:"customer_id" gorm:"not null"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`   // kucing, anjing, burung, dll.
	Breed      string     `json:"breed"`  // ras
	Gender     string     `json:"gender"` // jantan/betina
	BirthDate  *time.Time `json:"birth_date" gorm:"column:birth_date"`
	Color      string     `json:"color"`
	Weight     float32    `json:"weight"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Customer *User `gorm:"foreignKey:CustomerID;references:ID;constraint:OnDelete:CASCADE"`
}

type PetRequest struct {
	CustomerID uint    `json:"customer_id"`
	Name       string  `json:"name" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Breed      string  `json:"breed" binding:"required"`
	Gender     string  `json:"gender" binding:"required"`
	BirthDate  string  `json:"birth_date"`
	Color      string  `json:"color"`
	Weight     float32 `json:"weight"`
}

type CustomerResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PetResponse struct {
	ID         uint             `json:"id"`
	CustomerID uint             `json:"customer_id"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Breed      string           `json:"breed"`
	Gender     string           `json:"gender"`
	BirthDate  string           `json:"birth_date"`
	Color      string           `json:"color"`
	Weight     float32          `json:"weight"`
	Customer   CustomerResponse `json:"customer"`
}

type PetMiniResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
