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

// Register godoc: POST /api/v1/registrations
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

// GetAll godoc: GET /api/v1/registrations
func (h *RegistrationHandler) GetAll(c *gin.Context) {
	registrations, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar pendaftaran", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar pendaftaran berhasil diambil", registrations)
}

// GetByID godoc: GET /api/v1/registrations/:id
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

// GetByStudentID godoc: GET /api/v1/students/:id/registrations
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

// GetByCourseID godoc: GET /api/v1/courses/:id/registrations
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

// CancelRegistration godoc: PUT /api/v1/registrations/:id/cancel
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
