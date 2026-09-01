package repositories

import (
	"errors"

	"english-course-api/models"

	"gorm.io/gorm"
)

// ClassRepository interface untuk database operation entitas Class
type ClassRepository interface {
	Create(class *models.Class) error
	FindAll() ([]models.Class, error)
	FindByID(id uint) (*models.Class, error)
	FindByCourseID(courseID uint) ([]models.Class, error)
	Update(class *models.Class) error
	Delete(id uint) error
	GetStudentsByClassID(classID uint) ([]models.Student, error)
}

type classRepository struct {
	db *gorm.DB
}

// NewClassRepository menginisialisasi implementasi ClassRepository
func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(class *models.Class) error {
	return r.db.Create(class).Error
}

func (r *classRepository) FindAll() ([]models.Class, error) {
	var classes []models.Class
	err := r.db.Preload("Course").Order("id ASC").Find(&classes).Error
	return classes, err
}

func (r *classRepository) FindByID(id uint) (*models.Class, error) {
	var class models.Class
	err := r.db.Preload("Course").First(&class, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) FindByCourseID(courseID uint) ([]models.Class, error) {
	var classes []models.Class
	err := r.db.Where("course_id = ?", courseID).Order("id ASC").Find(&classes).Error
	return classes, err
}

func (r *classRepository) Update(class *models.Class) error {
	return r.db.Save(class).Error
}

func (r *classRepository) Delete(id uint) error {
	return r.db.Delete(&models.Class{}, id).Error
}

func (r *classRepository) GetStudentsByClassID(classID uint) ([]models.Student, error) {
	var students []models.Student
	// Join: students -> registrations -> class_placements
	err := r.db.
		Table("students").
		Joins("JOIN registrations ON registrations.student_id = students.id").
		Joins("JOIN class_placements ON class_placements.registration_id = registrations.id").
		Where("class_placements.class_id = ?", classID).
		Find(&students).Error

	return students, err
}
