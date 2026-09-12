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
// @Summary      Register a new student
// @Description  Creates a new student profile with unique email validation
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        request body services.CreateStudentRequest true "Student registration payload"
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
			utils.ErrorResponse(c, http.StatusConflict, utils.Translate(c, "student.email_conflict", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "student.create_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, utils.Translate(c, "student.created", nil), student)
}

// GetAll godoc
// @Summary      Get all students
// @Description  Returns a list of all registered students
// @Tags         Students
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Student}
// @Failure      500  {object}  utils.APIResponse
// @Router       /students [get]
func (h *StudentHandler) GetAll(c *gin.Context) {
	students, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "student.fetch_all_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "student.fetched_all", nil), students)
}

// GetByID godoc
// @Summary      Get student profile detail
// @Description  Returns detail of a student by ID
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "student.invalid_id", nil), nil)
		return
	}

	student, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrStudentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "student.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "student.fetch_detail_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "student.fetched_detail", nil), student)
}

// Update godoc
// @Summary      Update student profile
// @Description  Updates name, email, or phone number of a student
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "Student ID"
// @Param        request  body      services.UpdateStudentRequest true  "Update student payload"
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "student.invalid_id", nil), nil)
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
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "student.not_found", nil), nil)
			return
		}
		if errors.Is(err, services.ErrStudentEmailConflict) {
			utils.ErrorResponse(c, http.StatusConflict, utils.Translate(c, "student.email_conflict", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "student.update_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "student.updated", nil), student)
}

// Delete godoc
// @Summary      Delete student data
// @Description  Deletes student data by ID
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "student.invalid_id", nil), nil)
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrStudentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "student.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "student.delete_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "student.deleted", nil), nil)
}
