package repositories

import (
	"time"

	"gorm.io/gorm"
)

type BlacklistToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"unique;not null"`
	ExpiredAt time.Time `gorm:"not null"`
}

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"uniqueIndex"`
	UserID    uint   `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

func AddBlacklistToken(db *gorm.DB, token string, expiredAt time.Time) error {
	blacklist := BlacklistToken{
		Token:     token,
		ExpiredAt: expiredAt,
	}
	return db.Create(&blacklist).Error
}

func IsTokenBlacklisted(db *gorm.DB, token string) bool {
	var count int64
	db.Model(&BlacklistToken{}).Where("token = ?", token).Count(&count)
	return count > 0
}
