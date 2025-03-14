package models

import "time"

type User struct {
	ID              uint      `gorm:"primaryKey"`
	Name            string    `json:"name" binding:"required"`
	Email           string    `json:"email" gorm:"unique;not null" binding:"required,email"`
	Password        string    `json:"-" gorm:"column:password"`
	Role            string    `json:"role" gorm:"type:varchar(20);not null" binding:"required,oneof=admin doctor customer"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Phone_Number    string    `json:"phone_number" gorm:"unique;not null" binding:"required"`
	Wa_Phone_Number string    `json:"wa_phone_number" gorm:"unique;not null" binding:"required"`
}

type UserRegistrationRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" gorm:"unique;not null" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	Role            string `json:"role" gorm:"type:varchar(20);not null" binding:"required,oneof=admin doctor customer"`
	Phone_Number    string `json:"phone_number" gorm:"unique;not null" binding:"required"`
	Wa_Phone_Number string `json:"wa_phone_number" gorm:"unique;not null" binding:"required"`
}

type UserEditRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" gorm:"unique;not null" binding:"email"`
	// Password        string `json:"password" binding:"min=6"`
	Role string `json:"role" gorm:"type:varchar(20);not null" binding:"oneof=admin doctor customer"`
	// Phone_Number    string `json:"phone_number" gorm:"unique;not null"`
	// Wa_Phone_Number string `json:"wa_phone_number" gorm:"unique;not null"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
