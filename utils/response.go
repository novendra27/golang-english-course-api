package utils

import (
	"errors"
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

// ValidationErrorResponse menerjemahkan error dari validator.ValidationErrors menjadi pesan error multi-bahasa per-field
func ValidationErrorResponse(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errMap := make(map[string]string)
		for _, fe := range ve {
			errMap[fe.Field()] = TranslateValidationError(c, fe)
		}
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Success: false,
			Message: Translate(c, "system.validation_failed", nil),
			Errors:  errMap,
		})
		return
	}

	// Jika bukan tipe validator.ValidationErrors (misal invalid JSON syntax)
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: Translate(c, "system.invalid_json", nil),
		Errors:  err.Error(),
	})
}

