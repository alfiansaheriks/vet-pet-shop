package models

import "time"

type MedicalRecord struct {
	ID            uint      `gorm:"primaryKey"`
	PetID         uint      `gorm:"not null"`
	DoctorID      uint      `gorm:"not null"`
	CustomerID    uint      `gorm:"not null"`
	AppointmentID *uint     `gorm:"default:null"` // Optional
	VisitDate     time.Time `gorm:"not null"`
	Diagnosis     string    `gorm:"type:text"`
	Treatment     string    `gorm:"type:text"`
	Notes         string    `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Pet         *Pet                      `gorm:"foreignKey:PetID;references:ID;constraint:OnDelete:CASCADE;"`
	Doctor      *User                     `gorm:"foreignKey:DoctorID;references:ID;constraint:OnDelete:CASCADE;"`
	Customer    *User                     `gorm:"foreignKey:CustomerID;references:ID;constraint:OnDelete:CASCADE;"`
	Appointment *Appointment              `gorm:"foreignKey:AppointmentID;references:ID;constraint:OnDelete:CASCADE;"`
	Medications []MedicalRecordMedication `gorm:"foreignKey:MedicalRecordID;references:ID;constraint:OnDelete:CASCADE;"`
}

type MedicalRecordMedication struct {
	ID              uint   `gorm:"primaryKey"`
	MedicalRecordID uint   `gorm:"not null"`
	ProductID       uint   `gorm:"not null"` // FK ke tabel produk
	Dosage          string `gorm:"type:varchar(255)"`
	Notes           string `gorm:"type:text"` // opsional catatan tambahan

	// Relationships
	Product       *Product       `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE;"`
	MedicalRecord *MedicalRecord `gorm:"foreignKey:MedicalRecordID;references:ID;constraint:OnDelete:CASCADE;"`
}

type MedicalRecordRequest struct {
	PetID         uint   `json:"pet_id" binding:"required"`
	DoctorID      uint   `json:"doctor_id" binding:"required"`
	AppointmentID *uint  `json:"appointment_id"`
	VisitDate     string `json:"visit_date" binding:"required"`
	Diagnosis     string `json:"diagnosis"`
	Treatment     string `json:"treatment"`
	Notes         string `json:"notes"`
	Medications   []struct {
		ID        uint   `json:"id"`
		ProductID uint   `json:"product_id" binding:"required"`
		Dosage    string `json:"dosage" binding:"required"`
		Notes     string `json:"notes"`
	} `json:"medications"`
}

type MedicalRecordResponse struct {
	ID            uint   `json:"id"`
	PetID         uint   `json:"pet_id"`
	DoctorID      uint   `json:"doctor_id"`
	CustomerID    uint   `json:"customer_id"`
	AppointmentID *uint  `json:"appointment_id"`
	VisitDate     string `json:"visit_date"`
	Diagnosis     string `json:"diagnosis"`
	Treatment     string `json:"treatment"`
	Notes         string `json:"notes"`
	CreatedAt     string `json:"created_at"`
	Medications   []struct {
		ID          uint   `json:"id"`
		ProductID   uint   `json:"product_id"`
		ProductName string `json:"product_name"`
		Dosage      string `json:"dosage"`
		Notes       string `json:"notes"`
	} `json:"medications"`

	Pet      PetMiniResponse    `json:"pet"`
	Doctor   DoctorMiniResponse `json:"doctor"`
	Customer UserMiniResponse   `json:"customer"`
}

type Medications struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	Product   string `json:"product"`
	Dosage    string `json:"dosage"`
	Notes     string `json:"notes"`
}
