package services_test

import (
	"errors"
	"testing"

	"english-course-api/models"
	"english-course-api/services"
)

type mockClassRepo struct {
	classes map[uint]*models.Class
}

func (m *mockClassRepo) Create(c *models.Class) error                       { return nil }
func (m *mockClassRepo) FindAll() ([]models.Class, error)                   { return nil, nil }
func (m *mockClassRepo) FindByCourseID(courseID uint) ([]models.Class, error) { return nil, nil }
func (m *mockClassRepo) Update(c *models.Class) error                       { return nil }
func (m *mockClassRepo) Delete(id uint) error                               { return nil }
func (m *mockClassRepo) GetStudentsByClassID(classID uint) ([]models.Student, error) {
	return nil, nil
}
func (m *mockClassRepo) FindByID(id uint) (*models.Class, error) {
	if c, ok := m.classes[id]; ok {
		return c, nil
	}
	return nil, nil
}

type mockPlacementRepo struct {
	placements map[uint]*models.ClassPlacement
	countByClass map[uint]int64
}

func (m *mockPlacementRepo) CreatePlacement(p *models.ClassPlacement, isFull bool) error {
	p.ID = uint(len(m.placements) + 1)
	m.placements[p.ID] = p
	m.countByClass[p.ClassID]++
	return nil
}

func (m *mockPlacementRepo) FindAll() ([]models.ClassPlacement, error) { return nil, nil }
func (m *mockPlacementRepo) FindByID(id uint) (*models.ClassPlacement, error) {
	if p, ok := m.placements[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *mockPlacementRepo) FindByRegistrationID(regID uint) (*models.ClassPlacement, error) {
	for _, p := range m.placements {
		if p.RegistrationID == regID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockPlacementRepo) CountByClassID(classID uint) (int64, error) {
	return m.countByClass[classID], nil
}

// --- UNIT TESTS ---

func TestClassPlacementService_PlaceStudent(t *testing.T) {
	classRepo := &mockClassRepo{
		classes: map[uint]*models.Class{
			1: {ID: 1, CourseID: 1, Name: "Basic English A", Capacity: 2, Status: "open"},
			2: {ID: 2, CourseID: 2, Name: "Intermediate English A", Capacity: 10, Status: "open"},
			3: {ID: 3, CourseID: 1, Name: "Full Class", Capacity: 1, Status: "full"},
		},
	}

	regRepo := &mockRegistrationRepo{
		registrations: map[uint]*models.Registration{
			1: {ID: 1, StudentID: 1, CourseID: 1, Status: "registered"}, // Lunas & siap placement
			2: {ID: 2, StudentID: 2, CourseID: 1, Status: "pending"},    // Belum bayar
			3: {ID: 3, StudentID: 3, CourseID: 1, Status: "registered"}, // Lunas
		},
	}

	placementRepo := &mockPlacementRepo{
		placements:   make(map[uint]*models.ClassPlacement),
		countByClass: map[uint]int64{3: 1}, // Kelas 3 sudah terisi 1 orang (penuh)
	}

	service := services.NewClassPlacementService(placementRepo, regRepo, classRepo)

	t.Run("Success Placement", func(t *testing.T) {
		req := services.CreateClassPlacementRequest{RegistrationID: 1, ClassID: 1}
		placement, err := service.PlaceStudent(req)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if placement.ClassID != 1 || placement.RegistrationID != 1 {
			t.Errorf("unexpected placement data: %+v", placement)
		}
	})

	t.Run("Reject Unpaid Student (Pending Status)", func(t *testing.T) {
		req := services.CreateClassPlacementRequest{RegistrationID: 2, ClassID: 1}
		_, err := service.PlaceStudent(req)
		if !errors.Is(err, services.ErrPlacementPaymentRequired) {
			t.Errorf("expected ErrPlacementPaymentRequired, got: %v", err)
		}
	})

	t.Run("Reject Course Mismatch", func(t *testing.T) {
		// Reg 3 daftar Course 1, tapi ditempatkan di Class 2 (Course 2)
		req := services.CreateClassPlacementRequest{RegistrationID: 3, ClassID: 2}
		_, err := service.PlaceStudent(req)
		if !errors.Is(err, services.ErrPlacementCourseMismatch) {
			t.Errorf("expected ErrPlacementCourseMismatch, got: %v", err)
		}
	})

	t.Run("Reject Duplicate Placement for Same Registration", func(t *testing.T) {
		// Reg 1 sudah ditempatkan di test case 1
		req := services.CreateClassPlacementRequest{RegistrationID: 1, ClassID: 1}
		_, err := service.PlaceStudent(req)
		if !errors.Is(err, services.ErrPlacementAlreadyAssigned) {
			t.Errorf("expected ErrPlacementAlreadyAssigned, got: %v", err)
		}
	})

	t.Run("Reject Full Class Capacity", func(t *testing.T) {
		req := services.CreateClassPlacementRequest{RegistrationID: 3, ClassID: 3}
		_, err := service.PlaceStudent(req)
		if !errors.Is(err, services.ErrPlacementClassFull) {
			t.Errorf("expected ErrPlacementClassFull, got: %v", err)
		}
	})
}
