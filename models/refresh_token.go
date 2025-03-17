package models

import "time"

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"uniqueIndex"`
	UserID    uint   `gorm:"index"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
