package repositories

import (
	"time"
	"vet-pet-shop/models"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	DB *gorm.DB
}

func (r *RefreshTokenRepository) SaveRefreshToken(userID uint, token string, expiresAt time.Time) error {
	refreshToken := RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return r.DB.Create(&refreshToken).Error
}

func (r *RefreshTokenRepository) DeleteRefreshToken(token string) error {
	return r.DB.Where("token = ?", token).Delete(&RefreshToken{}).Error
}

func (r *RefreshTokenRepository) DeleteTokensByUserID(userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}

func (r *RefreshTokenRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.DB.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByUserID(userID uint) (RefreshToken, error) {
	var refreshToken RefreshToken
	err := r.DB.Where("user_id = ?", userID).First(&refreshToken).Error
	return refreshToken, err
}
