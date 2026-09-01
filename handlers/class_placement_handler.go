package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"english-course-api/services"
	"english-course-api/utils"

	"github.com/gin-gonic/gin"
)

type ClassPlacementHandler struct {
	service services.ClassPlacementService
}

func NewClassPlacementHandler(service services.ClassPlacementService) *ClassPlacementHandler {
	return &ClassPlacementHandler{service: service}
}

// PlaceStudent godoc: POST /api/v1/class-placements
func (h *ClassPlacementHandler) PlaceStudent(c *gin.Context) {
	var req services.CreateClassPlacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	placement, err := h.service.PlaceStudent(req)
	if err != nil {
		if errors.Is(err, services.ErrPlacementRegNotFound) || errors.Is(err, services.ErrPlacementClassNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, services.ErrPlacementPaymentRequired) ||
			errors.Is(err, services.ErrPlacementCourseMismatch) ||
			errors.Is(err, services.ErrPlacementClassClosed) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		if errors.Is(err, services.ErrPlacementAlreadyAssigned) || errors.Is(err, services.ErrPlacementClassFull) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menempatkan siswa ke kelas", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Student berhasil ditempatkan ke dalam kelas 🎉", placement)
}

// GetAll godoc: GET /api/v1/class-placements
func (h *ClassPlacementHandler) GetAll(c *gin.Context) {
	placements, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar penempatan kelas", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar penempatan kelas berhasil diambil", placements)
}

// GetByID godoc: GET /api/v1/class-placements/:id
func (h *ClassPlacementHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID penempatan tidak valid", nil)
		return
	}

	placement, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrPlacementNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail penempatan", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail penempatan kelas berhasil diambil", placement)
}
