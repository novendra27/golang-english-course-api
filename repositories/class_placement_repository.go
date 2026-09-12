// Package repositories provides data access layer abstraction and database operations via GORM.
package repositories

import (
	"errors"

	"english-course-api/models"

	"gorm.io/gorm"
)

// ClassPlacementRepository defines the data access contract for ClassPlacement entities.
type ClassPlacementRepository interface {
	CreatePlacement(placement *models.ClassPlacement, isClassFull bool) error
	FindAll() ([]models.ClassPlacement, error)
	FindByID(id uint) (*models.ClassPlacement, error)
	FindByRegistrationID(registrationID uint) (*models.ClassPlacement, error)
	CountByClassID(classID uint) (int64, error)
}

type classPlacementRepository struct {
	db *gorm.DB
}

// NewClassPlacementRepository creates a new ClassPlacementRepository instance.
func NewClassPlacementRepository(db *gorm.DB) ClassPlacementRepository {
	return &classPlacementRepository{db: db}
}

// CreatePlacement records class placement and marks class status as 'full' if capacity is met atomically.
func (r *classPlacementRepository) CreatePlacement(placement *models.ClassPlacement, isClassFull bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Persist ClassPlacement record
		if err := tx.Create(placement).Error; err != nil {
			return err
		}

		// 2. If class reaches max capacity, update its status to 'full'
		if isClassFull {
			if err := tx.Model(&models.Class{}).Where("id = ?", placement.ClassID).Update("status", "full").Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *classPlacementRepository) FindAll() ([]models.ClassPlacement, error) {
	var placements []models.ClassPlacement
	err := r.db.
		Preload("Registration.Student").
		Preload("Registration.Course").
		Preload("Registration.Payment").
		Preload("Class.Course").
		Order("id DESC").
		Find(&placements).Error
	return placements, err
}

func (r *classPlacementRepository) FindByID(id uint) (*models.ClassPlacement, error) {
	var placement models.ClassPlacement
	err := r.db.
		Preload("Registration.Student").
		Preload("Registration.Course").
		Preload("Registration.Payment").
		Preload("Class.Course").
		First(&placement, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &placement, nil
}

func (r *classPlacementRepository) FindByRegistrationID(registrationID uint) (*models.ClassPlacement, error) {
	var placement models.ClassPlacement
	err := r.db.
		Where("registration_id = ?", registrationID).
		First(&placement).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &placement, nil
}

func (r *classPlacementRepository) CountByClassID(classID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.ClassPlacement{}).Where("class_id = ?", classID).Count(&count).Error
	return count, err
}

