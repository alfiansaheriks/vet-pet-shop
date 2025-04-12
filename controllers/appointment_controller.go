package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"vet-pet-shop/config"
	"vet-pet-shop/models"
	"vet-pet-shop/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAppointment(c *gin.Context, db *gorm.DB) {
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var request models.AppointmentRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	id, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "ID not found"})
		return
	}
	customerID, ok := id.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "ID not found"})
		return
	}
	request.CustomerID = customerID

	var doctor models.BranchDoctor
	// doctorRepo := repositories.BranchDoctor{DB: tx}
	if err := tx.Where("user_id = ? AND branch_id = ?", request.DoctorID, request.BranchID).First(&doctor).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Doctor not found in this branch"})
		return
	}
	// check if the doctor is available
	var appointmentCount int64
	if err := tx.Model(&models.Appointment{}).Where("doctor_id = ? AND branch_id = ? AND date = ? AND time = ?", request.DoctorID, request.BranchID, request.Date, request.Time).Count(&appointmentCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if appointmentCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Doctor is not available at this time"})
		return
	}

	appointment := models.Appointment{
		CustomerID: request.CustomerID,
		DoctorID:   request.DoctorID,
		BranchID:   request.BranchID,
		Date:       request.Date,
		Time:       request.Time,
		Notes:      request.Notes,
		Status:     "pending",
		VisitType:  request.VisitType,
		CreatedAt:  time.Now(),
	}

	role, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Role not found"})
		return
	}
	// check if the user is a customer
	if role != "customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "You are not authorized to create an appointment"})
		return
	}

	appointmentRepo := repositories.AppointmentRepository{DB: tx}
	if err := appointmentRepo.CreateAppointment(&appointment); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	tx.Commit()

	if err := db.Preload("Customer").
		Preload("Doctor").
		Preload("Branch").
		First(&appointment, appointment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch appointment with relations"})
		return
	}

	response := models.AppointmentResponse{
		ID: appointment.ID,
		Customer: models.UserMiniResponse{
			ID:    appointment.Customer.ID,
			Name:  appointment.Customer.Name,
			Email: appointment.Customer.Email,
		},
		Doctor: models.DoctorMiniResponse{
			ID:    appointment.Doctor.ID,
			Name:  appointment.Doctor.Name,
			Email: appointment.Doctor.Email,
		},
		Branch:    *appointment.Branch,
		VisitType: appointment.VisitType,
		Date:      appointment.Date.Format("2006-01-02"),
		Time:      appointment.Time,
		Notes:     appointment.Notes,
		Status:    appointment.Status,
		CreatedAt: appointment.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Appointment created successfully", "data": response})
}

