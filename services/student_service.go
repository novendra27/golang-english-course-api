package services

import (
	"errors"

	"english-course-api/models"
	"english-course-api/repositories"
)

var (
	ErrStudentNotFound      = errors.New("data student tidak ditemukan")
	ErrStudentEmailConflict = errors.New("email student sudah digunakan")
)

// Request DTOs
type CreateStudentRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

type UpdateStudentRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

// StudentService interface
type StudentService interface {
	Create(req CreateStudentRequest) (*models.Student, error)
	GetAll() ([]models.Student, error)
	GetByID(id uint) (*models.Student, error)
	Update(id uint, req UpdateStudentRequest) (*models.Student, error)
	Delete(id uint) error
}

type studentService struct {
	repo repositories.StudentRepository
}

// NewStudentService membuat instance baru StudentService
func NewStudentService(repo repositories.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) Create(req CreateStudentRequest) (*models.Student, error) {
	// 1. Cek duplikasi email
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrStudentEmailConflict
	}

	student := &models.Student{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	if err := s.repo.Create(student); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *studentService) GetAll() ([]models.Student, error) {
	return s.repo.FindAll()
}

func (s *studentService) GetByID(id uint) (*models.Student, error) {
	student, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, ErrStudentNotFound
	}
	return student, nil
}

func (s *studentService) Update(id uint, req UpdateStudentRequest) (*models.Student, error) {
	student, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, ErrStudentNotFound
	}

	// Cek apakah email diubah dan sudah dipakai oleh student lain
	if req.Email != student.Email {
		existing, err := s.repo.FindByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrStudentEmailConflict
		}
	}

	student.Name = req.Name
	student.Email = req.Email
	student.Phone = req.Phone

	if err := s.repo.Update(student); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *studentService) Delete(id uint) error {
	student, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if student == nil {
		return ErrStudentNotFound
	}
	return s.repo.Delete(id)
}
