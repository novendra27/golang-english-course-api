package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"english-course-api/services"
	"english-course-api/utils"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	service services.RegistrationService
}

func NewRegistrationHandler(service services.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{service: service}
}

// Register godoc
// @Summary      Mendaftar ke kursus (Registration)
// @Description  Mendaftarkan siswa ke kursus dan secara otomatis membuat tagihan Payment pending
// @Tags         Registrations
// @Accept       json
// @Produce      json
// @Param        request body services.CreateRegistrationRequest true "Payload pendaftaran"
// @Success      201  {object}  utils.APIResponse{data=models.Registration}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      409  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /registrations [post]
func (h *RegistrationHandler) Register(c *gin.Context) {
	var req services.CreateRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	registration, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, services.ErrRegistrationStudentNotFound) || errors.Is(err, services.ErrRegistrationCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, services.ErrRegistrationCourseInactive) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		if errors.Is(err, services.ErrRegistrationAlreadyActive) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat registrasi", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Pendaftaran kursus berhasil dibuat", registration)
}

// GetAll godoc
// @Summary      Mengambil daftar seluruh pendaftaran
// @Description  Mengembalikan semua data registrasi beserta relasi student, course, dan payment
// @Tags         Registrations
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Registration}
// @Failure      500  {object}  utils.APIResponse
// @Router       /registrations [get]
func (h *RegistrationHandler) GetAll(c *gin.Context) {
	registrations, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar pendaftaran", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar pendaftaran berhasil diambil", registrations)
}

// GetByID godoc
// @Summary      Mengambil detail pendaftaran
// @Description  Mengembalikan detail registrasi berdasarkan ID
// @Tags         Registrations
// @Produce      json
// @Param        id   path      int  true  "Registration ID"
// @Success      200  {object}  utils.APIResponse{data=models.Registration}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /registrations/{id} [get]
func (h *RegistrationHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID pendaftaran tidak valid", nil)
		return
	}

	reg, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRegistrationNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail pendaftaran", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail pendaftaran berhasil diambil", reg)
}

// GetByStudentID godoc
// @Summary      Mengambil riwayat pendaftaran student
// @Description  Mengembalikan daftar kursus yang didaftarkan oleh student tertentu
// @Tags         Students
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  utils.APIResponse{data=[]models.Registration}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /students/{id}/registrations [get]
func (h *RegistrationHandler) GetByStudentID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID student tidak valid", nil)
		return
	}

	regs, err := h.service.GetByStudentID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRegistrationStudentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar pendaftaran student", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar pendaftaran student berhasil diambil", regs)
}

// GetByCourseID godoc
// @Summary      Mengambil daftar pendaftaran pada course
// @Description  Mengembalikan semua siswa yang terdaftar di course tertentu
// @Tags         Courses
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  utils.APIResponse{data=[]models.Registration}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /courses/{id}/registrations [get]
func (h *RegistrationHandler) GetByCourseID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID course tidak valid", nil)
		return
	}

	regs, err := h.service.GetByCourseID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRegistrationCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar pendaftaran course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar pendaftaran course berhasil diambil", regs)
}

// CancelRegistration godoc
// @Summary      Membatalkan pendaftaran
// @Description  Mengubah status registrasi menjadi cancelled jika belum selesai
// @Tags         Registrations
// @Produce      json
// @Param        id   path      int  true  "Registration ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /registrations/{id}/cancel [put]
func (h *RegistrationHandler) CancelRegistration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID pendaftaran tidak valid", nil)
		return
	}

	err = h.service.CancelRegistration(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRegistrationNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, services.ErrRegistrationCannotBeCanceled) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membatalkan pendaftaran", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pendaftaran berhasil dibatalkan", nil)
}