func GetAllAppointments(c *gin.Context, db *gorm.DB) {
	apppointmentRepo := repositories.AppointmentRepository{DB: db}
	appointments, err := apppointmentRepo.GetAllAppointments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	var response []models.AppointmentResponse
	for _, appointment := range appointments {
		response = append(response, models.AppointmentResponse{
			ID: appointment.ID,
			Customer: models.UserMiniResponse{
				ID:    appointment.Customer.ID,
				Name:  appointment.Customer.Name,
				Email: appointment.Customer.Email,
			},
			Doctor: models.DoctorMiniResponse{
				ID:    appointment.Doctor.ID,
				Name:  appointment.Doctor.Name,
				Email: appointment.Doctor.Email,
			},
			Branch:    *appointment.Branch,
			VisitType: appointment.VisitType,
			Date:      appointment.Date.Format("2006-01-02"),
			Notes:     appointment.Notes,
			Status:    appointment.Status,
			CreatedAt: appointment.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func GetAppointmentByID(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid appointment ID"})
		return
	}

	appointmentRepo := repositories.AppointmentRepository{DB: db}
	appointment, err := appointmentRepo.GetAppointmentById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if len(appointment) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "Appointment not found"})
		return
	}

	var response models.AppointmentResponse
	response = models.AppointmentResponse{
		ID: appointment[0].ID,
		Customer: models.UserMiniResponse{
			ID:    appointment[0].Customer.ID,
			Name:  appointment[0].Customer.Name,
			Email: appointment[0].Customer.Email,
		},
		Doctor: models.DoctorMiniResponse{
			ID:    appointment[0].Doctor.ID,
			Name:  appointment[0].Doctor.Name,
			Email: appointment[0].Doctor.Email,
		},
		Branch:    *appointment[0].Branch,
		VisitType: appointment[0].VisitType,
		Date:      appointment[0].Date.Format("2006-01-02"),
		Time:      appointment[0].Time,
		Notes:     appointment[0].Notes,
		Status:    appointment[0].Status,
		CreatedAt: appointment[0].CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func UpdateAppointment(c *gin.Context, db *gorm.DB) {
	var request models.AppointmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid appointment ID"})
		return
	}
	appointmentRepo := repositories.AppointmentRepository{DB: db}
	appointment, err := appointmentRepo.GetAppointmentById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	appointments := models.Appointment{
		ID:         appointment[0].ID,
		CustomerID: appointment[0].CustomerID,
		DoctorID:   appointment[0].DoctorID,
		BranchID:   appointment[0].BranchID,
		VisitType:  request.VisitType,
		Date:       request.Date,
		Time:       request.Time,
		Notes:      request.Notes,
		Status:     appointment[0].Status,
		CreatedAt:  appointment[0].CreatedAt,
		UpdatedAt:  time.Now(),
	}

	err = appointmentRepo.UpdateAppointment(appointments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if err := db.Preload("Customer").
		Preload("Doctor").
		Preload("Branch").
		First(&appointments, appointments.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch appointment with relations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Appointment updated successfully"})
}

func DeleteAppointment(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid appointment ID"})
		return
	}

	appointmentRepo := repositories.AppointmentRepository{DB: db}
	err = appointmentRepo.DeleteAppointment(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Appointment deleted successfully"})
}

func GetAppointmentByCustomerId(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid customer ID"})
		return
	}

	appointmentRepo := repositories.AppointmentRepository{DB: db}
	appointments, err := appointmentRepo.GetAppointmentByCustomerId(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if len(appointments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "success", "message": "You don't have any appointment"})
		return
	}

	var response []models.AppointmentResponse
	for _, appointment := range appointments {
		response = append(response, models.AppointmentResponse{
			ID: appointment.ID,
			Customer: models.UserMiniResponse{
				ID:    appointment.Customer.ID,
				Name:  appointment.Customer.Name,
				Email: appointment.Customer.Email,
			},
			Doctor: models.DoctorMiniResponse{
				ID:    appointment.Doctor.ID,
				Name:  appointment.Doctor.Name,
				Email: appointment.Doctor.Email,
			},
			Branch:    *appointment.Branch,
			VisitType: appointment.VisitType,
			Date:      appointment.Date.Format("2006-01-02"),
			Time:      appointment.Time,
			Notes:     appointment.Notes,
			Status:    appointment.Status,
			CreatedAt: appointment.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func GetAppointmentByDoctorId(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid doctor ID"})
		return
	}

	appointmentRepo := repositories.AppointmentRepository{DB: db}
	appointments, err := appointmentRepo.GetAppointmentByDoctorId(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if len(appointments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "success", "message": "You don't have any appointment"})
		return
	}

	var response []models.AppointmentResponse
	for _, appointment := range appointments {
		response = append(response, models.AppointmentResponse{
			ID: appointment.ID,
			Customer: models.UserMiniResponse{
				ID:    appointment.Customer.ID,
				Name:  appointment.Customer.Name,
				Email: appointment.Customer.Email,
			},
			Doctor: models.DoctorMiniResponse{
				ID:    appointment.Doctor.ID,
				Name:  appointment.Doctor.Name,
				Email: appointment.Doctor.Email,
			},
			Branch:    *appointment.Branch,
			VisitType: appointment.VisitType,
			Date:      appointment.Date.Format("2006-01-02"),
			Time:      appointment.Time,
			Notes:     appointment.Notes,
			Status:    appointment.Status,
			CreatedAt: appointment.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func GetAppointmentByBranchId(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid branch ID"})
		return
	}

	appointmentRepo := repositories.AppointmentRepository{DB: db}
	appointments, err := appointmentRepo.GetAppointmentByBranchId(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if len(appointments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "success", "message": "You don't have any appointment"})
		return
	}

	var response []models.AppointmentResponse
	for _, appointment := range appointments {
		response = append(response, models.AppointmentResponse{
			ID: appointment.ID,
			Customer: models.UserMiniResponse{
				ID:    appointment.Customer.ID,
				Name:  appointment.Customer.Name,
				Email: appointment.Customer.Email,
			},
			Doctor: models.DoctorMiniResponse{
				ID:    appointment.Doctor.ID,
				Name:  appointment.Doctor.Name,
				Email: appointment.Doctor.Email,
			},
			Branch:    *appointment.Branch,
			VisitType: appointment.VisitType,
			Date:      appointment.Date.Format("2006-01-02"),
			Time:      appointment.Time,
			Notes:     appointment.Notes,
			Status:    appointment.Status,
			CreatedAt: appointment.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func GetAppointmentByDate(c *gin.Context, db *gorm.DB) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Date is required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid date format"})
		return
	}

	fmt.Printf("date: %v\n", date)

	appointmentRepo := repositories.AppointmentRepository{DB: db}
	appointments, err := appointmentRepo.GetAppointmentByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if len(appointments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "success", "message": "You don't have any appointment"})
		return
	}

	var response []models.AppointmentResponse
	for _, appointment := range appointments {
		response = append(response, models.AppointmentResponse{
			ID: appointment.ID,
			Customer: models.UserMiniResponse{
				ID:    appointment.Customer.ID,
				Name:  appointment.Customer.Name,
				Email: appointment.Customer.Email,
			},
			Doctor: models.DoctorMiniResponse{
				ID:    appointment.Doctor.ID,
				Name:  appointment.Doctor.Name,
				Email: appointment.Doctor.Email,
			},
			Branch:    *appointment.Branch,
			VisitType: appointment.VisitType,
			Date:      appointment.Date.Format("2006-01-02"),
			Time:      appointment.Time,
			Notes:     appointment.Notes,
			Status:    appointment.Status,
			CreatedAt: appointment.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func GetAppointmentByStatus(c *gin.Context, db *gorm.DB) {
	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Status is required"})
		return
	}

	appointmentRepo := repositories.AppointmentRepository{DB: db}
	appointments, err := appointmentRepo.GetAppointmentByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if len(appointments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "success", "message": "You don't have any appointment data."})
		return
	}

	var response []models.AppointmentResponse
	for _, appointment := range appointments {
		response = append(response, models.AppointmentResponse{
			ID: appointment.ID,
			Customer: models.UserMiniResponse{
				ID:    appointment.Customer.ID,
				Name:  appointment.Customer.Name,
				Email: appointment.Customer.Email,
			},
			Doctor: models.DoctorMiniResponse{
				ID:    appointment.Doctor.ID,
				Name:  appointment.Doctor.Name,
				Email: appointment.Doctor.Email,
			},
			Branch:    *appointment.Branch,
			VisitType: appointment.VisitType,
			Date:      appointment.Date.Format("2006-01-02"),
			Time:      appointment.Time,
			Notes:     appointment.Notes,
			Status:    appointment.Status,
			CreatedAt: appointment.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func GetAppointmentByVisitType(c *gin.Context, db *gorm.DB) {
	visitType := c.Param("visit_type")
	if visitType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Visit type is required"})
		return
	}

	appointmentRepo := repositories.AppointmentRepository{DB: db}
	appointments, err := appointmentRepo.GetAppointmentByVisitType(visitType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if len(appointments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "success", "message": "You don't have any appointment"})
		return
	}

	var response []models.AppointmentResponse
	for _, appointment := range appointments {
		response = append(response, models.AppointmentResponse{
			ID: appointment.ID,
			Customer: models.UserMiniResponse{
				ID:    appointment.Customer.ID,
				Name:  appointment.Customer.Name,
				Email: appointment.Customer.Email,
			},
			Doctor: models.DoctorMiniResponse{
				ID:    appointment.Doctor.ID,
				Name:  appointment.Doctor.Name,
				Email: appointment.Doctor.Email,
			},
			Branch:    *appointment.Branch,
			VisitType: appointment.VisitType,
			Date:      appointment.Date.Format("2006-01-02"),
			Time:      appointment.Time,
			Notes:     appointment.Notes,
			Status:    appointment.Status,
			CreatedAt: appointment.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}
