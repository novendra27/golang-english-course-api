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
// @Summary      Create a new course
// @Description  Creates a new course catalog package with price and duration
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        request body services.CreateCourseRequest true "Course creation payload"
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
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "course.create_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, utils.Translate(c, "course.created", nil), course)
}

// GetAll godoc
// @Summary      Get all courses
// @Description  Returns a catalog of all available English courses
// @Tags         Courses
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Course}
// @Failure      500  {object}  utils.APIResponse
// @Router       /courses [get]
func (h *CourseHandler) GetAll(c *gin.Context) {
	courses, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "course.fetch_all_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "course.fetched_all", nil), courses)
}

// GetByID godoc
// @Summary      Get course detail with classes
// @Description  Returns detail of a specific course along with open classes under it
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "course.invalid_id", nil), nil)
		return
	}

	course, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "course.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "course.fetch_detail_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "course.fetched_detail", nil), course)
}

// Update godoc
// @Summary      Update course data
// @Description  Updates name, description, price, duration, or status of a course
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        id       path      int                         true  "Course ID"
// @Param        request  body      services.UpdateCourseRequest true  "Update course payload"
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "course.invalid_id", nil), nil)
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
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "course.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "course.update_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "course.updated", nil), course)
}

// Delete godoc
// @Summary      Delete course data
// @Description  Deletes a course by ID
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "course.invalid_id", nil), nil)
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "course.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "course.delete_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "course.deleted", nil), nil)
}
