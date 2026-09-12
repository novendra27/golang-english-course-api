package services_test

import (
	"errors"
	"testing"
	"time"

	"english-course-api/models"
	"english-course-api/services"
)

// Mock Repositories
type mockStudentRepo struct {
	students map[uint]*models.Student
}

func (m *mockStudentRepo) Create(s *models.Student) error           { return nil }
func (m *mockStudentRepo) FindAll() ([]models.Student, error)       { return nil, nil }
func (m *mockStudentRepo) FindByEmail(email string) (*models.Student, error) { return nil, nil }
func (m *mockStudentRepo) Update(s *models.Student) error           { return nil }
func (m *mockStudentRepo) Delete(id uint) error                     { return nil }
func (m *mockStudentRepo) FindByID(id uint) (*models.Student, error) {
	if s, ok := m.students[id]; ok {
		return s, nil
	}
	return nil, nil
}

type mockCourseRepo struct {
	courses map[uint]*models.Course
}

func (m *mockCourseRepo) Create(c *models.Course) error     { return nil }
func (m *mockCourseRepo) FindAll() ([]models.Course, error) { return nil, nil }
func (m *mockCourseRepo) Update(c *models.Course) error     { return nil }
func (m *mockCourseRepo) Delete(id uint) error               { return nil }
func (m *mockCourseRepo) FindByID(id uint) (*models.Course, error) {
	if c, ok := m.courses[id]; ok {
		return c, nil
	}
	return nil, nil
}

type mockRegistrationRepo struct {
	registrations map[uint]*models.Registration
	activeKey     string // "studentID_courseID"
}

func (m *mockRegistrationRepo) CreateWithPayment(reg *models.Registration, pay *models.Payment) error {
	reg.ID = uint(len(m.registrations) + 1)
	m.registrations[reg.ID] = reg
	return nil
}

func (m *mockRegistrationRepo) FindAll() ([]models.Registration, error) { return nil, nil }
func (m *mockRegistrationRepo) FindByStudentID(id uint) ([]models.Registration, error) { return nil, nil }
func (m *mockRegistrationRepo) FindByCourseID(id uint) ([]models.Registration, error) { return nil, nil }
func (m *mockRegistrationRepo) FindByID(id uint) (*models.Registration, error) {
	if r, ok := m.registrations[id]; ok {
		return r, nil
	}
	return nil, nil
}

func (m *mockRegistrationRepo) FindActiveRegistration(studentID, courseID uint) (*models.Registration, error) {
	for _, r := range m.registrations {
		if r.StudentID == studentID && r.CourseID == courseID && (r.Status == "pending" || r.Status == "registered") {
			return r, nil
		}
	}
	return nil, nil
}

func (m *mockRegistrationRepo) UpdateStatus(id uint, status string) error {
	if r, ok := m.registrations[id]; ok {
		r.Status = status
		return nil
	}
	return errors.New("not found")
}

// --- UNIT TESTS ---

func TestRegistrationService_Register(t *testing.T) {
	studentRepo := &mockStudentRepo{
		students: map[uint]*models.Student{
			1: {ID: 1, Name: "Budi", Email: "budi@example.com"},
		},
	}
	courseRepo := &mockCourseRepo{
		courses: map[uint]*models.Course{
			1: {ID: 1, Name: "Basic English", Price: 750000, Status: "active"},
			2: {ID: 2, Name: "Inactive Course", Price: 500000, Status: "inactive"},
		},
	}
	regRepo := &mockRegistrationRepo{
		registrations: make(map[uint]*models.Registration),
	}

	service := services.NewRegistrationService(regRepo, studentRepo, courseRepo)

	t.Run("Success Register", func(t *testing.T) {
		req := services.CreateRegistrationRequest{StudentID: 1, CourseID: 1}
		reg, err := service.Register(req)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if reg.Status != "pending" {
			t.Errorf("expected status 'pending', got: %s", reg.Status)
		}
		if reg.Payment == nil || reg.Payment.Amount != 750000 {
			t.Errorf("expected payment amount 750000, got: %v", reg.Payment)
		}
	})

	t.Run("Reject Duplicate Active Registration", func(t *testing.T) {
		req := services.CreateRegistrationRequest{StudentID: 1, CourseID: 1}
		_, err := service.Register(req)
		if !errors.Is(err, services.ErrRegistrationAlreadyActive) {
			t.Errorf("expected ErrRegistrationAlreadyActive, got: %v", err)
		}
	})

	t.Run("Reject Inactive Course", func(t *testing.T) {
		req := services.CreateRegistrationRequest{StudentID: 1, CourseID: 2}
		_, err := service.Register(req)
		if !errors.Is(err, services.ErrRegistrationCourseInactive) {
			t.Errorf("expected ErrRegistrationCourseInactive, got: %v", err)
		}
	})

	t.Run("Reject Non-Existent Student", func(t *testing.T) {
		req := services.CreateRegistrationRequest{StudentID: 999, CourseID: 1}
		_, err := service.Register(req)
		if !errors.Is(err, services.ErrRegistrationStudentNotFound) {
			t.Errorf("expected ErrRegistrationStudentNotFound, got: %v", err)
		}
	})
}

func TestRegistrationService_CancelRegistration(t *testing.T) {
	regRepo := &mockRegistrationRepo{
		registrations: map[uint]*models.Registration{
			1: {ID: 1, StudentID: 1, CourseID: 1, Status: "pending", RegistrationDate: time.Now()},
			2: {ID: 2, StudentID: 1, CourseID: 1, Status: "completed", RegistrationDate: time.Now()},
		},
	}
	service := services.NewRegistrationService(regRepo, &mockStudentRepo{}, &mockCourseRepo{})

	t.Run("Success Cancel Pending Registration", func(t *testing.T) {
		err := service.CancelRegistration(1)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if regRepo.registrations[1].Status != "cancelled" {
			t.Errorf("expected status 'cancelled', got: %s", regRepo.registrations[1].Status)
		}
	})

	t.Run("Reject Cancel Completed Registration", func(t *testing.T) {
		err := service.CancelRegistration(2)
		if !errors.Is(err, services.ErrRegistrationCannotBeCanceled) {
			t.Errorf("expected ErrRegistrationCannotBeCanceled, got: %v", err)
		}
	})
}
