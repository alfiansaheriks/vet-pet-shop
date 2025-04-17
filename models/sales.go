package models

import "time"

type SalesTransaction struct {
	ID            uint    `gorm:"primaryKey"`
	BranchID      uint    `gorm:"not null"`     // Cabang tempat transaksi
	CustomerID    *uint   `gorm:"default:null"` // Optional, jika customer terdaftar
	TotalPrice    float64 `gorm:"not null"`
	PaymentMethod string  `gorm:"type:varchar(100)"` // Cash, QRIS, etc
	Status        string  `gorm:"type:varchar(50)"`  // Paid / Unpaid
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Customer *User                  `gorm:"foreignKey:CustomerID;references:ID"`
	Branch   *Branch                `gorm:"foreignKey:BranchID;references:ID"`
	Items    []SalesTransactionItem `gorm:"foreignKey:TransactionID;references:ID;constraint:OnDelete:CASCADE"`
}

type SalesTransactionItem struct {
	ID            uint    `gorm:"primaryKey"`
	TransactionID uint    `gorm:"not null"`
	ProductID     uint    `gorm:"not null"`
	Quantity      uint    `gorm:"not null"`
	Price         float64 `gorm:"not null"` // Harga saat transaksi
	Subtotal      float64 `gorm:"not null"` // Quantity * Price

	Product     *Product          `gorm:"foreignKey:ProductID;references:ID"`
	Transaction *SalesTransaction `gorm:"foreignKey:TransactionID;references:ID"`
}
