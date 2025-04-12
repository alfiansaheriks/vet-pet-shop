package models

import (
	"time"
)

type Branch struct {
	Branch_ID uint        `gorm:"primaryKey"`
	Name      string      `gorm:"unique;not null" json:"name"`
	Address   string      `gorm:"not null" json:"address"`
	Phone     string      `gorm:"unique;not null" json:"phone"`
	Latitude  float64     `gorm:"not null" json:"latitude"`
	Longitude float64     `gorm:"not null" json:"longitude"`
	CreatedAt time.Time   `json:"-"`
	UpdatedAt time.Time   `json:"-"`
	Inventory []Inventory `json:"inventory,omitempty" gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;"`
}

type BranchDoctor struct {
	BranchDoctor_ID uint    `gorm:"primaryKey" json:"branch_doctor_id"`
	BranchID        uint    `json:"branch_id"`
	Branch          *Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;"`
	UserID          uint    `json:"user_id"`
	User            *User   `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type BranchDoctorRequest struct {
	BranchID uint `json:"branch_id" binding:"required"`
	UserID   uint `json:"user_id" binding:"required"`
	// UpdatedAt time.Time
}
