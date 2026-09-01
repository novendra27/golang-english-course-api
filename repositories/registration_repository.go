package repositories

import (
	"errors"

	"english-course-api/models"

	"gorm.io/gorm"
)

// RegistrationRepository interface untuk operasi database pendaftaran
type RegistrationRepository interface {
	CreateWithPayment(registration *models.Registration, payment *models.Payment) error
	FindAll() ([]models.Registration, error)
	FindByID(id uint) (*models.Registration, error)
	FindByStudentID(studentID uint) ([]models.Registration, error)
	FindByCourseID(courseID uint) ([]models.Registration, error)
	FindActiveRegistration(studentID, courseID uint) (*models.Registration, error)
	UpdateStatus(id uint, status string) error
}

type registrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationRepository{db: db}
}

// CreateWithPayment membuat data registrasi dan data payment secara atomik dalam 1 transaksi
func (r *registrationRepository) CreateWithPayment(registration *models.Registration, payment *models.Payment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Buat Registration
		if err := tx.Create(registration).Error; err != nil {
			return err
		}

		// 2. Set RegistrationID pada Payment dan buat Payment
		payment.RegistrationID = registration.ID
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *registrationRepository) FindAll() ([]models.Registration, error) {
	var registrations []models.Registration
	err := r.db.
		Preload("Student").
		Preload("Course").
		Preload("Payment").
		Preload("ClassPlacement.Class").
		Order("id DESC").
		Find(&registrations).Error
	return registrations, err
}

func (r *registrationRepository) FindByID(id uint) (*models.Registration, error) {
	var reg models.Registration
	err := r.db.
		Preload("Student").
		Preload("Course").
		Preload("Payment").
		Preload("ClassPlacement.Class").
		First(&reg, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reg, nil
}

func (r *registrationRepository) FindByStudentID(studentID uint) ([]models.Registration, error) {
	var registrations []models.Registration
	err := r.db.
		Preload("Course").
		Preload("Payment").
		Preload("ClassPlacement.Class").
		Where("student_id = ?", studentID).
		Order("id DESC").
		Find(&registrations).Error
	return registrations, err
}

func (r *registrationRepository) FindByCourseID(courseID uint) ([]models.Registration, error) {
	var registrations []models.Registration
	err := r.db.
		Preload("Student").
		Preload("Payment").
		Preload("ClassPlacement.Class").
		Where("course_id = ?", courseID).
		Order("id DESC").
		Find(&registrations).Error
	return registrations, err
}

func (r *registrationRepository) FindActiveRegistration(studentID, courseID uint) (*models.Registration, error) {
	var reg models.Registration
	err := r.db.
		Where("student_id = ? AND course_id = ? AND status IN ('pending', 'registered')", studentID, courseID).
		First(&reg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reg, nil
}

func (r *registrationRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Registration{}).Where("id = ?", id).Update("status", status).Error
}
