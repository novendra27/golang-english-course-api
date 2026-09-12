package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"english-course-api/services"
	"english-course-api/utils"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	service services.RegistrationService
}

func NewRegistrationHandler(service services.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{service: service}
}

// Register godoc
// @Summary      Register to a course
// @Description  Registers a student to a course and automatically generates a pending payment invoice
// @Tags         Registrations
// @Accept       json
// @Produce      json
// @Param        request body services.CreateRegistrationRequest true "Registration payload"
// @Success      201  {object}  utils.APIResponse{data=models.Registration}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      409  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /registrations [post]
func (h *RegistrationHandler) Register(c *gin.Context) {
	var req services.CreateRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	registration, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, services.ErrRegistrationStudentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "registration.student_not_found", nil), nil)
			return
		}
		if errors.Is(err, services.ErrRegistrationCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "registration.course_not_found", nil), nil)
			return
		}
		if errors.Is(err, services.ErrRegistrationCourseInactive) {
			utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "registration.course_inactive", nil), nil)
			return
		}
		if errors.Is(err, services.ErrRegistrationAlreadyActive) {
			utils.ErrorResponse(c, http.StatusConflict, utils.Translate(c, "registration.already_active", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "registration.create_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, utils.Translate(c, "registration.created", nil), registration)
}

// GetAll godoc
// @Summary      Get all registrations
// @Description  Returns all registration records along with student, course, and payment relations
// @Tags         Registrations
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Registration}
// @Failure      500  {object}  utils.APIResponse
// @Router       /registrations [get]
func (h *RegistrationHandler) GetAll(c *gin.Context) {
	registrations, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "registration.fetch_all_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "registration.fetched_all", nil), registrations)
}

// GetByID godoc
// @Summary      Get registration detail
// @Description  Returns detail of a registration by ID
// @Tags         Registrations
// @Produce      json
// @Param        id   path      int  true  "Registration ID"
// @Success      200  {object}  utils.APIResponse{data=models.Registration}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /registrations/{id} [get]
func (h *RegistrationHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "registration.invalid_id", nil), nil)
		return
	}

	reg, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRegistrationNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "registration.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "registration.fetch_detail_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "registration.fetched_detail", nil), reg)
}

// GetByStudentID godoc
// @Summary      Get student registration history
// @Description  Returns all registrations made by a specific student
// @Tags         Students
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  utils.APIResponse{data=[]models.Registration}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /students/{id}/registrations [get]
func (h *RegistrationHandler) GetByStudentID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "registration.invalid_student_id", nil), nil)
		return
	}

	regs, err := h.service.GetByStudentID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRegistrationStudentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "registration.student_not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "registration.student_regs_fetch_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "registration.student_regs_fetched", nil), regs)
}

// GetByCourseID godoc
// @Summary      Get registrations for a course
// @Description  Returns all students registered for a specific course
// @Tags         Courses
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  utils.APIResponse{data=[]models.Registration}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /courses/{id}/registrations [get]
func (h *RegistrationHandler) GetByCourseID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "registration.invalid_course_id", nil), nil)
		return
	}

	regs, err := h.service.GetByCourseID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRegistrationCourseNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "registration.course_not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "registration.course_regs_fetch_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "registration.course_regs_fetched", nil), regs)
}

// CancelRegistration godoc
// @Summary      Cancel a registration
// @Description  Updates registration status to cancelled if not yet completed
// @Tags         Registrations
// @Produce      json
// @Param        id   path      int  true  "Registration ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /registrations/{id}/cancel [put]
func (h *RegistrationHandler) CancelRegistration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "registration.invalid_id", nil), nil)
		return
	}

	err = h.service.CancelRegistration(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRegistrationNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "registration.not_found", nil), nil)
			return
		}
		if errors.Is(err, services.ErrRegistrationCannotBeCanceled) {
			utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "registration.cannot_cancel", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "registration.cancel_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "registration.cancelled", nil), nil)
}
