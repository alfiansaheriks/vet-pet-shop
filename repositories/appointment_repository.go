package repositories

import (
	"time"
	"vet-pet-shop/models"

	"gorm.io/gorm"
)

type AppointmentRepository struct {
	DB *gorm.DB
}

func (r *AppointmentRepository) CreateAppointment(appointment *models.Appointment) error {
	return r.DB.Create(appointment).Error
}
func (r *AppointmentRepository) GetAllAppointments() ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Preload("Customer").Preload("Doctor").Preload("Branch").Find(&appointments).Error
	return appointments, err
}
func (r *AppointmentRepository) GetAppointmentById(id uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Preload("Customer").Preload("Doctor").Preload("Branch").Where("id = ?", id).Find(&appointments).Error
	return appointments, err
}
func (r *AppointmentRepository) UpdateAppointment(appointment models.Appointment) error {
	return r.DB.Save(appointment).Error
}
func (r *AppointmentRepository) DeleteAppointment(id uint) error {
	return r.DB.Delete(&models.Appointment{}, id).Error
}
func (r *AppointmentRepository) GetAppointmentByCustomerId(customerId uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Preload("Customer").Preload("Doctor").Preload("Branch").Where("customer_id = ?", customerId).Find(&appointments).Error
	return appointments, err
}
func (r *AppointmentRepository) GetAppointmentByDoctorId(doctorId uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Preload("Customer").Preload("Doctor").Preload("Branch").Where("doctor_id = ?", doctorId).Find(&appointments).Error
	return appointments, err
}
func (r *AppointmentRepository) GetAppointmentByBranchId(branchId uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Preload("Customer").Preload("Doctor").Preload("Branch").Where("branch_id = ?", branchId).Find(&appointments).Error
	return appointments, err
}
func (r *AppointmentRepository) GetAppointmentByDate(date time.Time) ([]models.Appointment, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	var appointments []models.Appointment
	err := r.DB.Preload("Customer").Preload("Doctor").Preload("Branch").
		Where("date >= ? AND date < ?", start, end).
		Find(&appointments).Error

	return appointments, err
}
func (r *AppointmentRepository) GetAppointmentByDateTime(date time.Time, time string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Preload("Customer").Preload("Doctor").Preload("Branch").Where("date = ? AND time = ?", date, time).Find(&appointments).Error
	return appointments, err
}
func (r *AppointmentRepository) GetAppointmentByStatus(status string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Preload("Customer").Preload("Doctor").Preload("Branch").Where("status = ?", status).Find(&appointments).Error
	return appointments, err
}
func (r *AppointmentRepository) GetAppointmentByVisitType(visitType string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Preload("Customer").Preload("Doctor").Preload("Branch").Where("visit_type = ?", visitType).Find(&appointments).Error
	return appointments, err
}
