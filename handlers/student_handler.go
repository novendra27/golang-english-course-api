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

// Create godoc
// @Summary      Mendaftarkan student baru
// @Description  Membuat data peserta kursus baru dengan validasi email unik
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        request body services.CreateStudentRequest true "Payload pendaftaran student"
// @Success      201  {object}  utils.APIResponse{data=models.Student}
// @Failure      400  {object}  utils.APIResponse
// @Failure      409  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /students [post]
func (h *StudentHandler) Create(c *gin.Context) {
	var req services.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
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

// GetAll godoc
// @Summary      Mengambil daftar seluruh student
// @Description  Mengembalikan daftar semua peserta kursus yang terdaftar
// @Tags         Students
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Student}
// @Failure      500  {object}  utils.APIResponse
// @Router       /students [get]
func (h *StudentHandler) GetAll(c *gin.Context) {
	students, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar student", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar student berhasil diambil", students)
}

// GetByID godoc
// @Summary      Mengambil detail profil student
// @Description  Mengembalikan data detail student berdasarkan ID
// @Tags         Students
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  utils.APIResponse{data=models.Student}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /students/{id} [get]
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

// Update godoc
// @Summary      Mengubah data profil student
// @Description  Memperbarui nama, email, atau no telepon student
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "Student ID"
// @Param        request  body      services.UpdateStudentRequest true  "Payload update student"
// @Success      200  {object}  utils.APIResponse{data=models.Student}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      409  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /students/{id} [put]
func (h *StudentHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID student tidak valid", nil)
		return
	}

	var req services.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
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

// Delete godoc
// @Summary      Menghapus data student
// @Description  Menghapus data student berdasarkan ID
// @Tags         Students
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /students/{id} [delete]
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
