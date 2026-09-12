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
// @Summary      Place student into class
// @Description  Places a student with paid registration ('registered') into a class with capacity and course matching validation
// @Tags         Class Placements
// @Accept       json
// @Produce      json
// @Param        request body services.CreateClassPlacementRequest true "Class placement payload"
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
		if errors.Is(err, services.ErrPlacementRegNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "placement.reg_not_found", nil), nil)
			return
		}
		if errors.Is(err, services.ErrPlacementClassNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "placement.class_not_found", nil), nil)
			return
		}
		if errors.Is(err, services.ErrPlacementPaymentRequired) {
			utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "placement.payment_required", nil), nil)
			return
		}
		if errors.Is(err, services.ErrPlacementCourseMismatch) {
			utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "placement.course_mismatch", nil), nil)
			return
		}
		if errors.Is(err, services.ErrPlacementClassClosed) {
			utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "placement.class_closed", nil), nil)
			return
		}
		if errors.Is(err, services.ErrPlacementAlreadyAssigned) {
			utils.ErrorResponse(c, http.StatusConflict, utils.Translate(c, "placement.already_assigned", nil), nil)
			return
		}
		if errors.Is(err, services.ErrPlacementClassFull) {
			utils.ErrorResponse(c, http.StatusConflict, utils.Translate(c, "placement.class_full", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "placement.place_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, utils.Translate(c, "placement.placed_success", nil), placement)
}

// GetAll godoc
// @Summary      Get all class placements
// @Description  Returns all class placement records along with registration, student, course, and class relations
// @Tags         Class Placements
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.ClassPlacement}
// @Failure      500  {object}  utils.APIResponse
// @Router       /class-placements [get]
func (h *ClassPlacementHandler) GetAll(c *gin.Context) {
	placements, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "placement.fetch_all_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "placement.fetched_all", nil), placements)
}

// GetByID godoc
// @Summary      Get class placement detail
// @Description  Returns detail of a specific class placement by ID
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
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "placement.invalid_id", nil), nil)
		return
	}

	placement, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrPlacementNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "placement.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "placement.fetch_detail_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "placement.fetched_detail", nil), placement)
}
