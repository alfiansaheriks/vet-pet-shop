package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" gorm:"unique;not null" binding:"required,email"`
	Password  string    `json:"password,omitempty" binding:"required"`
	Role      string    `json:"role" gorm:"type:varchar(20);not null" binding:"required,oneof=admin doctor customer"`
	CreatedAt time.Time `json:"created_at"`
}
