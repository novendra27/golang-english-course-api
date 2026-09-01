package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"english-course-api/services"
	"english-course-api/utils"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	service services.ClassService
}

func NewClassHandler(service services.ClassService) *ClassHandler {
	return &ClassHandler{service: service}
}

// Create godoc: POST /api/v1/classes
func (h *ClassHandler) Create(c *gin.Context) {
	var req services.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payload request tidak valid", err.Error())
		return
	}

	class, err := h.service.Create(req)
	if err != nil {
		if errors.Is(err, services.ErrClassCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat data class", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Class berhasil dibuat", class)
}

// GetAll godoc: GET /api/v1/classes
func (h *ClassHandler) GetAll(c *gin.Context) {
	classes, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar class", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar class berhasil diambil", classes)
}

// GetByID godoc: GET /api/v1/classes/:id
func (h *ClassHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID class tidak valid", nil)
		return
	}

	class, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrClassNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail class", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail class berhasil diambil", class)
}

// Update godoc: PUT /api/v1/classes/:id
func (h *ClassHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID class tidak valid", nil)
		return
	}

	var req services.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payload request tidak valid", err.Error())
		return
	}

	class, err := h.service.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrClassNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, services.ErrClassCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data class", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data class berhasil diperbarui", class)
}

// Delete godoc: DELETE /api/v1/classes/:id
func (h *ClassHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID class tidak valid", nil)
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrClassNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus data class", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data class berhasil dihapus", nil)
}

// GetStudents godoc: GET /api/v1/classes/:id/students
func (h *ClassHandler) GetStudents(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID class tidak valid", nil)
		return
	}

	students, err := h.service.GetStudents(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrClassNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar siswa di kelas", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar siswa dalam kelas berhasil diambil", students)
}
