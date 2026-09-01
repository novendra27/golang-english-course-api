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

// GetAll godoc: GET /api/v1/payments
func (h *PaymentHandler) GetAll(c *gin.Context) {
	payments, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar tagihan pembayaran", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar tagihan pembayaran berhasil diambil", payments)
}

// GetByID godoc: GET /api/v1/payments/:id
func (h *PaymentHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID payment tidak valid", nil)
		return
	}

	payment, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrPaymentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail pembayaran", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail pembayaran berhasil diambil", payment)
}

// Pay godoc: POST /api/v1/payments/:id/pay
func (h *PaymentHandler) Pay(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID payment tidak valid", nil)
		return
	}

	var req services.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payload request tidak valid", err.Error())
		return
	}

	payment, err := h.service.Pay(uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrPaymentNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, services.ErrPaymentAlreadyPaid) || errors.Is(err, services.ErrPaymentInvalidStatus) || errors.Is(err, services.ErrPaymentAmountInvalid) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memproses pembayaran", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pembayaran berhasil diproses! Status registrasi aktif 🎉", payment)
}
