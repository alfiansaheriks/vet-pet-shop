package controllers

import (
	"net/http"
	"strconv"
	"time"
	"vet-pet-shop/models"
	"vet-pet-shop/repositories"
	"vet-pet-shop/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateMedicalRecord(c *gin.Context, db *gorm.DB) {
	var req models.MedicalRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visitDate, err := utils.ParseDatetoTimestamp(req.VisitDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid visit date format"})
		return
	}

	medicalRecordRepo := repositories.MedicalRecordRepository{DB: db}

	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to start transaction"})
		return
	}

	id, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to get user ID"})
		tx.Rollback()
		return
	}

	customerId := id.(uint)

	medicalRecord := models.MedicalRecord{
		PetID:         req.PetID,
		DoctorID:      req.DoctorID,
		CustomerID:    customerId,
		AppointmentID: req.AppointmentID,
		VisitDate:     visitDate,
		Diagnosis:     req.Diagnosis,
		Treatment:     req.Treatment,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
	}

	if err := medicalRecordRepo.CreateMedicalRecord(tx, &medicalRecord); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to create medical record"})
		return
	}

	var medications []models.MedicalRecordMedication
	for _, med := range req.Medications {
		medications = append(medications, models.MedicalRecordMedication{
			ProductID:       med.ProductID,
			Dosage:          med.Dosage,
			Notes:           med.Notes,
			MedicalRecordID: medicalRecord.ID,
		})
	}

	if err := medicalRecordRepo.CreateMedications(tx, medications); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to create medications"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to commit transaction"})
		return
	}

	db.Preload("Doctor").Preload("Pet").Preload("Customer").Preload("Medications.Product").
		First(&medicalRecord, medicalRecord.ID)

	var medsResp []struct {
		ID          uint   `json:"id"`
		ProductID   uint   `json:"product_id"`
		ProductName string `json:"product_name"`
		Dosage      string `json:"dosage"`
		Notes       string `json:"notes"`
	}

	for _, m := range medicalRecord.Medications {
		medsResp = append(medsResp, struct {
			ID          uint   `json:"id"`
			ProductID   uint   `json:"product_id"`
			ProductName string `json:"product_name"`
			Dosage      string `json:"dosage"`
			Notes       string `json:"notes"`
		}{
			ID:          m.ID,
			ProductID:   m.ProductID,
			ProductName: m.Product.Name,
			Dosage:      m.Dosage,
			Notes:       m.Notes,
		})
	}

	response := models.MedicalRecordResponse{
		ID:            medicalRecord.ID,
		PetID:         medicalRecord.PetID,
		DoctorID:      medicalRecord.DoctorID,
		CustomerID:    medicalRecord.CustomerID,
		AppointmentID: medicalRecord.AppointmentID,
		VisitDate:     medicalRecord.VisitDate.Format("2006-01-02"),
		Diagnosis:     medicalRecord.Diagnosis,
		Treatment:     medicalRecord.Treatment,
		Notes:         medicalRecord.Notes,
		CreatedAt:     medicalRecord.CreatedAt.Format("2006-01-02 15:04:05"),
		Medications:   medsResp,
		Doctor: models.DoctorMiniResponse{
			ID:   medicalRecord.Doctor.ID,
			Name: medicalRecord.Doctor.Name,
		},
		Pet: models.PetMiniResponse{
			ID:   medicalRecord.Pet.ID,
			Name: medicalRecord.Pet.Name,
			Type: medicalRecord.Pet.Type,
		},
		Customer: models.UserMiniResponse{
			ID:   medicalRecord.Customer.ID,
			Name: medicalRecord.Customer.Name,
		},
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Medical record created successfully", "data": response})
}

