package repositories

import (
	"vet-pet-shop/models"

	"gorm.io/gorm"
)

type MedicalRecordRepository struct {
	DB *gorm.DB
}

func (r *MedicalRecordRepository) CreateMedicalRecord(tx *gorm.DB, medicalRecord *models.MedicalRecord) error {
	return tx.Create(medicalRecord).Error
}

func (r *MedicalRecordRepository) CreateMedications(tx *gorm.DB, medications []models.MedicalRecordMedication) error {
	return tx.Create(&medications).Error
}

func (r *MedicalRecordRepository) GetMedicationsIDs(meds []models.MedicalRecordMedication) []uint {
	ids := make([]uint, len(meds))
	for i, med := range meds {
		ids[i] = med.ProductID
	}
	return ids
}

func (r *MedicalRecordRepository) GetAllMedicalRecords() ([]models.MedicalRecord, error) {
	var records []models.MedicalRecord
	err := r.DB.
		Preload("Pet").
		Preload("Doctor").
		Preload("Customer").
		Preload("Medications.Product").
		Find(&records).Error

	return records, err
}

func (r *MedicalRecordRepository) GetMedicalRecordById(id uint) (*models.MedicalRecord, error) {
	var medicalRecord models.MedicalRecord
	err := r.DB.Preload("Pet").Preload("Doctor").Preload("Customer").Preload("Medications.Product").Where("id = ?", id).First(&medicalRecord).Error
	if err != nil {
		return nil, err
	}
	return &medicalRecord, nil
}

func (r *MedicalRecordRepository) UpdateMedicalRecord(medicalRecord models.MedicalRecord) error {
	return r.DB.Save(&medicalRecord).Error
}

func (r *MedicalRecordRepository) DeleteMedicalRecord(id uint) error {
	// Hapus data terkait di tabel medical_record_medications
	if err := r.DB.Where("medical_record_id = ?", id).Delete(&models.MedicalRecordMedication{}).Error; err != nil {
		return err
	}

	// hapus data di tabel medical_records
	if err := r.DB.Where("id = ?", id).Delete(&models.MedicalRecord{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *MedicalRecordRepository) GetMedicalRecordsByPetId(petId uint) ([]models.MedicalRecord, error) {
	var medicalRecords []models.MedicalRecord
	err := r.DB.Preload("Pet").Preload("Doctor").Preload("Customer").Preload("Medications.Product").Where("pet_id = ?", petId).Find(&medicalRecords).Error
	return medicalRecords, err
}

func (r *MedicalRecordRepository) GetMedicalRecordsByCustomerId(customerId uint) ([]models.MedicalRecord, error) {
	var medicalRecords []models.MedicalRecord
	err := r.DB.Preload("Pet").Preload("Doctor").Preload("Customer").Preload("Medications.Product").Where("customer_id = ?", customerId).Find(&medicalRecords).Error
	return medicalRecords, err
}

func (r *MedicalRecordRepository) GetMedicalRecordByAppointmentId(appointmentId uint) ([]models.MedicalRecord, error) {
	var medicalRecord []models.MedicalRecord
	err := r.DB.Preload("Pet").Preload("Doctor").Preload("Customer").Preload("Medications.Product").Where("appointment_id = ?", appointmentId).First(&medicalRecord).Error
	if err != nil {
		return nil, err
	}
	return medicalRecord, nil
}

func (r *MedicalRecordRepository) GetMedicalRecordByDoctorId(doctorId uint) ([]models.MedicalRecord, error) {
	var medicalRecords []models.MedicalRecord
	err := r.DB.Preload("Pet").Preload("Doctor").Preload("Customer").Preload("Medications.Product").Where("doctor_id = ?", doctorId).Find(&medicalRecords).Error
	return medicalRecords, err
}
