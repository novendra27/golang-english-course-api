// Package services implements the core business logic and state machine workflows.
package services

import (
	"errors"
	"time"

	"english-course-api/models"
	"english-course-api/repositories"
)

var (
	ErrPlacementNotFound        = errors.New("class placement data not found")
	ErrPlacementRegNotFound     = errors.New("registration data not found")
	ErrPlacementClassNotFound   = errors.New("class data not found")
	ErrPlacementPaymentRequired = errors.New("student has not completed payment (registration status must be 'registered' / payment must be 'paid')")
	ErrPlacementCourseMismatch  = errors.New("selected class does not match the registered course")
	ErrPlacementAlreadyAssigned = errors.New("registration is already assigned to a class")
	ErrPlacementClassFull       = errors.New("class capacity is full")
	ErrPlacementClassClosed     = errors.New("class is closed for new placements")
)

type CreateClassPlacementRequest struct {
	RegistrationID uint `json:"registration_id" binding:"required"`
	ClassID        uint `json:"class_id" binding:"required"`
}

// ClassPlacementService defines the business logic operations for Class Placement workflows.
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

// NewClassPlacementService creates a new ClassPlacementService instance.
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
	// 1. Validate existence of Registration
	reg, err := s.regRepo.FindByID(req.RegistrationID)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		return nil, ErrPlacementRegNotFound
	}

	// 2. BUSINESS RULE: Registration status must be 'registered' (Payment is settled)
	if reg.Status != "registered" {
		return nil, ErrPlacementPaymentRequired
	}

	// 3. Validate existence of Class
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

	// 4. BUSINESS RULE: Registration Course must MATCH the Class Course
	if reg.CourseID != class.CourseID {
		return nil, ErrPlacementCourseMismatch
	}

	// 5. BUSINESS RULE: Prevent duplicate class placements for the same registration
	existingPlacement, err := s.placementRepo.FindByRegistrationID(req.RegistrationID)
	if err != nil {
		return nil, err
	}
	if existingPlacement != nil {
		return nil, ErrPlacementAlreadyAssigned
	}

	// 6. BUSINESS RULE: Check Class Capacity
	currentCount, err := s.placementRepo.CountByClassID(req.ClassID)
	if err != nil {
		return nil, err
	}

	if int(currentCount) >= class.Capacity {
		return nil, ErrPlacementClassFull
	}

	isClassFull := (int(currentCount)+1 >= class.Capacity)

	// 7. Persist ClassPlacement in an atomic database transaction
	now := time.Now()
	placement := &models.ClassPlacement{
		RegistrationID: req.RegistrationID,
		ClassID:        req.ClassID,
		PlacementDate:  now,
	}

	if err := s.placementRepo.CreatePlacement(placement, isClassFull); err != nil {
		return nil, err
	}

	// Attach relations for response
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

