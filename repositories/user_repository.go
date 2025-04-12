package repositories

import (
	"vet-pet-shop/config"
	"vet-pet-shop/models"
)

func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

func CreateUserContact(userContact *models.Contact) error {
	return config.DB.Create(userContact).Error
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.Find(&users).Error
	return users, err
}

func GetUserByID(id uint) (models.User, error) {
	var user models.User
	err := config.DB.Preload("Contact").First(&user, id).Error
	return user, err
}

func UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

func DeleteUser(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}
