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

// Create godoc
// @Summary      Membuat kelas baru di bawah course
// @Description  Membuka kelas baru dengan jadwal dan kapasitas tertentu
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        request body services.CreateClassRequest true "Payload data kelas"
// @Success      201  {object}  utils.APIResponse{data=models.Class}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /classes [post]
func (h *ClassHandler) Create(c *gin.Context) {
	var req services.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
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

// GetAll godoc
// @Summary      Mengambil daftar seluruh kelas
// @Description  Mengembalikan seluruh kelas beserta data Course-nya
// @Tags         Classes
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Class}
// @Failure      500  {object}  utils.APIResponse
// @Router       /classes [get]
func (h *ClassHandler) GetAll(c *gin.Context) {
	classes, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar class", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar class berhasil diambil", classes)
}

// GetByID godoc
// @Summary      Mengambil detail kelas
// @Description  Mengembalikan data spesifik kelas berdasarkan ID
// @Tags         Classes
// @Produce      json
// @Param        id   path      int  true  "Class ID"
// @Success      200  {object}  utils.APIResponse{data=models.Class}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /classes/{id} [get]
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

// Update godoc
// @Summary      Mengubah data kelas
// @Description  Memperbarui nama, kapasitas, jadwal, atau status kelas
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        id       path      int                        true  "Class ID"
// @Param        request  body      services.UpdateClassRequest true  "Payload update class"
// @Success      200  {object}  utils.APIResponse{data=models.Class}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /classes/{id} [put]
func (h *ClassHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID class tidak valid", nil)
		return
	}

	var req services.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
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

// Delete godoc
// @Summary      Menghapus data kelas
// @Description  Menghapus kelas berdasarkan ID
// @Tags         Classes
// @Produce      json
// @Param        id   path      int  true  "Class ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /classes/{id} [delete]
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

// GetStudents godoc
// @Summary      Mengambil daftar siswa dalam kelas
// @Description  Mengembalikan seluruh siswa yang telah ditempatkan di kelas ini
// @Tags         Classes
// @Produce      json
// @Param        id   path      int  true  "Class ID"
// @Success      200  {object}  utils.APIResponse{data=[]models.Student}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /classes/{id}/students [get]
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