func GetAllMedicalRecords(c *gin.Context, db *gorm.DB) {
	medicalRecordRepo := repositories.MedicalRecordRepository{DB: db}

	medicalRecords, err := medicalRecordRepo.GetAllMedicalRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch medical records"})
		return
	}

	var response []models.MedicalRecordResponse

	for _, record := range medicalRecords {
		var medsResp []struct {
			ID          uint   `json:"id"`
			ProductID   uint   `json:"product_id"`
			ProductName string `json:"product_name"`
			Dosage      string `json:"dosage"`
			Notes       string `json:"notes"`
		}

		for _, med := range record.Medications {
			medsResp = append(medsResp, struct {
				ID          uint   `json:"id"`
				ProductID   uint   `json:"product_id"`
				ProductName string `json:"product_name"`
				Dosage      string `json:"dosage"`
				Notes       string `json:"notes"`
			}{
				ID:          med.ID,
				ProductID:   med.ProductID,
				ProductName: med.Product.Name,
				Dosage:      med.Dosage,
				Notes:       med.Notes,
			})
		}

		response = append(response, models.MedicalRecordResponse{
			ID:            record.ID,
			PetID:         record.PetID,
			DoctorID:      record.DoctorID,
			CustomerID:    record.CustomerID,
			AppointmentID: record.AppointmentID,
			VisitDate:     utils.ParseTimeStampToDate(record.VisitDate),
			Diagnosis:     record.Diagnosis,
			Treatment:     record.Treatment,
			Notes:         record.Notes,
			CreatedAt:     utils.ParseTimeStampToDate(record.CreatedAt),
			Medications:   medsResp,
			Pet: models.PetMiniResponse{
				ID:   record.Pet.ID,
				Name: record.Pet.Name,
				Type: record.Pet.Type,
			},
			Doctor: models.DoctorMiniResponse{
				ID:   record.Doctor.ID,
				Name: record.Doctor.Name,
			},
			Customer: models.UserMiniResponse{
				ID:   record.Customer.ID,
				Name: record.Customer.Name,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func GetMedicalRecordById(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid medical record ID"})
		return
	}

	medicalRecordRepo := repositories.MedicalRecordRepository{DB: db}

	medicalRecord, err := medicalRecordRepo.GetMedicalRecordById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch medical record"})
		return
	}

	var medsResp []struct {
		ID          uint   `json:"id"`
		ProductID   uint   `json:"product_id"`
		ProductName string `json:"product_name"`
		Dosage      string `json:"dosage"`
		Notes       string `json:"notes"`
	}

	for _, med := range medicalRecord.Medications {
		medsResp = append(medsResp, struct {
			ID          uint   `json:"id"`
			ProductID   uint   `json:"product_id"`
			ProductName string `json:"product_name"`
			Dosage      string `json:"dosage"`
			Notes       string `json:"notes"`
		}{
			ID:          med.ID,
			ProductID:   med.ProductID,
			ProductName: med.Product.Name,
			Dosage:      med.Dosage,
			Notes:       med.Notes,
		})
	}

	response := models.MedicalRecordResponse{
		ID:            medicalRecord.ID,
		PetID:         medicalRecord.PetID,
		DoctorID:      medicalRecord.DoctorID,
		CustomerID:    medicalRecord.CustomerID,
		AppointmentID: medicalRecord.AppointmentID,
		VisitDate:     utils.ParseTimeStampToDate(medicalRecord.VisitDate),
		Diagnosis:     medicalRecord.Diagnosis,
		Treatment:     medicalRecord.Treatment,
		Notes:         medicalRecord.Notes,
		CreatedAt:     utils.ParseTimeStampToDate(medicalRecord.CreatedAt),
		Pet: models.PetMiniResponse{
			ID:   medicalRecord.Pet.ID,
			Name: medicalRecord.Pet.Name,
			Type: medicalRecord.Pet.Type,
		},
		Doctor: models.DoctorMiniResponse{
			ID:   medicalRecord.Doctor.ID,
			Name: medicalRecord.Doctor.Name,
		},
		Customer: models.UserMiniResponse{
			ID:   medicalRecord.Customer.ID,
			Name: medicalRecord.Customer.Name,
		},
		Medications: medsResp,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func UpdateMedicalRecord(c *gin.Context, db *gorm.DB) {
	recordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid medical record ID"})
		return
	}

	var req models.MedicalRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	visitDate, err := utils.ParseDatetoTimestamp(req.VisitDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid visit date format"})
		return
	}

	medicalRecordRepo := repositories.MedicalRecordRepository{DB: db}
	tx := db.Begin()

	record, err := medicalRecordRepo.GetMedicalRecordById(uint(recordID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch medical record"})
		tx.Rollback()
		return
	}

	record.PetID = req.PetID
	record.DoctorID = req.DoctorID
	record.AppointmentID = req.AppointmentID
	record.VisitDate = visitDate
	record.Diagnosis = req.Diagnosis
	record.Treatment = req.Treatment
	record.Notes = req.Notes

	if err := medicalRecordRepo.UpdateMedicalRecord(*record); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to update medical record"})
		return
	}

	existingIDs := map[uint]bool{}
	for _, med := range record.Medications {
		if med.ID != 0 {
			existingIDs[med.ID] = true
		}
	}

	if err := tx.Where("medical_record_id = ? AND id NOT IN ?", record.ID, medicalRecordRepo.GetMedicationsIDs(record.Medications)).Delete(&models.MedicalRecordMedication{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to delete medications"})
		return
	}

	for _, med := range req.Medications {
		if med.ID != 0 {
			if err := tx.Model(&models.MedicalRecordMedication{}).Where("id = ?", med.ID).Updates(models.MedicalRecordMedication{
				ProductID: med.ProductID,
				Dosage:    med.Dosage,
				Notes:     med.Notes,
			}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to update medication"})
				return
			}
		} else {
			newMed := models.MedicalRecordMedication{
				MedicalRecordID: record.ID,
				ProductID:       med.ProductID,
				Dosage:          med.Dosage,
				Notes:           med.Notes,
			}

			if err := medicalRecordRepo.CreateMedications(tx, []models.MedicalRecordMedication{newMed}); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to create medication"})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Medical record updated successfully"})
}

func DeleteMedicalRecord(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid medical record ID"})
		return
	}

	medicalRecordRepo := repositories.MedicalRecordRepository{DB: db}

	if err := medicalRecordRepo.DeleteMedicalRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to delete medical record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Medical record deleted successfully"})
}

func GetMedicalRecordsByPetId(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid pet ID"})
		return
	}

	medicalRecordRepo := repositories.MedicalRecordRepository{DB: db}
	medicalRecords, err := medicalRecordRepo.GetMedicalRecordsByPetId(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch medical records"})
		return
	}

	var response []models.MedicalRecordResponse
	for _, record := range medicalRecords {
		var medsResp []struct {
			ID          uint   `json:"id"`
			ProductID   uint   `json:"product_id"`
			ProductName string `json:"product_name"`
			Dosage      string `json:"dosage"`
			Notes       string `json:"notes"`
		}

		for _, med := range record.Medications {
			medsResp = append(medsResp, struct {
				ID          uint   `json:"id"`
				ProductID   uint   `json:"product_id"`
				ProductName string `json:"product_name"`
				Dosage      string `json:"dosage"`
				Notes       string `json:"notes"`
			}{
				ID:          med.ID,
				ProductID:   med.ProductID,
				ProductName: med.Product.Name,
				Dosage:      med.Dosage,
				Notes:       med.Notes,
			})
		}

		response = append(response, models.MedicalRecordResponse{
			ID:            record.ID,
			PetID:         record.PetID,
			DoctorID:      record.DoctorID,
			CustomerID:    record.CustomerID,
			AppointmentID: record.AppointmentID,
			VisitDate:     utils.ParseTimeStampToDate(record.VisitDate),
			Diagnosis:     record.Diagnosis,
			Treatment:     record.Treatment,
			Notes:         record.Notes,
			CreatedAt:     utils.ParseTimeStampToDate(record.CreatedAt),
			Pet: models.PetMiniResponse{
				ID:   record.Pet.ID,
				Name: record.Pet.Name,
				Type: record.Pet.Type,
			},
			Doctor: models.DoctorMiniResponse{
				ID:   record.Doctor.ID,
				Name: record.Doctor.Name,
			},
			Customer: models.UserMiniResponse{
				ID:   record.Customer.ID,
				Name: record.Customer.Name,
			},
			Medications: medsResp,
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data fetched!", "data": response})
}

func GetMedicalRecordsByCustomerId(c *gin.Context, db *gorm.DB) {
	id, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to get user ID"})
		return
	}

	customerId := id.(uint)

	medicalRecordRepo := repositories.MedicalRecordRepository{DB: db}
	medicalRecords, err := medicalRecordRepo.GetMedicalRecordsByCustomerId(uint(customerId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch medical records"})
		return
	}

	var response []models.MedicalRecordResponse
	for _, record := range medicalRecords {
		var medsResp []struct {
			ID          uint   `json:"id"`
			ProductID   uint   `json:"product_id"`
			ProductName string `json:"product_name"`
			Dosage      string `json:"dosage"`
			Notes       string `json:"notes"`
		}

		for _, med := range record.Medications {
			medsResp = append(medsResp, struct {
				ID          uint   `json:"id"`
				ProductID   uint   `json:"product_id"`
				ProductName string `json:"product_name"`
				Dosage      string `json:"dosage"`
				Notes       string `json:"notes"`
			}{
				ID:          med.ID,
				ProductID:   med.ProductID,
				ProductName: med.Product.Name,
				Dosage:      med.Dosage,
				Notes:       med.Notes,
			})
		}

		response = append(response, models.MedicalRecordResponse{
			ID:            record.ID,
			PetID:         record.PetID,
			DoctorID:      record.DoctorID,
			CustomerID:    record.CustomerID,
			AppointmentID: record.AppointmentID,
			VisitDate:     utils.ParseTimeStampToDate(record.VisitDate),
			Diagnosis:     record.Diagnosis,
			Treatment:     record.Treatment,
			Notes:         record.Notes,
			CreatedAt:     utils.ParseTimeStampToDate(record.CreatedAt),
			Pet: models.PetMiniResponse{
				ID:   record.Pet.ID,
				Name: record.Pet.Name,
				Type: record.Pet.Type,
			},
			Doctor: models.DoctorMiniResponse{
				ID:   record.Doctor.ID,
				Name: record.Doctor.Name,
			},
			Customer: models.UserMiniResponse{
				ID:   record.Customer.ID,
				Name: record.Customer.Name,
			},
			Medications: medsResp,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data fetched!", "data": response})
}

func GetMedicalRecordsByAppointmentId(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid appointment ID"})
		return
	}

	medicalRecordRepo := repositories.MedicalRecordRepository{DB: db}
	medicalRecords, err := medicalRecordRepo.GetMedicalRecordByAppointmentId(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch medical records"})
		return
	}

	var response []models.MedicalRecordResponse
	for _, record := range medicalRecords {
		var medsResp []struct {
			ID          uint   `json:"id"`
			ProductID   uint   `json:"product_id"`
			ProductName string `json:"product_name"`
			Dosage      string `json:"dosage"`
			Notes       string `json:"notes"`
		}

		for _, med := range record.Medications {
			medsResp = append(medsResp, struct {
				ID          uint   `json:"id"`
				ProductID   uint   `json:"product_id"`
				ProductName string `json:"product_name"`
				Dosage      string `json:"dosage"`
				Notes       string `json:"notes"`
			}{
				ID:          med.ID,
				ProductID:   med.ProductID,
				ProductName: med.Product.Name,
				Dosage:      med.Dosage,
				Notes:       med.Notes,
			})
		}

		response = append(response, models.MedicalRecordResponse{
			ID:            record.ID,
			PetID:         record.PetID,
			DoctorID:      record.DoctorID,
			CustomerID:    record.CustomerID,
			AppointmentID: record.AppointmentID,
			VisitDate:     utils.ParseTimeStampToDate(record.VisitDate),
			Diagnosis:     record.Diagnosis,
			Treatment:     record.Treatment,
			Notes:         record.Notes,
			CreatedAt:     utils.ParseTimeStampToDate(record.CreatedAt),
			Pet: models.PetMiniResponse{
				ID:   record.Pet.ID,
				Name: record.Pet.Name,
				Type: record.Pet.Type,
			},
			Doctor: models.DoctorMiniResponse{
				ID:   record.Doctor.ID,
				Name: record.Doctor.Name,
			},
			Customer: models.UserMiniResponse{
				ID:   record.Customer.ID,
				Name: record.Customer.Name,
			},
			Medications: medsResp,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data fetched!", "data": response})
}

func GetMedicalRecordByDoctorId(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid doctor ID"})
		return
	}

	medicalRecordRepo := repositories.MedicalRecordRepository{DB: db}
	medicalRecords, err := medicalRecordRepo.GetMedicalRecordByDoctorId(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch medical records"})
		return
	}

	var response []models.MedicalRecordResponse
	for _, record := range medicalRecords {
		var medsResp []struct {
			ID          uint   `json:"id"`
			ProductID   uint   `json:"product_id"`
			ProductName string `json:"product_name"`
			Dosage      string `json:"dosage"`
			Notes       string `json:"notes"`
		}

		for _, med := range record.Medications {
			medsResp = append(medsResp, struct {
				ID          uint   `json:"id"`
				ProductID   uint   `json:"product_id"`
				ProductName string `json:"product_name"`
				Dosage      string `json:"dosage"`
				Notes       string `json:"notes"`
			}{
				ID:          med.ID,
				ProductID:   med.ProductID,
				ProductName: med.Product.Name,
				Dosage:      med.Dosage,
				Notes:       med.Notes,
			})
		}

		response = append(response, models.MedicalRecordResponse{
			ID:            record.ID,
			PetID:         record.PetID,
			DoctorID:      record.DoctorID,
			CustomerID:    record.CustomerID,
			AppointmentID: record.AppointmentID,
			VisitDate:     utils.ParseTimeStampToDate(record.VisitDate),
			Diagnosis:     record.Diagnosis,
			Treatment:     record.Treatment,
			Notes:         record.Notes,
			CreatedAt:     utils.ParseTimeStampToDate(record.CreatedAt),
			Pet: models.PetMiniResponse{
				ID:   record.Pet.ID,
				Name: record.Pet.Name,
				Type: record.Pet.Type,
			},
			Doctor: models.DoctorMiniResponse{
				ID:   record.Doctor.ID,
				Name: record.Doctor.Name,
			},
			Customer: models.UserMiniResponse{
				ID:   record.Customer.ID,
				Name: record.Customer.Name,
			},
			Medications: medsResp,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data fetched!", "data": response})
}
