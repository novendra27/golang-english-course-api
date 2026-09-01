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

// PlaceStudent godoc
// @Summary      Menempatkan siswa ke dalam kelas
// @Description  Menempatkan siswa yang pendaftarannya sudah lunas ('registered') ke dalam kelas dengan validasi kapasitas dan kesesuaian course
// @Tags         Class Placements
// @Accept       json
// @Produce      json
// @Param        request body services.CreateClassPlacementRequest true "Payload penempatan kelas"
// @Success      201  {object}  utils.APIResponse{data=models.ClassPlacement}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      409  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /class-placements [post]
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

// GetAll godoc
// @Summary      Mengambil seluruh data penempatan kelas
// @Description  Mengembalikan daftar penempatan siswa beserta relasi registration, student, course, dan class
// @Tags         Class Placements
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.ClassPlacement}
// @Failure      500  {object}  utils.APIResponse
// @Router       /class-placements [get]
func (h *ClassPlacementHandler) GetAll(c *gin.Context) {
	placements, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar penempatan kelas", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar penempatan kelas berhasil diambil", placements)
}

// GetByID godoc
// @Summary      Mengambil detail penempatan kelas
// @Description  Mengembalikan data spesifik penempatan kelas berdasarkan ID
// @Tags         Class Placements
// @Produce      json
// @Param        id   path      int  true  "Placement ID"
// @Success      200  {object}  utils.APIResponse{data=models.ClassPlacement}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /class-placements/{id} [get]
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
