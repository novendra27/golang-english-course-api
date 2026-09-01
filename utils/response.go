package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APIResponse merepresentasikan format JSON response standar aplikasi
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse mengirimkan format response sukses dengan status code dan data
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirimkan format response error terstandarisasi
func ErrorResponse(c *gin.Context, statusCode int, message string, errDetails interface{}) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errDetails,
	})
}

// ValidationErrorResponse menerjemahkan error dari validator.ValidationErrors menjadi pesan error yang rapi per-field
func ValidationErrorResponse(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errMap := make(map[string]string)
		for _, fe := range ve {
			errMap[fe.Field()] = formatFieldError(fe)
		}
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Success: false,
			Message: "Validasi request gagal",
			Errors:  errMap,
		})
		return
	}

	// Jika bukan tipe validator.ValidationErrors (misal invalid JSON syntax)
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: "Format payload JSON tidak valid",
		Errors:  err.Error(),
	})
}

// formatFieldError menerjemahkan tag validasi validator/v10 menjadi pesan yang mudah dibaca
func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' wajib diisi", fe.Field())
	case "email":
		return fmt.Sprintf("Field '%s' harus berupa alamat email yang valid", fe.Field())
	case "gt":
		return fmt.Sprintf("Field '%s' harus lebih besar dari %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("Field '%s' harus lebih besar atau sama dengan %s", fe.Field(), fe.Param())
	case "lt":
		return fmt.Sprintf("Field '%s' harus lebih kecil dari %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("Field '%s' harus lebih kecil atau sama dengan %s", fe.Field(), fe.Param())
	case "min":
		return fmt.Sprintf("Field '%s' minimal memiliki panjang/nilai %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("Field '%s' maksimal memiliki panjang/nilai %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("Field '%s' harus salah satu dari [%s]", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("Field '%s' tidak valid (%s)", fe.Field(), fe.Tag())
	}
}
