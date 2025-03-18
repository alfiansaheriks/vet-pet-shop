package models

import (
	"time"
)

type Branch struct {
	Branch_ID uint    `gorm:"primaryKey" json:"branch_id"`
	Name      string  `gorm:"unique;not null" json:"name"`
	Address   string  `gorm:"not null" json:"address"`
	Phone     string  `gorm:"unique;not null" json:"phone"`
	Latitude  float64 `gorm:"not null" json:"latitude"`
	Longitude float64 `gorm:"not null" json:"longitude"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
