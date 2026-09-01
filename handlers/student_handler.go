package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"english-course-api/services"
	"english-course-api/utils"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	service services.StudentService
}

func NewStudentHandler(service services.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

// Create godoc: POST /api/v1/students
func (h *StudentHandler) Create(c *gin.Context) {
	var req services.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payload request tidak valid", err.Error())
		return
	}

	student, err := h.service.Create(req)
	if err != nil {
		if errors.Is(err, services.ErrStudentEmailConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat data student", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Student berhasil dibuat", student)
}

// GetAll godoc: GET /api/v1/students
func (h *StudentHandler) GetAll(c *gin.Context) {
	students, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar student", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar student berhasil diambil", students)
}

// GetByID godoc: GET /api/v1/students/:id
func (h *StudentHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID student tidak valid", nil)
		return
	}

	student, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrStudentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail student", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail student berhasil diambil", student)
}

// Update godoc: PUT /api/v1/students/:id
func (h *StudentHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID student tidak valid", nil)
		return
	}

	var req services.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payload request tidak valid", err.Error())
		return
	}

	student, err := h.service.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrStudentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, services.ErrStudentEmailConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data student", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data student berhasil diperbarui", student)
}

// Delete godoc: DELETE /api/v1/students/:id
func (h *StudentHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID student tidak valid", nil)
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrStudentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus data student", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data student berhasil dihapus", nil)
}
