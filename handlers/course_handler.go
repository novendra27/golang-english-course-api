package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"english-course-api/services"
	"english-course-api/utils"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service services.CourseService
}

func NewCourseHandler(service services.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

// Create godoc: POST /api/v1/courses
func (h *CourseHandler) Create(c *gin.Context) {
	var req services.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payload request tidak valid", err.Error())
		return
	}

	course, err := h.service.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat data course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Course berhasil dibuat", course)
}

// GetAll godoc: GET /api/v1/courses
func (h *CourseHandler) GetAll(c *gin.Context) {
	courses, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar course berhasil diambil", courses)
}

// GetByID godoc: GET /api/v1/courses/:id
func (h *CourseHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID course tidak valid", nil)
		return
	}

	course, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail course berhasil diambil", course)
}

// Update godoc: PUT /api/v1/courses/:id
func (h *CourseHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID course tidak valid", nil)
		return
	}

	var req services.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payload request tidak valid", err.Error())
		return
	}

	course, err := h.service.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data course berhasil diperbarui", course)
}

// Delete godoc: DELETE /api/v1/courses/:id
func (h *CourseHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID course tidak valid", nil)
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus data course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data course berhasil dihapus", nil)
}
