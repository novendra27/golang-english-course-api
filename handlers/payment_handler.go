package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"english-course-api/services"
	"english-course-api/utils"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service services.PaymentService
}

func NewPaymentHandler(service services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// GetAll godoc
// @Summary      Get all payment invoices
// @Description  Returns a list of all payment transaction records
// @Tags         Payments
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Payment}
// @Failure      500  {object}  utils.APIResponse
// @Router       /payments [get]
func (h *PaymentHandler) GetAll(c *gin.Context) {
	payments, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "payment.fetch_all_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "payment.fetched_all", nil), payments)
}

// GetByID godoc
// @Summary      Get payment invoice detail
// @Description  Returns information of a specific payment invoice by ID
// @Tags         Payments
// @Produce      json
// @Param        id   path      int  true  "Payment ID"
// @Success      200  {object}  utils.APIResponse{data=models.Payment}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /payments/{id} [get]
func (h *PaymentHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "payment.invalid_id", nil), nil)
		return
	}

	payment, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrPaymentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "payment.not_found", nil), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "payment.fetch_detail_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "payment.fetched_detail", nil), payment)
}

// Pay godoc
// @Summary      Process simulated payment
// @Description  Processes payment settlement and atomically transitions registration status to 'registered'
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id       path      int                           true  "Payment ID"
// @Param        request  body      services.ProcessPaymentRequest true  "Payment processing payload"
// @Success      200  {object}  utils.APIResponse{data=models.Payment}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /payments/{id}/pay [post]
func (h *PaymentHandler) Pay(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "payment.invalid_id", nil), nil)
		return
	}

	var req services.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	payment, err := h.service.Pay(uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrPaymentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, utils.Translate(c, "payment.not_found", nil), nil)
			return
		}
		if errors.Is(err, services.ErrPaymentAlreadyPaid) {
			utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "payment.already_paid", nil), err.Error())
			return
		}
		if errors.Is(err, services.ErrPaymentInvalidStatus) {
			utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "payment.invalid_status", nil), err.Error())
			return
		}
		if errors.Is(err, services.ErrPaymentAmountInvalid) {
			utils.ErrorResponse(c, http.StatusBadRequest, utils.Translate(c, "payment.amount_invalid", nil), err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.Translate(c, "payment.pay_failed", nil), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.Translate(c, "payment.paid_success", nil), payment)
}
