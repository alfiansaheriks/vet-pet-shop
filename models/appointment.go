package models

import "time"

type Appointment struct {
	ID         uint      `gorm:"primaryKey"`
	CustomerID uint      `json:"customer_id" gorm:"not null"`
	DoctorID   uint      `json:"doctor_id" gorm:"not null"`
	BranchID   uint      `json:"branch_id"`
	VisitType  string    `json:"visit_type" gorm:"not null" binding:"required,oneof=visit home"`
	Date       time.Time `json:"date" gorm:"not null" binding:"required"`
	Time       string    `json:"time"`
	Notes      string    `json:"notes" gorm:"not null"`
	Status     string    `json:"status" gorm:"not null" binding:"required,oneof=pending confirmed canceled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Customer   *User     `json:"customer,omitempty" gorm:"foreignKey:CustomerID;references:ID;constraint:OnDelete:CASCADE;"`
	Doctor     *User     `json:"doctor,omitempty" gorm:"foreignKey:DoctorID;references:ID;constraint:OnDelete:CASCADE;"`
	Branch     *Branch   `json:"branch,omitempty" gorm:"foreignKey:BranchID;references:Branch_ID;constraint:OnDelete:CASCADE;"`
}

type AppointmentRequest struct {
	CustomerID uint      `json:"customer_id"`
	DoctorID   uint      `json:"doctor_id" binding:"required"`
	BranchID   uint      `json:"branch_id"`
	VisitType  string    `json:"visit_type" binding:"required,oneof=visit home"`
	Date       time.Time `json:"date" binding:"required"`
	Time       string    `json:"time" binding:"required"`
	Notes      string    `json:"notes" binding:"required"`
}

type UserMiniResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DoctorMiniResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AppointmentResponse struct {
	Customer  UserMiniResponse   `json:"customer"`
	Doctor    DoctorMiniResponse `json:"doctor"`
	Branch    Branch             `json:"branch"`
	ID        uint               `json:"id"`
	Notes     string             `json:"notes"`
	Date      string             `json:"date"`
	Time      string             `json:"time"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	VisitType string             `json:"visit_type"`
}
