package models

import "time"

type MedicalRecord struct {
	ID                  uint       `gorm:"primaryKey"`
	AppointmentID       uint       `json:"appointment_id" gorm:"not null"` // relasi ke appointment
	PetName             string     `json:"pet_name" gorm:"not null"`
	PetType             string     `json:"pet_type"`             // anjing, kucing, dll
	PetBreed            string     `json:"pet_breed"`            // ras/spesies
	PetGender           string     `json:"pet_gender"`           // jantan/betina
	PetAge              string     `json:"pet_age"`              // bisa disimpan dalam format teks: 2 bulan, 1 tahun
	Complaint           string     `json:"complaint"`            // keluhan
	Weight              float32    `json:"weight"`               // berat badan
	Temperature         float32    `json:"temperature"`          // suhu
	PhysicalNotes       string     `json:"physical_notes"`       // deskripsi hasil pemeriksaan fisik
	Diagnosis           string     `json:"diagnosis"`            // diagnosis utama
	AdditionalDiagnosis string     `json:"additional_diagnosis"` // opsional
	Treatment           string     `json:"treatment"`            // tindakan (injeksi, dsb)
	Medication          string     `json:"medication"`           // daftar obat
	Instructions        string     `json:"instructions"`         // petunjuk pemakaian
	NextVisit           *time.Time `json:"next_visit"`           // kalau ada rencana kontrol
	DoctorNote          string     `json:"doctor_note"`          // saran dokter
	CreatedAt           time.Time
	UpdatedAt           time.Time

	Appointment *Appointment `gorm:"foreignKey:AppointmentID;constraint:OnDelete:CASCADE;"`
}
