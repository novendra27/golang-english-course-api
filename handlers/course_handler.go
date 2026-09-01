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

// Create godoc
// @Summary      Membuat katalog course baru
// @Description  Menambahkan paket kursus baru beserta harga dan durasi
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        request body services.CreateCourseRequest true "Payload data course"
// @Success      201  {object}  utils.APIResponse{data=models.Course}
// @Failure      400  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /courses [post]
func (h *CourseHandler) Create(c *gin.Context) {
	var req services.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	course, err := h.service.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat data course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Course berhasil dibuat", course)
}

// GetAll godoc
// @Summary      Mengambil daftar seluruh course
// @Description  Mengembalikan katalog semua kursus bahasa Inggris yang tersedia
// @Tags         Courses
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Course}
// @Failure      500  {object}  utils.APIResponse
// @Router       /courses [get]
func (h *CourseHandler) GetAll(c *gin.Context) {
	courses, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar course", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar course berhasil diambil", courses)
}

// GetByID godoc
// @Summary      Mengambil detail course beserta daftar kelasnya
// @Description  Mengembalikan detail spesifik course beserta kelas yang dibuka di bawahnya
// @Tags         Courses
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  utils.APIResponse{data=models.Course}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /courses/{id} [get]
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

// Update godoc
// @Summary      Mengubah data course
// @Description  Memperbarui nama, deskripsi, harga, durasi, atau status course
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        id       path      int                         true  "Course ID"
// @Param        request  body      services.UpdateCourseRequest true  "Payload update course"
// @Success      200  {object}  utils.APIResponse{data=models.Course}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /courses/{id} [put]
func (h *CourseHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID course tidak valid", nil)
		return
	}

	var req services.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
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

// Delete godoc
// @Summary      Menghapus data course
// @Description  Menghapus course berdasarkan ID
// @Tags         Courses
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /courses/{id} [delete]
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
