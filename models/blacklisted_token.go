package models

import "time"

type BlacklistedToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"unique;not null"`
	ExpiredAt time.Time `gorm:"not null"`
}
