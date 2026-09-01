package services

import (
	"errors"

	"english-course-api/models"
	"english-course-api/repositories"
)

var (
	ErrCourseNotFound = errors.New("data course tidak ditemukan")
)

// Request DTOs
type CreateCourseRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Duration    string  `json:"duration" binding:"required"`
	Status      string  `json:"status"` // default: 'active'
}

type UpdateCourseRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Duration    string  `json:"duration" binding:"required"`
	Status      string  `json:"status" binding:"required"`
}

// CourseService interface
type CourseService interface {
	Create(req CreateCourseRequest) (*models.Course, error)
	GetAll() ([]models.Course, error)
	GetByID(id uint) (*models.Course, error)
	Update(id uint, req UpdateCourseRequest) (*models.Course, error)
	Delete(id uint) error
}

type courseService struct {
	repo repositories.CourseRepository
}

// NewCourseService membuat instance baru CourseService
func NewCourseService(repo repositories.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) Create(req CreateCourseRequest) (*models.Course, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}

	course := &models.Course{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Duration:    req.Duration,
		Status:      status,
	}

	if err := s.repo.Create(course); err != nil {
		return nil, err
	}

	return course, nil
}

func (s *courseService) GetAll() ([]models.Course, error) {
	return s.repo.FindAll()
}

func (s *courseService) GetByID(id uint) (*models.Course, error) {
	course, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}
	return course, nil
}

func (s *courseService) Update(id uint, req UpdateCourseRequest) (*models.Course, error) {
	course, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}

	course.Name = req.Name
	course.Description = req.Description
	course.Price = req.Price
	course.Duration = req.Duration
	course.Status = req.Status

	if err := s.repo.Update(course); err != nil {
		return nil, err
	}

	return course, nil
}

func (s *courseService) Delete(id uint) error {
	course, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if course == nil {
		return ErrCourseNotFound
	}
	return s.repo.Delete(id)
}
