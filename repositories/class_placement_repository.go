package repositories

import (
	"errors"

	"english-course-api/models"

	"gorm.io/gorm"
)

// ClassPlacementRepository interface untuk operasi database penempatan kelas
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

func NewClassPlacementRepository(db *gorm.DB) ClassPlacementRepository {
	return &classPlacementRepository{db: db}
}

// CreatePlacement menyimpan penempatan kelas dan mengupdate status kelas menjadi 'full' jika kapasitas tercapai secara atomik
func (r *classPlacementRepository) CreatePlacement(placement *models.ClassPlacement, isClassFull bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan ClassPlacement
		if err := tx.Create(placement).Error; err != nil {
			return err
		}

		// 2. Jika kapasitas kelas telah penuh, update status kelas menjadi 'full'
		if isClassFull {
			if err := tx.Model(&models.Class{}).Where("id = ?", placement.ClassID).Update("status", "full").Error; err != nil {
				return err
			}
		}

		// 3. Update status Registration menjadi 'completed' (atau tetap 'registered' dengan placement aktif)
		// Registration tetap 'registered' atau 'completed'
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
