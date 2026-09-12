// Package repositories provides data access layer abstraction and database operations via GORM.
package repositories

import (
	"errors"
	"time"

	"english-course-api/models"

	"gorm.io/gorm"
)

// PaymentRepository defines the data access contract for Payment entities.
type PaymentRepository interface {
	FindAll() ([]models.Payment, error)
	FindByID(id uint) (*models.Payment, error)
	FindByRegistrationID(registrationID uint) (*models.Payment, error)
	ProcessPaymentSuccess(paymentID uint, method string, paidAt time.Time) error
}

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new PaymentRepository instance.
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

// ProcessPaymentSuccess transitions payment status to 'paid' and registration status to 'registered' atomically in one transaction.
func (r *paymentRepository) ProcessPaymentSuccess(paymentID uint, method string, paidAt time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch payment record
		var payment models.Payment
		if err := tx.First(&payment, paymentID).Error; err != nil {
			return err
		}

		// 2. Update payment status
		payment.Status = "paid"
		payment.PaymentMethod = method
		payment.PaymentDate = &paidAt
		if err := tx.Save(&payment).Error; err != nil {
			return err
		}

		// 3. Update associated registration status to 'registered'
		if err := tx.Model(&models.Registration{}).
			Where("id = ?", payment.RegistrationID).
			Update("status", "registered").Error; err != nil {
			return err
		}

		return nil
	})
}

