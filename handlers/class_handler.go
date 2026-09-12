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
// @Summary      Create a new class
// @Description  Opens a new class under a course with schedule and capacity
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        request body services.CreateClassRequest true "Class creation payload"
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
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "class.course_not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "class.create_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, utils.Translate(c, "class.created", nil), class)
}

// GetAll godoc
// @Summary      Get all classes
// @Description  Returns all classes along with their Course data
// @Tags         Classes
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Class}
// @Failure      500  {object}  utils.APIResponse
// @Router       /classes [get]
func (h *ClassHandler) GetAll(c *gin.Context) {
	classes, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "class.fetch_all_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "class.fetched_all", nil), classes)
}

// GetByID godoc
// @Summary      Get class detail
// @Description  Returns detail of a specific class by ID
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "class.invalid_id", nil), nil)
		return
	}

	class, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrClassNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "class.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "class.fetch_detail_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "class.fetched_detail", nil), class)
}

// Update godoc
// @Summary      Update class data
// @Description  Updates name, capacity, schedule, or status of a class
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        id       path      int                        true  "Class ID"
// @Param        request  body      services.UpdateClassRequest true  "Update class payload"
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "class.invalid_id", nil), nil)
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
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "class.not_found", nil), nil)
			return
		}
		if errors.Is(err, services.ErrClassCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "class.course_not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "class.update_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "class.updated", nil), class)
}

// Delete godoc
// @Summary      Delete class data
// @Description  Deletes a class by ID
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "class.invalid_id", nil), nil)
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrClassNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "class.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "class.delete_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "class.deleted", nil), nil)
}

// GetStudents godoc
// @Summary      Get students in a class
// @Description  Returns all students who have been placed into this class
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "class.invalid_id", nil), nil)
		return
	}

	students, err := h.service.GetStudents(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrClassNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "class.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "class.students_fetch_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "class.students_fetched", nil), students)
}
