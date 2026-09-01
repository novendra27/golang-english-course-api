package repositories

import (
	"errors"

	"english-course-api/models"

	"gorm.io/gorm"
)

// StudentRepository adalah interface abstraksi untuk operasi database entitas Student
type StudentRepository interface {
	Create(student *models.Student) error
	FindAll() ([]models.Student, error)
	FindByID(id uint) (*models.Student, error)
	FindByEmail(email string) (*models.Student, error)
	Update(student *models.Student) error
	Delete(id uint) error
}

type studentRepository struct {
	db *gorm.DB
}

// NewStudentRepository menginisialisasi implementasi StudentRepository
func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *studentRepository) FindAll() ([]models.Student, error) {
	var students []models.Student
	err := r.db.Order("id ASC").Find(&students).Error
	return students, err
}

func (r *studentRepository) FindByID(id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.First(&student, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) FindByEmail(email string) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("email = ?", email).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}

func (r *studentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Student{}, id).Error
}
