package services

import (
	"errors"
	"time"

	"english-course-api/models"
	"english-course-api/repositories"
)

var (
	ErrPlacementNotFound        = errors.New("data penempatan kelas tidak ditemukan")
	ErrPlacementRegNotFound     = errors.New("data registrasi tidak ditemukan")
	ErrPlacementClassNotFound   = errors.New("data kelas tidak ditemukan")
	ErrPlacementPaymentRequired = errors.New("student belum menyelesaikan pembayaran (status registrasi belum 'registered' / payment belum 'paid')")
	ErrPlacementCourseMismatch  = errors.New("kelas yang dipilih tidak sesuai dengan course yang didaftarkan")
	ErrPlacementAlreadyAssigned = errors.New("registrasi ini sudah pernah ditempatkan ke dalam kelas")
	ErrPlacementClassFull       = errors.New("kapasitas kelas sudah penuh")
	ErrPlacementClassClosed     = errors.New("kelas sedang ditutup untuk penempatan baru")
)

type CreateClassPlacementRequest struct {
	RegistrationID uint `json:"registration_id" binding:"required"`
	ClassID        uint `json:"class_id" binding:"required"`
}

type ClassPlacementService interface {
	PlaceStudent(req CreateClassPlacementRequest) (*models.ClassPlacement, error)
	GetAll() ([]models.ClassPlacement, error)
	GetByID(id uint) (*models.ClassPlacement, error)
}

type classPlacementService struct {
	placementRepo repositories.ClassPlacementRepository
	regRepo       repositories.RegistrationRepository
	classRepo     repositories.ClassRepository
}

func NewClassPlacementService(
	placementRepo repositories.ClassPlacementRepository,
	regRepo repositories.RegistrationRepository,
	classRepo repositories.ClassRepository,
) ClassPlacementService {
	return &classPlacementService{
		placementRepo: placementRepo,
		regRepo:       regRepo,
		classRepo:     classRepo,
	}
}

func (s *classPlacementService) PlaceStudent(req CreateClassPlacementRequest) (*models.ClassPlacement, error) {
	// 1. Validasi Keberadaan Registrasi
	reg, err := s.regRepo.FindByID(req.RegistrationID)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		return nil, ErrPlacementRegNotFound
	}

	// 2. ATURAN BISNIS: Registrasi harus berstatus 'registered' (Pembayaran sudah 'paid')
	if reg.Status != "registered" {
		return nil, ErrPlacementPaymentRequired
	}

	// 3. Validasi Keberadaan Kelas
	class, err := s.classRepo.FindByID(req.ClassID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, ErrPlacementClassNotFound
	}

	if class.Status == "closed" {
		return nil, ErrPlacementClassClosed
	}

	// 4. ATURAN BISNIS: Course pendaftaran harus SAMA dengan Course pada Kelas
	if reg.CourseID != class.CourseID {
		return nil, ErrPlacementCourseMismatch
	}

	// 5. ATURAN BISNIS: Cek apakah registrasi sudah pernah di-assign ke kelas
	existingPlacement, err := s.placementRepo.FindByRegistrationID(req.RegistrationID)
	if err != nil {
		return nil, err
	}
	if existingPlacement != nil {
		return nil, ErrPlacementAlreadyAssigned
	}

	// 6. ATURAN BISNIS: Cek Kapasitas Kelas
	currentCount, err := s.placementRepo.CountByClassID(req.ClassID)
	if err != nil {
		return nil, err
	}

	if int(currentCount) >= class.Capacity {
		return nil, ErrPlacementClassFull
	}

	isClassFull := (int(currentCount)+1 >= class.Capacity)

	// 7. Simpan Placement dalam Database Transaction
	now := time.Now()
	placement := &models.ClassPlacement{
		RegistrationID: req.RegistrationID,
		ClassID:        req.ClassID,
		PlacementDate:  now,
	}

	if err := s.placementRepo.CreatePlacement(placement, isClassFull); err != nil {
		return nil, err
	}

	// Attach relasi untuk response
	placement.Registration = reg
	placement.Class = class

	return placement, nil
}

func (s *classPlacementService) GetAll() ([]models.ClassPlacement, error) {
	return s.placementRepo.FindAll()
}

func (s *classPlacementService) GetByID(id uint) (*models.ClassPlacement, error) {
	placement, err := s.placementRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if placement == nil {
		return nil, ErrPlacementNotFound
	}
	return placement, nil
}
