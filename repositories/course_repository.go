// Package repositories provides data access layer abstraction and database operations via GORM.
package repositories

import (
	"errors"

	"english-course-api/models"

	"gorm.io/gorm"
)

// CourseRepository defines the data access contract for Course entities.
type CourseRepository interface {
	Create(course *models.Course) error
	FindAll() ([]models.Course, error)
	FindByID(id uint) (*models.Course, error)
	Update(course *models.Course) error
	Delete(id uint) error
}

type courseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new CourseRepository instance.
func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) FindAll() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Order("id ASC").Find(&courses).Error
	return courses, err
}

func (r *courseRepository) FindByID(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Classes").First(&course, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Course{}, id).Error
}

