package repositories

import (
	"errors"
	"time"

	"english-course-api/models"

	"gorm.io/gorm"
)

// PaymentRepository interface untuk operasi database pembayaran
type PaymentRepository interface {
	FindAll() ([]models.Payment, error)
	FindByID(id uint) (*models.Payment, error)
	FindByRegistrationID(registrationID uint) (*models.Payment, error)
	ProcessPaymentSuccess(paymentID uint, method string, paidAt time.Time) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) FindAll() ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.
		Preload("Registration.Student").
		Preload("Registration.Course").
		Order("id DESC").
		Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) FindByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.
		Preload("Registration.Student").
		Preload("Registration.Course").
		First(&payment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByRegistrationID(registrationID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.
		Where("registration_id = ?", registrationID).
		First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// ProcessPaymentSuccess mengubah status Payment menjadi 'paid' dan Registration menjadi 'registered' dalam 1 database transaction
func (r *paymentRepository) ProcessPaymentSuccess(paymentID uint, method string, paidAt time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Ambil data payment
		var payment models.Payment
		if err := tx.First(&payment, paymentID).Error; err != nil {
			return err
		}

		// 2. Update status Payment
		payment.Status = "paid"
		payment.PaymentMethod = method
		payment.PaymentDate = &paidAt
		if err := tx.Save(&payment).Error; err != nil {
			return err
		}

		// 3. Update status Registration menjadi 'registered'
		if err := tx.Model(&models.Registration{}).
			Where("id = ?", payment.RegistrationID).
			Update("status", "registered").Error; err != nil {
			return err
		}

		return nil
	})
}
