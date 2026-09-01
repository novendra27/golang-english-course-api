package services

import (
	"errors"
	"time"

	"english-course-api/models"
	"english-course-api/repositories"
)

var (
	ErrRegistrationNotFound         = errors.New("data registrasi tidak ditemukan")
	ErrRegistrationStudentNotFound  = errors.New("student tidak ditemukan")
	ErrRegistrationCourseNotFound   = errors.New("course tidak ditemukan")
	ErrRegistrationCourseInactive   = errors.New("course sedang tidak aktif")
	ErrRegistrationAlreadyActive    = errors.New("student masih memiliki pendaftaran aktif untuk course ini")
	ErrRegistrationCannotBeCanceled = errors.New("registrasi tidak dapat dibatalkan pada status saat ini")
)

type CreateRegistrationRequest struct {
	StudentID uint `json:"student_id" binding:"required"`
	CourseID  uint `json:"course_id" binding:"required"`
}

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
	// 1. Validasi keberadaan Student
	student, err := s.studentRepo.FindByID(req.StudentID)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, ErrRegistrationStudentNotFound
	}

	// 2. Validasi keberadaan Course & Status
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

	// 3. Validasi Duplikasi Registrasi Aktif
	activeReg, err := s.regRepo.FindActiveRegistration(req.StudentID, req.CourseID)
	if err != nil {
		return nil, err
	}
	if activeReg != nil {
		return nil, ErrRegistrationAlreadyActive
	}

	// 4. Siapkan Model Registrasi & Payment
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

	// Attach relasi untuk response
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

	if reg.Status == "completed" || reg.Status == "cancelled" {
		return ErrRegistrationCannotBeCanceled
	}

	return s.regRepo.UpdateStatus(id, "cancelled")
}
