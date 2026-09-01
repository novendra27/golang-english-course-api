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
// @Summary      Mengambil seluruh tagihan pembayaran
// @Description  Mengembalikan daftar semua data transaksi payment
// @Tags         Payments
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]models.Payment}
// @Failure      500  {object}  utils.APIResponse
// @Router       /payments [get]
func (h *PaymentHandler) GetAll(c *gin.Context) {
	payments, err := h.service.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar tagihan pembayaran", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar tagihan pembayaran berhasil diambil", payments)
}

// GetByID godoc
// @Summary      Mengambil detail tagihan pembayaran
// @Description  Mengembalikan informasi tagihan payment berdasarkan ID
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

// Pay godoc
// @Summary      Memproses simulasi pembayaran tagihan
// @Description  Melakukan pembayaran tagihan dan secara atomik mengubah status registrasi menjadi 'registered'
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id       path      int                           true  "Payment ID"
// @Param        request  body      services.ProcessPaymentRequest true  "Payload pembayaran"
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
		utils.ErrorResponse(c, http.StatusBadRequest, "ID payment tidak valid", nil)
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
