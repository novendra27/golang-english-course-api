// Package services implements the core business logic and state machine workflows.
package services

import (
	"errors"
	"time"

	"english-course-api/models"
	"english-course-api/repositories"
)

var (
	ErrPaymentNotFound      = errors.New("payment data not found")
	ErrPaymentAlreadyPaid   = errors.New("payment has already been settled")
	ErrPaymentInvalidStatus = errors.New("payment status is invalid for processing")
	ErrPaymentAmountInvalid = errors.New("payment amount does not match the invoice amount")
)

type ProcessPaymentRequest struct {
	PaymentMethod string  `json:"payment_method" binding:"required"` // e.g. "bank_transfer", "ewallet", "credit_card"
	Amount        float64 `json:"amount" binding:"required,gt=0"`
}

// PaymentService defines the business logic operations for Payment settlement.
type PaymentService interface {
	GetAll() ([]models.Payment, error)
	GetByID(id uint) (*models.Payment, error)
	Pay(paymentID uint, req ProcessPaymentRequest) (*models.Payment, error)
}

type paymentService struct {
	paymentRepo repositories.PaymentRepository
}

// NewPaymentService creates a new PaymentService instance.
func NewPaymentService(paymentRepo repositories.PaymentRepository) PaymentService {
	return &paymentService{paymentRepo: paymentRepo}
}

func (s *paymentService) GetAll() ([]models.Payment, error) {
	return s.paymentRepo.FindAll()
}

func (s *paymentService) GetByID(id uint) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

func (s *paymentService) Pay(paymentID uint, req ProcessPaymentRequest) (*models.Payment, error) {
	// 1. Fetch payment record
	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	// 2. Validate current payment status
	if payment.Status == "paid" {
		return nil, ErrPaymentAlreadyPaid
	}
	if payment.Status != "pending" {
		return nil, ErrPaymentInvalidStatus
	}

	// 3. Validate payment amount
	if req.Amount < payment.Amount {
		return nil, ErrPaymentAmountInvalid
	}

	now := time.Now()
	if err := s.paymentRepo.ProcessPaymentSuccess(paymentID, req.PaymentMethod, now); err != nil {
		return nil, err
	}

	// Fetch updated record after transaction
	return s.paymentRepo.FindByID(paymentID)
}

