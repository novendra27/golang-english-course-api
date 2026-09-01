package services

import (
	"errors"
	"time"

	"english-course-api/models"
	"english-course-api/repositories"
)

var (
	ErrPaymentNotFound      = errors.New("data payment tidak ditemukan")
	ErrPaymentAlreadyPaid   = errors.New("pembayaran sudah berstatus lunas (paid)")
	ErrPaymentInvalidStatus = errors.New("status pembayaran tidak valid untuk diproses")
	ErrPaymentAmountInvalid = errors.New("jumlah pembayaran tidak sesuai dengan tagihan")
)

type ProcessPaymentRequest struct {
	PaymentMethod string  `json:"payment_method" binding:"required"` // e.g. "bank_transfer", "ewallet", "credit_card"
	Amount        float64 `json:"amount" binding:"required,gt=0"`
}

type PaymentService interface {
	GetAll() ([]models.Payment, error)
	GetByID(id uint) (*models.Payment, error)
	Pay(paymentID uint, req ProcessPaymentRequest) (*models.Payment, error)
}

type paymentService struct {
	paymentRepo repositories.PaymentRepository
}

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
	// 1. Ambil data payment
	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	// 2. Validasi Status Pembayaran
	if payment.Status == "paid" {
		return nil, ErrPaymentAlreadyPaid
	}
	if payment.Status != "pending" {
		return nil, ErrPaymentInvalidStatus
	}

	// 3. Validasi Jumlah Pembayaran
	if req.Amount < payment.Amount {
		return nil, ErrPaymentAmountInvalid
	}

	now := time.Now()
	if err := s.paymentRepo.ProcessPaymentSuccess(paymentID, req.PaymentMethod, now); err != nil {
		return nil, err
	}

	// Ambil data terbaru setelah transaksi
	return s.paymentRepo.FindByID(paymentID)
}
