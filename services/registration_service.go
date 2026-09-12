// Package services implements the core business logic and state machine workflows.
package services

import (
	"errors"
	"time"

	"english-course-api/models"
	"english-course-api/repositories"
)

var (
	ErrRegistrationNotFound         = errors.New("registration data not found")
	ErrRegistrationStudentNotFound  = errors.New("student not found")
	ErrRegistrationCourseNotFound   = errors.New("course not found")
	ErrRegistrationCourseInactive   = errors.New("course is currently inactive")
	ErrRegistrationAlreadyActive    = errors.New("student already has an active registration for this course")
	ErrRegistrationCannotBeCanceled = errors.New("registration cannot be cancelled in its current status")
)

type CreateRegistrationRequest struct {
	StudentID uint `json:"student_id" binding:"required"`
	CourseID  uint `json:"course_id" binding:"required"`
}

// RegistrationService defines the business logic operations for Registration workflows.
type RegistrationService interface {
	Register(req CreateRegistrationRequest) (*models.Registration, error)
	GetAll() ([]models.Registration, error)
	GetByID(id uint) (*models.Registration, error)
	GetByStudentID(studentID uint) ([]models.Registration, error)
	GetByCourseID(courseID uint) ([]models.Registration, error)
	CancelRegistration(id uint) error
}

type registrationService struct {
	regRepo     repositories.RegistrationRepository
	studentRepo repositories.StudentRepository
	courseRepo  repositories.CourseRepository
}

// NewRegistrationService creates a new RegistrationService instance.
func NewRegistrationService(
	regRepo repositories.RegistrationRepository,
	studentRepo repositories.StudentRepository,
	courseRepo repositories.CourseRepository,
) RegistrationService {
	return &registrationService{
		regRepo:     regRepo,
		studentRepo: studentRepo,
		courseRepo:  courseRepo,
	}
}

func (s *registrationService) Register(req CreateRegistrationRequest) (*models.Registration, error) {
	// 1. Validate existence of Student
	student, err := s.studentRepo.FindByID(req.StudentID)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, ErrRegistrationStudentNotFound
	}

	// 2. Validate existence and active status of Course
	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrRegistrationCourseNotFound
	}
	if course.Status != "active" {
		return nil, ErrRegistrationCourseInactive
	}

	// 3. Prevent duplicate active registrations (status: 'pending' or 'registered')
	activeReg, err := s.regRepo.FindActiveRegistration(req.StudentID, req.CourseID)
	if err != nil {
		return nil, err
	}
	if activeReg != nil {
		return nil, ErrRegistrationAlreadyActive
	}

	// 4. Prepare Registration and initial Payment records
	now := time.Now()
	registration := &models.Registration{
		StudentID:        req.StudentID,
		CourseID:         req.CourseID,
		RegistrationDate: now,
		Status:           "pending",
	}

	payment := &models.Payment{
		Amount:        course.Price,
		PaymentMethod: "pending",
		Status:        "pending",
	}

	if err := s.regRepo.CreateWithPayment(registration, payment); err != nil {
		return nil, err
	}

	// Attach relations for response
	registration.Student = student
	registration.Course = course
	registration.Payment = payment

	return registration, nil
}

func (s *registrationService) GetAll() ([]models.Registration, error) {
	return s.regRepo.FindAll()
}

func (s *registrationService) GetByID(id uint) (*models.Registration, error) {
	reg, err := s.regRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		return nil, ErrRegistrationNotFound
	}
	return reg, nil
}

func (s *registrationService) GetByStudentID(studentID uint) ([]models.Registration, error) {
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, ErrRegistrationStudentNotFound
	}
	return s.regRepo.FindByStudentID(studentID)
}

func (s *registrationService) GetByCourseID(courseID uint) ([]models.Registration, error) {
	course, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrRegistrationCourseNotFound
	}
	return s.regRepo.FindByCourseID(courseID)
}

func (s *registrationService) CancelRegistration(id uint) error {
	reg, err := s.regRepo.FindByID(id)
	if err != nil {
		return err
	}
	if reg == nil {
		return ErrRegistrationNotFound
	}

	// Only 'pending' registrations can be cancelled
	if reg.Status != "pending" {
		return ErrRegistrationCannotBeCanceled
	}

	return s.regRepo.UpdateStatus(id, "cancelled")
}
