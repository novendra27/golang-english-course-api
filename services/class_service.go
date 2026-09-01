package services

import (
	"errors"

	"english-course-api/models"
	"english-course-api/repositories"
)

var (
	ErrClassNotFound       = errors.New("data class tidak ditemukan")
	ErrClassCourseNotFound = errors.New("course yang direferensikan tidak ditemukan")
)

// Request DTOs
type CreateClassRequest struct {
	CourseID uint   `json:"course_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity" binding:"required,gt=0"`
	Schedule string `json:"schedule" binding:"required"`
	Status   string `json:"status"` // 'open', 'full', 'closed'
}

type UpdateClassRequest struct {
	CourseID uint   `json:"course_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity" binding:"required,gt=0"`
	Schedule string `json:"schedule" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

// ClassService interface
type ClassService interface {
	Create(req CreateClassRequest) (*models.Class, error)
	GetAll() ([]models.Class, error)
	GetByID(id uint) (*models.Class, error)
	Update(id uint, req UpdateClassRequest) (*models.Class, error)
	Delete(id uint) error
	GetStudents(classID uint) ([]models.Student, error)
}

type classService struct {
	classRepo  repositories.ClassRepository
	courseRepo repositories.CourseRepository
}

// NewClassService membuat instance baru ClassService
func NewClassService(classRepo repositories.ClassRepository, courseRepo repositories.CourseRepository) ClassService {
	return &classService{
		classRepo:  classRepo,
		courseRepo: courseRepo,
	}
}

func (s *classService) Create(req CreateClassRequest) (*models.Class, error) {
	// 1. Validasi keberadaan Course
	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrClassCourseNotFound
	}

	status := req.Status
	if status == "" {
		status = "open"
	}

	class := &models.Class{
		CourseID: req.CourseID,
		Name:     req.Name,
		Capacity: req.Capacity,
		Schedule: req.Schedule,
		Status:   status,
	}

	if err := s.classRepo.Create(class); err != nil {
		return nil, err
	}

	// Attach course object ke response
	class.Course = course
	return class, nil
}

func (s *classService) GetAll() ([]models.Class, error) {
	return s.classRepo.FindAll()
}

func (s *classService) GetByID(id uint) (*models.Class, error) {
	class, err := s.classRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, ErrClassNotFound
	}
	return class, nil
}

func (s *classService) Update(id uint, req UpdateClassRequest) (*models.Class, error) {
	class, err := s.classRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, ErrClassNotFound
	}

	// Validasi CourseID jika diubah
	if req.CourseID != class.CourseID {
		course, err := s.courseRepo.FindByID(req.CourseID)
		if err != nil {
			return nil, err
		}
		if course == nil {
			return nil, ErrClassCourseNotFound
		}
		class.Course = course
	}

	class.CourseID = req.CourseID
	class.Name = req.Name
	class.Capacity = req.Capacity
	class.Schedule = req.Schedule
	class.Status = req.Status

	if err := s.classRepo.Update(class); err != nil {
		return nil, err
	}

	return class, nil
}

func (s *classService) Delete(id uint) error {
	class, err := s.classRepo.FindByID(id)
	if err != nil {
		return err
	}
	if class == nil {
		return ErrClassNotFound
	}
	return s.classRepo.Delete(id)
}

func (s *classService) GetStudents(classID uint) ([]models.Student, error) {
	class, err := s.classRepo.FindByID(classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, ErrClassNotFound
	}

	return s.classRepo.GetStudentsByClassID(classID)
}
